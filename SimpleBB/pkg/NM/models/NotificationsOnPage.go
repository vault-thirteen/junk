package models

import (
	"github.com/vault-thirteen/SimpleBB/pkg/common/models/UidList"
	cmr "github.com/vault-thirteen/SimpleBB/pkg/common/models/rpc"
)

type NotificationsOnPage struct {
	// Notification parameters. If pagination is used, these lists contain
	// information after the application of pagination.
	NotificationIds *ul.UidList    `json:"notificationIds"`
	Notifications   []Notification `json:"notifications"`
	PageData        *cmr.PageData  `json:"pageData,omitempty"`
}

func NewNotificationsOnPage() (nop *NotificationsOnPage) {
	nop = &NotificationsOnPage{}
	return nop
}
