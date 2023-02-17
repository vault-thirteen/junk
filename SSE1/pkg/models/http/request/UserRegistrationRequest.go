package request

import (
	_ "github.com/mailru/easyjson/gen"
	"github.com/vault-thirteen/junk/SSE1/pkg/models/user"
)

type UserRegistrationRequest struct {
	User UserRegistrationRequestUser `json:"user"`
}

func NewUserRegistrationRequest(
	httpRequestBody []byte,
) (request *UserRegistrationRequest, err error) {
	request = new(UserRegistrationRequest)
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

func (urr *UserRegistrationRequest) Validate() (err error) {
	err = user.ValidateUserAuthenticationName(urr.User.InternalName)
	if err != nil {
		return
	}
	err = user.ValidateUserPublicName(urr.User.PublicName)
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
