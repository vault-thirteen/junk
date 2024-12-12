package main

import (
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/vault-thirteen/SDLW/SDL"
	"github.com/vault-thirteen/SDLW/win32"
)

func main() {
	var err error
	err = win32.LoadLibrary()
	mustBeNoError(err)
	err = sdl.LoadLibrary(sdl.SdlDll)
	mustBeNoError(err)

	err = work()
	mustBeNoError(err)
}

func mustBeNoError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func work() (err error) {
	// Initialization and status check.
	if sdl.Init(sdl.INIT_EVERYTHING) != 0 {
		return sdl.GetError()
	}
	defer sdl.Quit()

	status := sdl.WasInit(0)
	fmt.Println(fmt.Sprintf("Initialization status: %v.", status))

	// Configuration Variables.
	sdl.SetHint(sdl.HINT_ACCELEROMETER_AS_JOYSTICK, "1")
	_ = sdl.GetHint(sdl.HINT_ACCELEROMETER_AS_JOYSTICK)
	sdl.ResetHints()

	// Error Handling.
	sdl.SetError("Test: %d.", 123)
	manualErr := sdl.GetError()
	if manualErr.Error() != "Test: 123." {
		return errors.New("errors do not work")
	}
	sdl.ClearError()
	manualErr = sdl.GetError()
	if manualErr != nil {
		return errors.New("errors do not work")
	}

	time.Sleep(time.Second * 1)

	// Log Handling.
	sdl.Log("Hello, World ! N=%d.", 123)

	// Querying SDL Version.
	fmt.Println(sdl.GetVersion().Text())
	fmt.Println(sdl.GetRevision())

	return nil
}
