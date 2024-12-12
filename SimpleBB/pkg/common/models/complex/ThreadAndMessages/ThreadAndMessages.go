package tam

import (
	"github.com/vault-thirteen/SimpleBB/pkg/common/interfaces/derived2"
	t "github.com/vault-thirteen/SimpleBB/pkg/common/models/complex/Thread"
	cmr "github.com/vault-thirteen/SimpleBB/pkg/common/models/rpc"
)

type threadAndMessages struct {
	Thread   derived2.IThread    `json:"thread"`
	Messages []derived2.IMessage `json:"messages"`
	PageData *cmr.PageData       `json:"pageData"`
}

func NewThreadAndMessages() (tam derived2.IThreadAndMessages) {
	return &threadAndMessages{
		Thread:   t.NewThread(),
		Messages: []derived2.IMessage{},
		PageData: &cmr.PageData{},
	}
}

// Emulated class members.
func (tam *threadAndMessages) GetThread() (thread derived2.IThread)        { return tam.Thread }
func (tam *threadAndMessages) SetThread(thread derived2.IThread)           { tam.Thread = thread }
func (tam *threadAndMessages) GetMessages() (messages []derived2.IMessage) { return tam.Messages }
func (tam *threadAndMessages) SetMessages(messages []derived2.IMessage)    { tam.Messages = messages }
func (tam *threadAndMessages) GetPageData() (pageData *cmr.PageData)       { return tam.PageData }
func (tam *threadAndMessages) SetPageData(pageData *cmr.PageData)          { tam.PageData = pageData }
