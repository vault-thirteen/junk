package mm

import m "github.com/vault-thirteen/SDLW/SDL/model"

/*
typedef enum

	{
	    MIX_INIT_FLAC   = 0x00000001,
	    MIX_INIT_MOD    = 0x00000002,
	    MIX_INIT_MP3    = 0x00000008,
	    MIX_INIT_OGG    = 0x00000010,
	    MIX_INIT_MID    = 0x00000020,
	    MIX_INIT_OPUS   = 0x00000040
	} MIX_InitFlags;
*/
type InitFlags = m.Int

const (
	INIT_FLAC = InitFlags(0x00000001)
	INIT_MOD  = InitFlags(0x00000002)
	INIT_MP3  = InitFlags(0x00000008)
	INIT_OGG  = InitFlags(0x00000010)
	INIT_MID  = InitFlags(0x00000020)
	INIT_OPUS = InitFlags(0x00000040)
)
