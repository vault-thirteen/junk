package derived1

import "github.com/vault-thirteen/SimpleBB/pkg/ACM/models"

type IUserData interface {
	// Emulated class members.
	GetUser() IUser
	SetUser(user IUser)
	GetSession() *models.Session
	SetSession(session *models.Session)
}
