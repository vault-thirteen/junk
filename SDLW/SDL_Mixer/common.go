package sdlm

import (
	"errors"

	"github.com/vault-thirteen/SDLW/dll"
	"golang.org/x/sys/windows"
)

const (
	SdlMixDll         = "SDL2_mixer.dll"
	DllFuncNamePrefix = "Mix_"
)

const (
	ErrExternalFunctionsNotSet = "external functions are not set"
)

var funcs = []dll.FuncMapping{
	{&fnLinked_Version, "Linked_Version"},
	{&fnInit, "Init"},
	{&fnQuit, "Quit"},
	{&fnOpenAudio, "OpenAudio"},
	{&fnOpenAudioDevice, "OpenAudioDevice"},
	{&fnQuerySpec, "QuerySpec"},
	{&fnAllocateChannels, "AllocateChannels"},
	{&fnLoadWAV_RW, "LoadWAV_RW"},
	{&fnLoadWAV, "LoadWAV"},
	{&fnLoadMUS, "LoadMUS"},
	{&fnLoadMUS_RW, "LoadMUS_RW"},
	{&fnLoadMUSType_RW, "LoadMUSType_RW"},
	{&fnQuickLoad_WAV, "QuickLoad_WAV"},
	{&fnQuickLoad_RAW, "QuickLoad_RAW"},
	{&fnFreeChunk, "FreeChunk"},
	{&fnFreeMusic, "FreeMusic"},
	{&fnGetNumChunkDecoders, "GetNumChunkDecoders"},
	{&fnGetChunkDecoder, "GetChunkDecoder"},
	{&fnHasChunkDecoder, "HasChunkDecoder"},
	{&fnGetNumMusicDecoders, "GetNumMusicDecoders"},
	{&fnGetMusicDecoder, "GetMusicDecoder"},
	{&fnHasMusicDecoder, "HasMusicDecoder"},
	{&fnGetMusicType, "GetMusicType"},
	{&fnGetMusicTitle, "GetMusicTitle"},
	{&fnGetMusicTitleTag, "GetMusicTitleTag"},
	{&fnGetMusicArtistTag, "GetMusicArtistTag"},
	{&fnGetMusicAlbumTag, "GetMusicAlbumTag"},
	{&fnGetMusicCopyrightTag, "GetMusicCopyrightTag"},
	{&fnGetMusicHookData, "GetMusicHookData"},
	{&fnRegisterEffect, "RegisterEffect"},
	{&fnUnregisterEffect, "UnregisterEffect"},
	{&fnUnregisterAllEffects, "UnregisterAllEffects"},
	{&fnSetPanning, "SetPanning"},
	{&fnSetPosition, "SetPosition"},
	{&fnSetDistance, "SetDistance"},
	{&fnSetReverseStereo, "SetReverseStereo"},
	{&fnReserveChannels, "ReserveChannels"},
	{&fnGroupChannel, "GroupChannel"},
	{&fnGroupChannels, "GroupChannels"},
	{&fnGroupAvailable, "GroupAvailable"},
	{&fnGroupCount, "GroupCount"},
	{&fnGroupOldest, "GroupOldest"},
	{&fnGroupNewer, "GroupNewer"},
	{&fnPlayChannel, "PlayChannel"},
	{&fnPlayChannelTimed, "PlayChannelTimed"},
	{&fnPlayMusic, "PlayMusic"},
	{&fnFadeInMusic, "FadeInMusic"},
	{&fnFadeInMusicPos, "FadeInMusicPos"},
	{&fnFadeInChannel, "FadeInChannel"},
	{&fnFadeInChannelTimed, "FadeInChannelTimed"},
	{&fnVolume, "Volume"},
	{&fnVolumeChunk, "VolumeChunk"},
	{&fnVolumeMusic, "VolumeMusic"},
	{&fnGetMusicVolume, "GetMusicVolume"},
	{&fnMasterVolume, "MasterVolume"},
	{&fnHaltChannel, "HaltChannel"},
	{&fnHaltGroup, "HaltGroup"},
	{&fnHaltMusic, "HaltMusic"},
	{&fnExpireChannel, "ExpireChannel"},
	{&fnFadeOutChannel, "FadeOutChannel"},
	{&fnFadeOutGroup, "FadeOutGroup"},
	{&fnFadeOutMusic, "FadeOutMusic"},
	{&fnFadingMusic, "FadingMusic"},
	{&fnFadingChannel, "FadingChannel"},
	{&fnPause, "Pause"},
	{&fnResume, "Resume"},
	{&fnPaused, "Paused"},
	{&fnPauseMusic, "PauseMusic"},
	{&fnResumeMusic, "ResumeMusic"},
	{&fnRewindMusic, "RewindMusic"},
	{&fnPausedMusic, "PausedMusic"},
	{&fnModMusicJumpToOrder, "ModMusicJumpToOrder"},
	{&fnSetMusicPosition, "SetMusicPosition"},
	{&fnGetMusicPosition, "GetMusicPosition"},
	{&fnMusicDuration, "MusicDuration"},
	{&fnGetMusicLoopStartTime, "GetMusicLoopStartTime"},
	{&fnGetMusicLoopEndTime, "GetMusicLoopEndTime"},
	{&fnGetMusicLoopLengthTime, "GetMusicLoopLengthTime"},
	{&fnPlaying, "Playing"},
	{&fnPlayingMusic, "PlayingMusic"},
	{&fnSetMusicCMD, "SetMusicCMD"},
	{&fnSetSynchroValue, "SetSynchroValue"},
	{&fnGetSynchroValue, "GetSynchroValue"},
	{&fnSetSoundFonts, "SetSoundFonts"},
	{&fnGetSoundFonts, "GetSoundFonts"},
	{&fnSetTimidityCfg, "SetTimidityCfg"},
	{&fnGetTimidityCfg, "GetTimidityCfg"},
	{&fnGetChunk, "GetChunk"},
	{&fnCloseAudio, "CloseAudio"},

	// Manually added functions.
	{&fnSetPostMix, "SetPostMix"},
	{&fnHookMusic, "HookMusic"},
	{&fnHookMusicFinished, "HookMusicFinished"},
	{&fnChannelFinished, "ChannelFinished"},
	{&fnEachSoundFont, "EachSoundFont"},
}

var (
	sdlMixDll windows.Handle

	fnLinked_Version         uintptr
	fnInit                   uintptr
	fnQuit                   uintptr
	fnOpenAudio              uintptr
	fnOpenAudioDevice        uintptr
	fnQuerySpec              uintptr
	fnAllocateChannels       uintptr
	fnLoadWAV_RW             uintptr
	fnLoadWAV                uintptr
	fnLoadMUS                uintptr
	fnLoadMUS_RW             uintptr
	fnLoadMUSType_RW         uintptr
	fnQuickLoad_WAV          uintptr
	fnQuickLoad_RAW          uintptr
	fnFreeChunk              uintptr
	fnFreeMusic              uintptr
	fnGetNumChunkDecoders    uintptr
	fnGetChunkDecoder        uintptr
	fnHasChunkDecoder        uintptr
	fnGetNumMusicDecoders    uintptr
	fnGetMusicDecoder        uintptr
	fnHasMusicDecoder        uintptr
	fnGetMusicType           uintptr
	fnGetMusicTitle          uintptr
	fnGetMusicTitleTag       uintptr
	fnGetMusicArtistTag      uintptr
	fnGetMusicAlbumTag       uintptr
	fnGetMusicCopyrightTag   uintptr
	fnGetMusicHookData       uintptr
	fnRegisterEffect         uintptr
	fnUnregisterEffect       uintptr
	fnUnregisterAllEffects   uintptr
	fnSetPanning             uintptr
	fnSetPosition            uintptr
	fnSetDistance            uintptr
	fnSetReverseStereo       uintptr
	fnReserveChannels        uintptr
	fnGroupChannel           uintptr
	fnGroupChannels          uintptr
	fnGroupAvailable         uintptr
	fnGroupCount             uintptr
	fnGroupOldest            uintptr
	fnGroupNewer             uintptr
	fnPlayChannel            uintptr
	fnPlayChannelTimed       uintptr
	fnPlayMusic              uintptr
	fnFadeInMusic            uintptr
	fnFadeInMusicPos         uintptr
	fnFadeInChannel          uintptr
	fnFadeInChannelTimed     uintptr
	fnVolume                 uintptr
	fnVolumeChunk            uintptr
	fnVolumeMusic            uintptr
	fnGetMusicVolume         uintptr
	fnMasterVolume           uintptr
	fnHaltChannel            uintptr
	fnHaltGroup              uintptr
	fnHaltMusic              uintptr
	fnExpireChannel          uintptr
	fnFadeOutChannel         uintptr
	fnFadeOutGroup           uintptr
	fnFadeOutMusic           uintptr
	fnFadingMusic            uintptr
	fnFadingChannel          uintptr
	fnPause                  uintptr
	fnResume                 uintptr
	fnPaused                 uintptr
	fnPauseMusic             uintptr
	fnResumeMusic            uintptr
	fnRewindMusic            uintptr
	fnPausedMusic            uintptr
	fnModMusicJumpToOrder    uintptr
	fnSetMusicPosition       uintptr
	fnGetMusicPosition       uintptr
	fnMusicDuration          uintptr
	fnGetMusicLoopStartTime  uintptr
	fnGetMusicLoopEndTime    uintptr
	fnGetMusicLoopLengthTime uintptr
	fnPlaying                uintptr
	fnPlayingMusic           uintptr
	fnSetMusicCMD            uintptr
	fnSetSynchroValue        uintptr
	fnGetSynchroValue        uintptr
	fnSetSoundFonts          uintptr
	fnGetSoundFonts          uintptr
	fnSetTimidityCfg         uintptr
	fnGetTimidityCfg         uintptr
	fnGetChunk               uintptr
	fnCloseAudio             uintptr

	// Manually added functions.
	fnSetPostMix        uintptr
	fnHookMusic         uintptr
	fnHookMusicFinished uintptr
	fnChannelFinished   uintptr
	fnEachSoundFont     uintptr
)

var (
	extFnGetError uintptr
)

type ExternalFunctions struct {
	ExtFnGetError uintptr
}

// LoadLibrary loads the library and its functions.
func LoadLibrary(dllFile string, extFuncs *ExternalFunctions) (err error) {
	err = saveExternalFunctions(extFuncs)
	if err != nil {
		return err
	}

	err = dll.LoadLibrary(dllFile, &sdlMixDll, funcs, DllFuncNamePrefix)
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

// saveExternalFunctions saves external functions.
func saveExternalFunctions(extFuncs *ExternalFunctions) (err error) {
	if extFuncs == nil {
		return errors.New(ErrExternalFunctionsNotSet)
	}

	extFnGetError = extFuncs.ExtFnGetError

	return nil
}
