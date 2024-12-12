package rpc

import (
	"github.com/vault-thirteen/SimpleBB/pkg/MM/models"
	"github.com/vault-thirteen/SimpleBB/pkg/common/interfaces/derived2"
	base2 "github.com/vault-thirteen/SimpleBB/pkg/common/models/base"
	rpc2 "github.com/vault-thirteen/SimpleBB/pkg/common/models/rpc"
	cm "github.com/vault-thirteen/SimpleBB/pkg/common/models/simple"
)

// Ping.

type PingParams = rpc2.PingParams
type PingResult = rpc2.PingResult

// Section.

type AddSectionParams struct {
	rpc2.CommonParams

	// Identifier of a parent section containing this section.
	// Null means that this section is a root section.
	// Only a single root section can exist.
	Parent *base2.Id `json:"parent"`

	// Name of this section.
	Name cm.Name `json:"name"`
}
type AddSectionResult struct {
	rpc2.CommonResult

	// ID of the created section.
	SectionId base2.Id `json:"sectionId"`
}

type ChangeSectionNameParams struct {
	rpc2.CommonParams

	// Identifier of a section.
	SectionId base2.Id `json:"sectionId"`

	// Name of this section.
	Name cm.Name `json:"name"`
}
type ChangeSectionNameResult = rpc2.CommonResultWithSuccess

type ChangeSectionParentParams struct {
	rpc2.CommonParams

	// Identifier of a section.
	SectionId base2.Id `json:"sectionId"`

	// Identifier of a parent section containing this section.
	Parent base2.Id `json:"parent"`
}
type ChangeSectionParentResult = rpc2.CommonResultWithSuccess

type GetSectionParams struct {
	rpc2.CommonParams

	SectionId base2.Id `json:"sectionId"`
}
type GetSectionResult struct {
	rpc2.CommonResult

	Section derived2.ISection `json:"section"`
}

type MoveSectionUpParams struct {
	rpc2.CommonParams

	// Identifier of a section.
	SectionId base2.Id `json:"sectionId"`
}
type MoveSectionUpResult = rpc2.CommonResultWithSuccess

type MoveSectionDownParams struct {
	rpc2.CommonParams

	// Identifier of a section.
	SectionId base2.Id `json:"sectionId"`
}
type MoveSectionDownResult = rpc2.CommonResultWithSuccess

type DeleteSectionParams struct {
	rpc2.CommonParams

	SectionId base2.Id `json:"sectionId"`
}
type DeleteSectionResult = rpc2.CommonResultWithSuccess

// Forum.

type AddForumParams struct {
	rpc2.CommonParams

	// Identifier of a section containing this forum.
	SectionId base2.Id `json:"sectionId"`

	// Name of this forum.
	Name cm.Name `json:"name"`
}
type AddForumResult struct {
	rpc2.CommonResult

	// ID of the created forum.
	ForumId base2.Id `json:"forumId"`
}

type ChangeForumNameParams struct {
	rpc2.CommonParams

	ForumId base2.Id `json:"forumId"`

	// New name.
	Name cm.Name `json:"name"`
}
type ChangeForumNameResult = rpc2.CommonResultWithSuccess

type ChangeForumSectionParams struct {
	rpc2.CommonParams

	// Identifier of this forum.
	ForumId base2.Id `json:"forumId"`

	// Identifier of a section containing this forum.
	SectionId base2.Id `json:"sectionId"`
}
type ChangeForumSectionResult = rpc2.CommonResultWithSuccess

type GetForumParams struct {
	rpc2.CommonParams

	ForumId base2.Id `json:"forumId"`
}
type GetForumResult struct {
	rpc2.CommonResult

	Forum derived2.IForum `json:"forum"`
}

type MoveForumUpParams struct {
	rpc2.CommonParams

	// Identifier of a forum.
	ForumId base2.Id `json:"forumId"`
}
type MoveForumUpResult = rpc2.CommonResultWithSuccess

type MoveForumDownParams struct {
	rpc2.CommonParams

	// Identifier of a forum.
	ForumId base2.Id `json:"forumId"`
}
type MoveForumDownResult = rpc2.CommonResultWithSuccess

type DeleteForumParams struct {
	rpc2.CommonParams

	ForumId base2.Id `json:"forumId"`
}
type DeleteForumResult = rpc2.CommonResultWithSuccess

// Thread.

type AddThreadParams struct {
	rpc2.CommonParams

	// ID of a forum containing this thread.
	ForumId base2.Id `json:"forumId"`

	// Thread name.
	Name cm.Name `json:"name"`
}
type AddThreadResult struct {
	rpc2.CommonResult

	// ID of the created forum.
	ThreadId base2.Id `json:"threadId"`
}

type ChangeThreadNameParams struct {
	rpc2.CommonParams

	ThreadId base2.Id `json:"threadId"`

	// New name.
	Name cm.Name `json:"name"`
}
type ChangeThreadNameResult = rpc2.CommonResultWithSuccess

type ChangeThreadForumParams struct {
	rpc2.CommonParams

	ThreadId base2.Id `json:"threadId"`

	// ID of a new parent forum.
	ForumId base2.Id `json:"forumId"`
}
type ChangeThreadForumResult = rpc2.CommonResultWithSuccess

type GetThreadParams struct {
	rpc2.CommonParams

	ThreadId base2.Id `json:"threadId"`
}
type GetThreadResult struct {
	rpc2.CommonResult

	Thread derived2.IThread `json:"thread"`
}

type GetThreadNamesByIdsParams struct {
	rpc2.CommonParams

	ThreadIds []base2.Id `json:"threadIds"`
}
type GetThreadNamesByIdsResult struct {
	rpc2.CommonResult

	ThreadIds   []base2.Id `json:"threadIds"`
	ThreadNames []cm.Name  `json:"threadNames"`
}

type MoveThreadUpParams struct {
	rpc2.CommonParams

	// Identifier of a thread.
	ThreadId base2.Id `json:"threadId"`
}
type MoveThreadUpResult = rpc2.CommonResultWithSuccess

type MoveThreadDownParams struct {
	rpc2.CommonParams

	// Identifier of a thread.
	ThreadId base2.Id `json:"threadId"`
}
type MoveThreadDownResult = rpc2.CommonResultWithSuccess

type DeleteThreadParams struct {
	rpc2.CommonParams

	ThreadId base2.Id `json:"threadId"`
}
type DeleteThreadResult = rpc2.CommonResultWithSuccess

type ThreadExistsSParams struct {
	rpc2.CommonParams
	rpc2.DKeyParams

	ThreadId base2.Id `json:"threadId"`
}
type ThreadExistsSResult struct {
	rpc2.CommonResult

	Exists base2.Flag `json:"exists"`
}

// Message.

type AddMessageParams struct {
	rpc2.CommonParams

	// ID of a thread containing this message.
	ThreadId base2.Id `json:"threadId"`

	// Message text.
	Text base2.Text `json:"text"`
}
type AddMessageResult struct {
	rpc2.CommonResult

	// ID of the created message.
	MessageId base2.Id `json:"messageId"`
}

type ChangeMessageTextParams struct {
	rpc2.CommonParams

	MessageId base2.Id `json:"messageId"`

	// New text.
	Text base2.Text `json:"text"`
}
type ChangeMessageTextResult = rpc2.CommonResultWithSuccess

type ChangeMessageThreadParams struct {
	rpc2.CommonParams

	MessageId base2.Id `json:"messageId"`

	// ID of a new parent thread.
	ThreadId base2.Id `json:"threadId"`
}
type ChangeMessageThreadResult = rpc2.CommonResultWithSuccess

type GetMessageParams struct {
	rpc2.CommonParams

	MessageId base2.Id `json:"messageId"`
}
type GetMessageResult struct {
	rpc2.CommonResult

	Message derived2.IMessage `json:"message"`
}

type GetLatestMessageOfThreadParams struct {
	rpc2.CommonParams

	ThreadId base2.Id `json:"threadId"`
}
type GetLatestMessageOfThreadResult struct {
	rpc2.CommonResult

	Message derived2.IMessage `json:"message"`
}

type DeleteMessageParams struct {
	rpc2.CommonParams

	MessageId base2.Id `json:"messageId"`
}
type DeleteMessageResult = rpc2.CommonResultWithSuccess

// Composite objects.

type ListThreadAndMessagesParams struct {
	rpc2.CommonParams

	ThreadId base2.Id `json:"threadId"`
}
type ListThreadAndMessagesResult struct {
	rpc2.CommonResult

	ThreadAndMessages derived2.IThreadAndMessages `json:"tam"`
}

type ListThreadAndMessagesOnPageParams struct {
	rpc2.CommonParams

	ThreadId base2.Id    `json:"threadId"`
	Page     base2.Count `json:"page"`
}
type ListThreadAndMessagesOnPageResult struct {
	rpc2.CommonResult

	ThreadAndMessagesOnPage derived2.IThreadAndMessages `json:"tamop"`
}

type ListForumAndThreadsParams struct {
	rpc2.CommonParams

	ForumId base2.Id `json:"forumId"`
}
type ListForumAndThreadsResult struct {
	rpc2.CommonResult

	ForumAndThreads derived2.IForumAndThreads `json:"fat"`
}

type ListForumAndThreadsOnPageParams struct {
	rpc2.CommonParams

	ForumId base2.Id    `json:"forumId"`
	Page    base2.Count `json:"page"`
}
type ListForumAndThreadsOnPageResult struct {
	rpc2.CommonResult

	ForumAndThreadsOnPage derived2.IForumAndThreads `json:"fatop"`
}

type ListSectionsAndForumsParams struct {
	rpc2.CommonParams
}
type ListSectionsAndForumsResult struct {
	rpc2.CommonResult

	SectionsAndForums *models.SectionsAndForums `json:"saf"`
}

// Other.

type GetDKeyParams struct {
	rpc2.CommonParams
}
type GetDKeyResult struct {
	rpc2.CommonResult

	DKey base2.Text `json:"dKey"`
}

type ShowDiagnosticDataParams struct{}
type ShowDiagnosticDataResult struct {
	rpc2.CommonResult
	rpc2.RequestsCount
}

type TestParams struct {
	N uint `json:"n"`
}
type TestResult struct {
	rpc2.CommonResult
}
