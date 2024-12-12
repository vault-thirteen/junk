package ud

import (
	"github.com/vault-thirteen/SimpleBB/pkg/ACM/models"
	"github.com/vault-thirteen/SimpleBB/pkg/common/interfaces/derived1"
	"github.com/vault-thirteen/SimpleBB/pkg/common/models/complex/User"
)

type userData struct {
	User    derived1.IUser
	Session *models.Session
}

func NewUserData() (ud derived1.IUserData) {
	return &userData{
		User:    u.NewUser(),
		Session: &models.Session{},
	}
}

// Emulated class members.
func (ud *userData) GetUser() derived1.IUser {
	return ud.User
}
func (ud *userData) SetUser(user derived1.IUser) {
	ud.User = user
}
func (ud *userData) GetSession() *models.Session {
	return ud.Session
}
func (ud *userData) SetSession(session *models.Session) {
	ud.Session = session
}
