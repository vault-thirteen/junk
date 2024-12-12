package fat

import (
	"github.com/vault-thirteen/SimpleBB/pkg/common/interfaces/derived2"
	f "github.com/vault-thirteen/SimpleBB/pkg/common/models/complex/Forum"
	cmr "github.com/vault-thirteen/SimpleBB/pkg/common/models/rpc"
)

type forumAndThreads struct {
	Forum    derived2.IForum    `json:"forum"`
	Threads  []derived2.IThread `json:"threads"`
	PageData *cmr.PageData      `json:"pageData,omitempty"`
}

func NewForumAndThreads() (fat derived2.IForumAndThreads) {
	return &forumAndThreads{
		Forum:    f.NewForum(),
		Threads:  []derived2.IThread{},
		PageData: &cmr.PageData{},
	}
}

// Emulated class members.
func (fat *forumAndThreads) GetForum() (forum derived2.IForum)        { return fat.Forum }
func (fat *forumAndThreads) SetForum(forum derived2.IForum)           { fat.Forum = forum }
func (fat *forumAndThreads) GetThreads() (threads []derived2.IThread) { return fat.Threads }
func (fat *forumAndThreads) SetThreads(threads []derived2.IThread)    { fat.Threads = threads }
func (fat *forumAndThreads) GetPageData() (pageData *cmr.PageData)    { return fat.PageData }
func (fat *forumAndThreads) SetPageData(pageData *cmr.PageData)       { fat.PageData = pageData }
