package dbo

import (
	"github.com/vault-thirteen/SimpleBB/pkg/common/models/UidList"
)

func (dbo *DatabaseObject) dbQuery_ReadMessagesById(messageIds ul.UidList) (query string, err error) {
	var vs string
	vs, err = messageIds.ValuesString()
	if err != nil {
		return "", err
	}

	return `SELECT Id, ThreadId, Text, TextChecksum, CreatorUserId, CreatorTime, EditorUserId, EditorTime FROM ` + dbo.tableNames.Messages + ` WHERE Id IN (` + vs + `) ORDER BY FIND_IN_SET(Id, '` + vs + `');`, nil
}

func (dbo *DatabaseObject) dbQuery_ReadMessageLinksById(messageIds ul.UidList) (query string, err error) {
	var vs string
	vs, err = messageIds.ValuesString()
	if err != nil {
		return "", err
	}

	return `SELECT Id, ThreadId FROM ` + dbo.tableNames.Messages + ` WHERE Id IN (` + vs + `) ORDER BY FIND_IN_SET(Id, '` + vs + `');`, nil
}

func (dbo *DatabaseObject) dbQuery_ReadThreadNamesByIds(threadIds ul.UidList) (query string, err error) {
	var vs string
	vs, err = threadIds.ValuesString()
	if err != nil {
		return "", err
	}

	return `SELECT Name FROM ` + dbo.tableNames.Threads + ` WHERE Id IN (` + vs + `) ORDER BY FIND_IN_SET(Id, '` + vs + `');`, nil
}

func (dbo *DatabaseObject) dbQuery_ReadThreadsById(threadIds ul.UidList) (query string, err error) {
	var vs string
	vs, err = threadIds.ValuesString()
	if err != nil {
		return "", err
	}

	return `SELECT Id, ForumId, Name, Messages, CreatorUserId, CreatorTime, EditorUserId, EditorTime FROM ` + dbo.tableNames.Threads + ` WHERE Id IN (` + vs + `) ORDER BY FIND_IN_SET(Id, '` + vs + `');`, nil
}
