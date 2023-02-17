package request

import (
	_ "github.com/mailru/easyjson/gen"
	"github.com/vault-thirteen/junk/SSE1/pkg/models/user"
)

type UserDisablingRequest struct {
	User UserDisablingRequestUser `json:"user"`
}

func NewUserDisablingRequest(
	httpRequestBody []byte,
) (request *UserDisablingRequest, err error) {
	request = new(UserDisablingRequest)
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

func (urr *UserDisablingRequest) Validate() (err error) {
	err = user.ValidateUserAuthenticationName(urr.User.InternalName)
	if err != nil {
		return
	}
	err = user.ValidateUserAuthenticationPassword(urr.User.Password)
	if err != nil {
		return
	}
	err = user.ValidateUserRegistrationSecretCode(urr.User.SecretCode)
	if err != nil {
		return
	}
	return
}
