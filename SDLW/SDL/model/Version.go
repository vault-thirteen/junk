package m

import (
	"fmt"
)

// Version
/*
typedef struct SDL_version
{
	Uint8 major;
	Uint8 minor;
	Uint8 patch;
} SDL_version;
*/
// SDL_version.h
type Version struct {
	Major Uint8
	Minor Uint8
	Patch Uint8
}

func (v Version) Text() string {
	return fmt.Sprintf("%d.%d.%d", v.Major, v.Minor, v.Patch)
}
