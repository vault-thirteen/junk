package nt

import (
	nm "github.com/vault-thirteen/SimpleBB/pkg/NM/models"
	base2 "github.com/vault-thirteen/SimpleBB/pkg/common/models/base"
)

type NotificationTemplate struct {
	Id           base2.Id                       `json:"id"`
	Name         base2.Text                     `json:"name"`
	FormatString *nm.FormatString               `json:"formatString"`
	Arguments    *NotificationTemplateArguments `json:"arguments"`
}
