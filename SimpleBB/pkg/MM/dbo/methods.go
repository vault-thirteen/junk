package dbo

// Due to the large number of methods, they are sorted alphabetically.

import (
	"database/sql"
	"github.com/vault-thirteen/SimpleBB/pkg/common/interfaces/derived1"
	"github.com/vault-thirteen/SimpleBB/pkg/common/interfaces/derived2"
	"github.com/vault-thirteen/SimpleBB/pkg/common/models/UidList"
	base2 "github.com/vault-thirteen/SimpleBB/pkg/common/models/base"
	complex2 "github.com/vault-thirteen/SimpleBB/pkg/common/models/complex/Forum"
	"github.com/vault-thirteen/SimpleBB/pkg/common/models/complex/Message"
	"github.com/vault-thirteen/SimpleBB/pkg/common/models/complex/Section"
	"github.com/vault-thirteen/SimpleBB/pkg/common/models/complex/Thread"
	dbo2 "github.com/vault-thirteen/SimpleBB/pkg/common/models/dbo"
	cm "github.com/vault-thirteen/SimpleBB/pkg/common/models/simple"
	cms "github.com/vault-thirteen/SimpleBB/pkg/common/models/sql"
	"time"

	mm "github.com/vault-thirteen/SimpleBB/pkg/MM/models"
	ae "github.com/vault-thirteen/auxie/errors"
)

func (dbo *DatabaseObject) CountForumsById(forumId base2.Id) (n base2.Count, err error) {
	row := dbo.DatabaseObject.PreparedStatement(DbPsid_CountForumsById).QueryRow(forumId)

	n, err = cms.NewNonNullValueFromScannableSource[base2.Count](row)
	if err != nil {
		return dbo2.CountOnError, err
	}

	return n, nil
}

func (dbo *DatabaseObject) CountMessagesById(messageId base2.Id) (n base2.Count, err error) {
	row := dbo.DatabaseObject.PreparedStatement(DbPsid_CountMessagesById).QueryRow(messageId)

	n, err = cms.NewNonNullValueFromScannableSource[base2.Count](row)
	if err != nil {
		return dbo2.CountOnError, err
	}

	return n, nil
}

func (dbo *DatabaseObject) CountRootSections() (n base2.Count, err error) {
	row := dbo.DatabaseObject.PreparedStatement(DbPsid_CountRootSections).QueryRow()

	n, err = cms.NewNonNullValueFromScannableSource[base2.Count](row)
	if err != nil {
		return dbo2.CountOnError, err
	}

	return n, nil
}

func (dbo *DatabaseObject) CountSectionsById(sectionId base2.Id) (n base2.Count, err error) {
	row := dbo.DatabaseObject.PreparedStatement(DbPsid_CountSectionsById).QueryRow(sectionId)

	n, err = cms.NewNonNullValueFromScannableSource[base2.Count](row)
	if err != nil {
		return dbo2.CountOnError, err
	}

	return n, nil
}

func (dbo *DatabaseObject) CountThreadsById(threadId base2.Id) (n base2.Count, err error) {
	row := dbo.DatabaseObject.PreparedStatement(DbPsid_CountThreadsById).QueryRow(threadId)

	n, err = cms.NewNonNullValueFromScannableSource[base2.Count](row)
	if err != nil {
		return dbo2.CountOnError, err
	}

	return n, nil
}

func (dbo *DatabaseObject) DeleteForumById(forumId base2.Id) (err error) {
	var result sql.Result
	result, err = dbo.DatabaseObject.PreparedStatement(DbPsid_DeleteForumById).Exec(forumId)
	if err != nil {
		return err
	}

	return dbo2.CheckRowsAffected(result, 1)
}

func (dbo *DatabaseObject) DeleteMessageById(messageId base2.Id) (err error) {
	var result sql.Result
	result, err = dbo.DatabaseObject.PreparedStatement(DbPsid_DeleteMessageById).Exec(messageId)
	if err != nil {
		return err
	}

	return dbo2.CheckRowsAffected(result, 1)
}

func (dbo *DatabaseObject) DeleteSectionById(sectionId base2.Id) (err error) {
	var result sql.Result
	result, err = dbo.DatabaseObject.PreparedStatement(DbPsid_DeleteSectionById).Exec(sectionId)
	if err != nil {
		return err
	}

	return dbo2.CheckRowsAffected(result, 1)
}

func (dbo *DatabaseObject) DeleteThreadById(threadId base2.Id) (err error) {
	var result sql.Result
	result, err = dbo.DatabaseObject.PreparedStatement(DbPsid_DeleteThreadById).Exec(threadId)
	if err != nil {
		return err
	}

	return dbo2.CheckRowsAffected(result, 1)
}

func (dbo *DatabaseObject) GetForumById(forumId base2.Id) (forum derived2.IForum, err error) {
	row := dbo.DatabaseObject.PreparedStatement(DbPsid_GetForumById).QueryRow(forumId)

	forum, err = complex2.NewForumFromScannableSource(row)
	if err != nil {
		return nil, err
	}

	return forum, nil
}

func (dbo *DatabaseObject) GetForumSectionById(forumId base2.Id) (sectionId base2.Id, err error) {
	row := dbo.DatabaseObject.PreparedStatement(DbPsid_GetForumSectionById).QueryRow(forumId)

	sectionId, err = cms.NewNonNullValueFromScannableSource[base2.Id](row)
	if err != nil {
		return dbo2.IdOnError, err
	}

	return sectionId, nil
}

func (dbo *DatabaseObject) GetForumThreadsById(forumId base2.Id) (threads *ul.UidList, err error) {
	threads = ul.New()
	err = dbo.DatabaseObject.PreparedStatement(DbPsid_GetForumThreadsById).QueryRow(forumId).Scan(threads)
	if err != nil {
		return nil, err
	}

	return threads, nil
}

func (dbo *DatabaseObject) GetMessageById(messageId base2.Id) (message derived2.IMessage, err error) {
	row := dbo.DatabaseObject.PreparedStatement(DbPsid_GetMessageById).QueryRow(messageId)

	message, err = m.NewMessageFromScannableSource(row)
	if err != nil {
		return nil, err
	}

	return message, nil
}

func (dbo *DatabaseObject) GetMessageCreatorAndTimeById(messageId base2.Id) (creatorUserId base2.Id, ToC time.Time, ToE *time.Time, err error) {
	// N.B.: ToC can not be null, but ToE can be null !
	err = dbo.DatabaseObject.PreparedStatement(DbPsid_GetMessageCreatorAndTimeById).QueryRow(messageId).Scan(&creatorUserId, &ToC, &ToE)
	if err != nil {
		return dbo2.IdOnError, time.Time{}, nil, err
	}

	return creatorUserId, ToC, ToE, nil
}

func (dbo *DatabaseObject) GetMessageThreadById(messageId base2.Id) (threadId base2.Id, err error) {
	row := dbo.DatabaseObject.PreparedStatement(DbPsid_GetMessageThreadById).QueryRow(messageId)

	threadId, err = cms.NewNonNullValueFromScannableSource[base2.Id](row)
	if err != nil {
		return dbo2.IdOnError, err
	}

	return threadId, nil
}

func (dbo *DatabaseObject) GetSectionById(sectionId base2.Id) (section derived2.ISection, err error) {
	row := dbo.DatabaseObject.PreparedStatement(DbPsid_GetSectionById).QueryRow(sectionId)

	section, err = s.NewSectionFromScannableSource(row)
	if err != nil {
		return nil, err
	}

	return section, nil
}

func (dbo *DatabaseObject) GetSectionChildTypeById(sectionId base2.Id) (childType derived1.ISectionChildType, err error) {
	row := dbo.DatabaseObject.PreparedStatement(DbPsid_GetSectionChildTypeById).QueryRow(sectionId)

	childType, err = cms.NewNonNullValueFromScannableSource[derived1.ISectionChildType](row)
	if err != nil {
		return childType, err
	}

	return childType, nil
}

func (dbo *DatabaseObject) GetSectionChildrenById(sectionId base2.Id) (children *ul.UidList, err error) {
	children = ul.New()
	err = dbo.DatabaseObject.PreparedStatement(DbPsid_GetSectionChildrenById).QueryRow(sectionId).Scan(children)
	if err != nil {
		return nil, err
	}

	return children, nil
}

func (dbo *DatabaseObject) GetSectionParentById(sectionId base2.Id) (parent *base2.Id, err error) {
	row := dbo.DatabaseObject.PreparedStatement(DbPsid_GetSectionParentById).QueryRow(sectionId)
	return cms.NewValueFromScannableSource[base2.Id](row)
}

func (dbo *DatabaseObject) GetThreadById(threadId base2.Id) (thread derived2.IThread, err error) {
	row := dbo.DatabaseObject.PreparedStatement(DbPsid_GetThreadByIdM).QueryRow(threadId)

	thread, err = t.NewThreadFromScannableSource(row)
	if err != nil {
		return nil, err
	}

	return thread, nil
}

func (dbo *DatabaseObject) GetThreadForumById(threadId base2.Id) (forumId base2.Id, err error) {
	row := dbo.DatabaseObject.PreparedStatement(DbPsid_GetThreadForumById).QueryRow(threadId)

	forumId, err = cms.NewNonNullValueFromScannableSource[base2.Id](row)
	if err != nil {
		return dbo2.IdOnError, err
	}

	return forumId, nil
}

func (dbo *DatabaseObject) GetThreadMessagesById(threadId base2.Id) (messages *ul.UidList, err error) {
	messages = ul.New()
	err = dbo.DatabaseObject.PreparedStatement(DbPsid_GetThreadMessagesById).QueryRow(threadId).Scan(messages)
	if err != nil {
		return nil, err
	}

	return messages, nil
}

func (dbo *DatabaseObject) InsertNewForum(sectionId base2.Id, name cm.Name, creatorUserId base2.Id) (lastInsertedId base2.Id, err error) {
	var result sql.Result
	result, err = dbo.DatabaseObject.PreparedStatement(DbPsid_InsertNewForum).Exec(sectionId, name, creatorUserId)
	if err != nil {
		return dbo2.LastInsertedIdOnError, err
	}

	return dbo2.CheckRowsAffectedAndGetLastInsertedId(result, 1)
}

func (dbo *DatabaseObject) InsertNewMessage(parentThread base2.Id, messageText base2.Text, textChecksum []byte, creatorUserId base2.Id) (lastInsertedId base2.Id, err error) {
	var result sql.Result
	result, err = dbo.DatabaseObject.PreparedStatement(DbPsid_InsertNewMessage).Exec(parentThread, messageText, textChecksum, creatorUserId)
	if err != nil {
		return dbo2.LastInsertedIdOnError, err
	}

	return dbo2.CheckRowsAffectedAndGetLastInsertedId(result, 1)
}

func (dbo *DatabaseObject) InsertNewSection(parent *base2.Id, name cm.Name, creatorUserId base2.Id) (lastInsertedId base2.Id, err error) {
	var result sql.Result
	result, err = dbo.DatabaseObject.PreparedStatement(DbPsid_InsertNewSection).Exec(parent, name, creatorUserId)
	if err != nil {
		return dbo2.LastInsertedIdOnError, err
	}

	return dbo2.CheckRowsAffectedAndGetLastInsertedId(result, 1)
}

func (dbo *DatabaseObject) InsertNewThread(parentForum base2.Id, threadName cm.Name, creatorUserId base2.Id) (lastInsertedId base2.Id, err error) {
	var result sql.Result
	result, err = dbo.DatabaseObject.PreparedStatement(DbPsid_InsertNewThread).Exec(parentForum, threadName, creatorUserId)
	if err != nil {
		return dbo2.LastInsertedIdOnError, err
	}

	return dbo2.CheckRowsAffectedAndGetLastInsertedId(result, 1)
}

func (dbo *DatabaseObject) ReadForums() (forums []derived2.IForum, err error) {
	var rows *sql.Rows
	rows, err = dbo.DatabaseObject.PreparedStatement(DbPsid_ReadForums).Query()
	if err != nil {
		return nil, err
	}

	defer func() {
		derr := rows.Close()
		if derr != nil {
			err = ae.Combine(err, derr)
		}
	}()

	return complex2.NewForumArrayFromRows(rows)
}

func (dbo *DatabaseObject) ReadMessagesById(messageIds *ul.UidList) (messages []derived2.IMessage, err error) {
	if messageIds == nil {
		return []derived2.IMessage{}, nil
	}

	var query string
	query, err = dbo.dbQuery_ReadMessagesById(*messageIds)
	if err != nil {
		return nil, err
	}

	var rows *sql.Rows
	rows, err = dbo.DatabaseObject.DB().Query(query)
	if err != nil {
		return nil, err
	}

	defer func() {
		derr := rows.Close()
		if derr != nil {
			err = ae.Combine(err, derr)
		}
	}()

	return m.NewMessageArrayFromRows(rows)
}

func (dbo *DatabaseObject) ReadMessageLinksById(messageIds *ul.UidList) (messageLinks []mm.MessageLink, err error) {
	if messageIds == nil {
		return []mm.MessageLink{}, nil
	}

	var query string
	query, err = dbo.dbQuery_ReadMessageLinksById(*messageIds)
	if err != nil {
		return nil, err
	}

	var rows *sql.Rows
	rows, err = dbo.DatabaseObject.DB().Query(query)
	if err != nil {
		return nil, err
	}

	defer func() {
		derr := rows.Close()
		if derr != nil {
			err = ae.Combine(err, derr)
		}
	}()

	return mm.NewMessageLinkArrayFromRows(rows)
}

func (dbo *DatabaseObject) ReadSections() (sections []derived2.ISection, err error) {
	var rows *sql.Rows
	rows, err = dbo.DatabaseObject.PreparedStatement(DbPsid_ReadSections).Query()
	if err != nil {
		return nil, err
	}

	defer func() {
		derr := rows.Close()
		if derr != nil {
			err = ae.Combine(err, derr)
		}
	}()

	return s.NewSectionArrayFromRows(rows)
}

func (dbo *DatabaseObject) ReadThreadLinks() (threadLinks []mm.ThreadLink, err error) {
	var rows *sql.Rows
	rows, err = dbo.DatabaseObject.PreparedStatement(DbPsid_ReadThreadLinks).Query()
	if err != nil {
		return nil, err
	}

	defer func() {
		derr := rows.Close()
		if derr != nil {
			err = ae.Combine(err, derr)
		}
	}()

	return mm.NewThreadLinkArrayFromRows(rows)
}

func (dbo *DatabaseObject) ReadThreadNamesByIds(threadIds ul.UidList) (threadNames []cm.Name, err error) {
	var query string
	query, err = dbo.dbQuery_ReadThreadNamesByIds(threadIds)
	if err != nil {
		return nil, err
	}

	var rows *sql.Rows
	rows, err = dbo.DB().Query(query)
	if err != nil {
		return nil, err
	}

	defer func() {
		derr := rows.Close()
		if derr != nil {
			err = ae.Combine(err, derr)
		}
	}()

	return cms.NewArrayFromScannableSource[cm.Name](rows)
}

func (dbo *DatabaseObject) ReadThreadsById(threadIds *ul.UidList) (threads []derived2.IThread, err error) {
	if threadIds == nil {
		return []derived2.IThread{}, nil
	}

	var query string
	query, err = dbo.dbQuery_ReadThreadsById(*threadIds)
	if err != nil {
		return nil, err
	}

	var rows *sql.Rows
	rows, err = dbo.DB().Query(query)
	if err != nil {
		return nil, err
	}

	defer func() {
		derr := rows.Close()
		if derr != nil {
			err = ae.Combine(err, derr)
		}
	}()

	return t.NewThreadArrayFromRows(rows)
}

func (dbo *DatabaseObject) SetForumNameById(forumId base2.Id, name cm.Name, editorUserId base2.Id) (err error) {
	var result sql.Result
	result, err = dbo.DatabaseObject.PreparedStatement(DbPsid_SetForumNameById).Exec(name, editorUserId, forumId)
	if err != nil {
		return err
	}

	return dbo2.CheckRowsAffected(result, 1)
}

func (dbo *DatabaseObject) SetForumSectionById(forumId base2.Id, sectionId base2.Id, editorUserId base2.Id) (err error) {
	var result sql.Result
	result, err = dbo.DatabaseObject.PreparedStatement(DbPsid_SetForumSectionById).Exec(sectionId, editorUserId, forumId)
	if err != nil {
		return err
	}

	return dbo2.CheckRowsAffected(result, 1)
}

func (dbo *DatabaseObject) SetForumThreadsById(forumId base2.Id, threads *ul.UidList) (err error) {
	var result sql.Result
	result, err = dbo.DatabaseObject.PreparedStatement(DbPsid_SetForumThreadsById).Exec(threads, forumId)
	if err != nil {
		return err
	}

	return dbo2.CheckRowsAffected(result, 1)
}

func (dbo *DatabaseObject) SetMessageTextById(messageId base2.Id, text base2.Text, textChecksum []byte, editorUserId base2.Id) (err error) {
	var result sql.Result
	result, err = dbo.DatabaseObject.PreparedStatement(DbPsid_SetMessageTextById).Exec(text, textChecksum, editorUserId, messageId)
	if err != nil {
		return err
	}

	return dbo2.CheckRowsAffected(result, 1)
}

func (dbo *DatabaseObject) SetMessageThreadById(messageId base2.Id, threadId base2.Id, editorUserId base2.Id) (err error) {
	var result sql.Result
	result, err = dbo.DatabaseObject.PreparedStatement(DbPsid_SetMessageThreadById).Exec(threadId, editorUserId, messageId)
	if err != nil {
		return err
	}

	return dbo2.CheckRowsAffected(result, 1)
}

func (dbo *DatabaseObject) SetSectionChildTypeById(sectionId base2.Id, childType derived1.ISectionChildType) (err error) {
	var result sql.Result
	result, err = dbo.DatabaseObject.PreparedStatement(DbPsid_SetSectionChildTypeById).Exec(childType, sectionId)
	if err != nil {
		return err
	}

	return dbo2.CheckRowsAffected(result, 1)
}

func (dbo *DatabaseObject) SetSectionChildrenById(sectionId base2.Id, children *ul.UidList) (err error) {
	var result sql.Result
	result, err = dbo.DatabaseObject.PreparedStatement(DbPsid_SetSectionChildrenById).Exec(children, sectionId)
	if err != nil {
		return err
	}

	return dbo2.CheckRowsAffected(result, 1)
}

func (dbo *DatabaseObject) SetSectionNameById(sectionId base2.Id, name cm.Name, editorUserId base2.Id) (err error) {
	var result sql.Result
	result, err = dbo.DatabaseObject.PreparedStatement(DbPsid_SetSectionNameById).Exec(name, editorUserId, sectionId)
	if err != nil {
		return err
	}

	return dbo2.CheckRowsAffected(result, 1)
}

func (dbo *DatabaseObject) SetSectionParentById(sectionId base2.Id, parent base2.Id, editorUserId base2.Id) (err error) {
	var result sql.Result
	result, err = dbo.DatabaseObject.PreparedStatement(DbPsid_SetSectionParentById).Exec(parent, editorUserId, sectionId)
	if err != nil {
		return err
	}

	return dbo2.CheckRowsAffected(result, 1)
}

func (dbo *DatabaseObject) SetThreadForumById(threadId base2.Id, forumId base2.Id, editorUserId base2.Id) (err error) {
	var result sql.Result
	result, err = dbo.DatabaseObject.PreparedStatement(DbPsid_SetThreadForumById).Exec(forumId, editorUserId, threadId)
	if err != nil {
		return err
	}

	return dbo2.CheckRowsAffected(result, 1)
}

func (dbo *DatabaseObject) SetThreadMessagesById(threadId base2.Id, messages *ul.UidList) (err error) {
	var result sql.Result
	result, err = dbo.DatabaseObject.PreparedStatement(DbPsid_SetThreadMessagesById).Exec(messages, threadId)
	if err != nil {
		return err
	}

	return dbo2.CheckRowsAffected(result, 1)
}

func (dbo *DatabaseObject) SetThreadNameById(threadId base2.Id, name cm.Name, editorUserId base2.Id) (err error) {
	var result sql.Result
	result, err = dbo.DatabaseObject.PreparedStatement(DbPsid_SetThreadNameById).Exec(name, editorUserId, threadId)
	if err != nil {
		return err
	}

	return dbo2.CheckRowsAffected(result, 1)
}
