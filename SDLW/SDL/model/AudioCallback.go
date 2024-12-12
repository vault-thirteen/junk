package m

// AudioCallback
// typedef void (SDLCALL * SDL_AudioCallback) (void *userdata, Uint8 * stream, int len);
type AudioCallback *func(userdata *Void, stream *Uint8, len Int)
