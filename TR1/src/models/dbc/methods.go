package dbc

import (
	"errors"
	"gorm.io/gorm"
	"time"

	"github.com/vault-thirteen/BytePackedPassword"
	"github.com/vault-thirteen/TR1/src/models/common"
)

const (
	CountOnError = -1
	TmpForumPos  = -1
)

const (
	UserRoleName_Author = "author"
	UserRoleName_Writer = "writer"
	UserRoleName_Reader = "reader"
	UserRoleName_User   = "user"
)

const (
	Err_InvalidUserRoleName = "invalid user role name"
	Err_NoRowsUpdated       = "no rows updated"
)

// Common methods.

func (dbc *DbController) countAllItems(model *gorm.DB) (n int, err error) {
	var n64 int64
	tx := model.Count(&n64)
	if tx.Error != nil {
		return CountOnError, tx.Error
	}
	return int(n64), nil
}
func (dbc *DbController) listItemsOnPage(model *gorm.DB, page int, dst any) (err error) {
	tx := model.Limit(dbc.pageSize).Offset((page - 1) * dbc.pageSize).Find(dst)
	if tx.Error != nil {
		return tx.Error
	}
	return nil
}
func (dbc *DbController) listItemsOnPageWithTotalCount(model *gorm.DB, page int, dst any) (totalCount int, err error) {
	totalCount, err = dbc.countAllItems(model)
	if err != nil {
		return CountOnError, err
	}

	err = dbc.listItemsOnPage(model, page, dst)
	if err != nil {
		return CountOnError, err
	}

	return totalCount, nil
}

// User registration.

func (dbc *DbController) IsUserNameFree(userName string) (isFree bool, err error) {
	var n int64
	tx := dbc.db.Model(&cm.User{}).Where("name = ?", userName).Count(&n)
	if tx.Error != nil {
		return false, tx.Error
	}
	if n > 0 {
		return false, nil
	}
	return true, nil
}
func (dbc *DbController) IsUserEmailFree(userEmail string) (isFree bool, err error) {
	var n int64
	tx := dbc.db.Model(&cm.User{}).Where("email = ?", userEmail).Count(&n)
	if tx.Error != nil {
		return false, tx.Error
	}
	if n > 0 {
		return false, nil
	}
	return true, nil
}
func (dbc *DbController) ExistsRegistrationRequestWithUserName(userName string) (exists bool, err error) {
	var n int64
	tx := dbc.db.Model(&cm.RegistrationRequest{}).Where("user_name = ?", userName).Count(&n)
	if tx.Error != nil {
		return false, tx.Error
	}
	if n > 0 {
		return true, nil
	}
	return false, nil
}
func (dbc *DbController) ExistsRegistrationRequestWithUserEmail(userEmail string) (exists bool, err error) {
	var n int64
	tx := dbc.db.Model(&cm.RegistrationRequest{}).Where("user_email = ?", userEmail).Count(&n)
	if tx.Error != nil {
		return false, tx.Error
	}
	if n > 0 {
		return true, nil
	}
	return false, nil
}
func (dbc *DbController) CreateRegistrationRequest(rr cm.RegistrationRequest) (err error) {
	tx := dbc.db.Create(&rr)
	if tx.Error != nil {
		return tx.Error
	}
	return nil
}
func (dbc *DbController) FindRegistrationRequestNRFA(rr *cm.RegistrationRequest) (err error) {
	tx := dbc.db.First(rr, "request_id = ? AND is_ready_for_approval = ?", rr.RequestId, false)
	if tx.Error != nil {
		return tx.Error
	}
	return nil
}
func (dbc *DbController) MarkRegistrationRequestAsReadyForApproval(rr *cm.RegistrationRequest) (err error) {
	tx := dbc.db.Model(&rr).Update("is_ready_for_approval", true)
	if tx.Error != nil {
		return tx.Error
	}
	return nil
}
func (dbc *DbController) CreateUser(user *cm.User, password string) (err error) {
	var pwd = &cm.Password{
		UserId: user.Id,
	}

	pwd.Bytes, err = bpp.PackSymbols([]rune(password))
	if err != nil {
		return err
	}

	user.Password = pwd

	tx := dbc.db.Create(user)
	if tx.Error != nil {
		return tx.Error
	}

	return nil
}
func (dbc *DbController) DeleteRegistrationRequestRFA(rr *cm.RegistrationRequest) (err error) {
	tx := dbc.db.Where("is_ready_for_approval = ?", true).Delete(&rr)
	if tx.Error != nil {
		return tx.Error
	}
	return nil
}
func (dbc *DbController) DeleteRegistrationRequestNRFA(rr *cm.RegistrationRequest) (err error) {
	tx := dbc.db.Where("is_ready_for_approval = ?", false).Delete(&rr)
	if tx.Error != nil {
		return tx.Error
	}
	return nil
}
func (dbc *DbController) GetFirstOutdatedRegistrationRequest(edgeTime time.Time) (rrs []cm.RegistrationRequest, err error) {
	tx := dbc.db.Limit(1).Where("is_ready_for_approval = ? AND created_at <= ?", false, edgeTime).Find(&rrs)
	if tx.Error != nil {
		return nil, tx.Error
	}
	return rrs, nil
}
func (dbc *DbController) ListRegistrationRequestsRFA(page int) (rrs []cm.RegistrationRequest, totalCount int, err error) {
	model := dbc.db.Model(&cm.RegistrationRequest{}).Where("is_ready_for_approval = ?", true).Omit("UserPassword")

	totalCount, err = dbc.listItemsOnPageWithTotalCount(model, page, &rrs)
	if err != nil {
		return nil, totalCount, err
	}

	return rrs, totalCount, nil
}
func (dbc *DbController) GetRegistrationRequestRFA(userEmail string, rr *cm.RegistrationRequest) (err error) {
	tx := dbc.db.First(rr, "is_ready_for_approval = ? AND user_email = ?", true, userEmail)
	if tx.Error != nil {
		return tx.Error
	}
	return nil
}

// User logging in.

func (dbc *DbController) ExistsLogInRequestWithUserEmail(userEmail string) (exists bool, err error) {
	var n int64
	tx := dbc.db.Model(&cm.LogInRequest{}).Where("user_email = ?", userEmail).Count(&n)
	if tx.Error != nil {
		return false, tx.Error
	}
	if n > 0 {
		return true, nil
	}
	return false, nil
}
func (dbc *DbController) GetUserByEmailAbleToLogIn(user *cm.User) (err error) {
	tx := dbc.db.First(user, "email = ? AND can_log_in = ?", user.Email, true)
	if tx.Error != nil {
		return tx.Error
	}
	return nil
}
func (dbc *DbController) ExistsSessionWithUserId(user *cm.User) (exists bool, err error) {
	var n int64
	tx := dbc.db.Model(&cm.Session{}).Where("user_id = ?", user.Id).Count(&n)
	if tx.Error != nil {
		return false, tx.Error
	}
	if n > 0 {
		return true, nil
	}
	return false, nil
}
func (dbc *DbController) CreateLogInRequest(lir cm.LogInRequest) (err error) {
	tx := dbc.db.Create(&lir)
	if tx.Error != nil {
		return tx.Error
	}
	return nil
}
func (dbc *DbController) GetFirstOutdatedLogInRequest(edgeTime time.Time) (lirs []cm.LogInRequest, err error) {
	tx := dbc.db.Limit(1).Where("created_at <= ?", edgeTime).Find(&lirs)
	if tx.Error != nil {
		return nil, tx.Error
	}
	return lirs, nil
}
func (dbc *DbController) DeleteOldLogInRequest(lir *cm.LogInRequest) (err error) {
	tx := dbc.db.Delete(&lir)
	if tx.Error != nil {
		return tx.Error
	}
	return nil
}
func (dbc *DbController) FindLogInRequest(lir *cm.LogInRequest) (err error) {
	tx := dbc.db.First(lir, "request_id = ?", lir.RequestId)
	if tx.Error != nil {
		return tx.Error
	}
	return nil
}
func (dbc *DbController) GetUserByIdAbleToLogIn(user *cm.User) (err error) {
	tx := dbc.db.Preload("Password").First(user, "id = ? AND can_log_in = ?", user.Id, true)
	if tx.Error != nil {
		return tx.Error
	}
	return nil
}
func (dbc *DbController) CreateSession(session *cm.Session) (err error) {
	tx := dbc.db.Create(session)
	if tx.Error != nil {
		return tx.Error
	}
	return nil
}
func (dbc *DbController) DeleteLogInRequest(lir *cm.LogInRequest) (err error) {
	tx := dbc.db.Where("id = ?", lir.Id).Delete(&lir)
	if tx.Error != nil {
		return tx.Error
	}
	return nil
}
func (dbc *DbController) CreateLogEvent(le *cm.LogEvent) (err error) {
	tx := dbc.db.Create(le)
	if tx.Error != nil {
		return tx.Error
	}
	return nil
}

// User authorisation.

func (dbc *DbController) GetUserWithSessionByIdAbleToLogIn(user *cm.User) (err error) {
	tx := dbc.db.Preload("Session").First(user, "id = ? AND can_log_in = ?", user.Id, true)
	if tx.Error != nil {
		return tx.Error
	}
	return nil
}
func (dbc *DbController) CreateLogOutRequest(lor cm.LogOutRequest) (err error) {
	tx := dbc.db.Create(&lor)
	if tx.Error != nil {
		return tx.Error
	}
	return nil
}
func (dbc *DbController) GetFirstOutdatedLogOutRequest(edgeTime time.Time) (lors []cm.LogOutRequest, err error) {
	tx := dbc.db.Limit(1).Where("created_at <= ?", edgeTime).Find(&lors)
	if tx.Error != nil {
		return nil, tx.Error
	}
	return lors, nil
}
func (dbc *DbController) DeleteOldLogOutRequest(lor *cm.LogOutRequest) (err error) {
	tx := dbc.db.Delete(&lor)
	if tx.Error != nil {
		return tx.Error
	}
	return nil
}
func (dbc *DbController) FindLogOutRequest(lor *cm.LogOutRequest) (err error) {
	tx := dbc.db.First(lor, "request_id = ?", lor.RequestId)
	if tx.Error != nil {
		return tx.Error
	}
	return nil
}
func (dbc *DbController) DeleteSession(session *cm.Session) (err error) {
	tx := dbc.db.Delete(session)
	if tx.Error != nil {
		return tx.Error
	}
	return nil
}
func (dbc *DbController) DeleteLogOutRequest(lor *cm.LogOutRequest) (err error) {
	tx := dbc.db.Where("id = ?", lor.Id).Delete(&lor)
	if tx.Error != nil {
		return tx.Error
	}
	return nil
}
func (dbc *DbController) GetFirstOutdatedSession(edgeTime time.Time) (ss []cm.Session, err error) {
	tx := dbc.db.Limit(1).Where("created_at <= ?", edgeTime).Find(&ss)
	if tx.Error != nil {
		return nil, tx.Error
	}
	return ss, nil
}
func (dbc *DbController) DeleteOldSession(s *cm.Session) (err error) {
	tx := dbc.db.Delete(&s)
	if tx.Error != nil {
		return tx.Error
	}
	return nil
}

// E-mail change.

func (dbc *DbController) CreateEmailChangeRequest(ecr cm.EmailChangeRequest) (err error) {
	tx := dbc.db.Create(&ecr)
	if tx.Error != nil {
		return tx.Error
	}
	return nil
}
func (dbc *DbController) ExistsEmailChangeRequestWithNewEmail(newEmail string) (exists bool, err error) {
	var n int64
	tx := dbc.db.Model(&cm.EmailChangeRequest{}).Where("new_email = ?", newEmail).Count(&n)
	if tx.Error != nil {
		return false, tx.Error
	}
	if n > 0 {
		return true, nil
	}
	return false, nil
}
func (dbc *DbController) FindEmailChangeRequest(ecr *cm.EmailChangeRequest) (err error) {
	tx := dbc.db.First(ecr, "request_id = ?", ecr.RequestId)
	if tx.Error != nil {
		return tx.Error
	}
	return nil
}
func (dbc *DbController) SaveUserEmail(user *cm.User, email string) (err error) {
	user.Email = email

	tx := dbc.db.Save(user)
	if tx.Error != nil {
		return tx.Error
	}
	return nil
}
func (dbc *DbController) DeleteEmailChangeRequest(ecr *cm.EmailChangeRequest) (err error) {
	tx := dbc.db.Where("id = ?", ecr.Id).Delete(&ecr)
	if tx.Error != nil {
		return tx.Error
	}
	return nil
}
func (dbc *DbController) GetFirstOutdatedEmailChangeRequest(edgeTime time.Time) (ecrs []cm.EmailChangeRequest, err error) {
	tx := dbc.db.Limit(1).Where("created_at <= ?", edgeTime).Find(&ecrs)
	if tx.Error != nil {
		return nil, tx.Error
	}
	return ecrs, nil
}
func (dbc *DbController) DeleteOldEmailChangeRequest(ecr *cm.EmailChangeRequest) (err error) {
	tx := dbc.db.Delete(&ecr)
	if tx.Error != nil {
		return tx.Error
	}
	return nil
}

// Password change.

func (dbc *DbController) ExistsPasswordChangeRequestWithUserId(user *cm.User) (exists bool, err error) {
	var n int64
	tx := dbc.db.Model(&cm.PasswordChangeRequest{}).Where("user_id = ?", user.Id).Count(&n)
	if tx.Error != nil {
		return false, tx.Error
	}
	if n > 0 {
		return true, nil
	}
	return false, nil
}
func (dbc *DbController) CreatePasswordChangeRequest(pcr cm.PasswordChangeRequest) (err error) {
	tx := dbc.db.Create(&pcr)
	if tx.Error != nil {
		return tx.Error
	}
	return nil
}
func (dbc *DbController) GetFirstOutdatedPasswordChangeRequest(edgeTime time.Time) (pcrs []cm.PasswordChangeRequest, err error) {
	tx := dbc.db.Limit(1).Where("created_at <= ?", edgeTime).Find(&pcrs)
	if tx.Error != nil {
		return nil, tx.Error
	}
	return pcrs, nil
}
func (dbc *DbController) DeleteOldPasswordChangeRequest(pcr *cm.PasswordChangeRequest) (err error) {
	tx := dbc.db.Delete(&pcr)
	if tx.Error != nil {
		return tx.Error
	}
	return nil
}
func (dbc *DbController) FindPasswordChangeRequest(pcr *cm.PasswordChangeRequest) (err error) {
	tx := dbc.db.First(pcr, "request_id = ?", pcr.RequestId)
	if tx.Error != nil {
		return tx.Error
	}
	return nil
}
func (dbc *DbController) SaveUserPassword(user *cm.User, password string) (err error) {
	user.Password.Bytes, err = bpp.PackSymbols([]rune(password))
	if err != nil {
		return err
	}

	//tx := dbc.db.Save(user) <- This does not work in modern versions of GORM !
	tx := dbc.db.Session(&gorm.Session{FullSaveAssociations: true}).Save(user)
	if tx.Error != nil {
		return tx.Error
	}
	return nil
}
func (dbc *DbController) DeletePasswordChangeRequest(pcr *cm.PasswordChangeRequest) (err error) {
	tx := dbc.db.Where("id = ?", pcr.Id).Delete(&pcr)
	if tx.Error != nil {
		return tx.Error
	}
	return nil
}

// User & Session.

func (dbc *DbController) FindUserSession(user *cm.User) (session *cm.Session, err error) {
	session = new(cm.Session)
	tx := dbc.db.First(session, "user_id = ?", user.Id)
	if tx.Error != nil {
		return nil, tx.Error
	}
	return session, nil
}
func (dbc *DbController) ListUsers(page int) (users []cm.User, totalCount int, err error) {
	model := dbc.db.Model(&cm.User{})

	totalCount, err = dbc.listItemsOnPageWithTotalCount(model, page, &users)
	if err != nil {
		return nil, totalCount, err
	}

	return users, totalCount, nil
}
func (dbc *DbController) ListUserSessions(page int) (sessions []cm.Session, totalCount int, err error) {
	model := dbc.db.Preload("User").Model(&cm.Session{})

	totalCount, err = dbc.listItemsOnPageWithTotalCount(model, page, &sessions)
	if err != nil {
		return nil, totalCount, err
	}

	return sessions, totalCount, nil
}
func (dbc *DbController) GetUserName(user *cm.User) (err error) {
	tmpUser := new(cm.User)
	tx := dbc.db.First(tmpUser, "id = ?", user.Id)
	if tx.Error != nil {
		return tx.Error
	}
	user.Name = tmpUser.Name
	return nil
}
func (dbc *DbController) GetUserRoles(user *cm.User) (err error) {
	tmpUser := new(cm.User)
	tx := dbc.db.First(tmpUser, "id = ?", user.Id)
	if tx.Error != nil {
		return tx.Error
	}
	user.Roles = tmpUser.Roles
	return nil
}
func (dbc *DbController) GetUserParameters(user *cm.User) (err error) {
	tx := dbc.db.First(user, "id = ?", user.Id)
	if tx.Error != nil {
		return tx.Error
	}
	return nil
}
func (dbc *DbController) SetUserRole(user *cm.User, roleName string, newValue bool) (err error) {
	var tx *gorm.DB

	switch roleName {
	case UserRoleName_Author:
		tx = dbc.db.Model(user).Where("id = ?", user.Id).Update("can_create_thread", newValue)

	case UserRoleName_Writer:
		tx = dbc.db.Model(user).Where("id = ?", user.Id).Update("can_write_message", newValue)

	case UserRoleName_Reader:
		tx = dbc.db.Model(user).Where("id = ?", user.Id).Update("can_read", newValue)

	case UserRoleName_User:
		tx = dbc.db.Model(user).Where("id = ?", user.Id).Update("can_log_in", newValue)

	default:
		return errors.New(Err_InvalidUserRoleName)
	}

	if tx.Error != nil {
		return tx.Error
	}

	return nil
}
func (dbc *DbController) SetUserBanTime(user *cm.User) (err error) {
	var tx *gorm.DB
	tx = dbc.db.Model(user).Where("id = ?", user.Id).Update("ban_time", user.BanTime)
	if tx.Error != nil {
		return tx.Error
	}

	return nil
}

// Forum.

func (dbc *DbController) getForumMaxPos() (pos int, err error) {
	var forum cm.Forum
	tx := dbc.db.Select("pos").Order("pos DESC").First(&forum)
	if tx.Error != nil {
		return -1, tx.Error
	}
	return forum.Pos, nil
}
func (dbc *DbController) getForumPreviousPos(curPos int) (pos int, err error) {
	var forum cm.Forum
	tx := dbc.db.Select("pos").Where("pos < ?", curPos).Order("pos DESC").First(&forum)
	if tx.Error != nil {
		return -1, tx.Error
	}
	return forum.Pos, nil
}
func (dbc *DbController) getForumNextPos(curPos int) (pos int, err error) {
	var forum cm.Forum
	tx := dbc.db.Select("pos").Where("pos > ?", curPos).Order("pos ASC").First(&forum)
	if tx.Error != nil {
		return -1, tx.Error
	}
	return forum.Pos, nil
}
func (dbc *DbController) addForum(forum *cm.Forum) (err error) {
	tx := dbc.db.Create(forum)
	if tx.Error != nil {
		return tx.Error
	}
	return nil
}
func (dbc *DbController) getForumById(forumIn *cm.Forum) (forumOut *cm.Forum, err error) {
	var forum cm.Forum
	tx := dbc.db.First(&forum, forumIn.Id)
	if tx.Error != nil {
		return nil, tx.Error
	}
	return &forum, nil
}
func (dbc *DbController) updateForumPos(oldPos, newPos int) (err error) {
	tx := dbc.db.Model(&cm.Forum{}).Where("pos = ?", oldPos).Update("pos", newPos)
	if tx.Error != nil {
		return tx.Error
	}
	if tx.RowsAffected != 1 {
		return errors.New(Err_NoRowsUpdated)
	}
	return nil
}
func (dbc *DbController) AddForum(forumIn *cm.Forum) (forumOut *cm.Forum, err error) {
	var forumsCount int
	var forumMaxPos = 0
	var forum cm.Forum

	err = dbc.db.Transaction(func(tx *gorm.DB) error {
		var txErr error
		forumsCount, txErr = dbc.countAllItems(dbc.db.Model(&cm.Forum{}))
		if txErr != nil {
			return txErr
		}

		if forumsCount > 0 {
			forumMaxPos, txErr = dbc.getForumMaxPos()
			if txErr != nil {
				return txErr
			}
		}

		forum = cm.Forum{
			Name: forumIn.Name,
			Pos:  forumMaxPos + 1,
		}

		txErr = dbc.addForum(&forum)
		if txErr != nil {
			return txErr
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return &forum, nil
}
func (dbc *DbController) GetForum(forumIn *cm.Forum) (forumOut *cm.Forum, err error) {
	return dbc.getForumById(forumIn)
}
func (dbc *DbController) ListAllForums() (forums []cm.Forum, err error) {
	tx := dbc.db.Order("pos asc").Find(&forums)
	if tx.Error != nil {
		return nil, tx.Error
	}
	return forums, nil
}
func (dbc *DbController) ChangeForumName(forum *cm.Forum) (err error) {
	tx := dbc.db.Model(&cm.Forum{}).Where("id = ?", forum.Id).Update("name", forum.Name)
	if tx.Error != nil {
		return tx.Error
	}
	if tx.RowsAffected != 1 {
		return errors.New(Err_NoRowsUpdated)
	}
	return nil
}
func (dbc *DbController) MoveForumUp(forum *cm.Forum) (err error) {
	forum, err = dbc.getForumById(forum)
	if err != nil {
		return err
	}
	var curPos = forum.Pos

	var prevPos int
	prevPos, err = dbc.getForumPreviousPos(curPos)
	if err != nil {
		return err
	}

	// Swap values. Since SQL standard does not support such functionality, we
	// "invent a wheel" here, i.e. use a temporary value in order not to break
	// the 'unique' restriction.
	err = dbc.db.Transaction(func(tx *gorm.DB) error {
		var txErr error
		txErr = dbc.updateForumPos(prevPos, TmpForumPos)
		if txErr != nil {
			return txErr
		}

		txErr = dbc.updateForumPos(curPos, prevPos)
		if txErr != nil {
			return txErr
		}

		txErr = dbc.updateForumPos(TmpForumPos, curPos)
		if txErr != nil {
			return txErr
		}

		return nil
	})
	if err != nil {
		return err
	}

	return nil
}
func (dbc *DbController) MoveForumDown(forum *cm.Forum) (err error) {
	forum, err = dbc.getForumById(forum)
	if err != nil {
		return err
	}
	var curPos = forum.Pos

	var nextPos int
	nextPos, err = dbc.getForumNextPos(curPos)
	if err != nil {
		return err
	}

	// Swap values. Since SQL standard does not support such functionality, we
	// "invent a wheel" here, i.e. use a temporary value in order not to break
	// the 'unique' restriction.
	err = dbc.db.Transaction(func(tx *gorm.DB) error {
		var txErr error
		txErr = dbc.updateForumPos(nextPos, TmpForumPos)
		if txErr != nil {
			return txErr
		}

		txErr = dbc.updateForumPos(curPos, nextPos)
		if txErr != nil {
			return txErr
		}

		txErr = dbc.updateForumPos(TmpForumPos, curPos)
		if txErr != nil {
			return txErr
		}

		return nil
	})
	if err != nil {
		return err
	}

	return nil
}
func (dbc *DbController) DeleteForum(forum *cm.Forum) (err error) {
	tx := dbc.db.Delete(forum)
	if tx.Error != nil {
		return tx.Error
	}
	if tx.RowsAffected != 1 {
		return errors.New(Err_NoRowsUpdated)
	}
	return nil
}

// Thread.

func (dbc *DbController) getThreadById(threadIn *cm.Thread) (threadOut *cm.Thread, err error) {
	var thread cm.Thread
	tx := dbc.db.First(&thread, threadIn.Id)
	if tx.Error != nil {
		return nil, tx.Error
	}
	return &thread, nil
}
func (dbc *DbController) AddThread(forum *cm.Forum, threadIn *cm.Thread) (threadOut *cm.Thread, err error) {
	thread := cm.Thread{
		Name:    threadIn.Name,
		ForumId: forum.Id,
	}

	tx := dbc.db.Create(&thread)
	if tx.Error != nil {
		return nil, tx.Error
	}

	return &thread, nil
}
func (dbc *DbController) GetThread(threadIn *cm.Thread) (threadOut *cm.Thread, err error) {
	return dbc.getThreadById(threadIn)
}
func (dbc *DbController) ListThreads(forum *cm.Forum, page int) (threads []cm.Thread, totalCount int, err error) {
	model := dbc.db.Model(&cm.Thread{}).Where("forum_id = ?", forum.Id).Order("updated_at DESC")

	totalCount, err = dbc.listItemsOnPageWithTotalCount(model, page, &threads)
	if err != nil {
		return nil, totalCount, err
	}

	return threads, totalCount, nil
}
func (dbc *DbController) ChangeThreadName(thread *cm.Thread) (err error) {
	tx := dbc.db.Model(&cm.Thread{}).Where("id = ?", thread.Id).Update("name", thread.Name)
	if tx.Error != nil {
		return tx.Error
	}
	if tx.RowsAffected != 1 {
		return errors.New(Err_NoRowsUpdated)
	}
	return nil
}
func (dbc *DbController) ChangeThreadForum(thread *cm.Thread) (err error) {
	tx := dbc.db.Model(&cm.Thread{}).Where("id = ?", thread.Id).Update("forum_id", thread.ForumId)
	if tx.Error != nil {
		return tx.Error
	}
	if tx.RowsAffected != 1 {
		return errors.New(Err_NoRowsUpdated)
	}
	return nil
}
func (dbc *DbController) DeleteThread(thread *cm.Thread) (err error) {
	tx := dbc.db.Delete(thread)
	if tx.Error != nil {
		return tx.Error
	}
	if tx.RowsAffected != 1 {
		return errors.New(Err_NoRowsUpdated)
	}
	return nil
}
func (dbc *DbController) TouchThread(thread *cm.Thread) (err error) {
	tx := dbc.db.Model(&cm.Thread{}).Where("id = ?", thread.Id).Update("updated_at", time.Now())
	if tx.Error != nil {
		return tx.Error
	}
	if tx.RowsAffected != 1 {
		return errors.New(Err_NoRowsUpdated)
	}
	return nil
}

// Message.

func (dbc *DbController) CountThreadMessages(thread *cm.Thread) (n int, err error) {
	var n64 int64
	tx := dbc.db.Model(&cm.Message{}).Where("thread_id = ?", thread.Id).Count(&n64)
	if tx.Error != nil {
		return CountOnError, tx.Error
	}
	return int(n64), nil
}
func (dbc *DbController) GetThreadLastMessage(thread *cm.Thread) (messageOut *cm.Message, err error) {
	var message cm.Message
	tx := dbc.db.Where("thread_id = ?", thread.Id).Order("created_at DESC").First(&message)
	if tx.Error != nil {
		return nil, tx.Error
	}
	return &message, nil
}
func (dbc *DbController) AddMessage(user *cm.User, thread *cm.Thread, messageIn *cm.Message) (messageOut *cm.Message, err error) {
	message := cm.Message{
		Text:      messageIn.Text,
		ThreadId:  thread.Id,
		CreatorId: user.Id,
	}

	tx := dbc.db.Create(&message)
	if tx.Error != nil {
		return nil, tx.Error
	}

	return &message, nil
}
func (dbc *DbController) GetMessage(messageIn *cm.Message) (messageOut *cm.Message, err error) {
	var message cm.Message
	tx := dbc.db.First(&message, messageIn.Id)
	if tx.Error != nil {
		return nil, tx.Error
	}
	return &message, nil
}
func (dbc *DbController) ListMessages(thread *cm.Thread, page int) (messages []cm.Message, totalCount int, err error) {
	model := dbc.db.Model(&cm.Message{}).Where("thread_id = ?", thread.Id).Order("created_at ASC")

	totalCount, err = dbc.listItemsOnPageWithTotalCount(model, page, &messages)
	if err != nil {
		return nil, totalCount, err
	}

	return messages, totalCount, nil
}
func (dbc *DbController) ChangeMessageText(user *cm.User, messageIn *cm.Message) (err error) {
	tx := dbc.db.Model(&cm.Message{}).Where("id = ?", messageIn.Id).Updates(
		map[string]interface{}{
			"text":      messageIn.Text,
			"editor_id": user.Id,
		},
	)
	if tx.Error != nil {
		return tx.Error
	}
	if tx.RowsAffected != 1 {
		return errors.New(Err_NoRowsUpdated)
	}
	return nil
}
func (dbc *DbController) ChangeMessageThread(user *cm.User, messageIn *cm.Message) (err error) {
	tx := dbc.db.Model(&cm.Message{}).Where("id = ?", messageIn.Id).Updates(
		map[string]interface{}{
			"thread_id": messageIn.ThreadId,
			"editor_id": user.Id,
		},
	)
	if tx.Error != nil {
		return tx.Error
	}
	if tx.RowsAffected != 1 {
		return errors.New(Err_NoRowsUpdated)
	}
	return nil
}
func (dbc *DbController) DeleteMessage(message *cm.Message) (err error) {
	tx := dbc.db.Delete(message)
	if tx.Error != nil {
		return tx.Error
	}
	if tx.RowsAffected != 1 {
		return errors.New(Err_NoRowsUpdated)
	}
	return nil
}
