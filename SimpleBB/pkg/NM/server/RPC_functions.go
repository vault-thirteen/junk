package server

import (
	"fmt"
	am "github.com/vault-thirteen/SimpleBB/pkg/ACM/rpc"
	rpc2 "github.com/vault-thirteen/SimpleBB/pkg/NM/rpc"
	"github.com/vault-thirteen/SimpleBB/pkg/common/models/UidList"
	base2 "github.com/vault-thirteen/SimpleBB/pkg/common/models/base"
	"github.com/vault-thirteen/SimpleBB/pkg/common/models/complex/EnumValue"
	"github.com/vault-thirteen/SimpleBB/pkg/common/models/complex/IncidentType"
	"github.com/vault-thirteen/SimpleBB/pkg/common/models/complex/Resource"
	cm "github.com/vault-thirteen/SimpleBB/pkg/common/models/complex/SystemEvent"
	set "github.com/vault-thirteen/SimpleBB/pkg/common/models/complex/SystemEventType"
	rpc3 "github.com/vault-thirteen/SimpleBB/pkg/common/models/rpc"
	c "github.com/vault-thirteen/SimpleBB/pkg/common/models/server"
	"sync"

	jrm1 "github.com/vault-thirteen/JSON-RPC-M1"
	nm "github.com/vault-thirteen/SimpleBB/pkg/NM/models"
)

// RPC functions.

// Notification.

// addNotification creates a new notification.
// This method is used to send notifications by administrators.
func (srv *Server) addNotification(p *rpc2.AddNotificationParams) (result *rpc2.AddNotificationResult, re *jrm1.RpcError) {
	// Check parameters.
	if p.UserId == 0 {
		return nil, jrm1.NewRpcErrorByUser(RpcErrorCode_UserIdIsNotSet, RpcErrorMsg_UserIdIsNotSet, nil)
	}

	if len(p.Text) == 0 {
		return nil, jrm1.NewRpcErrorByUser(RpcErrorCode_TextIsNotSet, RpcErrorMsg_TextIsNotSet, nil)
	}

	var userRoles *am.GetSelfRolesResult
	userRoles, re = srv.mustBeAnAuthToken(p.Auth)
	if re != nil {
		return nil, re
	}

	// Check permissions.
	if !userRoles.User.GetUserParameters().GetRoles().IsAdministrator {
		return nil, jrm1.NewRpcErrorByUser(c.RpcErrorCode_Permission, c.RpcErrorMsg_Permission, nil)
	}

	srv.dbo.LockForWriting()
	defer srv.dbo.UnlockAfterWriting()

	insertedNotificationId, err := srv.dbo.InsertNewNotification(p.UserId, p.Text)
	if err != nil {
		return nil, srv.databaseError(err)
	}

	result = &rpc2.AddNotificationResult{
		NotificationId: insertedNotificationId,
	}

	return result, nil
}

// addNotificationS creates a new notification.
// This method is used to send notifications by the system.
func (srv *Server) addNotificationS(p *rpc2.AddNotificationSParams) (result *rpc2.AddNotificationSResult, re *jrm1.RpcError) {
	// Check parameters.
	if p.UserId == 0 {
		return nil, jrm1.NewRpcErrorByUser(RpcErrorCode_UserIdIsNotSet, RpcErrorMsg_UserIdIsNotSet, nil)
	}

	if len(p.Text) == 0 {
		return nil, jrm1.NewRpcErrorByUser(RpcErrorCode_TextIsNotSet, RpcErrorMsg_TextIsNotSet, nil)
	}

	re = srv.mustBeNoAuth(p.Auth)
	if re != nil {
		return nil, re
	}

	// Check the DKey.
	if !srv.dKeyI.CheckString(p.DKey.ToString()) {
		srv.incidentManager.ReportIncident(ev.NewEnumValue(it.IncidentType_WrongDKey), "", nil)
		return nil, jrm1.NewRpcErrorByUser(c.RpcErrorCode_Permission, c.RpcErrorMsg_Permission, nil)
	}

	srv.dbo.LockForWriting()
	defer srv.dbo.UnlockAfterWriting()

	insertedNotificationId, err := srv.dbo.InsertNewNotification(p.UserId, p.Text)
	if err != nil {
		return nil, srv.databaseError(err)
	}

	result = &rpc2.AddNotificationSResult{
		NotificationId: insertedNotificationId,
	}

	return result, nil
}

// sendNotificationIfPossibleS tries to send a notification to a user when it
// is possible. This method is used by the system.
func (srv *Server) sendNotificationIfPossibleS(p *rpc2.SendNotificationIfPossibleSParams) (result *rpc2.SendNotificationIfPossibleSResult, re *jrm1.RpcError) {
	// Check parameters.
	if p.UserId == 0 {
		return nil, jrm1.NewRpcErrorByUser(RpcErrorCode_UserIdIsNotSet, RpcErrorMsg_UserIdIsNotSet, nil)
	}

	if len(p.Text) == 0 {
		return nil, jrm1.NewRpcErrorByUser(RpcErrorCode_TextIsNotSet, RpcErrorMsg_TextIsNotSet, nil)
	}

	re = srv.mustBeNoAuth(p.Auth)
	if re != nil {
		return nil, re
	}

	// Check the DKey.
	if !srv.dKeyI.CheckString(p.DKey.ToString()) {
		srv.incidentManager.ReportIncident(ev.NewEnumValue(it.IncidentType_WrongDKey), "", nil)
		return nil, jrm1.NewRpcErrorByUser(c.RpcErrorCode_Permission, c.RpcErrorMsg_Permission, nil)
	}

	return srv.sendNotificationIfPossibleH(p.UserId, p.Text)
}

// getNotification reads a notification.
func (srv *Server) getNotification(p *rpc2.GetNotificationParams) (result *rpc2.GetNotificationResult, re *jrm1.RpcError) {
	// Check parameters.
	if p.NotificationId == 0 {
		return nil, jrm1.NewRpcErrorByUser(RpcErrorCode_NotificationIdIsNotSet, RpcErrorMsg_NotificationIdIsNotSet, nil)
	}

	var userRoles *am.GetSelfRolesResult
	userRoles, re = srv.mustBeAnAuthToken(p.Auth)
	if re != nil {
		return nil, re
	}

	// Check permissions.
	if !userRoles.User.GetUserParameters().GetRoles().CanLogIn {
		return nil, jrm1.NewRpcErrorByUser(c.RpcErrorCode_Permission, c.RpcErrorMsg_Permission, nil)
	}

	srv.dbo.LockForReading()
	defer srv.dbo.UnlockAfterReading()

	// Read the notification.
	notification, err := srv.dbo.GetNotificationById(p.NotificationId)
	if err != nil {
		return nil, srv.databaseError(err)
	}
	if notification == nil {
		return nil, jrm1.NewRpcErrorByUser(RpcErrorCode_NotificationIsNotFound, RpcErrorMsg_NotificationIsNotFound, nil)
	}

	// Check the recipient.
	if notification.UserId != userRoles.User.GetUserParameters().GetId() {
		srv.incidentManager.ReportIncident(ev.NewEnumValue(it.IncidentType_ReadingNotificationOfOtherUsers), userRoles.User.GetUserParameters().GetEmail(), p.Auth.UserIPAB)

		return nil, jrm1.NewRpcErrorByUser(c.RpcErrorCode_Permission, c.RpcErrorMsg_Permission, nil)
	}

	// All clear.
	result = &rpc2.GetNotificationResult{
		Notification: notification,
	}

	return result, nil
}

// getNotifications reads all notifications for a user.
func (srv *Server) getNotifications(p *rpc2.GetNotificationsParams) (result *rpc2.GetNotificationsResult, re *jrm1.RpcError) {
	var userRoles *am.GetSelfRolesResult
	userRoles, re = srv.mustBeAnAuthToken(p.Auth)
	if re != nil {
		return nil, re
	}

	// Check permissions.
	if !userRoles.User.GetUserParameters().GetRoles().CanLogIn {
		return nil, jrm1.NewRpcErrorByUser(c.RpcErrorCode_Permission, c.RpcErrorMsg_Permission, nil)
	}

	srv.dbo.LockForReading()
	defer srv.dbo.UnlockAfterReading()

	// Get notifications.
	notifications, err := srv.dbo.GetAllNotificationsByUserId(userRoles.User.GetUserParameters().GetId())
	if err != nil {
		return nil, srv.databaseError(err)
	}

	result = &rpc2.GetNotificationsResult{
		Notifications: notifications,
	}

	result.NotificationIds, err = ul.NewFromArray(nm.ListNotificationIds(notifications))
	if err != nil {
		srv.logError(err)
		return nil, jrm1.NewRpcErrorByUser(c.RpcErrorCode_UidList, fmt.Sprintf(c.RpcErrorMsgF_UidList, err.Error()), nil)
	}

	return result, nil
}

// getNotificationsOnPage reads notifications for a user on the selected page.
func (srv *Server) getNotificationsOnPage(p *rpc2.GetNotificationsOnPageParams) (result *rpc2.GetNotificationsOnPageResult, re *jrm1.RpcError) {
	// Check parameters.
	if p.Page == 0 {
		return nil, jrm1.NewRpcErrorByUser(RpcErrorCode_PageIsNotSet, RpcErrorMsg_PageIsNotSet, nil)
	}

	var userRoles *am.GetSelfRolesResult
	userRoles, re = srv.mustBeAnAuthToken(p.Auth)
	if re != nil {
		return nil, re
	}

	// Check permissions.
	if !userRoles.User.GetUserParameters().GetRoles().CanLogIn {
		return nil, jrm1.NewRpcErrorByUser(c.RpcErrorCode_Permission, c.RpcErrorMsg_Permission, nil)
	}

	srv.dbo.LockForReading()
	defer srv.dbo.UnlockAfterReading()

	// Get notifications on page.
	notifications, err := srv.dbo.GetNotificationsByUserIdOnPage(userRoles.User.GetUserParameters().GetId(), p.Page, srv.settings.SystemSettings.PageSize)
	if err != nil {
		return nil, srv.databaseError(err)
	}

	// Count all notifications.
	var allNotificationsCount base2.Count
	allNotificationsCount, err = srv.dbo.CountAllNotificationsByUserId(userRoles.User.GetUserParameters().GetId())
	if err != nil {
		return nil, srv.databaseError(err)
	}

	var notificationIds *ul.UidList
	notificationIds, err = ul.NewFromArray(nm.ListNotificationIds(notifications))
	if err != nil {
		srv.logError(err)
		return nil, jrm1.NewRpcErrorByUser(c.RpcErrorCode_UidList, fmt.Sprintf(c.RpcErrorMsgF_UidList, err.Error()), nil)
	}

	result = &rpc2.GetNotificationsOnPageResult{
		NotificationsOnPage: &nm.NotificationsOnPage{
			NotificationIds: notificationIds,
			Notifications:   notifications,
			PageData: &rpc3.PageData{
				PageNumber:  p.Page,
				TotalPages:  base2.CalculateTotalPages(allNotificationsCount, srv.settings.SystemSettings.PageSize),
				PageSize:    srv.settings.SystemSettings.PageSize,
				ItemsOnPage: base2.Count(len(notifications)),
				TotalItems:  allNotificationsCount,
			},
		},
	}

	return result, nil
}

// getUnreadNotifications reads all unread notifications for a user.
func (srv *Server) getUnreadNotifications(p *rpc2.GetUnreadNotificationsParams) (result *rpc2.GetUnreadNotificationsResult, re *jrm1.RpcError) {
	var userRoles *am.GetSelfRolesResult
	userRoles, re = srv.mustBeAnAuthToken(p.Auth)
	if re != nil {
		return nil, re
	}

	// Check permissions.
	if !userRoles.User.GetUserParameters().GetRoles().CanLogIn {
		return nil, jrm1.NewRpcErrorByUser(c.RpcErrorCode_Permission, c.RpcErrorMsg_Permission, nil)
	}

	srv.dbo.LockForReading()
	defer srv.dbo.UnlockAfterReading()

	// Get notifications.
	notifications, err := srv.dbo.GetUnreadNotifications(userRoles.User.GetUserParameters().GetId())
	if err != nil {
		return nil, srv.databaseError(err)
	}

	result = &rpc2.GetUnreadNotificationsResult{
		Notifications: notifications,
	}

	result.NotificationIds, err = ul.NewFromArray(nm.ListNotificationIds(notifications))
	if err != nil {
		srv.logError(err)
		return nil, jrm1.NewRpcErrorByUser(c.RpcErrorCode_UidList, fmt.Sprintf(c.RpcErrorMsgF_UidList, err.Error()), nil)
	}

	return result, nil
}

// countUnreadNotifications counts unread notifications for a user.
func (srv *Server) countUnreadNotifications(p *rpc2.CountUnreadNotificationsParams) (result *rpc2.CountUnreadNotificationsResult, re *jrm1.RpcError) {
	var userRoles *am.GetSelfRolesResult
	userRoles, re = srv.mustBeAnAuthToken(p.Auth)
	if re != nil {
		return nil, re
	}

	// Check permissions.
	if !userRoles.User.GetUserParameters().GetRoles().CanLogIn {
		return nil, jrm1.NewRpcErrorByUser(c.RpcErrorCode_Permission, c.RpcErrorMsg_Permission, nil)
	}

	srv.dbo.LockForReading()
	defer srv.dbo.UnlockAfterReading()

	// Count unread notifications.
	n, err := srv.dbo.CountUnreadNotificationsByUserId(userRoles.User.GetUserParameters().GetId())
	if err != nil {
		return nil, srv.databaseError(err)
	}

	result = &rpc2.CountUnreadNotificationsResult{
		UNC: n,
	}

	return result, nil
}

// markNotificationAsRead marks a notification as read by its recipient.
func (srv *Server) markNotificationAsRead(p *rpc2.MarkNotificationAsReadParams) (result *rpc2.MarkNotificationAsReadResult, re *jrm1.RpcError) {
	// Check parameters.
	if p.NotificationId == 0 {
		return nil, jrm1.NewRpcErrorByUser(RpcErrorCode_NotificationIdIsNotSet, RpcErrorMsg_NotificationIdIsNotSet, nil)
	}

	var userRoles *am.GetSelfRolesResult
	userRoles, re = srv.mustBeAnAuthToken(p.Auth)
	if re != nil {
		return nil, re
	}

	// Check permissions.
	if !userRoles.User.GetUserParameters().GetRoles().CanLogIn {
		return nil, jrm1.NewRpcErrorByUser(c.RpcErrorCode_Permission, c.RpcErrorMsg_Permission, nil)
	}

	srv.dbo.LockForWriting()
	defer srv.dbo.UnlockAfterWriting()

	// Get the notification to see its real recipient.
	notification, err := srv.dbo.GetNotificationById(p.NotificationId)
	if err != nil {
		return nil, srv.databaseError(err)
	}
	if notification == nil {
		return nil, jrm1.NewRpcErrorByUser(RpcErrorCode_NotificationIsNotFound, RpcErrorMsg_NotificationIsNotFound, nil)
	}

	// Check the recipient and status.
	if notification.UserId != userRoles.User.GetUserParameters().GetId() {
		srv.incidentManager.ReportIncident(ev.NewEnumValue(it.IncidentType_ReadingNotificationOfOtherUsers), userRoles.User.GetUserParameters().GetEmail(), p.Auth.UserIPAB)

		return nil, jrm1.NewRpcErrorByUser(c.RpcErrorCode_Permission, c.RpcErrorMsg_Permission, nil)
	}

	if notification.IsRead {
		return nil, jrm1.NewRpcErrorByUser(RpcErrorCode_NotificationIsAlreadyRead, RpcErrorMsg_NotificationIsAlreadyRead, nil)
	}

	// Make the mark.
	err = srv.dbo.MarkNotificationAsRead(p.NotificationId, userRoles.User.GetUserParameters().GetId())
	if err != nil {
		return nil, srv.databaseError(err)
	}

	result = &rpc2.MarkNotificationAsReadResult{
		Success: rpc3.Success{
			OK: true,
		},
	}
	return result, nil
}

// deleteNotification removes a notification.
func (srv *Server) deleteNotification(p *rpc2.DeleteNotificationParams) (result *rpc2.DeleteNotificationResult, re *jrm1.RpcError) {
	// Check parameters.
	if p.NotificationId == 0 {
		return nil, jrm1.NewRpcErrorByUser(RpcErrorCode_NotificationIdIsNotSet, RpcErrorMsg_NotificationIdIsNotSet, nil)
	}

	var userRoles *am.GetSelfRolesResult
	userRoles, re = srv.mustBeAnAuthToken(p.Auth)
	if re != nil {
		return nil, re
	}

	// Check permissions.
	if !userRoles.User.GetUserParameters().GetRoles().IsAdministrator {
		return nil, jrm1.NewRpcErrorByUser(c.RpcErrorCode_Permission, c.RpcErrorMsg_Permission, nil)
	}

	srv.dbo.LockForWriting()
	defer srv.dbo.UnlockAfterWriting()

	// Delete the notification.
	err := srv.dbo.DeleteNotificationById(p.NotificationId)
	if err != nil {
		return nil, srv.databaseError(err)
	}

	result = &rpc2.DeleteNotificationResult{
		Success: rpc3.Success{
			OK: true,
		},
	}
	return result, nil
}

// Resource.

// addResource creates a new resource.
func (srv *Server) addResource(p *rpc2.AddResourceParams) (result *rpc2.AddResourceResult, re *jrm1.RpcError) {
	var userRoles *am.GetSelfRolesResult
	userRoles, re = srv.mustBeAnAuthToken(p.Auth)
	if re != nil {
		return nil, re
	}

	// Check permissions.
	if !userRoles.User.GetUserParameters().GetRoles().IsAdministrator {
		return nil, jrm1.NewRpcErrorByUser(c.RpcErrorCode_Permission, c.RpcErrorMsg_Permission, nil)
	}

	resource := res.NewResourceFromValue(p.Resource)
	if resource == nil {
		return nil, jrm1.NewRpcErrorByUser(RpcErrorCode_ResourceIsNotValid, RpcErrorMsg_ResourceIsNotValid, nil)
	}

	srv.dbo.LockForWriting()
	defer srv.dbo.UnlockAfterWriting()

	insertedNotificationId, err := srv.dbo.AddResource(resource)
	if err != nil {
		return nil, srv.databaseError(err)
	}

	result = &rpc2.AddResourceResult{
		ResourceId: insertedNotificationId,
	}

	return result, nil
}

// getResource reads a resource.
func (srv *Server) getResource(p *rpc2.GetResourceParams) (result *rpc2.GetResourceResult, re *jrm1.RpcError) {
	// Check parameters.
	if p.ResourceId == 0 {
		return nil, jrm1.NewRpcErrorByUser(RpcErrorCode_ResourceIdIsNotSet, RpcErrorMsg_ResourceIdIsNotSet, nil)
	}

	var userRoles *am.GetSelfRolesResult
	userRoles, re = srv.mustBeAnAuthToken(p.Auth)
	if re != nil {
		return nil, re
	}

	// Check permissions.
	if !userRoles.User.GetUserParameters().GetRoles().CanLogIn {
		return nil, jrm1.NewRpcErrorByUser(c.RpcErrorCode_Permission, c.RpcErrorMsg_Permission, nil)
	}

	srv.dbo.LockForReading()
	defer srv.dbo.UnlockAfterReading()

	// Read the resource.
	resource, err := srv.dbo.GetResourceById(p.ResourceId)
	if err != nil {
		return nil, srv.databaseError(err)
	}
	if resource == nil {
		return nil, jrm1.NewRpcErrorByUser(RpcErrorCode_ResourceIsNotFound, RpcErrorMsg_ResourceIsNotFound, nil)
	}

	result = &rpc2.GetResourceResult{
		Resource: resource,
	}

	return result, nil
}

// getResource reads the value of a resource.
func (srv *Server) getResourceValue(p *rpc2.GetResourceValueParams) (result *rpc2.GetResourceValueResult, re *jrm1.RpcError) {
	// Check parameters.
	if p.ResourceId == 0 {
		return nil, jrm1.NewRpcErrorByUser(RpcErrorCode_ResourceIdIsNotSet, RpcErrorMsg_ResourceIdIsNotSet, nil)
	}

	var userRoles *am.GetSelfRolesResult
	userRoles, re = srv.mustBeAnAuthToken(p.Auth)
	if re != nil {
		return nil, re
	}

	// Check permissions.
	if !userRoles.User.GetUserParameters().GetRoles().CanLogIn {
		return nil, jrm1.NewRpcErrorByUser(c.RpcErrorCode_Permission, c.RpcErrorMsg_Permission, nil)
	}

	srv.dbo.LockForReading()
	defer srv.dbo.UnlockAfterReading()

	// Read the resource.
	resource, err := srv.dbo.GetResourceById(p.ResourceId)
	if err != nil {
		return nil, srv.databaseError(err)
	}
	if resource == nil {
		return nil, jrm1.NewRpcErrorByUser(RpcErrorCode_ResourceIsNotFound, RpcErrorMsg_ResourceIsNotFound, nil)
	}

	result = &rpc2.GetResourceValueResult{
		Resource: nm.ResourceWithValue{
			Id:    p.ResourceId,
			Value: resource.GetValue(),
		},
	}

	return result, nil
}

func (srv *Server) getListOfAllResourcesOnPage(p *rpc2.GetListOfAllResourcesOnPageParams) (result *rpc2.GetListOfAllResourcesOnPageResult, re *jrm1.RpcError) {
	// Check parameters.
	if p.Page == 0 {
		return nil, jrm1.NewRpcErrorByUser(RpcErrorCode_PageIsNotSet, RpcErrorMsg_PageIsNotSet, nil)
	}

	var userRoles *am.GetSelfRolesResult
	userRoles, re = srv.mustBeAnAuthToken(p.Auth)
	if re != nil {
		return nil, re
	}

	// Check permissions.
	if !userRoles.User.GetUserParameters().GetRoles().IsReader {
		return nil, jrm1.NewRpcErrorByUser(c.RpcErrorCode_Permission, c.RpcErrorMsg_Permission, nil)
	}

	srv.dbo.LockForReading()
	defer srv.dbo.UnlockAfterReading()

	// Get resource IDs on page.
	resourceIds, err := srv.dbo.GetResourceIdsOnPage(p.Page, srv.settings.SystemSettings.PageSize)
	if err != nil {
		return nil, srv.databaseError(err)
	}

	// Count all resources.
	var allResourcesCount base2.Count
	allResourcesCount, err = srv.dbo.CountAllResources()
	if err != nil {
		return nil, srv.databaseError(err)
	}

	result = &rpc2.GetListOfAllResourcesOnPageResult{
		ResourcesOnPage: &nm.ResourcesOnPage{
			ResourceIds: resourceIds,
			PageData: &rpc3.PageData{
				PageNumber:  p.Page,
				TotalPages:  base2.CalculateTotalPages(allResourcesCount, srv.settings.SystemSettings.PageSize),
				PageSize:    srv.settings.SystemSettings.PageSize,
				ItemsOnPage: base2.Count(len(resourceIds)),
				TotalItems:  allResourcesCount,
			},
		},
	}

	return result, nil
}

// deleteResource removes a resource.
func (srv *Server) deleteResource(p *rpc2.DeleteResourceParams) (result *rpc2.DeleteResourceResult, re *jrm1.RpcError) {
	// Check parameters.
	if p.ResourceId == 0 {
		return nil, jrm1.NewRpcErrorByUser(RpcErrorCode_ResourceIdIsNotSet, RpcErrorMsg_ResourceIdIsNotSet, nil)
	}

	var userRoles *am.GetSelfRolesResult
	userRoles, re = srv.mustBeAnAuthToken(p.Auth)
	if re != nil {
		return nil, re
	}

	// Check permissions.
	if !userRoles.User.GetUserParameters().GetRoles().IsAdministrator {
		return nil, jrm1.NewRpcErrorByUser(c.RpcErrorCode_Permission, c.RpcErrorMsg_Permission, nil)
	}

	srv.dbo.LockForWriting()
	defer srv.dbo.UnlockAfterWriting()

	// Delete the resource.
	err := srv.dbo.DeleteResourceById(p.ResourceId)
	if err != nil {
		return nil, srv.databaseError(err)
	}

	result = &rpc2.DeleteResourceResult{
		Success: rpc3.Success{
			OK: true,
		},
	}
	return result, nil
}

// Other.

// processSystemEventS processes a system event. This method is used by the
// system.
func (srv *Server) processSystemEventS(p *rpc2.ProcessSystemEventSParams) (result *rpc2.ProcessSystemEventSResult, re *jrm1.RpcError) {
	// Check parameters.
	se, err := cm.NewSystemEventWithData(p.SystemEventData)
	if err != nil {
		return nil, jrm1.NewRpcErrorByUser(RpcErrorCode_SystemEvent, RpcErrorMsg_SystemEvent, nil)
	}

	re = srv.mustBeNoAuth(p.Auth)
	if re != nil {
		return nil, re
	}

	// Check the DKey.
	if !srv.dKeyI.CheckString(p.DKey.ToString()) {
		srv.incidentManager.ReportIncident(ev.NewEnumValue(it.IncidentType_WrongDKey), "", nil)
		return nil, jrm1.NewRpcErrorByUser(c.RpcErrorCode_Permission, c.RpcErrorMsg_Permission, nil)
	}

	re = srv.saveSystemEventH(se)
	if re != nil {
		return nil, re
	}

	seType := se.GetSystemEventData().GetType()
	switch seType.AsInt() {
	case set.SystemEventType_ThreadParentChange:
		re = srv.processSystemEvent_ThreadParentChange(se)
	case set.SystemEventType_ThreadNameChange:
		re = srv.processSystemEvent_ThreadNameChange(se)
	case set.SystemEventType_ThreadDeletion:
		re = srv.processSystemEvent_ThreadDeletion(se)
	case set.SystemEventType_ThreadNewMessage:
		re = srv.processSystemEvent_ThreadNewMessage(se)
	case set.SystemEventType_ThreadMessageEdit:
		re = srv.processSystemEvent_ThreadMessageEdit(se)
	case set.SystemEventType_ThreadMessageDeletion:
		re = srv.processSystemEvent_ThreadMessageDeletion(se)
	case set.SystemEventType_MessageTextEdit:
		re = srv.processSystemEvent_MessageTextEdit(se)
	case set.SystemEventType_MessageParentChange:
		re = srv.processSystemEvent_MessageParentChange(se)
	case set.SystemEventType_MessageDeletion:
		re = srv.processSystemEvent_MessageDeletion(se)

	default:
		return nil, jrm1.NewRpcErrorByUser(RpcErrorCode_SystemEvent, RpcErrorMsg_SystemEvent, nil)
	}

	if re != nil {
		return nil, re
	}

	result = &rpc2.ProcessSystemEventSResult{
		Success: rpc3.Success{
			OK: true,
		},
	}
	return result, nil
}

func (srv *Server) getDKey(p *rpc2.GetDKeyParams) (result *rpc2.GetDKeyResult, re *jrm1.RpcError) {
	re = srv.mustBeNoAuth(p.Auth)
	if re != nil {
		return nil, re
	}

	result = &rpc2.GetDKeyResult{
		DKey: base2.Text(srv.dKeyI.GetString()),
	}

	return result, nil
}

func (srv *Server) showDiagnosticData() (result *rpc2.ShowDiagnosticDataResult, re *jrm1.RpcError) {
	trc, src := srv.js.GetRequestsCount()

	result = &rpc2.ShowDiagnosticDataResult{
		RequestsCount: rpc3.RequestsCount{
			TotalRequestsCount:      base2.Text(trc),
			SuccessfulRequestsCount: base2.Text(src),
		},
	}

	return result, nil
}

func (srv *Server) test(p *rpc2.TestParams) (result *rpc2.TestResult, re *jrm1.RpcError) {
	result = &rpc2.TestResult{}

	var wg = new(sync.WaitGroup)
	var errChan = make(chan error, p.N)

	for i := uint(1); i <= p.N; i++ {
		wg.Add(1)
		go srv.doTestA(wg, errChan)
	}
	wg.Wait()
	close(errChan)

	for err := range errChan {
		if err != nil {
			srv.logError(err)
			return nil, jrm1.NewRpcErrorByUser(RpcErrorCode_TestError, fmt.Sprintf(RpcErrorMsgF_TestError, err.Error()), nil)
		}
	}

	return result, nil
}
