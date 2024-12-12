package win32

import (
	"syscall"
	"unsafe"
)

var (
	modkernel32 = syscall.NewLazyDLL("kernel32.dll")

	procGetLastError   = modkernel32.NewProc("GetLastError")
	procGetStdHandle   = modkernel32.NewProc("GetStdHandle")
	procGetConsoleMode = modkernel32.NewProc("GetConsoleMode")
	procSetConsoleMode = modkernel32.NewProc("SetConsoleMode")
)

func GetLastError() uint32 {
	ret, _, _ := procGetLastError.Call()
	return uint32(ret)
}

func GetStdHandle(nStdHandle DWORD) HANDLE {
	ret, _, _ := procGetStdHandle.Call(
		uintptr(nStdHandle),
	)

	return HANDLE(ret)
}

func GetConsoleMode(hConsoleHandle HANDLE, lpMode *DWORD) bool {
	ret, _, _ := procGetConsoleMode.Call(
		uintptr(hConsoleHandle),
		uintptr(unsafe.Pointer(lpMode)))

	return ret != 0
}

func SetConsoleMode(hConsoleHandle HANDLE, dwMode DWORD) bool {
	ret, _, _ := procSetConsoleMode.Call(
		uintptr(hConsoleHandle),
		uintptr(dwMode))

	return ret != 0
}
