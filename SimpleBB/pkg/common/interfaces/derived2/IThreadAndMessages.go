package derived2

import cmr "github.com/vault-thirteen/SimpleBB/pkg/common/models/rpc"

type IThreadAndMessages interface {
	// Emulated class members.
	GetThread() (thread IThread)
	SetThread(thread IThread)
	GetMessages() (messages []IMessage)
	SetMessages(messages []IMessage)
	GetPageData() (pageData *cmr.PageData)
	SetPageData(pageData *cmr.PageData)
}
