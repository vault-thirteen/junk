package server

import (
	"fmt"
	am "github.com/vault-thirteen/SimpleBB/pkg/ACM/rpc"
	rpc2 "github.com/vault-thirteen/SimpleBB/pkg/SM/rpc"
	base2 "github.com/vault-thirteen/SimpleBB/pkg/common/models/base"
	rpc3 "github.com/vault-thirteen/SimpleBB/pkg/common/models/rpc"
	c "github.com/vault-thirteen/SimpleBB/pkg/common/models/server"
	"sync"

	jrm1 "github.com/vault-thirteen/JSON-RPC-M1"
	sm "github.com/vault-thirteen/SimpleBB/pkg/SM/models"
)

// RPC functions.

// Subscription.

// addSubscription creates a subscription.
func (srv *Server) addSubscription(p *rpc2.AddSubscriptionParams) (result *rpc2.AddSubscriptionResult, re *jrm1.RpcError) {
	// Check parameters.
	if p.ThreadId == 0 {
		return nil, jrm1.NewRpcErrorByUser(RpcErrorCode_ThreadIdIsNotSet, RpcErrorMsg_ThreadIdIsNotSet, nil)
	}
	if p.UserId == 0 {
		return nil, jrm1.NewRpcErrorByUser(RpcErrorCode_UserIdIsNotSet, RpcErrorMsg_UserIdIsNotSet, nil)
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
	if userRoles.User.GetUserParameters().GetId() != p.UserId {
		return nil, jrm1.NewRpcErrorByUser(c.RpcErrorCode_Permission, c.RpcErrorMsg_Permission, nil)
	}

	// Check existence of the thread.
	var threadExists base2.Flag
	threadExists, re = srv.checkIfThreadExists(p.ThreadId)
	if re != nil {
		return nil, re
	}

	// If thread does not exist, we can not subscribe to it.
	if !threadExists {
		return nil, jrm1.NewRpcErrorByUser(RpcErrorCode_ThreadDoesNotExist, RpcErrorMsg_ThreadDoesNotExist, nil)
	}

	srv.dbo.LockForWriting()
	defer srv.dbo.UnlockAfterWriting()

	// Read Subscriptions. If they are not initialised, initialise them. Other
	// methods reading subscriptions should not initialise un-initialised
	// subscriptions.
	var usr *sm.UserSubscriptionsRecord
	var err error
	usr, err = srv.dbo.GetUserSubscriptions(p.UserId)
	if err != nil {
		return nil, srv.databaseError(err)
	}

	// If no real record exists, create it.
	if usr.Id == sm.IdForVirtualUserSubscriptionsRecord {
		err = srv.dbo.InitUserSubscriptions(p.UserId)
		if err != nil {
			return nil, srv.databaseError(err)
		}

		usr, err = srv.dbo.GetUserSubscriptions(p.UserId)
		if err != nil {
			return nil, srv.databaseError(err)
		}
	}

	var tsr *sm.ThreadSubscriptionsRecord
	tsr, err = srv.dbo.GetThreadSubscriptions(p.ThreadId)
	if err != nil {
		return nil, srv.databaseError(err)
	}

	// If no real record exists, create it.
	if tsr == nil {
		err = srv.dbo.InitThreadSubscriptions(p.ThreadId)
		if err != nil {
			return nil, srv.databaseError(err)
		}

		tsr, err = srv.dbo.GetThreadSubscriptions(p.ThreadId)
		if err != nil {
			return nil, srv.databaseError(err)
		}
	}

	// Add items.
	err = usr.Threads.AddItem(p.ThreadId, false)
	if err != nil {
		srv.logError(err)
		return nil, jrm1.NewRpcErrorByUser(c.RpcErrorCode_UidList, fmt.Sprintf(c.RpcErrorMsgF_UidList, err.Error()), nil)
	}

	err = tsr.Users.AddItem(p.UserId, false)
	if err != nil {
		srv.logError(err)
		return nil, jrm1.NewRpcErrorByUser(c.RpcErrorCode_UidList, fmt.Sprintf(c.RpcErrorMsgF_UidList, err.Error()), nil)
	}

	// Save changes.
	err = srv.dbo.SaveUserSubscriptions(usr)
	if err != nil {
		return nil, srv.databaseError(err)
	}

	err = srv.dbo.SaveThreadSubscriptions(tsr)
	if err != nil {
		return nil, srv.databaseError(err)
	}

	result = &rpc2.AddSubscriptionResult{
		Success: rpc3.Success{
			OK: true,
		},
	}
	return result, nil
}

// isSelfSubscribed checks whether the caller user has a subscription to the thread.
func (srv *Server) isSelfSubscribed(p *rpc2.IsSelfSubscribedParams) (result *rpc2.IsSelfSubscribedResult, re *jrm1.RpcError) {
	// Check parameters.
	if p.ThreadId == 0 {
		return nil, jrm1.NewRpcErrorByUser(RpcErrorCode_ThreadIdIsNotSet, RpcErrorMsg_ThreadIdIsNotSet, nil)
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

	var isSubscribed base2.Flag
	isSubscribed, re = srv.isUserSubscribedH(userRoles.User.GetUserParameters().GetId(), p.ThreadId)
	if re != nil {
		return nil, re
	}

	result = &rpc2.IsSelfSubscribedResult{
		UserId:       userRoles.User.GetUserParameters().GetId(),
		ThreadId:     p.ThreadId,
		IsSubscribed: isSubscribed,
	}

	return result, nil
}

// isUserSubscribed checks whether the user has a subscription to the thread.
func (srv *Server) isUserSubscribed(p *rpc2.IsUserSubscribedParams) (result *rpc2.IsUserSubscribedResult, re *jrm1.RpcError) {
	// Check parameters.
	if p.ThreadId == 0 {
		return nil, jrm1.NewRpcErrorByUser(RpcErrorCode_ThreadIdIsNotSet, RpcErrorMsg_ThreadIdIsNotSet, nil)
	}
	if p.UserId == 0 {
		return nil, jrm1.NewRpcErrorByUser(RpcErrorCode_UserIdIsNotSet, RpcErrorMsg_UserIdIsNotSet, nil)
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
	if userRoles.User.GetUserParameters().GetId() != p.UserId {
		return nil, jrm1.NewRpcErrorByUser(c.RpcErrorCode_Permission, c.RpcErrorMsg_Permission, nil)
	}

	var isSubscribed base2.Flag
	isSubscribed, re = srv.isUserSubscribedH(p.UserId, p.ThreadId)
	if re != nil {
		return nil, re
	}

	result = &rpc2.IsUserSubscribedResult{
		UserId:       p.UserId,
		ThreadId:     p.ThreadId,
		IsSubscribed: isSubscribed,
	}

	return result, nil
}

// isUserSubscribedS checks whether the user has a subscription to the thread.
// This method is used by the system.
func (srv *Server) isUserSubscribedS(p *rpc2.IsUserSubscribedSParams) (result *rpc2.IsUserSubscribedSResult, re *jrm1.RpcError) {
	// Check parameters.
	if p.ThreadId == 0 {
		return nil, jrm1.NewRpcErrorByUser(RpcErrorCode_ThreadIdIsNotSet, RpcErrorMsg_ThreadIdIsNotSet, nil)
	}
	if p.UserId == 0 {
		return nil, jrm1.NewRpcErrorByUser(RpcErrorCode_UserIdIsNotSet, RpcErrorMsg_UserIdIsNotSet, nil)
	}

	re = srv.mustBeNoAuth(p.Auth)
	if re != nil {
		return nil, re
	}

	// Check the DKey.
	if !srv.dKeyI.CheckString(p.DKey.ToString()) {
		return nil, jrm1.NewRpcErrorByUser(c.RpcErrorCode_Permission, c.RpcErrorMsg_Permission, nil)
	}

	var isSubscribed base2.Flag
	isSubscribed, re = srv.isUserSubscribedH(p.UserId, p.ThreadId)
	if re != nil {
		return nil, re
	}

	result = &rpc2.IsUserSubscribedSResult{
		UserId:       p.UserId,
		ThreadId:     p.ThreadId,
		IsSubscribed: isSubscribed,
	}

	return result, nil
}

// countSelfSubscriptions counts subscriptions of the current user.
func (srv *Server) countSelfSubscriptions(p *rpc2.CountSelfSubscriptionsParams) (result *rpc2.CountSelfSubscriptionsResult, re *jrm1.RpcError) {
	var userRoles *am.GetSelfRolesResult
	userRoles, re = srv.mustBeAnAuthToken(p.Auth)
	if re != nil {
		return nil, re
	}

	// Check permissions.
	if !userRoles.User.GetUserParameters().GetRoles().IsReader {
		return nil, jrm1.NewRpcErrorByUser(c.RpcErrorCode_Permission, c.RpcErrorMsg_Permission, nil)
	}

	result = &rpc2.CountSelfSubscriptionsResult{}

	result.UserSubscriptionsCount, re = srv.countSelfSubscriptionsH(userRoles.User.GetUserParameters().GetId())
	if re != nil {
		return nil, re
	}

	return result, nil
}

// getSelfSubscriptions reads subscriptions of the current user.
func (srv *Server) getSelfSubscriptions(p *rpc2.GetSelfSubscriptionsParams) (result *rpc2.GetSelfSubscriptionsResult, re *jrm1.RpcError) {
	var userRoles *am.GetSelfRolesResult
	userRoles, re = srv.mustBeAnAuthToken(p.Auth)
	if re != nil {
		return nil, re
	}

	// Check permissions.
	if !userRoles.User.GetUserParameters().GetRoles().IsReader {
		return nil, jrm1.NewRpcErrorByUser(c.RpcErrorCode_Permission, c.RpcErrorMsg_Permission, nil)
	}

	var usr *sm.UserSubscriptionsRecord
	usr, re = srv.getUserSubscriptionsRecordH(userRoles.User.GetUserParameters().GetId())
	if re != nil {
		return nil, re
	}

	result = &rpc2.GetSelfSubscriptionsResult{
		UserSubscriptions: sm.NewUserSubscriptions(userRoles.User.GetUserParameters().GetId(), usr.Threads, 0, 0),
	}

	return result, nil
}

// getSelfSubscriptionsOnPage reads subscriptions of the current user on the selected page.
func (srv *Server) getSelfSubscriptionsOnPage(p *rpc2.GetSelfSubscriptionsOnPageParams) (result *rpc2.GetSelfSubscriptionsOnPageResult, re *jrm1.RpcError) {
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

	var usr *sm.UserSubscriptionsRecord
	usr, re = srv.getUserSubscriptionsRecordH(userRoles.User.GetUserParameters().GetId())
	if re != nil {
		return nil, re
	}

	result = &rpc2.GetSelfSubscriptionsOnPageResult{
		UserSubscriptions: sm.NewUserSubscriptions(userRoles.User.GetUserParameters().GetId(), usr.Threads, p.Page, srv.settings.SystemSettings.PageSize),
	}

	return result, nil
}

// getUserSubscriptions reads user subscriptions.
func (srv *Server) getUserSubscriptions(p *rpc2.GetUserSubscriptionsParams) (result *rpc2.GetUserSubscriptionsResult, re *jrm1.RpcError) {
	// Check parameters.
	if p.UserId == 0 {
		return nil, jrm1.NewRpcErrorByUser(RpcErrorCode_UserIdIsNotSet, RpcErrorMsg_UserIdIsNotSet, nil)
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
	if userRoles.User.GetUserParameters().GetId() != p.UserId {
		return nil, jrm1.NewRpcErrorByUser(c.RpcErrorCode_Permission, c.RpcErrorMsg_Permission, nil)
	}

	var usr *sm.UserSubscriptionsRecord
	usr, re = srv.getUserSubscriptionsRecordH(p.UserId)
	if re != nil {
		return nil, re
	}

	result = &rpc2.GetUserSubscriptionsResult{
		UserSubscriptions: sm.NewUserSubscriptions(p.UserId, usr.Threads, 0, 0),
	}

	return result, nil
}

// getUserSubscriptionsOnPage reads user subscriptions on the selected page.
func (srv *Server) getUserSubscriptionsOnPage(p *rpc2.GetUserSubscriptionsOnPageParams) (result *rpc2.GetUserSubscriptionsOnPageResult, re *jrm1.RpcError) {
	// Check parameters.
	if p.UserId == 0 {
		return nil, jrm1.NewRpcErrorByUser(RpcErrorCode_UserIdIsNotSet, RpcErrorMsg_UserIdIsNotSet, nil)
	}
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
	if userRoles.User.GetUserParameters().GetId() != p.UserId {
		return nil, jrm1.NewRpcErrorByUser(c.RpcErrorCode_Permission, c.RpcErrorMsg_Permission, nil)
	}

	var usr *sm.UserSubscriptionsRecord
	usr, re = srv.getUserSubscriptionsRecordH(p.UserId)
	if re != nil {
		return nil, re
	}

	result = &rpc2.GetUserSubscriptionsOnPageResult{
		UserSubscriptions: sm.NewUserSubscriptions(p.UserId, usr.Threads, p.Page, srv.settings.SystemSettings.PageSize),
	}

	return result, nil
}

// getThreadSubscribersS reads a list of users subscribed to the specified
// thread. This method is used by the system.
func (srv *Server) getThreadSubscribersS(p *rpc2.GetThreadSubscribersSParams) (result *rpc2.GetThreadSubscribersSResult, re *jrm1.RpcError) {
	// Check parameters.
	if p.ThreadId == 0 {
		return nil, jrm1.NewRpcErrorByUser(RpcErrorCode_ThreadIdIsNotSet, RpcErrorMsg_ThreadIdIsNotSet, nil)
	}

	re = srv.mustBeNoAuth(p.Auth)
	if re != nil {
		return nil, re
	}

	// Check the DKey.
	if !srv.dKeyI.CheckString(p.DKey.ToString()) {
		return nil, jrm1.NewRpcErrorByUser(c.RpcErrorCode_Permission, c.RpcErrorMsg_Permission, nil)
	}

	// Check existence of the thread.
	var threadExists base2.Flag
	threadExists, re = srv.checkIfThreadExists(p.ThreadId)
	if re != nil {
		return nil, re
	}

	// If thread does not exist, we do not even set the thread ID in result.
	if !threadExists {
		result = &rpc2.GetThreadSubscribersSResult{
			ThreadSubscriptions: nil,
		}
		return result, nil
	}

	srv.dbo.LockForReading()
	defer srv.dbo.UnlockAfterReading()

	// Read Subscriptions.
	var tsr *sm.ThreadSubscriptionsRecord
	var err error
	tsr, err = srv.dbo.GetThreadSubscriptions(p.ThreadId)
	if err != nil {
		return nil, srv.databaseError(err)
	}

	result = &rpc2.GetThreadSubscribersSResult{
		ThreadSubscriptions: tsr,
	}

	return result, nil
}

// deleteSelfSubscription deletes a subscription of the caller user.
func (srv *Server) deleteSelfSubscription(p *rpc2.DeleteSelfSubscriptionParams) (result *rpc2.DeleteSelfSubscriptionResult, re *jrm1.RpcError) {
	// Check parameters.
	if p.ThreadId == 0 {
		return nil, jrm1.NewRpcErrorByUser(RpcErrorCode_ThreadIdIsNotSet, RpcErrorMsg_ThreadIdIsNotSet, nil)
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

	// Delete subscription.
	var s = &sm.Subscription{
		ThreadId: p.ThreadId,
		UserId:   userRoles.User.GetUserParameters().GetId(),
	}
	re = srv.deleteSubscriptionH(s)
	if re != nil {
		return nil, re
	}

	result = &rpc2.DeleteSelfSubscriptionResult{
		Success: rpc3.Success{
			OK: true,
		},
	}
	return result, nil
}

// deleteSubscription deletes a subscription.
func (srv *Server) deleteSubscription(p *rpc2.DeleteSubscriptionParams) (result *rpc2.DeleteSubscriptionResult, re *jrm1.RpcError) {
	// Check parameters.
	if p.ThreadId == 0 {
		return nil, jrm1.NewRpcErrorByUser(RpcErrorCode_ThreadIdIsNotSet, RpcErrorMsg_ThreadIdIsNotSet, nil)
	}
	if p.UserId == 0 {
		return nil, jrm1.NewRpcErrorByUser(RpcErrorCode_UserIdIsNotSet, RpcErrorMsg_UserIdIsNotSet, nil)
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
	if userRoles.User.GetUserParameters().GetId() != p.UserId {
		return nil, jrm1.NewRpcErrorByUser(c.RpcErrorCode_Permission, c.RpcErrorMsg_Permission, nil)
	}

	// Delete subscription.
	var s = &sm.Subscription{
		ThreadId: p.ThreadId,
		UserId:   p.UserId,
	}
	re = srv.deleteSubscriptionH(s)
	if re != nil {
		return nil, re
	}

	result = &rpc2.DeleteSubscriptionResult{
		Success: rpc3.Success{
			OK: true,
		},
	}
	return result, nil
}

// deleteSubscriptionS deletes a subscription. This method is used by the
// system.
func (srv *Server) deleteSubscriptionS(p *rpc2.DeleteSubscriptionSParams) (result *rpc2.DeleteSubscriptionSResult, re *jrm1.RpcError) {
	// Check parameters.
	if p.ThreadId == 0 {
		return nil, jrm1.NewRpcErrorByUser(RpcErrorCode_ThreadIdIsNotSet, RpcErrorMsg_ThreadIdIsNotSet, nil)
	}
	if p.UserId == 0 {
		return nil, jrm1.NewRpcErrorByUser(RpcErrorCode_UserIdIsNotSet, RpcErrorMsg_UserIdIsNotSet, nil)
	}

	re = srv.mustBeNoAuth(p.Auth)
	if re != nil {
		return nil, re
	}

	// Check the DKey.
	if !srv.dKeyI.CheckString(p.DKey.ToString()) {
		return nil, jrm1.NewRpcErrorByUser(c.RpcErrorCode_Permission, c.RpcErrorMsg_Permission, nil)
	}

	// Delete subscription.
	var s = &sm.Subscription{
		ThreadId: p.ThreadId,
		UserId:   p.UserId,
	}
	re = srv.deleteSubscriptionH(s)
	if re != nil {
		return nil, re
	}

	result = &rpc2.DeleteSubscriptionSResult{
		Success: rpc3.Success{
			OK: true,
		},
	}
	return result, nil
}

// clearThreadSubscriptionsS clears subscriptions of an existing thread. This
// method is used by the system.
func (srv *Server) clearThreadSubscriptionsS(p *rpc2.ClearThreadSubscriptionsSParams) (result *rpc2.ClearThreadSubscriptionsSResult, re *jrm1.RpcError) {
	// Check parameters.
	if p.ThreadId == 0 {
		return nil, jrm1.NewRpcErrorByUser(RpcErrorCode_ThreadIdIsNotSet, RpcErrorMsg_ThreadIdIsNotSet, nil)
	}

	re = srv.mustBeNoAuth(p.Auth)
	if re != nil {
		return nil, re
	}

	// Check the DKey.
	if !srv.dKeyI.CheckString(p.DKey.ToString()) {
		return nil, jrm1.NewRpcErrorByUser(c.RpcErrorCode_Permission, c.RpcErrorMsg_Permission, nil)
	}

	// Check existence of the thread.
	var threadExists base2.Flag
	threadExists, re = srv.checkIfThreadExists(p.ThreadId)
	if re != nil {
		return nil, re
	}

	if !threadExists {
		return nil, jrm1.NewRpcErrorByUser(RpcErrorCode_ThreadDoesNotExist, RpcErrorMsg_ThreadDoesNotExist, nil)
	}

	// Clear subscriptions.
	re = srv.clearThreadSubscriptionsH(p.ThreadId)
	if re != nil {
		return nil, re
	}

	result = &rpc2.ClearThreadSubscriptionsSResult{
		Success: rpc3.Success{
			OK: true,
		},
	}
	return result, nil
}

// Other.

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
