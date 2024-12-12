package server

import (
	"bytes"
	"context"
	"fmt"
	am "github.com/vault-thirteen/SimpleBB/pkg/ACM/rpc"
	rpc2 "github.com/vault-thirteen/SimpleBB/pkg/MM/rpc"
	nm "github.com/vault-thirteen/SimpleBB/pkg/NM/rpc"
	"github.com/vault-thirteen/SimpleBB/pkg/common/interfaces/derived2"
	"github.com/vault-thirteen/SimpleBB/pkg/common/models/UidList"
	base2 "github.com/vault-thirteen/SimpleBB/pkg/common/models/base"
	cn "github.com/vault-thirteen/SimpleBB/pkg/common/models/net"
	rpc3 "github.com/vault-thirteen/SimpleBB/pkg/common/models/rpc"
	server2 "github.com/vault-thirteen/SimpleBB/pkg/common/models/server"
	"log"
	"sync"
	"time"

	jrm1 "github.com/vault-thirteen/JSON-RPC-M1"
	ac "github.com/vault-thirteen/SimpleBB/pkg/ACM/client"
	nc "github.com/vault-thirteen/SimpleBB/pkg/NM/client"
	ah "github.com/vault-thirteen/auxie/hash"
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

// getDKeyForNM receives a DKey from Notification module.
func (srv *Server) getDKeyForNM() (dKey *base2.Text, re *jrm1.RpcError) {
	params := nm.GetDKeyParams{}
	result := new(nm.GetDKeyResult)
	var err error
	re, err = srv.nmServiceClient.MakeRequest(context.Background(), nc.FuncGetDKey, params, result)
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

// sendNotificationToUser sends a system notification to a user.
func (srv *Server) sendNotificationToUser(userId base2.Id, text base2.Text) (re *jrm1.RpcError) {
	params := nm.AddNotificationSParams{
		DKeyParams: rpc3.DKeyParams{
			// DKey is set during module start-up, so it is non-null.
			DKey: *srv.dKeyForNM,
		},
		UserId: userId,
		Text:   text,
	}
	result := new(nm.AddNotificationSResult)
	var err error
	re, err = srv.nmServiceClient.MakeRequest(context.Background(), nc.FuncAddNotificationS, params, result)
	if err != nil {
		srv.logError(err)
		return jrm1.NewRpcErrorByUser(server2.RpcErrorCode_RPCCall, server2.RpcErrorMsg_RPCCall, nil)
	}
	if re != nil {
		return re
	}

	return nil
}

func (srv *Server) getMessageTextChecksum(msgText base2.Text) (checksum []byte) {
	x, _ := ah.CalculateCrc32(msgText.AsBytes())
	return x[:]
}

func (srv *Server) checkMessageTextChecksum(msgText base2.Text, checksum []byte) (ok bool) {
	return bytes.Compare(srv.getMessageTextChecksum(msgText), checksum) == 0
}

// canUserEditMessage checks whether a user (specified by the 'userRoles'
// argument) can edit a message (specified as an 'message' argument).
func (srv *Server) canUserEditMessage(userRoles *am.GetSelfRolesResult, message derived2.IMessage) (ok bool) {
	// Are you stupid, kido ?
	if (userRoles == nil) || (message == nil) {
		return false
	}

	// Moderators have extended rights to edit messages of any users.
	if userRoles.User.GetUserParameters().GetRoles().IsModerator {
		return true
	}

	// Writers can edit their own messages.
	if !userRoles.User.GetUserParameters().GetRoles().IsWriter {
		return false
	}

	// User can not edit messages created by other users.
	if userRoles.User.GetUserParameters().GetId() != message.GetEventData().GetCreatorUserId() {
		return false
	}

	// User can edit its own messages if they are not too old.
	if time.Now().Before(srv.getMessageMaxEditTime(message)) {
		return true
	}

	return false
}

// canUserAddMessage checks whether a user (specified by the 'userRoles'
// argument) can add a new message into a thread in case when there is a
// [latest] message in the thread (specified as an 'latestMessageInThread'
// argument). If the thread is empty, i.e. no latest message is available, it
// must be set as null.
func (srv *Server) canUserAddMessage(userRoles *am.GetSelfRolesResult, latestMessageInThread derived2.IMessage) (ok bool) {
	// Are you stupid, kido ?
	if userRoles == nil {
		return false
	}

	// Only writers can add new messages.
	if !userRoles.User.GetUserParameters().GetRoles().IsWriter {
		return false
	}

	// If there is no latest message in the thread, the thread is empty.
	if latestMessageInThread == nil {
		return true
	}

	// If the latest message in the thread was written by another [that] user,
	// this user can add a new message.
	if latestMessageInThread.GetEventData().GetCreatorUserId() != userRoles.User.GetUserParameters().GetId() {
		return true
	}

	// If the user is trying to add its another message into a thread, we must
	// check for collision. If the latest message can be edited, it should be
	// edited, and another new message is not allowed.
	if time.Now().Before(srv.getMessageMaxEditTime(latestMessageInThread)) {
		return false
	}

	return true
}

// getLatestMessageOfThreadH is a helper function used by other functions to
// get the latest message of a thread.
func (srv *Server) getLatestMessageOfThreadH(threadId base2.Id) (message derived2.IMessage, re *jrm1.RpcError) {
	// Get the thread.
	var thread derived2.IThread
	var err error
	thread, err = srv.dbo.GetThreadById(threadId)
	if err != nil {
		return nil, srv.databaseError(err)
	}

	if thread == nil {
		return nil, jrm1.NewRpcErrorByUser(RpcErrorCode_ThreadIsNotFound, RpcErrorMsg_ThreadIsNotFound, nil)
	}

	if thread.GetMessages() == nil {
		return nil, nil
	}

	latestMessageId := thread.GetMessages().LastElement()
	if latestMessageId == nil {
		return nil, nil
	}

	// Read the message.
	message, err = srv.dbo.GetMessageById(*latestMessageId)
	if err != nil {
		return nil, srv.databaseError(err)
	}

	if message == nil {
		return nil, jrm1.NewRpcErrorByUser(RpcErrorCode_MessageIsNotFound, RpcErrorMsg_MessageIsNotFound, nil)
	}

	return message, nil
}

// getMessageMaxEditTime returns the time border after which message editing is
// not allowed for an ordinary (non-moderator) user.
func (srv *Server) getMessageMaxEditTime(message derived2.IMessage) time.Time {
	return message.GetLastTouchTime().Add(time.Second * time.Duration(srv.settings.SystemSettings.MessageEditTime))
}

// deleteThreadH is a helper function used by other functions to delete a
// thread.
func (srv *Server) deleteThreadH(threadId base2.Id) (re *jrm1.RpcError) {
	srv.dbo.LockForWriting()
	defer srv.dbo.UnlockAfterWriting()

	// Read the thread.
	var thread derived2.IThread
	var err error
	thread, err = srv.dbo.GetThreadById(threadId)
	if err != nil {
		return srv.databaseError(err)
	}

	if thread == nil {
		return jrm1.NewRpcErrorByUser(RpcErrorCode_ThreadIsNotFound, RpcErrorMsg_ThreadIsNotFound, nil)
	}

	// Check for derived1.
	if thread.GetMessages().Size() > 0 {
		return jrm1.NewRpcErrorByUser(RpcErrorCode_ThreadIsNotEmpty, RpcErrorMsg_ThreadIsNotEmpty, nil)
	}

	// Update the link.
	var linkThreads *ul.UidList
	linkThreads, err = srv.dbo.GetForumThreadsById(thread.GetForumId())
	if err != nil {
		return srv.databaseError(err)
	}

	err = linkThreads.RemoveItem(threadId)
	if err != nil {
		srv.logError(err)
		return jrm1.NewRpcErrorByUser(server2.RpcErrorCode_UidList, fmt.Sprintf(server2.RpcErrorMsgF_UidList, err.Error()), nil)
	}

	err = srv.dbo.SetForumThreadsById(thread.GetForumId(), linkThreads)
	if err != nil {
		return srv.databaseError(err)
	}

	// Delete the thread.
	err = srv.dbo.DeleteThreadById(threadId)
	if err != nil {
		return srv.databaseError(err)
	}

	return nil
}

// changeThreadForumH is a helper function used by other functions to move a
// thread from an old forum to a new forum.
func (srv *Server) changeThreadForumH(threadId base2.Id, newForumId base2.Id, userId base2.Id) (re *jrm1.RpcError) {
	srv.dbo.LockForWriting()
	defer srv.dbo.UnlockAfterWriting()

	var n base2.Count
	var err error
	n, err = srv.dbo.CountThreadsById(threadId)
	if err != nil {
		return srv.databaseError(err)
	}

	if n == 0 {
		return jrm1.NewRpcErrorByUser(RpcErrorCode_ThreadIsNotFound, RpcErrorMsg_ThreadIsNotFound, nil)
	}

	// Ensure that an old parent exists.
	var oldParent base2.Id
	oldParent, err = srv.dbo.GetThreadForumById(threadId)
	if err != nil {
		return srv.databaseError(err)
	}

	n, err = srv.dbo.CountForumsById(oldParent)
	if err != nil {
		return srv.databaseError(err)
	}

	if n == 0 {
		return jrm1.NewRpcErrorByUser(RpcErrorCode_ForumIsNotFound, RpcErrorMsg_ForumIsNotFound, nil)
	}

	// Ensure that a new parent exists.
	n, err = srv.dbo.CountForumsById(newForumId)
	if err != nil {
		return srv.databaseError(err)
	}

	if n == 0 {
		return jrm1.NewRpcErrorByUser(RpcErrorCode_ForumIsNotFound, RpcErrorMsg_ForumIsNotFound, nil)
	}

	// Update the moved thread.
	err = srv.dbo.SetThreadForumById(threadId, newForumId, userId)
	if err != nil {
		return srv.databaseError(err)
	}

	// Update the new link.
	var threadsR *ul.UidList
	threadsR, err = srv.dbo.GetForumThreadsById(newForumId)
	if err != nil {
		return srv.databaseError(err)
	}

	err = threadsR.AddItem(threadId, false)
	if err != nil {
		srv.logError(err)
		return jrm1.NewRpcErrorByUser(server2.RpcErrorCode_UidList, fmt.Sprintf(server2.RpcErrorMsgF_UidList, err.Error()), nil)
	}

	err = srv.dbo.SetForumThreadsById(newForumId, threadsR)
	if err != nil {
		return srv.databaseError(err)
	}

	// Update the old link.
	var threadsL *ul.UidList
	threadsL, err = srv.dbo.GetForumThreadsById(oldParent)
	if err != nil {
		return srv.databaseError(err)
	}

	err = threadsL.RemoveItem(threadId)
	if err != nil {
		srv.logError(err)
		return jrm1.NewRpcErrorByUser(server2.RpcErrorCode_UidList, fmt.Sprintf(server2.RpcErrorMsgF_UidList, err.Error()), nil)
	}

	err = srv.dbo.SetForumThreadsById(oldParent, threadsL)
	if err != nil {
		return srv.databaseError(err)
	}

	return nil
}

// changeThreadNameH is a helper function used by other functions to rename a
// thread.
func (srv *Server) changeThreadNameH(threadId base2.Id, newThreadName base2.Text, userId base2.Id) (re *jrm1.RpcError) {
	srv.dbo.LockForWriting()
	defer srv.dbo.UnlockAfterWriting()

	var n base2.Count
	var err error
	n, err = srv.dbo.CountThreadsById(threadId)
	if err != nil {
		return srv.databaseError(err)
	}

	if n == 0 {
		return jrm1.NewRpcErrorByUser(RpcErrorCode_ThreadIsNotFound, RpcErrorMsg_ThreadIsNotFound, nil)
	}

	err = srv.dbo.SetThreadNameById(threadId, newThreadName, userId)
	if err != nil {
		return srv.databaseError(err)
	}

	return nil
}

// addMessageH is a helper function used by other functions to inserts a new
// message into a thread.
func (srv *Server) addMessageH(threadId base2.Id, messageText base2.Text, userRoles *am.GetSelfRolesResult) (result *rpc2.AddMessageResult, re *jrm1.RpcError) {
	srv.dbo.LockForWriting()
	defer srv.dbo.UnlockAfterWriting()

	// Get the latest message in the thread.
	var latestMessageInThread derived2.IMessage
	latestMessageInThread, re = srv.getLatestMessageOfThreadH(threadId)
	if re != nil {
		return nil, re
	}

	// Check permissions (Part II).
	canIAddMessage := srv.canUserAddMessage(userRoles, latestMessageInThread)
	if !canIAddMessage {
		return nil, jrm1.NewRpcErrorByUser(server2.RpcErrorCode_Permission, server2.RpcErrorMsg_Permission, nil)
	}

	// Ensure that a parent exists.
	var err error
	var n base2.Count
	n, err = srv.dbo.CountThreadsById(threadId)
	if err != nil {
		return nil, srv.databaseError(err)
	}

	if n == 0 {
		return nil, jrm1.NewRpcErrorByUser(RpcErrorCode_ThreadIsNotFound, RpcErrorMsg_ThreadIsNotFound, nil)
	}

	// Insert a message and link it with its thread.
	var parentMessages *ul.UidList
	parentMessages, err = srv.dbo.GetThreadMessagesById(threadId)
	if err != nil {
		return nil, srv.databaseError(err)
	}

	messageTextChecksum := srv.getMessageTextChecksum(messageText)

	var insertedMessageId base2.Id
	insertedMessageId, err = srv.dbo.InsertNewMessage(threadId, messageText, messageTextChecksum, userRoles.User.GetUserParameters().GetId())
	if err != nil {
		return nil, srv.databaseError(err)
	}

	err = parentMessages.AddItem(insertedMessageId, false)
	if err != nil {
		srv.logError(err)
		return nil, jrm1.NewRpcErrorByUser(server2.RpcErrorCode_UidList, fmt.Sprintf(server2.RpcErrorMsgF_UidList, err.Error()), nil)
	}

	err = srv.dbo.SetThreadMessagesById(threadId, parentMessages)
	if err != nil {
		return nil, srv.databaseError(err)
	}

	// Update thread's position if needed.
	if srv.settings.SystemSettings.NewThreadsAtTop {
		var messageThread derived2.IThread
		messageThread, err = srv.dbo.GetThreadById(threadId)
		if err != nil {
			return nil, srv.databaseError(err)
		}

		if messageThread == nil {
			return nil, jrm1.NewRpcErrorByUser(RpcErrorCode_ThreadIsNotFound, RpcErrorMsg_ThreadIsNotFound, nil)
		}

		var threads *ul.UidList
		threads, err = srv.dbo.GetForumThreadsById(messageThread.GetForumId())
		if err != nil {
			return nil, srv.databaseError(err)
		}

		var isAlreadyRaised bool
		isAlreadyRaised, err = threads.RaiseItem(threadId)
		if err != nil {
			srv.logError(err)
			return nil, jrm1.NewRpcErrorByUser(server2.RpcErrorCode_UidList, fmt.Sprintf(server2.RpcErrorMsgF_UidList, err.Error()), nil)
		}

		// Update the list if it has been changed.
		if !isAlreadyRaised {
			err = srv.dbo.SetForumThreadsById(messageThread.GetForumId(), threads)
			if err != nil {
				return nil, srv.databaseError(err)
			}
		}
	}

	result = &rpc2.AddMessageResult{
		MessageId: insertedMessageId,
	}

	return result, nil
}

// changeMessageTextH is a helper function used by other functions to change
// text of a message.
func (srv *Server) changeMessageTextH(messageId base2.Id, newText base2.Text, userRoles *am.GetSelfRolesResult) (initialMessage derived2.IMessage, re *jrm1.RpcError) {
	srv.dbo.LockForWriting()
	defer srv.dbo.UnlockAfterWriting()

	// Get the edited message.
	var err error
	initialMessage, err = srv.dbo.GetMessageById(messageId)
	if err != nil {
		return nil, srv.databaseError(err)
	}

	if initialMessage == nil {
		return nil, jrm1.NewRpcErrorByUser(RpcErrorCode_MessageIsNotFound, RpcErrorMsg_MessageIsNotFound, nil)
	}

	// Check permissions.
	canIEditMessage := srv.canUserEditMessage(userRoles, initialMessage)
	if !canIEditMessage {
		return nil, jrm1.NewRpcErrorByUser(server2.RpcErrorCode_Permission, server2.RpcErrorMsg_Permission, nil)
	}

	// Edit the message.
	messageTextChecksum := srv.getMessageTextChecksum(newText)

	err = srv.dbo.SetMessageTextById(messageId, newText, messageTextChecksum, userRoles.User.GetUserParameters().GetId())
	if err != nil {
		return nil, srv.databaseError(err)
	}

	return initialMessage, nil
}

// changeMessageThreadH is a helper function used by other functions to move a
// message from an old thread to a new thread.
func (srv *Server) changeMessageThreadH(messageId base2.Id, newThreadId base2.Id, userRoles *am.GetSelfRolesResult) (initialMessage derived2.IMessage, re *jrm1.RpcError) {
	srv.dbo.LockForWriting()
	defer srv.dbo.UnlockAfterWriting()

	var n base2.Count
	var err error
	n, err = srv.dbo.CountMessagesById(messageId)
	if err != nil {
		return nil, srv.databaseError(err)
	}

	if n == 0 {
		return nil, jrm1.NewRpcErrorByUser(RpcErrorCode_MessageIsNotFound, RpcErrorMsg_MessageIsNotFound, nil)
	}

	// Get the edited message.
	initialMessage, err = srv.dbo.GetMessageById(messageId)
	if err != nil {
		return nil, srv.databaseError(err)
	}

	// Ensure that an old parent exists.
	var oldParent base2.Id
	oldParent, err = srv.dbo.GetMessageThreadById(messageId)
	if err != nil {
		return nil, srv.databaseError(err)
	}

	n, err = srv.dbo.CountThreadsById(oldParent)
	if err != nil {
		return nil, srv.databaseError(err)
	}

	if n == 0 {
		return nil, jrm1.NewRpcErrorByUser(RpcErrorCode_ThreadIsNotFound, RpcErrorMsg_ThreadIsNotFound, nil)
	}

	// Ensure that a new parent exists.
	n, err = srv.dbo.CountThreadsById(newThreadId)
	if err != nil {
		return nil, srv.databaseError(err)
	}

	if n == 0 {
		return nil, jrm1.NewRpcErrorByUser(RpcErrorCode_ThreadIsNotFound, RpcErrorMsg_ThreadIsNotFound, nil)
	}

	// Update the moved message.
	err = srv.dbo.SetMessageThreadById(messageId, newThreadId, userRoles.User.GetUserParameters().GetId())
	if err != nil {
		return nil, srv.databaseError(err)
	}

	// Update the new link.
	var messagesR *ul.UidList
	messagesR, err = srv.dbo.GetThreadMessagesById(newThreadId)
	if err != nil {
		return nil, srv.databaseError(err)
	}

	err = messagesR.AddItem(messageId, false)
	if err != nil {
		srv.logError(err)
		return nil, jrm1.NewRpcErrorByUser(server2.RpcErrorCode_UidList, fmt.Sprintf(server2.RpcErrorMsgF_UidList, err.Error()), nil)
	}

	err = srv.dbo.SetThreadMessagesById(newThreadId, messagesR)
	if err != nil {
		return nil, srv.databaseError(err)
	}

	// Update the old link.
	var messagesL *ul.UidList
	messagesL, err = srv.dbo.GetThreadMessagesById(oldParent)
	if err != nil {
		return nil, srv.databaseError(err)
	}

	err = messagesL.RemoveItem(messageId)
	if err != nil {
		srv.logError(err)
		return nil, jrm1.NewRpcErrorByUser(server2.RpcErrorCode_UidList, fmt.Sprintf(server2.RpcErrorMsgF_UidList, err.Error()), nil)
	}

	err = srv.dbo.SetThreadMessagesById(oldParent, messagesL)
	if err != nil {
		return nil, srv.databaseError(err)
	}

	return initialMessage, nil
}

// deleteMessageH is a helper function used by other functions to remove a
// message.
func (srv *Server) deleteMessageH(messageId base2.Id) (initialMessage derived2.IMessage, re *jrm1.RpcError) {
	srv.dbo.LockForWriting()
	defer srv.dbo.UnlockAfterWriting()

	// Read the message.
	var err error
	initialMessage, err = srv.dbo.GetMessageById(messageId)
	if err != nil {
		return nil, srv.databaseError(err)
	}

	if initialMessage == nil {
		return nil, jrm1.NewRpcErrorByUser(RpcErrorCode_MessageIsNotFound, RpcErrorMsg_MessageIsNotFound, nil)
	}

	// Update the link.
	var linkMessages *ul.UidList
	linkMessages, err = srv.dbo.GetThreadMessagesById(initialMessage.GetThreadId())
	if err != nil {
		return nil, srv.databaseError(err)
	}

	err = linkMessages.RemoveItem(messageId)
	if err != nil {
		srv.logError(err)
		return nil, jrm1.NewRpcErrorByUser(server2.RpcErrorCode_UidList, fmt.Sprintf(server2.RpcErrorMsgF_UidList, err.Error()), nil)
	}

	err = srv.dbo.SetThreadMessagesById(initialMessage.GetThreadId(), linkMessages)
	if err != nil {
		return nil, srv.databaseError(err)
	}

	// Delete the message.
	err = srv.dbo.DeleteMessageById(messageId)
	if err != nil {
		return nil, srv.databaseError(err)
	}

	return initialMessage, nil
}

// reportSystemEvent reports the system event to the notification module.
func (srv *Server) reportSystemEvent(se derived2.ISystemEvent) (re *jrm1.RpcError) {
	if se == nil {
		return jrm1.NewRpcErrorByUser(server2.RpcErrorCode_SystemEvent, server2.RpcErrorMsg_SystemEvent, nil)
	}

	params := nm.ProcessSystemEventSParams{
		DKeyParams: rpc3.DKeyParams{
			// DKey is set during module start-up, so it is non-null.
			DKey: *srv.dKeyForNM,
		},
		SystemEventData: se.GetSystemEventData(),
	}
	result := new(nm.ProcessSystemEventSResult)
	var err error
	re, err = srv.nmServiceClient.MakeRequest(context.Background(), nc.FuncProcessSystemEventS, params, result)
	if err != nil {
		srv.logError(err)
		return jrm1.NewRpcErrorByUser(server2.RpcErrorCode_RPCCall, server2.RpcErrorMsg_RPCCall, nil)
	}
	if re != nil {
		return re
	}

	return nil
}
