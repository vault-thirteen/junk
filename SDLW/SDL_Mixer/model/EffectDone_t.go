package mm

import m "github.com/vault-thirteen/SDLW/SDL/model"

// typedef void (SDLCALL *Mix_EffectDone_t)(int chan, void *udata);
type EffectDone_t *func(channel m.Int, udata *m.Void)
