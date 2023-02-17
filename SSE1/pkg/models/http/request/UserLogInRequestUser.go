package request

type UserLogInRequestUser struct {
	InternalName string `json:"internal_name"`
	Password     string `json:"password"`
}
