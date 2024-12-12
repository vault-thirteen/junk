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

ButtonName = {
	PaginatorPrev: "<",
	PaginatorNext: ">",
	BackA: "Go Back",
	Accept: "Accept",
	Reject: "Reject",
	LogOut: "Log Out",
	EnableRole: "Enable Role",
	DisableRole: "Disable Role",
	Proceed: "Proceed",
	CreateRootSection: "Create",
	CreateNormalSection: "Create",
	ChangeSectionName: "Change Name",
	ChangeSectionParent: "Change Parent",
	MoveSectionUp: "Move Up",
	MoveSectionDown: "Move Down",
	DeleteSection: "Delete",
	CreateForum: "Create",
	ChangeForumName: "Change Name",
	ChangeForumParent: "Change Parent",
	MoveForumUp: "Move Up",
	MoveForumDown: "Move Down",
	DeleteForum: "Delete",
	CreateThread: "Create",
	ChangeThreadName: "Change Name",
	ChangeThreadParent: "Change Parent",
	MoveThreadUp: "Move Up",
	MoveThreadDown: "Move Down",
	DeleteThread: "Delete",
	CreateMessage: "Create",
	ChangeMessageText: "Change Text",
	ChangeMessageParent: "Change Parent",
	DeleteMessage: "Delete",
	CreateNotification: "Create",
	DeleteNotification: "Delete",
	CreateResource: "Create",
	DeleteResource: "Delete",
}

ButtonClass = {
	PaginatorPrev: "btnPrev",
	PaginatorNext: "btnNext",
	BackA: "btnBack",
	Accept: "btnAccept",
	Reject: "btnReject",
	LogOut: "btnLogOut",
	EnableRole: "btnEnableRole",
	DisableRole: "btnDisableRole",
	Proceed: "btnProceed",
	CreateRootSection: "btnCreateRootSection",
	CreateNormalSection: "btnCreateNormalSection",
	ChangeSectionName: "btnChangeSectionName",
	ChangeSectionParent: "btnChangeSectionParent",
	MoveSection: "btnMoveSection",
	DeleteSection: "btnDeleteSection",
	CreateForum: "btnCreateForum",
	ChangeForumName: "btnChangeForumName",
	ChangeForumParent: "btnChangeForumParent",
	MoveForum: "btnMoveForum",
	DeleteForum: "btnDeleteForum",
	CreateThread: "btnCreateThread",
	ChangeThreadName: "btnChangeThreadName",
	ChangeThreadParent: "btnChangeThreadParent",
	MoveThread: "btnMoveThread",
	DeleteThread: "btnDeleteThread",
	CreateMessage: "btnCreateMessage",
	ChangeMessageText: "btnChangeMessageText",
	ChangeMessageParent: "btnChangeMessageParent",
	DeleteMessage: "btnDeleteMessage",
	CreateNotification: "btnCreateNotification",
	DeleteNotification: "btnDeleteNotification",
	CreateResource: "btnCreateResource",
	DeleteResource: "btnDeleteResource",
}

PageZoneClass = {
	Paginator: "paginator",
	SubpageTitleA: "subpageTitle",
}

EventHandlerVariant = {
	UserListPrevA: "userListPrev",
	UserListNextA: "userListNext",
	LoggedUserListPrevA: "loggedUserListPrev",
	LoggedUserListNextA: "loggedUserListNext",
	RrfaListPrevA: "rrfaListPrev",
	RrfaListNextA: "rrfaListNext",
	ResourceListPrevA: "resourceListPrev",
	ResourceListNextA: "resourceListNext",
}

// Global variables.
class GlobalVariablesContainer {
	constructor(areSettingsReady, settings, id, page, pages) {
		this.AreSettingsReady = areSettingsReady;
		this.Settings = settings;
		this.Id = id;
		this.Page = page;
		this.Pages = pages;
	}
}

mca_gvc = new GlobalVariablesContainer(false, null, 0, 0);

// Settings.

function isSettingsUpdateNeeded() {
	return true;
}

function saveSettings(s) {
	mca_gvc.Settings = s;
	mca_gvc.AreSettingsReady = true;
}

function getSettings() {
	if (!mca_gvc.AreSettingsReady) {
		console.error(Err.Settings);
		return null;
	}

	return mca_gvc.Settings;
}

// Entry point.
async function onPageLoad() {
	// Settings initialisation.
	let ok = await updateSettingsIfNeeded();
	if (!ok) {
		return;
	}

	// Select a page.
	let curPage = window.location.search;
	let sp = new URLSearchParams(curPage);

	if (sp.has(Qpn.RegistrationsReadyForApproval)) {
		if (!preparePageVariable(sp)) {
			return;
		}
		await showPage_RegistrationsReadyForApproval();
		return;
	}

	if (sp.has(Qpn.ListOfLoggedUsers)) {
		if (!preparePageVariable(sp)) {
			return;
		}
		await showPage_ListOfLoggedUsers();
		return;
	}

	if (sp.has(Qpn.ListOfUsers)) {
		if (!preparePageVariable(sp)) {
			return;
		}
		await showPage_ListOfUsers();
		return;
	}

	if (sp.has(Qpn.ListOfResources)) {
		if (!preparePageVariable(sp)) {
			return;
		}
		await showPage_ListOfResources();
		return;
	}

	if (sp.has(Qpn.UserPage)) {
		if (!prepareIdVariable(sp)) {
			return;
		}
		await showPage_UserPage();
		return;
	}

	if (sp.has(Qpn.ManagerOfSections)) {
		await showPage_ManagerOfSections();
		return;
	}

	if (sp.has(Qpn.ManagerOfForums)) {
		await showPage_ManagerOfForums();
		return;
	}

	if (sp.has(Qpn.ManagerOfThreads)) {
		await showPage_ManagerOfThreads();
		return;
	}

	if (sp.has(Qpn.ManagerOfMessages)) {
		await showPage_ManagerOfMessages();
		return;
	}

	if (sp.has(Qpn.ManagerOfNotifications)) {
		await showPage_ManagerOfNotifications();
		return;
	}

	if (sp.has(Qpn.ManagerOfResources)) {
		await showPage_ManagerOfResources();
		return;
	}

	showPage_MainMenu();
}

async function showPage_RegistrationsReadyForApproval() {
	let res = await getListOfRegistrationsReadyForApproval(mca_gvc.Page);
	let pageCount = res.pageData.totalPages;
	if (!preparePageNumber(pageCount)) {
		return
	}

	let rrfas = jsonToRrfas(res.rrfa);

	// Draw.
	let p = document.getElementById("subpage");
	showBlock(p);
	addBtnBack(p);
	addTitle(p, "List of Registrations ready for Approval");
	addPaginator(p, mca_gvc.Page, mca_gvc.Pages, EventHandlerVariant.RrfaListPrevA, EventHandlerVariant.RrfaListNextA);
	addDiv(p, "subpageListOfRRFA");
	await fillListOfRRFA("subpageListOfRRFA", rrfas);
}

async function showPage_ListOfLoggedUsers() {
	let res = await getListOfLoggedUsersOnPage(mca_gvc.Page);
	let pageCount = res.pageData.totalPages;
	if (!preparePageNumber(pageCount)) {
		return
	}

	let userIds = res.loggedUserIds;

	// Draw.
	let p = document.getElementById("subpage");
	showBlock(p);
	addBtnBack(p);
	addTitle(p, "List of logged-in Users");
	addPaginator(p, mca_gvc.Page, mca_gvc.Pages, EventHandlerVariant.LoggedUserListPrevA, EventHandlerVariant.LoggedUserListNextA);
	addDiv(p, "subpageListOfLoggedUsers");
	await fillListOfLoggedUsers("subpageListOfLoggedUsers", userIds);
}

async function showPage_ListOfUsers() {
	let resA = await getListOfAllUsersOnPage(mca_gvc.Page);
	let pageCount = resA.pageData.totalPages;
	if (!preparePageNumber(pageCount)) {
		return
	}

	let userIds = resA.userIds;
	let resB = await getListOfLoggedUsers();
	let loggedUserIds = resB.loggedUserIds;

	// Draw.
	let p = document.getElementById("subpage");
	showBlock(p);
	addBtnBack(p);
	addTitle(p, "List of All Users");
	addPaginator(p, mca_gvc.Page, mca_gvc.Pages, EventHandlerVariant.UserListPrevA, EventHandlerVariant.UserListNextA);
	addDiv(p, "subpageListOfUsers");
	await fillListOfUsers("subpageListOfUsers", userIds, loggedUserIds);
}

async function showPage_ListOfResources() {
	let res = await getListOfAllResourcesOnPage(mca_gvc.Page);
	let pageCount = res.rop.pageData.totalPages;
	if (!preparePageNumber(pageCount)) {
		return
	}

	let resourceIds = res.rop.resourceIds;

	// Draw.
	let p = document.getElementById("subpage");
	showBlock(p);
	addBtnBack(p);
	addTitle(p, "List of Resources");
	addPaginator(p, mca_gvc.Page, mca_gvc.Pages, EventHandlerVariant.ResourceListPrevA, EventHandlerVariant.ResourceListNextA);
	addDiv(p, "subpageListOfResources");
	await fillListOfResources("subpageListOfResources", resourceIds);
}

async function showPage_UserPage() {
	let userId = mca_gvc.Id;
	let resA = await viewUserParameters(userId);
	let user = jsonToUser(resA.user);

	let resB = await isUserLoggedIn(userId);
	let userLogInState = resB.isUserLoggedIn;
	let resC = await getSelfRoles();
	let selfUser = jsonToUser(resC.user);

	// Draw.
	let p = document.getElementById("subpage");
	showBlock(p);
	addBtnBack(p);
	addTitle(p, "User Page");
	addDiv(p, "subpageUserPage");
	fillUserPage("subpageUserPage", user, userLogInState, selfUser);
}

async function showPage_ManagerOfSections() {
	// Draw.
	let p = document.getElementById("subpage");
	showBlock(p);
	addBtnBack(p);
	addTitle(p, "Management of Sections");
	addDiv(p, "sectionManager");
	fillSectionManager("sectionManager");
}

async function showPage_ManagerOfForums() {
	// Draw.
	let p = document.getElementById("subpage");
	showBlock(p);
	addBtnBack(p);
	addTitle(p, "Management of Forums");
	addDiv(p, "forumManager");
	fillForumManager("forumManager");
}

async function showPage_ManagerOfThreads() {
	// Draw.
	let p = document.getElementById("subpage");
	showBlock(p);
	addBtnBack(p);
	addTitle(p, "Management of Threads");
	addDiv(p, "threadManager");
	fillThreadManager("threadManager");
}

async function showPage_ManagerOfMessages() {
	// Draw.
	let p = document.getElementById("subpage");
	showBlock(p);
	addBtnBack(p);
	addTitle(p, "Management of Messages");
	addDiv(p, "messageManager");
	fillMessageManager("messageManager");
}

async function showPage_ManagerOfNotifications() {
	// Draw.
	let p = document.getElementById("subpage");
	showBlock(p);
	addBtnBack(p);
	addTitle(p, "Management of Notifications");
	addDiv(p, "notificationManager");
	fillNotificationManager("notificationManager");
}

async function showPage_ManagerOfResources() {
	// Draw.
	let p = document.getElementById("subpage");
	showBlock(p);
	addBtnBack(p);
	addTitle(p, "Management of Resources");
	addDiv(p, "resourceManager");
	fillResourceManager("resourceManager");
}

function showPage_MainMenu() {
	let tbl = document.getElementById("acpMenu");
	showTable(tbl);
}

// Event handlers.

async function onGoRegApprovalClick(btn) {
	await redirectToSubPageA(false, Qp.Prefix + Qpn.RegistrationsReadyForApproval);
}

async function onGoLoggedUsersClick(btn) {
	await redirectToSubPageA(false, Qp.Prefix + Qpn.ListOfLoggedUsers);
}

async function onGoListAllUsersClick(btn) {
	await redirectToSubPageA(false, Qp.Prefix + Qpn.ListOfUsers);
}

async function onGoListResourcesClick(btn) {
	await redirectToSubPageA(false, Qp.Prefix + Qpn.ListOfResources);
}

async function onGoManageSectionsClick(btn) {
	await redirectToSubPageA(false, Qp.Prefix + Qpn.ManagerOfSections);
}

async function onGoManageForumsClick(btn) {
	await redirectToSubPageA(false, Qp.Prefix + Qpn.ManagerOfForums);
}

async function onGoManageThreadsClick(btn) {
	await redirectToSubPageA(false, Qp.Prefix + Qpn.ManagerOfThreads);
}

async function onGoManageMessagesClick(btn) {
	await redirectToSubPageA(false, Qp.Prefix + Qpn.ManagerOfMessages);
}

async function onGoManageNotificationsClick(btn) {
	await redirectToSubPageA(false, Qp.Prefix + Qpn.ManagerOfNotifications);
}

async function onGoManageResourcesClick(btn) {
	await redirectToSubPageA(false, Qp.Prefix + Qpn.ManagerOfResources);
}

async function onBtnPrevClick_resources(btn) {
	if (mca_gvc.Page <= 1) {
		console.error(Err.PreviousPageDoesNotExist);
		return;
	}

	mca_gvc.Page--;
	let url = composeUrlForAdminPageA(Qpn.ListOfResources, mca_gvc.Page);
	await redirectPage(false, url);
}

async function onBtnNextClick_resources(btn) {
	if (mca_gvc.Page >= mca_gvc.Pages) {
		console.error(Err.NextPageDoesNotExist);
		return;
	}

	mca_gvc.Page++;
	let url = composeUrlForAdminPageA(Qpn.ListOfResources, mca_gvc.Page);
	await redirectPage(false, url);
}

async function onBtnPrevClick_rrfa(btn) {
	if (mca_gvc.Page <= 1) {
		console.error(Err.PreviousPageDoesNotExist);
		return;
	}

	mca_gvc.Page--;
	let url = composeUrlForAdminPageA(Qpn.RegistrationsReadyForApproval, mca_gvc.Page);
	await redirectPage(false, url);
}

async function onBtnNextClick_rrfa(btn) {
	if (mca_gvc.Page >= mca_gvc.Pages) {
		console.error(Err.NextPageDoesNotExist);
		return;
	}

	mca_gvc.Page++;
	let url = composeUrlForAdminPageA(Qpn.RegistrationsReadyForApproval, mca_gvc.Page);
	await redirectPage(false, url);
}

async function onBtnPrevClick_logged(btn) {
	if (mca_gvc.Page <= 1) {
		console.error(Err.PreviousPageDoesNotExist);
		return;
	}

	mca_gvc.Page--;
	let url = composeUrlForAdminPageA(Qpn.ListOfLoggedUsers, mca_gvc.Page);
	await redirectPage(false, url);
}

async function onBtnNextClick_logged(btn) {
	if (mca_gvc.Page >= mca_gvc.Pages) {
		console.error(Err.NextPageDoesNotExist);
		return;
	}

	mca_gvc.Page++;
	let url = composeUrlForAdminPageA(Qpn.ListOfLoggedUsers, mca_gvc.Page);
	await redirectPage(false, url);
}

async function onBtnPrevClick_userList(btn) {
	if (mca_gvc.Page <= 1) {
		console.error(Err.PreviousPageDoesNotExist);
		return;
	}

	mca_gvc.Page--;
	let url = composeUrlForAdminPageA(Qpn.ListOfUsers, mca_gvc.Page);
	await redirectPage(false, url);
}

async function onBtnNextClick_userList(btn) {
	if (mca_gvc.Page >= mca_gvc.Pages) {
		console.error(Err.NextPageDoesNotExist);
		return;
	}

	mca_gvc.Page++;
	let url = composeUrlForAdminPageA(Qpn.ListOfUsers, mca_gvc.Page);
	await redirectPage(false, url);
}


async function onBtnAcceptClick(btn) {
	let tr = btn.parentElement.parentElement;
	let reqEmail = tr.children[3].textContent;
	let res = await approveAndRegisterUser(reqEmail);
	if (!res.ok) {
		return;
	}
	hideBlock(tr);
}

async function onBtnRejectClick(btn) {
	let tr = btn.parentElement.parentElement;
	let reqId = Number(tr.children[1].textContent);
	let res = await rejectRegistrationRequest(reqId);
	if (!res.ok) {
		return;
	}
	hideBlock(tr);
}

async function onBtnLogOutClick(btn) {
	let tr = btn.parentElement.parentElement;
	let userId = Number(tr.children[1].textContent);
	let res = await logUserOutA(userId);
	if (!res.ok) {
		return;
	}
	hideBlock(tr);
}

async function onBtnLogOutUPClick(userId) {
	let res = await logUserOutA(userId);
	if (!res.ok) {
		return;
	}
	await reloadPage(false);
}

async function onBtnEnableRoleUPClick(role, userId) {
	let res;
	switch (role) {
		case UserRole.Author:
			res = await setUserRoleAuthor(userId, true);
			break;

		case UserRole.Writer:
			res = await setUserRoleWriter(userId, true);
			break;

		case UserRole.Reader:
			res = await setUserRoleReader(userId, true);
			break;

		case UserRole.Logging:
			res = await unbanUser(userId);
			break;

		default:
			return;
	}

	if (!res.ok) {
		return;
	}
	await reloadPage(false);
}

async function onBtnDisableRoleUPClick(role, userId) {
	let res;
	switch (role) {
		case UserRole.Author:
			res = await setUserRoleAuthor(userId, false);
			break;

		case UserRole.Writer:
			res = await setUserRoleWriter(userId, false);
			break;

		case UserRole.Reader:
			res = await setUserRoleReader(userId, false);
			break;

		case UserRole.Logging:
			res = await banUser(userId);
			break;

		default:
			return;
	}

	if (!res.ok) {
		return;
	}
	await reloadPage(false);
}


function onSectionManagerBtnProceedClick(btn) {
	let selectedActionIdx = getSelectedActionIdxBPC(btn);
	if (selectedActionIdx == null) {
		return;
	}

	btn.disabled = true;
	disableParentFormBPC(btn);

	// Draw.
	let sm = document.getElementById("sectionManager");
	let fs = newFieldset();
	sm.appendChild(fs);

	let d = newDiv();
	d.className = "title";
	d.textContent = "Section Parameters";
	fs.appendChild(d);

	switch (selectedActionIdx) {
		case 1: // Create a root section.
			d = newDiv();
			d.innerHTML = htmlInputParameterName("");
			fs.appendChild(d);
			d = newDiv();
			d.innerHTML = '<input type="button" class="' + ButtonClass.CreateRootSection + '" value="' + ButtonName.CreateRootSection + '" onclick="onBtnCreateRootSectionClick(this)">';
			fs.appendChild(d);
			break;

		case 2: // Create a normal section.
			d = newDiv();
			d.innerHTML = htmlInputParameterName("");
			fs.appendChild(d);
			d = newDiv();
			d.innerHTML = htmlInputParameterParent("ID of a parent section");
			fs.appendChild(d);
			d = newDiv();
			d.innerHTML = '<input type="button" class="' + ButtonClass.CreateNormalSection + '" value="' + ButtonName.CreateNormalSection + '" onclick="onBtnCreateNormalSectionClick(this)">';
			fs.appendChild(d);
			break;

		case 3: // Change section's name.
			d = newDiv();
			d.innerHTML = htmlInputParameterId("ID of the changed section");
			fs.appendChild(d);
			d = newDiv();
			d.innerHTML = htmlInputParameterNewName("New name of the section");
			fs.appendChild(d);
			d = newDiv();
			d.innerHTML = '<input type="button" class="' + ButtonClass.ChangeSectionName + '" value="' + ButtonName.ChangeSectionName + '" onclick="onBtnChangeSectionNameClick(this)">';
			fs.appendChild(d);
			break;

		case 4: // Change section's parent.
			d = newDiv();
			d.innerHTML = htmlInputParameterId("ID of the changed section");
			fs.appendChild(d);
			d = newDiv();
			d.innerHTML = htmlInputParameterNewParent("ID of the new parent");
			fs.appendChild(d);
			d = newDiv();
			d.innerHTML = '<input type="button" class="' + ButtonClass.ChangeSectionParent + '" value="' + ButtonName.ChangeSectionParent + '" onclick="onBtnChangeSectionParentClick(this)">';
			fs.appendChild(d);
			break;

		case 5: // Move section up & down.
			d = newDiv();
			d.innerHTML = htmlInputParameterId("ID of the moved section");
			fs.appendChild(d);
			d = newDiv();
			d.innerHTML = '<input type="button" class="' + ButtonClass.MoveSection + '" value="' + ButtonName.MoveSectionUp + '" onclick="onBtnMoveSectionUpClick(this)">' +
				'<span class="subpageSpacerA">&nbsp;</span>' +
				'<input type="button" class="' + ButtonClass.MoveSection + '" value="' + ButtonName.MoveSectionDown + '" onclick="onBtnMoveSectionDownClick(this)">';
			fs.appendChild(d);
			break;

		case 6: // Delete a section.
			d = newDiv();
			d.innerHTML = htmlInputParameterId("ID of the section to delete");
			fs.appendChild(d);
			d = newDiv();
			d.innerHTML = '<input type="button" class="' + ButtonClass.DeleteSection + '" value="' + ButtonName.DeleteSection + '" onclick="onBtnDeleteSectionClick(this)">';
			fs.appendChild(d);
			break;
	}
}

async function onBtnCreateRootSectionClick(btn) {
	// Input.
	let pp = btn.parentNode.parentNode;
	let name = pp.childNodes[1].childNodes[1].value;
	if (name.length < 1) {
		console.error(Err.NameIsNotSet);
		return;
	}

	// Work.
	let res = await addSection(null, name);
	let sectionId = res.sectionId;
	disableParentForm(btn, pp, false);
	let txt = "A root section was created. ID=" + sectionId.toString() + ".";
	showActionSuccess(btn, txt);
	await reloadPage(true);
}

async function onBtnCreateNormalSectionClick(btn) {
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
	let res = await addSection(parent, name);
	let sectionId = res.sectionId;
	disableParentForm(btn, pp, false);
	let txt = "A normal section was created. ID=" + sectionId.toString() + ".";
	showActionSuccess(btn, txt);
	await reloadPage(true);
}

async function onBtnChangeSectionNameClick(btn) {
	// Input.
	let pp = btn.parentNode.parentNode;
	let sectionId = Number(pp.childNodes[1].childNodes[1].value);
	if (sectionId < 1) {
		console.error(Err.IdNotSet);
		return;
	}
	let newName = pp.childNodes[2].childNodes[1].value;
	if (newName.length < 1) {
		console.error(Err.NameIsNotSet);
		return;
	}

	// Work.
	let res = await changeSectionName(sectionId, newName);
	if (res.ok !== true) {
		return;
	}
	disableParentForm(btn, pp, false);
	let txt = "Section name was changed.";
	showActionSuccess(btn, txt);
	await reloadPage(true);
}

async function onBtnChangeSectionParentClick(btn) {
	// Input.
	let pp = btn.parentNode.parentNode;
	let sectionId = Number(pp.childNodes[1].childNodes[1].value);
	if (sectionId < 1) {
		console.error(Err.IdNotSet);
		return;
	}
	let newParent = Number(pp.childNodes[2].childNodes[1].value);
	if (newParent < 1) {
		console.error(Err.IdNotSet);
		return;
	}

	// Work.
	let res = await changeSectionParent(sectionId, newParent);
	if (res.ok !== true) {
		return;
	}
	disableParentForm(btn, pp, false);
	let txt = "Section was moved to a new parent.";
	showActionSuccess(btn, txt);
	await reloadPage(true);
}

async function onBtnMoveSectionUpClick(btn) {
	// Input.
	let pp = btn.parentNode.parentNode;
	let sectionId = Number(pp.childNodes[1].childNodes[1].value);
	if (sectionId < 1) {
		console.error(Err.IdNotSet);
		return;
	}

	// Work.
	let res = await moveSectionUp(sectionId);
	if (res.ok !== true) {
		return;
	}
	disableParentForm(btn, pp, true);
	let txt = "Section was moved up.";
	showActionSuccess(btn, txt);
	await reloadPage(true);
}

async function onBtnMoveSectionDownClick(btn) {
	// Input.
	let pp = btn.parentNode.parentNode;
	let sectionId = Number(pp.childNodes[1].childNodes[1].value);
	if (sectionId < 1) {
		console.error(Err.IdNotSet);
		return;
	}

	// Work.
	let res = await moveSectionDown(sectionId);
	if (res.ok !== true) {
		return;
	}
	disableParentForm(btn, pp, true);
	let txt = "Section was moved down.";
	showActionSuccess(btn, txt);
	await reloadPage(true);
}

async function onBtnDeleteSectionClick(btn) {
	// Input.
	let pp = btn.parentNode.parentNode;
	let sectionId = Number(pp.childNodes[1].childNodes[1].value);
	if (sectionId < 1) {
		console.error(Err.IdNotSet);
		return;
	}

	// Work.
	let res = await deleteSection(sectionId);
	if (res.ok !== true) {
		return;
	}
	disableParentForm(btn, pp, false);
	let txt = "Section was deleted.";
	showActionSuccess(btn, txt);
	await reloadPage(true);
}


function onForumManagerBtnProceedClick(btn) {
	let selectedActionIdx = getSelectedActionIdxBPC(btn);
	if (selectedActionIdx == null) {
		return;
	}

	btn.disabled = true;
	disableParentFormBPC(btn);

	// Draw.
	let fm = document.getElementById("forumManager");
	let fs = newFieldset();
	fm.appendChild(fs);

	let d = newDiv();
	d.className = "title";
	d.textContent = "Forum Parameters";
	fs.appendChild(d);

	switch (selectedActionIdx) {
		case 1: // Create a forum.
			d = newDiv();
			d.innerHTML = htmlInputParameterName("");
			fs.appendChild(d);
			d = newDiv();
			d.innerHTML = htmlInputParameterParent("ID of a parent section");
			fs.appendChild(d);
			d = newDiv();
			d.innerHTML = '<input type="button" class="' + ButtonClass.CreateForum + '" value="' + ButtonName.CreateForum + '" onclick="onBtnCreateForumClick(this)">';
			fs.appendChild(d);
			break;

		case 2: // Change forum's name.
			d = newDiv();
			d.innerHTML = htmlInputParameterId("ID of the changed forum");
			fs.appendChild(d);
			d = newDiv();
			d.innerHTML = htmlInputParameterNewName("New name of the forum");
			fs.appendChild(d);
			d = newDiv();
			d.innerHTML = '<input type="button" class="' + ButtonClass.ChangeForumName + '" value="' + ButtonName.ChangeForumName + '" onclick="onBtnChangeForumNameClick(this)">';
			fs.appendChild(d);
			break;

		case 3: // Change forum's parent.
			d = newDiv();
			d.innerHTML = htmlInputParameterId("ID of the changed forum");
			fs.appendChild(d);
			d = newDiv();
			d.innerHTML = htmlInputParameterNewParent("ID of the new parent");
			fs.appendChild(d);
			d = newDiv();
			d.innerHTML = '<input type="button" class="' + ButtonClass.ChangeForumParent + '" value="' + ButtonName.ChangeForumParent + '" onclick="onBtnChangeForumParentClick(this)">';
			fs.appendChild(d);
			break;

		case 4: // Move forum up & down.
			d = newDiv();
			d.innerHTML = htmlInputParameterId("ID of the moved forum");
			fs.appendChild(d);
			d = newDiv();
			d.innerHTML = '<input type="button" class="' + ButtonClass.MoveForum + '" value="' + ButtonName.MoveForumUp + '" onclick="onBtnMoveForumUpClick(this)">' +
				'<span class="subpageSpacerA">&nbsp;</span>' +
				'<input type="button" class="' + ButtonClass.MoveForum + '" value="' + ButtonName.MoveForumDown + '" onclick="onBtnMoveForumDownClick(this)">';
			fs.appendChild(d);
			break;

		case 5: // Delete a forum.
			d = newDiv();
			d.innerHTML = htmlInputParameterId("ID of the forum to delete");
			fs.appendChild(d);
			d = newDiv();
			d.innerHTML = '<input type="button" class="' + ButtonClass.DeleteForum + '" value="' + ButtonName.DeleteForum + '" onclick="onBtnDeleteForumClick(this)">';
			fs.appendChild(d);
			break;
	}
}

async function onBtnCreateForumClick(btn) {
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
	let res = await addForum(parent, name);
	let forumId = res.forumId;
	disableParentForm(btn, pp, false);
	let txt = "A forum was created. ID=" + forumId.toString() + ".";
	showActionSuccess(btn, txt);
	await reloadPage(true);
}

async function onBtnChangeForumNameClick(btn) {
	// Input.
	let pp = btn.parentNode.parentNode;
	let forumId = Number(pp.childNodes[1].childNodes[1].value);
	if (forumId < 1) {
		console.error(Err.IdNotSet);
		return;
	}
	let newName = pp.childNodes[2].childNodes[1].value;
	if (newName.length < 1) {
		console.error(Err.NameIsNotSet);
		return;
	}

	// Work.
	let res = await changeForumName(forumId, newName);
	if (res.ok !== true) {
		return;
	}
	disableParentForm(btn, pp, false);
	let txt = "Forum name was changed.";
	showActionSuccess(btn, txt);
	await reloadPage(true);
}

async function onBtnChangeForumParentClick(btn) {
	// Input.
	let pp = btn.parentNode.parentNode;
	let forumId = Number(pp.childNodes[1].childNodes[1].value);
	if (forumId < 1) {
		console.error(Err.IdNotSet);
		return;
	}
	let newParent = Number(pp.childNodes[2].childNodes[1].value);
	if (newParent < 1) {
		console.error(Err.IdNotSet);
		return;
	}

	// Work.
	let res = await changeForumSection(forumId, newParent);
	if (res.ok !== true) {
		return;
	}
	disableParentForm(btn, pp, false);
	let txt = "Forum was moved to a new parent.";
	showActionSuccess(btn, txt);
	await reloadPage(true);
}

async function onBtnMoveForumUpClick(btn) {
	// Input.
	let pp = btn.parentNode.parentNode;
	let forumId = Number(pp.childNodes[1].childNodes[1].value);
	if (forumId < 1) {
		console.error(Err.IdNotSet);
		return;
	}

	// Work.
	let res = await moveForumUp(forumId);
	if (res.ok !== true) {
		return;
	}
	disableParentForm(btn, pp, true);
	let txt = "Forum was moved up.";
	showActionSuccess(btn, txt);
	await reloadPage(true);
}

async function onBtnMoveForumDownClick(btn) {
	// Input.
	let pp = btn.parentNode.parentNode;
	let forumId = Number(pp.childNodes[1].childNodes[1].value);
	if (forumId < 1) {
		console.error(Err.IdNotSet);
		return;
	}

	// Work.
	let res = await moveForumDown(forumId);
	if (res.ok !== true) {
		return;
	}
	disableParentForm(btn, pp, true);
	let txt = "Forum was moved down.";
	showActionSuccess(btn, txt);
	await reloadPage(true);
}

async function onBtnDeleteForumClick(btn) {
	// Input.
	let pp = btn.parentNode.parentNode;
	let forumId = Number(pp.childNodes[1].childNodes[1].value);
	if (forumId < 1) {
		console.error(Err.IdNotSet);
		return;
	}

	// Work.
	let res = await deleteForum(forumId);
	if (res.ok !== true) {
		return;
	}
	disableParentForm(btn, pp, false);
	let txt = "Forum was deleted.";
	showActionSuccess(btn, txt);
	await reloadPage(true);
}


function onThreadManagerBtnProceedClick(btn) {
	let selectedActionIdx = getSelectedActionIdxBPC(btn);
	if (selectedActionIdx == null) {
		return;
	}

	btn.disabled = true;
	disableParentFormBPC(btn);

	// Draw.
	let tm = document.getElementById("threadManager");
	let fs = newFieldset();
	tm.appendChild(fs);

	let d = newDiv();
	d.className = "title";
	d.textContent = "Thread Parameters";
	fs.appendChild(d);

	switch (selectedActionIdx) {
		case 1: // Create a thread.
			d = newDiv();
			d.innerHTML = htmlInputParameterName("");
			fs.appendChild(d);
			d = newDiv();
			d.innerHTML = htmlInputParameterParent("ID of a parent forum");
			fs.appendChild(d);
			d = newDiv();
			d.innerHTML = '<input type="button" class="' + ButtonClass.CreateThread + '" value="' + ButtonName.CreateThread + '" onclick="onBtnCreateThreadClick(this)">';
			fs.appendChild(d);
			break;

		case 2: // Change thread's name.
			d = newDiv();
			d.innerHTML = htmlInputParameterId("ID of the changed thread");
			fs.appendChild(d);
			d = newDiv();
			d.innerHTML = htmlInputParameterNewName("New name of the thread");
			fs.appendChild(d);
			d = newDiv();
			d.innerHTML = '<input type="button" class="' + ButtonClass.ChangeThreadName + '" value="' + ButtonName.ChangeThreadName + '" onclick="onBtnChangeThreadNameClick(this)">';
			fs.appendChild(d);
			break;

		case 3: // Change thread's parent.
			d = newDiv();
			d.innerHTML = htmlInputParameterId("ID of the changed thread");
			fs.appendChild(d);
			d = newDiv();
			d.innerHTML = htmlInputParameterNewParent("ID of the new parent");
			fs.appendChild(d);
			d = newDiv();
			d.innerHTML = '<input type="button" class="' + ButtonClass.ChangeThreadParent + '" value="' + ButtonName.ChangeThreadParent + '" onclick="onBtnChangeThreadParentClick(this)">';
			fs.appendChild(d);
			break;

		case 4: // Move thread up & down.
			d = newDiv();
			d.innerHTML = htmlInputParameterId("ID of the moved thread");
			fs.appendChild(d);
			d = newDiv();
			d.innerHTML = '<input type="button" class="' + ButtonClass.MoveThread + '" value="' + ButtonName.MoveThreadUp + '" onclick="onBtnMoveThreadUpClick(this)">' +
				'<span class="subpageSpacerA">&nbsp;</span>' +
				'<input type="button" class="' + ButtonClass.MoveThread + '" value="' + ButtonName.MoveThreadDown + '" onclick="onBtnMoveThreadDownClick(this)">';
			fs.appendChild(d);
			break;

		case 5: // Delete a thread.
			d = newDiv();
			d.innerHTML = htmlInputParameterId("ID of the thread to delete");
			fs.appendChild(d);
			d = newDiv();
			d.innerHTML = '<input type="button" class="' + ButtonClass.DeleteThread + '" value="' + ButtonName.DeleteThread + '" onclick="onBtnDeleteThreadClick(this)">';
			fs.appendChild(d);
			break;
	}
}

async function onBtnCreateThreadClick(btn) {
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

async function onBtnChangeThreadNameClick(btn) {
	// Input.
	let pp = btn.parentNode.parentNode;
	let threadId = Number(pp.childNodes[1].childNodes[1].value);
	if (threadId < 1) {
		console.error(Err.IdNotSet);
		return;
	}
	let newName = pp.childNodes[2].childNodes[1].value;
	if (newName.length < 1) {
		console.error(Err.NameIsNotSet);
		return;
	}

	// Work.
	let res = await changeThreadName(threadId, newName);
	if (res.ok !== true) {
		return;
	}
	disableParentForm(btn, pp, false);
	let txt = "Thread name was changed.";
	showActionSuccess(btn, txt);
	await reloadPage(true);
}

async function onBtnChangeThreadParentClick(btn) {
	// Input.
	let pp = btn.parentNode.parentNode;
	let threadId = Number(pp.childNodes[1].childNodes[1].value);
	if (threadId < 1) {
		console.error(Err.IdNotSet);
		return;
	}
	let newParent = Number(pp.childNodes[2].childNodes[1].value);
	if (newParent < 1) {
		console.error(Err.IdNotSet);
		return;
	}

	// Work.
	let res = await changeThreadForum(threadId, newParent);
	if (res.ok !== true) {
		return;
	}
	disableParentForm(btn, pp, false);
	let txt = "Thread was moved to a new parent.";
	showActionSuccess(btn, txt);
	await reloadPage(true);
}

async function onBtnMoveThreadUpClick(btn) {
	// Input.
	let pp = btn.parentNode.parentNode;
	let threadId = Number(pp.childNodes[1].childNodes[1].value);
	if (threadId < 1) {
		console.error(Err.IdNotSet);
		return;
	}

	// Work.
	let res = await moveThreadUp(threadId);
	if (res.ok !== true) {
		return;
	}
	disableParentForm(btn, pp, true);
	let txt = "Thread was moved up.";
	showActionSuccess(btn, txt);
	await reloadPage(true);
}

async function onBtnMoveThreadDownClick(btn) {
	// Input.
	let pp = btn.parentNode.parentNode;
	let threadId = Number(pp.childNodes[1].childNodes[1].value);
	if (threadId < 1) {
		console.error(Err.IdNotSet);
		return;
	}

	// Work.
	let res = await moveThreadDown(threadId);
	if (res.ok !== true) {
		return;
	}
	disableParentForm(btn, pp, true);
	let txt = "Thread was moved down.";
	showActionSuccess(btn, txt);
	await reloadPage(true);
}

async function onBtnDeleteThreadClick(btn) {
	// Input.
	let pp = btn.parentNode.parentNode;
	let threadId = Number(pp.childNodes[1].childNodes[1].value);
	if (threadId < 1) {
		console.error(Err.IdNotSet);
		return;
	}

	// Work.
	let res = await deleteThread(threadId);
	if (res.ok !== true) {
		return;
	}
	disableParentForm(btn, pp, false);
	let txt = "Thread was deleted.";
	showActionSuccess(btn, txt);
	await reloadPage(true);
}


function onMessageManagerBtnProceedClick(btn) {
	let selectedActionIdx = getSelectedActionIdxBPC(btn);
	if (selectedActionIdx == null) {
		return;
	}

	btn.disabled = true;
	disableParentFormBPC(btn);

	// Draw.
	let mm = document.getElementById("messageManager");
	let fs = newFieldset();
	mm.appendChild(fs);

	let d = newDiv();
	d.className = "title";
	d.textContent = "Message Parameters";
	fs.appendChild(d);

	switch (selectedActionIdx) {
		case 1: // Create a message.
			d = newDiv();
			d.innerHTML = htmlInputParameterText("");
			fs.appendChild(d);
			d = newDiv();
			d.innerHTML = htmlInputParameterParent("ID of a parent thread");
			fs.appendChild(d);
			d = newDiv();
			d.innerHTML = '<input type="button" class="' + ButtonClass.CreateMessage + '" value="' + ButtonName.CreateMessage + '" onclick="onBtnCreateMessageClick(this)">';
			fs.appendChild(d);
			break;

		case 2: // Change message's text.
			d = newDiv();
			d.innerHTML = htmlInputParameterId("ID of the changed message");
			fs.appendChild(d);
			d = newDiv();
			d.innerHTML = htmlInputParameterNewText("New text of the message");
			fs.appendChild(d);
			d = newDiv();
			d.innerHTML = '<input type="button" class="' + ButtonClass.ChangeMessageText + '" value="' + ButtonName.ChangeMessageText + '" onclick="onBtnChangeMessageTextClick(this)">';
			fs.appendChild(d);
			break;

		case 3: // Change message's parent.
			d = newDiv();
			d.innerHTML = htmlInputParameterId("ID of the changed message");
			fs.appendChild(d);
			d = newDiv();
			d.innerHTML = htmlInputParameterNewParent("ID of the new parent");
			fs.appendChild(d);
			d = newDiv();
			d.innerHTML = '<input type="button" class="' + ButtonClass.ChangeMessageParent + '" value="' + ButtonName.ChangeMessageParent + '" onclick="onBtnChangeMessageParentClick(this)">';
			fs.appendChild(d);
			break;

		case 4: // Delete a message.
			d = newDiv();
			d.innerHTML = htmlInputParameterId("ID of the message to delete");
			fs.appendChild(d);
			d = newDiv();
			d.innerHTML = '<input type="button" class="' + ButtonClass.DeleteMessage + '" value="' + ButtonName.DeleteMessage + '" onclick="onBtnDeleteMessageClick(this)">';
			fs.appendChild(d);
			break;
	}
}

async function onBtnCreateMessageClick(btn) {
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

async function onBtnChangeMessageTextClick(btn) {
	// Input.
	let pp = btn.parentNode.parentNode;
	let messageId = Number(pp.childNodes[1].childNodes[1].value);
	if (messageId < 1) {
		console.error(Err.IdNotSet);
		return;
	}
	let newText = pp.childNodes[2].childNodes[1].value;
	if (newText.length < 1) {
		console.error(Err.TextIsNotSet);
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

async function onBtnChangeMessageParentClick(btn) {
	// Input.
	let pp = btn.parentNode.parentNode;
	let messageId = Number(pp.childNodes[1].childNodes[1].value);
	if (messageId < 1) {
		console.error(Err.IdNotSet);
		return;
	}
	let newParent = Number(pp.childNodes[2].childNodes[1].value);
	if (newParent < 1) {
		console.error(Err.IdNotSet);
		return;
	}

	// Work.
	let res = await changeMessageThread(messageId, newParent);
	if (res.ok !== true) {
		return;
	}
	disableParentForm(btn, pp, false);
	let txt = "Message was moved to a new parent.";
	showActionSuccess(btn, txt);
	await reloadPage(true);
}

async function onBtnDeleteMessageClick(btn) {
	// Input.
	let pp = btn.parentNode.parentNode;
	let messageId = Number(pp.childNodes[1].childNodes[1].value);
	if (messageId < 1) {
		console.error(Err.IdNotSet);
		return;
	}

	// Work.
	let res = await deleteMessage(messageId);
	if (res.ok !== true) {
		return;
	}
	disableParentForm(btn, pp, false);
	let txt = "Message was deleted.";
	showActionSuccess(btn, txt);
	await reloadPage(true);
}


function onNotificationManagerBtnProceedClick(btn) {
	let selectedActionIdx = getSelectedActionIdxBPC(btn);
	if (selectedActionIdx == null) {
		return;
	}

	btn.disabled = true;
	disableParentFormBPC(btn);

	// Draw.
	let nm = document.getElementById("notificationManager");
	let fs = newFieldset();
	nm.appendChild(fs);

	let d = newDiv();
	d.className = "title";
	d.textContent = "Notification Parameters";
	fs.appendChild(d);

	switch (selectedActionIdx) {
		case 1: // Create a notification.
			d = newDiv();
			d.innerHTML = htmlInputParameterText("");
			fs.appendChild(d);
			d = newDiv();
			d.innerHTML = htmlInputParameterUser("ID of a user");
			fs.appendChild(d);
			d = newDiv();
			d.innerHTML = '<input type="button" class="' + ButtonClass.CreateNotification + '" value="' + ButtonName.CreateNotification + '" onclick="onBtnCreateNotificationClick(this)">';
			fs.appendChild(d);
			break;

		case 2: // Delete a notification.
			d = newDiv();
			d.innerHTML = htmlInputParameterId("ID of the notification to delete");
			fs.appendChild(d);
			d = newDiv();
			d.innerHTML = '<input type="button" class="' + ButtonClass.DeleteNotification + '" value="' + ButtonName.DeleteNotification + '" onclick="onBtnDeleteNotificationClick(this)">';
			fs.appendChild(d);
			break;
	}
}

async function onBtnCreateNotificationClick(btn) {
	// Input.
	let pp = btn.parentNode.parentNode;
	let text = pp.childNodes[1].childNodes[1].value;
	if (text.length < 1) {
		console.error(Err.TextIsNotSet);
		return;
	}
	let userId = Number(pp.childNodes[2].childNodes[1].value);
	if (userId < 1) {
		console.error(Err.IdNotSet);
		return;
	}

	// Work.
	let res = await addNotification(userId, text);
	let notificationId = res.notificationId;
	disableParentForm(btn, pp, false);
	let txt = "A notification was created. ID=" + notificationId.toString() + ".";
	showActionSuccess(btn, txt);
	await reloadPage(true);
}

async function onBtnDeleteNotificationClick(btn) {
	// Input.
	let pp = btn.parentNode.parentNode;
	let notificationId = Number(pp.childNodes[1].childNodes[1].value);
	if (notificationId < 1) {
		console.error(Err.IdNotSet);
		return;
	}

	// Work.
	let res = await deleteNotification(notificationId);
	if (res.ok !== true) {
		return;
	}
	disableParentForm(btn, pp, false);
	let txt = "Notification was deleted.";
	showActionSuccess(btn, txt);
	await reloadPage(true);
}


function onResourceManagerBtnProceedClick(btn) {
	let selectedActionIdx = getSelectedActionIdxBPC(btn);
	if (selectedActionIdx == null) {
		return;
	}

	btn.disabled = true;
	disableParentFormBPC(btn);

	// Draw.
	let nm = document.getElementById("resourceManager");
	let fs = newFieldset();
	nm.appendChild(fs);

	let d = newDiv();
	d.className = "title";
	d.textContent = "Resource Parameters";
	fs.appendChild(d);

	switch (selectedActionIdx) {
		case 1: // Create a resource.
			d = newDiv();
			d.innerHTML = htmlInputParameterType("Resource type: 1 = text, 2 = number.");
			fs.appendChild(d);
			d = newDiv();
			d.innerHTML = htmlInputParameterText("Value: text or number.");
			fs.appendChild(d);
			d = newDiv();
			d.innerHTML = '<input type="button" class="' + ButtonClass.CreateResource + '" value="' + ButtonName.CreateResource + '" onclick="onBtnCreateResourceClick(this)">';
			fs.appendChild(d);
			break;

		case 2: // Delete a resource.
			d = newDiv();
			d.innerHTML = htmlInputParameterId("ID of the resource to delete");
			fs.appendChild(d);
			d = newDiv();
			d.innerHTML = '<input type="button" class="' + ButtonClass.DeleteResource + '" value="' + ButtonName.DeleteResource + '" onclick="onBtnDeleteResourceClick(this)">';
			fs.appendChild(d);
			break;
	}
}

async function onBtnCreateResourceClick(btn) {
	// Input.
	let pp = btn.parentNode.parentNode;
	let resourceType = Number(pp.childNodes[1].childNodes[1].value);
	let resource = new Resource(null, resourceType);
	if (!resource.checkType()) {
		console.error(Err.TypeIsNotSet);
		return;
	}
	let resourceValue = pp.childNodes[2].childNodes[1].value;
	if (resourceValue.length < 1) {
		console.error(Err.ValueNotSet);
		return;
	}
	if (!resource.setValue(resourceValue)) {
		return;
	}

	// Work.
	let res = await addResource(resource.getValue());
	let resourceId = res.resourceId;
	disableParentForm(btn, pp, false);
	let txt = "A resource was created. ID=" + resourceId.toString() + ".";
	showActionSuccess(btn, txt);
	await reloadPage(true);
}

async function onBtnDeleteResourceClick(btn) {
	// Input.
	let pp = btn.parentNode.parentNode;
	let resourceId = Number(pp.childNodes[1].childNodes[1].value);
	if (resourceId < 1) {
		console.error(Err.IdNotSet);
		return;
	}

	// Work.
	let res = await deleteResource(resourceId);
	if (res.ok !== true) {
		return;
	}
	disableParentForm(btn, pp, false);
	let txt = "Resource was deleted.";
	showActionSuccess(btn, txt);
	await reloadPage(true);
}


// Other functions.

async function fillListOfRRFA(elClass, rrfas) {
	let div = document.getElementById(elClass);
	div.innerHTML = "";
	let tbl = newTable();
	tbl.className = elClass;

	// Header.
	let tr = newTr();
	let ths = ["#", "ID", "PreRegTime", "E-Mail", "Name", "Actions"];
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

	// Cells.
	let rrfa;
	for (let i = 0; i < rrfas.length; i++) {
		rrfa = rrfas[i];

		// Fill data.
		tr = newTr();
		let tds = [];
		for (let j = 0; j < ths.length; j++) {
			tds.push("");
		}

		tds[0] = (i + 1).toString();
		tds[1] = rrfa.Id.toString();
		tds[2] = prettyTime(rrfa.PreRegTime);
		tds[3] = rrfa.Email;
		tds[4] = rrfa.Name;
		tds[5] = '<input type="button" class="' + ButtonClass.Accept + '" value="' + ButtonName.Accept + '" onclick="onBtnAcceptClick(this)">' +
			'<span class="subpageSpacerA">&nbsp;</span>' +
			'<input type="button" class="' + ButtonClass.Reject + '" value="' + ButtonName.Reject + '" onclick="onBtnRejectClick(this)">';

		let td;
		for (let j = 0; j < tds.length; j++) {
			td = newTd();

			if (j === 0) {
				td.className = "numCol";
			}

			if (j !== 5) {
				td.textContent = tds[j];
			} else {
				td.innerHTML = tds[j];
			}
			tr.appendChild(td);
		}

		tbl.appendChild(tr);
	}

	div.appendChild(tbl);
}

async function fillListOfLoggedUsers(elClass, userIds) {
	let div = document.getElementById(elClass);
	div.innerHTML = "";
	let tbl = newTable();
	tbl.className = elClass;

	// Header.
	let tr = newTr();
	let ths = [
		"#", "ID", "E-Mail", "Name", "IP Address", "Log Time", "Actions"];
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

	let columnsWithLink = [1, 2, 3];

	// Cells.
	let resA, userId, user, resB, userSession;
	for (let i = 0; i < userIds.length; i++) {
		userId = userIds[i];

		// Get user parameters.
		resA = await viewUserParameters(userId);
		user = jsonToUser(resA.user);

		// Get user session.
		resB = await getUserSession(userId);
		userSession = jsonToSession(resB.session);

		// Fill data.
		tr = newTr();
		let tds = [];
		for (let j = 0; j < ths.length; j++) {
			tds.push("");
		}

		tds[0] = (i + 1).toString();
		tds[1] = userId.toString();
		tds[2] = user.Email;
		tds[3] = user.Name;
		tds[4] = userSession.UserIPA;
		tds[5] = prettyTime(userSession.StartTime);
		tds[6] = '<input type="button" class="' + ButtonClass.LogOut + '" value="' + ButtonName.LogOut + '" onclick="onBtnLogOutClick(this)">';

		let td, url;
		for (let j = 0; j < tds.length; j++) {
			url = composeUrlForUserPageA(userId.toString());
			td = newTd();

			if (j === 0) {
				td.className = "numCol";
			}

			if (columnsWithLink.includes(j)) {
				td.innerHTML = '<a href="' + url + '">' + tds[j] + '</a>';
			} else if (j === 6) {
				td.innerHTML = tds[j];
			} else {
				td.textContent = tds[j];
			}
			tr.appendChild(td);
		}

		tbl.appendChild(tr);
	}

	div.appendChild(tbl);
}

async function fillListOfUsers(elClass, userIds, loggedUserIds) {
	let div = document.getElementById(elClass);
	div.innerHTML = "";
	let tbl = newTable();
	tbl.className = elClass;

	// Header.
	let tr = newTr();
	let ths = [
		"#", "ID", "IsLoggedIn", "E-Mail", "Name", "RegTime", "ApprovalTime",
		"LastBadLogInTime", "LastBadActionTime", "BanTime", "CanLogIn",
		"IsReader", "IsWriter", "IsAuthor", "IsModerator", "IsAdministrator"];
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

	let columnsWithLink = [1, 3, 4];

	// Cells.
	let userId, res, user, isUserLoggedIn;
	for (let i = 0; i < userIds.length; i++) {
		userId = userIds[i];

		// Get user parameters.
		res = await viewUserParameters(userId);
		user = jsonToUser(res.user);
		isUserLoggedIn = loggedUserIds.includes(userId);

		// Fill data.
		tr = newTr();
		let tds = [];
		for (let j = 0; j < ths.length; j++) {
			tds.push("");
		}

		tds[0] = (i + 1).toString();
		tds[1] = userId.toString();
		tds[2] = booleanToString(isUserLoggedIn);
		tds[3] = user.Email;
		tds[4] = user.Name;
		tds[5] = prettyTime(user.RegTime);
		tds[6] = prettyTime(user.ApprovalTime);
		tds[7] = prettyTime(user.LastBadLogInTime);
		tds[8] = prettyTime(user.LastBadActionTime);
		tds[9] = prettyTime(user.BanTime);
		tds[10] = booleanToString(user.Roles.CanLogIn);
		tds[11] = booleanToString(user.Roles.IsReader);
		tds[12] = booleanToString(user.Roles.IsWriter);
		tds[13] = booleanToString(user.Roles.IsAuthor);
		tds[14] = booleanToString(user.Roles.IsModerator);
		tds[15] = booleanToString(user.Roles.IsAdministrator);

		let td, url;
		for (let j = 0; j < tds.length; j++) {
			url = composeUrlForUserPageA(userId.toString());
			td = newTd();

			if (j === 0) {
				td.className = "numCol";
			}

			if (columnsWithLink.includes(j)) {
				td.innerHTML = '<a href="' + url + '">' + tds[j] + '</a>';
			} else {
				td.textContent = tds[j];
			}
			tr.appendChild(td);
		}

		tbl.appendChild(tr);
	}

	div.appendChild(tbl);
}

async function fillListOfResources(elClass, resourceIds) {
	let div = document.getElementById(elClass);
	div.innerHTML = "";
	let tbl = newTable();
	tbl.className = elClass;

	// Header.
	let tr = newTr();
	let ths = ["#", "ID", "Type", "Text", "Number", "ToC"];
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

	// Cells.
	let resourceId, res, resource;
	for (let i = 0; i < resourceIds.length; i++) {
		resourceId = resourceIds[i];

		// Get resource parameters.
		res = await getResource(resourceId);
		resource = jsonToResource(res.resource);

		// Fill data.
		tr = newTr();
		let tds = [];
		for (let j = 0; j < ths.length; j++) {
			tds.push("");
		}

		tds[0] = (i + 1).toString();
		tds[1] = resourceId.toString();
		tds[2] = resource.Type.toString();
		tds[3] = resource.Text;
		tds[4] = resource.Number;
		tds[5] = prettyTime(resource.ToC);

		let td;
		for (let j = 0; j < tds.length; j++) {
			td = newTd();

			if (j === 0) {
				td.className = "numCol";
			}

			td.textContent = tds[j];
			tr.appendChild(td);
		}

		tbl.appendChild(tr);
	}

	div.appendChild(tbl);
}

function fillUserPage(elClass, user, userLogInState, selfUser) {
	let div = document.getElementById(elClass);
	div.innerHTML = "";
	let tbl = newTable();
	tbl.className = elClass;

	// Header.
	let tr = newTr();
	let ths = ["#", "Field Name", "Value", "Actions"];
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

	let userId = mca_gvc.Id;
	let fieldNames = [
		"ID", "E-Mail", "Name", "PreRegTime", "RegTime", "ApprovalTime",
		"IsAdministrator", "IsModerator", "IsAuthor", "IsWriter", "IsReader",
		"CanLogIn", "LastBadLogInTime", "BanTime", "LastBadActionTime",
		"IsLoggedIn"
	];
	let fieldValues = [
		userId.toString(),
		user.Email,
		user.Name,
		prettyTime(user.PreRegTime),
		prettyTime(user.RegTime),
		prettyTime(user.ApprovalTime),
		booleanToString(user.Roles.IsAdministrator),
		booleanToString(user.Roles.IsModerator),
		booleanToString(user.Roles.IsAuthor),
		booleanToString(user.Roles.IsWriter),
		booleanToString(user.Roles.IsReader),
		booleanToString(user.Roles.CanLogIn),
		prettyTime(user.LastBadLogInTime),
		prettyTime(user.BanTime),
		prettyTime(user.LastBadActionTime),
		booleanToString(userLogInState),
	];

	// Rows.
	let tds, td, actions;
	for (let i = 0; i < fieldNames.length; i++) {
		tr = newTr();

		tds = [];
		for (let j = 0; j < ths.length; j++) {
			tds.push("");
		}

		tds[0] = (i + 1).toString();
		tds[1] = fieldNames[i];
		tds[2] = fieldValues[i];

		switch (fieldNames[i]) {
			case "IsAuthor":
				if (user.Roles.IsAuthor) {
					actions = '<input type="button" class="' + ButtonClass.DisableRole + '" value="' + ButtonName.DisableRole + '" ' +
						'onclick="onBtnDisableRoleUPClick(\'' + UserRole.Author + '\',' + userId + ')">';
				} else {
					actions = '<input type="button" class="' + ButtonClass.EnableRole + '" value="' + ButtonName.EnableRole + '" ' +
						'onclick="onBtnEnableRoleUPClick(\'' + UserRole.Author + '\',' + userId + ')">';
				}
				break;

			case "IsWriter":
				if (user.Roles.IsWriter) {
					actions = '<input type="button" class="' + ButtonClass.DisableRole + '" value="' + ButtonName.DisableRole + '" ' +
						'onclick="onBtnDisableRoleUPClick(\'' + UserRole.Writer + '\',' + userId + ')">';
				} else {
					actions = '<input type="button" class="' + ButtonClass.EnableRole + '" value="' + ButtonName.EnableRole + '" ' +
						'onclick="onBtnEnableRoleUPClick(\'' + UserRole.Writer + '\',' + userId + ')">';
				}
				break;

			case "IsReader":
				if (user.Roles.IsReader) {
					actions = '<input type="button" class="' + ButtonClass.DisableRole + '" value="' + ButtonName.DisableRole + '" ' +
						'onclick="onBtnDisableRoleUPClick(\'' + UserRole.Reader + '\',' + userId + ')">';
				} else {
					actions = '<input type="button" class="' + ButtonClass.EnableRole + '" value="' + ButtonName.EnableRole + '" ' +
						'onclick="onBtnEnableRoleUPClick(\'' + UserRole.Reader + '\',' + userId + ')">';
				}
				break;

			case "CanLogIn":
				if (selfUser.Id === user.Id) {
					actions = "";
					break;
				}
				if (user.Roles.CanLogIn) {
					actions = '<input type="button" class="' + ButtonClass.DisableRole + '" value="' + ButtonName.DisableRole + '" ' +
						'onclick="onBtnDisableRoleUPClick(\'' + UserRole.Logging + '\',' + userId + ')">';
				} else {
					actions = '<input type="button" class="' + ButtonClass.EnableRole + '" value="' + ButtonName.EnableRole + '" ' +
						'onclick="onBtnEnableRoleUPClick(\'' + UserRole.Logging + '\',' + userId + ')">';
				}
				break;

			case "IsLoggedIn":
				if (userLogInState) {
					actions = '<input type="button" class="' + ButtonClass.LogOut + '" value="' + ButtonName.LogOut + '" onclick="onBtnLogOutUPClick(' + userId + ')">';
				} else {
					actions = "";
				}
				break;

			default:
				actions = "";
		}
		tds[3] = actions;

		let jLast = tds.length - 1;
		for (let j = 0; j < tds.length; j++) {
			td = newTd();
			if (j === 0) {
				td.className = "numCol";
			}
			if (j === jLast) {
				td.innerHTML = tds[j];
			} else {
				td.textContent = tds[j];
			}
			tr.appendChild(td);
		}

		tbl.appendChild(tr);
	}

	div.appendChild(tbl);
}

function fillSectionManager(elClass) {
	let div = document.getElementById(elClass);
	div.innerHTML = "";
	let fs = newFieldset();
	div.appendChild(fs);

	let actionNames = ["Select an action", "Create a root section", "Create a normal section",
		"Change section's name", "Change section's parent", "Move section up & down", "Delete a section"];
	createRadioButtonsForActions(fs, actionNames);
	let d = newDiv();
	d.innerHTML = '<input type="button" class="' + ButtonClass.Proceed + '" value="' + ButtonName.Proceed + '" onclick="onSectionManagerBtnProceedClick(this)">';
	fs.appendChild(d);
}

function fillForumManager(elClass) {
	let div = document.getElementById(elClass);
	div.innerHTML = "";
	let fs = newFieldset();
	div.appendChild(fs);

	let actionNames = ["Select an action", "Create a forum", "Change forums's name",
		"Change forums's parent", "Move forum up & down", "Delete a forum"];
	createRadioButtonsForActions(fs, actionNames);
	let d = newDiv();
	d.innerHTML = '<input type="button" class="' + ButtonClass.Proceed + '" value="' + ButtonName.Proceed + '" onclick="onForumManagerBtnProceedClick(this)">';
	fs.appendChild(d);
}

function fillThreadManager(elClass) {
	let div = document.getElementById(elClass);
	div.innerHTML = "";
	let fs = newFieldset();
	div.appendChild(fs);

	let actionNames = ["Select an action", "Create a thread", "Change thread's name",
		"Change thread's parent", "Move thread up & down", "Delete a thread"];
	createRadioButtonsForActions(fs, actionNames);
	let d = newDiv();
	d.innerHTML = '<input type="button" class="' + ButtonClass.Proceed + '" value="' + ButtonName.Proceed + '" onclick="onThreadManagerBtnProceedClick(this)">';
	fs.appendChild(d);
}

function fillMessageManager(elClass) {
	let div = document.getElementById(elClass);
	div.innerHTML = "";
	let fs = newFieldset();
	div.appendChild(fs);

	let actionNames = ["Select an action",
		"Create a message", "Change message's text", "Change message's parent", "Delete a message"];
	createRadioButtonsForActions(fs, actionNames);
	let d = newDiv();
	d.innerHTML = '<input type="button" class="' + ButtonClass.Proceed + '" value="' + ButtonName.Proceed + '" onclick="onMessageManagerBtnProceedClick(this)">';
	fs.appendChild(d);
}

function fillNotificationManager(elClass) {
	let div = document.getElementById(elClass);
	div.innerHTML = "";
	let fs = newFieldset();
	div.appendChild(fs);

	let actionNames = ["Select an action", "Create a notification", "Delete a notification"];
	createRadioButtonsForActions(fs, actionNames);
	let d = newDiv();
	d.innerHTML = '<input type="button" class="' + ButtonClass.Proceed + '" value="' + ButtonName.Proceed + '" onclick="onNotificationManagerBtnProceedClick(this)">';
	fs.appendChild(d);
}

function fillResourceManager(elClass) {
	let div = document.getElementById(elClass);
	div.innerHTML = "";
	let fs = newFieldset();
	div.appendChild(fs);

	let actionNames = ["Select an action", "Create a resource", "Delete a resource"];
	createRadioButtonsForActions(fs, actionNames);
	let d = newDiv();
	d.innerHTML = '<input type="button" class="' + ButtonClass.Proceed + '" value="' + ButtonName.Proceed + '" onclick="onResourceManagerBtnProceedClick(this)">';
	fs.appendChild(d);
}

function addBtnBack(el) {
	let btn = newInput();
	btn.type = "button";
	btn.className = ButtonClass.BackA;
	btn.value = ButtonName.BackA;
	btn.addEventListener("click", async (e) => {
		await redirectToMainMenuA(false);
	})
	el.appendChild(btn);
}

function addTitle(el, text) {
	let div = newDiv();
	let cn = PageZoneClass.SubpageTitleA
	div.className = cn;
	div.id = cn;
	div.textContent = text;
	el.appendChild(div);
}

function addDiv(el, x) {
	let div = newDiv();
	div.className = x;
	div.id = x;
	el.appendChild(div);
}

function createRadioButtonsForActions(fs, actionNames) {
	for (let i = 0; i < actionNames.length; i++) {
		let d = newDiv();
		if (i === 0) {
			d.className = "title";
			d.textContent = actionNames[i];
		} else {
			d.innerHTML = '<input type="radio" name="action" id="action_' + i + '" value="' + actionNames[i] + '" />' +
				'<label class="action" for="action_' + i + '">' + actionNames[i] + '</label>';

		}
		fs.appendChild(d);
	}
}

function showActionSuccess(btn, txt) {
	let ppp = btn.parentNode.parentNode.parentNode;
	let d = newDiv();
	d.className = "actionSuccess";
	d.textContent = txt;
	ppp.appendChild(d);
}

function disableParentFormBPC(btn) {
	let pp = btn.parentNode.parentNode;
	for (let i = 0; i < pp.childNodes.length; i++) {
		let ch = pp.childNodes[i];
		ch.childNodes[0].disabled = true;
	}
}

function getSelectedActionIdxBPC(btn) {
	let selectedActionIdx = 0;
	let pp = btn.parentNode.parentNode;
	for (let i = 0; i < pp.childNodes.length; i++) {
		let ch = pp.childNodes[i];
		if (ch.childNodes[0].checked === true) {
			selectedActionIdx = i;
			break;
		}
	}
	if (selectedActionIdx < 1) {
		return null;
	}
	return selectedActionIdx;
}

function htmlInputParameterId(hint) {
	return htmlInputParameter("ID", hint);
}

function htmlInputParameterName(hint) {
	return htmlInputParameter("Name", hint);
}

function htmlInputParameterNewName(hint) {
	return htmlInputParameter("New Name", hint);
}

function htmlInputParameterParent(hint) {
	return htmlInputParameter("Parent", hint);
}

function htmlInputParameterNewParent(hint) {
	return htmlInputParameter("New Parent", hint);
}

function htmlInputParameterText(hint) {
	return htmlInputParameter("Text", hint);
}

function htmlInputParameterType(hint) {
	return htmlInputParameter("Type", hint);
}

function htmlInputParameterNewText(hint) {
	return htmlInputParameter("New Text", hint);
}

function htmlInputParameterUser(hint) {
	return htmlInputParameter("User", hint);
}

function htmlInputParameter(name, hint) {
	let nameLC = name.toLowerCase().replaceAll(" ", "_");

	let label;
	if (hint.length > 0) {
		label = `<label class="parameter" for="` + nameLC + `" title="` + hint + `">` + name + `</label>`;
	} else {
		label = `<label class="parameter" for="` + nameLC + `">` + name + `</label>`;
	}

	let input = `<input class="parameter" type="text" name="` + nameLC + `" id="` + nameLC + `" value="" />`;

	return label + input;
}
