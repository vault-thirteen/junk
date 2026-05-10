package cm

import "net"

type RegistrationRequest struct {
	MetaData
	Id int `gorm:"primarykey" json:"id"`

	// Fields requested by a potential user.
	UserName     string `gorm:"uniqueIndex,size:255" json:"userName"`
	UserEmail    string `gorm:"uniqueIndex,size:255" json:"userEmail"`
	UserPassword string `json:"userPassword,omitempty"`

	// System fields.
	RequestId          string `gorm:"uniqueIndex,size:255" json:"requestId"`
	UserIPAB           net.IP `json:"userIPAB"`
	CaptchaId          string `json:"captchaId"`
	VerificationCode   string `json:"verificationCode"`
	IsReadyForApproval bool   `gorm:"index" json:"isReadyForApproval"`
}
