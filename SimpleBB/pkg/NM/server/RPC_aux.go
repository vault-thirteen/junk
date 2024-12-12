package server

import (
	"context"
	"fmt"
	jrm1 "github.com/vault-thirteen/JSON-RPC-M1"
	ac "github.com/vault-thirteen/SimpleBB/pkg/ACM/client"
	am "github.com/vault-thirteen/SimpleBB/pkg/ACM/rpc"
	nm "github.com/vault-thirteen/SimpleBB/pkg/NM/rpc"
	sc "github.com/vault-thirteen/SimpleBB/pkg/SM/client"
	sm "github.com/vault-thirteen/SimpleBB/pkg/SM/models"
	rpc2 "github.com/vault-thirteen/SimpleBB/pkg/SM/rpc"
	"github.com/vault-thirteen/SimpleBB/pkg/common/interfaces/derived2"
	base2 "github.com/vault-thirteen/SimpleBB/pkg/common/models/base"
	"github.com/vault-thirteen/SimpleBB/pkg/common/models/complex/EnumValue"
	"github.com/vault-thirteen/SimpleBB/pkg/common/models/complex/IncidentType"
	set "github.com/vault-thirteen/SimpleBB/pkg/common/models/complex/SystemEventType"
	cn "github.com/vault-thirteen/SimpleBB/pkg/common/models/net"
	rpc3 "github.com/vault-thirteen/SimpleBB/pkg/common/models/rpc"
	server2 "github.com/vault-thirteen/SimpleBB/pkg/common/models/server"
	"log"
	"sync"
)

// Auxiliary functions used in RPC functions.

// logError logs error if debug mode is enabled.
func (srv *Server) logError(err error) {
	if err == nil {
		return
	}

	if srv.settings.SystemSettings.IsDebugMode {
		log.Println(err)
	}
}

// processDatabaseError processes a database error.
func (srv *Server) processDatabaseError(err error) {
	if err == nil {
		return
	}

	if server2.IsNetworkError(err) {
		log.Println(fmt.Sprintf(server2.ErrFDatabaseNetwork, err.Error()))
		*(srv.dbErrors) <- err
	} else {
		srv.logError(err)
	}

	return
}

// databaseError processes the database error and returns an RPC error.
func (srv *Server) databaseError(err error) (re *jrm1.RpcError) {
	srv.processDatabaseError(err)
	return jrm1.NewRpcErrorByUser(server2.RpcErrorCode_Database, server2.RpcErrorMsg_Database, err)
}

// Token-related functions.

// mustBeNoAuth ensures that authorisation is not used.
func (srv *Server) mustBeNoAuth(auth *rpc3.Auth) (re *jrm1.RpcError) {
	if auth != nil {
		srv.incidentManager.ReportIncident(ev.NewEnumValue(it.IncidentType_IllegalAccessAttempt), "", auth.UserIPAB)
		return jrm1.NewRpcErrorByUser(server2.RpcErrorCode_Permission, server2.RpcErrorMsg_Permission, nil)
	}

	return nil
}

// mustBeAuthUserIPA ensures that user's IP address is set. If it is not set,
// an error is returned and the caller of this function must stop and return
// this error.
func (srv *Server) mustBeAuthUserIPA(auth *rpc3.Auth) (re *jrm1.RpcError) {
	if (auth == nil) ||
		(len(auth.UserIPA) == 0) {
		return jrm1.NewRpcErrorByUser(server2.RpcErrorCode_Authorisation, server2.RpcErrorMsg_Authorisation, nil)
	}

	var err error
	auth.UserIPAB, err = cn.ParseIPA(auth.UserIPA)
	if err != nil {
		srv.logError(err)
		return jrm1.NewRpcErrorByUser(server2.RpcErrorCode_Authorisation, server2.RpcErrorMsg_Authorisation, nil)
	}

	return nil
}

// mustBeNoAuthToken ensures that an authorisation token is not present. If the
// token is present, an error is returned and the caller of this function must
// stop and return this error.
func (srv *Server) mustBeNoAuthToken(auth *rpc3.Auth) (re *jrm1.RpcError) {
	re = srv.mustBeAuthUserIPA(auth)
	if re != nil {
		return re
	}

	if len(auth.Token) > 0 {
		srv.incidentManager.ReportIncident(ev.NewEnumValue(it.IncidentType_IllegalAccessAttempt), "", nil)
		return jrm1.NewRpcErrorByUser(server2.RpcErrorCode_Permission, server2.RpcErrorMsg_Permission, nil)
	}

	return nil
}

// mustBeAnAuthToken ensures that an authorisation token is present and is
// valid. If the token is absent or invalid, an error is returned and the caller
// of this function must stop and return this error. User data is returned when
// token is valid.
func (srv *Server) mustBeAnAuthToken(auth *rpc3.Auth) (userRoles *am.GetSelfRolesResult, re *jrm1.RpcError) {
	re = srv.mustBeAuthUserIPA(auth)
	if re != nil {
		return nil, re
	}

	if len(auth.Token) == 0 {
		return nil, jrm1.NewRpcErrorByUser(server2.RpcErrorCode_Authorisation, server2.RpcErrorMsg_Authorisation, nil)
	}

	userRoles, re = srv.getUserSelfRoles(auth)
	if re != nil {
		return nil, re
	}

	return userRoles, nil
}

// Other functions.

// getUserSelfRoles reads roles of the RPC caller (user) from ACM module.
func (srv *Server) getUserSelfRoles(auth *rpc3.Auth) (userRoles *am.GetSelfRolesResult, re *jrm1.RpcError) {
	var params = am.GetSelfRolesParams{
		CommonParams: rpc3.CommonParams{
			Auth: auth,
		},
	}

	userRoles = new(am.GetSelfRolesResult)
	var err error
	re, err = srv.acmServiceClient.MakeRequest(context.Background(), ac.FuncGetSelfRoles, params, userRoles)
	if err != nil {
		srv.logError(err)
		return nil, jrm1.NewRpcErrorByUser(server2.RpcErrorCode_RPCCall, server2.RpcErrorMsg_RPCCall, nil)
	}
	if re != nil {
		return nil, jrm1.NewRpcErrorByUser(server2.RpcErrorCode_Authorisation, server2.RpcErrorMsg_Authorisation, nil)
	}

	return userRoles, nil
}

func (srv *Server) doTestA(wg *sync.WaitGroup, errChan chan error) {
	defer wg.Done()

	var ap = am.TestParams{}

	var ar = new(am.TestResult)
	re, err := srv.acmServiceClient.MakeRequest(context.Background(), ac.FuncTest, ap, ar)
	if err != nil {
		errChan <- err
	}
	if re != nil {
		errChan <- re.AsError()
	}
}

// getDKeyForSM receives a DKey from Subscription module.
func (srv *Server) getDKeyForSM() (dKey *base2.Text, re *jrm1.RpcError) {
	params := rpc2.GetDKeyParams{}
	result := new(rpc2.GetDKeyResult)
	var err error
	re, err = srv.smServiceClient.MakeRequest(context.Background(), sc.FuncGetDKey, params, result)
	if err != nil {
		srv.logError(err)
		return nil, jrm1.NewRpcErrorByUser(server2.RpcErrorCode_RPCCall, server2.RpcErrorMsg_RPCCall, nil)
	}
	if re != nil {
		return nil, re
	}

	// DKey must be non-empty.
	if len(result.DKey) == 0 {
		return nil, jrm1.NewRpcErrorByUser(server2.RpcErrorCode_ModuleSynchronisation, server2.RpcErrorMsg_ModuleSynchronisation, nil)
	}

	return &result.DKey, nil
}

func tryGetSystemEventThreadId(se derived2.ISystemEvent) (threadId base2.Id, re *jrm1.RpcError) {
	x := se.GetSystemEventData().GetThreadId()
	if x == nil {
		return threadId, jrm1.NewRpcErrorByUser(RpcErrorCode_ThreadIdIsNotSet, RpcErrorMsg_ThreadIdIsNotSet, nil)
	}

	return *x, nil
}

func tryGetSystemEventUserId(se derived2.ISystemEvent) (userId base2.Id, re *jrm1.RpcError) {
	x := se.GetSystemEventData().GetUserId()
	if x == nil {
		return userId, jrm1.NewRpcErrorByUser(RpcErrorCode_UserIdIsNotSet, RpcErrorMsg_UserIdIsNotSet, nil)
	}

	return *x, nil
}

func tryGetSystemEventCreatorId(se derived2.ISystemEvent) (creatorId base2.Id, re *jrm1.RpcError) {
	x := se.GetSystemEventData().GetCreator()
	if x == nil {
		return creatorId, jrm1.NewRpcErrorByUser(RpcErrorCode_CreatorIsNotSet, RpcErrorMsg_CreatorIsNotSet, nil)
	}

	return *x, nil
}

func (srv *Server) processSystemEvent_ThreadParentChange(se derived2.ISystemEvent) (re *jrm1.RpcError) {
	return srv.sendNotificationsToThreadSubscribers(se)
}

func (srv *Server) processSystemEvent_ThreadNameChange(se derived2.ISystemEvent) (re *jrm1.RpcError) {
	return srv.sendNotificationsToThreadSubscribers(se)
}

func (srv *Server) processSystemEvent_ThreadDeletion(se derived2.ISystemEvent) (re *jrm1.RpcError) {
	var threadId base2.Id
	threadId, re = tryGetSystemEventThreadId(se)
	if re != nil {
		return re
	}

	re = srv.sendNotificationsToThreadSubscribers(se)
	if re != nil {
		return re
	}

	// Ask the SM module to clear the subscriptions.
	re = srv.clearSubscriptionsOfDeletedThread(threadId)
	if re != nil {
		return re
	}

	return nil
}

func (srv *Server) processSystemEvent_ThreadNewMessage(se derived2.ISystemEvent) (re *jrm1.RpcError) {
	return srv.sendNotificationsToThreadSubscribers(se)
}

func (srv *Server) processSystemEvent_ThreadMessageEdit(se derived2.ISystemEvent) (re *jrm1.RpcError) {
	return srv.sendNotificationsToThreadSubscribers(se)
}

func (srv *Server) processSystemEvent_ThreadMessageDeletion(se derived2.ISystemEvent) (re *jrm1.RpcError) {
	return srv.sendNotificationsToThreadSubscribers(se)
}

func (srv *Server) processSystemEvent_MessageTextEdit(se derived2.ISystemEvent) (re *jrm1.RpcError) {
	var threadId, userId, creatorId base2.Id
	threadId, re = tryGetSystemEventThreadId(se)
	if re != nil {
		return re
	}
	userId, re = tryGetSystemEventUserId(se)
	if re != nil {
		return re
	}
	creatorId, re = tryGetSystemEventCreatorId(se)
	if re != nil {
		return re
	}

	var isSubscribed base2.Flag
	isSubscribed, re = srv.isUserSubscribed(threadId, creatorId)
	if re != nil {
		return re
	}

	// If user is subscribed to the thread, it does not need the second
	// notification about this message.
	if isSubscribed {
		return nil
	}

	// The actor does not need a notification about self action.
	if userId == creatorId {
		return nil
	}

	return srv.sendNotificationToCreator(se)
}

func (srv *Server) processSystemEvent_MessageParentChange(se derived2.ISystemEvent) (re *jrm1.RpcError) {
	var userId, creatorId base2.Id
	userId, re = tryGetSystemEventUserId(se)
	if re != nil {
		return re
	}
	creatorId, re = tryGetSystemEventCreatorId(se)
	if re != nil {
		return re
	}

	// The actor does not need a notification about self action.
	if userId == creatorId {
		return nil
	}

	return srv.sendNotificationToCreator(se)
}

func (srv *Server) processSystemEvent_MessageDeletion(se derived2.ISystemEvent) (re *jrm1.RpcError) {
	var threadId, userId, creatorId base2.Id
	threadId, re = tryGetSystemEventThreadId(se)
	if re != nil {
		return re
	}
	userId, re = tryGetSystemEventUserId(se)
	if re != nil {
		return re
	}
	creatorId, re = tryGetSystemEventCreatorId(se)
	if re != nil {
		return re
	}

	var isSubscribed base2.Flag
	isSubscribed, re = srv.isUserSubscribed(threadId, creatorId)
	if re != nil {
		return re
	}

	// If user is subscribed to the thread, it does not need the second
	// notification about this message.
	if isSubscribed {
		return nil
	}

	// The actor does not need a notification about self action.
	if userId == creatorId {
		return nil
	}

	return srv.sendNotificationToCreator(se)
}

// sendNotificationsToThreadSubscribers sends notifications to thread
// subscribers.
func (srv *Server) sendNotificationsToThreadSubscribers(se derived2.ISystemEvent) (re *jrm1.RpcError) {
	var tsr *sm.ThreadSubscriptionsRecord
	tsr, re = srv.getThreadSubscribers(*se.GetSystemEventData().GetThreadId())
	if re != nil {
		return re
	}

	var notificationText base2.Text
	notificationText, re = srv.composeNotificationText(se)
	if re != nil {
		return re
	}

	if tsr != nil {
		for _, userId := range tsr.Users.AsArray() {
			// The performer of the action does not need the notification.
			if se.GetSystemEventData().GetUserId() != nil {
				if userId == *se.GetSystemEventData().GetUserId() {
					continue
				}
			}

			_, re = srv.sendNotificationIfPossibleH(userId, notificationText)
			if re != nil {
				return re
			}
		}
	}

	return nil
}

// sendNotificationToCreator sends a notification to the initial creator of the
// object.
func (srv *Server) sendNotificationToCreator(se derived2.ISystemEvent) (re *jrm1.RpcError) {
	var notificationText base2.Text
	notificationText, re = srv.composeNotificationText(se)
	if re != nil {
		return re
	}

	_, re = srv.sendNotificationIfPossibleH(*se.GetSystemEventData().GetCreator(), notificationText)
	if re != nil {
		return re
	}

	return nil
}

// composeNotificationText creates a text for notification about the system
// event.
func (srv *Server) composeNotificationText(se derived2.ISystemEvent) (text base2.Text, re *jrm1.RpcError) {
	switch se.GetSystemEventData().GetType().AsInt() {
	case set.SystemEventType_ThreadParentChange:
		// Template: FUT.
		text = base2.Text(fmt.Sprintf("A user (%d) has moved the thread (%d) into another forum.", *se.GetSystemEventData().GetUserId(), *se.GetSystemEventData().GetThreadId()))

	case set.SystemEventType_ThreadNameChange:
		// Template: FUT.
		text = base2.Text(fmt.Sprintf("A user (%d) has renamed the thread (%d).", *se.GetSystemEventData().GetUserId(), *se.GetSystemEventData().GetThreadId()))

	case set.SystemEventType_ThreadDeletion:
		// Template: FUT.
		text = base2.Text(fmt.Sprintf("A user (%d) has deleted the thread (%d).", *se.GetSystemEventData().GetUserId(), *se.GetSystemEventData().GetThreadId()))

	case set.SystemEventType_ThreadNewMessage:
		// Template: FUMT.
		text = base2.Text(fmt.Sprintf("A user (%d) has added a new message (%d) into the thread (%d).", *se.GetSystemEventData().GetUserId(), *se.GetSystemEventData().GetMessageId(), *se.GetSystemEventData().GetThreadId()))

	case set.SystemEventType_ThreadMessageEdit:
		// Template: FUMT.
		text = base2.Text(fmt.Sprintf("A user (%d) has edited a message (%d) in the thread (%d).", *se.GetSystemEventData().GetUserId(), *se.GetSystemEventData().GetMessageId(), *se.GetSystemEventData().GetThreadId()))

	case set.SystemEventType_ThreadMessageDeletion:
		// Template: FUMT.
		text = base2.Text(fmt.Sprintf("A user (%d) has deleted a message (%d) from the thread (%d).", *se.GetSystemEventData().GetUserId(), *se.GetSystemEventData().GetMessageId(), *se.GetSystemEventData().GetThreadId()))

	case set.SystemEventType_MessageTextEdit:
		// Template: FMTU.
		text = base2.Text(fmt.Sprintf("Your message (%d) in the thread (%d) was edited by a user (%d).", *se.GetSystemEventData().GetMessageId(), *se.GetSystemEventData().GetThreadId(), *se.GetSystemEventData().GetUserId()))

	case set.SystemEventType_MessageParentChange:
		// Template: FMTU.
		text = base2.Text(fmt.Sprintf("Your message (%d) in the thread (%d) was moved into another thread by a user (%d).", *se.GetSystemEventData().GetMessageId(), *se.GetSystemEventData().GetThreadId(), *se.GetSystemEventData().GetUserId()))

	case set.SystemEventType_MessageDeletion:
		// Template: FMTU.
		text = base2.Text(fmt.Sprintf("Your message (%d) in the thread (%d) was deleted by a user (%d).", *se.GetSystemEventData().GetMessageId(), *se.GetSystemEventData().GetThreadId(), *se.GetSystemEventData().GetUserId()))

	default:
		return "", jrm1.NewRpcErrorByUser(RpcErrorCode_SystemEvent, RpcErrorMsg_SystemEvent, nil)
	}

	return text, nil
}

// clearSubscriptionsOfDeletedThread clears remains of subscriptions of a
// deleted thread.
func (srv *Server) clearSubscriptionsOfDeletedThread(threadId base2.Id) (re *jrm1.RpcError) {
	params := rpc2.ClearThreadSubscriptionsSParams{
		DKeyParams: rpc3.DKeyParams{
			// DKey is set during module start-up, so it is non-null.
			DKey: *srv.dKeyForSM,
		},
		ThreadId: threadId,
	}
	result := new(rpc2.ClearThreadSubscriptionsSResult)
	var err error
	re, err = srv.smServiceClient.MakeRequest(context.Background(), sc.FuncClearThreadSubscriptionsS, params, result)
	if err != nil {
		srv.logError(err)
		return jrm1.NewRpcErrorByUser(server2.RpcErrorCode_RPCCall, server2.RpcErrorMsg_RPCCall, nil)
	}
	if re != nil {
		return re
	}

	return nil
}

// getThreadSubscribers gets a list of users subscribed to the thread.
func (srv *Server) getThreadSubscribers(threadId base2.Id) (tsr *sm.ThreadSubscriptionsRecord, re *jrm1.RpcError) {
	params := rpc2.GetThreadSubscribersSParams{
		DKeyParams: rpc3.DKeyParams{
			// DKey is set during module start-up, so it is non-null.
			DKey: *srv.dKeyForSM,
		},
		ThreadId: threadId,
	}
	result := new(rpc2.GetThreadSubscribersSResult)
	var err error
	re, err = srv.smServiceClient.MakeRequest(context.Background(), sc.FuncGetThreadSubscribersS, params, result)
	if err != nil {
		srv.logError(err)
		return nil, jrm1.NewRpcErrorByUser(server2.RpcErrorCode_RPCCall, server2.RpcErrorMsg_RPCCall, nil)
	}
	if re != nil {
		return nil, re
	}

	return result.ThreadSubscriptions, nil
}

// isUserSubscribed checks whether the user is subscribed to the thread.
func (srv *Server) isUserSubscribed(threadId base2.Id, userId base2.Id) (isSubscribed base2.Flag, re *jrm1.RpcError) {
	params := rpc2.IsUserSubscribedSParams{
		DKeyParams: rpc3.DKeyParams{
			// DKey is set during module start-up, so it is non-null.
			DKey: *srv.dKeyForSM,
		},
		ThreadId: threadId,
		UserId:   userId,
	}
	result := new(rpc2.IsUserSubscribedSResult)
	var err error
	re, err = srv.smServiceClient.MakeRequest(context.Background(), sc.FuncIsUserSubscribedS, params, result)
	if err != nil {
		srv.logError(err)
		return false, jrm1.NewRpcErrorByUser(server2.RpcErrorCode_RPCCall, server2.RpcErrorMsg_RPCCall, nil)
	}
	if re != nil {
		return false, re
	}

	return result.IsSubscribed, nil
}

// sendNotificationIfPossibleH is a helper function which tries to send a
// notification to a user when it is possible.
func (srv *Server) sendNotificationIfPossibleH(userId base2.Id, text base2.Text) (result *nm.SendNotificationIfPossibleSResult, re *jrm1.RpcError) {
	srv.dbo.LockForWriting()
	defer srv.dbo.UnlockAfterWriting()

	var err error
	var unc base2.Count
	unc, err = srv.dbo.CountUnreadNotificationsByUserId(userId)
	if err != nil {
		return nil, srv.databaseError(err)
	}

	// Notification box is full.
	if unc >= srv.settings.SystemSettings.NotificationCountLimit {
		result = &nm.SendNotificationIfPossibleSResult{
			IsSent: false,
		}
		return result, nil
	}

	var insertedNotificationId base2.Id
	insertedNotificationId, err = srv.dbo.InsertNewNotification(userId, text)
	if err != nil {
		return nil, srv.databaseError(err)
	}

	result = &nm.SendNotificationIfPossibleSResult{
		IsSent:         true,
		NotificationId: insertedNotificationId,
	}

	return result, nil
}

// saveSystemEventH saves the system event into database.
func (srv *Server) saveSystemEventH(se derived2.ISystemEvent) (re *jrm1.RpcError) {
	srv.dbo.LockForWriting()
	defer srv.dbo.UnlockAfterWriting()

	var err error
	err = srv.dbo.SaveSystemEvent(se)
	if err != nil {
		return srv.databaseError(err)
	}

	return nil
}
