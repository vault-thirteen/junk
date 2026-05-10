window.onpageshow = function (event) {
    if (event.persisted) {
        // Unfortunately, JavaScript does not reload a page when you click
        // "Go Back" button in your web browser. Every year programmers invent
        // a new "wheel" to fix this bug. And every year old working solutions
        // stop working and new ones are invented. This circus looks infinite,
        // but in reality it will end as soon as this evil programming language
        // dies. Please, do not support JavaScript and its developers in any
        // means possible. Please, let this evil "technology" to die.
        console.info(Msg.JavaScriptMustDie);
        window.location.reload();
    }
};

// Names of JavaScript storage variables.
const Varname = {
    // Settings.
    Settings_LoadTime: "settings_LoadTime",
    Settings_Version: "settings_Version",
    Settings_TTL: "settings_TTL",
    Settings_SiteName: "settings_SiteName",
    Settings_SiteDomain: "settings_SiteDomain",
    Settings_SessionMaxDuration: "settings_SessionMaxDuration",
    Settings_MessageEditTime: "settings_MessageEditTime",
    Settings_PageSize: "settings_PageSize",
}

// Page sections.
const id_td =
    {
        pageHeader: "pageHeader",
        pageContent: "pageContent",
        pageFooter: "pageFooter",
    };

// IDs of various sections.
const id_table =
    {
        changeEmail: "changeEmail",
        changePassword: "changePassword",
        logIn: "logIn",
        logOut: "logOut",
        register: "register",
    };

// Page content types.
const pageContent =
    {
        changeEmail: "changeEmail",
        changePassword: "changePassword",
        logIn: "logIn",
        logOut: "logOut",
        mainPage: "mainPage",
        register: "register",
    };

class InteractiveTable {
    constructor(table) {
        this.table = table;
    }

    getFieldValue(i) {
        return this.table.rows[i].children[1].children[0].value;
    }

    getFieldCheckedState(i) {
        return this.table.rows[i].children[1].children[0].checked;
    }

    setFieldValue(i, value) {
        this.table.rows[i].children[1].children[0].value = value;
    }

    disableField(i) {
        disableElement(this.table.rows[i].children[1].children[0]);
    }

    hideRow(i) {
        hideElement(this.table.rows[i]);
    }

    showRow(i) {
        showElement(this.table.rows[i]);
    }

    setImage(i, src) {
        this.table.rows[i].children[1].children[0].src = src;
    }
}

// Common basic functions.

function isNumber(x) {
    return typeof x === 'number';
}

function isNumericString(str) {
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
    console.error(Err.BooleanToString, b);
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

function countSymbol(str, symbol) {
    let x = str.replaceAll(symbol, '');
    return str.length - x.length;
}

function validateEmailAddress(x) {
    if (typeof x !== 'string') {
        return false
    }
    if (x.length < 3) {
        return false
    }
    if (countSymbol(x, '@') !== 1) {
        return false;
    }
    let atPos = x.indexOf('@');
    if ((atPos === 0) || (atPos === x.length - 1)) {
        return false
    }
    return true;
}

async function redirectPage(wait, url) {
    if (wait) {
        await sleep(delay.redirect * 1000);
    }

    document.location.href = url;
}

// Settings.

class Settings {
    constructor(version, ttl, siteName, siteDomain, sessionMaxDuration, messageEditTime, pageSize) {
        this.Version = version;
        this.TTL = ttl;
        this.SiteName = siteName;
        this.SiteDomain = siteDomain;
        this.SessionMaxDuration = sessionMaxDuration;
        this.MessageEditTime = messageEditTime;
        this.PageSize = pageSize;
    }
}

async function updateSettingsIfNeeded() {
    if (isSettingsUpdateNeeded()) {
        return await updateSettings();
    }
    return true;
}

function isSettingsUpdateNeeded() {
    let settingsLoadTimeStr = sessionStorage.getItem(Varname.Settings_LoadTime);
    if (settingsLoadTimeStr == null) {
        return true;
    }

    let settingsTtlStr = sessionStorage.getItem(Varname.Settings_TTL);
    if (settingsTtlStr == null) {
        return true;
    }
    let settingsTtl = Number(settingsTtlStr);

    let timeNow = getCurrentTimestamp();
    let settingsAge = timeNow - Number(settingsLoadTimeStr);
    if (settingsAge >= settingsTtl) {
        return true;
    }

    return false;
}

async function updateSettings() {
    let resp = await fetchSettings();
    let s = jsonToSettings(resp);
    console.info(Msg.NewSettingsReceived + s.Version.toString() + Msg.Dot);

    // Save the settings for future usage.
    saveSettings(s);
    return true;
}

async function fetchSettings() {
    let data = await fetch(path.settings);
    return await data.json();
}

function jsonToSettings(x) {
    return new Settings(
        x.version,
        x.ttl,
        x.siteName,
        x.siteDomain,
        x.sessionMaxDuration,
        x.messageEditTime,
        x.pageSize,
    );
}

function saveSettings(s) {
    sessionStorage.setItem(Varname.Settings_Version, s.Version);
    sessionStorage.setItem(Varname.Settings_TTL, s.TTL);
    sessionStorage.setItem(Varname.Settings_SiteName, s.SiteName);
    sessionStorage.setItem(Varname.Settings_SiteDomain, s.SiteDomain);
    sessionStorage.setItem(Varname.Settings_SessionMaxDuration, s.SessionMaxDuration.toString());
    sessionStorage.setItem(Varname.Settings_MessageEditTime, s.MessageEditTime.toString());
    sessionStorage.setItem(Varname.Settings_PageSize, s.PageSize.toString());

    let timeNow = getCurrentTimestamp();
    sessionStorage.setItem(Varname.Settings_LoadTime, timeNow.toString());
}

function getSettings() {
    let settingsLoadTime = sessionStorage.getItem(Varname.Settings_LoadTime);
    if (settingsLoadTime == null) {
        console.error(Err.Settings);
        return null;
    }

    return new Settings(
        sessionStorage.getItem(Varname.Settings_Version),
        sessionStorage.getItem(Varname.Settings_TTL),
        sessionStorage.getItem(Varname.Settings_SiteName),
        sessionStorage.getItem(Varname.Settings_SiteDomain),
        sessionStorage.getItem(Varname.Settings_SessionMaxDuration),
        sessionStorage.getItem(Varname.Settings_MessageEditTime),
        sessionStorage.getItem(Varname.Settings_PageSize),
    );
}

// Entry point.
async function onPageLoad() {
    // Settings initialisation.
    let ok = await updateSettingsIfNeeded();
    if (!ok) {
        return;
    }
    let settings = getSettings();

    drawPageHeader(settings);
    drawPageFooter(settings);

    let urlParams = new URLSearchParams(document.location.search);
    let ap = urlParams.get(url_parameter.action);

    switch (ap) {
        case ActionPage.LogIn:
            drawPageContent(settings, pageContent.logIn);
            return;

        case ActionPage.LogOut:
            drawPageContent(settings, pageContent.logOut);
            return;

        case ActionPage.Register:
            drawPageContent(settings, pageContent.register);
            return;

        case ActionPage.ChangePassword:
            drawPageContent(settings, pageContent.changePassword);
            return;

        case ActionPage.ChangeEmail:
            drawPageContent(settings, pageContent.changeEmail);
            return;
    }

    //TODO
    let selfRoles = await getSelfRoles();
    if (selfRoles == null) {
        if (lastHttpStatusCode === httpStatusCode.NotAuthorised) {
            await redirectPage(true, makeUrl_ActionPage(ActionPage.LogIn));
            return;
        }
        console.log(Msg.LastHttpStatusCode + lastHttpStatusCode);
        return;
    }
    console.log(selfRoles);

    switch (ap) {
        case null:
            drawPageContent(settings, pageContent.mainPage);
            return;

        default:
            console.error(Err.UnknownActionPage, ap);
            return;
    }
}

// UI functions.

function hideElement(el) {
    el.style.display = "none";
}

function showElement(el) {
    switch (el.tagName.toLowerCase()) {
        case "tr":
            el.style.display = "table-row";
            return;

        default:
            console.error(Err.UnknownElementType, el.tagName);
            return;
    }
}

function enableElement(el) {
    el.disabled = false;
}

function disableElement(el) {
    el.disabled = true;
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

function drawPageHeader(settings) {
    let ph = document.getElementById(id_td.pageHeader);
    ph.textContent = settings.SiteName + " " + "header";
}

function drawPageFooter(settings) {
    let pf = document.getElementById(id_td.pageFooter);
    pf.textContent = settings.SiteName + " " + "footer";
}

function drawPageContent(settings, contentType) {
    let pc = document.getElementById(id_td.pageContent);

    switch (contentType) {
        case pageContent.logIn:
            drawPageContent_LogIn(settings, pc);
            return;

        case pageContent.logOut:
            drawPageContent_LogOut(settings, pc);
            return;

        case pageContent.register:
            drawPageContent_Register(settings, pc);
            return;

        case pageContent.changePassword:
            drawPageContent_ChangePassword(settings, pc);
            return;

        case pageContent.changeEmail:
            drawPageContent_ChangeEmail(settings, pc);
            return;

        default:
            console.error(Err.UnknownPageContentType, contentType);
            return;
    }
}

function drawPageContent_ChangeEmail(settings, pc) {
    pc.innerHTML = `
<table id="changeEmail">
    <tr>
        <td colspan="2">
            Fill the form below to change your e-mail address. <br>
            <br>
        </td>
    </tr>
    <tr>
        <td class="fieldName">Current Password</td>
        <td>
            <input type="password" name="cur_pwd"/>
        </td>
    </tr>
    <tr>
        <td class="fieldName">New E-Mail</td>
        <td>
            <input type="text" name="new_email"/>
        </td>
    </tr>
    <tr>
        <td colspan="2" class="taCenter">
            <input type="button" name="change_email_proceed_1" value=" Proceed " onClick="on_change_email_proceed_1_click(this)" class="proceedButton"/>
        </td>
    </tr>
    <tr>
        <td class="fieldName">Captcha Question</td>
        <td>
            <img alt="captcha_question" src=""/>
        </td>
    </tr>
    <tr>
         <td class="fieldName">Captcha Answer</td>
        <td>
            <input type="text" name="captcha_answer"/>
        </td>
    </tr>
    <tr>
         <td class="fieldName">Verification Code (old mail)</td>
        <td>
            <input type="text" name="verification_code_old"/>
        </td>
    </tr>
    <tr>
         <td class="fieldName">Verification Code (new mail)</td>
        <td>
            <input type="text" name="verification_code_new"/>
        </td>
    </tr>
    <tr>
        <td colspan="2" class="taCenter">
            <input type="button" name="change_email_proceed_2" value=" Proceed " onClick="on_change_email_proceed_2_click(this)" class="proceedButton"/>
        </td>
    </tr>
    <tr>
         <td class="fieldName">Request ID</td>
        <td>
            <input type="text" name="request_id"/>
        </td>
    </tr>
    <tr>
         <td class="fieldName">Auth Data</td>
        <td>
            <input type="text" name="auth_data"/>
        </td>
    </tr>
</table>`;

    let tbl = document.getElementById(id_table.changeEmail);
    for (let i = 0; i < tbl.rows.length; i++) {
        if (i > 3) {
            hideElement(tbl.rows[i]);
        }
    }
}

function drawPageContent_ChangePassword(settings, pc) {
    pc.innerHTML = `
<table id="changePassword">
    <tr>
        <td colspan="2">
            Fill the form below to change your password. <br>
            <br>
        </td>
    </tr>
    <tr>
         <td class="fieldName">Current Password</td>
        <td>
            <input type="password" name="cur_pwd"/>
        </td>
    </tr>
    <tr>
         <td class="fieldName">New Password</td>
        <td>
            <input type="password" name="new_pwd_1"/>
        </td>
    </tr>
    <tr>
         <td class="fieldName">New Password (again)</td>
        <td>
            <input type="password" name="new_pwd_2"/>
        </td>
    </tr>
    <tr>
        <td colspan="2" class="taCenter">
            <input type="button" name="change_password_proceed_1" value=" Proceed " onClick="on_change_password_proceed_1_click(this)" class="proceedButton"/>
        </td>
    </tr>
    <tr>
         <td class="fieldName">Captcha Question</td>
        <td>
            <img alt="captcha_question" src=""/>
        </td>
    </tr>
    <tr>
         <td class="fieldName">Captcha Answer</td>
        <td>
            <input type="text" name="captcha_answer"/>
        </td>
    </tr>
    <tr>
         <td class="fieldName">Verification Code</td>
        <td>
            <input type="text" name="verification_code"/>
        </td>
    </tr>
    <tr>
        <td colspan="2" class="taCenter">
            <input type="button" name="change_password_proceed_2" value=" Proceed " onClick="on_change_password_proceed_2_click(this)" class="proceedButton"/>
        </td>
    </tr>
    <tr>
         <td class="fieldName">Request ID</td>
        <td>
            <input type="text" name="request_id"/>
        </td>
    </tr>
    <tr>
         <td class="fieldName">Auth Data</td>
        <td>
            <input type="text" name="auth_data"/>
        </td>
    </tr>
</table>`;

    let tbl = document.getElementById(id_table.changePassword);
    for (let i = 0; i < tbl.rows.length; i++) {
        if (i > 4) {
            hideElement(tbl.rows[i]);
        }
    }
}

function drawPageContent_LogIn(settings, pc) {
    pc.innerHTML = `
<table id="logIn">
    <tr>
        <td colspan="2">
            In order to use this website, you must be logged into the system. <br>
            If you have no account, <a href="/?a=register">click here</a> to register one. <br>
            If you have an account, log in using the form below. <br>
            <br>
        </td>
    </tr>
    <tr>
         <td class="fieldName">E-Mail</td>
        <td>
            <input type="text" name="user_email"/>
        </td>
    </tr>
    <tr>
        <td colspan="2" class="taCenter">
            <input type="button" name="log_in_proceed_1" value=" Proceed " onClick="on_log_in_proceed_1_click(this)" class="proceedButton"/>
        </td>
    </tr>
    <tr>
         <td class="fieldName">Captcha Question</td>
        <td>
            <img alt="captcha_question" src=""/>
        </td>
    </tr>
    <tr>
         <td class="fieldName">Captcha Answer</td>
        <td>
            <input type="text" name="captcha_answer"/>
        </td>
    </tr>
    <tr>
         <td class="fieldName">Verification Code</td>
        <td>
            <input type="text" name="verification_code"/>
        </td>
    </tr>
    <tr>
         <td class="fieldName">Password</td>
        <td>
            <input type="password" name="user_pwd"/>
        </td>
    </tr>
    <tr>
        <td colspan="2" class="taCenter">
            <input type="button" name="log_in_proceed_2" value=" Proceed " onClick="on_log_in_proceed_2_click(this)" class="proceedButton"/>
        </td>
    </tr>
    <tr>
         <td class="fieldName">Request ID</td>
        <td>
            <input type="text" name="request_id"/>
        </td>
    </tr>
    <tr>
         <td class="fieldName">Auth Data</td>
        <td>
            <input type="text" name="auth_data"/>
        </td>
    </tr>
</table>`;

    let tbl = document.getElementById(id_table.logIn);
    for (let i = 0; i < tbl.rows.length; i++) {
        if (i > 2) {
            hideElement(tbl.rows[i]);
        }
    }
}

function drawPageContent_LogOut(settings, pc) {
    pc.innerHTML = `
<table id="logOut">
    <tr>
        <td colspan="2">
            If you really want to log out of the system, confirm you decision. <br>
            <br>
        </td>
    </tr>
    <tr>
        <td colspan="2" class="taCenter">
            <input type="button" name="log_out_proceed_1" value=" Proceed " onClick="on_log_out_proceed_1_click(this)" class="proceedButton"/>
        </td>
    </tr>
    <tr>
         <td class="fieldName">Request ID</td>
        <td>
            <input type="text" name="request_id"/>
        </td>
    </tr>
    <tr>
         <td class="fieldName">Are you sure</td>
        <td>
            <input type="checkbox" name="are_you_sure"/>
        </td>
    </tr>
    <tr>
        <td colspan="2" class="taCenter">
            <input type="button" name="log_in_proceed_2" value=" Proceed " onClick="on_log_out_proceed_2_click(this)" class="proceedButton"/>
        </td>
    </tr>
</table>`;

    let tbl = document.getElementById(id_table.logOut);
    for (let i = 0; i < tbl.rows.length; i++) {
        if (i > 1) {
            hideElement(tbl.rows[i]);
        }
    }
}

function drawPageContent_Register(settings, pc) {
    pc.innerHTML = `
<table id="register">
    <tr>
        <td colspan="2">
            Fill the form below to register a new account. <br>
            <br>
        </td>
    </tr>
    <tr>
         <td class="fieldName">Name</td>
        <td>
            <input type="text" name="user_name"/>
        </td>
    </tr>
    <tr>
         <td class="fieldName">E-Mail</td>
        <td>
            <input type="text" name="user_email"/>
        </td>
    </tr>
    <tr>
         <td class="fieldName">Password</td>
        <td>
            <input type="password" name="user_pwd_1"/>
        </td>
    </tr>
    <tr>
         <td class="fieldName">Password (again)</td>
        <td>
            <input type="password" name="user_pwd_2"/>
        </td>
    </tr>
    <tr>
        <td colspan="2" class="taCenter">
            <input type="button" name="register_proceed_1" value=" Proceed " onClick="on_register_proceed_1_click(this)" class="proceedButton"/>
        </td>
    </tr>
    <tr>
         <td class="fieldName">Captcha Question</td>
        <td>
            <img alt="captcha_question" src=""/>
        </td>
    </tr>
    <tr>
         <td class="fieldName">Captcha Answer</td>
        <td>
            <input type="text" name="captcha_answer"/>
        </td>
    </tr>
    <tr>
         <td class="fieldName">Verification Code</td>
        <td>
            <input type="text" name="verification_code"/>
        </td>
    </tr>
    <tr>
        <td colspan="2" class="taCenter">
            <input type="button" name="register_proceed_2" value=" Proceed " onClick="on_register_proceed_2_click(this)" class="proceedButton"/>
        </td>
    </tr>
    <tr>
         <td class="fieldName">Request ID</td>
        <td>
            <input type="text" name="request_id"/>
        </td>
    </tr>
</table>`;

    let tbl = document.getElementById(id_table.register);
    for (let i = 0; i < tbl.rows.length; i++) {
        if (i > 5) {
            hideElement(tbl.rows[i]);
        }
    }
}

// Event handlers.

async function on_change_email_proceed_1_click(e) {
    // Get input data.
    let it = new InteractiveTable(e.parentNode.parentNode.parentNode);
    let newEmail = it.getFieldValue(2);

    // Check data.
    let ok = validateEmailAddress(newEmail);
    if (!ok) {
        console.error(Err.EmailAddressIsNotValid);
        return;
    }

    // Work.
    let res = await startEmailChange(newEmail);
    if (res == null) {
        return;
    }

    // Show results.
    it.disableField(1); // Current password.
    it.disableField(2); // New e-mail.
    it.hideRow(3); // First button.
    it.showRow(4); // Captcha image.
    it.setImage(4, makeUrl_CaptchaImage(res.captchaId));
    it.showRow(5); // Captcha answer.
    it.showRow(6); // Verification code (old mail).
    it.showRow(7); // Verification code (new mail).
    it.showRow(8); // Second button.
    it.setFieldValue(9, res.requestId); // RequestId.
    it.disableField(9);
    it.hideRow(9);
    it.setFieldValue(10, res.authData); // AuthData.
    it.disableField(10);
    it.hideRow(10);
}

async function on_change_email_proceed_2_click(e) {
    // Get input data.
    let it = new InteractiveTable(e.parentNode.parentNode.parentNode);
    let curPassword = it.getFieldValue(1);
    let captchaAnswer = it.getFieldValue(5);
    let vCodeOldEmail = it.getFieldValue(6);
    let vCodeNewEmail = it.getFieldValue(7);
    let requestId = it.getFieldValue(9);
    let saltBA = base64ToByteArray(it.getFieldValue(10));

    // Check data.
    if (curPassword.length === 0) {
        console.error(Err.PasswordIsNotSet);
        return;
    }
    if (!isPasswordAllowed(curPassword)) {
        console.error(Err.PasswordIsNotAllowed);
        return;
    }
    if (captchaAnswer.length === 0) {
        console.error(Err.CaptchaAnswerIsNotSet);
        return;
    }
    if (vCodeOldEmail.length === 0) {
        console.error(Err.VerificationCodeIsNotSet);
        return;
    }
    if (vCodeNewEmail.length === 0) {
        console.error(Err.VerificationCodeIsNotSet);
        return;
    }
    if (requestId.length === 0) {
        console.error(Err.RequestIdIsNotSet);
        return;
    }

    // Prepare data.
    let keyBA = makeHashKey(curPassword, saltBA);
    let authChallengeResponse = byteArrayToBase64(keyBA);

    // Work.
    let res = await confirmEmailChange(requestId, captchaAnswer, vCodeOldEmail, vCodeNewEmail, authChallengeResponse);
    if (res == null) {
        return;
    }
    if (!isSuccessfulResult(res)) {
        return;
    }

    // Show results.
    await redirectPage(true, path.root);
}

async function on_change_password_proceed_1_click(e) {
    // Get input data.
    let it = new InteractiveTable(e.parentNode.parentNode.parentNode);
    let newPassword = it.getFieldValue(2);
    let newPassword2 = it.getFieldValue(3);

    // Check data.
    if (newPassword !== newPassword2) {
        console.error(Err.PasswordIsDifferent);
        return;
    }
    if (!isPasswordAllowed(newPassword)) {
        console.error(Err.PasswordIsNotAllowed);
        return;
    }

    // Work.
    let res = await startPasswordChange(newPassword, newPassword2);
    if (res == null) {
        return;
    }

    // Show results.
    it.disableField(1); // Current password.
    it.disableField(2); // New password.
    it.disableField(3); // New password (again).
    it.hideRow(4); // First button.
    it.showRow(5); // Captcha image.
    it.setImage(5, makeUrl_CaptchaImage(res.captchaId));
    it.showRow(6); // Captcha answer.
    it.showRow(7); // Verification code.
    it.showRow(8); // Second button.
    it.setFieldValue(9, res.requestId); // RequestId.
    it.disableField(9);
    it.hideRow(9);
    it.setFieldValue(10, res.authData); // AuthData.
    it.disableField(10);
    it.hideRow(10);
}

async function on_change_password_proceed_2_click(e) {
    // Get input data.
    let it = new InteractiveTable(e.parentNode.parentNode.parentNode);
    let curPassword = it.getFieldValue(1);
    let captchaAnswer = it.getFieldValue(6);
    let vCode = it.getFieldValue(7);
    let requestId = it.getFieldValue(9);
    let saltBA = base64ToByteArray(it.getFieldValue(10));

    // Check data.
    if (curPassword.length === 0) {
        console.error(Err.PasswordIsNotSet);
        return;
    }
    if (!isPasswordAllowed(curPassword)) {
        console.error(Err.PasswordIsNotAllowed);
        return;
    }
    if (captchaAnswer.length === 0) {
        console.error(Err.CaptchaAnswerIsNotSet);
        return;
    }
    if (vCode.length === 0) {
        console.error(Err.VerificationCodeIsNotSet);
        return;
    }
    if (requestId.length === 0) {
        console.error(Err.RequestIdIsNotSet);
        return;
    }

    // Prepare data.
    let keyBA = makeHashKey(curPassword, saltBA);
    let authChallengeResponse = byteArrayToBase64(keyBA);

    // Work.
    let res = await confirmPasswordChange(requestId, captchaAnswer, vCode, authChallengeResponse);
    if (res == null) {
        return;
    }
    if (!isSuccessfulResult(res)) {
        return;
    }

    // Show results.
    await redirectPage(true, path.root);
}

async function on_log_in_proceed_1_click(e) {
    // Get input data.
    let it = new InteractiveTable(e.parentNode.parentNode.parentNode);
    let email = it.getFieldValue(1);

    // Check data.
    let ok = validateEmailAddress(email);
    if (!ok) {
        console.error(Err.EmailAddressIsNotValid);
        return;
    }

    // Work.
    let res = await startLogIn(email);
    if (res == null) {
        return;
    }

    // Show results.
    it.disableField(1); // E-mail.
    it.hideRow(2); // First button.
    it.showRow(3); // Captcha image.
    it.setImage(3, makeUrl_CaptchaImage(res.captchaId));
    it.showRow(4); // Captcha answer.
    it.showRow(5); // Verification code.
    it.showRow(6); // Password.
    it.showRow(7); // Second button.
    it.setFieldValue(8, res.requestId); // RequestId.
    it.disableField(8);
    it.hideRow(8);
    it.setFieldValue(9, res.authData); // AuthData.
    it.disableField(9);
    it.hideRow(9);
}

async function on_log_in_proceed_2_click(e) {
    // Get input data.
    let it = new InteractiveTable(e.parentNode.parentNode.parentNode);
    let captchaAnswer = it.getFieldValue(4);
    let vCode = it.getFieldValue(5);
    let pwd = it.getFieldValue(6);
    let requestId = it.getFieldValue(8);
    let saltBA = base64ToByteArray(it.getFieldValue(9));

    // Check data.
    if (captchaAnswer.length === 0) {
        console.error(Err.CaptchaAnswerIsNotSet);
        return;
    }
    if (vCode.length === 0) {
        console.error(Err.VerificationCodeIsNotSet);
        return;
    }
    if (pwd.length === 0) {
        console.error(Err.PasswordIsNotSet);
        return;
    }
    if (!isPasswordAllowed(pwd)) {
        console.error(Err.PasswordIsNotAllowed);
        return;
    }
    if (requestId.length === 0) {
        console.error(Err.RequestIdIsNotSet);
        return;
    }

    // Prepare data.
    let keyBA = makeHashKey(pwd, saltBA);
    let authChallengeResponse = byteArrayToBase64(keyBA);

    // Work.
    let res = await confirmLogIn(requestId, captchaAnswer, vCode, authChallengeResponse);
    if (res == null) {
        return;
    }
    if (!isSuccessfulResult(res)) {
        return;
    }

    // Show results.
    await redirectPage(true, path.root);
}

async function on_log_out_proceed_1_click(e) {
    // Get input data.
    let it = new InteractiveTable(e.parentNode.parentNode.parentNode);

    // Work.
    let res = await startLogOut();
    if (res == null) {
        return;
    }

    // Show results.
    it.hideRow(1); // First button.
    it.setFieldValue(2, res.requestId); // RequestId.
    it.disableField(2);
    it.hideRow(2);
    it.showRow(3); // Are you sure.
    it.showRow(4); // Second button.
}

async function on_log_out_proceed_2_click(e) {
    // Get input data.
    let it = new InteractiveTable(e.parentNode.parentNode.parentNode);
    let requestId = it.getFieldValue(2);
    let areYouSure = it.getFieldCheckedState(3);

    // Check data.
    if (requestId.length === 0) {
        console.error(Err.RequestIdIsNotSet);
        return;
    }
    if (!areYouSure) {
        return;
    }

    // Work.
    let res = await confirmLogOut(requestId, areYouSure);
    if (res == null) {
        return;
    }
    if (!isSuccessfulResult(res)) {
        return;
    }

    // Show results.
    await redirectPage(true, path.root);
}

async function on_register_proceed_1_click(e) {
    // Get input data.
    let it = new InteractiveTable(e.parentNode.parentNode.parentNode);
    let name = it.getFieldValue(1);
    let email = it.getFieldValue(2);
    let password = it.getFieldValue(3);
    let password2 = it.getFieldValue(4);

    // Check data.
    if (name.length === 0) {
        console.error(Err.NameIsNotSet);
        return;
    }
    let ok = validateEmailAddress(email);
    if (!ok) {
        console.error(Err.EmailAddressIsNotValid);
        return;
    }
    if (password !== password2) {
        console.error(Err.PasswordIsDifferent);
        return;
    }
    if (!isPasswordAllowed(password)) {
        console.error(Err.PasswordIsNotAllowed);
        return;
    }

    // Work.
    let res = await startRegistration(name, email, password);
    if (res == null) {
        return;
    }

    // Show results.
    it.disableField(1); // Name.
    it.disableField(2); // E-mail.
    it.disableField(3); // Password.
    it.disableField(4); // Password #2.
    it.hideRow(5); // First button.
    it.showRow(6); // Captcha image.
    it.setImage(6, makeUrl_CaptchaImage(res.captchaId));
    it.showRow(7); // Captcha answer.
    it.showRow(8); // Verification code.
    it.showRow(9); // Second button.
    it.setFieldValue(10, res.requestId); // RequestId.
    it.disableField(10);
    it.hideRow(10);
}

async function on_register_proceed_2_click(e) {
    // Get input data.
    let it = new InteractiveTable(e.parentNode.parentNode.parentNode);
    let captchaAnswer = it.getFieldValue(7);
    let vCode = it.getFieldValue(8);
    let requestId = it.getFieldValue(10);

    // Check data.
    if (captchaAnswer.length === 0) {
        console.error(Err.CaptchaAnswerIsNotSet);
        return;
    }
    if (vCode.length === 0) {
        console.error(Err.VerificationCodeIsNotSet);
        return;
    }
    if (requestId.length === 0) {
        console.error(Err.RequestIdIsNotSet);
        return;
    }

    // Work.
    let res = await confirmRegistration(requestId, captchaAnswer, vCode);
    if (res == null) {
        return;
    }
    if (!isSuccessfulResult(res)) {
        return;
    }

    // Show results.
    await redirectPage(true, path.root);
}
