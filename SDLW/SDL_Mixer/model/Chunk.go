package mm

import m "github.com/vault-thirteen/SDLW/SDL/model"

// Chunk
/*
typedef struct Mix_Chunk {
    int allocated;
    Uint8 *abuf;
    Uint32 alen;
    Uint8 volume;
} Mix_Chunk;
*/
type Chunk struct {
	Allocated m.Int
	Abuf      *m.Uint8
	Alen      m.Uint32
	Volume    m.Uint8
}
