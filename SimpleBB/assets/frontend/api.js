// Settings.
settingsPath = "settings.json";
rootPath = "/";
adminPage = "/admin";
redirectDelay = 3;

// Pages and Query Parameters.
Qp = {
	ChangeEmailStep1: "?changeEmail1",
	ChangeEmailStep2: "?changeEmail2",
	ChangeEmailStep3: "?changeEmail3",
	ChangePwdStep1: "?changePwd1",
	ChangePwdStep2: "?changePwd2",
	ChangePwdStep3: "?changePwd3",
	LogInStep1: "?login1",
	LogInStep2: "?login2",
	LogInStep3: "?login3",
	LogInStep4: "?login4",
	LogOutStep1: "?logout1",
	LogOutStep2: "?logout2",
	Prefix: "?",
	RegistrationStep1: "?reg1",
	RegistrationStep2: "?reg2",
	RegistrationStep3: "?reg3",
	RegistrationStep4: "?reg4",
	SelfPage: "?selfPage",
}

Qpn = {
	Forum: "forum",
	Id: "id",
	ListOfLoggedUsers: "listOfLoggedUsers",
	ListOfResources: "listOfResources",
	ListOfUsers: "listOfUsers",
	ManagerOfSections: "manageSections",
	ManagerOfForums: "manageForums",
	ManagerOfThreads: "managerOfThreads",
	ManagerOfMessages: "managerOfMessages",
	ManagerOfNotifications: "managerOfNotifications",
	ManagerOfResources: "managerOfResources",
	Message: "message",
	Notifications: "notifications",
	Page: "page",
	RegistrationsReadyForApproval: "registrationsReadyForApproval",
	Section: "section",
	SubscriptionsPage: "subscriptions",
	Thread: "thread",
	UserPage: "userPage",
}

// Action names.
ActionName = {
	AddForum: "addForum",
	AddMessage: "addMessage",
	AddNotification: "addNotification",
	AddResource: "addResource",
	AddSection: "addSection",
	AddSubscription: "addSubscription",
	AddThread: "addThread",
	ApproveAndRegisterUser: "approveAndRegisterUser",
	BanUser: "banUser",
	ChangeEmail: "changeEmail",
	ChangeForumName: "changeForumName",
	ChangeForumSection: "changeForumSection",
	ChangeMessageText: "changeMessageText",
	ChangeMessageThread: "changeMessageThread",
	ChangePwd: "changePassword",
	ChangeSectionName: "changeSectionName",
	ChangeSectionParent: "changeSectionParent",
	ChangeThreadName: "changeThreadName",
	ChangeThreadForum: "changeThreadForum",
	CountSelfSubscriptions: "countSelfSubscriptions",
	CountUnreadNotifications: "countUnreadNotifications",
	DeleteForum: "deleteForum",
	DeleteMessage: "deleteMessage",
	DeleteNotification: "deleteNotification",
	DeleteResource: "deleteResource",
	DeleteSection: "deleteSection",
	DeleteSelfSubscription: "deleteSelfSubscription",
	DeleteThread: "deleteThread",
	GetLatestMessageOfThread: "getLatestMessageOfThread",
	GetListOfAllResourcesOnPage: "getListOfAllResourcesOnPage",
	GetListOfAllUsers: "getListOfAllUsers",
	GetListOfAllUsersOnPage: "getListOfAllUsersOnPage",
	GetListOfLoggedUsers: "getListOfLoggedUsers",
	GetListOfLoggedUsersOnPage: "getListOfLoggedUsersOnPage",
	GetListOfRegistrationsReadyForApproval: "getListOfRegistrationsReadyForApproval",
	GetMessage: "getMessage",
	GetNotifications: "getNotifications",
	GetNotificationsOnPage: "getNotificationsOnPage",
	GetResource: "getResource",
	GetResourceValue: "getResourceValue",
	GetSelfRoles: "getSelfRoles",
	GetSelfSubscriptions: "getSelfSubscriptions",
	GetSelfSubscriptionsOnPage: "getSelfSubscriptionsOnPage",
	GetThread: "getThread",
	GetThreadNamesByIds: "getThreadNamesByIds",
	GetUserName: "getUserName",
	GetUserSession: "getUserSession",
	IsSelfSubscribed: "isSelfSubscribed",
	IsUserLoggedIn: "isUserLoggedIn",
	ListForumAndThreadsOnPage: "listForumAndThreadsOnPage",
	ListSectionsAndForums: "listSectionsAndForums",
	ListThreadAndMessagesOnPage: "listThreadAndMessagesOnPage",
	LogUserIn: "logUserIn",
	LogUserOut: "logUserOut",
	LogUserOutA: "logUserOutA",
	MoveForumDown: "moveForumDown",
	MoveForumUp: "moveForumUp",
	MoveSectionDown: "moveSectionDown",
	MoveSectionUp: "moveSectionUp",
	MoveThreadDown: "moveThreadDown",
	MoveThreadUp: "moveThreadUp",
	MarkNotificationAsRead: "markNotificationAsRead",
	RegisterUser: "registerUser",
	RejectRegistrationRequest: "rejectRegistrationRequest",
	SetUserRoleAuthor: "setUserRoleAuthor",
	SetUserRoleReader: "setUserRoleReader",
	SetUserRoleWriter: "setUserRoleWriter",
	UnbanUser: "unbanUser",
	ViewUserParameters: "viewUserParameters",
}

// Errors.
Err = {
	ActionMismatch: "action mismatch",
	Client: "client error",
	DuplicateMapKey: "duplicate map key",
	ElementTypeUnsupported: "unsupported element type",
	IdNotSet: "ID is not set",
	IdNotFound: "ID is not found",
	MessageNotFound: "message is not found",
	NameIsNotSet: "name is not set",
	NextPageDoesNotExist: "next page does not exist",
	NextStepUnknown: "unknown next step",
	NotOk: "something went wrong",
	PageNotSet: "page is not set",
	PageNotFound: "page is not found",
	ParentIsNotSet: "parent is not set",
	PasswordNotValid: "password is not valid",
	PreviousPageDoesNotExist: "previous page does not exist",
	ResourceType: "unsupported resource type",
	RootSectionNotFound: "root section is not found",
	RpcError: "RPC error",
	SectionNotFound: "section is not found",
	Server: "server error",
	Settings: "settings error",
	TextIsNotSet: "text is not set",
	TypeIsNotSet: "type is not set",
	ThreadNotFound: "thread is not found",
	Unknown: "unknown error",
	UnknownVariant: "unknown variant",
	ValueNotSet: "value is not set",
	WebTokenIsNotSet: "web token is not set",
}

// Messages.
Msg = {
	GenericErrorPrefix: "Error: ",
	Redirecting: "Redirecting. Please wait ...",
}

// User role names.
UserRole = {
	Author: "author",
	Writer: "writer",
	Reader: "reader",
	Logging: "logging",
}

class ApiRequest {
	constructor(action, parameters) {
		this.Action = action;
		this.Parameters = parameters;
	}
}

class ApiResponse {
	constructor(isOk, jsonObject, statusCode, errorText) {
		this.IsOk = isOk;
		this.JsonObject = jsonObject;
		this.StatusCode = statusCode;
		this.ErrorText = errorText;
	}
}

// Request parameters.

class Parameters_AddForum {
	constructor(parent, name) {
		this.SectionId = parent;
		this.Name = name;
	}
}

class Parameters_AddMessage {
	constructor(parent, text) {
		this.ThreadId = parent;
		this.Text = text;
	}
}

class Parameters_AddNotification {
	constructor(userId, text) {
		this.UserId = userId;
		this.Text = text;
	}
}

class Parameters_AddResource {
	constructor(resource) {
		this.Resource = resource;
	}
}

class Parameters_AddSection {
	constructor(parent, name) {
		this.Parent = parent;
		this.Name = name;
	}
}

class Parameters_AddSubscription {
	constructor(threadId, userId) {
		this.ThreadId = threadId;
		this.UserId = userId;
	}
}

class Parameters_AddThread {
	constructor(parent, name) {
		this.ForumId = parent;
		this.Name = name;
	}
}

class Parameters_ApproveAndRegisterUser {
	constructor(email) {
		this.Email = email;
	}
}

class Parameters_BanUser {
	constructor(userId) {
		this.UserId = userId;
	}
}

class Parameters_ChangeEmail1 {
	constructor(stepN, newEmail) {
		this.StepN = stepN;
		this.NewEmail = newEmail;
	}
}

class Parameters_ChangeEmail2 {
	constructor(stepN, requestId, authChallengeResponse, verificationCodeOld, verificationCodeNew, captchaAnswer) {
		this.StepN = stepN;
		this.RequestId = requestId;
		this.AuthChallengeResponse = authChallengeResponse;
		this.VerificationCodeOld = verificationCodeOld;
		this.VerificationCodeNew = verificationCodeNew;
		this.CaptchaAnswer = captchaAnswer;
	}
}

class Parameters_ChangeForumName {
	constructor(forumId, name) {
		this.ForumId = forumId;
		this.Name = name;
	}
}

class Parameters_ChangeForumSection {
	constructor(forumId, sectionId) {
		this.ForumId = forumId;
		this.SectionId = sectionId; // Parent.
	}
}

class Parameters_ChangeMessageText {
	constructor(messageId, text) {
		this.MessageId = messageId;
		this.Text = text;
	}
}

class Parameters_ChangeMessageThread {
	constructor(messageId, threadId) {
		this.MessageId = messageId;
		this.ThreadId = threadId; // Parent.
	}
}

class Parameters_ChangePwd1 {
	constructor(stepN, newPassword) {
		this.StepN = stepN;
		this.NewPassword = newPassword;
	}
}

class Parameters_ChangePwd2 {
	constructor(stepN, requestId, authChallengeResponse, verificationCode, captchaAnswer) {
		this.StepN = stepN;
		this.RequestId = requestId;
		this.AuthChallengeResponse = authChallengeResponse;
		this.VerificationCode = verificationCode;
		this.CaptchaAnswer = captchaAnswer;
	}
}

class Parameters_ChangeSectionName {
	constructor(sectionId, name) {
		this.SectionId = sectionId;
		this.Name = name;
	}
}

class Parameters_ChangeSectionParent {
	constructor(sectionId, parent) {
		this.SectionId = sectionId;
		this.Parent = parent;
	}
}

class Parameters_ChangeThreadForum {
	constructor(threadId, forumId) {
		this.ThreadId = threadId;
		this.ForumId = forumId; // Parent.
	}
}

class Parameters_ChangeThreadName {
	constructor(threadId, name) {
		this.ThreadId = threadId;
		this.Name = name;
	}
}

class Parameters_DeleteForum {
	constructor(forumId) {
		this.ForumId = forumId;
	}
}

class Parameters_DeleteMessage {
	constructor(messageId) {
		this.MessageId = messageId;
	}
}

class Parameters_DeleteNotification {
	constructor(notificationId) {
		this.NotificationId = notificationId;
	}
}

class Parameters_DeleteResource {
	constructor(resourceId) {
		this.ResourceId = resourceId;
	}
}

class Parameters_DeleteSection {
	constructor(sectionId) {
		this.SectionId = sectionId;
	}
}

class Parameters_DeleteSelfSubscription {
	constructor(threadId) {
		this.ThreadId = threadId;
	}
}

class Parameters_DeleteThread {
	constructor(threadId) {
		this.ThreadId = threadId;
	}
}

class Parameters_GetLatestMessageOfThread {
	constructor(threadId) {
		this.ThreadId = threadId;
	}
}

class Parameters_GetListOfAllResourcesOnPage {
	constructor(page) {
		this.Page = page;
	}
}

class Parameters_GetListOfAllUsersOnPage {
	constructor(page) {
		this.Page = page;
	}
}

class Parameters_GetListOfLoggedUsersOnPage {
	constructor(page) {
		this.Page = page;
	}
}

class Parameters_GetListOfRegistrationsReadyForApproval {
	constructor(page) {
		this.Page = page;
	}
}

class Parameters_GetMessage {
	constructor(messageId) {
		this.MessageId = messageId;
	}
}

class Parameters_GetNotificationsOnPage {
	constructor(page) {
		this.Page = page;
	}
}

class Parameters_GetResource {
	constructor(resourceId) {
		this.ResourceId = resourceId;
	}
}

class Parameters_GetResourceValue {
	constructor(resourceId) {
		this.ResourceId = resourceId;
	}
}

class Parameters_GetSelfSubscriptionsOnPage {
	constructor(page) {
		this.Page = page;
	}
}

class Parameters_GetThread {
	constructor(threadId) {
		this.ThreadId = threadId;
	}
}

class Parameters_GetThreadNamesByIds {
	constructor(threadIds) {
		this.ThreadIds = threadIds;
	}
}

class Parameters_GetUserName {
	constructor(userId) {
		this.UserId = userId;
	}
}

class Parameters_GetUserSession {
	constructor(userId) {
		this.UserId = userId;
	}
}

class Parameters_IsSelfSubscribed {
	constructor(threadId) {
		this.ThreadId = threadId;
	}
}

class Parameters_IsUserLoggedIn {
	constructor(userId) {
		this.UserId = userId;
	}
}

class Parameters_ListForumAndThreadsOnPage {
	constructor(forumId, page) {
		this.ForumId = forumId;
		this.Page = page;
	}
}

class Parameters_ListThreadAndMessagesOnPage {
	constructor(threadId, page) {
		this.ThreadId = threadId;
		this.Page = page;
	}
}

class Parameters_LogIn1 {
	constructor(stepN, email) {
		this.StepN = stepN;
		this.Email = email;
	}
}

class Parameters_LogIn2 {
	constructor(stepN, email, requestId, captchaAnswer, authChallengeResponse) {
		this.StepN = stepN;
		this.Email = email;
		this.RequestId = requestId;
		this.CaptchaAnswer = captchaAnswer;
		this.AuthChallengeResponse = authChallengeResponse;
	}
}

class Parameters_LogIn3 {
	constructor(stepN, email, requestId, verificationCode) {
		this.StepN = stepN;
		this.Email = email;
		this.RequestId = requestId;
		this.VerificationCode = verificationCode;
	}
}

class Parameters_LogUserOutA {
	constructor(userId) {
		this.UserId = userId;
	}
}

class Parameters_MarkNotificationAsRead {
	constructor(notificationId) {
		this.NotificationId = notificationId;
	}
}

class Parameters_MoveForumUp {
	constructor(forumId) {
		this.ForumId = forumId;
	}
}

class Parameters_MoveForumDown {
	constructor(forumId) {
		this.ForumId = forumId;
	}
}

class Parameters_MoveSectionUp {
	constructor(sectionId) {
		this.SectionId = sectionId;
	}
}

class Parameters_MoveSectionDown {
	constructor(sectionId) {
		this.SectionId = sectionId;
	}
}

class Parameters_MoveThreadUp {
	constructor(threadId) {
		this.ThreadId = threadId;
	}
}

class Parameters_MoveThreadDown {
	constructor(threadId) {
		this.ThreadId = threadId;
	}
}

class Parameters_RegisterUser1 {
	constructor(stepN, email) {
		this.StepN = stepN;
		this.Email = email;
	}
}

class Parameters_RegisterUser2 {
	constructor(stepN, email, verificationCode) {
		this.StepN = stepN;
		this.Email = email;
		this.VerificationCode = verificationCode;
	}
}

class Parameters_RegisterUser3 {
	constructor(stepN, email, verificationCode, name, pwd) {
		this.StepN = stepN;
		this.Email = email;
		this.VerificationCode = verificationCode;
		this.Name = name;
		this.Password = pwd;
	}
}

class Parameters_RejectRegistrationRequest {
	constructor(registrationRequestId) {
		this.RegistrationRequestId = registrationRequestId;
	}
}

class Parameters_SetUserRoleAuthor {
	constructor(userId, isRoleEnabled) {
		this.UserId = userId;
		this.IsRoleEnabled = isRoleEnabled;
	}
}

class Parameters_SetUserRoleReader {
	constructor(userId, isRoleEnabled) {
		this.UserId = userId;
		this.IsRoleEnabled = isRoleEnabled;
	}
}

class Parameters_SetUserRoleWriter {
	constructor(userId, isRoleEnabled) {
		this.UserId = userId;
		this.IsRoleEnabled = isRoleEnabled;
	}
}

class Parameters_UnbanUser {
	constructor(userId) {
		this.UserId = userId;
	}
}

class Parameters_ViewUserParameters {
	constructor(userId) {
		this.UserId = userId;
	}
}

// API methods.

async function sendApiRequest(data) {
	let settings = getSettings();
	let url = rootPath + settings.ApiFolder;
	let ri = {
		method: "POST",
		body: JSON.stringify(data)
	};
	let resp = await fetch(url, ri);
	let result;
	if (resp.status === 200) {
		result = new ApiResponse(true, await resp.json(), resp.status, null);
		return result;
	} else {
		result = new ApiResponse(false, null, resp.status, await resp.text());
		if (result.ErrorText.length === 0) {
			result.ErrorText = createErrorTextByStatusCode(result.StatusCode);
		}
		console.error(result.ErrorText);
		return result;
	}
}

async function sendApiRequestAndGetResult(reqData) {
	let actionName = reqData.Action;

	let resp = await sendApiRequest(reqData);
	if (!resp.IsOk) {
		console.error(composeErrorText(resp.ErrorText));
		return null;
	}

	let jo = resp.JsonObject;
	if (jo == null) {
		console.error(composeErrorText(Err.RpcError));
		return null;
	}

	if (jo.action !== actionName) {
		console.error(composeErrorText(Err.ActionMismatch));
		return null;
	}

	return jo.result;
}

async function addForum(parent, name) {
	let params = new Parameters_AddForum(parent, name);
	let reqData = new ApiRequest(ActionName.AddForum, params);
	return await sendApiRequestAndGetResult(reqData);
}

async function addMessage(parent, text) {
	let params = new Parameters_AddMessage(parent, text);
	let reqData = new ApiRequest(ActionName.AddMessage, params);
	return await sendApiRequestAndGetResult(reqData);
}

async function addNotification(userId, text) {
	let params = new Parameters_AddNotification(userId, text);
	let reqData = new ApiRequest(ActionName.AddNotification, params);
	return await sendApiRequestAndGetResult(reqData);
}

async function addResource(resource) {
	let params = new Parameters_AddResource(resource);
	let reqData = new ApiRequest(ActionName.AddResource, params);
	return await sendApiRequestAndGetResult(reqData);
}

async function addSection(parent, name) {
	let params = new Parameters_AddSection(parent, name);
	let reqData = new ApiRequest(ActionName.AddSection, params);
	return await sendApiRequestAndGetResult(reqData);
}

async function addSubscription(threadId, userId) {
	let params = new Parameters_AddSubscription(threadId, userId);
	let reqData = new ApiRequest(ActionName.AddSubscription, params);
	return await sendApiRequestAndGetResult(reqData);
}

async function addThread(parent, name) {
	let params = new Parameters_AddThread(parent, name);
	let reqData = new ApiRequest(ActionName.AddThread, params);
	return await sendApiRequestAndGetResult(reqData);
}

async function approveAndRegisterUser(email) {
	let params = new Parameters_ApproveAndRegisterUser(email);
	let reqData = new ApiRequest(ActionName.ApproveAndRegisterUser, params);
	return await sendApiRequestAndGetResult(reqData);
}

async function banUser(userId) {
	let params = new Parameters_BanUser(userId);
	let reqData = new ApiRequest(ActionName.BanUser, params);
	return await sendApiRequestAndGetResult(reqData);
}

async function changeEmail1(stepN, newEmail) {
	let params = new Parameters_ChangeEmail1(stepN, newEmail);
	let reqData = new ApiRequest(ActionName.ChangeEmail, params);
	return await sendApiRequestAndGetResult(reqData);
}

async function changeEmail2(stepN, requestId, authChallengeResponse, verificationCodeOld, verificationCodeNew, captchaAnswer) {
	let params = new Parameters_ChangeEmail2(stepN, requestId, authChallengeResponse, verificationCodeOld, verificationCodeNew, captchaAnswer);
	let reqData = new ApiRequest(ActionName.ChangeEmail, params);
	return await sendApiRequestAndGetResult(reqData);
}

async function changeForumName(forumId, name) {
	let params = new Parameters_ChangeForumName(forumId, name);
	let reqData = new ApiRequest(ActionName.ChangeForumName, params);
	return await sendApiRequestAndGetResult(reqData);
}

async function changeForumSection(forumId, newParent) {
	let params = new Parameters_ChangeForumSection(forumId, newParent);
	let reqData = new ApiRequest(ActionName.ChangeForumSection, params);
	return await sendApiRequestAndGetResult(reqData);
}

async function changeMessageText(messageId, text) {
	let params = new Parameters_ChangeMessageText(messageId, text);
	let reqData = new ApiRequest(ActionName.ChangeMessageText, params);
	return await sendApiRequestAndGetResult(reqData);
}

async function changeMessageThread(messageId, newParent) {
	let params = new Parameters_ChangeMessageThread(messageId, newParent);
	let reqData = new ApiRequest(ActionName.ChangeMessageThread, params);
	return await sendApiRequestAndGetResult(reqData);
}

async function changePwd1(stepN, newPassword) {
	let params = new Parameters_ChangePwd1(stepN, newPassword);
	let reqData = new ApiRequest(ActionName.ChangePwd, params);
	return await sendApiRequestAndGetResult(reqData);
}

async function changePwd2(stepN, requestId, authChallengeResponse, verificationCode, captchaAnswer) {
	let params = new Parameters_ChangePwd2(stepN, requestId, authChallengeResponse, verificationCode, captchaAnswer);
	let reqData = new ApiRequest(ActionName.ChangePwd, params);
	return await sendApiRequestAndGetResult(reqData);
}

async function changeSectionName(sectionId, name) {
	let params = new Parameters_ChangeSectionName(sectionId, name);
	let reqData = new ApiRequest(ActionName.ChangeSectionName, params);
	return await sendApiRequestAndGetResult(reqData);
}

async function changeSectionParent(sectionId, newParent) {
	let params = new Parameters_ChangeSectionParent(sectionId, newParent);
	let reqData = new ApiRequest(ActionName.ChangeSectionParent, params);
	return await sendApiRequestAndGetResult(reqData);
}

async function changeThreadForum(threadId, newParent) {
	let params = new Parameters_ChangeThreadForum(threadId, newParent);
	let reqData = new ApiRequest(ActionName.ChangeThreadForum, params);
	return await sendApiRequestAndGetResult(reqData);
}

async function changeThreadName(threadId, name) {
	let params = new Parameters_ChangeThreadName(threadId, name);
	let reqData = new ApiRequest(ActionName.ChangeThreadName, params);
	return await sendApiRequestAndGetResult(reqData);
}

async function countSelfSubscriptions() {
	let reqData = new ApiRequest(ActionName.CountSelfSubscriptions, {});
	return await sendApiRequestAndGetResult(reqData);
}

async function countUnreadNotifications() {
	let reqData = new ApiRequest(ActionName.CountUnreadNotifications, {});
	return await sendApiRequestAndGetResult(reqData);
}

async function deleteForum(forumId) {
	let params = new Parameters_DeleteForum(forumId);
	let reqData = new ApiRequest(ActionName.DeleteForum, params);
	return await sendApiRequestAndGetResult(reqData);
}

async function deleteMessage(messageId) {
	let params = new Parameters_DeleteMessage(messageId);
	let reqData = new ApiRequest(ActionName.DeleteMessage, params);
	return await sendApiRequestAndGetResult(reqData);
}

async function deleteNotification(notificationId) {
	let params = new Parameters_DeleteNotification(notificationId);
	let reqData = new ApiRequest(ActionName.DeleteNotification, params);
	return await sendApiRequestAndGetResult(reqData);
}

async function deleteResource(resourceId) {
	let params = new Parameters_DeleteResource(resourceId);
	let reqData = new ApiRequest(ActionName.DeleteResource, params);
	return await sendApiRequestAndGetResult(reqData);
}

async function deleteSection(sectionId) {
	let params = new Parameters_DeleteSection(sectionId);
	let reqData = new ApiRequest(ActionName.DeleteSection, params);
	return await sendApiRequestAndGetResult(reqData);
}

async function deleteSelfSubscription(threadId) {
	let params = new Parameters_DeleteSelfSubscription(threadId);
	let reqData = new ApiRequest(ActionName.DeleteSelfSubscription, params);
	return await sendApiRequestAndGetResult(reqData);
}

async function deleteThread(threadId) {
	let params = new Parameters_DeleteThread(threadId);
	let reqData = new ApiRequest(ActionName.DeleteThread, params);
	return await sendApiRequestAndGetResult(reqData);
}

async function getLatestMessageOfThread(threadId) {
	let params = new Parameters_GetLatestMessageOfThread(threadId);
	let reqData = new ApiRequest(ActionName.GetLatestMessageOfThread, params);
	return await sendApiRequestAndGetResult(reqData);
}

async function getListOfAllResourcesOnPage(pageN) {
	let params = new Parameters_GetListOfAllResourcesOnPage(pageN);
	let reqData = new ApiRequest(ActionName.GetListOfAllResourcesOnPage, params);
	return await sendApiRequestAndGetResult(reqData);
}

async function getListOfAllUsers() {
	let reqData = new ApiRequest(ActionName.GetListOfAllUsers, {});
	return await sendApiRequestAndGetResult(reqData);
}

async function getListOfAllUsersOnPage(pageN) {
	let params = new Parameters_GetListOfAllUsersOnPage(pageN);
	let reqData = new ApiRequest(ActionName.GetListOfAllUsersOnPage, params);
	return await sendApiRequestAndGetResult(reqData);
}

async function getListOfLoggedUsers() {
	let reqData = new ApiRequest(ActionName.GetListOfLoggedUsers, {});
	return await sendApiRequestAndGetResult(reqData);
}

async function getListOfLoggedUsersOnPage(pageN) {
	let params = new Parameters_GetListOfLoggedUsersOnPage(pageN);
	let reqData = new ApiRequest(ActionName.GetListOfLoggedUsersOnPage, params);
	return await sendApiRequestAndGetResult(reqData);
}

async function getListOfRegistrationsReadyForApproval(pageN) {
	let params = new Parameters_GetListOfRegistrationsReadyForApproval(pageN);
	let reqData = new ApiRequest(ActionName.GetListOfRegistrationsReadyForApproval, params);
	return await sendApiRequestAndGetResult(reqData);
}

async function getMessage(messageId) {
	let params = new Parameters_GetMessage(messageId);
	let reqData = new ApiRequest(ActionName.GetMessage, params);
	return await sendApiRequestAndGetResult(reqData);
}

async function getNotifications() {
	let reqData = new ApiRequest(ActionName.GetNotifications, {});
	return await sendApiRequestAndGetResult(reqData);
}

async function getNotificationsOnPage(page) {
	let params = new Parameters_GetNotificationsOnPage(page);
	let reqData = new ApiRequest(ActionName.GetNotificationsOnPage, params);
	return await sendApiRequestAndGetResult(reqData);
}

async function getResource(resourceId) {
	let params = new Parameters_GetResource(resourceId);
	let reqData = new ApiRequest(ActionName.GetResource, params);
	return await sendApiRequestAndGetResult(reqData);
}

async function getResourceValue(resourceId) {
	let params = new Parameters_GetResourceValue(resourceId);
	let reqData = new ApiRequest(ActionName.GetResourceValue, params);
	return await sendApiRequestAndGetResult(reqData);
}

async function getSelfRoles() {
	let reqData = new ApiRequest(ActionName.GetSelfRoles, {});
	return await sendApiRequestAndGetResult(reqData);
}

async function getSelfSubscriptions() {
	let reqData = new ApiRequest(ActionName.GetSelfSubscriptions, {});
	return await sendApiRequestAndGetResult(reqData);
}

async function getSelfSubscriptionsOnPage(page) {
	let params = new Parameters_GetSelfSubscriptionsOnPage(page);
	let reqData = new ApiRequest(ActionName.GetSelfSubscriptionsOnPage, params);
	return await sendApiRequestAndGetResult(reqData);
}

async function getThread(threadId) {
	let params = new Parameters_GetThread(threadId);
	let reqData = new ApiRequest(ActionName.GetThread, params);
	return await sendApiRequestAndGetResult(reqData);
}

async function getThreadNamesByIds(threadIds) {
	let params = new Parameters_GetThreadNamesByIds(threadIds);
	let reqData = new ApiRequest(ActionName.GetThreadNamesByIds, params);
	return await sendApiRequestAndGetResult(reqData);
}

async function getUserName(userId) {
	let params = new Parameters_GetUserName(userId);
	let reqData = new ApiRequest(ActionName.GetUserName, params);
	return await sendApiRequestAndGetResult(reqData);
}

async function getUserSession(userId) {
	let params = new Parameters_GetUserSession(userId);
	let reqData = new ApiRequest(ActionName.GetUserSession, params);
	return await sendApiRequestAndGetResult(reqData);
}

async function isSelfSubscribed(threadId) {
	let params = new Parameters_IsSelfSubscribed(threadId);
	let reqData = new ApiRequest(ActionName.IsSelfSubscribed, params);
	return await sendApiRequestAndGetResult(reqData);
}

async function isUserLoggedIn(userId) {
	let params = new Parameters_IsUserLoggedIn(userId);
	let reqData = new ApiRequest(ActionName.IsUserLoggedIn, params);
	return await sendApiRequestAndGetResult(reqData);
}

async function listForumAndThreadsOnPage(forumId, page) {
	let params = new Parameters_ListForumAndThreadsOnPage(forumId, page);
	let reqData = new ApiRequest(ActionName.ListForumAndThreadsOnPage, params);
	return await sendApiRequestAndGetResult(reqData);
}

async function listSectionsAndForums() {
	let reqData = new ApiRequest(ActionName.ListSectionsAndForums, {});
	return await sendApiRequestAndGetResult(reqData);
}

async function listThreadAndMessagesOnPage(threadId, page) {
	let params = new Parameters_ListThreadAndMessagesOnPage(threadId, page);
	let reqData = new ApiRequest(ActionName.ListThreadAndMessagesOnPage, params);
	return await sendApiRequestAndGetResult(reqData);
}

async function logUserIn1(stepN, email) {
	let params = new Parameters_LogIn1(stepN, email);
	let reqData = new ApiRequest(ActionName.LogUserIn, params);
	return await sendApiRequestAndGetResult(reqData);
}

async function logUserIn2(stepN, email, requestId, captchaAnswer, authChallengeResponse) {
	let params = new Parameters_LogIn2(stepN, email, requestId, captchaAnswer, authChallengeResponse);
	let reqData = new ApiRequest(ActionName.LogUserIn, params);
	return await sendApiRequestAndGetResult(reqData);
}

async function logUserIn3(stepN, email, requestId, verificationCode) {
	let params = new Parameters_LogIn3(stepN, email, requestId, verificationCode);
	let reqData = new ApiRequest(ActionName.LogUserIn, params);
	return await sendApiRequestAndGetResult(reqData);
}

async function logUserOut1() {
	let reqData = new ApiRequest(ActionName.LogUserOut, {});
	return await sendApiRequestAndGetResult(reqData);
}

async function logUserOutA(userId) {
	let params = new Parameters_LogUserOutA(userId);
	let reqData = new ApiRequest(ActionName.LogUserOutA, params);
	return await sendApiRequestAndGetResult(reqData);
}

async function markNotificationAsRead(notificationId) {
	let params = new Parameters_MarkNotificationAsRead(notificationId);
	let reqData = new ApiRequest(ActionName.MarkNotificationAsRead, params);
	return await sendApiRequestAndGetResult(reqData);
}

async function moveForumUp(forumId) {
	let params = new Parameters_MoveForumUp(forumId);
	let reqData = new ApiRequest(ActionName.MoveForumUp, params);
	return await sendApiRequestAndGetResult(reqData);
}

async function moveForumDown(forumId) {
	let params = new Parameters_MoveForumDown(forumId);
	let reqData = new ApiRequest(ActionName.MoveForumDown, params);
	return await sendApiRequestAndGetResult(reqData);
}

async function moveSectionUp(sectionId) {
	let params = new Parameters_MoveSectionUp(sectionId);
	let reqData = new ApiRequest(ActionName.MoveSectionUp, params);
	return await sendApiRequestAndGetResult(reqData);
}

async function moveSectionDown(sectionId) {
	let params = new Parameters_MoveSectionDown(sectionId);
	let reqData = new ApiRequest(ActionName.MoveSectionDown, params);
	return await sendApiRequestAndGetResult(reqData);
}

async function moveThreadUp(threadId) {
	let params = new Parameters_MoveThreadUp(threadId);
	let reqData = new ApiRequest(ActionName.MoveThreadUp, params);
	return await sendApiRequestAndGetResult(reqData);
}

async function moveThreadDown(threadId) {
	let params = new Parameters_MoveThreadDown(threadId);
	let reqData = new ApiRequest(ActionName.MoveThreadDown, params);
	return await sendApiRequestAndGetResult(reqData);
}

async function registerUser1(stepN, email) {
	let params = new Parameters_RegisterUser1(stepN, email);
	let reqData = new ApiRequest(ActionName.RegisterUser, params);
	return await sendApiRequestAndGetResult(reqData);
}

async function registerUser2(stepN, email, verificationCode) {
	let params = new Parameters_RegisterUser2(stepN, email, verificationCode);
	let reqData = new ApiRequest(ActionName.RegisterUser, params);
	return await sendApiRequestAndGetResult(reqData);
}

async function registerUser3(stepN, email, verificationCode, name, pwd) {
	let params = new Parameters_RegisterUser3(stepN, email, verificationCode, name, pwd);
	let reqData = new ApiRequest(ActionName.RegisterUser, params);
	return await sendApiRequestAndGetResult(reqData);
}

async function rejectRegistrationRequest(registrationRequestId) {
	let params = new Parameters_RejectRegistrationRequest(registrationRequestId);
	let reqData = new ApiRequest(ActionName.RejectRegistrationRequest, params);
	return await sendApiRequestAndGetResult(reqData);
}

async function setUserRoleAuthor(userId, roleValue) {
	let params = new Parameters_SetUserRoleAuthor(userId, roleValue);
	let reqData = new ApiRequest(ActionName.SetUserRoleAuthor, params);
	return await sendApiRequestAndGetResult(reqData);
}

async function setUserRoleReader(userId, roleValue) {
	let params = new Parameters_SetUserRoleReader(userId, roleValue);
	let reqData = new ApiRequest(ActionName.SetUserRoleReader, params);
	return await sendApiRequestAndGetResult(reqData);
}

async function setUserRoleWriter(userId, roleValue) {
	let params = new Parameters_SetUserRoleWriter(userId, roleValue);
	let reqData = new ApiRequest(ActionName.SetUserRoleWriter, params);
	return await sendApiRequestAndGetResult(reqData);
}

async function unbanUser(userId) {
	let params = new Parameters_UnbanUser(userId);
	let reqData = new ApiRequest(ActionName.UnbanUser, params);
	return await sendApiRequestAndGetResult(reqData);
}

async function viewUserParameters(userId) {
	let params = new Parameters_ViewUserParameters(userId);
	let reqData = new ApiRequest(ActionName.ViewUserParameters, params);
	return await sendApiRequestAndGetResult(reqData);
}

// Various API helpers.

async function getSelfSubscriptionsWithThreadNamesPaginated(pageNumber) {
	let resA = await getSelfSubscriptionsOnPage(pageNumber);
	let us = resA.userSubscriptions;
	let userId = us.subscriber;
	let threadIds = us.subscriptions;

	let threadNames = [];
	if (threadIds.length > 0) {
		let resB = await getThreadNamesByIds(threadIds);
		threadNames = resB.threadNames;
	}

	let pageData = jsonToPageData(resA.userSubscriptions.pageData);

	return new SubscriptionsWithThreadNamesOnPage(userId, threadIds, threadNames, pageData);
}

// Models.

class EventParameters {
	constructor(userId, time) {
		this.UserId = userId;
		this.Time = time;
	}
}

class Forum {
	constructor(id, sectionId, name, threads, creator, editor) {
		this.Id = id;
		this.SectionId = sectionId;
		this.Name = name;
		this.Threads = threads;
		this.Creator = creator;
		this.Editor = editor;
	}
}

class Message {
	constructor(id, threadId, text, textChecksum, creator, editor) {
		this.Id = id;
		this.ThreadId = threadId;
		this.Text = text;
		this.TextChecksum = textChecksum;
		this.Creator = creator;
		this.Editor = editor;
	}

	getLastTouchTime() {
		if (this.Editor.Time == null) {
			return new Date(this.Creator.Time);
		}
		return new Date(this.Editor.Time);
	}

	getMaxEditTime(settings) {
		let lastTouchTime = this.getLastTouchTime();
		return addTimeSec(lastTouchTime, settings.MessageEditTime);
	}
}

class Notification {
	constructor(id, userId, text, toc, isRead, tor) {
		this.Id = id;
		this.UserId = userId;
		this.Text = text;
		this.ToC = toc;
		this.IsRead = isRead;
		this.ToR = tor;
	}
}

class OptionalEventParameters {
	constructor(userId, time) {
		this.UserId = userId;
		this.Time = time;
	}
}

class PageData {
	constructor(pageNumber, totalPages, pageSize, itemsOnPage, totalItems) {
		this.PageNumber = pageNumber;
		this.TotalPages = totalPages;
		this.PageSize = pageSize;
		this.ItemsOnPage = itemsOnPage;
		this.TotalItems = totalItems;
	}
}

class Password {
	constructor(pwd) {
		this.Pwd = pwd;
	}

	check() {
		if (this.Pwd.length < 16) {
			return false;
		}

		if ((this.Pwd.length % 4) !== 0) {
			return false;
		}

		let symbol;
		for (let i = 0; i < this.Pwd.length; i++) {
			symbol = this.Pwd.charAt(i);
			if (!this.checkSymbol(symbol)) {
				return false;
			}
		}

		return true
	}

	checkSymbol(symbol) {
		switch (symbol) {
			case '0':
			case '1':
			case '2':
			case '3':
			case '4':
			case '5':
			case '6':
			case '7':
			case '8':
			case '9':
				return true;

			case 'A':
			case 'B':
			case 'C':
			case 'D':
			case 'E':
			case 'F':
			case 'G':
			case 'H':
			case 'I':
			case 'J':
			case 'K':
			case 'L':
			case 'M':
			case 'N':
			case 'O':
			case 'P':
			case 'Q':
			case 'R':
			case 'S':
			case 'T':
			case 'U':
			case 'V':
			case 'W':
			case 'X':
			case 'Y':
			case 'Z':
				return true;

			case ' ':
			case '!':
			case '"':
			case '#':
			case '$':
			case '%':
			case '&':
			case "'":
			case '(':
			case ')':
			case '*':
			case '+':
			case ',':
			case '-':
			case '.':
			case '/':
			case ':':
			case ';':
			case '<':
			case '=':
			case '>':
			case '?':
			case '@':
			case '[':
			case "\\":
			case ']':
			case '^':
			case '_':
				return true;
		}

		return false;
	}
}

class Resource {
	constructor(id, type, text, number, toc) {
		this.Id = id;
		this.Type = type;
		this.Text = text;
		this.Number = number;
		this.ToC = toc;
	}

	checkType() {
		switch (this.Type) {
			case 1:
				return true;
			case 2:
				return true;
			default:
				return false;
		}
	}

	setValue(v) {
		switch (this.Type) {
			case 1:
				this.Text = String(v);
				return true;

			case 2:
				if (!isNumber(v) && !isNumeric(v)) {
					return false;
				}
				this.Number = Number(v);
				return true;

			default:
				console.error(Err.ResourceType, this.Type);
				return false;
		}
	}

	getValue() {
		switch (this.Type) {
			case 1:
				return this.Text

			case 2:
				return this.Number;

			default:
				console.error(Err.ResourceType, this.Type);
				return null;
		}
	}
}

class ResourceValue {
	constructor(id, value) {
		this.Id = id;
		this.Value = value;
	}
}

class RRFA {
	constructor(id, preRegTime, email, name) {
		this.Id = id;
		this.PreRegTime = preRegTime;
		this.Email = email;
		this.Name = name;
	}
}

class Section {
	constructor(id, parent, childType, children, name, creator, editor) {
		this.Id = id;
		this.Parent = parent;
		this.ChildType = childType;
		this.Children = children;
		this.Name = name;
		this.Creator = creator;
		this.Editor = editor;
	}
}

class SectionNode {
	constructor(section, level) {
		this.Section = section;
		this.Level = level;
	}
}

class Session {
	constructor(id, userId, startTime, userIPA) {
		this.Id = id;
		this.UserId = userId;
		this.StartTime = startTime;
		this.UserIPA = userIPA;
	}
}

class Settings {
	constructor(version, productVersion, siteName, siteDomain, captchaFolder, sessionMaxDuration, messageEditTime, pageSize, apiFolder, publicSettingsFileName, isFrontEndEnabled, frontEndStaticFilesFolder, notificationCountLimit) {
		this.Version = version;
		this.ProductVersion = productVersion;
		this.SiteName = siteName;
		this.SiteDomain = siteDomain;
		this.CaptchaFolder = captchaFolder;
		this.SessionMaxDuration = sessionMaxDuration;
		this.MessageEditTime = messageEditTime;
		this.PageSize = pageSize;
		this.ApiFolder = apiFolder;
		this.PublicSettingsFileName = publicSettingsFileName;
		this.IsFrontEndEnabled = isFrontEndEnabled;
		this.FrontEndStaticFilesFolder = frontEndStaticFilesFolder;
		this.NotificationCountLimit = notificationCountLimit;
	}

	check() {
		if ((this.PublicSettingsFileName !== settingsPath)) {
			console.error(Err.Settings);
			return false;
		}
		return true;
	}
}

class SubscriptionsWithThreadNamesOnPage {
	constructor(userId, threadIds, threadNames, pageData) {
		this.UserId = userId;
		this.ThreadIds = threadIds;
		this.ThreadNames = threadNames;
		this.PageData = pageData;
	}
}

class Thread {
	constructor(id, forumId, name, messages, creator, editor) {
		this.Id = id;
		this.ForumId = forumId;
		this.Name = name;
		this.Messages = messages;
		this.Creator = creator;
		this.Editor = editor;
	}
}

class User {
	constructor(id, preRegTime, email, name, approvalTime, regTime, roles, lastBadLogInTime, banTime, lastBadActionTime) {
		this.Id = id;
		this.PreRegTime = preRegTime;
		this.Email = email;
		this.Name = name;
		this.ApprovalTime = approvalTime;
		this.RegTime = regTime;
		this.Roles = roles;
		this.LastBadLogInTime = lastBadLogInTime;
		this.BanTime = banTime;
		this.LastBadActionTime = lastBadActionTime;
	}

	canAddMessage(latestMessageInThread) {
		if (!this.Roles.IsWriter) {
			return false;
		}

		if (latestMessageInThread == null) {
			return true;
		}

		if (latestMessageInThread.Creator.UserId !== this.Id) {
			return true;
		}

		let messageMaxEditTime = latestMessageInThread.getMaxEditTime(getSettings());
		if (Date.now() < messageMaxEditTime) {
			return false;
		}

		return true;
	}

	canEditMessage(message) {
		if (this.Roles.IsModerator) {
			return true;
		}

		if (!this.Roles.IsWriter) {
			return false;
		}

		if (this.Id !== message.Creator.UserId) {
			return false;
		}

		let messageMaxEditTime = message.getMaxEditTime(getSettings());
		if (Date.now() < messageMaxEditTime) {
			return true
		}

		return false;
	}
}

class UserRoles {
	constructor(isAdministrator, isModerator, isAuthor, isWriter, isReader, canLogIn) {
		this.IsAdministrator = isAdministrator;
		this.IsModerator = isModerator;
		this.IsAuthor = isAuthor;
		this.IsWriter = isWriter;
		this.IsReader = isReader;
		this.CanLogIn = canLogIn;
	}
}

// Methods related to models.

function jsonToEventParameters(x) {
	if (x == null) return null;
	return new EventParameters(x.userId, x.time);
}

function jsonToForum(x) {
	if (x == null) return null;
	let creator = jsonToEventParameters(x.creator);
	let editor = jsonToOptionalEventParameters(x.editor);
	return new Forum(x.id, x.sectionId, x.name, x.threads, creator, editor);
}

function jsonToForums(x) {
	if (x == null) return null;
	let fs = [];
	let f;
	for (let i = 0; i < x.length; i++) {
		f = jsonToForum(x[i]);
		fs.push(f);
	}
	return fs;
}

function jsonToMessage(x) {
	if (x == null) return null;
	let creator = jsonToEventParameters(x.creator);
	let editor = jsonToOptionalEventParameters(x.editor);
	return new Message(x.id, x.threadId, x.text, x.textChecksum, creator, editor);
}

function jsonToMessages(x) {
	if (x == null) return null;
	let ms = [];
	let m;
	for (let i = 0; i < x.length; i++) {
		m = jsonToMessage(x[i]);
		ms.push(m);
	}
	return ms;
}

function jsonToNotification(x) {
	if (x == null) return null;
	return new Notification(x.id, x.userId, x.text, x.toc, x.isRead, x.tor);
}

function jsonToNotifications(x) {
	if (x == null) return null;
	let ns = [];
	let n;
	for (let i = 0; i < x.length; i++) {
		n = jsonToNotification(x[i]);
		ns.push(n);
	}
	return ns;
}

function jsonToOptionalEventParameters(x) {
	if (x == null) return null;
	return new OptionalEventParameters(x.userId, x.time);
}

function jsonToPageData(x) {
	if (x == null) return null;
	return new PageData(x.pageNumber, x.totalPages, x.pageSize, x.itemsOnPage, x.totalItems);
}

function jsonToResource(x) {
	if (x == null) return null;
	return new Resource(x.id, x.type, x.text, x.number, x.toc);
}

function jsonToResourceValue(x) {
	if (x == null) return null;
	return new ResourceValue(x.id, x.value);
}

function jsonToRrfa(x) {
	if (x == null) return null;
	return new RRFA(x.id, x.preRegTime, x.email, x.name);
}

function jsonToRrfas(x) {
	if (x == null) return null;
	let rs = [];
	let r;
	for (let i = 0; i < x.length; i++) {
		r = jsonToRrfa(x[i]);
		rs.push(r);
	}
	return rs;
}

function jsonToSection(x) {
	if (x == null) return null;
	let creator = jsonToEventParameters(x.creator);
	let editor = jsonToOptionalEventParameters(x.editor);
	return new Section(x.id, x.parent, x.childType, x.children, x.name, creator, editor);
}

function jsonToSections(x) {
	if (x == null) return null;
	let ss = [];
	let s;
	for (let i = 0; i < x.length; i++) {
		s = jsonToSection(x[i]);
		ss.push(s);
	}
	return ss;
}

function jsonToSession(x) {
	if (x == null) return null;
	return new Session(x.id, x.userId, x.startTime, x.userIPA);
}

function jsonToThread(x) {
	if (x == null) return null;
	let creator = jsonToEventParameters(x.creator);
	let editor = jsonToOptionalEventParameters(x.editor);
	return new Thread(x.id, x.forumId, x.name, x.messages, creator, editor);
}

function jsonToThreads(x) {
	if (x == null) return null;
	let ts = [];
	let t;
	for (let i = 0; i < x.length; i++) {
		t = jsonToThread(x[i]);
		ts.push(t);
	}
	return ts;
}

function jsonToUser(x) {
	if (x == null) return null;
	let userRoles = jsonToUserRoles(x.roles);
	return new User(x.id, x.preRegTime, x.email, x.name, x.approvalTime, x.regTime, userRoles, x.lastBadLogInTime, x.banTime, x.lastBadActionTime);
}

function jsonToUserRoles(x) {
	if (x == null) return null;
	return new UserRoles(x.isAdministrator, x.isModerator, x.isAuthor, x.isWriter, x.isReader, x.canLogIn);
}

function putArrayItemsIntoMap(a) {
	let m = new Map();
	if (a == null) {
		return m;
	}

	let key;
	for (let i = 0; i < a.length; i++) {
		key = a[i].Id;
		if (m.has(key)) {
			console.error(Err.DuplicateMapKey);
			return null;
		}
		m.set(key, a[i]);
	}
	return m;
}

function getRootSectionIdx(sections) {
	for (let i = 0; i < sections.length; i++) {
		if (sections[i].Parent == null) {
			return i;
		}
		return null;
	}
}

// createTreeOfSections creates a tree of sections.
// 'nodes' is an output parameter.
function createTreeOfSections(section, sectionsMap, level, nodes) {
	if (section == null) {
		return;
	}

	nodes.push(new SectionNode(section, level));

	if (section.ChildType !== SectionChildType.Section) {
		return;
	}

	let subSectionIds = section.Children;
	let subSection;
	level++;
	for (let i = 0; i < subSectionIds.length; i++) {
		subSection = sectionsMap.get(subSectionIds[i]);
		createTreeOfSections(subSection, sectionsMap, level, nodes);
	}
}

function findCurrentNodeLevel(allNodes, sectionId) {
	let node;
	for (let i = 0; i < allNodes.length; i++) {
		node = allNodes[i];
		if (node.Section.Id === sectionId) {
			return node.Level;
		}
	}
}

function composeNotificationShortText(fullText) {
	let wordsCountMax = 8;
	let segmenter = new Intl.Segmenter([], {granularity: 'word'});
	let segmentedText = segmenter.segment(fullText);
	let words = [...segmentedText].filter(s => s.isWordLike).map(s => s.segment);
	if (words.length <= wordsCountMax) {
		return fullText;
	}
	let shortTextWOP = words.slice(1, wordsCountMax + 1).join(" ");
	return shortTextWOP + " ...";
}

async function composeMessageHeaderText(message) {
	let url = composeUrlForMessage(message.Id);
	let creatorName = await getUserNameById(message.Creator.UserId);
	let txt = "<a href='" + url + "'>" + prettyTime(message.Creator.Time) + "</a>" +
		' by <span class="messageCreatorName">' + creatorName + '</span>';

	if (message.Editor.UserId != null) {
		let editorName = await getUserNameById(message.Editor.UserId);
		txt += ', edited by <span class="messageEditorName">' + editorName + '</span>' +
			' on <span class="messageEditorTime">' + prettyTime(message.Editor.Time) + '</span>';
	}

	return txt;
}

// Common methods for settings.

async function fetchSettings() {
	let data = await fetch(rootPath + settingsPath);
	return await data.json();
}

function jsonToSettings(x) {
	return new Settings(
		x.version,
		x.productVersion,
		x.siteName,
		x.siteDomain,
		x.captchaFolder,
		x.sessionMaxDuration,
		x.messageEditTime,
		x.pageSize,
		x.apiFolder,
		x.publicSettingsFileName,
		x.isFrontEndEnabled,
		x.frontEndStaticFilesFolder,
		x.notificationCountLimit,
	);
}

async function updateSettingsIfNeeded() {
	if (isSettingsUpdateNeeded()) {
		return await updateSettings();
	}
	return true;
}

async function updateSettings() {
	let resp = await fetchSettings();
	let s = jsonToSettings(resp);
	if (!s.check()) {
		return false;
	}
	console.info('Received settings. Version: ' + s.Version.toString() + ".");

	// Save the settings for future usage.
	saveSettings(s);
	return true;
}

// Various common functions.

function isNumber(x) {
	return typeof x === 'number';
}

function isNumeric(str) {
	if (typeof str != "string") {
		return false
	}

	return !isNaN(str) && !isNaN(parseFloat(str))
}

function booleanToString(b) {
	if (b === true) {
		return "Yes";
	}
	if (b === false) {
		return "No";
	}
	console.error("boolToText:", b);
	return null;
}

function stringToBoolean(s) {
	if (s == null) {
		return null;
	}

	let x = s.trim().toLowerCase();

	switch (x) {
		case "true":
			return true;

		case "false":
			return false;

		case "yes":
		case "1":
			return true;

		case "no":
		case "0":
			return false;

		default:
			return JSON.parse(x);
	}
}

function getCurrentTimestamp() {
	return Math.floor(Date.now() / 1000);
}

async function sleep(ms) {
	await new Promise(r => setTimeout(r, ms));
}

function addTimeSec(t, deltaSec) {
	return new Date(t.getTime() + deltaSec * 1000);
}

function createErrorTextByStatusCode(statusCode) {
	if ((statusCode >= 400) && (statusCode <= 499)) {
		return Msg.GenericErrorPrefix + Err.Client + " (" + statusCode.toString() + ")";
	}
	if ((statusCode >= 500) && (statusCode <= 599)) {
		return Msg.GenericErrorPrefix + Err.Server + " (" + statusCode.toString() + ")";
	}
	return Msg.GenericErrorPrefix + Err.Unknown + " (" + statusCode.toString() + ")";
}

function composeErrorText(errMsg) {
	return Msg.GenericErrorPrefix + errMsg.trim() + ".";
}

function prepareIdVariable(sp) {
	if (!sp.has(Qpn.Id)) {
		console.error(Err.IdNotSet);
		return false;
	}

	let xId = Number(sp.get(Qpn.Id));
	if (xId <= 0) {
		console.error(Err.IdNotFound);
		return false;
	}

	mca_gvc.Id = xId;
	return true;
}

function preparePageVariable(sp) {
	let pageNumber;
	if (!sp.has(Qpn.Page)) {
		pageNumber = 1;
	} else {
		pageNumber = Number(sp.get(Qpn.Page));
		if (pageNumber <= 0) {
			console.error(Err.PageNotFound);
			return false;
		}
	}

	mca_gvc.Page = pageNumber;
	return true;
}

function prettyTime(timeStr) {
	if (timeStr == null) {
		return "";
	}
	if (timeStr.length === 0) {
		return "";
	}

	let t = new Date(timeStr);
	let monthN = t.getUTCMonth() + 1; // Months in JavaScript start with 0 !

	return t.getUTCDate().toString().padStart(2, '0') + "." +
		monthN.toString().padStart(2, '0') + "." +
		t.getUTCFullYear().toString().padStart(4, '0') + " " +
		t.getUTCHours().toString().padStart(2, '0') + ":" +
		t.getUTCMinutes().toString().padStart(2, '0');
}

function escapeHtml(text) {
	let div = document.createElement('div');
	div.textContent = text;
	return div.innerHTML;
}

function processMessageText(msgText) {
	let txt = escapeHtml(msgText);
	txt = txt.replaceAll("\r\n", '<br>');
	txt = txt.replaceAll("\n", '<br>');
	txt = txt.replaceAll("\r", '<br>');
	return txt;
}

function preparePageNumber(pageCount) {
	// Repair the page count.
	if ((pageCount === undefined) || (pageCount === 0)) {
		pageCount = 1;
	}
	mca_gvc.Pages = pageCount;

	// Check page number for overflow.
	let pageNumber = mca_gvc.Page;
	if (pageNumber > pageCount) {
		console.error(Err.PageNotFound);
		return false;
	}

	return true;
}

// URL composition.

function composeCaptchaImageUrl(captchaId) {
	let settings = getSettings();
	let captchaPath = rootPath + settings.CaptchaFolder;
	return captchaPath + Qp.Prefix + Qpn.Id + "=" + captchaId;
}

function composeUrlForEntity(entityName, entityId) {
	return Qp.Prefix + entityName + "&" + Qpn.Id + "=" + entityId;
}

function composeUrlForEntityPage(entityName, entityId, page) {
	return Qp.Prefix + entityName + "&" + Qpn.Id + "=" + entityId + "&" + Qpn.Page + "=" + page;
}

function composeUrlForForum(forumId) {
	return composeUrlForEntity(Qpn.Forum, forumId);
}

function composeUrlForForumPage(forumId, page) {
	return composeUrlForEntityPage(Qpn.Forum, forumId, page);
}

function composeUrlForFuncPage(func, page) {
	return Qp.Prefix + func + "&" + Qpn.Page + "=" + page;
}

function composeUrlForMessage(messageId) {
	return composeUrlForEntity(Qpn.Message, messageId);
}

function composeUrlForNotificationsPage(page) {
	return composeUrlForFuncPage(Qpn.Notifications, page);
}

function composeUrlForSection(sectionId) {
	return composeUrlForEntity(Qpn.Section, sectionId);
}

function composeUrlForSubscriptionsPage(page) {
	return composeUrlForFuncPage(Qpn.SubscriptionsPage, page);
}

function composeUrlForThread(threadId) {
	return composeUrlForEntity(Qpn.Thread, threadId);
}

function composeUrlForThreadPage(threadId, page) {
	return composeUrlForEntityPage(Qpn.Thread, threadId, page);
}

function composeUrlForUserPageA(userId) {
	return Qp.Prefix + Qpn.UserPage + "&" + Qpn.Id + "=" + userId;
}

function composeUrlForAdminPageA(func, page) {
	return Qp.Prefix + func + "&" + Qpn.Page + "=" + page;
}

// Redirect & page reloading.

async function reloadPage(wait) {
	if (wait) {
		await sleep(redirectDelay * 1000);
	}
	location.reload();
}

// redirectToRelativePath redirects to a page with a relative path.
// E.g., if a relative path is 'x', then URL is '/x', where '/' is a front end
// root.
async function redirectToRelativePath(wait, relPath) {
	let url = rootPath + relPath;
	await redirectPage(wait, url);
}

async function redirectToMainPage(wait) {
	await redirectPage(wait, rootPath);
}

async function redirectPage(wait, url) {
	if (wait) {
		await sleep(redirectDelay * 1000);
	}

	document.location.href = url;
}

async function redirectToMainMenuA(wait) {
	let url = adminPage;
	await redirectPage(wait, url);
}

async function redirectToSubPageA(wait, qp) {
	let url = adminPage + qp;
	await redirectPage(wait, url);
}

// UI functions.

function addPaginator(el, pageNumber, pageCount, variantPrev, variantNext) {
	let div = document.createElement("DIV");
	let cn = PageZoneClass.Paginator;
	div.className = cn;
	div.id = cn;

	let s = document.createElement("span");
	s.textContent = "Page " + pageNumber + " of " + pageCount + " ";
	div.appendChild(s);

	let btnPrev = document.createElement("input");
	btnPrev.type = "button";
	btnPrev.className = ButtonClass.PaginatorPrev;
	btnPrev.id = ButtonClass.PaginatorPrev;
	btnPrev.value = ButtonName.PaginatorPrev;
	addClickEventHandler(btnPrev, variantPrev);
	div.appendChild(btnPrev);

	s = document.createElement("span");
	s.className = "spacerA";
	s.innerHTML = "&nbsp;";
	div.appendChild(s);

	let btnNext = document.createElement("input");
	btnNext.type = "button";
	btnNext.className = ButtonClass.PaginatorNext;
	btnNext.id = ButtonClass.PaginatorNext;
	btnNext.value = ButtonName.PaginatorNext;
	addClickEventHandler(btnNext, variantNext);
	div.appendChild(btnNext);

	el.appendChild(div);
}

function disableButton(btn) {
	switch (btn.tagName) {
		case "INPUT":
			btn.value = "";
			btn.disabled = true;
			btn.style.display = "none";
			return;

		default:
			console.error(Err.ElementTypeUnsupported);
	}
}

function disableParentForm(btn, pp, ignoreButton) {
	if (!ignoreButton) {
		btn.disabled = true;
	}

	let el;
	for (i = 0; i < pp.childNodes.length; i++) {
		let ch = pp.childNodes[i];
		for (j = 0; j < ch.childNodes.length; j++) {
			el = ch.childNodes[j];

			if (el !== btn) {
				el.disabled = true;
			} else {
				if (!ignoreButton) {
					el.disabled = true;
				}
			}
		}
	}
}

function showTable(tbl) {
	tbl.style.display = "table";
}

function showBlock(block) {
	block.style.display = "block";
}

function hideBlock(block) {
	block.style.display = "none";
}

function showActionSuccess(btn, txt) {
	let ppp = btn.parentNode.parentNode.parentNode;
	let d = document.createElement("DIV");
	d.className = "actionSuccess";
	d.textContent = txt;
	ppp.appendChild(d);
}

function addClickEventHandler(btn, ehVariant) {
	let en = "click";

	switch (ehVariant) {
		case EventHandlerVariant.ForumPagePrev:
			btn.addEventListener(en, async (e) => {
				await onBtnPrevClick_forumPage(btn);
			});
			return;

		case EventHandlerVariant.ForumPageNext:
			btn.addEventListener(en, async (e) => {
				await onBtnNextClick_forumPage(btn);
			});
			return;

		case EventHandlerVariant.ThreadPagePrev:
			btn.addEventListener(en, async (e) => {
				await onBtnPrevClick_threadPage(btn);
			});
			return;

		case EventHandlerVariant.ThreadPageNext:
			btn.addEventListener(en, async (e) => {
				await onBtnNextClick_threadPage(btn);
			});
			return;

		case EventHandlerVariant.SubscriptionsPrev:
			btn.addEventListener(en, async (e) => {
				await onBtnPrevClick_subscriptionsPage(btn);
			});
			return;

		case EventHandlerVariant.SubscriptionsNext:
			btn.addEventListener(en, async (e) => {
				await onBtnNextClick_subscriptionsPage(btn);
			});
			return;

		case EventHandlerVariant.NotificationsPagePrev:
			btn.addEventListener(en, async (e) => {
				await onBtnPrevClick_notificationsPage(btn);
			});
			return;

		case EventHandlerVariant.NotificationsPageNext:
			btn.addEventListener(en, async (e) => {
				await onBtnNextClick_notificationsPage(btn);
			});
			return;

		case EventHandlerVariant.UserListPrevA:
			btn.addEventListener(en, async (e) => {
				await onBtnPrevClick_userList(btn);
			});
			return;

		case EventHandlerVariant.UserListNextA:
			btn.addEventListener(en, async (e) => {
				await onBtnNextClick_userList(btn);
			});
			return;

		case EventHandlerVariant.LoggedUserListPrevA:
			btn.addEventListener(en, async (e) => {
				await onBtnPrevClick_logged(btn);
			});
			return;

		case EventHandlerVariant.LoggedUserListNextA:
			btn.addEventListener(en, async (e) => {
				await onBtnNextClick_logged(btn);
			});
			return;

		case EventHandlerVariant.RrfaListPrevA:
			btn.addEventListener(en, async (e) => {
				await onBtnPrevClick_rrfa(btn);
			});
			return;

		case EventHandlerVariant.RrfaListNextA:
			btn.addEventListener(en, async (e) => {
				await onBtnNextClick_rrfa(btn);
			});
			return;

		case EventHandlerVariant.ResourceListPrevA:
			btn.addEventListener(en, async (e) => {
				await onBtnPrevClick_resources(btn);
			});
			return;

		case EventHandlerVariant.ResourceListNextA:
			btn.addEventListener(en, async (e) => {
				await onBtnNextClick_resources(btn);
			});
			return;

		default:
			console.error(Err.UnknownVariant);
	}
}

function newDiv() {
	return document.createElement("DIV");
}

function newFieldset() {
	return document.createElement("FIELDSET");
}

function newTable() {
	return document.createElement("TABLE");
}

function newTr() {
	return document.createElement("TR");
}

function newTh() {
	return document.createElement("TH");
}

function newTd() {
	return document.createElement("TD");
}

function newInput() {
	return document.createElement("INPUT");
}
