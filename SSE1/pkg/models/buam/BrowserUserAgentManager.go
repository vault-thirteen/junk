package buam

import (
	"errors"
	"strconv"

	cache "github.com/vault-thirteen/Cache"
	"github.com/vault-thirteen/junk/SSE1/pkg/interfaces/storage"
)

// Errors.
const (
	ErrStorageIsNotSet = "Storage is not set"
	ErrTypeCasting     = "Type Casting Failure"
)

// A Manager working with Browsers' 'User Agent' Fields.
type BrowserUserAgentManager struct {
	storage storage.IStorage
	cache   *cache.Cache[string, string]
}

// Manager's Constructor.
func NewBrowserUserAgentManager(
	storage storage.IStorage,
) (buam *BrowserUserAgentManager, err error) {
	if storage == nil {
		err = errors.New(ErrStorageIsNotSet)
		return
	}
	buam = &BrowserUserAgentManager{
		storage: storage,
		cache:   cache.NewCache[string, string](100, 0, 1*60),
	}
	return
}

// Returns an Id of textual Representation of a 'User Agent' Setting.
// Uses the in-Memory Cache and the external Storage.
func (buam *BrowserUserAgentManager) GetBrowserUserAgentId(
	browserUserAgentName string,
) (id uint, err error) {

	// Try to get an Id from the Cache.
	var idAsInterface interface{}
	idAsInterface, err = buam.cache.GetRecord(browserUserAgentName)
	if err == nil {
		var ok bool
		id, ok = idAsInterface.(uint)
		if !ok {
			err = errors.New(ErrTypeCasting)
			return
		}
		return
	}

	// Search in the Database and update the Cache.
	id, err = buam.storage.GetBrowserUserAgentId(browserUserAgentName)
	if err != nil {
		return
	}
	err = buam.cache.AddRecord(browserUserAgentName, strconv.FormatUint(uint64(id), 10))
	if err != nil {
		return
	}
	return
}
