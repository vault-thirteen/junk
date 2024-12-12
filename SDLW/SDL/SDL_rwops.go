package sdl

// SDL_rwops.h.

import (
	"syscall"
	"unsafe"

	m "github.com/vault-thirteen/SDLW/SDL/model"
	"golang.org/x/sys/windows"
)

// RWFromFile
/*
SDL_RWops* SDL_RWFromFile(const char *file,
                          const char *mode);
*/
// https://wiki.libsdl.org/SDL2/SDL_RWFromFile
func RWFromFile(file string, mode string) *m.RWops {
	var err error
	var cpFile *byte
	cpFile, err = windows.BytePtrFromString(file)
	mustBeNoError(err)
	var cpMode *byte
	cpMode, err = windows.BytePtrFromString(mode)
	mustBeNoError(err)

	ret, _, _ := syscall.SyscallN(fnRWFromFile, uintptr(unsafe.Pointer(cpFile)), uintptr(unsafe.Pointer(cpMode)))
	return (*m.RWops)(unsafe.Pointer(ret))
}
