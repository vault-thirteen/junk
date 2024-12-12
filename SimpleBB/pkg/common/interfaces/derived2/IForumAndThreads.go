package derived2

import cmr "github.com/vault-thirteen/SimpleBB/pkg/common/models/rpc"

type IForumAndThreads interface {
	// Emulated class members.
	GetForum() (forum IForum)
	SetForum(forum IForum)
	GetThreads() (threads []IThread)
	SetThreads(threads []IThread)
	GetPageData() (pageData *cmr.PageData)
	SetPageData(pageData *cmr.PageData)
}
