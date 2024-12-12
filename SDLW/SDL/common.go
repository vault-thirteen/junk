package sdl

import (
	"github.com/vault-thirteen/SDLW/dll"
	"golang.org/x/sys/windows"
)

const (
	SdlDll            = "SDL2.dll"
	DllFuncNamePrefix = "SDL_"
)

var funcs = []dll.FuncMapping{
	// Initialization and Shutdown.
	{&fnInit, "Init"},
	{&fnInitSubSystem, "InitSubSystem"},
	{&fnQuit, "Quit"},
	{&fnQuitSubSystem, "QuitSubSystem"},
	{&fnSetMainReady, "SetMainReady"},
	{&fnRegisterApp, "RegisterApp"},
	{&fnUnregisterApp, "UnregisterApp"},
	{&fnWasInit, "WasInit"},

	// Configuration Variables.
	{&fnSetHintWithPriority, "SetHintWithPriority"},
	{&fnSetHint, "SetHint"},
	{&fnResetHint, "ResetHint"},
	{&fnResetHints, "ResetHints"},
	{&fnGetHint, "GetHint"},
	{&fnGetHintBoolean, "GetHintBoolean"},
	{&fnAddHintCallback, "AddHintCallback"},
	{&fnDelHintCallback, "DelHintCallback"},
	{&fnClearHints, "ClearHints"},

	// Error Handling.
	{&fnClearError, "ClearError"},
	{&fnGetError, "GetError"},
	{&fnGetErrorMsg, "GetErrorMsg"},
	{&fnSetError, "SetError"},
	{&fnError, "Error"},

	// Log Handling.
	{&fnLogSetAllPriority, "LogSetAllPriority"},
	{&fnLogSetPriority, "LogSetPriority"},
	{&fnLogGetPriority, "LogGetPriority"},
	{&fnLogResetPriorities, "LogResetPriorities"},
	{&fnLog, "Log"},
	{&fnLogVerbose, "LogVerbose"},
	{&fnLogDebug, "LogDebug"},
	{&fnLogInfo, "LogInfo"},
	{&fnLogWarn, "LogWarn"},
	{&fnLogError, "LogError"},
	{&fnLogCritical, "LogCritical"},
	{&fnLogMessage, "LogMessage"},
	{&fnLogMessageV, "LogMessageV"},
	{&fnLogGetOutputFunction, "LogGetOutputFunction"},
	{&fnLogSetOutputFunction, "LogSetOutputFunction"},

	// Assertions.
	{&fnReportAssertion, "ReportAssertion"},
	{&fnSetAssertionHandler, "SetAssertionHandler"},
	{&fnGetDefaultAssertionHandler, "GetDefaultAssertionHandler"},
	{&fnGetAssertionHandler, "GetAssertionHandler"},
	{&fnGetAssertionReport, "GetAssertionReport"},
	{&fnResetAssertionReport, "ResetAssertionReport"},

	// Querying SDL Version.
	{&fnGetVersion, "GetVersion"},
	{&fnGetRevision, "GetRevision"},
	{&fnGetRevisionNumber, "GetRevisionNumber"},

	// Audio.
	{&fnGetNumAudioDrivers, "GetNumAudioDrivers"},
	{&fnGetAudioDriver, "GetAudioDriver"},
	{&fnAudioInit, "AudioInit"},
	{&fnAudioQuit, "AudioQuit"},
	{&fnGetCurrentAudioDriver, "GetCurrentAudioDriver"},
	{&fnOpenAudio, "OpenAudio"},
	{&fnGetNumAudioDevices, "GetNumAudioDevices"},
	{&fnGetAudioDeviceName, "GetAudioDeviceName"},
	{&fnGetAudioDeviceSpec, "GetAudioDeviceSpec"},
	{&fnGetDefaultAudioInfo, "GetDefaultAudioInfo"},
	{&fnOpenAudioDevice, "OpenAudioDevice"},
	{&fnGetAudioStatus, "GetAudioStatus"},
	{&fnGetAudioDeviceStatus, "GetAudioDeviceStatus"},
	{&fnPauseAudio, "PauseAudio"},
	{&fnPauseAudioDevice, "PauseAudioDevice"},
	{&fnLoadWAV_RW, "LoadWAV_RW"},
	{&fnFreeWAV, "FreeWAV"},
	{&fnBuildAudioCVT, "BuildAudioCVT"},
	{&fnConvertAudio, "ConvertAudio"},
	{&fnNewAudioStream, "NewAudioStream"},
	{&fnAudioStreamPut, "AudioStreamPut"},
	{&fnAudioStreamGet, "AudioStreamGet"},
	{&fnAudioStreamAvailable, "AudioStreamAvailable"},
	{&fnAudioStreamFlush, "AudioStreamFlush"},
	{&fnAudioStreamClear, "AudioStreamClear"},
	{&fnFreeAudioStream, "FreeAudioStream"},
	{&fnMixAudio, "MixAudio"},
	{&fnMixAudioFormat, "MixAudioFormat"},
	{&fnQueueAudio, "QueueAudio"},
	{&fnDequeueAudio, "DequeueAudio"},
	{&fnGetQueuedAudioSize, "GetQueuedAudioSize"},
	{&fnClearQueuedAudio, "ClearQueuedAudio"},
	{&fnLockAudio, "LockAudio"},
	{&fnLockAudioDevice, "LockAudioDevice"},
	{&fnUnlockAudio, "UnlockAudio"},
	{&fnUnlockAudioDevice, "UnlockAudioDevice"},
	{&fnCloseAudio, "CloseAudio"},
	{&fnCloseAudioDevice, "CloseAudioDevice"},

	// RWOps.
	{&fnRWFromFile, "RWFromFile"},
}

var (
	sdlDll windows.Handle

	// Initialization and Shutdown.
	fnInit          uintptr
	fnInitSubSystem uintptr
	fnQuit          uintptr
	fnQuitSubSystem uintptr
	fnSetMainReady  uintptr
	fnRegisterApp   uintptr
	fnUnregisterApp uintptr
	fnWasInit       uintptr

	// Configuration Variables.
	fnSetHintWithPriority uintptr
	fnSetHint             uintptr
	fnResetHint           uintptr
	fnResetHints          uintptr
	fnGetHint             uintptr
	fnGetHintBoolean      uintptr
	fnAddHintCallback     uintptr
	fnDelHintCallback     uintptr
	fnClearHints          uintptr

	// Error Handling.
	fnClearError  uintptr
	fnGetError    uintptr
	fnGetErrorMsg uintptr
	fnSetError    uintptr
	fnError       uintptr

	// Log Handling.
	fnLogSetAllPriority    uintptr
	fnLogSetPriority       uintptr
	fnLogGetPriority       uintptr
	fnLogResetPriorities   uintptr
	fnLog                  uintptr
	fnLogVerbose           uintptr
	fnLogDebug             uintptr
	fnLogInfo              uintptr
	fnLogWarn              uintptr
	fnLogError             uintptr
	fnLogCritical          uintptr
	fnLogMessage           uintptr
	fnLogMessageV          uintptr
	fnLogGetOutputFunction uintptr
	fnLogSetOutputFunction uintptr

	// Assertions.
	fnReportAssertion            uintptr
	fnSetAssertionHandler        uintptr
	fnGetDefaultAssertionHandler uintptr
	fnGetAssertionHandler        uintptr
	fnGetAssertionReport         uintptr
	fnResetAssertionReport       uintptr

	// Querying SDL Version.
	fnGetVersion        uintptr
	fnGetRevision       uintptr
	fnGetRevisionNumber uintptr

	// Audio.
	fnGetNumAudioDrivers    uintptr
	fnGetAudioDriver        uintptr
	fnAudioInit             uintptr
	fnAudioQuit             uintptr
	fnGetCurrentAudioDriver uintptr
	fnOpenAudio             uintptr
	fnGetNumAudioDevices    uintptr
	fnGetAudioDeviceName    uintptr
	fnGetAudioDeviceSpec    uintptr
	fnGetDefaultAudioInfo   uintptr
	fnOpenAudioDevice       uintptr
	fnGetAudioStatus        uintptr
	fnGetAudioDeviceStatus  uintptr
	fnPauseAudio            uintptr
	fnPauseAudioDevice      uintptr
	fnLoadWAV_RW            uintptr
	fnFreeWAV               uintptr
	fnBuildAudioCVT         uintptr
	fnConvertAudio          uintptr
	fnNewAudioStream        uintptr
	fnAudioStreamPut        uintptr
	fnAudioStreamGet        uintptr
	fnAudioStreamAvailable  uintptr
	fnAudioStreamFlush      uintptr
	fnAudioStreamClear      uintptr
	fnFreeAudioStream       uintptr
	fnMixAudio              uintptr
	fnMixAudioFormat        uintptr
	fnQueueAudio            uintptr
	fnDequeueAudio          uintptr
	fnGetQueuedAudioSize    uintptr
	fnClearQueuedAudio      uintptr
	fnLockAudio             uintptr
	fnLockAudioDevice       uintptr
	fnUnlockAudio           uintptr
	fnUnlockAudioDevice     uintptr
	fnCloseAudio            uintptr
	fnCloseAudioDevice      uintptr

	// RWOps.
	fnRWFromFile uintptr
)

// LoadLibrary loads the library and its functions.
func LoadLibrary(dllFile string) (err error) {
	err = dll.LoadLibrary(dllFile, &sdlDll, funcs, DllFuncNamePrefix)
	if err != nil {
		return err
	}

	return nil
}

// BytePtrFromStringP converts a Go string into a C string.
// If something goes wrong, it panics.
func BytePtrFromStringP(s string) (cpS *byte) {
	var err error
	cpS, err = windows.BytePtrFromString(s)
	mustBeNoError(err)
	return cpS
}

// mustBeNoError panics if an error occurs.
func mustBeNoError(err error) {
	if err != nil {
		panic(err)
	}
}

// GetFnGetError gets the handle of SDL 'GetError' function.
// This handle is used by various SDL extensions.
func GetFnGetError() uintptr {
	return fnGetError
}
