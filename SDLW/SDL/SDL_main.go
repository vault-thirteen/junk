package sdl

// SDL_main.h.

import (
	"syscall"
	"unsafe"

	m "github.com/vault-thirteen/SDLW/SDL/model"
)

/* Automatically added functions */

//extern DECLSPEC void SDLCALL SDL_SetMainReady(void);
func SetMainReady() {
	_, _, _ = syscall.SyscallN(fnSetMainReady)
}

//extern DECLSPEC int SDLCALL SDL_RegisterApp(const char *name, Uint32 style, void *hInst);
func RegisterApp(name string, style m.Uint32, hInst *m.Void) m.Int {
	ret, _, _ := syscall.SyscallN(fnRegisterApp, uintptr(unsafe.Pointer(BytePtrFromStringP(name))), uintptr(style), uintptr(unsafe.Pointer(hInst)))
	return (m.Int)(ret)
}

//extern DECLSPEC void SDLCALL SDL_UnregisterApp(void);
func UnregisterApp() {
	_, _, _ = syscall.SyscallN(fnUnregisterApp)
}
