package win32

import (
	"errors"

	"golang.org/x/sys/windows"
)

func Enable_console_colours() (err error) {
	var hConsole windows.Handle
	hConsole = GetStdHandle(STD_OUTPUT_HANDLE)
	if hConsole == INVALID_HANDLE_VALUE {
		return errors.New("GetStdHandle returned invalid handle")
	}

	var dwMode DWORD
	if !GetConsoleMode(hConsole, &dwMode) {
		return errors.New("GetConsoleMode failed")
	}

	dwMode |= ENABLE_VIRTUAL_TERMINAL_PROCESSING

	if !SetConsoleMode(hConsole, dwMode) {
		return errors.New("SetConsoleMode failed")
	}

	return nil
}
