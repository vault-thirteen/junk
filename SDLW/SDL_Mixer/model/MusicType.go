package mm

import m "github.com/vault-thirteen/SDLW/SDL/model"

// MusicType
/*
typedef enum {
    MUS_NONE,
    MUS_CMD,
    MUS_WAV,
    MUS_MOD,
    MUS_MID,
    MUS_OGG,
    MUS_MP3,
    MUS_MP3_MAD_UNUSED,
    MUS_FLAC,
    MUS_MODPLUG_UNUSED,
    MUS_OPUS
} Mix_MusicType;
*/
type MusicType = m.Int

const (
	MUS_NONE           = MusicType(0)
	MUS_CMD            = MusicType(1)
	MUS_WAV            = MusicType(2)
	MUS_MOD            = MusicType(3)
	MUS_MID            = MusicType(4)
	MUS_OGG            = MusicType(5)
	MUS_MP3            = MusicType(6)
	MUS_MP3_MAD_UNUSED = MusicType(7)
	MUS_FLAC           = MusicType(8)
	MUS_MODPLUG_UNUSED = MusicType(9)
	MUS_OPUS           = MusicType(10)
)
