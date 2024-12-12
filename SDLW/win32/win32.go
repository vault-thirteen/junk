package win32

import (
	"github.com/vault-thirteen/SDLW/dll"
	"golang.org/x/sys/windows"
)

const (
	WinKernelDll      = "kernel32.dll"
	DllFuncNamePrefix = ""
)

var funcs = []dll.FuncMapping{
	{&fnGetLastError, "GetLastError"},
	{&fnSetLastError, "SetLastError"},
}

var (
	kernelDll      windows.Handle
	fnGetLastError uintptr
	fnSetLastError uintptr
)

// LoadLibrary loads the library and its functions.
func LoadLibrary() (err error) {
	err = dll.LoadLibrary(WinKernelDll, &kernelDll, funcs, DllFuncNamePrefix)
	if err != nil {
		return err
	}

	return nil
}
