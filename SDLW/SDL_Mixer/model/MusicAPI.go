package mm

import m "github.com/vault-thirteen/SDLW/SDL/model"

// MusicAPI
/*
typedef enum
{
    MIX_MUSIC_CMD,
    MIX_MUSIC_WAVE,
    MIX_MUSIC_MODPLUG,
    MIX_MUSIC_FLUIDSYNTH,
    MIX_MUSIC_TIMIDITY,
    MIX_MUSIC_NATIVEMIDI,
    MIX_MUSIC_OGG,
    MIX_MUSIC_DRMP3,
    MIX_MUSIC_MPG123,
    MIX_MUSIC_DRFLAC,
    MIX_MUSIC_FLAC,
    MIX_MUSIC_OPUS,
    MIX_MUSIC_LIBXMP,
    MIX_MUSIC_LAST
} Mix_MusicAPI;
*/
type MusicAPI = m.Int

const (
	MIX_MUSIC_CMD        = MusicAPI(0)
	MIX_MUSIC_WAVE       = MusicAPI(1)
	MIX_MUSIC_MODPLUG    = MusicAPI(2)
	MIX_MUSIC_FLUIDSYNTH = MusicAPI(3)
	MIX_MUSIC_TIMIDITY   = MusicAPI(4)
	MIX_MUSIC_NATIVEMIDI = MusicAPI(5)
	MIX_MUSIC_OGG        = MusicAPI(6)
	MIX_MUSIC_DRMP3      = MusicAPI(7)
	MIX_MUSIC_MPG123     = MusicAPI(8)
	MIX_MUSIC_DRFLAC     = MusicAPI(9)
	MIX_MUSIC_FLAC       = MusicAPI(10)
	MIX_MUSIC_OPUS       = MusicAPI(11)
	MIX_MUSIC_LIBXMP     = MusicAPI(12)
	MIX_MUSIC_LAST       = MusicAPI(13)
)
