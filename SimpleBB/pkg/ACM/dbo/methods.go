package dbo

// Due to the large number of methods, they are sorted alphabetically.

import (
	"database/sql"
	base22 "github.com/vault-thirteen/SimpleBB/pkg/common/interfaces/base2"
	"github.com/vault-thirteen/SimpleBB/pkg/common/interfaces/derived1"
	"github.com/vault-thirteen/SimpleBB/pkg/common/interfaces/derived2"
	"github.com/vault-thirteen/SimpleBB/pkg/common/models/UserRoles"
	base2 "github.com/vault-thirteen/SimpleBB/pkg/common/models/base"
	"github.com/vault-thirteen/SimpleBB/pkg/common/models/complex/User"
	"github.com/vault-thirteen/SimpleBB/pkg/common/models/complex/UserParameters"
	dbo2 "github.com/vault-thirteen/SimpleBB/pkg/common/models/dbo"
	"github.com/vault-thirteen/SimpleBB/pkg/common/models/simple"
	cms "github.com/vault-thirteen/SimpleBB/pkg/common/models/sql"
	"net"
	"time"

	am "github.com/vault-thirteen/SimpleBB/pkg/ACM/models"
	ae "github.com/vault-thirteen/auxie/errors"
)

func (dbo *DatabaseObject) ApprovePreRegUser(email simple.Email) (err error) {
	var result sql.Result
	result, err = dbo.DatabaseObject.PreparedStatement(DbPsid_ApprovePreRegUser).Exec(email)
	if err != nil {
		return err
	}

	return dbo2.CheckRowsAffected(result, 1)
}

func (dbo *DatabaseObject) ApproveUserByEmail(email simple.Email) (err error) {
	var result sql.Result
	result, err = dbo.PreparedStatement(DbPsid_ApprovePreRegUserEmail).Exec(email)
	if err != nil {
		return err
	}

	return dbo2.CheckRowsAffected(result, 1)
}

func (dbo *DatabaseObject) AttachVerificationCodeToPreRegUser(email simple.Email, code simple.VerificationCode) (err error) {
	var result sql.Result
	result, err = dbo.PreparedStatement(DbPsid_AttachVerificationCodeToPreRegUser).Exec(code, email)
	if err != nil {
		return err
	}

	return dbo2.CheckRowsAffected(result, 1)
}

func (dbo *DatabaseObject) AttachVerificationCodeToPreSession(userId base2.Id, requestId simple.RequestId, code simple.VerificationCode) (err error) {
	var result sql.Result
	result, err = dbo.PreparedStatement(DbPsid_AttachVerificationCodeToPreSession).Exec(code, requestId, userId)
	if err != nil {
		return err
	}

	return dbo2.CheckRowsAffected(result, 1)
}

func (dbo *DatabaseObject) CheckVerificationCodeForLogIn(requestId simple.RequestId, code simple.VerificationCode) (ok bool, err error) {
	row := dbo.PreparedStatement(DbPsid_CheckVerificationCodeForLogIn).QueryRow(requestId, code)

	var n int
	n, err = cms.NewNonNullValueFromScannableSource[int](row)
	if err != nil {
		return false, err
	}

	if n != 1 {
		return false, nil
	}

	return true, nil
}

func (dbo *DatabaseObject) CheckVerificationCodeForPwdChange(requestId simple.RequestId, code simple.VerificationCode) (ok bool, err error) {
	row := dbo.PreparedStatement(DbPsid_CheckVerificationCodeForPwdChange).QueryRow(requestId, code)

	var n int
	n, err = cms.NewNonNullValueFromScannableSource[int](row)
	if err != nil {
		return false, err
	}

	if n != 1 {
		return false, nil
	}

	return true, nil
}

func (dbo *DatabaseObject) CheckVerificationCodeForPreReg(email simple.Email, code simple.VerificationCode) (ok bool, err error) {
	row := dbo.PreparedStatement(DbPsid_CheckVerificationCodeForPreReg).QueryRow(email, code)

	var n int
	n, err = cms.NewNonNullValueFromScannableSource[int](row)
	if err != nil {
		return false, err
	}

	if n != 1 {
		return false, nil
	}

	return true, nil
}

func (dbo *DatabaseObject) CheckVerificationCodesForEmailChange(requestId simple.RequestId, codeOld simple.VerificationCode, codeNew simple.VerificationCode) (ok bool, err error) {
	row := dbo.PreparedStatement(DbPsid_CheckVerificationCodesForEmailChange).QueryRow(requestId, codeOld, codeNew)

	var n int
	n, err = cms.NewNonNullValueFromScannableSource[int](row)
	if err != nil {
		return false, err
	}

	if n != 1 {
		return false, nil
	}

	return true, nil
}

func (dbo *DatabaseObject) CountAllUsers() (n base2.Count, err error) {
	row := dbo.PreparedStatement(DbPsid_CountAllUsers).QueryRow()

	n, err = cms.NewNonNullValueFromScannableSource[base2.Count](row)
	if err != nil {
		return dbo2.CountOnError, err
	}

	return n, nil
}

func (dbo *DatabaseObject) CountRegistrationsReadyForApproval() (n base2.Count, err error) {
	row := dbo.PreparedStatement(DbPsid_CountRegistrationsReadyForApproval).QueryRow()

	n, err = cms.NewNonNullValueFromScannableSource[base2.Count](row)
	if err != nil {
		return dbo2.CountOnError, err
	}

	return n, nil
}

func (dbo *DatabaseObject) CountEmailChangesByUserId(userId base2.Id) (n base2.Count, err error) {
	row := dbo.PreparedStatement(DbPsid_CountEmailChangesByUserId).QueryRow(userId)

	n, err = cms.NewNonNullValueFromScannableSource[base2.Count](row)
	if err != nil {
		return dbo2.CountOnError, err
	}

	return n, nil
}

func (dbo *DatabaseObject) CountLoggedUsers() (n base2.Count, err error) {
	row := dbo.PreparedStatement(DbPsid_CountLoggedUsers).QueryRow()

	n, err = cms.NewNonNullValueFromScannableSource[base2.Count](row)
	if err != nil {
		return dbo2.CountOnError, err
	}

	return n, nil
}

func (dbo *DatabaseObject) CountPasswordChangesByUserId(userId base2.Id) (n base2.Count, err error) {
	row := dbo.PreparedStatement(DbPsid_CountPasswordChangesByUserId).QueryRow(userId)

	n, err = cms.NewNonNullValueFromScannableSource[base2.Count](row)
	if err != nil {
		return dbo2.CountOnError, err
	}

	return n, nil
}

func (dbo *DatabaseObject) CountPreSessionsByUserEmail(email simple.Email) (n base2.Count, err error) {
	row := dbo.PreparedStatement(DbPsid_CountPreSessionsByUserEmail).QueryRow(email)

	n, err = cms.NewNonNullValueFromScannableSource[base2.Count](row)
	if err != nil {
		return dbo2.CountOnError, err
	}

	return n, nil
}

func (dbo *DatabaseObject) CountSessionsByUserEmail(email simple.Email) (n base2.Count, err error) {
	row := dbo.PreparedStatement(DbPsid_CountSessionsByUserEmail).QueryRow(email)

	n, err = cms.NewNonNullValueFromScannableSource[base2.Count](row)
	if err != nil {
		return dbo2.CountOnError, err
	}

	return n, nil
}

func (dbo *DatabaseObject) CountSessionsByUserId(userId base2.Id) (n base2.Count, err error) {
	row := dbo.PreparedStatement(DbPsid_CountSessionsByUserId).QueryRow(userId)

	n, err = cms.NewNonNullValueFromScannableSource[base2.Count](row)
	if err != nil {
		return dbo2.CountOnError, err
	}

	return n, nil
}

func (dbo *DatabaseObject) CountUsersWithEmailAbleToLogIn(email simple.Email) (n base2.Count, err error) {
	row := dbo.PreparedStatement(DbPsid_CountUsersWithEmailAbleToLogIn).QueryRow(email)

	n, err = cms.NewNonNullValueFromScannableSource[base2.Count](row)
	if err != nil {
		return dbo2.CountOnError, err
	}

	return n, nil
}

func (dbo *DatabaseObject) CountUsersWithEmail(email simple.Email) (n base2.Count, err error) {
	row := dbo.PreparedStatement(DbPsid_CountUsersWithEmail).QueryRow(email, email)

	n, err = cms.NewNonNullValueFromScannableSource[base2.Count](row)
	if err != nil {
		return dbo2.CountOnError, err
	}

	return n, nil
}

func (dbo *DatabaseObject) CountUsersWithName(name simple.Name) (n base2.Count, err error) {
	row := dbo.PreparedStatement(DbPsid_CountUsersWithName).QueryRow(name, name)

	n, err = cms.NewNonNullValueFromScannableSource[base2.Count](row)
	if err != nil {
		return dbo2.CountOnError, err
	}

	return n, nil
}

func (dbo *DatabaseObject) CreateEmailChangeRequest(ecr *am.EmailChange) (err error) {
	var result sql.Result
	result, err = dbo.PreparedStatement(DbPsid_CreateEmailChangeRequest).Exec(
		ecr.UserId,
		ecr.RequestId,
		ecr.UserIPAB,
		ecr.AuthDataBytes,
		ecr.IsCaptchaRequired,
		ecr.CaptchaId,
		ecr.VerificationCodeOld,
		true,
		ecr.NewEmail,
		ecr.VerificationCodeNew,
		true,
	)
	if err != nil {
		return err
	}

	return dbo2.CheckRowsAffected(result, 1)
}

func (dbo *DatabaseObject) CreatePasswordChangeRequest(pcr *am.PasswordChange) (err error) {
	var result sql.Result
	result, err = dbo.PreparedStatement(DbPsid_CreatePasswordChangeRequest).Exec(
		pcr.UserId,
		pcr.RequestId,
		pcr.UserIPAB,
		pcr.AuthDataBytes,
		pcr.IsCaptchaRequired,
		pcr.CaptchaId,
		pcr.VerificationCode,
		true,
		pcr.NewPasswordBytes,
	)
	if err != nil {
		return err
	}

	return dbo2.CheckRowsAffected(result, 1)
}

func (dbo *DatabaseObject) CreatePreSession(userId base2.Id, requestId simple.RequestId, userIPAB net.IP, pwdSalt []byte, isCaptchaRequired base2.Flag, captchaId *simple.CaptchaId) (err error) {
	var result sql.Result
	result, err = dbo.PreparedStatement(DbPsid_CreatePreSession).Exec(userId, requestId, userIPAB, pwdSalt, isCaptchaRequired, captchaId)
	if err != nil {
		return err
	}

	return dbo2.CheckRowsAffected(result, 1)
}

func (dbo *DatabaseObject) CreateSession(userId base2.Id, userIPAB net.IP) (lastInsertedId base2.Id, err error) {
	var result sql.Result
	result, err = dbo.PreparedStatement(DbPsid_CreateSession).Exec(userId, userIPAB)
	if err != nil {
		return dbo2.LastInsertedIdOnError, err
	}

	return dbo2.CheckRowsAffectedAndGetLastInsertedId(result, 1)
}

func (dbo *DatabaseObject) DeleteAbandonedPreSessions() (err error) {
	timeBorder := time.Now().Add(-time.Duration(dbo.sp.PreSessionExpirationTime) * time.Second)

	_, err = dbo.PreparedStatement(DbPsid_DeleteAbandonedPreSessions).Exec(timeBorder)
	if err != nil {
		return err
	}

	// Affected rows are not checked while they may be none or multiple.

	return nil
}

func (dbo *DatabaseObject) DeleteEmailChangeByRequestId(requestId simple.RequestId) (err error) {
	var result sql.Result
	result, err = dbo.PreparedStatement(DbPsid_DeleteEmailChangeByRequestId).Exec(requestId)
	if err != nil {
		return err
	}

	return dbo2.CheckRowsAffected(result, 1)
}

func (dbo *DatabaseObject) DeletePasswordChangeByRequestId(requestId simple.RequestId) (err error) {
	var result sql.Result
	result, err = dbo.PreparedStatement(DbPsid_DeletePasswordChangeByRequestId).Exec(requestId)
	if err != nil {
		return err
	}

	return dbo2.CheckRowsAffected(result, 1)
}

func (dbo *DatabaseObject) DeletePreRegUserIfNotApprovedByEmail(email simple.Email) (err error) {
	var result sql.Result
	result, err = dbo.PreparedStatement(DbPsid_DeletePreRegUserIfNotApprovedByEmail).Exec(email)
	if err != nil {
		return err
	}

	return dbo2.CheckRowsAffected(result, 1)
}

func (dbo *DatabaseObject) DeletePreSessionByRequestId(requestId simple.RequestId) (err error) {
	var result sql.Result
	result, err = dbo.PreparedStatement(DbPsid_DeletePreSessionByRequestId).Exec(requestId)
	if err != nil {
		return err
	}

	return dbo2.CheckRowsAffected(result, 1)
}

func (dbo *DatabaseObject) DeleteSession(sessionId base2.Id, userId base2.Id, userIPAB net.IP) (err error) {
	var result sql.Result
	result, err = dbo.PreparedStatement(DbPsid_DeleteSession).Exec(sessionId, userId, userIPAB)
	if err != nil {
		return err
	}

	return dbo2.CheckRowsAffected(result, 1)
}

func (dbo *DatabaseObject) DeleteSessionByUserId(userId base2.Id) (err error) {
	var result sql.Result
	result, err = dbo.PreparedStatement(DbPsid_DeleteSessionByUserId).Exec(userId)
	if err != nil {
		return err
	}

	return dbo2.CheckRowsAffected(result, 1)
}

func (dbo *DatabaseObject) GetEmailChangeByRequestId(requestId simple.RequestId) (ecr *am.EmailChange, err error) {
	row := dbo.PreparedStatement(DbPsid_GetEmailChangeByRequestId).QueryRow(requestId)
	return am.NewEmailChangeFromScannableSource(row)
}

func (dbo *DatabaseObject) GetListOfAllUsers() (userIds []base2.Id, err error) {
	var rows *sql.Rows
	rows, err = dbo.PreparedStatement(DbPsid_GetListOfAllUsers).Query()
	if err != nil {
		return nil, err
	}

	defer func() {
		derr := rows.Close()
		if derr != nil {
			err = ae.Combine(err, derr)
		}
	}()

	return cms.NewArrayFromScannableSource[base2.Id](rows)
}

func (dbo *DatabaseObject) GetListOfAllUsersOnPage(pageNumber base2.Count, pageSize base2.Count) (userIds []base2.Id, err error) {
	var rows *sql.Rows
	rows, err = dbo.PreparedStatement(DbPsid_GetListOfAllUsersOnPage).Query(pageSize, (pageNumber-1)*pageSize)
	if err != nil {
		return nil, err
	}

	defer func() {
		derr := rows.Close()
		if derr != nil {
			err = ae.Combine(err, derr)
		}
	}()

	return cms.NewArrayFromScannableSource[base2.Id](rows)
}

func (dbo *DatabaseObject) GetListOfLoggedUsers() (userIds []base2.Id, err error) {
	var rows *sql.Rows
	rows, err = dbo.PreparedStatement(DbPsid_GetListOfLoggedUsers).Query()
	if err != nil {
		return nil, err
	}

	defer func() {
		derr := rows.Close()
		if derr != nil {
			err = ae.Combine(err, derr)
		}
	}()

	return cms.NewArrayFromScannableSource[base2.Id](rows)
}

func (dbo *DatabaseObject) GetListOfLoggedUsersOnPage(pageNumber base2.Count, pageSize base2.Count) (userIds []base2.Id, err error) {
	var rows *sql.Rows
	rows, err = dbo.PreparedStatement(DbPsid_GetListOfLoggedUsersOnPage).Query(pageSize, (pageNumber-1)*pageSize)
	if err != nil {
		return nil, err
	}

	defer func() {
		derr := rows.Close()
		if derr != nil {
			err = ae.Combine(err, derr)
		}
	}()

	return cms.NewArrayFromScannableSource[base2.Id](rows)
}

func (dbo *DatabaseObject) GetListOfRegistrationsReadyForApproval(pageNumber base2.Count, pageSize base2.Count) (rrfas []am.RegistrationReadyForApproval, err error) {
	var rows *sql.Rows
	rows, err = dbo.PreparedStatement(DbPsid_GetListOfRegistrationsReadyForApproval).Query(pageSize, (pageNumber-1)*pageSize)
	if err != nil {
		return nil, err
	}

	defer func() {
		derr := rows.Close()
		if derr != nil {
			err = ae.Combine(err, derr)
		}
	}()

	return am.NewRegistrationReadyForApprovalArrayFromRows(rows)
}

func (dbo *DatabaseObject) GetPasswordChangeByRequestId(requestId simple.RequestId) (pcr *am.PasswordChange, err error) {
	row := dbo.PreparedStatement(DbPsid_GetPasswordChangeByRequestId).QueryRow(requestId)
	return am.NewPasswordChangeFromScannableSource(row)
}

func (dbo *DatabaseObject) GetPreSessionByRequestId(requestId simple.RequestId) (preSession *am.PreSession, err error) {
	row := dbo.PreparedStatement(DbPsid_GetPreSessionByRequestId).QueryRow(requestId)
	return am.NewPreSessionFromScannableSource(row)
}

func (dbo *DatabaseObject) GetSessionByUserId(userId base2.Id) (session *am.Session, err error) {
	row := dbo.PreparedStatement(DbPsid_GetSessionByUserId).QueryRow(userId)
	return am.NewSessionFromScannableSource(row)
}

func (dbo *DatabaseObject) GetUserNameById(userId base2.Id) (userName *simple.Name, err error) {
	row := dbo.PreparedStatement(DbPsid_GetUserNameById).QueryRow(userId)
	return cms.NewValueFromScannableSource[simple.Name](row)
}

func (dbo *DatabaseObject) GetUserById(userId base2.Id) (user derived1.IUser, err error) {
	row := dbo.PreparedStatement(DbPsid_GetUserById).QueryRow(userId)
	return u.NewUserFromScannableSource(row)
}

func (dbo *DatabaseObject) GetUserIdByEmail(email simple.Email) (userId base2.Id, err error) {
	row := dbo.PreparedStatement(DbPsid_GetUserIdByEmail).QueryRow(email)

	userId, err = cms.NewNonNullValueFromScannableSource[base2.Id](row)
	if err != nil {
		return dbo2.IdOnError, err
	}

	return userId, nil
}

func (dbo *DatabaseObject) GetUserLastBadActionTimeById(userId base2.Id) (lastBadActionTime *time.Time, err error) {
	var user = u.NewUser()
	var uParams = user.GetUserParameters()

	err = dbo.PreparedStatement(DbPsid_GetUserLastBadActionTimeById).QueryRow(userId).Scan(uParams.GetLastBadActionTimePtr())
	if err != nil {
		return nil, err
	}

	return uParams.GetLastBadActionTime(), nil
}

func (dbo *DatabaseObject) GetUserLastBadLogInTimeByEmail(email simple.Email) (lastBadLogInTime *time.Time, err error) {
	var user = u.NewUser()
	var uParams = user.GetUserParameters()

	err = dbo.PreparedStatement(DbPsid_GetUserLastBadLogInTimeByEmail).QueryRow(email).Scan(uParams.GetLastBadLogInTimePtr())
	if err != nil {
		return nil, err
	}

	return uParams.GetLastBadLogInTime(), nil
}

func (dbo *DatabaseObject) GetUserPasswordById(userId base2.Id) (password *[]byte, err error) {
	row := dbo.PreparedStatement(DbPsid_GetUserPasswordById).QueryRow(userId)
	return cms.NewValueFromScannableSource[[]byte](row)
}

func (dbo *DatabaseObject) GetUserRolesById(userId base2.Id) (roles *ur.UserRoles, err error) {
	row := dbo.PreparedStatement(DbPsid_GetUserRolesById).QueryRow(userId)
	return ur.NewUserRolesFromScannableSource(row)
}

func (dbo *DatabaseObject) InsertPreRegisteredUser(email simple.Email) (err error) {
	var result sql.Result
	result, err = dbo.PreparedStatement(DbPsid_InsertPreRegisteredUser).Exec(email)
	if err != nil {
		return err
	}

	return dbo2.CheckRowsAffected(result, 1)
}

func (dbo *DatabaseObject) RegisterPreRegUser(email simple.Email) (err error) {
	// Part 1.
	var result sql.Result
	result, err = dbo.PreparedStatement(DbPsid_RegisterPreRegUserP1).Exec(email)
	if err != nil {
		return err
	}

	err = dbo2.CheckRowsAffected(result, 1)
	if err != nil {
		return err
	}

	// Part 2.
	result, err = dbo.PreparedStatement(DbPsid_RegisterPreRegUserP2).Exec(email)
	if err != nil {
		return err
	}

	return dbo2.CheckRowsAffected(result, 1)
}

func (dbo *DatabaseObject) RejectRegistrationRequest(id base2.Id) (err error) {
	var result sql.Result
	result, err = dbo.DatabaseObject.PreparedStatement(DbPsid_RejectRegistrationRequest).Exec(id)
	if err != nil {
		return err
	}

	return dbo2.CheckRowsAffected(result, 1)
}

func (dbo *DatabaseObject) SaveIncident(module derived1.IModule, incidentType derived1.IIncidentType, email simple.Email, userIPAB net.IP) (err error) {
	var result sql.Result
	result, err = dbo.PreparedStatement(DbPsid_SaveIncident).Exec(module, incidentType, email, userIPAB)
	if err != nil {
		return err
	}

	return dbo2.CheckRowsAffected(result, 1)
}

func (dbo *DatabaseObject) SaveIncidentWithoutUserIPA(module derived1.IModule, incidentType derived1.IIncidentType, email simple.Email) (err error) {
	var result sql.Result
	result, err = dbo.PreparedStatement(DbPsid_SaveIncidentWithoutUserIPA).Exec(module, incidentType, email)
	if err != nil {
		return err
	}

	return dbo2.CheckRowsAffected(result, 1)
}

func (dbo *DatabaseObject) SaveLogEvent(logEvent derived2.ILogEvent) (err error) {
	var result sql.Result
	result, err = dbo.PreparedStatement(DbPsid_SaveLogEvent).Exec(
		logEvent.GetType(),
		logEvent.GetUserId(),
		logEvent.GetEmail(),
		logEvent.GetUserIPAB(),
		logEvent.GetAdminId(),
	)
	if err != nil {
		return err
	}

	return dbo2.CheckRowsAffected(result, 1)
}

func (dbo *DatabaseObject) SetEmailChangeVFlags(userId base2.Id, requestId simple.RequestId, ecvf *am.EmailChangeVerificationFlags) (err error) {
	var result sql.Result
	result, err = dbo.PreparedStatement(DbPsid_SetEmailChangeVFlags).Exec(
		ecvf.IsVerifiedByCaptcha,
		ecvf.IsVerifiedByPassword,
		ecvf.IsVerifiedByOldEmail,
		ecvf.IsVerifiedByNewEmail,
		requestId,
		userId,
	)
	if err != nil {
		return err
	}

	return dbo2.CheckRowsAffected(result, 1)
}

func (dbo *DatabaseObject) SetPasswordChangeVFlags(userId base2.Id, requestId simple.RequestId, pcvf *am.PasswordChangeVerificationFlags) (err error) {
	var result sql.Result
	result, err = dbo.PreparedStatement(DbPsid_SetPasswordChangeVFlags).Exec(
		pcvf.IsVerifiedByCaptcha,
		pcvf.IsVerifiedByPassword,
		pcvf.IsVerifiedByEmail,
		requestId,
		userId,
	)
	if err != nil {
		return err
	}

	return dbo2.CheckRowsAffected(result, 1)
}

func (dbo *DatabaseObject) SetPreRegUserData(email simple.Email, code simple.VerificationCode, name simple.Name, password []byte) (err error) {
	var result sql.Result
	result, err = dbo.PreparedStatement(DbPsid_SetPreRegUserData).Exec(name, password, email, code)
	if err != nil {
		return err
	}

	return dbo2.CheckRowsAffected(result, 1)
}

func (dbo *DatabaseObject) SetPreRegUserEmailSendStatus(emailSendStatus base2.Flag, email simple.Email) (err error) {
	var result sql.Result
	result, err = dbo.PreparedStatement(DbPsid_SetPreRegUserEmailSendStatus).Exec(emailSendStatus, email)
	if err != nil {
		return err
	}

	return dbo2.CheckRowsAffected(result, 1)
}

func (dbo *DatabaseObject) SetPreSessionCaptchaFlags(userId base2.Id, requestId simple.RequestId, isVerifiedByCaptcha base2.Flag) (err error) {
	var result sql.Result
	result, err = dbo.PreparedStatement(DbPsid_SetPreSessionCaptchaFlag).Exec(isVerifiedByCaptcha, requestId, userId)
	if err != nil {
		return err
	}

	return dbo2.CheckRowsAffected(result, 1)
}

func (dbo *DatabaseObject) SetPreSessionEmailSendStatus(userId base2.Id, requestId simple.RequestId, emailSendStatus base2.Flag) (err error) {
	var result sql.Result
	result, err = dbo.PreparedStatement(DbPsid_SetPreSessionEmailSendStatus).Exec(emailSendStatus, requestId, userId)
	if err != nil {
		return err
	}

	return dbo2.CheckRowsAffected(result, 1)
}

func (dbo *DatabaseObject) SetPreSessionPasswordFlag(userId base2.Id, requestId simple.RequestId, isVerifiedByPassword base2.Flag) (err error) {
	var result sql.Result
	result, err = dbo.PreparedStatement(DbPsid_SetPreSessionPasswordFlag).Exec(isVerifiedByPassword, requestId, userId)
	if err != nil {
		return err
	}

	return dbo2.CheckRowsAffected(result, 1)
}

func (dbo *DatabaseObject) SetPreSessionVerificationFlag(userId base2.Id, requestId simple.RequestId, isVerifiedByEmail base2.Flag) (err error) {
	var result sql.Result
	result, err = dbo.PreparedStatement(DbPsid_SetPreSessionVerificationFlag).Exec(isVerifiedByEmail, requestId, userId)
	if err != nil {
		return err
	}

	return dbo2.CheckRowsAffected(result, 1)
}

func (dbo *DatabaseObject) SetUserEmail(userId base2.Id, email simple.Email, newEmail simple.Email) (err error) {
	var result sql.Result
	result, err = dbo.PreparedStatement(DbPsid_SetUserEmail).Exec(newEmail, userId, email)
	if err != nil {
		return err
	}

	return dbo2.CheckRowsAffected(result, 1)
}

func (dbo *DatabaseObject) SetUserPassword(userId base2.Id, email simple.Email, newPassword []byte) (err error) {
	var result sql.Result
	result, err = dbo.PreparedStatement(DbPsid_SetUserPassword).Exec(newPassword, userId, email)
	if err != nil {
		return err
	}

	return dbo2.CheckRowsAffected(result, 1)
}

func (dbo *DatabaseObject) SetUserRoleAuthor(userId base2.Id, isRoleEnabled base2.Flag) (err error) {
	var result sql.Result
	result, err = dbo.PreparedStatement(DbPsid_SetUserRoleAuthor).Exec(isRoleEnabled, userId)
	if err != nil {
		return err
	}

	return dbo2.CheckRowsAffected(result, 1)
}

func (dbo *DatabaseObject) SetUserRoleCanLogIn(userId base2.Id, isRoleEnabled base2.Flag) (err error) {
	var result sql.Result
	result, err = dbo.PreparedStatement(DbPsid_SetUserRoleCanLogIn).Exec(isRoleEnabled, userId)
	if err != nil {
		return err
	}

	return dbo2.CheckRowsAffected(result, 1)
}

func (dbo *DatabaseObject) SetUserRoleReader(userId base2.Id, isRoleEnabled base2.Flag) (err error) {
	var result sql.Result
	result, err = dbo.PreparedStatement(DbPsid_SetUserRoleReader).Exec(isRoleEnabled, userId)
	if err != nil {
		return err
	}

	return dbo2.CheckRowsAffected(result, 1)
}

func (dbo *DatabaseObject) SetUserRoleWriter(userId base2.Id, isRoleEnabled base2.Flag) (err error) {
	var result sql.Result
	result, err = dbo.PreparedStatement(DbPsid_SetUserRoleWriter).Exec(isRoleEnabled, userId)
	if err != nil {
		return err
	}

	return dbo2.CheckRowsAffected(result, 1)
}

func (dbo *DatabaseObject) UpdatePreSessionRequestId(userId base2.Id, requestIdOld simple.RequestId, requestIdNew simple.RequestId) (err error) {
	var result sql.Result
	result, err = dbo.PreparedStatement(DbPsid_UpdatePreSessionRequestId).Exec(requestIdNew, requestIdOld, userId)
	if err != nil {
		return err
	}

	return dbo2.CheckRowsAffected(result, 1)
}

func (dbo *DatabaseObject) UpdateUserBanTime(userId base2.Id) (err error) {
	var result sql.Result
	result, err = dbo.PreparedStatement(DbPsid_UpdateUserBanTime).Exec(userId)
	if err != nil {
		return err
	}

	return dbo2.CheckRowsAffected(result, 1)
}

func (dbo *DatabaseObject) UpdateUserLastBadActionTimeById(userId base2.Id) (err error) {
	var result sql.Result
	result, err = dbo.PreparedStatement(DbPsid_UpdateUserLastBadActionTimeById).Exec(userId)
	if err != nil {
		return err
	}

	return dbo2.CheckRowsAffected(result, 1)
}

func (dbo *DatabaseObject) UpdateUserLastBadLogInTimeByEmail(email simple.Email) (err error) {
	var result sql.Result
	result, err = dbo.PreparedStatement(DbPsid_UpdateUserLastBadLogInTimeByEmail).Exec(email)
	if err != nil {
		return err
	}

	return dbo2.CheckRowsAffected(result, 1)
}

func (dbo *DatabaseObject) ViewUserParametersById(userId base2.Id) (userParameters base22.IUserParameters, err error) {
	row := dbo.PreparedStatement(DbPsid_GetUserParametersById).QueryRow(userId)
	return up.NewUserParametersFromScannableSource(row)
}
