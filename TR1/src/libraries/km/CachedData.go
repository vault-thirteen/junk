package km

import "time"

type CachedData struct {
	UserId         int
	SessionId      int
	ExpirationTime int64
}

func (cd *CachedData) IsGood() bool {
	return time.Now().Unix() < cd.ExpirationTime
}
