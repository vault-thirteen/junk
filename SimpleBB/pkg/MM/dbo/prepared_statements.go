package dbo

import (
	"database/sql"
	"fmt"
)

// Indices of prepared statements.
const (
	DbPsid_ReadSections                 = 0
	DbPsid_InsertNewForum               = 1
	DbPsid_CountForumsById              = 2
	DbPsid_DeleteSectionById            = 3
	DbPsid_GetSectionById               = 4
	DbPsid_SetForumNameById             = 5
	DbPsid_SetSectionChildTypeById      = 6
	DbPsid_SetForumSectionById          = 7
	DbPsid_GetForumSectionById          = 8
	DbPsid_InsertNewThread              = 9
	DbPsid_GetForumThreadsById          = 10
	DbPsid_SetForumThreadsById          = 11
	DbPsid_SetThreadNameById            = 12
	DbPsid_GetThreadForumById           = 13
	DbPsid_SetThreadForumById           = 14
	DbPsid_CountThreadsById             = 15
	DbPsid_GetThreadMessagesById        = 16
	DbPsid_InsertNewMessage             = 17
	DbPsid_SetThreadMessagesById        = 18
	DbPsid_SetMessageTextById           = 19
	DbPsid_GetMessageThreadById         = 20
	DbPsid_SetMessageThreadById         = 21
	DbPsid_GetMessageCreatorAndTimeById = 22
	DbPsid_GetMessageById               = 23
	DbPsid_DeleteMessageById            = 24
	DbPsid_GetThreadByIdM               = 25
	DbPsid_DeleteThreadById             = 26
	DbPsid_GetForumById                 = 27
	DbPsid_DeleteForumById              = 28
	DbPsid_ReadForums                   = 29
	DbPsid_CountRootSections            = 30
	DbPsid_InsertNewSection             = 31
	DbPsid_CountSectionsById            = 32
	DbPsid_GetSectionChildrenById       = 33
	DbPsid_SetSectionChildrenById       = 34
	DbPsid_SetSectionNameById           = 35
	DbPsid_GetSectionParentById         = 36
	DbPsid_SetSectionParentById         = 37
	DbPsid_GetSectionChildTypeById      = 38
	DbPsid_CountMessagesById            = 39
	DbPsid_ReadThreadLinks              = 40
)

func (dbo *DatabaseObject) makePreparedStatementQueryStrings() (qs []string) {
	var q string
	qs = make([]string, 0)

	// 0.
	q = fmt.Sprintf(`SELECT Id, Parent, ChildType, Children, Name, CreatorUserId, CreatorTime, EditorUserId, EditorTime FROM %s;`, dbo.tableNames.Sections)
	qs = append(qs, q)

	// 1.
	q = fmt.Sprintf(`INSERT INTO %s (SectionId, NAME, CreatorUserId, CreatorTime) VALUES (?, ?, ?, Now());`, dbo.tableNames.Forums)
	qs = append(qs, q)

	// 2.
	q = fmt.Sprintf(`SELECT COUNT(Id) FROM %s WHERE Id = ?;`, dbo.tableNames.Forums)
	qs = append(qs, q)

	// 3.
	q = fmt.Sprintf(`DELETE FROM %s WHERE Id = ? AND ((Children IS NULL) OR (JSON_LENGTH(JSON_EXTRACT(Children, "$")) = 0));`, dbo.tableNames.Sections)
	qs = append(qs, q)

	// 4.
	q = fmt.Sprintf(`SELECT Id, Parent, ChildType, Children, Name, CreatorUserId, CreatorTime, EditorUserId, EditorTime FROM %s WHERE Id = ?;`, dbo.tableNames.Sections)
	qs = append(qs, q)

	// 5.
	q = fmt.Sprintf(`UPDATE %s SET NAME = ?, EditorUserId = ?, EditorTime = Now() WHERE Id = ?;`, dbo.tableNames.Forums)
	qs = append(qs, q)

	// 6.
	q = fmt.Sprintf(`UPDATE %s SET ChildType = ? WHERE Id = ?;`, dbo.tableNames.Sections)
	qs = append(qs, q)

	// 7.
	q = fmt.Sprintf(`UPDATE %s SET SectionId = ?, EditorUserId = ?, EditorTime = Now() WHERE Id = ?;`, dbo.tableNames.Forums)
	qs = append(qs, q)

	// 8.
	q = fmt.Sprintf(`SELECT SectionId FROM %s WHERE Id = ?;`, dbo.tableNames.Forums)
	qs = append(qs, q)

	// 9.
	q = fmt.Sprintf(`INSERT INTO %s (ForumId, NAME, CreatorUserId, CreatorTime) VALUES (?, ?, ?, Now());`, dbo.tableNames.Threads)
	qs = append(qs, q)

	// 10.
	q = fmt.Sprintf(`SELECT Threads FROM %s WHERE Id = ?;`, dbo.tableNames.Forums)
	qs = append(qs, q)

	// 11.
	q = fmt.Sprintf(`UPDATE %s SET Threads = ? WHERE Id = ?;`, dbo.tableNames.Forums)
	qs = append(qs, q)

	// 12.
	q = fmt.Sprintf(`UPDATE %s SET NAME = ?, EditorUserId = ?, EditorTime = Now() WHERE Id = ?;`, dbo.tableNames.Threads)
	qs = append(qs, q)

	// 13.
	q = fmt.Sprintf(`SELECT ForumId FROM %s WHERE Id = ?;`, dbo.tableNames.Threads)
	qs = append(qs, q)

	// 14.
	q = fmt.Sprintf(`UPDATE %s SET ForumId = ?, EditorUserId = ?, EditorTime = Now() WHERE Id = ?;`, dbo.tableNames.Threads)
	qs = append(qs, q)

	// 15.
	q = fmt.Sprintf(`SELECT COUNT(Id) FROM %s WHERE Id = ?;`, dbo.tableNames.Threads)
	qs = append(qs, q)

	// 16.
	q = fmt.Sprintf(`SELECT Messages FROM %s WHERE Id = ?;`, dbo.tableNames.Threads)
	qs = append(qs, q)

	// 17.
	q = fmt.Sprintf(`INSERT INTO %s (ThreadId, TEXT, TextChecksum, CreatorUserId, CreatorTime) VALUES (?, ?, ?, ?, Now());`, dbo.tableNames.Messages)
	qs = append(qs, q)

	// 18.
	q = fmt.Sprintf(`UPDATE %s SET Messages = ? WHERE Id = ?;`, dbo.tableNames.Threads)
	qs = append(qs, q)

	// 19.
	q = fmt.Sprintf(`UPDATE %s SET TEXT = ?, TextChecksum = ?, EditorUserId = ?, EditorTime = Now() WHERE Id = ?;`, dbo.tableNames.Messages)
	qs = append(qs, q)

	// 20.
	q = fmt.Sprintf(`SELECT ThreadId FROM %s WHERE Id = ?;`, dbo.tableNames.Messages)
	qs = append(qs, q)

	// 21.
	q = fmt.Sprintf(`UPDATE %s SET ThreadId = ?, EditorUserId = ?, EditorTime = Now() WHERE Id = ?;`, dbo.tableNames.Messages)
	qs = append(qs, q)

	// 22.
	q = fmt.Sprintf(`SELECT CreatorUserId, CreatorTime, EditorTime FROM %s WHERE Id = ?;`, dbo.tableNames.Messages)
	qs = append(qs, q)

	// 23.
	q = fmt.Sprintf(`SELECT Id, ThreadId, Text, TextChecksum, CreatorUserId, CreatorTime, EditorUserId, EditorTime FROM %s WHERE Id = ?;`, dbo.tableNames.Messages)
	qs = append(qs, q)

	// 24.
	q = fmt.Sprintf(`DELETE FROM %s WHERE Id = ?;`, dbo.tableNames.Messages)
	qs = append(qs, q)

	// 25.
	q = fmt.Sprintf(`SELECT Id, ForumId, Name, Messages, CreatorUserId, CreatorTime, EditorUserId, EditorTime FROM %s WHERE Id = ?;`, dbo.tableNames.Threads)
	qs = append(qs, q)

	// 26.
	q = fmt.Sprintf(`DELETE FROM %s WHERE Id = ? AND ((Messages IS NULL) OR (JSON_LENGTH(JSON_EXTRACT(Messages, "$")) = 0));`, dbo.tableNames.Threads)
	qs = append(qs, q)

	// 27.
	q = fmt.Sprintf(`SELECT Id, SectionId, Name, Threads, CreatorUserId, CreatorTime, EditorUserId, EditorTime FROM %s WHERE Id = ?;`, dbo.tableNames.Forums)
	qs = append(qs, q)

	// 28.
	q = fmt.Sprintf(`DELETE FROM %s WHERE Id = ? AND ((Threads IS NULL) OR (JSON_LENGTH(JSON_EXTRACT(Threads, "$")) = 0));`, dbo.tableNames.Forums)
	qs = append(qs, q)

	// 29.
	q = fmt.Sprintf(`SELECT Id, SectionId, Name, Threads, CreatorUserId, CreatorTime, EditorUserId, EditorTime FROM %s;`, dbo.tableNames.Forums)
	qs = append(qs, q)

	// 30.
	q = fmt.Sprintf(`SELECT COUNT(Id) FROM %s WHERE Parent IS NULL;`, dbo.tableNames.Sections)
	qs = append(qs, q)

	// 31.
	q = fmt.Sprintf(`INSERT INTO %s (Parent, NAME, CreatorUserId, CreatorTime) VALUES (?, ?, ?, Now());`, dbo.tableNames.Sections)
	qs = append(qs, q)

	// 32.
	q = fmt.Sprintf(`SELECT COUNT(Id) FROM %s WHERE Id = ?;`, dbo.tableNames.Sections)
	qs = append(qs, q)

	// 33.
	q = fmt.Sprintf(`SELECT Children FROM %s WHERE Id = ?;`, dbo.tableNames.Sections)
	qs = append(qs, q)

	// 34.
	q = fmt.Sprintf(`UPDATE %s SET Children = ? WHERE Id = ?;`, dbo.tableNames.Sections)
	qs = append(qs, q)

	// 35.
	q = fmt.Sprintf(`UPDATE %s SET NAME = ?, EditorUserId = ?, EditorTime = Now() WHERE Id = ?;`, dbo.tableNames.Sections)
	qs = append(qs, q)

	// 36.
	q = fmt.Sprintf(`SELECT Parent FROM %s WHERE Id = ?;`, dbo.tableNames.Sections)
	qs = append(qs, q)

	// 37.
	q = fmt.Sprintf(`UPDATE %s SET Parent = ?, EditorUserId = ?, EditorTime = Now() WHERE Id = ?;`, dbo.tableNames.Sections)
	qs = append(qs, q)

	// 38.
	q = fmt.Sprintf(`SELECT ChildType FROM %s WHERE Id = ?;`, dbo.tableNames.Sections)
	qs = append(qs, q)

	// 39.
	q = fmt.Sprintf(`SELECT COUNT(Id) FROM %s WHERE Id = ?;`, dbo.tableNames.Messages)
	qs = append(qs, q)

	// 40.
	q = fmt.Sprintf(`SELECT Id, ForumId, Messages FROM %s;`, dbo.tableNames.Threads)
	qs = append(qs, q)

	return qs
}

func (dbo *DatabaseObject) GetPreparedStatementByIndex(i int) (ps *sql.Stmt) {
	return dbo.DatabaseObject.PreparedStatement(i)
}
