package c

import (
	cc "github.com/vault-thirteen/SimpleBB/pkg/common/models/Client"
)

// List of supported functions.
const (
	// Ping.
	FuncPing = cc.FuncPing

	// Section.
	FuncAddSection          = "AddSection"
	FuncChangeSectionName   = "ChangeSectionName"
	FuncChangeSectionParent = "ChangeSectionParent"
	FuncGetSection          = "GetSection"
	FuncMoveSectionUp       = "MoveSectionUp"
	FuncMoveSectionDown     = "MoveSectionDown"
	FuncDeleteSection       = "DeleteSection"

	// Forum.
	FuncAddForum           = "AddForum"
	FuncChangeForumName    = "ChangeForumName"
	FuncChangeForumSection = "ChangeForumSection"
	FuncGetForum           = "GetForum"
	FuncMoveForumUp        = "MoveForumUp"
	FuncMoveForumDown      = "MoveForumDown"
	FuncDeleteForum        = "DeleteForum"

	// Thread.
	FuncAddThread           = "AddThread"
	FuncChangeThreadName    = "ChangeThreadName"
	FuncChangeThreadForum   = "ChangeThreadForum"
	FuncGetThread           = "GetThread"
	FuncGetThreadNamesByIds = "GetThreadNamesByIds"
	FuncMoveThreadUp        = "MoveThreadUp"
	FuncMoveThreadDown      = "MoveThreadDown"
	FuncDeleteThread        = "DeleteThread"
	FuncThreadExistsS       = "ThreadExistsS"

	// Message.
	FuncAddMessage               = "AddMessage"
	FuncChangeMessageText        = "ChangeMessageText"
	FuncChangeMessageThread      = "ChangeMessageThread"
	FuncGetMessage               = "GetMessage"
	FuncGetLatestMessageOfThread = "GetLatestMessageOfThread"
	FuncDeleteMessage            = "DeleteMessage"

	// Composite objects.
	FuncListThreadAndMessages       = "ListThreadAndMessages"
	FuncListThreadAndMessagesOnPage = "ListThreadAndMessagesOnPage"
	FuncListForumAndThreads         = "ListForumAndThreads"
	FuncListForumAndThreadsOnPage   = "ListForumAndThreadsOnPage"
	FuncListSectionsAndForums       = "ListSectionsAndForums"

	// Other.
	FuncGetDKey            = "GetDKey"
	FuncShowDiagnosticData = cc.FuncShowDiagnosticData
	FuncTest               = "Test"
)
