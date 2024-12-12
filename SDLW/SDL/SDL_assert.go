package sdl

// SDL_assert.h.

import (
	"syscall"
	"unsafe"

	m "github.com/vault-thirteen/SDLW/SDL/model"
)

/* Automatically added functions */

//extern DECLSPEC void SDLCALL SDL_SetAssertionHandler(SDL_AssertionHandler handler, void *userdata);
func SetAssertionHandler(handler uintptr, userdata *m.Void) {
	_, _, _ = syscall.SyscallN(fnSetAssertionHandler, handler, uintptr(unsafe.Pointer(userdata)))
}

//extern DECLSPEC SDL_AssertionHandler SDLCALL SDL_GetDefaultAssertionHandler(void);
func GetDefaultAssertionHandler() m.AssertionHandler {
	ret, _, _ := syscall.SyscallN(fnGetDefaultAssertionHandler)
	return (m.AssertionHandler)(unsafe.Pointer(ret))
}

//extern DECLSPEC SDL_AssertionHandler SDLCALL SDL_GetAssertionHandler(void **puserdata);
func GetAssertionHandler(puserdata **m.Void) m.AssertionHandler {
	ret, _, _ := syscall.SyscallN(fnGetAssertionHandler, uintptr(unsafe.Pointer(puserdata)))
	return (m.AssertionHandler)(unsafe.Pointer(ret))
}

//extern DECLSPEC const SDL_AssertData * SDLCALL SDL_GetAssertionReport(void);
func GetAssertionReport() *m.AssertData {
	ret, _, _ := syscall.SyscallN(fnGetAssertionReport)
	return (*m.AssertData)(unsafe.Pointer(ret))
}

//extern DECLSPEC void SDLCALL SDL_ResetAssertionReport(void);
func ResetAssertionReport() {
	_, _, _ = syscall.SyscallN(fnResetAssertionReport)
}
