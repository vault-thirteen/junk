package mm

import m "github.com/vault-thirteen/SDLW/SDL/model"

// EffectDone
// typedef void (SDLCALL *Mix_EffectDone_t)(int chan, void *udata);
type EffectDone *func(channel m.Int, udata *m.Void)
