package sdl

// SDL_audio.h.

import (
	"syscall"
	"unsafe"

	"github.com/vault-thirteen/SDLW/SDL/model"
	"golang.org/x/sys/windows"
)

// Audio flags.
const (
	AUDIO_MASK_BITSIZE  = m.Uint16(0xFF)
	AUDIO_MASK_DATATYPE = m.Uint16(1 << 8)
	AUDIO_MASK_ENDIAN   = m.Uint16(1 << 12)
	AUDIO_MASK_SIGNED   = m.Uint16(1 << 15)
)

func AUDIO_BITSIZE(x m.Uint16) m.Uint16 {
	return x & AUDIO_MASK_BITSIZE
}
func AUDIO_ISFLOAT(x m.Uint16) m.Uint16 {
	return x & AUDIO_MASK_DATATYPE
}
func AUDIO_ISBIGENDIAN(x m.Uint16) m.Uint16 {
	return x & AUDIO_MASK_ENDIAN
}
func AUDIO_ISSIGNED(x m.Uint16) m.Uint16 {
	return x & AUDIO_MASK_SIGNED
}

// Audio format flags.
const (
	AUDIO_U8     = 0x0008
	AUDIO_S8     = 0x8008
	AUDIO_U16LSB = 0x0010
	AUDIO_S16LSB = 0x8010
	AUDIO_U16MSB = 0x1010
	AUDIO_S16MSB = 0x9010
	AUDIO_S16    = AUDIO_S16LSB
	AUDIO_U16    = AUDIO_U16LSB

	// Integer Samples.
	AUDIO_S32LSB = 0x8020
	AUDIO_S32MSB = 0x9020
	AUDIO_S32    = AUDIO_S32LSB

	// Float Samples.
	AUDIO_F32LSB = 0x8120
	AUDIO_F32MSB = 0x9120
	AUDIO_F32    = AUDIO_F32LSB

	MIX_MAXVOLUME = 128
)

// Which audio format changes are allowed when opening a device.
const (
	AUDIO_ALLOW_FREQUENCY_CHANGE = 0x00000001
	AUDIO_ALLOW_FORMAT_CHANGE    = 0x00000002
	AUDIO_ALLOW_CHANNELS_CHANGE  = 0x00000004
	AUDIO_ALLOW_SAMPLES_CHANGE   = 0x00000008
	AUDIO_ALLOW_ANY_CHANGE       = AUDIO_ALLOW_FREQUENCY_CHANGE | AUDIO_ALLOW_FORMAT_CHANGE | AUDIO_ALLOW_CHANNELS_CHANGE | AUDIO_ALLOW_SAMPLES_CHANGE
)

/* Manually added functions */

// LoadWAV
/*
SDL_AudioSpec* SDL_LoadWAV(const char*    file,
                           SDL_AudioSpec* spec,
                           Uint8**        audio_buf,
                           Uint32*        audio_len)
*/
// #define SDL_LoadWAV(file, spec, audio_buf, audio_len) SDL_LoadWAV_RW(SDL_RWFromFile(file, "rb"),1, spec,audio_buf,audio_len)
// This is a C macro. It is not directly accessible via the DLL file.
// https://wiki.libsdl.org/SDL2/SDL_LoadWAV
func LoadWAV(file string, spec *m.AudioSpec, audioBuf **m.Uint8, audioLen *m.Uint32) *m.AudioSpec {
	ops := RWFromFile(file, "rb")
	if ops == nil {
		panic(ops)
	}

	ret, _, _ := syscall.SyscallN(fnLoadWAV_RW, uintptr(unsafe.Pointer(ops)), 0, uintptr(unsafe.Pointer(spec)), uintptr(unsafe.Pointer(audioBuf)), uintptr(unsafe.Pointer(audioLen)))
	return (*m.AudioSpec)(unsafe.Pointer(ret))
}

/* Automatically added functions */

//extern DECLSPEC int SDLCALL SDL_GetNumAudioDrivers(void);
func GetNumAudioDrivers() m.Int {
	ret, _, _ := syscall.SyscallN(fnGetNumAudioDrivers)
	return (m.Int)(ret)
}

//extern DECLSPEC const char *SDLCALL SDL_GetAudioDriver(int index);
func GetAudioDriver(index m.Int) string {
	ret, _, _ := syscall.SyscallN(fnGetAudioDriver, uintptr(index))
	return windows.BytePtrToString((*byte)(unsafe.Pointer(ret)))
}

//extern DECLSPEC int SDLCALL SDL_AudioInit(const char *driver_name);
func AudioInit(driver_name string) m.Int {
	ret, _, _ := syscall.SyscallN(fnAudioInit, uintptr(unsafe.Pointer(BytePtrFromStringP(driver_name))))
	return (m.Int)(ret)
}

//extern DECLSPEC void SDLCALL SDL_AudioQuit(void);
func AudioQuit() {
	_, _, _ = syscall.SyscallN(fnAudioQuit)
}

//extern DECLSPEC const char *SDLCALL SDL_GetCurrentAudioDriver(void);
func GetCurrentAudioDriver() string {
	ret, _, _ := syscall.SyscallN(fnGetCurrentAudioDriver)
	return windows.BytePtrToString((*byte)(unsafe.Pointer(ret)))
}

//extern DECLSPEC int SDLCALL SDL_OpenAudio(SDL_AudioSpec * desired, SDL_AudioSpec * obtained);
func OpenAudio(desired *m.AudioSpec, obtained *m.AudioSpec) m.Int {
	ret, _, _ := syscall.SyscallN(fnOpenAudio, uintptr(unsafe.Pointer(desired)), uintptr(unsafe.Pointer(obtained)))
	return (m.Int)(ret)
}

//extern DECLSPEC int SDLCALL SDL_GetNumAudioDevices(int iscapture);
func GetNumAudioDevices(iscapture m.Int) m.Int {
	ret, _, _ := syscall.SyscallN(fnGetNumAudioDevices, uintptr(iscapture))
	return (m.Int)(ret)
}

//extern DECLSPEC const char *SDLCALL SDL_GetAudioDeviceName(int index, int iscapture);
func GetAudioDeviceName(index m.Int, iscapture m.Int) string {
	ret, _, _ := syscall.SyscallN(fnGetAudioDeviceName, uintptr(index), uintptr(iscapture))
	return windows.BytePtrToString((*byte)(unsafe.Pointer(ret)))
}

//extern DECLSPEC int SDLCALL SDL_GetAudioDeviceSpec(int index, int iscapture, SDL_AudioSpec *spec);
func GetAudioDeviceSpec(index m.Int, iscapture m.Int, spec *m.AudioSpec) m.Int {
	ret, _, _ := syscall.SyscallN(fnGetAudioDeviceSpec, uintptr(index), uintptr(iscapture), uintptr(unsafe.Pointer(spec)))
	return (m.Int)(ret)
}

//extern DECLSPEC int SDLCALL SDL_GetDefaultAudioInfo(char **name, SDL_AudioSpec *spec, int iscapture);
func GetDefaultAudioInfo(name **m.Char, spec *m.AudioSpec, iscapture m.Int) m.Int {
	ret, _, _ := syscall.SyscallN(fnGetDefaultAudioInfo, uintptr(unsafe.Pointer(name)), uintptr(unsafe.Pointer(spec)), uintptr(iscapture))
	return (m.Int)(ret)
}

//extern DECLSPEC SDL_AudioDeviceID SDLCALL SDL_OpenAudioDevice(const char *device, int iscapture, const SDL_AudioSpec *desired, SDL_AudioSpec *obtained, int allowed_changes);
func OpenAudioDevice(device string, iscapture m.Int, desired *m.AudioSpec, obtained *m.AudioSpec, allowed_changes m.Int) m.AudioDeviceID {
	ret, _, _ := syscall.SyscallN(fnOpenAudioDevice, uintptr(unsafe.Pointer(BytePtrFromStringP(device))), uintptr(iscapture), uintptr(unsafe.Pointer(desired)), uintptr(unsafe.Pointer(obtained)), uintptr(allowed_changes))
	return (m.AudioDeviceID)(ret)
}

//extern DECLSPEC SDL_AudioStatus SDLCALL SDL_GetAudioStatus(void);
func GetAudioStatus() m.AudioStatus {
	ret, _, _ := syscall.SyscallN(fnGetAudioStatus)
	return (m.AudioStatus)(ret)
}

//extern DECLSPEC SDL_AudioStatus SDLCALL SDL_GetAudioDeviceStatus(SDL_AudioDeviceID dev);
func GetAudioDeviceStatus(dev m.AudioDeviceID) m.AudioStatus {
	ret, _, _ := syscall.SyscallN(fnGetAudioDeviceStatus, uintptr(dev))
	return (m.AudioStatus)(ret)
}

//extern DECLSPEC void SDLCALL SDL_PauseAudio(int pause_on);
func PauseAudio(pause_on m.Int) {
	_, _, _ = syscall.SyscallN(fnPauseAudio, uintptr(pause_on))
}

//extern DECLSPEC void SDLCALL SDL_PauseAudioDevice(SDL_AudioDeviceID dev, int pause_on);
func PauseAudioDevice(dev m.AudioDeviceID, pause_on m.Int) {
	_, _, _ = syscall.SyscallN(fnPauseAudioDevice, uintptr(dev), uintptr(pause_on))
}

//extern DECLSPEC SDL_AudioSpec *SDLCALL SDL_LoadWAV_RW(SDL_RWops * src, int freesrc, SDL_AudioSpec * spec, Uint8 ** audio_buf, Uint32 * audio_len);
func LoadWAV_RW(src *m.RWops, freesrc m.Int, spec *m.AudioSpec, audio_buf **m.Uint8, audio_len *m.Uint32) *m.AudioSpec {
	ret, _, _ := syscall.SyscallN(fnLoadWAV_RW, uintptr(unsafe.Pointer(src)), uintptr(freesrc), uintptr(unsafe.Pointer(spec)), uintptr(unsafe.Pointer(audio_buf)), uintptr(unsafe.Pointer(audio_len)))
	return (*m.AudioSpec)(unsafe.Pointer(ret))
}

//extern DECLSPEC void SDLCALL SDL_FreeWAV(Uint8 * audio_buf);
func FreeWAV(audio_buf *m.Uint8) {
	_, _, _ = syscall.SyscallN(fnFreeWAV, uintptr(unsafe.Pointer(audio_buf)))
}

//extern DECLSPEC int SDLCALL SDL_BuildAudioCVT(SDL_AudioCVT * cvt, SDL_AudioFormat src_format, Uint8 src_channels, int src_rate, SDL_AudioFormat dst_format, Uint8 dst_channels, int dst_rate);
func BuildAudioCVT(cvt *m.AudioCVT, src_format m.AudioFormat, src_channels m.Uint8, src_rate m.Int, dst_format m.AudioFormat, dst_channels m.Uint8, dst_rate m.Int) m.Int {
	ret, _, _ := syscall.SyscallN(fnBuildAudioCVT, uintptr(unsafe.Pointer(cvt)), uintptr(src_format), uintptr(src_channels), uintptr(src_rate), uintptr(dst_format), uintptr(dst_channels), uintptr(dst_rate))
	return (m.Int)(ret)
}

//extern DECLSPEC int SDLCALL SDL_ConvertAudio(SDL_AudioCVT * cvt);
func ConvertAudio(cvt *m.AudioCVT) m.Int {
	ret, _, _ := syscall.SyscallN(fnConvertAudio, uintptr(unsafe.Pointer(cvt)))
	return (m.Int)(ret)
}

//extern DECLSPEC SDL_AudioStream * SDLCALL SDL_NewAudioStream(const SDL_AudioFormat src_format, const Uint8 src_channels, const int src_rate, const SDL_AudioFormat dst_format, const Uint8 dst_channels, const int dst_rate);
func NewAudioStream(src_format m.AudioFormat, src_channels m.Uint8, src_rate m.Int, dst_format m.AudioFormat, dst_channels m.Uint8, dst_rate m.Int) *m.AudioStream {
	ret, _, _ := syscall.SyscallN(fnNewAudioStream, uintptr(src_format), uintptr(src_channels), uintptr(src_rate), uintptr(dst_format), uintptr(dst_channels), uintptr(dst_rate))
	return (*m.AudioStream)(unsafe.Pointer(ret))
}

//extern DECLSPEC int SDLCALL SDL_AudioStreamPut(SDL_AudioStream *stream, const void *buf, int len);
func AudioStreamPut(stream *m.AudioStream, buf *m.Void, len m.Int) m.Int {
	ret, _, _ := syscall.SyscallN(fnAudioStreamPut, uintptr(unsafe.Pointer(stream)), uintptr(unsafe.Pointer(buf)), uintptr(len))
	return (m.Int)(ret)
}

//extern DECLSPEC int SDLCALL SDL_AudioStreamGet(SDL_AudioStream *stream, void *buf, int len);
func AudioStreamGet(stream *m.AudioStream, buf *m.Void, len m.Int) m.Int {
	ret, _, _ := syscall.SyscallN(fnAudioStreamGet, uintptr(unsafe.Pointer(stream)), uintptr(unsafe.Pointer(buf)), uintptr(len))
	return (m.Int)(ret)
}

//extern DECLSPEC int SDLCALL SDL_AudioStreamAvailable(SDL_AudioStream *stream);
func AudioStreamAvailable(stream *m.AudioStream) m.Int {
	ret, _, _ := syscall.SyscallN(fnAudioStreamAvailable, uintptr(unsafe.Pointer(stream)))
	return (m.Int)(ret)
}

//extern DECLSPEC int SDLCALL SDL_AudioStreamFlush(SDL_AudioStream *stream);
func AudioStreamFlush(stream *m.AudioStream) m.Int {
	ret, _, _ := syscall.SyscallN(fnAudioStreamFlush, uintptr(unsafe.Pointer(stream)))
	return (m.Int)(ret)
}

//extern DECLSPEC void SDLCALL SDL_AudioStreamClear(SDL_AudioStream *stream);
func AudioStreamClear(stream *m.AudioStream) {
	_, _, _ = syscall.SyscallN(fnAudioStreamClear, uintptr(unsafe.Pointer(stream)))
}

//extern DECLSPEC void SDLCALL SDL_FreeAudioStream(SDL_AudioStream *stream);
func FreeAudioStream(stream *m.AudioStream) {
	_, _, _ = syscall.SyscallN(fnFreeAudioStream, uintptr(unsafe.Pointer(stream)))
}

//extern DECLSPEC void SDLCALL SDL_MixAudio(Uint8 * dst, const Uint8 * src, Uint32 len, int volume);
func MixAudio(dst *m.Uint8, src *m.Uint8, len m.Uint32, volume m.Int) {
	_, _, _ = syscall.SyscallN(fnMixAudio, uintptr(unsafe.Pointer(dst)), uintptr(unsafe.Pointer(src)), uintptr(len), uintptr(volume))
}

//extern DECLSPEC void SDLCALL SDL_MixAudioFormat(Uint8 * dst, const Uint8 * src, SDL_AudioFormat format, Uint32 len, int volume);
func MixAudioFormat(dst *m.Uint8, src *m.Uint8, format m.AudioFormat, len m.Uint32, volume m.Int) {
	_, _, _ = syscall.SyscallN(fnMixAudioFormat, uintptr(unsafe.Pointer(dst)), uintptr(unsafe.Pointer(src)), uintptr(format), uintptr(len), uintptr(volume))
}

//extern DECLSPEC int SDLCALL SDL_QueueAudio(SDL_AudioDeviceID dev, const void *data, Uint32 len);
func QueueAudio(dev m.AudioDeviceID, data *m.Void, len m.Uint32) m.Int {
	ret, _, _ := syscall.SyscallN(fnQueueAudio, uintptr(dev), uintptr(unsafe.Pointer(data)), uintptr(len))
	return (m.Int)(ret)
}

//extern DECLSPEC Uint32 SDLCALL SDL_DequeueAudio(SDL_AudioDeviceID dev, void *data, Uint32 len);
func DequeueAudio(dev m.AudioDeviceID, data *m.Void, len m.Uint32) m.Uint32 {
	ret, _, _ := syscall.SyscallN(fnDequeueAudio, uintptr(dev), uintptr(unsafe.Pointer(data)), uintptr(len))
	return (m.Uint32)(ret)
}

//extern DECLSPEC Uint32 SDLCALL SDL_GetQueuedAudioSize(SDL_AudioDeviceID dev);
func GetQueuedAudioSize(dev m.AudioDeviceID) m.Uint32 {
	ret, _, _ := syscall.SyscallN(fnGetQueuedAudioSize, uintptr(dev))
	return (m.Uint32)(ret)
}

//extern DECLSPEC void SDLCALL SDL_ClearQueuedAudio(SDL_AudioDeviceID dev);
func ClearQueuedAudio(dev m.AudioDeviceID) {
	_, _, _ = syscall.SyscallN(fnClearQueuedAudio, uintptr(dev))
}

//extern DECLSPEC void SDLCALL SDL_LockAudio(void);
func LockAudio() {
	_, _, _ = syscall.SyscallN(fnLockAudio)
}

//extern DECLSPEC void SDLCALL SDL_LockAudioDevice(SDL_AudioDeviceID dev);
func LockAudioDevice(dev m.AudioDeviceID) {
	_, _, _ = syscall.SyscallN(fnLockAudioDevice, uintptr(dev))
}

//extern DECLSPEC void SDLCALL SDL_UnlockAudio(void);
func UnlockAudio() {
	_, _, _ = syscall.SyscallN(fnUnlockAudio)
}

//extern DECLSPEC void SDLCALL SDL_UnlockAudioDevice(SDL_AudioDeviceID dev);
func UnlockAudioDevice(dev m.AudioDeviceID) {
	_, _, _ = syscall.SyscallN(fnUnlockAudioDevice, uintptr(dev))
}

//extern DECLSPEC void SDLCALL SDL_CloseAudio(void);
func CloseAudio() {
	_, _, _ = syscall.SyscallN(fnCloseAudio)
}

//extern DECLSPEC void SDLCALL SDL_CloseAudioDevice(SDL_AudioDeviceID dev);
func CloseAudioDevice(dev m.AudioDeviceID) {
	_, _, _ = syscall.SyscallN(fnCloseAudioDevice, uintptr(dev))
}
