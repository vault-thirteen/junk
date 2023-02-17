package request

import (
	"github.com/vault-thirteen/junk/SSE1/pkg/models/user"
)

type UserLogOutRequest struct {
	User UserLogOutRequestUser `json:"user"`
}

func NewUserLogOutRequest(
	httpRequestBody []byte,
) (request *UserLogOutRequest, err error) {
	request = new(UserLogOutRequest)
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

func (ulor *UserLogOutRequest) Validate() (err error) {
	err = user.ValidateUserAuthenticationName(ulor.User.InternalName)
	if err != nil {
		return
	}
	return
}
