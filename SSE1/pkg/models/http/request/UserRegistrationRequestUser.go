package request

type UserRegistrationRequestUser struct {
	InternalName string `json:"internal_name"`
	PublicName   string `json:"public_name"`
	Password     string `json:"password"`
	SecretCode   string `json:"secret_code"`
}
