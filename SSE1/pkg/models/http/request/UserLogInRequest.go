package request

import (
	_ "github.com/mailru/easyjson/gen"
	"github.com/vault-thirteen/junk/SSE1/pkg/models/user"
)

type UserLogInRequest struct {
	User    UserLogInRequestUser  `json:"user"`
	Machine UserLogRequestMachine `json:"-"`
}

func NewUserLogInRequest(
	httpRequestBody []byte,

) (request *UserLogInRequest, err error) {
	request = new(UserLogInRequest)
	err = request.UnmarshalJSON(httpRequestBody)
	if err != nil {
		return
	}
	err = request.Validate()
	if err != nil {
		return
	}
	return
}

func (ulir *UserLogInRequest) Validate() (err error) {
	err = user.ValidateUserAuthenticationName(ulir.User.InternalName)
	if err != nil {
		return
	}
	err = user.ValidateUserAuthenticationPassword(ulir.User.Password)
	if err != nil {
		return
	}
	return
}
