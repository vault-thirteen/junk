package sdl

// SDL_log.h.

import (
	"syscall"
	"unsafe"

	m "github.com/vault-thirteen/SDLW/SDL/model"
)

const MAX_LOG_MESSAGE = 4096

/* Manually added functions */

// Log
// void SDL_Log(SDL_PRINTF_FORMAT_STRING const char *fmt, ...) SDL_PRINTF_VARARG_FUNC(1);
// https://wiki.libsdl.org/SDL2/SDL_Log
func Log(fmt string, args ...uintptr) {
	// Golang wants all the args to be in a single array. Let it be so.
	argz := make([]uintptr, 0, len(args)+1)
	argz = append(argz, uintptr(unsafe.Pointer(BytePtrFromStringP(fmt))))
	for _, arg := range args {
		argz = append(argz, arg)
	}

	_, _, _ = syscall.SyscallN(fnLog, argz...)
}

// LogVerbose
// void SDL_LogVerbose(int category, SDL_PRINTF_FORMAT_STRING const char *fmt, ...) SDL_PRINTF_VARARG_FUNC(2);
// https://wiki.libsdl.org/SDL2/SDL_LogVerbose
func LogVerbose(category m.Int, fmt string, args ...uintptr) {
	// Golang wants all the args to be in a single array. Let it be so.
	argz := make([]uintptr, 0, len(args)+2)
	argz = append(argz, uintptr(category))
	argz = append(argz, uintptr(unsafe.Pointer(BytePtrFromStringP(fmt))))
	for _, arg := range args {
		argz = append(argz, arg)
	}

	_, _, _ = syscall.SyscallN(fnLogVerbose, argz...)
}

// LogDebug
// void SDL_LogDebug(int category, SDL_PRINTF_FORMAT_STRING const char *fmt, ...) SDL_PRINTF_VARARG_FUNC(2);
// https://wiki.libsdl.org/SDL2/SDL_LogDebug
func LogDebug(category m.Int, fmt string, args ...uintptr) {
	// Golang wants all the args to be in a single array. Let it be so.
	argz := make([]uintptr, 0, len(args)+2)
	argz = append(argz, uintptr(category))
	argz = append(argz, uintptr(unsafe.Pointer(BytePtrFromStringP(fmt))))
	for _, arg := range args {
		argz = append(argz, arg)
	}

	_, _, _ = syscall.SyscallN(fnLogDebug, argz...)
}

// LogInfo
// void SDL_LogInfo(int category, SDL_PRINTF_FORMAT_STRING const char *fmt, ...) SDL_PRINTF_VARARG_FUNC(2);
// https://wiki.libsdl.org/SDL2/SDL_LogInfo
func LogInfo(category m.Int, fmt string, args ...uintptr) {
	// Golang wants all the args to be in a single array. Let it be so.
	argz := make([]uintptr, 0, len(args)+2)
	argz = append(argz, uintptr(category))
	argz = append(argz, uintptr(unsafe.Pointer(BytePtrFromStringP(fmt))))
	for _, arg := range args {
		argz = append(argz, arg)
	}

	_, _, _ = syscall.SyscallN(fnLogInfo, argz...)
}

// LogWarn
// void SDL_LogWarn(int category, SDL_PRINTF_FORMAT_STRING const char *fmt, ...) SDL_PRINTF_VARARG_FUNC(2);
// https://wiki.libsdl.org/SDL2/SDL_LogWarn
func LogWarn(category m.Int, fmt string, args ...uintptr) {
	// Golang wants all the args to be in a single array. Let it be so.
	argz := make([]uintptr, 0, len(args)+2)
	argz = append(argz, uintptr(category))
	argz = append(argz, uintptr(unsafe.Pointer(BytePtrFromStringP(fmt))))
	for _, arg := range args {
		argz = append(argz, arg)
	}

	_, _, _ = syscall.SyscallN(fnLogWarn, argz...)
}

// LogError
// void SDL_LogError(int category, SDL_PRINTF_FORMAT_STRING const char *fmt, ...) SDL_PRINTF_VARARG_FUNC(2);
// https://wiki.libsdl.org/SDL2/SDL_LogError
func LogError(category m.Int, fmt string, args ...uintptr) {
	// Golang wants all the args to be in a single array. Let it be so.
	argz := make([]uintptr, 0, len(args)+2)
	argz = append(argz, uintptr(category))
	argz = append(argz, uintptr(unsafe.Pointer(BytePtrFromStringP(fmt))))
	for _, arg := range args {
		argz = append(argz, arg)
	}

	_, _, _ = syscall.SyscallN(fnLogError, argz...)
}

// LogCritical
// void SDL_LogCritical(int category, SDL_PRINTF_FORMAT_STRING const char *fmt, ...) SDL_PRINTF_VARARG_FUNC(2);
// https://wiki.libsdl.org/SDL2/SDL_LogCritical
func LogCritical(category m.Int, fmt string, args ...uintptr) {
	// Golang wants all the args to be in a single array. Let it be so.
	argz := make([]uintptr, 0, len(args)+2)
	argz = append(argz, uintptr(category))
	argz = append(argz, uintptr(unsafe.Pointer(BytePtrFromStringP(fmt))))
	for _, arg := range args {
		argz = append(argz, arg)
	}

	_, _, _ = syscall.SyscallN(fnLogCritical, argz...)
}

// LogMessage
/*
void SDL_LogMessage(int category,
                    SDL_LogPriority priority,
                    SDL_PRINTF_FORMAT_STRING const char *fmt, ...) SDL_PRINTF_VARARG_FUNC(3);
*/
// https://wiki.libsdl.org/SDL2/SDL_LogMessage
func LogMessage(category m.Int, priority m.LogPriority, fmt string, args ...uintptr) {
	// Golang wants all the args to be in a single array. Let it be so.
	argz := make([]uintptr, 0, len(args)+3)
	argz = append(argz, uintptr(category))
	argz = append(argz, uintptr(priority))
	argz = append(argz, uintptr(unsafe.Pointer(BytePtrFromStringP(fmt))))
	for _, arg := range args {
		argz = append(argz, arg)
	}

	_, _, _ = syscall.SyscallN(fnLogMessage, argz...)
}

/* Automatically added functions */

//extern DECLSPEC void SDLCALL SDL_LogSetAllPriority(SDL_LogPriority priority);
func LogSetAllPriority(priority m.LogPriority) {
	_, _, _ = syscall.SyscallN(fnLogSetAllPriority, uintptr(priority))
}

//extern DECLSPEC void SDLCALL SDL_LogSetPriority(int category, SDL_LogPriority priority);
func LogSetPriority(category m.Int, priority m.LogPriority) {
	_, _, _ = syscall.SyscallN(fnLogSetPriority, uintptr(category), uintptr(priority))
}

//extern DECLSPEC SDL_LogPriority SDLCALL SDL_LogGetPriority(int category);
func LogGetPriority(category m.Int) m.LogPriority {
	ret, _, _ := syscall.SyscallN(fnLogGetPriority, uintptr(category))
	return (m.LogPriority)(ret)
}

//extern DECLSPEC void SDLCALL SDL_LogResetPriorities(void);
func LogResetPriorities() {
	_, _, _ = syscall.SyscallN(fnLogResetPriorities)
}

//extern DECLSPEC void SDLCALL SDL_LogMessageV(int category, SDL_LogPriority priority, const char *fmt, va_list ap);
//func LogMessageV(category m.Int, priority m.LogPriority, fmt string, ap va_list) {
//	_, _, _ = syscall.SyscallN(fnLogMessageV, uintptr(category), uintptr(priority), uintptr(unsafe.Pointer(BytePtrFromStringP(fmt))), uintptr(ap))
//}

//extern DECLSPEC void SDLCALL SDL_LogGetOutputFunction(SDL_LogOutputFunction *callback, void **userdata);
func LogGetOutputFunction(callback *m.LogOutputFunction, userdata **m.Void) {
	_, _, _ = syscall.SyscallN(fnLogGetOutputFunction, uintptr(unsafe.Pointer(callback)), uintptr(unsafe.Pointer(userdata)))
}

//extern DECLSPEC void SDLCALL SDL_LogSetOutputFunction(SDL_LogOutputFunction callback, void *userdata);
func LogSetOutputFunction(callback m.LogOutputFunction, userdata *m.Void) {
	_, _, _ = syscall.SyscallN(fnLogSetOutputFunction, uintptr(unsafe.Pointer(callback)), uintptr(unsafe.Pointer(userdata)))
}
