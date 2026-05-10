package win32

import (
	"syscall"
	"unsafe"

	"github.com/vault-thirteen/TR1/src/models/common/dll"
	"golang.org/x/sys/windows"
)

// #define STD_OUTPUT_HANDLE ((DWORD)-11) // WinBase.h:820.
const STD_OUTPUT_HANDLE = 1<<32 - 11

// #define ENABLE_VIRTUAL_TERMINAL_PROCESSING  0x0004
const ENABLE_VIRTUAL_TERMINAL_PROCESSING = 4

// WINBASEAPI BOOL WINAPI GetConsoleMode(_In_ HANDLE hConsoleHandle, _Out_ LPDWORD lpMode);
func GetConsoleMode(consoleHandle windows.Handle, lpMode *DWORD) bool {
	ret, _, callErr := syscall.SyscallN(fnGetConsoleMode, uintptr(consoleHandle), uintptr(unsafe.Pointer(lpMode)))
	dll.MustBeNoCallError(callErr)
	return ret != 0
}

// WINBASEAPI BOOL WINAPI SetConsoleMode(_In_ HANDLE hConsoleHandle, _In_ DWORD dwMode);
func SetConsoleMode(consoleHandle windows.Handle, dwMode DWORD) bool {
	ret, _, callErr := syscall.SyscallN(fnSetConsoleMode, uintptr(consoleHandle), uintptr(dwMode))
	dll.MustBeNoCallError(callErr)
	return ret != 0
}
