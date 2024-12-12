package mm

import m "github.com/vault-thirteen/SDLW/SDL/model"

// typedef void (SDLCALL *Mix_EffectFunc_t)(int chan, void *stream, int len, void *udata);
type EffectFunc_t *func(channel m.Int, stream *m.Void, len m.Int, udata *m.Void)
