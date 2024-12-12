package win32

import (
	"fmt"
)

func EnableVirtualTerminalProcessing() (err error) {
	var errNum uint32

	// Get the Mode.
	var streamStdOut = GetStdHandle(STD_OUTPUT_HANDLE)
	var outConsoleMode DWORD
	if !GetConsoleMode(streamStdOut, &outConsoleMode) {
		errNum = GetLastError()
		return fmt.Errorf("failed to get output console mode, error %d", errNum)
	}

	// Update the Mode.
	outConsoleMode = outConsoleMode | ENABLE_VIRTUAL_TERMINAL_PROCESSING
	if !SetConsoleMode(streamStdOut, outConsoleMode) {
		errNum = GetLastError()
		return fmt.Errorf("failed to set output console mode, error %d", errNum)
	}

	return nil
}
