package nt

import (
	cmb "github.com/vault-thirteen/SimpleBB/pkg/common/models/base"
)

type NotificationTemplateArguments struct {
	MessageId  *cmb.Id `json:"messageId"`
	ResourceId *cmb.Id `json:"resourceId"`
	ThreadId   *cmb.Id `json:"threadId"`
	UserId     *cmb.Id `json:"userId"`
}
