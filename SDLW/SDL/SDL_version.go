package sdl

// SDL_version.h.

import (
	"syscall"
	"unsafe"

	m "github.com/vault-thirteen/SDLW/SDL/model"
	"golang.org/x/sys/windows"
)

/* Manually added functions */

//extern DECLSPEC void SDLCALL SDL_GetVersion(SDL_version * ver);
func GetVersion() (ver *m.Version) {
	ver = new(m.Version)
	_, _, _ = syscall.SyscallN(fnGetVersion, uintptr(unsafe.Pointer(ver)))
	return ver
}

/* Automatically added functions */

//extern DECLSPEC const char *SDLCALL SDL_GetRevision(void);
func GetRevision() string {
	ret, _, _ := syscall.SyscallN(fnGetRevision)
	return windows.BytePtrToString((*byte)(unsafe.Pointer(ret)))
}

//extern SDL_DEPRECATED DECLSPEC int SDLCALL SDL_GetRevisionNumber(void);
func GetRevisionNumber() m.Int {
	ret, _, _ := syscall.SyscallN(fnGetRevisionNumber)
	return (m.Int)(ret)
}
