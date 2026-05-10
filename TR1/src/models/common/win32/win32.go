package win32

import (
	"github.com/vault-thirteen/TR1/src/models/common/dll"
	bt "github.com/vault-thirteen/auxie/BasicTypes"
	"golang.org/x/sys/windows"
)

const (
	WinKernelDll      = "kernel32.dll"
	User32Dll         = "user32.dll"
	DllFuncNamePrefix = ""
)

// typedef unsigned long DWORD; // minwindef.h:156.
type DWORD bt.DWord

var kernelFuncs = []dll.FuncMapping{
	{&fnGetLastError, "GetLastError"},
	{&fnSetLastError, "SetLastError"},
	{&fnGetStdHandle, "GetStdHandle"},
	{&fnGetConsoleMode, "GetConsoleMode"},
	{&fnSetConsoleMode, "SetConsoleMode"},
}

var user32Funcs = []dll.FuncMapping{}

var (
	kernelDll      windows.Handle
	fnGetLastError uintptr
	fnSetLastError uintptr

	user32Dll        windows.Handle
	fnGetStdHandle   uintptr
	fnGetConsoleMode uintptr
	fnSetConsoleMode uintptr
)

// LoadLibrary loads the library and its functions.
func LoadLibrary() (err error) {
	err = dll.LoadLibrary(WinKernelDll, &kernelDll, kernelFuncs, DllFuncNamePrefix)
	if err != nil {
		return err
	}

	err = dll.LoadLibrary(User32Dll, &user32Dll, user32Funcs, DllFuncNamePrefix)
	if err != nil {
		return err
	}

	return nil
}
