package dll

import (
	"errors"
	"fmt"
	"runtime"

	"golang.org/x/sys/windows"
)

const (
	ErrDllFileIsNotSet = "DLL file is not set"
)

// LoadLibrary loads the DLL library and its functions.
func LoadLibrary(
	dllFile string,
	h *windows.Handle,
	funcMappings []FuncMapping,
	funcNamePrefix string,
) (err error) {
	if len(dllFile) == 0 {
		return errors.New(ErrDllFileIsNotSet)
	}

	fmt.Println(fmt.Sprintf("Loading library: %v.", dllFile))
	*h, err = windows.LoadLibrary(dllFile)
	if err != nil {
		return err
	}

	fmt.Printf("Loading functions: ")
	for _, fm := range funcMappings {
		fmt.Printf("[%s] ", fm.FunctionName)
		*(fm.Fn), err = windows.GetProcAddress(*h, funcNamePrefix+fm.FunctionName)
		if err != nil {
			return err
		}
	}
	fmt.Println()

	return nil
}

// MustBeNoCallError checks for a call error.
// If call error is set, it finds a name of the function which started this
// function and panics about that function's failure.
func MustBeNoCallError(callErr windows.Errno) {
	if !errors.Is(callErr, windows.ERROR_SUCCESS) {
		caller, _, _, _ := runtime.Caller(1) // Skip 1 function up in stack: mustBeNoCallError -> caller.
		fn := runtime.FuncForPC(caller)
		panic(fmt.Sprintf("%s syscall failed: %v", fn.Name(), callErr.Error()))
	}
}
