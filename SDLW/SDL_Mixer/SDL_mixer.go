package sdlm

// SDL_mixer.h.

import (
	"fmt"
	"syscall"
	"unsafe"

	m "github.com/vault-thirteen/SDLW/SDL/model"
	mm "github.com/vault-thirteen/SDLW/SDL_Mixer/model"
	"golang.org/x/sys/windows"
)

const (
	CHANNELS          = 8
	DEFAULT_FREQUENCY = 44100
	DEFAULT_CHANNELS  = 2
)

/* Experimental code */

// Test
// This tests shows that Golang is a 'komplete Scheiße'. What does that mean ?
// Here we are using system calls without any type casts. We use pointers
// returned by a DLL library directly, i.e. we use pointers returned by C code,
// and we do not type cast them into Go pointers except one thing. We are using
// the Go's `uintptr` type, because system calls use them, and we have no other
// choice other than this. Well, this very simple experiment shows that a
// pointer to Mix_Music type is not NULL according to the non-zero printed
// value. What happens next ? This pointer is given to another DLL function.
// In C code this function returns a positive value of 503, but in Golang it
// returns zero (0). What does it mean ? Golang is a 'komplete Scheiße'.
// I and will not create an issue on Golang's GitHub page because they will say
// as they usually do that they are great and we are nothing. This small
// example shows how "great" is Golang in reality.
// Golang ist komplete Scheiße. Das stimmt.
func Test(filePath string) {
	music, _, _ := syscall.SyscallN(fnLoadMUS, uintptr(unsafe.Pointer(BytePtrFromStringP(filePath))))
	fmt.Println("music:", music)
	duration, _, _ := syscall.SyscallN(fnMusicDuration, music)
	fmt.Println("duration:", duration)
}

/* Manually added functions */

// extern DECLSPEC void SDLCALL Mix_SetPostMix(void (SDLCALL *mix_func)(void *udata, Uint8 *stream, int len), void *arg);
func SetPostMix(mix_func mm.MixFunc, arg *m.Void) {
	_, _, _ = syscall.SyscallN(fnSetPostMix, uintptr(unsafe.Pointer(mix_func)), uintptr(unsafe.Pointer(arg)))
}

// extern DECLSPEC void SDLCALL Mix_HookMusic(void (SDLCALL *mix_func)(void *udata, Uint8 *stream, int len), void *arg);
func HookMusic(mix_func mm.MixFunc, arg *m.Void) {
	_, _, _ = syscall.SyscallN(fnHookMusic, uintptr(unsafe.Pointer(mix_func)), uintptr(unsafe.Pointer(arg)))
}

// extern DECLSPEC void SDLCALL Mix_HookMusicFinished(void (SDLCALL *music_finished)(void));
func HookMusicFinished(music_finished mm.MusicFinished) {
	_, _, _ = syscall.SyscallN(fnHookMusicFinished, uintptr(unsafe.Pointer(music_finished)))
}

// extern DECLSPEC void SDLCALL Mix_ChannelFinished(void (SDLCALL *channel_finished)(int channel));
func ChannelFinished(channel_finished mm.ChannelFinished) {
	_, _, _ = syscall.SyscallN(fnChannelFinished, uintptr(unsafe.Pointer(channel_finished)))
}

// extern DECLSPEC int SDLCALL Mix_EachSoundFont(int (SDLCALL *function)(const char*, void*), void *data);
func EachSoundFont(function mm.Function, data *m.Void) m.Int {
	ret, _, _ := syscall.SyscallN(fnEachSoundFont, uintptr(unsafe.Pointer(function)), uintptr(unsafe.Pointer(data)))
	return m.Int(ret)
}

/* Automatically added functions */

//extern DECLSPEC const SDL_version * SDLCALL Mix_Linked_Version(void);
func Linked_Version() *m.Version {
	ret, _, _ := syscall.SyscallN(fnLinked_Version)
	return (*m.Version)(unsafe.Pointer(ret))
}

//extern DECLSPEC int SDLCALL Mix_Init(int flags);
func Init(flags m.Int) m.Int {
	ret, _, _ := syscall.SyscallN(fnInit, uintptr(flags))
	return (m.Int)(ret)
}

//extern DECLSPEC void SDLCALL Mix_Quit(void);
func Quit() {
	_, _, _ = syscall.SyscallN(fnQuit)
}

//extern DECLSPEC int SDLCALL Mix_OpenAudio(int frequency, Uint16 format, int channels, int chunksize);
func OpenAudio(frequency m.Int, format m.Uint16, channels m.Int, chunksize m.Int) m.Int {
	ret, _, _ := syscall.SyscallN(fnOpenAudio, uintptr(frequency), uintptr(format), uintptr(channels), uintptr(chunksize))
	return (m.Int)(ret)
}

//extern DECLSPEC int SDLCALL Mix_OpenAudioDevice(int frequency, Uint16 format, int channels, int chunksize, const char* device, int allowed_changes);
func OpenAudioDevice(frequency m.Int, format m.Uint16, channels m.Int, chunksize m.Int, device string, allowed_changes m.Int) m.Int {
	ret, _, _ := syscall.SyscallN(fnOpenAudioDevice, uintptr(frequency), uintptr(format), uintptr(channels), uintptr(chunksize), uintptr(unsafe.Pointer(BytePtrFromStringP(device))), uintptr(allowed_changes))
	return (m.Int)(ret)
}

//extern DECLSPEC int SDLCALL Mix_QuerySpec(int *frequency, Uint16 *format, int *channels);
func QuerySpec(frequency *m.Int, format *m.Uint16, channels *m.Int) m.Int {
	ret, _, _ := syscall.SyscallN(fnQuerySpec, uintptr(unsafe.Pointer(frequency)), uintptr(unsafe.Pointer(format)), uintptr(unsafe.Pointer(channels)))
	return (m.Int)(ret)
}

//extern DECLSPEC int SDLCALL Mix_AllocateChannels(int numchans);
func AllocateChannels(numchans m.Int) m.Int {
	ret, _, _ := syscall.SyscallN(fnAllocateChannels, uintptr(numchans))
	return (m.Int)(ret)
}

//extern DECLSPEC Mix_Chunk * SDLCALL Mix_LoadWAV_RW(SDL_RWops *src, int freesrc);
func LoadWAV_RW(src *m.RWops, freesrc m.Int) *mm.Chunk {
	ret, _, _ := syscall.SyscallN(fnLoadWAV_RW, uintptr(unsafe.Pointer(src)), uintptr(freesrc))
	return (*mm.Chunk)(unsafe.Pointer(ret))
}

//extern DECLSPEC Mix_Chunk * SDLCALL Mix_LoadWAV(const char *file);
func LoadWAV(file string) *mm.Chunk {
	ret, _, _ := syscall.SyscallN(fnLoadWAV, uintptr(unsafe.Pointer(BytePtrFromStringP(file))))
	return (*mm.Chunk)(unsafe.Pointer(ret))
}

//extern DECLSPEC Mix_Music * SDLCALL Mix_LoadMUS(const char *file);
func LoadMUS(file string) *mm.Music {
	ret, _, _ := syscall.SyscallN(fnLoadMUS, uintptr(unsafe.Pointer(BytePtrFromStringP(file))))
	return (*mm.Music)(unsafe.Pointer(ret))
}

//extern DECLSPEC Mix_Music * SDLCALL Mix_LoadMUS_RW(SDL_RWops *src, int freesrc);
func LoadMUS_RW(src *m.RWops, freesrc m.Int) *mm.Music {
	ret, _, _ := syscall.SyscallN(fnLoadMUS_RW, uintptr(unsafe.Pointer(src)), uintptr(freesrc))
	return (*mm.Music)(unsafe.Pointer(ret))
}

//extern DECLSPEC Mix_Music * SDLCALL Mix_LoadMUSType_RW(SDL_RWops *src, Mix_MusicType type, int freesrc);
func LoadMUSType_RW(src *m.RWops, musicType mm.MusicType, freesrc m.Int) *mm.Music {
	ret, _, _ := syscall.SyscallN(fnLoadMUSType_RW, uintptr(unsafe.Pointer(src)), uintptr(musicType), uintptr(freesrc))
	return (*mm.Music)(unsafe.Pointer(ret))
}

//extern DECLSPEC Mix_Chunk * SDLCALL Mix_QuickLoad_WAV(Uint8 *mem);
func QuickLoad_WAV(mem *m.Uint8) *mm.Chunk {
	ret, _, _ := syscall.SyscallN(fnQuickLoad_WAV, uintptr(unsafe.Pointer(mem)))
	return (*mm.Chunk)(unsafe.Pointer(ret))
}

//extern DECLSPEC Mix_Chunk * SDLCALL Mix_QuickLoad_RAW(Uint8 *mem, Uint32 len);
func QuickLoad_RAW(mem *m.Uint8, len m.Uint32) *mm.Chunk {
	ret, _, _ := syscall.SyscallN(fnQuickLoad_RAW, uintptr(unsafe.Pointer(mem)), uintptr(len))
	return (*mm.Chunk)(unsafe.Pointer(ret))
}

//extern DECLSPEC void SDLCALL Mix_FreeChunk(Mix_Chunk *chunk);
func FreeChunk(chunk *mm.Chunk) {
	_, _, _ = syscall.SyscallN(fnFreeChunk, uintptr(unsafe.Pointer(chunk)))
}

//extern DECLSPEC void SDLCALL Mix_FreeMusic(Mix_Music *music);
func FreeMusic(music *mm.Music) {
	_, _, _ = syscall.SyscallN(fnFreeMusic, uintptr(unsafe.Pointer(music)))
}

//extern DECLSPEC int SDLCALL Mix_GetNumChunkDecoders(void);
func GetNumChunkDecoders() m.Int {
	ret, _, _ := syscall.SyscallN(fnGetNumChunkDecoders)
	return (m.Int)(ret)
}

//extern DECLSPEC const char * SDLCALL Mix_GetChunkDecoder(int index);
func GetChunkDecoder(index m.Int) string {
	ret, _, _ := syscall.SyscallN(fnGetChunkDecoder, uintptr(index))
	return windows.BytePtrToString((*byte)(unsafe.Pointer(ret)))
}

//extern DECLSPEC SDL_bool SDLCALL Mix_HasChunkDecoder(const char *name);
func HasChunkDecoder(name string) m.Bool {
	ret, _, _ := syscall.SyscallN(fnHasChunkDecoder, uintptr(unsafe.Pointer(BytePtrFromStringP(name))))
	return (m.Bool)(ret)
}

//extern DECLSPEC int SDLCALL Mix_GetNumMusicDecoders(void);
func GetNumMusicDecoders() m.Int {
	ret, _, _ := syscall.SyscallN(fnGetNumMusicDecoders)
	return (m.Int)(ret)
}

//extern DECLSPEC const char * SDLCALL Mix_GetMusicDecoder(int index);
func GetMusicDecoder(index m.Int) string {
	ret, _, _ := syscall.SyscallN(fnGetMusicDecoder, uintptr(index))
	return windows.BytePtrToString((*byte)(unsafe.Pointer(ret)))
}

//extern DECLSPEC SDL_bool SDLCALL Mix_HasMusicDecoder(const char *name);
func HasMusicDecoder(name string) m.Bool {
	ret, _, _ := syscall.SyscallN(fnHasMusicDecoder, uintptr(unsafe.Pointer(BytePtrFromStringP(name))))
	return (m.Bool)(ret)
}

//extern DECLSPEC Mix_MusicType SDLCALL Mix_GetMusicType(const Mix_Music *music);
func GetMusicType(music *mm.Music) mm.MusicType {
	ret, _, _ := syscall.SyscallN(fnGetMusicType, uintptr(unsafe.Pointer(music)))
	return (mm.MusicType)(ret)
}

//extern DECLSPEC const char *SDLCALL Mix_GetMusicTitle(const Mix_Music *music);
func GetMusicTitle(music *mm.Music) string {
	ret, _, _ := syscall.SyscallN(fnGetMusicTitle, uintptr(unsafe.Pointer(music)))
	return windows.BytePtrToString((*byte)(unsafe.Pointer(ret)))
}

//extern DECLSPEC const char *SDLCALL Mix_GetMusicTitleTag(const Mix_Music *music);
func GetMusicTitleTag(music *mm.Music) string {
	ret, _, _ := syscall.SyscallN(fnGetMusicTitleTag, uintptr(unsafe.Pointer(music)))
	return windows.BytePtrToString((*byte)(unsafe.Pointer(ret)))
}

//extern DECLSPEC const char *SDLCALL Mix_GetMusicArtistTag(const Mix_Music *music);
func GetMusicArtistTag(music *mm.Music) string {
	ret, _, _ := syscall.SyscallN(fnGetMusicArtistTag, uintptr(unsafe.Pointer(music)))
	return windows.BytePtrToString((*byte)(unsafe.Pointer(ret)))
}

//extern DECLSPEC const char *SDLCALL Mix_GetMusicAlbumTag(const Mix_Music *music);
func GetMusicAlbumTag(music *mm.Music) string {
	ret, _, _ := syscall.SyscallN(fnGetMusicAlbumTag, uintptr(unsafe.Pointer(music)))
	return windows.BytePtrToString((*byte)(unsafe.Pointer(ret)))
}

//extern DECLSPEC const char *SDLCALL Mix_GetMusicCopyrightTag(const Mix_Music *music);
func GetMusicCopyrightTag(music *mm.Music) string {
	ret, _, _ := syscall.SyscallN(fnGetMusicCopyrightTag, uintptr(unsafe.Pointer(music)))
	return windows.BytePtrToString((*byte)(unsafe.Pointer(ret)))
}

//extern DECLSPEC void * SDLCALL Mix_GetMusicHookData(void);
func GetMusicHookData() *m.Void {
	ret, _, _ := syscall.SyscallN(fnGetMusicHookData)
	return (*m.Void)(unsafe.Pointer(ret))
}

//extern DECLSPEC int SDLCALL Mix_RegisterEffect(int chan, Mix_EffectFunc_t f, Mix_EffectDone_t d, void *arg);
func RegisterEffect(channel m.Int, f mm.EffectFunc, d mm.EffectDone, arg *m.Void) m.Int {
	ret, _, _ := syscall.SyscallN(fnRegisterEffect, uintptr(channel), uintptr(unsafe.Pointer(f)), uintptr(unsafe.Pointer(d)), uintptr(unsafe.Pointer(arg)))
	return (m.Int)(ret)
}

//extern DECLSPEC int SDLCALL Mix_UnregisterEffect(int channel, Mix_EffectFunc_t f);
func UnregisterEffect(channel m.Int, f mm.EffectFunc) m.Int {
	ret, _, _ := syscall.SyscallN(fnUnregisterEffect, uintptr(channel), uintptr(unsafe.Pointer(f)))
	return (m.Int)(ret)
}

//extern DECLSPEC int SDLCALL Mix_UnregisterAllEffects(int channel);
func UnregisterAllEffects(channel m.Int) m.Int {
	ret, _, _ := syscall.SyscallN(fnUnregisterAllEffects, uintptr(channel))
	return (m.Int)(ret)
}

//extern DECLSPEC int SDLCALL Mix_SetPanning(int channel, Uint8 left, Uint8 right);
func SetPanning(channel m.Int, left m.Uint8, right m.Uint8) m.Int {
	ret, _, _ := syscall.SyscallN(fnSetPanning, uintptr(channel), uintptr(left), uintptr(right))
	return (m.Int)(ret)
}

//extern DECLSPEC int SDLCALL Mix_SetPosition(int channel, Sint16 angle, Uint8 distance);
func SetPosition(channel m.Int, angle m.Int16, distance m.Uint8) m.Int {
	ret, _, _ := syscall.SyscallN(fnSetPosition, uintptr(channel), uintptr(angle), uintptr(distance))
	return (m.Int)(ret)
}

//extern DECLSPEC int SDLCALL Mix_SetDistance(int channel, Uint8 distance);
func SetDistance(channel m.Int, distance m.Uint8) m.Int {
	ret, _, _ := syscall.SyscallN(fnSetDistance, uintptr(channel), uintptr(distance))
	return (m.Int)(ret)
}

//extern DECLSPEC int SDLCALL Mix_SetReverseStereo(int channel, int flip);
func SetReverseStereo(channel m.Int, flip m.Int) m.Int {
	ret, _, _ := syscall.SyscallN(fnSetReverseStereo, uintptr(channel), uintptr(flip))
	return (m.Int)(ret)
}

//extern DECLSPEC int SDLCALL Mix_ReserveChannels(int num);
func ReserveChannels(num m.Int) m.Int {
	ret, _, _ := syscall.SyscallN(fnReserveChannels, uintptr(num))
	return (m.Int)(ret)
}

//extern DECLSPEC int SDLCALL Mix_GroupChannel(int which, int tag);
func GroupChannel(which m.Int, tag m.Int) m.Int {
	ret, _, _ := syscall.SyscallN(fnGroupChannel, uintptr(which), uintptr(tag))
	return (m.Int)(ret)
}

//extern DECLSPEC int SDLCALL Mix_GroupChannels(int from, int to, int tag);
func GroupChannels(from m.Int, to m.Int, tag m.Int) m.Int {
	ret, _, _ := syscall.SyscallN(fnGroupChannels, uintptr(from), uintptr(to), uintptr(tag))
	return (m.Int)(ret)
}

//extern DECLSPEC int SDLCALL Mix_GroupAvailable(int tag);
func GroupAvailable(tag m.Int) m.Int {
	ret, _, _ := syscall.SyscallN(fnGroupAvailable, uintptr(tag))
	return (m.Int)(ret)
}

//extern DECLSPEC int SDLCALL Mix_GroupCount(int tag);
func GroupCount(tag m.Int) m.Int {
	ret, _, _ := syscall.SyscallN(fnGroupCount, uintptr(tag))
	return (m.Int)(ret)
}

//extern DECLSPEC int SDLCALL Mix_GroupOldest(int tag);
func GroupOldest(tag m.Int) m.Int {
	ret, _, _ := syscall.SyscallN(fnGroupOldest, uintptr(tag))
	return (m.Int)(ret)
}

//extern DECLSPEC int SDLCALL Mix_GroupNewer(int tag);
func GroupNewer(tag m.Int) m.Int {
	ret, _, _ := syscall.SyscallN(fnGroupNewer, uintptr(tag))
	return (m.Int)(ret)
}

//extern DECLSPEC int SDLCALL Mix_PlayChannel(int channel, Mix_Chunk *chunk, int loops);
func PlayChannel(channel m.Int, chunk *mm.Chunk, loops m.Int) m.Int {
	ret, _, _ := syscall.SyscallN(fnPlayChannel, uintptr(channel), uintptr(unsafe.Pointer(chunk)), uintptr(loops))
	return (m.Int)(ret)
}

//extern DECLSPEC int SDLCALL Mix_PlayChannelTimed(int channel, Mix_Chunk *chunk, int loops, int ticks);
func PlayChannelTimed(channel m.Int, chunk *mm.Chunk, loops m.Int, ticks m.Int) m.Int {
	ret, _, _ := syscall.SyscallN(fnPlayChannelTimed, uintptr(channel), uintptr(unsafe.Pointer(chunk)), uintptr(loops), uintptr(ticks))
	return (m.Int)(ret)
}

//extern DECLSPEC int SDLCALL Mix_PlayMusic(Mix_Music *music, int loops);
func PlayMusic(music *mm.Music, loops m.Int) m.Int {
	ret, _, _ := syscall.SyscallN(fnPlayMusic, uintptr(unsafe.Pointer(music)), uintptr(loops))
	return (m.Int)(ret)
}

//extern DECLSPEC int SDLCALL Mix_FadeInMusic(Mix_Music *music, int loops, int ms);
func FadeInMusic(music *mm.Music, loops m.Int, ms m.Int) m.Int {
	ret, _, _ := syscall.SyscallN(fnFadeInMusic, uintptr(unsafe.Pointer(music)), uintptr(loops), uintptr(ms))
	return (m.Int)(ret)
}

//extern DECLSPEC int SDLCALL Mix_FadeInMusicPos(Mix_Music *music, int loops, int ms, double position);
func FadeInMusicPos(music *mm.Music, loops m.Int, ms m.Int, position m.Double) m.Int {
	ret, _, _ := syscall.SyscallN(fnFadeInMusicPos, uintptr(unsafe.Pointer(music)), uintptr(loops), uintptr(ms), uintptr(position))
	return (m.Int)(ret)
}

//extern DECLSPEC int SDLCALL Mix_FadeInChannel(int channel, Mix_Chunk *chunk, int loops, int ms);
func FadeInChannel(channel m.Int, chunk *mm.Chunk, loops m.Int, ms m.Int) m.Int {
	ret, _, _ := syscall.SyscallN(fnFadeInChannel, uintptr(channel), uintptr(unsafe.Pointer(chunk)), uintptr(loops), uintptr(ms))
	return (m.Int)(ret)
}

//extern DECLSPEC int SDLCALL Mix_FadeInChannelTimed(int channel, Mix_Chunk *chunk, int loops, int ms, int ticks);
func FadeInChannelTimed(channel m.Int, chunk *mm.Chunk, loops m.Int, ms m.Int, ticks m.Int) m.Int {
	ret, _, _ := syscall.SyscallN(fnFadeInChannelTimed, uintptr(channel), uintptr(unsafe.Pointer(chunk)), uintptr(loops), uintptr(ms), uintptr(ticks))
	return (m.Int)(ret)
}

//extern DECLSPEC int SDLCALL Mix_Volume(int channel, int volume);
func Volume(channel m.Int, volume m.Int) m.Int {
	ret, _, _ := syscall.SyscallN(fnVolume, uintptr(channel), uintptr(volume))
	return (m.Int)(ret)
}

//extern DECLSPEC int SDLCALL Mix_VolumeChunk(Mix_Chunk *chunk, int volume);
func VolumeChunk(chunk *mm.Chunk, volume m.Int) m.Int {
	ret, _, _ := syscall.SyscallN(fnVolumeChunk, uintptr(unsafe.Pointer(chunk)), uintptr(volume))
	return (m.Int)(ret)
}

//extern DECLSPEC int SDLCALL Mix_VolumeMusic(int volume);
func VolumeMusic(volume m.Int) m.Int {
	ret, _, _ := syscall.SyscallN(fnVolumeMusic, uintptr(volume))
	return (m.Int)(ret)
}

//extern DECLSPEC int SDLCALL Mix_GetMusicVolume(Mix_Music *music);
func GetMusicVolume(music *mm.Music) m.Int {
	ret, _, _ := syscall.SyscallN(fnGetMusicVolume, uintptr(unsafe.Pointer(music)))
	return (m.Int)(ret)
}

//extern DECLSPEC int SDLCALL Mix_MasterVolume(int volume);
func MasterVolume(volume m.Int) m.Int {
	ret, _, _ := syscall.SyscallN(fnMasterVolume, uintptr(volume))
	return (m.Int)(ret)
}

//extern DECLSPEC int SDLCALL Mix_HaltChannel(int channel);
func HaltChannel(channel m.Int) m.Int {
	ret, _, _ := syscall.SyscallN(fnHaltChannel, uintptr(channel))
	return (m.Int)(ret)
}

//extern DECLSPEC int SDLCALL Mix_HaltGroup(int tag);
func HaltGroup(tag m.Int) m.Int {
	ret, _, _ := syscall.SyscallN(fnHaltGroup, uintptr(tag))
	return (m.Int)(ret)
}

//extern DECLSPEC int SDLCALL Mix_HaltMusic(void);
func HaltMusic() m.Int {
	ret, _, _ := syscall.SyscallN(fnHaltMusic)
	return (m.Int)(ret)
}

//extern DECLSPEC int SDLCALL Mix_ExpireChannel(int channel, int ticks);
func ExpireChannel(channel m.Int, ticks m.Int) m.Int {
	ret, _, _ := syscall.SyscallN(fnExpireChannel, uintptr(channel), uintptr(ticks))
	return (m.Int)(ret)
}

//extern DECLSPEC int SDLCALL Mix_FadeOutChannel(int which, int ms);
func FadeOutChannel(which m.Int, ms m.Int) m.Int {
	ret, _, _ := syscall.SyscallN(fnFadeOutChannel, uintptr(which), uintptr(ms))
	return (m.Int)(ret)
}

//extern DECLSPEC int SDLCALL Mix_FadeOutGroup(int tag, int ms);
func FadeOutGroup(tag m.Int, ms m.Int) m.Int {
	ret, _, _ := syscall.SyscallN(fnFadeOutGroup, uintptr(tag), uintptr(ms))
	return (m.Int)(ret)
}

//extern DECLSPEC int SDLCALL Mix_FadeOutMusic(int ms);
func FadeOutMusic(ms m.Int) m.Int {
	ret, _, _ := syscall.SyscallN(fnFadeOutMusic, uintptr(ms))
	return (m.Int)(ret)
}

//extern DECLSPEC Mix_Fading SDLCALL Mix_FadingMusic(void);
func FadingMusic() mm.Fading {
	ret, _, _ := syscall.SyscallN(fnFadingMusic)
	return (mm.Fading)(ret)
}

//extern DECLSPEC Mix_Fading SDLCALL Mix_FadingChannel(int which);
func FadingChannel(which m.Int) mm.Fading {
	ret, _, _ := syscall.SyscallN(fnFadingChannel, uintptr(which))
	return (mm.Fading)(ret)
}

//extern DECLSPEC void SDLCALL Mix_Pause(int channel);
func Pause(channel m.Int) {
	_, _, _ = syscall.SyscallN(fnPause, uintptr(channel))
}

//extern DECLSPEC void SDLCALL Mix_Resume(int channel);
func Resume(channel m.Int) {
	_, _, _ = syscall.SyscallN(fnResume, uintptr(channel))
}

//extern DECLSPEC int SDLCALL Mix_Paused(int channel);
func Paused(channel m.Int) m.Int {
	ret, _, _ := syscall.SyscallN(fnPaused, uintptr(channel))
	return (m.Int)(ret)
}

//extern DECLSPEC void SDLCALL Mix_PauseMusic(void);
func PauseMusic() {
	_, _, _ = syscall.SyscallN(fnPauseMusic)
}

//extern DECLSPEC void SDLCALL Mix_ResumeMusic(void);
func ResumeMusic() {
	_, _, _ = syscall.SyscallN(fnResumeMusic)
}

//extern DECLSPEC void SDLCALL Mix_RewindMusic(void);
func RewindMusic() {
	_, _, _ = syscall.SyscallN(fnRewindMusic)
}

//extern DECLSPEC int SDLCALL Mix_PausedMusic(void);
func PausedMusic() m.Int {
	ret, _, _ := syscall.SyscallN(fnPausedMusic)
	return (m.Int)(ret)
}

//extern DECLSPEC int SDLCALL Mix_ModMusicJumpToOrder(int order);
func ModMusicJumpToOrder(order m.Int) m.Int {
	ret, _, _ := syscall.SyscallN(fnModMusicJumpToOrder, uintptr(order))
	return (m.Int)(ret)
}

//extern DECLSPEC int SDLCALL Mix_SetMusicPosition(double position);
func SetMusicPosition(position m.Double) m.Int {
	ret, _, _ := syscall.SyscallN(fnSetMusicPosition, uintptr(position))
	return (m.Int)(ret)
}

//extern DECLSPEC double SDLCALL Mix_GetMusicPosition(Mix_Music *music);
func GetMusicPosition(music *mm.Music) m.Double {
	ret, _, _ := syscall.SyscallN(fnGetMusicPosition, uintptr(unsafe.Pointer(music)))
	return (m.Double)(ret)
}

//extern DECLSPEC double SDLCALL Mix_MusicDuration(Mix_Music *music);
func MusicDuration(music *mm.Music) m.Double {
	ret, _, _ := syscall.SyscallN(fnMusicDuration, uintptr(unsafe.Pointer(music)))
	return (m.Double)(ret)
}

//extern DECLSPEC double SDLCALL Mix_GetMusicLoopStartTime(Mix_Music *music);
func GetMusicLoopStartTime(music *mm.Music) m.Double {
	ret, _, _ := syscall.SyscallN(fnGetMusicLoopStartTime, uintptr(unsafe.Pointer(music)))
	return (m.Double)(ret)
}

//extern DECLSPEC double SDLCALL Mix_GetMusicLoopEndTime(Mix_Music *music);
func GetMusicLoopEndTime(music *mm.Music) m.Double {
	ret, _, _ := syscall.SyscallN(fnGetMusicLoopEndTime, uintptr(unsafe.Pointer(music)))
	return (m.Double)(ret)
}

//extern DECLSPEC double SDLCALL Mix_GetMusicLoopLengthTime(Mix_Music *music);
func GetMusicLoopLengthTime(music *mm.Music) m.Double {
	ret, _, _ := syscall.SyscallN(fnGetMusicLoopLengthTime, uintptr(unsafe.Pointer(music)))
	return (m.Double)(ret)
}

//extern DECLSPEC int SDLCALL Mix_Playing(int channel);
func Playing(channel m.Int) m.Int {
	ret, _, _ := syscall.SyscallN(fnPlaying, uintptr(channel))
	return (m.Int)(ret)
}

//extern DECLSPEC int SDLCALL Mix_PlayingMusic(void);
func PlayingMusic() m.Int {
	ret, _, _ := syscall.SyscallN(fnPlayingMusic)
	return (m.Int)(ret)
}

//extern DECLSPEC int SDLCALL Mix_SetMusicCMD(const char *command);
func SetMusicCMD(command string) m.Int {
	ret, _, _ := syscall.SyscallN(fnSetMusicCMD, uintptr(unsafe.Pointer(BytePtrFromStringP(command))))
	return (m.Int)(ret)
}

//extern DECLSPEC int SDLCALL Mix_SetSynchroValue(int value);
func SetSynchroValue(value m.Int) m.Int {
	ret, _, _ := syscall.SyscallN(fnSetSynchroValue, uintptr(value))
	return (m.Int)(ret)
}

//extern DECLSPEC int SDLCALL Mix_GetSynchroValue(void);
func GetSynchroValue() m.Int {
	ret, _, _ := syscall.SyscallN(fnGetSynchroValue)
	return (m.Int)(ret)
}

//extern DECLSPEC int SDLCALL Mix_SetSoundFonts(const char *paths);
func SetSoundFonts(paths string) m.Int {
	ret, _, _ := syscall.SyscallN(fnSetSoundFonts, uintptr(unsafe.Pointer(BytePtrFromStringP(paths))))
	return (m.Int)(ret)
}

//extern DECLSPEC const char* SDLCALL Mix_GetSoundFonts(void);
func GetSoundFonts() string {
	ret, _, _ := syscall.SyscallN(fnGetSoundFonts)
	return windows.BytePtrToString((*byte)(unsafe.Pointer(ret)))
}

//extern DECLSPEC int SDLCALL Mix_SetTimidityCfg(const char *path);
func SetTimidityCfg(path string) m.Int {
	ret, _, _ := syscall.SyscallN(fnSetTimidityCfg, uintptr(unsafe.Pointer(BytePtrFromStringP(path))))
	return (m.Int)(ret)
}

//extern DECLSPEC const char* SDLCALL Mix_GetTimidityCfg(void);
func GetTimidityCfg() string {
	ret, _, _ := syscall.SyscallN(fnGetTimidityCfg)
	return windows.BytePtrToString((*byte)(unsafe.Pointer(ret)))
}

//extern DECLSPEC Mix_Chunk * SDLCALL Mix_GetChunk(int channel);
func GetChunk(channel m.Int) *mm.Chunk {
	ret, _, _ := syscall.SyscallN(fnGetChunk, uintptr(channel))
	return (*mm.Chunk)(unsafe.Pointer(ret))
}

//extern DECLSPEC void SDLCALL Mix_CloseAudio(void);
func CloseAudio() {
	_, _, _ = syscall.SyscallN(fnCloseAudio)
}
