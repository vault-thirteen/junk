package s

import (
	"errors"
	base2 "github.com/vault-thirteen/SimpleBB/pkg/common/models/base"
	c "github.com/vault-thirteen/SimpleBB/pkg/common/models/server"
)

// SystemSettings are system settings.
type SystemSettings struct {
	DKeySize        base2.Count `json:"dKeySize"`
	MessageEditTime base2.Count `json:"messageEditTime"`
	PageSize        base2.Count `json:"pageSize"`

	// NewThreadsAtTop parameter controls how new and updated threads are
	// placed inside forums. If set to 'True', then following will happen:
	// 1. New threads will be added to the start (top) of the list of forum's
	// threads instead of being added to the end (bottom) of the list;
	// 2. New messages added to threads will update the thread moving it to the
	// start (top) position of the list of forum's threads.
	// If set to 'False', then new threads are added to the end (bottom) of the
	// list and thread's new messages do not update thread's position in the
	// list.
	NewThreadsAtTop base2.Flag `json:"newThreadsAtTop"`

	IsDebugMode base2.Flag `json:"isDebugMode"`
}

func (s SystemSettings) Check() (err error) {
	if (s.DKeySize == 0) ||
		(s.MessageEditTime == 0) ||
		(s.PageSize == 0) {
		return errors.New(c.MsgSystemSettingError)
	}

	return nil
}
