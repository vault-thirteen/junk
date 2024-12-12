window.onpageshow = function (event) {
	if (event.persisted) {
		// Unfortunately, JavaScript does not reload a page when you click
		// "Go Back" button in your web browser. Every year programmers invent
		// a new "wheel" to fix this bug. And every year old working solutions
		// stop working and new ones are invented. This circus looks infinite,
		// but in reality it will end as soon as this evil programming language
		// dies. Please, do not support JavaScript and its developers in any
		// means possible. Please, let this evil "technology" to die.
		console.info("JavaScript must die. This pseudo language is a big mockery and ridicule of people. This is not a joke. This is truth.");
		window.location.reload();
	}
};

settingsExpirationDuration = 60;
userNameForUnknownUser = "Unknown user";

// Names of JavaScript storage variables.
Varname = {
	SettingsLoadTime: "settingsLoadTime",
	LogInTime: "logInTime",
	SettingsVersion: "settingsVersion",
	SettingsProductVersion: "settingsProductVersion",
	SettingsSiteName: "settingsSiteName",
	SettingsSiteDomain: "settingsSiteDomain",
	SettingsCaptchaFolder: "settingsCaptchaFolder",
	SettingsSessionMaxDuration: "settingsSessionMaxDuration",
	SettingsMessageEditTime: "settingsMessageEditTime",
	SettingsPageSize: "settingsPageSize",
	SettingsApiFolder: "settingsApiFolder",
	SettingsPublicSettingsFileName: "settingsPublicSettingsFileName",
	SettingsIsFrontEndEnabled: "settingsIsFrontEndEnabled",
	SettingsFrontEndStaticFilesFolder: "settingsFrontEndStaticFilesFolder",
	SettingsNotificationCountLimit: "settingsNotificationCountLimit",
	IsLoggedIn: "isLoggedIn",
	RegistrationEmail: "registrationEmail",
	RegistrationVcode: "registrationVcode",
	LogInEmail: "logInEmail",
	LogInRequestId: "logInRequestId",
	LogInAuthDataBytes: "logInAuthDataBytes",
	LogInIsCaptchaNeeded: "logInIsCaptchaNeeded",
	LogInCaptchaId: "logInCaptchaId",
	ChangeEmailRequestId: "changeEmailRequestId",
	ChangeEmailAuthDataBytes: "changeEmailAuthDataBytes",
	ChangeEmailIsCaptchaNeeded: "changeEmailIsCaptchaNeeded",
	ChangeEmailCaptchaId: "changeEmailCaptchaId",
	ChangePwdRequestId: "changePwdRequestId",
	ChangePwdAuthDataBytes: "changePwdAuthDataBytes",
	ChangePwdIsCaptchaNeeded: "changePwdIsCaptchaNeeded",
	ChangePwdCaptchaId: "changePwdCaptchaId",
}

// Form Input Elements.
Fi = {
	id1: "f1i",
	id2: "f2i",
	id3: "f3i",
	id4: "f4i",
	id4_errflag: "f4ief",
	id5: "f5i",
	id6: "f6i",
	id7: "f7i",
	id7_image: "f7ii",
	id8: "f8i",
	id9: "f9i",
	id10: "f10i",
	id11: "f11i",
	id11_image: "f11ii",
	id12: "f12i",
	id13: "f13i",
	id14: "f14i",
	id15: "f15i",
	id16: "f16i",
	id16_image: "f16ii",
	id17: "f17i",
	id18: "f18i",
	id19: "f19i",
	id20: "f20i",
	id21_tr: "f21tr",
	id22: "f22i",
}

// Section settings.
SectionChildType = {
	None: 0,
	Section: 1,
	Forum: 2,
}
SectionMarginDelta = 10;

ButtonName = {
	BackToRoot: "🏠",
	BackToSection: "Go back",
	BackToForum: "Go back",
	BackToThread: "Go back",
	StartNewThread: "Start a new Thread",
	SubscribeToThread: "Subscribe",
	AddMessage: "Add a Message",
	EditMessage: "Edit Message",
	Notifications: " ☼ ",
	Account: "Account",
	PaginatorPrev: "<",
	PaginatorNext: ">",
	ShowFullNotificationU: "Show",
	MarkNotificationAsReadU: "Read",
	DeleteNotificationU: "DEL",
	Unsubscribe: "Unsubscribe",
	ConfirmThreadStart: "Confirm",
	ConfirmMessageCreation: "Send Message",
	ConfirmMessageEdit: "Save Message",
}

ButtonClass = {
	BackToRoot: "btnGoToIndex",
	BackToSection: "btnGoToSection",
	BackToForum: "btnGoToForum",
	BackToThread: "btnGoToThread",
	StartNewThread: "btnStartNewThread",
	SubscribeToThread: "btnSubscribe",
	AddMessage: "btnAddMessage",
	EditMessage: "btnEditMessage",
	NotificationsBoxFull: "btnNotificationsBoxFull",
	NotificationsOn: "btnNotificationsOn",
	NotificationsOff: "btnNotificationsOff",
	Account: "btnAccount",
	PaginatorPrev: "btnPrev",
	PaginatorNext: "btnNext",
	ShowFullNotificationU: "btnShowFullNotificationU",
	MarkNotificationAsReadU: "btnMarkNotificationAsReadU",
	DeleteNotificationU: "btnDeleteNotificationU",
	Unsubscribe: "btnUnsubscribe",
	ConfirmThreadStart: "btnConfirmThreadStart",
	ConfirmMessageCreation: "btnConfirmMessageCreation",
	ConfirmMessageEdit: "btnConfirmMessageEdit",
}

PageZoneClass = {
	L: "L",
	R: "R",
	ActionPanel: "actionPanel",
	PageHead: "pageHead",
	PageTitle: "pageTitle",
	Paginator: "paginator",
	BottomActionPanel: "bottomActionPanel",
}

Hint = {
	NoNewNotifications: "All notifications are read.",
	NewNotifications: "New notifications are available.",
	NotificationBoxIsFull: "Notification box is full. You are not receiving new notifications.",
}

ObjectType = {
	Section: "Section",
	Forum: "Forum",
	Thread: "Thread",
	Message: "Message",
}

EventHandlerVariant = {
	ForumPagePrev: "ForumPagePrev",
	ForumPageNext: "ForumPageNext",
	ThreadPagePrev: "ThreadPagePrev",
	ThreadPageNext: "ThreadPageNext",
	SubscriptionsPrev: "SubscriptionsPrev",
	SubscriptionsNext: "SubscriptionsNext",
	NotificationsPagePrev: "NotificationsPagePrev",
	NotificationsPageNext: "NotificationsPageNext",
}

UserRole = {
	Administrator: "Administrator",
	Moderator: "Moderator",
	Author: "Author",
	Writer: "Writer",
	Reader: "Reader",
}

// Global variables.
class GlobalVariablesContainer {
	constructor(id, page, pages, unc) {
		this.Id = id;
		this.Page = page;
		this.Pages = pages;
		this.UNC = unc;
	}
}

class UserNameCache {
	constructor() {
		this.m = new Map();
	}

	async GetName(userId) {
		if (this.m.has(userId)) {
			return this.m.get(userId);
		}

		let res = await getUserName(userId);
		if (res == null) {
			return userNameForUnknownUser;
		}
		let user = jsonToUser(res.user);
		this.m.set(userId, user.Name);
		return user.Name;
	}
}

async function getUserNameById(userId) {
	return await mca_gvc.UNC.GetName(userId);
}

mca_gvc = new GlobalVariablesContainer(0, 0, 0, new UserNameCache());

// Settings.

function isSettingsUpdateNeeded() {
	let settingsLoadTimeStr = sessionStorage.getItem(Varname.SettingsLoadTime);
	if (settingsLoadTimeStr == null) {
		return true;
	}

	let timeNow = getCurrentTimestamp();
	let settingsAge = timeNow - Number(settingsLoadTimeStr);
	if (settingsAge >= settingsExpirationDuration) {
		return true;
	}

	return false;
}

function saveSettings(s) {
	sessionStorage.setItem(Varname.SettingsVersion, s.Version.toString());
	sessionStorage.setItem(Varname.SettingsProductVersion, s.ProductVersion);
	sessionStorage.setItem(Varname.SettingsSiteName, s.SiteName);
	sessionStorage.setItem(Varname.SettingsSiteDomain, s.SiteDomain);
	sessionStorage.setItem(Varname.SettingsCaptchaFolder, s.CaptchaFolder);
	sessionStorage.setItem(Varname.SettingsSessionMaxDuration, s.SessionMaxDuration.toString());
	sessionStorage.setItem(Varname.SettingsMessageEditTime, s.MessageEditTime.toString());
	sessionStorage.setItem(Varname.SettingsPageSize, s.PageSize.toString());
	sessionStorage.setItem(Varname.SettingsApiFolder, s.ApiFolder);
	sessionStorage.setItem(Varname.SettingsPublicSettingsFileName, s.PublicSettingsFileName);
	sessionStorage.setItem(Varname.SettingsIsFrontEndEnabled, s.IsFrontEndEnabled.toString());
	sessionStorage.setItem(Varname.SettingsFrontEndStaticFilesFolder, s.FrontEndStaticFilesFolder);
	sessionStorage.setItem(Varname.SettingsNotificationCountLimit, s.NotificationCountLimit.toString());

	let timeNow = getCurrentTimestamp();
	sessionStorage.setItem(Varname.SettingsLoadTime, timeNow.toString());
}

function getSettings() {
	let settingsLoadTime = sessionStorage.getItem(Varname.SettingsLoadTime);
	if (settingsLoadTime == null) {
		console.error(Err.Settings);
		return null;
	}

	return new Settings(
		sessionStorage.getItem(Varname.SettingsVersion),
		sessionStorage.getItem(Varname.SettingsProductVersion),
		sessionStorage.getItem(Varname.SettingsSiteName),
		sessionStorage.getItem(Varname.SettingsSiteDomain),
		sessionStorage.getItem(Varname.SettingsCaptchaFolder),
		sessionStorage.getItem(Varname.SettingsSessionMaxDuration),
		sessionStorage.getItem(Varname.SettingsMessageEditTime),
		sessionStorage.getItem(Varname.SettingsPageSize),
		sessionStorage.getItem(Varname.SettingsApiFolder),
		sessionStorage.getItem(Varname.SettingsPublicSettingsFileName),
		sessionStorage.getItem(Varname.SettingsIsFrontEndEnabled),
		sessionStorage.getItem(Varname.SettingsFrontEndStaticFilesFolder),
		sessionStorage.getItem(Varname.SettingsNotificationCountLimit),
	);
}

// Logged status.

function isLoggedIn(settings) {
	let isLoggedInStr = localStorage.getItem(Varname.IsLoggedIn);
	let isLoggedIn;

	if (isLoggedInStr === null) {
		isLoggedIn = false;
		localStorage.setItem(Varname.IsLoggedIn, isLoggedIn.toString());
		return false;
	}

	isLoggedIn = stringToBoolean(isLoggedInStr);
	if (!isLoggedIn) {
		return false;
	}

	// Check if the session is not closed by timeout.
	let logInTime = Number(localStorage.getItem(Varname.LogInTime));
	let timeNow = getCurrentTimestamp();
	let sessionAge = timeNow - logInTime;
	if (sessionAge > settings.SessionMaxDuration) {
		isLoggedIn = false;
		localStorage.setItem(Varname.IsLoggedIn, isLoggedIn.toString());
		return false;
	}

	return true;
}

function saveLogInStatus() {
	let isLoggedIn = true;
	localStorage.setItem(Varname.IsLoggedIn, isLoggedIn.toString());
	let timeNow = getCurrentTimestamp();
	localStorage.setItem(Varname.LogInTime, timeNow.toString());
}

function saveLogOutStatus() {
	let isLoggedIn = false;
	localStorage.setItem(Varname.IsLoggedIn, isLoggedIn.toString());
	localStorage.removeItem(Varname.LogInTime);
}

// Entry point.
async function onPageLoad() {
	// Settings initialisation.
	let ok = await updateSettingsIfNeeded();
	if (!ok) {
		return;
	}
	let settings = getSettings();

	let curPage = window.location.search;

	// Redirect to registration.
	switch (curPage) {
		case Qp.RegistrationStep1:
			await showReg1Form();
			return;

		case Qp.RegistrationStep2:
			await showReg2Form();
			return;

		case Qp.RegistrationStep3:
			await showReg3Form();
			return;

		case Qp.RegistrationStep4:
			await showReg4Form();
			return;
	}

	// Redirect to logging in.
	switch (curPage) {
		case Qp.LogInStep1:
			await showLogIn1Form();
			return;

		case Qp.LogInStep2:
			await showLogIn2Form();
			return;

		case Qp.LogInStep3:
			await showLogIn3Form();
			return;

		case Qp.LogInStep4:
			await showLogIn4Form();
			await redirectToMainPage(true);
			return;
	}

	let isLoggedInB = isLoggedIn(settings);
	if (!isLoggedInB) {
		await showLogIn1Form();
		return;
	}

	// Pages for logged users.
	switch (curPage) {
		case Qp.LogOutStep1:
			await showLogOut1Form();
			return;

		case Qp.LogOutStep2:
			await showLogOut2Form();
			return;

		case Qp.ChangeEmailStep1:
			await showChangeEmail1Form();
			return;

		case Qp.ChangeEmailStep2:
			await showChangeEmail2Form();
			return;

		case Qp.ChangeEmailStep3:
			await showChangeEmail3Form();
			return;

		case Qp.ChangePwdStep1:
			await showChangePwd1Form();
			return;

		case Qp.ChangePwdStep2:
			await showChangePwd2Form();
			return;

		case Qp.ChangePwdStep3:
			await showChangePwd3Form();
			return;

		case Qp.SelfPage:
			await showUserPage();
			return;
	}

	let sp = new URLSearchParams(curPage);

	// Notifications.
	if (sp.has(Qpn.Notifications)) {
		if (!preparePageVariable(sp)) {
			return;
		}
		await showPage_Notifications();
		return;
	}

	// Subscriptions.
	if (sp.has(Qpn.SubscriptionsPage)) {
		if (!preparePageVariable(sp)) {
			return;
		}
		await showPage_Subscriptions();
		return;
	}

	// Show the bulletin board.
	if (sp.has(Qpn.Section)) {
		if (!prepareIdVariable(sp)) {
			return;
		}
		await showSection();
		return;
	}

	if (sp.has(Qpn.Forum)) {
		if ((!prepareIdVariable(sp)) || (!preparePageVariable(sp))) {
			return;
		}
		await showForum();
		return;
	}

	if (sp.has(Qpn.Thread)) {
		if ((!prepareIdVariable(sp)) || (!preparePageVariable(sp))) {
			return;
		}
		await showThread();
		return;
	}

	if (sp.has(Qpn.Message)) {
		if ((!prepareIdVariable(sp))) {
			return;
		}
		await showMessage();
		return;
	}

	await showBB();
}

async function showReg1Form() {
	let settings = getSettings();

	// Draw.
	let p = document.getElementById("divReg1");
	showBlock(p);
	addActionPanel(p, true);
	await addPageHead(p, settings.SiteName, true);
}

async function showReg2Form() {
	let settings = getSettings();

	// Draw.
	let p = document.getElementById("divReg2");
	showBlock(p);
	addActionPanel(p, true);
	await addPageHead(p, settings.SiteName, true);
}

async function showReg3Form() {
	let settings = getSettings();

	// Draw.
	let p = document.getElementById("divReg3");
	showBlock(p);
	addActionPanel(p, true);
	await addPageHead(p, settings.SiteName, true);
}

async function showReg4Form() {
	let settings = getSettings();

	// Draw.
	let p = document.getElementById("divReg4");
	showBlock(p);
	addActionPanel(p, true);
	await addPageHead(p, settings.SiteName, true);
}

async function showLogIn1Form() {
	let settings = getSettings();

	// Draw.
	let p = document.getElementById("divLogIn1");
	showBlock(p);
	addActionPanel(p, true);
	await addPageHead(p, settings.SiteName, true);
}

async function showLogIn2Form() {
	// Captcha (optional).
	let isCaptchaNeeded = stringToBoolean(sessionStorage.getItem(Varname.LogInIsCaptchaNeeded));
	let captchaId = sessionStorage.getItem(Varname.LogInCaptchaId);
	let cptImageTr = document.getElementById("formHolderLogIn2CaptchaImage");
	let cptImage = document.getElementById(Fi.id7_image);
	let cptAnswerTr = document.getElementById("formHolderLogIn2CaptchaAnswer");
	let cptAnswer = document.getElementById(Fi.id7);
	let settings = getSettings();

	// Draw.
	let p = document.getElementById("divLogIn2");
	showBlock(p);
	addActionPanel(p, true);
	await addPageHead(p, settings.SiteName, true);
	setCaptchaInputsVisibility(isCaptchaNeeded, captchaId, cptImageTr, cptImage, cptAnswerTr, cptAnswer);
}

async function showLogIn3Form() {
	let settings = getSettings();

	// Draw.
	let p = document.getElementById("divLogIn3");
	showBlock(p);
	addActionPanel(p, true);
	await addPageHead(p, settings.SiteName, true);
}

async function showLogIn4Form() {
	let settings = getSettings();

	// Draw.
	let p = document.getElementById("divLogIn4");
	showBlock(p);
	addActionPanel(p, true);
	await addPageHead(p, settings.SiteName, true);
}

async function showLogOut1Form() {
	let settings = getSettings();

	// Draw.
	let p = document.getElementById("divLogOut1");
	showBlock(p);
	addActionPanel(p, true);
	await addPageHead(p, settings.SiteName, true);
}

async function showLogOut2Form() {
	let settings = getSettings();

	// Draw.
	let p = document.getElementById("divLogOut2");
	showBlock(p);
	addActionPanel(p, true);
	await addPageHead(p, settings.SiteName, true);
}

async function showChangeEmail1Form() {
	let settings = getSettings();

	// Draw.
	let p = document.getElementById("divChangeEmail1");
	showBlock(p);
	addActionPanel(p, true);
	await addPageHead(p, settings.SiteName, true);
}

async function showChangeEmail2Form() {
	// Captcha (optional).
	let isCaptchaNeeded = stringToBoolean(sessionStorage.getItem(Varname.ChangeEmailIsCaptchaNeeded));
	let captchaId = sessionStorage.getItem(Varname.ChangeEmailCaptchaId);
	let cptImageTr = document.getElementById("formHolderChangeEmail2CaptchaImage");
	let cptImage = document.getElementById(Fi.id11_image);
	let cptAnswerTr = document.getElementById("formHolderChangeEmail2CaptchaAnswer");
	let cptAnswer = document.getElementById(Fi.id11);
	let settings = getSettings();

	// Draw.
	let p = document.getElementById("divChangeEmail2");
	showBlock(p);
	addActionPanel(p, true);
	await addPageHead(p, settings.SiteName, true);
	setCaptchaInputsVisibility(isCaptchaNeeded, captchaId, cptImageTr, cptImage, cptAnswerTr, cptAnswer);
}

async function showChangeEmail3Form() {
	let settings = getSettings();

	// Draw.
	let p = document.getElementById("divChangeEmail3");
	showBlock(p);
	addActionPanel(p, true);
	await addPageHead(p, settings.SiteName, true);
}

async function showChangePwd1Form() {
	let settings = getSettings();

	// Draw.
	let p = document.getElementById("divChangePwd1");
	showBlock(p);
	addActionPanel(p, true);
	await addPageHead(p, settings.SiteName, true);
}

async function showChangePwd2Form() {
	// Captcha (optional).
	let isCaptchaNeeded = stringToBoolean(sessionStorage.getItem(Varname.ChangePwdIsCaptchaNeeded));
	let captchaId = sessionStorage.getItem(Varname.ChangePwdCaptchaId);
	let cptImageTr = document.getElementById("formHolderChangePwd2CaptchaImage");
	let cptImage = document.getElementById(Fi.id16_image);
	let cptAnswerTr = document.getElementById("formHolderChangePwd2CaptchaAnswer");
	let cptAnswer = document.getElementById(Fi.id16);
	let settings = getSettings();

	// Draw.
	let p = document.getElementById("divChangePwd2");
	showBlock(p);
	addActionPanel(p, true);
	await addPageHead(p, settings.SiteName, true);
	setCaptchaInputsVisibility(isCaptchaNeeded, captchaId, cptImageTr, cptImage, cptAnswerTr, cptAnswer);
}

async function showChangePwd3Form() {
	let settings = getSettings();

	// Draw.
	let p = document.getElementById("divChangePwd3");
	showBlock(p);
	addActionPanel(p, true);
	await addPageHead(p, settings.SiteName, true);
}

async function showUserPage() {
	let res = await getSelfRoles();
	let user = jsonToUser(res.user);
	let settings = getSettings();

	// Draw.
	let p = document.getElementById("divUserPage");
	showBlock(p);
	addActionPanel(p, true);
	await addPageHead(p, settings.SiteName, true);
	await drawUserPage(user);
}

async function showPage_Notifications() {
	let res = await getNotificationsOnPage(mca_gvc.Page);
	let pageCount = res.nop.pageData.totalPages;
	if (!preparePageNumber(pageCount)) {
		return
	}

	let notifications = jsonToNotifications(res.nop.notifications);
	let settings = getSettings();

	// Draw.
	let p = document.getElementById("divBB");
	showBlock(p);
	await addPageHead(p, settings.SiteName, false);
	addActionPanel(p, false);
	drawPageTitle(p, 'Notifications');
	addPaginator(p, mca_gvc.Page, mca_gvc.Pages, EventHandlerVariant.NotificationsPagePrev, EventHandlerVariant.NotificationsPageNext);
	drawNotifications(p, 'notificationList', notifications);
}

async function showPage_Subscriptions() {
	let swtnop = await getSelfSubscriptionsWithThreadNamesPaginated(mca_gvc.Page);
	let pageCount = swtnop.PageData.TotalPages;
	if (!preparePageNumber(pageCount)) {
		return
	}

	let settings = getSettings();

	// Draw.
	let p = document.getElementById("divBB");
	showBlock(p);
	await addPageHead(p, settings.SiteName, false);
	addActionPanel(p, false);
	drawPageTitle(p, 'Subscriptions');
	addPaginator(p, mca_gvc.Page, mca_gvc.Pages, EventHandlerVariant.SubscriptionsPrev, EventHandlerVariant.SubscriptionsNext);
	drawSubscriptions(p, 'subscriptionsList', swtnop);
}

async function showSection() {
	let sectionId = mca_gvc.Id;
	let res = await listSectionsAndForums();
	let sections = jsonToSections(res.saf.sections);
	let forums = jsonToForums(res.saf.forums);
	let rootSectionIdx = getRootSectionIdx(sections);
	if (rootSectionIdx == null) {
		console.error(Err.RootSectionNotFound);
	}
	let rootSection = sections[rootSectionIdx];
	let sectionsMap = putArrayItemsIntoMap(sections);
	if (sectionsMap == null) {
		return;
	}
	let forumsMap = putArrayItemsIntoMap(forums);
	if (forumsMap == null) {
		return;
	}
	let allNodes = [];
	createTreeOfSections(rootSection, sectionsMap, 1, allNodes);
	let nodes = [];
	if (!sectionsMap.has(sectionId)) {
		console.error(Err.SectionNotFound);
		return;
	}
	let curSection = sectionsMap.get(sectionId);
	let curLevel = findCurrentNodeLevel(allNodes, sectionId);
	createTreeOfSections(curSection, sectionsMap, curLevel, nodes);
	let settings = getSettings();

	// Draw.
	let p = document.getElementById("divBB");
	showBlock(p);
	await addPageHead(p, settings.SiteName, false);
	if (curSection.Parent != null) {
		addActionPanel(p, false, ObjectType.Section, curSection.Parent);
	} else {
		addActionPanel(p, false);
	}
	drawSectionsAndForums(p, nodes, forumsMap);
}

async function showForum() {
	let forumId = mca_gvc.Id;
	let res = await listForumAndThreadsOnPage(forumId, mca_gvc.Page);
	let pageCount = res.fatop.pageData.totalPages;
	if (!preparePageNumber(pageCount)) {
		return
	}

	let forum = jsonToForum(res.fatop.forum);
	if (forum.Id !== forumId) {
		return;
	}
	let threads = jsonToThreads(res.fatop.threads);
	let threadsMap = putArrayItemsIntoMap(threads);
	if (threadsMap == null) {
		return;
	}
	let settings = getSettings();

	// Draw.
	let p = document.getElementById("divBB");
	showBlock(p);
	await addPageHead(p, settings.SiteName, false);
	addActionPanel(p, false, ObjectType.Section, forum.SectionId);
	drawPageTitle(p, forum.Name);
	addPaginator(p, mca_gvc.Page, mca_gvc.Pages, EventHandlerVariant.ForumPagePrev, EventHandlerVariant.ForumPageNext);
	drawForumThreads(p, forum, threadsMap);
	await addBottomActionPanel(p, ObjectType.Forum, forum);
}

async function showThread() {
	let threadId = mca_gvc.Id;
	let res = await listThreadAndMessagesOnPage(threadId, mca_gvc.Page);
	let pageCount = res.tamop.pageData.totalPages;
	if (!preparePageNumber(pageCount)) {
		return
	}

	let thread = jsonToThread(res.tamop.thread);
	if (thread.Id !== threadId) {
		return;
	}
	let messages = jsonToMessages(res.tamop.messages);
	let messagesMap = putArrayItemsIntoMap(messages);
	if (messagesMap == null) {
		return;
	}
	let settings = getSettings();

	// Draw.
	let p = document.getElementById("divBB");
	showBlock(p);
	await addPageHead(p, settings.SiteName, false);
	addActionPanel(p, false, ObjectType.Forum, thread.ForumId);
	drawPageTitle(p, thread.Name);
	addPaginator(p, mca_gvc.Page, mca_gvc.Pages, EventHandlerVariant.ThreadPagePrev, EventHandlerVariant.ThreadPageNext);
	await drawThreadMessages(p, thread, messagesMap);
	await addBottomActionPanel(p, ObjectType.Thread, thread);
}

async function showMessage() {
	let messageId = mca_gvc.Id;
	let res = await getMessage(messageId);
	let message = jsonToMessage(res.message);
	let settings = getSettings();

	// Draw.
	let p = document.getElementById("divBB");
	showBlock(p);
	await addPageHead(p, settings.SiteName, false);
	addActionPanel(p, false, ObjectType.Thread, message.ThreadId);
	await drawMessage(p, message);
	await addBottomActionPanel(p, ObjectType.Message, message);
}

async function showBB() {
	let res = await listSectionsAndForums();
	let sections = jsonToSections(res.saf.sections);
	let forums = jsonToForums(res.saf.forums);
	let rootSectionIdx = getRootSectionIdx(sections);
	if (rootSectionIdx == null) {
		console.error(Err.RootSectionNotFound);
	}
	let rootSection = sections[rootSectionIdx];
	let sectionsMap = putArrayItemsIntoMap(sections);
	if (sectionsMap == null) {
		return;
	}
	let forumsMap = putArrayItemsIntoMap(forums);
	if (forumsMap == null) {
		return;
	}
	let nodes = [];
	createTreeOfSections(rootSection, sectionsMap, 1, nodes);
	let settings = getSettings();

	// Draw.
	let p = document.getElementById("divBB");
	showBlock(p);
	await addPageHead(p, settings.SiteName, false);
	addActionPanel(p, false);
	drawSectionsAndForums(p, nodes, forumsMap);
}

// Event handlers.

async function onReg1Submit(btn) {
	// Send the request.
	let h3Field = document.getElementById("header3TextReg1");
	let errField = document.getElementById("header4TextReg1");
	let email = document.getElementById(Fi.id1).value;
	let res = await registerUser1(1, email);
	let nextStep = res.nextStep;
	if (nextStep !== 2) {
		errField.innerHTML = composeErrorText(Err.NextStepUnknown);
		return;
	}
	errField.innerHTML = "";

	// Save some non-sensitive input data into browser for the next page.
	sessionStorage.setItem(Varname.RegistrationEmail, email);

	// Redirect to next step.
	disableButton(btn);
	h3Field.innerHTML = Msg.Redirecting;
	await redirectToRelativePath(true, Qp.RegistrationStep2);
}

async function onReg2Submit(btn) {
	// Send the request.
	let h3Field = document.getElementById("header3TextReg2");
	let errField = document.getElementById("header4TextReg2");
	let email = sessionStorage.getItem(Varname.RegistrationEmail);
	let vcode = document.getElementById(Fi.id2).value;
	let res = await registerUser2(2, email, vcode);
	let nextStep = res.nextStep;
	if (nextStep !== 3) {
		errField.innerHTML = composeErrorText(Err.NextStepUnknown);
		return;
	}
	errField.innerHTML = "";

	// Save some non-sensitive input data into browser for the next page.
	sessionStorage.setItem(Varname.RegistrationVcode, vcode);

	// Redirect to next step.
	disableButton(btn);
	h3Field.innerHTML = Msg.Redirecting;
	await redirectToRelativePath(true, Qp.RegistrationStep3);
}

async function onReg3Submit(btn) {
	// Check the input.
	let pwdStr = document.getElementById(Fi.id4).value;
	let pwd = new Password(pwdStr);
	let pwdErrFlag = document.getElementById(Fi.id4_errflag);
	if (pwd.check()) {
		pwdErrFlag.className = "flag_none";
	} else {
		pwdErrFlag.className = "flag_error";
		return;
	}

	// Send the request.
	let h3Field = document.getElementById("header3TextReg3");
	let errField = document.getElementById("header4TextReg3");
	let email = sessionStorage.getItem(Varname.RegistrationEmail);
	let vcode = sessionStorage.getItem(Varname.RegistrationVcode);
	let name = document.getElementById(Fi.id3).value;
	let res = await registerUser3(3, email, vcode, name, pwd.Pwd);
	let nextStep = res.nextStep;
	if ((nextStep !== 4) && (nextStep !== 0)) {
		errField.innerHTML = composeErrorText(Err.NextStepUnknown);
		return;
	}
	errField.innerHTML = "";

	// Clear saved input data from browser.
	sessionStorage.removeItem(Varname.RegistrationEmail);
	sessionStorage.removeItem(Varname.RegistrationVcode);

	// Redirect to next step.
	disableButton(btn);
	h3Field.innerHTML = Msg.Redirecting;
	await redirectToRelativePath(true, Qp.RegistrationStep4);
}

async function onLogIn1Submit(btn) {
	// Send the request.
	let h3Field = document.getElementById("header3TextLogIn1");
	let errField = document.getElementById("header4TextLogIn1");
	let email = document.getElementById(Fi.id5).value;
	let res = await logUserIn1(1, email);
	let nextStep = res.nextStep;
	if (nextStep !== 2) {
		errField.innerHTML = composeErrorText(Err.NextStepUnknown);
		return;
	}
	errField.innerHTML = "";

	// Save some non-sensitive input data into browser for the next page.
	sessionStorage.setItem(Varname.LogInEmail, email);
	let requestId = res.requestId;
	sessionStorage.setItem(Varname.LogInRequestId, requestId);
	let authDataBytes = res.authDataBytes;
	sessionStorage.setItem(Varname.LogInAuthDataBytes, authDataBytes);
	let isCaptchaNeeded = res.isCaptchaNeeded;
	sessionStorage.setItem(Varname.LogInIsCaptchaNeeded, isCaptchaNeeded.toString());
	let captchaId = res.captchaId;
	sessionStorage.setItem(Varname.LogInCaptchaId, captchaId);

	// Redirect to next step.
	disableButton(btn);
	h3Field.innerHTML = Msg.Redirecting;
	await redirectToRelativePath(true, Qp.LogInStep2);
}

async function onLogIn2Submit(btn) {
	let errField = document.getElementById("header4TextLogIn2");
	let h3Field = document.getElementById("header3TextLogIn2");

	// Captcha (optional).
	let captchaAnswer = document.getElementById(Fi.id7).value;

	// Secret.
	let authDataBytes = sessionStorage.getItem(Varname.LogInAuthDataBytes);
	let saltBA = base64ToByteArray(authDataBytes);
	let pwd = document.getElementById(Fi.id6).value;
	if (!isPasswordAllowed(pwd)) {
		errField.innerHTML = composeErrorText(Err.PasswordNotValid);
		return;
	}
	let keyBA = makeHashKey(pwd, saltBA);
	let authChallengeResponse = byteArrayToBase64(keyBA);

	// Send the request.
	let email = sessionStorage.getItem(Varname.LogInEmail);
	let requestId = sessionStorage.getItem(Varname.LogInRequestId);
	let res = await logUserIn2(2, email, requestId, captchaAnswer, authChallengeResponse);
	let nextStep = res.nextStep;
	if (nextStep !== 3) {
		errField.innerHTML = composeErrorText(Err.NextStepUnknown);
		return;
	}
	errField.innerHTML = "";

	// Clear saved input data from browser.
	sessionStorage.removeItem(Varname.LogInAuthDataBytes);
	sessionStorage.removeItem(Varname.LogInIsCaptchaNeeded);
	sessionStorage.removeItem(Varname.LogInCaptchaId);

	// Save some non-sensitive input data into browser for the next page.
	let newRequestId = res.requestId;
	sessionStorage.setItem(Varname.LogInRequestId, newRequestId);

	// Redirect to next step.
	disableButton(btn);
	h3Field.innerHTML = Msg.Redirecting;
	await redirectToRelativePath(true, Qp.LogInStep3);
}

async function onLogIn3Submit(btn) {
	let errField = document.getElementById("header4TextLogIn3");
	let h3Field = document.getElementById("header3TextLogIn3");

	// Send the request.
	let vcode = document.getElementById(Fi.id8).value;
	let email = sessionStorage.getItem(Varname.LogInEmail);
	let requestId = sessionStorage.getItem(Varname.LogInRequestId);
	let res = await logUserIn3(3, email, requestId, vcode);
	let nextStep = res.nextStep;
	let isWebTokenSet = res.isWebTokenSet;
	if (nextStep !== 0) {
		errField.innerHTML = composeErrorText(Err.NextStepUnknown);
		return;
	}
	if (!isWebTokenSet) {
		errField.innerHTML = composeErrorText(Err.WebTokenIsNotSet);
		return;
	}
	errField.innerHTML = "";

	// Clear saved input data from browser.
	sessionStorage.removeItem(Varname.LogInEmail);
	sessionStorage.removeItem(Varname.LogInRequestId);

	// Save the 'log' flag.
	saveLogInStatus();

	// Redirect to next step.
	disableButton(btn);
	h3Field.innerHTML = Msg.Redirecting;
	await redirectToMainPage(true);
}

async function onLogOut1Submit(btn) {
	let errField = document.getElementById("header4TextLogOut1");
	let h3Field = document.getElementById("header3TextLogOut1");
	let res = await logUserOut1();
	let ok = res.ok;
	if (!ok) {
		errField.innerHTML = composeErrorText(Err.NotOk);
		return;
	}
	errField.innerHTML = "";

	// Save the 'log' flag.
	saveLogOutStatus();

	// Redirect to next step.
	disableButton(btn);
	h3Field.innerHTML = Msg.Redirecting;
	await redirectToRelativePath(true, Qp.LogOutStep2);
}

async function onChangeEmail1Submit(btn) {
	// Send the request.
	let h3Field = document.getElementById("header3TextChangeEmail1");
	let errField = document.getElementById("header4TextChangeEmail1");
	let newEmail = document.getElementById(Fi.id9).value;
	let res = changeEmail1(1, newEmail);
	let nextStep = res.nextStep;
	if (nextStep !== 2) {
		errField.innerHTML = composeErrorText(Err.NextStepUnknown);
		return;
	}
	errField.innerHTML = "";

	// Save some non-sensitive input data into browser for the next page.
	let requestId = res.requestId;
	sessionStorage.setItem(Varname.ChangeEmailRequestId, requestId);
	let authDataBytes = res.authDataBytes;
	sessionStorage.setItem(Varname.ChangeEmailAuthDataBytes, authDataBytes);
	let isCaptchaNeeded = res.isCaptchaNeeded;
	sessionStorage.setItem(Varname.ChangeEmailIsCaptchaNeeded, isCaptchaNeeded.toString());
	let captchaId = res.captchaId;
	sessionStorage.setItem(Varname.ChangeEmailCaptchaId, captchaId);

	// Redirect to next step.
	disableButton(btn);
	h3Field.innerHTML = Msg.Redirecting;
	await redirectToRelativePath(true, Qp.ChangeEmailStep2);
}

async function onChangeEmail2Submit(btn) {
	let h3Field = document.getElementById("header3TextChangeEmail2");
	let errField = document.getElementById("header4TextChangeEmail2");

	// Captcha (optional).
	let captchaAnswer = document.getElementById(Fi.id11).value;

	// Secret.
	let authDataBytes = sessionStorage.getItem(Varname.ChangeEmailAuthDataBytes);
	let saltBA = base64ToByteArray(authDataBytes);
	let pwd = document.getElementById(Fi.id10).value;
	if (!isPasswordAllowed(pwd)) {
		errField.innerHTML = composeErrorText(Err.PasswordNotValid);
		return;
	}
	let keyBA = makeHashKey(pwd, saltBA);
	let authChallengeResponse = byteArrayToBase64(keyBA);

	// Send the request.
	let requestId = sessionStorage.getItem(Varname.ChangeEmailRequestId);
	let vCodeOld = document.getElementById(Fi.id12).value;
	let vCodeNew = document.getElementById(Fi.id13).value;
	let res = changeEmail2(2, requestId, authChallengeResponse, vCodeOld, vCodeNew, captchaAnswer);
	let nextStep = res.nextStep;
	if (nextStep !== 0) {
		errField.innerHTML = composeErrorText(Err.NextStepUnknown);
		return;
	}
	let ok = res.ok;
	if (!ok) {
		errField.innerHTML = composeErrorText(Err.NotOk);
		return;
	}
	errField.innerHTML = "";

	// Clear saved input data from browser.
	sessionStorage.removeItem(Varname.ChangeEmailRequestId);
	sessionStorage.removeItem(Varname.ChangeEmailAuthDataBytes);
	sessionStorage.removeItem(Varname.ChangeEmailIsCaptchaNeeded);
	sessionStorage.removeItem(Varname.ChangeEmailCaptchaId);

	// Save the 'log' flag.
	saveLogOutStatus();

	// Redirect to next step.
	disableButton(btn);
	h3Field.innerHTML = Msg.Redirecting;
	await redirectToRelativePath(true, Qp.ChangeEmailStep3);
}

async function onChangePwd1Submit(btn) {
	// Send the request.
	let h3Field = document.getElementById("header3TextChangePwd1");
	let errField = document.getElementById("header4TextChangePwd1");
	let newPwd = document.getElementById(Fi.id14).value;
	let res = await changePwd1(1, newPwd);
	let nextStep = res.nextStep;
	if (nextStep !== 2) {
		errField.innerHTML = composeErrorText(Err.NextStepUnknown);
		return;
	}
	errField.innerHTML = "";

	// Save some non-sensitive input data into browser for the next page.
	let requestId = res.requestId;
	sessionStorage.setItem(Varname.ChangePwdRequestId, requestId);
	let authDataBytes = res.authDataBytes;
	sessionStorage.setItem(Varname.ChangePwdAuthDataBytes, authDataBytes);
	let isCaptchaNeeded = res.isCaptchaNeeded;
	sessionStorage.setItem(Varname.ChangePwdIsCaptchaNeeded, isCaptchaNeeded.toString());
	let captchaId = res.captchaId;
	sessionStorage.setItem(Varname.ChangePwdCaptchaId, captchaId);

	// Redirect to next step.
	disableButton(btn);
	h3Field.innerHTML = Msg.Redirecting;
	await redirectToRelativePath(true, Qp.ChangePwdStep2);
}

async function onChangePwd2Submit(btn) {
	let h3Field = document.getElementById("header3TextChangePwd2");
	let errField = document.getElementById("header4TextChangePwd2");

	// Captcha (optional).
	let captchaAnswer = document.getElementById(Fi.id16).value;

	// Secret.
	let authDataBytes = sessionStorage.getItem(Varname.ChangePwdAuthDataBytes);
	let saltBA = base64ToByteArray(authDataBytes);
	let pwd = document.getElementById(Fi.id15).value;
	if (!isPasswordAllowed(pwd)) {
		errField.innerHTML = composeErrorText(Err.PasswordNotValid);
		return;
	}
	let keyBA = makeHashKey(pwd, saltBA);
	let authChallengeResponse = byteArrayToBase64(keyBA);

	// Send the request.
	let requestId = sessionStorage.getItem(Varname.ChangePwdRequestId);
	let vcode = document.getElementById(Fi.id17).value;
	let res = await changePwd2(2, requestId, authChallengeResponse, vcode, captchaAnswer);
	let nextStep = res.nextStep;
	if (nextStep !== 0) {
		errField.innerHTML = composeErrorText(Err.NextStepUnknown);
		return;
	}
	let ok = res.ok;
	if (!ok) {
		errField.innerHTML = composeErrorText(Err.NotOk);
		return;
	}
	errField.innerHTML = "";

	// Clear saved input data from browser.
	sessionStorage.removeItem(Varname.ChangePwdRequestId);
	sessionStorage.removeItem(Varname.ChangePwdAuthDataBytes);
	sessionStorage.removeItem(Varname.ChangePwdIsCaptchaNeeded);
	sessionStorage.removeItem(Varname.ChangePwdCaptchaId);

	// Save the 'log' flag.
	saveLogOutStatus();

	// Redirect to next step.
	disableButton(btn);
	h3Field.innerHTML = Msg.Redirecting;
	await redirectToRelativePath(true, Qp.ChangePwdStep3);
}

async function onBtnPrevClick_forumPage(btn) {
	if (mca_gvc.Page <= 1) {
		console.error(Err.PreviousPageDoesNotExist);
		return;
	}

	mca_gvc.Page--;
	let url = composeUrlForForumPage(mca_gvc.Id, mca_gvc.Page);
	await redirectPage(false, url);
}

async function onBtnNextClick_forumPage(btn) {
	if (mca_gvc.Page >= mca_gvc.Pages) {
		console.error(Err.NextPageDoesNotExist);
		return;
	}

	mca_gvc.Page++;
	let url = composeUrlForForumPage(mca_gvc.Id, mca_gvc.Page);
	await redirectPage(false, url);
}

async function onBtnPrevClick_threadPage(btn) {
	if (mca_gvc.Page <= 1) {
		console.error(Err.PreviousPageDoesNotExist);
		return;
	}

	mca_gvc.Page--;
	let url = composeUrlForThreadPage(mca_gvc.Id, mca_gvc.Page);
	await redirectPage(false, url);
}

async function onBtnNextClick_threadPage(btn) {
	if (mca_gvc.Page >= mca_gvc.Pages) {
		console.error(Err.NextPageDoesNotExist);
		return;
	}

	mca_gvc.Page++;
	let url = composeUrlForThreadPage(mca_gvc.Id, mca_gvc.Page);
	await redirectPage(false, url);
}

async function onBtnPrevClick_subscriptionsPage(btn) {
	if (mca_gvc.Page <= 1) {
		console.error(Err.PreviousPageDoesNotExist);
		return;
	}

	mca_gvc.Page--;
	let url = composeUrlForSubscriptionsPage(mca_gvc.Page);
	await redirectPage(false, url);
}

async function onBtnNextClick_subscriptionsPage(btn) {
	if (mca_gvc.Page >= mca_gvc.Pages) {
		console.error(Err.NextPageDoesNotExist);
		return;
	}

	mca_gvc.Page++;
	let url = composeUrlForSubscriptionsPage(mca_gvc.Page);
	await redirectPage(false, url);
}

async function onBtnPrevClick_notificationsPage(btn) {
	if (mca_gvc.Page <= 1) {
		console.error(Err.PreviousPageDoesNotExist);
		return;
	}

	mca_gvc.Page--;
	let url = composeUrlForNotificationsPage(mca_gvc.Page);
	await redirectPage(false, url);
}

async function onBtnNextClick_notificationsPage(btn) {
	if (mca_gvc.Page >= mca_gvc.Pages) {
		console.error(Err.NextPageDoesNotExist);
		return;
	}

	mca_gvc.Page++;
	let url = composeUrlForNotificationsPage(mca_gvc.Page);
	await redirectPage(false, url);
}

async function onBtnChangeEmailClick(btn) {
	let url = Qp.ChangeEmailStep1;
	await redirectPage(false, url);
}

async function onBtnChangePwdClick(btn) {
	let url = Qp.ChangePwdStep1;
	await redirectPage(false, url);
}

async function onBtnLogOutSelfClick(btn) {
	let url = Qp.LogOutStep1;
	await redirectPage(false, url);
}

async function onBtnManageSubscriptionsClick(btn) {
	let url = Qp.Prefix + Qpn.SubscriptionsPage;
	await redirectPage(false, url);
}

async function onBtnAccountClick(btn) {
	let url = Qp.SelfPage;
	await redirectPage(false, url);
}

async function onBtnNotificationsClick(btn) {
	let url = Qp.Prefix + Qpn.Notifications;
	await redirectPage(false, url);
}

async function onBtnGoToIndexClick(btn) {
	await redirectToMainPage();
}

async function onBtnGoToThreadClick(parentId) {
	let url = composeUrlForThread(parentId);
	await redirectPage(false, url);
}

async function onBtnGoToForumClick(parentId) {
	let url = composeUrlForForum(parentId);
	await redirectPage(false, url);
}

async function onBtnGoToSectionClick(parentId) {
	let url = composeUrlForSection(parentId);
	await redirectPage(false, url);
}

async function onBtnStartNewThreadClick(btn, panelCN, forumId) {
	disableButton(btn);
	let p = document.getElementById(panelCN);
	let div = newDiv();
	div.className = "newThreadCreation";
	div.id = "newThreadCreation";
	p.appendChild(div);

	// Draw.
	let d = newDiv();
	d.className = "title";
	d.textContent = "New Thread Parameters";
	div.appendChild(d);
	d = newDiv();
	d.innerHTML = '<label class="parameter" for="name">Name</label>' +
		'<input type="text" name="name" id="name" value="" class="newThreadName" />';
	div.appendChild(d);
	d = newDiv();
	d.innerHTML = '<label class="parameter" for="parent" title="ID of a parent forum" hidden="hidden">Parent</label>' +
		'<input type="text" name="parent" id="parent" value="' + forumId + '" readonly="readonly" hidden="hidden"/>';
	div.appendChild(d);
	d = newDiv();
	d.innerHTML = '<input type="button" class="' + ButtonClass.ConfirmThreadStart + '" value="' + ButtonName.ConfirmThreadStart + '" onclick="onBtnConfirmThreadStartClick(this)">';
	div.appendChild(d);
}

async function onBtnConfirmThreadStartClick(btn) {
	// Input.
	let pp = btn.parentNode.parentNode;
	let name = pp.childNodes[1].childNodes[1].value;
	if (name.length < 1) {
		console.error(Err.NameIsNotSet);
		return;
	}
	let parent = Number(pp.childNodes[2].childNodes[1].value);
	if (parent < 1) {
		console.error(Err.ParentIsNotSet);
		return;
	}

	// Work.
	let res = await addThread(parent, name);
	let threadId = res.threadId;
	disableParentForm(btn, pp, false);
	let txt = "A thread was created. ID=" + threadId.toString() + ".";
	showActionSuccess(btn, txt);
	await reloadPage(true);
}

function onBtnAddMessageClick(btn, panelCN, threadId) {
	disableButton(btn);
	let p = document.getElementById(panelCN);
	let div = newDiv();
	div.className = "newMessageCreation";
	div.id = "newMessageCreation";
	p.appendChild(div);

	// Draw.
	let d = newDiv();
	d.className = "title";
	d.textContent = "New Message Parameters";
	div.appendChild(d);
	d = newDiv();
	d.innerHTML = '<label class="parameter" for="txt">Text</label>' +
		'<textarea name="txt" id="txt" class="newMessageText"></textarea>';
	div.appendChild(d);
	d = newDiv();
	d.innerHTML = '<label class="parameter" for="parent" title="ID of a parent thread" hidden="hidden">Parent</label>' +
		'<input type="text" name="parent" id="parent" value="' + threadId + '" readonly="readonly" hidden="hidden"/>';
	div.appendChild(d);
	d = newDiv();
	d.innerHTML = '<input type="button" class="' + ButtonClass.ConfirmMessageCreation + '" value="' + ButtonName.ConfirmMessageCreation + '" onclick="onBtnConfirmMessageCreationClick(this)">';
	div.appendChild(d);
}

async function onBtnConfirmMessageCreationClick(btn) {
	// Input.
	let pp = btn.parentNode.parentNode;
	let text = pp.childNodes[1].childNodes[1].value;
	if (text.length < 1) {
		console.error(Err.TextIsNotSet);
		return;
	}
	let parent = Number(pp.childNodes[2].childNodes[1].value);
	if (parent < 1) {
		console.error(Err.ParentIsNotSet);
		return;
	}

	// Work.
	let res = await addMessage(parent, text);
	let messageId = res.messageId;
	disableParentForm(btn, pp, false);
	let txt = "A message was created. ID=" + messageId.toString() + ".";
	showActionSuccess(btn, txt);
	await reloadPage(true);
}

async function onBtnSubscribeClick(btn, threadId, userId) {
	let res = await addSubscription(threadId, userId);
	if (res.ok !== true) {
		return;
	}
	disableButton(btn);
}

async function onBtnEditMessageClick(btn, panelCN, messageId) {
	// Get edited message.
	let res = await getMessage(messageId);
	let message = jsonToMessages(res.message);

	disableButton(btn);
	let p = document.getElementById(panelCN);
	let div = newDiv();
	div.className = "messageEditing";
	div.id = "messageEditing";
	p.appendChild(div);

	// Draw.
	let d = newDiv();
	d.className = "title";
	d.textContent = "Message Editing";
	div.appendChild(d);
	d = newDiv();
	d.innerHTML = '<label class="parameter" for="txt">Text</label>' +
		'<textarea name="txt" id="txt" class="newMessageText">' + escapeHtml(message.Text) + '</textarea>';
	div.appendChild(d);
	d = newDiv();
	d.innerHTML = '<label class="parameter" for="id" title="ID of edited message" hidden="hidden">ID</label>' +
		'<input type="text" name="id" id="id" value="' + messageId + '" readonly="readonly" hidden="hidden"/>';
	div.appendChild(d);
	d = newDiv();
	d.innerHTML = '<input type="button" class="' + ButtonClass.ConfirmMessageEdit + '" value="' + ButtonName.ConfirmMessageEdit + '" onclick="onBtnConfirmMessageEditClick(this)">';
	div.appendChild(d);
}

async function onBtnConfirmMessageEditClick(btn) {
	// Input.
	let pp = btn.parentNode.parentNode;
	let newText = pp.childNodes[1].childNodes[1].value;
	if (newText.length < 1) {
		console.error(Err.TextIsNotSet);
		return;
	}
	let messageId = Number(pp.childNodes[2].childNodes[1].value);
	if (messageId < 1) {
		console.error(Err.IdNotSet);
		return;
	}

	// Work.
	let res = await changeMessageText(messageId, newText);
	if (res.ok !== true) {
		return;
	}
	disableParentForm(btn, pp, false);
	let txt = "Message text was changed.";
	showActionSuccess(btn, txt);
	await reloadPage(true);
}

async function onBtnMarkNotificationAsReadClick(btn, notificationId) {
	let res = await markNotificationAsRead(notificationId);
	if (res.ok !== true) {
		return;
	}

	// Change the 'unread' class in the table's row.
	let tr = btn.parentNode.parentNode;
	let td;
	for (let i = 0; i < tr.childNodes.length; i++) {
		td = tr.childNodes[i];
		if (td.className === 'unread') {
			td.className = 'read';
		}
	}

	// Hide the button.
	disableButton(btn);
}

async function onBtnDeleteNotificationClick(btn, notificationId) {
	let res = await deleteNotification(notificationId);
	if (res.ok !== true) {
		return;
	}

	// Change the 'unread' class in the table's row.
	let tr = btn.parentNode.parentNode;
	let td;
	for (let i = 0; i < tr.childNodes.length; i++) {
		td = tr.childNodes[i];
		if (td.className === '') {
			td.className = 'deleted';
		}
	}

	// Hide the button.
	disableButton(btn);
}

function onBtnShowFullNotificationClick(btn) {
	let tr = btn.parentNode.parentNode;
	let td = tr.childNodes[2];
	let subtableTbody = td.childNodes[0].childNodes[0];
	let childShort = subtableTbody.childNodes[0].childNodes[0];
	let childFull = subtableTbody.childNodes[1].childNodes[0];

	if (childShort.className === 'visible') {
		btn.value = "Hide";
		childShort.className = 'hidden';
		childFull.className = 'visible';
	} else {
		btn.value = "Show";
		childFull.className = 'hidden';
		childShort.className = 'visible';
	}
}

async function onBtnUnsubscribeClick(btn) {
	let tr = btn.parentElement.parentElement;
	let threadId = Number(tr.children[1].textContent);
	let res = await deleteSelfSubscription(threadId);
	if (!res.ok) {
		return;
	}
	tr.style.display = "none";
}

// Other functions.

async function addPageHead(el, text, atTop) {
	let settings = getSettings();
	let isLoggedInB = isLoggedIn(settings);
	let unreadNotificationsCount = -1;
	if (isLoggedInB) {
		unreadNotificationsCount = (await countUnreadNotifications()).unc;
	}

	// Draw.
	let div = newDiv();
	let cn = PageZoneClass.PageHead;
	div.className = cn;
	div.id = cn;
	let tbl = newTable();
	let tr = newTr();
	let tdL = newTd();
	tdL.className = cn + PageZoneClass.L;
	tdL.id = cn + PageZoneClass.L;
	tdL.textContent = text;
	tr.appendChild(tdL);
	let tdR = newTd();
	tdR.className = cn + PageZoneClass.R;
	tdR.id = cn + PageZoneClass.R;
	if (!isLoggedInB) {
		tdR.cssText = '';
	} else {
		let bcn;
		let bt;
		if (unreadNotificationsCount > 0) {
			if (unreadNotificationsCount >= settings.NotificationCountLimit) {
				bcn = ButtonClass.NotificationsBoxFull;
				bt = Hint.NotificationBoxIsFull;
			} else {
				bcn = ButtonClass.NotificationsOn;
				bt = Hint.NewNotifications;
			}
		} else {
			bcn = ButtonClass.NotificationsOff;
			bt = Hint.NoNewNotifications;
		}

		tdR.innerHTML = '<table><tr>' +
			'<td><input type="button" value="' + ButtonName.Notifications + '" class="' + bcn + '" onclick="onBtnNotificationsClick(this)" title="' + bt + '" /></td>'
			+ '<td><input type="button" value="' + ButtonName.Account + '" class="' + ButtonClass.Account + '" onclick="onBtnAccountClick(this)" /></td>' +
			'</tr></table>';
	}
	tr.appendChild(tdR);
	tbl.appendChild(tr);
	div.appendChild(tbl);

	if (atTop) {
		el.insertBefore(div, el.firstChild);
	} else {
		el.appendChild(div);
	}
}

function addActionPanel(el, atTop, parentType, parentId) {
	let div = newDiv();
	let cn = PageZoneClass.ActionPanel;
	div.className = cn;
	div.id = cn;
	let tbl = newTable();
	let tr = newTr();
	let td = newTd();
	td.innerHTML = '<form><input type="button" value="' + ButtonName.BackToRoot + '" class="' + ButtonClass.BackToRoot + '" onclick="onBtnGoToIndexClick(this)" /></form>';
	tr.appendChild(td);

	if (parentType != null) {
		switch (parentType) {
			case ObjectType.Thread:
				td = newTd();
				td.innerHTML = '<span class="spacerA">&nbsp;</span>' +
					'<input type="button" value="' + ButtonName.BackToThread + '" class="' + ButtonClass.BackToThread + '" onclick="onBtnGoToThreadClick(' + parentId + ')" />';
				tr.appendChild(td);
				break;

			case ObjectType.Forum:
				td = newTd();
				td.innerHTML = '<span class="spacerA">&nbsp;</span>' +
					'<input type="button" value="' + ButtonName.BackToForum + '" class="' + ButtonClass.BackToForum + '" onclick="onBtnGoToForumClick(' + parentId + ')" />';
				tr.appendChild(td);
				break;

			case ObjectType.Section:
				td = newTd();
				td.innerHTML = '<span class="spacerA">&nbsp;</span>' +
					'<input type="button" value="' + ButtonName.BackToSection + '" class="' + ButtonClass.BackToSection + '" onclick="onBtnGoToSectionClick(' + parentId + ')" />';
				tr.appendChild(td);
				break;
		}
	}

	tbl.appendChild(tr);
	div.appendChild(tbl);

	if (atTop) {
		el.insertBefore(div, el.firstChild);
	} else {
		el.appendChild(div);
	}
}

async function addBottomActionPanel(el, objectType, object) {
	let res = await getSelfRoles();
	let user = jsonToUser(res.user);
	let cn = PageZoneClass.BottomActionPanel;
	let cnL = cn + PageZoneClass.L;
	let cnR = cn + PageZoneClass.R;

	let div = newDiv();
	div.className = cn;
	div.id = cn;
	let tbl = newTable();
	let tr = newTr();
	tbl.appendChild(tr);
	let tdL = newTd();
	tdL.id = cnL;
	tdL.className = cnL;
	tr.appendChild(tdL);
	let tdR = newTd();
	tdR.id = cnR;
	tdR.className = cnR;
	tr.appendChild(tdR);
	div.appendChild(tbl);
	el.appendChild(div);

	switch (objectType) {
		case ObjectType.Forum:
			await drawBottomActionPanelButtonsForForum(object, tdL, tdR, user, cn);
			break;

		case ObjectType.Thread:
			await drawBottomActionPanelButtonsForThread(object, tdL, tdR, user, cn);
			break;

		case ObjectType.Message:
			await drawBottomActionPanelButtonsForMessage(object, tdL, tdR, user, cn);
			break;
	}
}

function drawPageTitle(p, title) {
	let d = newDiv();
	d.className = PageZoneClass.PageTitle;
	let ml = SectionMarginDelta;
	d.style.cssText = "margin-left: " + ml + "px";
	d.textContent = title;
	p.appendChild(d);
}

async function drawBottomActionPanelButtonsForForum(forum, tdL, tdR, user, panelClass) {
	if (user.Roles.IsAuthor) {
		tdR.innerHTML = '<form><input type="button" value="' + ButtonName.StartNewThread + '" class="' + ButtonClass.StartNewThread + '" ' +
			'onclick="onBtnStartNewThreadClick(this, \'' + panelClass + '\', ' + forum.Id + ')" /></form>';
	}
}

async function drawBottomActionPanelButtonsForThread(thread, tdL, tdR, user, panelClass) {
	let resA = await getLatestMessageOfThread(thread.Id);
	let latestMessageInThread = jsonToMessage(resA.message);
	let resB = await isSelfSubscribed(thread.Id);
	let isUserSubscribed = resB.isSubscribed;

	let tbl = newTable();
	let tr = newTr();

	if (!isUserSubscribed) {
		let td = newTd();
		td.innerHTML += '<form><input type="button" value="' + ButtonName.SubscribeToThread + '" class="' + ButtonClass.SubscribeToThread + '" ' +
			'onclick="onBtnSubscribeClick(this, ' + thread.Id + ', ' + user.Id + ')" /></form>';
		tr.appendChild(td);
	}

	if (user.canAddMessage(latestMessageInThread)) {
		let td = newTd();
		td.innerHTML += '<form><input type="button" value="' + ButtonName.AddMessage + '" class="' + ButtonClass.AddMessage + '" ' +
			'onclick="onBtnAddMessageClick(this, \'' + panelClass + '\', ' + thread.Id + ')" /></form>';
		tr.appendChild(td);
	}

	tbl.appendChild(tr);
	tdR.appendChild(tbl);
}

async function drawBottomActionPanelButtonsForMessage(message, tdL, tdR, user, panelClass) {
	if (user.canEditMessage(message)) {
		tdR.innerHTML = '<form><input type="button" value="' + ButtonName.EditMessage + '" class="' + ButtonClass.EditMessage + '" ' +
			'onclick="onBtnEditMessageClick(this, \'' + panelClass + '\', ' + message.Id + ')" /></form>';
	}
}

async function drawMessage(p, message) {
	// Header.
	let divMsgHdr = newDiv();
	divMsgHdr.className = "messageHeader";
	divMsgHdr.id = "messageHeader_" + message.Id;
	let ml = SectionMarginDelta * 2;
	divMsgHdr.style.cssText = "margin-left: " + ml + "px";
	divMsgHdr.innerHTML = await composeMessageHeaderText(message);
	p.appendChild(divMsgHdr);

	// Body.
	let divMsgBody = newDiv();
	divMsgBody.className = "messageBody";
	divMsgBody.id = "messageBody_" + message.Id;
	ml = SectionMarginDelta * 2;
	divMsgBody.style.cssText = "margin-left: " + ml + "px";
	divMsgBody.innerHTML = processMessageText(message.Text);
	p.appendChild(divMsgBody);
}

async function drawThreadMessages(p, thread, messagesMap) {
	let messageIds = thread.Messages;
	if (messageIds == null) {
		return
	}

	let divMsgHdr, divMsgBody, ml, messageId, message, txt;
	for (let i = 0; i < messageIds.length; i++) {
		messageId = messageIds[i];
		if (!messagesMap.has(messageId)) {
			console.error(Err.MessageNotFound);
			return;
		}
		message = messagesMap.get(messageId);

		// Header.
		divMsgHdr = newDiv();
		divMsgHdr.className = "messageHeader";
		divMsgHdr.id = "messageHeader_" + message.Id;
		ml = SectionMarginDelta * 2;
		divMsgHdr.style.cssText = "margin-left: " + ml + "px";
		txt = await composeMessageHeaderText(message);
		divMsgHdr.innerHTML = txt;
		p.appendChild(divMsgHdr);

		// Body.
		divMsgBody = newDiv();
		divMsgBody.className = "messageBody";
		divMsgBody.id = "messageBody_" + message.Id;
		ml = SectionMarginDelta * 2;
		divMsgBody.style.cssText = "margin-left: " + ml + "px";
		divMsgBody.innerHTML = processMessageText(message.Text);
		p.appendChild(divMsgBody);
	}
}

function drawForumThreads(p, forum, threadsMap) {
	let divThread, ml, url, threadId, thread;
	let threadIds = forum.Threads;
	if (threadIds != null) {
		for (let i = 0; i < threadIds.length; i++) {
			threadId = threadIds[i];
			if (!threadsMap.has(threadId)) {
				console.error(Err.ThreadNotFound);
				return;
			}
			thread = threadsMap.get(threadId);

			divThread = newDiv();
			divThread.className = "thread";
			divThread.id = "thread_" + thread.Id;
			ml = SectionMarginDelta * 2;
			divThread.style.cssText = "margin-left: " + ml + "px";
			url = composeUrlForThread(threadId);
			divThread.innerHTML = "<a href='" + url + "'>" + thread.Name + "</a>";
			p.appendChild(divThread);
		}
	}
}

function drawSectionsAndForums(p, nodes, forumsMap) {
	let node, divSection, divForum, ml, url, forumId, forum, sectionForums;
	for (let i = 0; i < nodes.length; i++) {
		node = nodes[i];

		divSection = newDiv();
		divSection.className = "section";
		divSection.id = "section_" + node.Section.Id;
		ml = SectionMarginDelta * node.Level;
		divSection.style.cssText = "margin-left: " + ml + "px";
		url = composeUrlForSection(node.Section.Id);
		divSection.innerHTML = "<a href='" + url + "'>" + node.Section.Name + "</a>";
		p.appendChild(divSection);

		if (node.Section.ChildType === SectionChildType.Forum) {
			sectionForums = node.Section.Children;
		} else {
			sectionForums = [];
		}
		for (let j = 0; j < sectionForums.length; j++) {
			forumId = sectionForums[j];
			forum = forumsMap.get(forumId);

			divForum = newDiv();
			divForum.className = "forum";
			divForum.id = "forum_" + forumId;
			ml = SectionMarginDelta * (node.Level + 1);
			divForum.style.cssText = "margin-left: " + ml + "px";
			url = composeUrlForForum(forumId);
			divForum.innerHTML = "<a href='" + url + "'>" + forum.Name + "</a>";
			p.appendChild(divForum);
		}
	}
}

function drawSubscriptions(p, elClass, swtnop) {
	let tbl = newTable();
	tbl.className = elClass;

	// Header.
	let tr = newTr();
	let ths = ["#", "Thread ID", "Thread Name", "Actions"];
	let th;
	for (let i = 0; i < ths.length; i++) {
		th = newTh();
		if (i === 0) {
			th.className = "numCol";
		}
		th.textContent = ths[i];
		tr.appendChild(th);
	}
	tbl.appendChild(tr);
	p.appendChild(tbl);

	let columnsWithLink = [1, 2];
	let columnsWithHtml = [3];

	if (swtnop.ThreadIds == null) {
		return;
	}

	// Cells.
	let threadId, threadName;
	for (let i = 0; i < swtnop.ThreadIds.length; i++) {
		threadId = swtnop.ThreadIds[i];
		threadName = swtnop.ThreadNames[i];

		// Fill data.
		tr = newTr();
		let tds = [];
		for (let j = 0; j < ths.length; j++) {
			tds.push("");
		}

		tds[0] = (i + 1).toString();
		tds[1] = threadId.toString();
		tds[2] = threadName;
		tds[3] = '<input type="button" class="' + ButtonClass.Unsubscribe + '" value="' + ButtonName.Unsubscribe + '" onclick="onBtnUnsubscribeClick(this)">';

		let td, url;
		for (let j = 0; j < tds.length; j++) {
			url = composeUrlForThread(threadId);
			td = newTd();

			if (j === 0) {
				td.className = "numCol";
			}

			if (columnsWithLink.includes(j)) {
				td.innerHTML = '<a href="' + url + '">' + tds[j] + '</a>';
			} else {
				if (columnsWithHtml.includes(j)) {
					td.innerHTML = tds[j];
				} else {
					td.textContent = tds[j];
				}
			}
			tr.appendChild(td);
		}

		tbl.appendChild(tr);
	}
}

function drawNotifications(p, elClass, notifications) {
	let tbl = newTable();
	tbl.className = elClass;

	// Header.
	let tr = newTr();
	let ths = ["#", "Time", "Text", "Actions"];
	let th;
	for (let i = 0; i < ths.length; i++) {
		th = newTh();
		if (i === 0) {
			th.className = "numCol";
		}
		th.textContent = ths[i];
		tr.appendChild(th);
	}
	tbl.appendChild(tr);
	p.appendChild(tbl);

	let columnsWithHtml = [2, 3];

	// Cells.
	let notification, actions;
	for (let i = 0; i < notifications.length; i++) {
		notification = notifications[i];

		tr = newTr();
		let tds = [];
		for (let j = 0; j < ths.length; j++) {
			tds.push("");
		}

		tds[0] = (i + 1).toString();
		tds[1] = prettyTime(notification.ToC);
		tds[2] = splitNotificationTextCell(notification.Text);

		actions = '<input type="button" class="' + ButtonClass.ShowFullNotificationU + '" value="' + ButtonName.ShowFullNotificationU + '" onclick="onBtnShowFullNotificationClick(this)">';
		if (!notification.IsRead) {
			actions += '<input type="button" class="' + ButtonClass.MarkNotificationAsReadU + '" value="' + ButtonName.MarkNotificationAsReadU + '" onclick="onBtnMarkNotificationAsReadClick(this, ' + notification.Id + ')">';
		} else {
			actions += '<input type="button" class="' + ButtonClass.DeleteNotificationU + '" value="' + ButtonName.DeleteNotificationU + '" onclick="onBtnDeleteNotificationClick(this, ' + notification.Id + ')">';
		}
		tds[3] = actions;

		let td;
		let jLast = tds.length - 1;
		for (let j = 0; j < tds.length; j++) {
			td = newTd();

			if (j === 0) {
				td.className = "numCol";
			} else if (j === jLast) {
				td.className = "lastCol";
			} else {
				if (!notification.IsRead) {
					td.className = "unread";
				}
				if (j === 1) {
					td.className += " col2";
				}
			}

			if (columnsWithHtml.includes(j)) {
				td.innerHTML = tds[j];
			} else {
				td.textContent = tds[j];
			}
			tr.appendChild(td);
		}

		tbl.appendChild(tr);
	}
}

async function drawUserPage(user) {
	let subscriptionsCount = (await countSelfSubscriptions()).userSubscriptionsCount;

	document.getElementById(Fi.id18).value = user.Name;
	document.getElementById(Fi.id19).value = user.Email;
	document.getElementById(Fi.id20).value = prettyTime(user.RegTime);

	let roles = [];
	if (user.Roles.IsAdministrator) {
		roles.push(UserRole.Administrator);
	}
	if (user.Roles.IsModerator) {
		roles.push(UserRole.Moderator);
	}
	if (user.Roles.IsAuthor) {
		roles.push(UserRole.Author);
	}
	if (user.Roles.IsWriter) {
		roles.push(UserRole.Writer);
	}
	if (user.Roles.IsReader) {
		roles.push(UserRole.Reader);
	}
	let rolesHtml = "";
	let cn = "userPageRole";
	for (let i = 0; i < roles.length; i++) {
		if (roles[i] !== UserRole.Administrator) {
			rolesHtml += '<span class="' + cn + '">' + roles[i] + '</span>';
		} else {
			rolesHtml += '<span class="' + cn + '"><a href="' + adminPage + '" target="_blank" rel="noopener noreferrer">' + roles[i] + '</a></span>';
		}
	}
	let tr = document.getElementById(Fi.id21_tr);
	tr.children[1].innerHTML = rolesHtml;

	document.getElementById(Fi.id22).value = subscriptionsCount.toString();
}

function splitNotificationTextCell(fullText) {
	let shortText = composeNotificationShortText(fullText);
	return '<table>' +
		'<tr><td class="visible">' + shortText + '</td></tr>' +
		'<tr><td class="hidden">' + fullText + '</td></tr>' +
		'</table>';
}

function setCaptchaInputsVisibility(isCaptchaNeeded, captchaId, cptImageTr, cptImage, cptAnswerTr, cptAnswer) {
	if (isCaptchaNeeded) {
		cptImageTr.style.display = "table-row";
		cptAnswerTr.style.display = "table-row";
		if (captchaId.length > 0) {
			cptImage.src = composeCaptchaImageUrl(captchaId);
		}
		cptAnswer.enabled = true;
	} else {
		cptImageTr.style.display = "none";
		cptAnswerTr.style.display = "none";
		cptAnswer.enabled = false;
	}
}
