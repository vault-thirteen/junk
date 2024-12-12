package dbo

import (
	"database/sql"
	"fmt"
	cdbo "github.com/vault-thirteen/SimpleBB/pkg/common/models/dbo"
)

// Indices of prepared statements.
const (
	DbPsid_CountUsersWithEmail                    = 0
	DbPsid_InsertPreRegisteredUser                = 1
	DbPsid_AttachVerificationCodeToPreRegUser     = 2
	DbPsid_CheckVerificationCodeForPreReg         = 3
	DbPsid_DeletePreRegUserIfNotApprovedByEmail   = 4
	DbPsid_ApprovePreRegUserEmail                 = 5
	DbPsid_SetPreRegUserData                      = 6
	DbPsid_ApprovePreRegUser                      = 7
	DbPsid_RegisterPreRegUserP1                   = 8
	DbPsid_RegisterPreRegUserP2                   = 9
	DbPsid_CountUsersWithName                     = 10
	DbPsid_ClearPreRegUsersTable                  = 11
	DbPsid_CountUsersWithEmailAbleToLogIn         = 12
	DbPsid_DeleteAbandonedPreSessions             = 13
	DbPsid_CountSessionsByUserEmail               = 14
	DbPsid_CountPreSessionsByUserEmail            = 15
	DbPsid_GetUserLastBadLogInTimeByEmail         = 16
	DbPsid_CreatePreSession                       = 17
	DbPsid_GetUserIdByEmail                       = 18
	DbPsid_UpdateUserLastBadLogInTimeByEmail      = 19
	DbPsid_GetPreSessionByRequestId               = 20
	DbPsid_GetUserPasswordById                    = 21
	DbPsid_DeletePreSessionByRequestId            = 22
	DbPsid_SetPreSessionCaptchaFlag               = 23
	DbPsid_SetPreSessionPasswordFlag              = 24
	DbPsid_AttachVerificationCodeToPreSession     = 25
	DbPsid_UpdatePreSessionRequestId              = 26
	DbPsid_CheckVerificationCodeForLogIn          = 27
	DbPsid_SetPreSessionVerificationFlag          = 28
	DbPsid_CreateSession                          = 29
	DbPsid_ClearSessions                          = 30
	DbPsid_GetUserById                            = 31
	DbPsid_GetSessionByUserId                     = 32
	DbPsid_DeleteSession                          = 33
	DbPsid_SaveIncident                           = 34
	DbPsid_SaveIncidentWithoutUserIPA             = 35
	DbPsid_GetListOfLoggedUsers                   = 36
	DbPsid_CountSessionsByUserId                  = 37
	DbPsid_GetUserRolesById                       = 38
	DbPsid_GetUserParametersById                  = 39
	DbPsid_SetUserRoleAuthor                      = 40
	DbPsid_SetUserRoleWriter                      = 41
	DbPsid_SetUserRoleReader                      = 42
	DbPsid_SetUserRoleCanLogIn                    = 43
	DbPsid_DeleteSessionByUserId                  = 44
	DbPsid_UpdateUserBanTime                      = 45
	DbPsid_SetPreRegUserEmailSendStatus           = 46
	DbPsid_SetPreSessionEmailSendStatus           = 47
	DbPsid_ClearPasswordChangesTable              = 48
	DbPsid_CountPasswordChangesByUserId           = 49
	DbPsid_UpdateUserLastBadActionTimeById        = 50
	DbPsid_GetUserLastBadActionTimeById           = 51
	DbPsid_CreatePasswordChangeRequest            = 52
	DbPsid_GetPasswordChangeByRequestId           = 53
	DbPsid_DeletePasswordChangeByRequestId        = 54
	DbPsid_CheckVerificationCodeForPwdChange      = 55
	DbPsid_SetPasswordChangeVFlags                = 56
	DbPsid_SetUserPassword                        = 57
	DbPsid_CountEmailChangesByUserId              = 58
	DbPsid_CreateEmailChangeRequest               = 59
	DbPsid_GetEmailChangeByRequestId              = 60
	DbPsid_DeleteEmailChangeByRequestId           = 61
	DbPsid_CheckVerificationCodesForEmailChange   = 62
	DbPsid_SetEmailChangeVFlags                   = 63
	DbPsid_SetUserEmail                           = 64
	DbPsid_SaveLogEvent                           = 65
	DbPsid_ClearEmailChangesTable                 = 66
	DbPsid_CountAllUsers                          = 67
	DbPsid_GetListOfAllUsersOnPage                = 68
	DbPsid_CountRegistrationsReadyForApproval     = 69
	DbPsid_GetListOfRegistrationsReadyForApproval = 70
	DbPsid_RejectRegistrationRequest              = 71
	DbPsid_GetUserNameById                        = 72
	DbPsid_GetListOfLoggedUsersOnPage             = 73
	DbPsid_CountLoggedUsers                       = 74
	DbPsid_GetListOfAllUsers                      = 75
)

func (dbo *DatabaseObject) makePreparedStatementQueryStrings() (qs []string) {
	var q string
	qs = make([]string, 0)

	// 0.
	q = fmt.Sprintf(`SELECT (SELECT COUNT(Email) FROM %s WHERE Email = ?) + (SELECT COUNT(Email) FROM %s WHERE Email = ?) AS n;`, dbo.tableNames.Users, dbo.tableNames.PreRegisteredUsers)
	qs = append(qs, q)

	// 1.
	q = fmt.Sprintf(`INSERT INTO %s (Email) VALUES (?);`, dbo.tableNames.PreRegisteredUsers)
	qs = append(qs, q)

	// 2.
	q = fmt.Sprintf(`UPDATE %s SET VerificationCode = ? WHERE Email = ?;`, dbo.tableNames.PreRegisteredUsers)
	qs = append(qs, q)

	// 3.
	q = fmt.Sprintf(`SELECT COUNT(Id) FROM %s WHERE Email = ? AND VerificationCode = ?;`, dbo.tableNames.PreRegisteredUsers)
	qs = append(qs, q)

	// 4.
	q = fmt.Sprintf(`DELETE FROM %s WHERE Email = ? AND IsEmailApproved = FALSE;`, dbo.tableNames.PreRegisteredUsers)
	qs = append(qs, q)

	// 5.
	q = fmt.Sprintf(`UPDATE %s SET IsEmailApproved = TRUE WHERE Email = ? AND IsEmailSent IS TRUE;`, dbo.tableNames.PreRegisteredUsers)
	qs = append(qs, q)

	// 6.
	q = fmt.Sprintf(`UPDATE %s SET NAME = ?, PASSWORD = ?, IsReadyForApproval = TRUE WHERE Email = ? AND VerificationCode = ? AND IsReadyForApproval = FALSE;`, dbo.tableNames.PreRegisteredUsers)
	qs = append(qs, q)

	// 7.
	q = fmt.Sprintf(`UPDATE %s SET IsApproved = TRUE, ApprovalTime = Now() WHERE Email = ? AND IsEmailSent IS TRUE AND IsEmailApproved IS TRUE AND IsReadyForApproval IS TRUE;`, dbo.tableNames.PreRegisteredUsers)
	qs = append(qs, q)

	// 8.
	q = fmt.Sprintf(`INSERT INTO %s (PreRegTime, Email, NAME, PASSWORD, ApprovalTime, RegTime, IsReader, CanLogIn) SELECT PreRegTime, Email, NAME, PASSWORD, ApprovalTime, Now(), TRUE, TRUE FROM %s AS pru WHERE pru.Email = ? AND pru.IsApproved IS TRUE;`, dbo.tableNames.Users, dbo.tableNames.PreRegisteredUsers)
	qs = append(qs, q)

	// 9.
	q = fmt.Sprintf(`DELETE FROM %s WHERE Email = ? AND IsApproved IS TRUE;`, dbo.tableNames.PreRegisteredUsers)
	qs = append(qs, q)

	// 10.
	q = fmt.Sprintf(`SELECT (SELECT COUNT(Name) FROM %s WHERE NAME = ?) + (SELECT COUNT(Name) FROM %s WHERE NAME = ?) AS n;`, dbo.tableNames.Users, dbo.tableNames.PreRegisteredUsers)
	qs = append(qs, q)

	// 11.
	q = fmt.Sprintf(`DELETE FROM %s WHERE IsReadyForApproval = FALSE AND PreRegTime < ?;`, dbo.tableNames.PreRegisteredUsers)
	qs = append(qs, q)

	// 12.
	q = fmt.Sprintf(`SELECT COUNT(Email) FROM %s WHERE Email = ? AND CanLogIn IS TRUE;`, dbo.tableNames.Users)
	qs = append(qs, q)

	// 13.
	q = fmt.Sprintf(`DELETE FROM %s WHERE TimeOfCreation < ?;`, dbo.tableNames.PreSessions)
	qs = append(qs, q)

	// 14.
	q = fmt.Sprintf(`SELECT COUNT(Id) FROM %s WHERE UserId = (SELECT Id FROM %s WHERE Email = ?);`, dbo.tableNames.Sessions, dbo.tableNames.Users)
	qs = append(qs, q)

	// 15.
	q = fmt.Sprintf(`SELECT COUNT(Id) FROM %s WHERE UserId = (SELECT Id FROM %s WHERE Email = ?);`, dbo.tableNames.PreSessions, dbo.tableNames.Users)
	qs = append(qs, q)

	// 16.
	q = fmt.Sprintf(`SELECT LastBadLogInTime FROM %s WHERE Email = ?;`, dbo.tableNames.Users)
	qs = append(qs, q)

	// 17.
	q = fmt.Sprintf(`INSERT INTO %s (UserId, RequestId, UserIPAB, AuthDataBytes, IsCaptchaRequired, CaptchaId) VALUES (?, ?, ?, ?, ?, ?);`, dbo.tableNames.PreSessions)
	qs = append(qs, q)

	// 18.
	q = fmt.Sprintf(`SELECT Id FROM %s WHERE Email = ?;`, dbo.tableNames.Users)
	qs = append(qs, q)

	// 19.
	q = fmt.Sprintf(`UPDATE %s SET LastBadLogInTime = Now() WHERE Email = ?;`, dbo.tableNames.Users)
	qs = append(qs, q)

	// 20.
	q = fmt.Sprintf(`SELECT Id, UserId, TimeOfCreation, RequestId, UserIPAB, AuthDataBytes, IsCaptchaRequired, CaptchaId, IsVerifiedByCaptcha, IsVerifiedByPassword, VerificationCode, IsEmailSent, IsVerifiedByEmail FROM %s WHERE RequestId = ?;`, dbo.tableNames.PreSessions)
	qs = append(qs, q)

	// 21.
	q = fmt.Sprintf(`SELECT Password FROM %s WHERE Id = ?;`, dbo.tableNames.Users)
	qs = append(qs, q)

	// 22.
	q = fmt.Sprintf(`DELETE FROM %s WHERE RequestId = ?;`, dbo.tableNames.PreSessions)
	qs = append(qs, q)

	// 23.
	q = fmt.Sprintf(`UPDATE %s SET IsVerifiedByCaptcha = ? WHERE RequestId = ? AND UserId = ?;`, dbo.tableNames.PreSessions)
	qs = append(qs, q)

	// 24.
	q = fmt.Sprintf(`UPDATE %s SET IsVerifiedByPassword = ? WHERE RequestId = ? AND UserId = ?;`, dbo.tableNames.PreSessions)
	qs = append(qs, q)

	// 25.
	q = fmt.Sprintf(`UPDATE %s SET VerificationCode = ? WHERE RequestId = ? AND UserId = ?;`, dbo.tableNames.PreSessions)
	qs = append(qs, q)

	// 26.
	q = fmt.Sprintf(`UPDATE %s SET RequestId = ? WHERE RequestId = ? AND UserId = ?;`, dbo.tableNames.PreSessions)
	qs = append(qs, q)

	// 27.
	q = fmt.Sprintf(`SELECT COUNT(Id) FROM %s WHERE RequestId = ? AND VerificationCode = ? AND IsEmailSent IS TRUE;`, dbo.tableNames.PreSessions)
	qs = append(qs, q)

	// 28.
	q = fmt.Sprintf(`UPDATE %s SET IsVerifiedByEmail = ? WHERE RequestId = ? AND UserId = ? AND IsEmailSent IS TRUE;`, dbo.tableNames.PreSessions)
	qs = append(qs, q)

	// 29.
	q = fmt.Sprintf(`INSERT INTO %s (UserId, UserIPAB) VALUES (?, ?);`, dbo.tableNames.Sessions)
	qs = append(qs, q)

	// 30.
	q = fmt.Sprintf(`DELETE FROM %s WHERE StartTime < ?;`, dbo.tableNames.Sessions)
	qs = append(qs, q)

	// 31.
	q = fmt.Sprintf(`SELECT Id, PreRegTime, Email, Name, ApprovalTime, RegTime, IsAuthor, IsWriter, IsReader, CanLogIn, LastBadLogInTime, BanTime, LastBadActionTime FROM %s WHERE Id = ?;`, dbo.tableNames.Users)
	qs = append(qs, q)

	// 32.
	q = fmt.Sprintf(`SELECT Id, UserId, StartTime, UserIPAB FROM %s WHERE UserId = ?;`, dbo.tableNames.Sessions)
	qs = append(qs, q)

	// 33.
	q = fmt.Sprintf(`DELETE FROM %s WHERE Id = ? AND UserId = ? AND UserIPAB = ?;`, dbo.tableNames.Sessions)
	qs = append(qs, q)

	// 34.
	q = fmt.Sprintf(cdbo.Query_SaveIncident, dbo.tableNames.Incidents)
	qs = append(qs, q)

	// 35.
	q = fmt.Sprintf(cdbo.Query_SaveIncidentWithoutUserIPA, dbo.tableNames.Incidents)
	qs = append(qs, q)

	// 36.
	q = fmt.Sprintf(`SELECT UserId FROM %s;`, dbo.tableNames.Sessions)
	qs = append(qs, q)

	// 37.
	q = fmt.Sprintf(`SELECT COUNT(Id) FROM %s WHERE UserId = ?;`, dbo.tableNames.Sessions)
	qs = append(qs, q)

	// 38.
	q = fmt.Sprintf(`SELECT IsAuthor, IsWriter, IsReader, CanLogIn FROM %s WHERE Id = ?;`, dbo.tableNames.Users)
	qs = append(qs, q)

	// 39.
	q = fmt.Sprintf(`SELECT Id, PreRegTime, Email, Name, ApprovalTime, RegTime, IsAuthor, IsWriter, IsReader, CanLogIn, LastBadLogInTime, BanTime, LastBadActionTime FROM %s WHERE Id = ?;`, dbo.tableNames.Users)
	qs = append(qs, q)

	// 40.
	q = fmt.Sprintf(`UPDATE %s SET IsAuthor = ? WHERE Id = ?;`, dbo.tableNames.Users)
	qs = append(qs, q)

	// 41.
	q = fmt.Sprintf(`UPDATE %s SET IsWriter = ? WHERE Id = ?;`, dbo.tableNames.Users)
	qs = append(qs, q)

	// 42.
	q = fmt.Sprintf(`UPDATE %s SET IsReader = ? WHERE Id = ?;`, dbo.tableNames.Users)
	qs = append(qs, q)

	// 43.
	q = fmt.Sprintf(`UPDATE %s SET CanLogIn = ? WHERE Id = ?;`, dbo.tableNames.Users)
	qs = append(qs, q)

	// 44.
	q = fmt.Sprintf(`DELETE FROM %s WHERE UserId = ?;`, dbo.tableNames.Sessions)
	qs = append(qs, q)

	// 45.
	q = fmt.Sprintf(`UPDATE %s SET BanTime = Now() WHERE Id = ?;`, dbo.tableNames.Users)
	qs = append(qs, q)

	// 46.
	q = fmt.Sprintf(`UPDATE %s SET IsEmailSent = ? WHERE Email = ?;`, dbo.tableNames.PreRegisteredUsers)
	qs = append(qs, q)

	// 47.
	q = fmt.Sprintf(`UPDATE %s SET IsEmailSent = ? WHERE RequestId = ? AND UserId = ?;`, dbo.tableNames.PreSessions)
	qs = append(qs, q)

	// 48.
	q = fmt.Sprintf(`DELETE FROM %s WHERE TimeOfCreation < ?;`, dbo.tableNames.PasswordChanges)
	qs = append(qs, q)

	// 49.
	q = fmt.Sprintf(`SELECT COUNT(Id) FROM %s WHERE UserId = ?;`, dbo.tableNames.PasswordChanges)
	qs = append(qs, q)

	// 50.
	q = fmt.Sprintf(`UPDATE %s SET LastBadActionTime = Now() WHERE Id = ?;`, dbo.tableNames.Users)
	qs = append(qs, q)

	// 51.
	q = fmt.Sprintf(`SELECT LastBadActionTime FROM %s WHERE Id = ?;`, dbo.tableNames.Users)
	qs = append(qs, q)

	// 52.
	q = fmt.Sprintf(`INSERT INTO %s (UserId, RequestId, UserIPAB, AuthDataBytes, IsCaptchaRequired, CaptchaId, VerificationCode, IsEmailSent, NewPassword) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?);`, dbo.tableNames.PasswordChanges)
	qs = append(qs, q)

	// 53.
	q = fmt.Sprintf(`SELECT Id, UserId, TimeOfCreation, RequestId, UserIPAB, AuthDataBytes, IsCaptchaRequired, CaptchaId, IsVerifiedByCaptcha, IsVerifiedByPassword, VerificationCode, IsEmailSent, IsVerifiedByEmail, NewPassword FROM %s WHERE RequestId = ?;`, dbo.tableNames.PasswordChanges)
	qs = append(qs, q)

	// 54.
	q = fmt.Sprintf(`DELETE FROM %s WHERE RequestId = ?;`, dbo.tableNames.PasswordChanges)
	qs = append(qs, q)

	// 55.
	q = fmt.Sprintf(`SELECT COUNT(Id) FROM %s WHERE RequestId = ? AND VerificationCode = ? AND IsEmailSent IS TRUE;`, dbo.tableNames.PasswordChanges)
	qs = append(qs, q)

	// 56.
	q = fmt.Sprintf(`UPDATE %s SET IsVerifiedByCaptcha = ?, IsVerifiedByPassword = ?, IsVerifiedByEmail = ? WHERE RequestId = ? AND UserId = ?;`, dbo.tableNames.PasswordChanges)
	qs = append(qs, q)

	// 57.
	q = fmt.Sprintf(`UPDATE %s SET PASSWORD = ? WHERE Id = ? AND Email = ?;`, dbo.tableNames.Users)
	qs = append(qs, q)

	// 58.
	q = fmt.Sprintf(`SELECT COUNT(Id) FROM %s WHERE UserId = ?;`, dbo.tableNames.EmailChanges)
	qs = append(qs, q)

	// 59.
	q = fmt.Sprintf(`INSERT INTO %s (UserId, RequestId, UserIPAB, AuthDataBytes, IsCaptchaRequired, CaptchaId, VerificationCodeOld, IsOldEmailSent, NewEmail, VerificationCodeNew, IsNewEmailSent) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?);`, dbo.tableNames.EmailChanges)
	qs = append(qs, q)

	// 60.
	q = fmt.Sprintf(`SELECT Id, UserId, TimeOfCreation, RequestId, UserIPAB, AuthDataBytes, IsCaptchaRequired, CaptchaId, IsVerifiedByCaptcha, IsVerifiedByPassword, VerificationCodeOld, IsOldEmailSent, IsVerifiedByOldEmail, NewEmail, VerificationCodeNew, IsNewEmailSent, IsVerifiedByNewEmail FROM %s WHERE RequestId = ?;`, dbo.tableNames.EmailChanges)
	qs = append(qs, q)

	// 61.
	q = fmt.Sprintf(`DELETE FROM %s WHERE RequestId = ?;`, dbo.tableNames.EmailChanges)
	qs = append(qs, q)

	// 62.
	q = fmt.Sprintf(`SELECT COUNT(Id) FROM %s WHERE RequestId = ? AND VerificationCodeOld = ? AND IsOldEmailSent IS TRUE AND VerificationCodeNew = ? AND IsNewEmailSent IS TRUE;`, dbo.tableNames.EmailChanges)
	qs = append(qs, q)

	// 63.
	q = fmt.Sprintf(`UPDATE %s SET IsVerifiedByCaptcha = ?, IsVerifiedByPassword = ?, IsVerifiedByOldEmail = ?, IsVerifiedByNewEmail = ? WHERE RequestId = ? AND UserId = ?;`, dbo.tableNames.EmailChanges)
	qs = append(qs, q)

	// 64.
	q = fmt.Sprintf(`UPDATE %s SET Email = ? WHERE Id = ? AND Email = ?;`, dbo.tableNames.Users)
	qs = append(qs, q)

	// 65.
	q = fmt.Sprintf(`INSERT INTO %s (Type, UserId, Email, UserIPAB, AdminId) VALUES (?, ?, ?, ?, ?);`, dbo.tableNames.LogEvents)
	qs = append(qs, q)

	// 66.
	q = fmt.Sprintf(`DELETE FROM %s WHERE TimeOfCreation < ?;`, dbo.tableNames.EmailChanges)
	qs = append(qs, q)

	// 67.
	q = fmt.Sprintf(`SELECT COUNT(Id) FROM %s;`, dbo.tableNames.Users)
	qs = append(qs, q)

	// 68.
	q = fmt.Sprintf(`SELECT Id FROM %s ORDER BY Id LIMIT ? OFFSET ?;`, dbo.tableNames.Users)
	qs = append(qs, q)

	// 69.
	q = fmt.Sprintf(`SELECT COUNT(Id) FROM %s WHERE IsReadyForApproval IS TRUE;`, dbo.tableNames.PreRegisteredUsers)
	qs = append(qs, q)

	// 70.
	q = fmt.Sprintf(`SELECT Id, PreRegTime, Email, Name FROM %s WHERE IsReadyForApproval IS TRUE ORDER BY PreRegTime LIMIT ? OFFSET ?;`, dbo.tableNames.PreRegisteredUsers)
	qs = append(qs, q)

	// 71.
	q = fmt.Sprintf(`DELETE FROM %s WHERE Id = ?;`, dbo.tableNames.PreRegisteredUsers)
	qs = append(qs, q)

	// 72.
	q = fmt.Sprintf(`SELECT Name FROM %s WHERE Id = ?;`, dbo.tableNames.Users)
	qs = append(qs, q)

	// 73.
	q = fmt.Sprintf(`SELECT UserId FROM %s LIMIT ? OFFSET ?;`, dbo.tableNames.Sessions)
	qs = append(qs, q)

	// 74.
	q = fmt.Sprintf(`SELECT COUNT(UserId) FROM %s;`, dbo.tableNames.Sessions)
	qs = append(qs, q)

	// 75.
	q = fmt.Sprintf(`SELECT Id FROM %s ORDER BY Id;`, dbo.tableNames.Users)
	qs = append(qs, q)

	return qs
}

func (dbo *DatabaseObject) GetPreparedStatementByIndex(i int) (ps *sql.Stmt) {
	return dbo.DatabaseObject.PreparedStatement(i)
}
