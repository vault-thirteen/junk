package sdl

// SDL_error.h.

import (
	"errors"
	"syscall"
	"unsafe"

	m "github.com/vault-thirteen/SDLW/SDL/model"
	"golang.org/x/sys/windows"
)

/* Manually added functions */

// GetError
// const char* SDL_GetError(void);
// https://wiki.libsdl.org/SDL2/SDL_GetError
func GetError() (err error) {
	cpErrText, _, _ := syscall.SyscallN(fnGetError)

	errText := windows.BytePtrToString((*byte)(unsafe.Pointer(cpErrText)))
	if len(errText) == 0 {
		return nil
	}

	return errors.New(errText)
}

// SetError
// int SDL_SetError(SDL_PRINTF_FORMAT_STRING const char *fmt, ...) SDL_PRINTF_VARARG_FUNC(1);
// https://wiki.libsdl.org/SDL2/SDL_SetError
func SetError(errFormat string, args ...uintptr) {
	// While this function always returns -1, the return is useless.

	// Golang wants all the args to be in a single array. Let it be so.
	argz := make([]uintptr, 0, len(args)+1)
	argz = append(argz, uintptr(unsafe.Pointer(BytePtrFromStringP(errFormat))))
	for _, arg := range args {
		argz = append(argz, arg)
	}

	_, _, _ = syscall.SyscallN(fnSetError, argz...)
}

/* Automatically added functions */

//extern DECLSPEC char * SDLCALL SDL_GetErrorMsg(char *errstr, int maxlen);
func GetErrorMsg(errstr string, maxlen m.Int) string {
	ret, _, _ := syscall.SyscallN(fnGetErrorMsg, uintptr(unsafe.Pointer(BytePtrFromStringP(errstr))), uintptr(maxlen))
	return windows.BytePtrToString((*byte)(unsafe.Pointer(ret)))
}

//extern DECLSPEC void SDLCALL SDL_ClearError(void);
func ClearError() {
	_, _, _ = syscall.SyscallN(fnClearError)
}

//extern DECLSPEC int SDLCALL SDL_Error(SDL_errorcode code);
func Error(code m.Errorcode) m.Int {
	ret, _, _ := syscall.SyscallN(fnError, uintptr(code))
	return (m.Int)(ret)
}
