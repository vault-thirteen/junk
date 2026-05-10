package cm

type Roles struct {
	CanLogIn        bool `json:"canLogIn"`
	CanRead         bool `json:"canRead"`
	CanWriteMessage bool `json:"canWriteMessage"`
	CanCreateThread bool `json:"canCreateThread"`
	IsModerator     bool `json:"isModerator" gorm:"-"`
	IsAdministrator bool `json:"isAdministrator" gorm:"-"`
}
