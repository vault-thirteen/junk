package request

type UserDisablingRequestUser struct {
	InternalName string `json:"internal_name"`
	Password     string `json:"password"`
	SecretCode   string `json:"secret_code"`
}
