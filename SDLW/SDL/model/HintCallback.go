package m

// HintCallback
// void (SDLCALL *SDL_HintCallback)(void *userdata, const char *name, const char *oldValue, const char *newValue);
// SDL_hints.h.
type HintCallback *func(userdata *Void, name *Char, oldValue *Char, newValue *Char)
