package win32

import (
	"syscall"

	"github.com/vault-thirteen/TR1/src/models/common/dll"
	bt "github.com/vault-thirteen/auxie/BasicTypes"
)

// GetLastError
// _Post_equals_last_error_ DWORD GetLastError();
// https://learn.microsoft.com/en-us/windows/win32/api/errhandlingapi/nf-errhandlingapi-getlasterror
func GetLastError() (err bt.DWord) {
	ret, _, callErr := syscall.SyscallN(fnGetLastError)
	dll.MustBeNoCallError(callErr)
	return bt.DWord(ret)
}

// SetLastError
// void SetLastError( [in] DWORD dwErrCode );
// https://learn.microsoft.com/en-us/windows/win32/api/errhandlingapi/nf-errhandlingapi-setlasterror
func SetLastError(err bt.DWord) {
	_, _, callErr := syscall.SyscallN(fnSetLastError, uintptr(err))
	dll.MustBeNoCallError(callErr)
}
