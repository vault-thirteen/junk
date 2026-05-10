package cm

import (
	"log"
	"runtime"
	"strings"
	"sync"
	"sync/atomic"

	"github.com/vault-thirteen/TR1/src/models/common/win32"
)

type DllLoader struct {
	guard              *sync.Mutex
	windowsApiIsLoaded atomic.Bool
}

func NewDllLoader() (dl *DllLoader) {
	dl = &DllLoader{
		guard: new(sync.Mutex),
	}

	dl.windowsApiIsLoaded.Store(false)

	return dl
}

func (dl *DllLoader) Init() (err error) {
	dl.guard.Lock()
	defer dl.guard.Unlock()

	os := runtime.GOOS
	log.Println("Operating system:", os)

	switch strings.ToLower(os) {
	case "windows":
		return dl.loadWindowsApi()
	default:
		log.Println("Current operating system is not fully supported. Some features may be unavailable.")
	}

	return nil
}

func (dl *DllLoader) loadWindowsApi() (err error) {
	if dl.windowsApiIsLoaded.Load() {
		return nil
	}

	err = win32.LoadLibrary()
	if err != nil {
		return err
	}

	dl.windowsApiIsLoaded.Store(true)

	return nil
}
