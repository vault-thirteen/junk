package win32

import (
	"syscall"

	"github.com/vault-thirteen/TR1/src/models/common/dll"
	"golang.org/x/sys/windows"
)

// typedef void *HANDLE; // WinBase.h:712.
type HANDLE uintptr //unsafe.Pointer

// #define INVALID_HANDLE_VALUE ((HANDLE)(LONG_PTR)-1)
const INVALID_HANDLE_VALUE = windows.InvalidHandle

// WINBASEAPI HANDLE WINAPI GetStdHandle(_In_ DWORD nStdHandle);
func GetStdHandle(handle DWORD) (h windows.Handle) {
	ret, _, callErr := syscall.SyscallN(fnGetStdHandle, uintptr(handle))
	dll.MustBeNoCallError(callErr)
	return windows.Handle(ret)
}
