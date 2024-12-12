package mm

import m "github.com/vault-thirteen/SDLW/SDL/model"

// void (SDLCALL *mix_func)(void *udata, Uint8 *stream, int len)
type MixFunc *func(udata *m.Void, stream *m.Uint8, len m.Int)
