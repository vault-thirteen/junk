package sdlm

import (
	"errors"
	"syscall"
	"unsafe"

	"golang.org/x/sys/windows"
)

// GetError
// const char* SDL_GetError(void);
// https://wiki.libsdl.org/SDL2/SDL_GetError
// #define Mix_GetError    SDL_GetError
func GetError() (err error) {
	cpErrText, _, _ := syscall.SyscallN(extFnGetError)

	errText := windows.BytePtrToString((*byte)(unsafe.Pointer(cpErrText)))
	if len(errText) == 0 {
		return nil
	}

	return errors.New(errText)
}
