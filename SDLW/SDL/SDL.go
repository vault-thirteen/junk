package sdl

// SDL.h.

import (
	"syscall"

	m "github.com/vault-thirteen/SDLW/SDL/model"
)

const (
	INIT_TIMER          m.Uint32 = 0x00000001
	INIT_AUDIO          m.Uint32 = 0x00000010
	INIT_VIDEO          m.Uint32 = 0x00000020
	INIT_JOYSTICK       m.Uint32 = 0x00000200
	INIT_HAPTIC         m.Uint32 = 0x00001000
	INIT_GAMECONTROLLER m.Uint32 = 0x00002000
	INIT_EVENTS         m.Uint32 = 0x00004000
	INIT_SENSOR         m.Uint32 = 0x00008000
	INIT_NOPARACHUTE    m.Uint32 = 0x00100000
	INIT_EVERYTHING     m.Uint32 = INIT_TIMER | INIT_AUDIO | INIT_VIDEO | INIT_EVENTS |
		INIT_JOYSTICK | INIT_HAPTIC | INIT_GAMECONTROLLER | INIT_SENSOR
)

/* Automatically added functions */

//extern DECLSPEC int SDLCALL SDL_Init(Uint32 flags);
func Init(flags m.Uint32) m.Int {
	ret, _, _ := syscall.SyscallN(fnInit, uintptr(flags))
	return (m.Int)(ret)
}

//extern DECLSPEC int SDLCALL SDL_InitSubSystem(Uint32 flags);
func InitSubSystem(flags m.Uint32) m.Int {
	ret, _, _ := syscall.SyscallN(fnInitSubSystem, uintptr(flags))
	return (m.Int)(ret)
}

//extern DECLSPEC void SDLCALL SDL_QuitSubSystem(Uint32 flags);
func QuitSubSystem(flags m.Uint32) {
	_, _, _ = syscall.SyscallN(fnQuitSubSystem, uintptr(flags))
}

//extern DECLSPEC Uint32 SDLCALL SDL_WasInit(Uint32 flags);
func WasInit(flags m.Uint32) m.Uint32 {
	ret, _, _ := syscall.SyscallN(fnWasInit, uintptr(flags))
	return (m.Uint32)(ret)
}

//extern DECLSPEC void SDLCALL SDL_Quit(void);
func Quit() {
	_, _, _ = syscall.SyscallN(fnQuit)
}
