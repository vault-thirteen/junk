package m

// AudioFilter
// typedef void (SDLCALL * SDL_AudioFilter) (struct SDL_AudioCVT * cvt, SDL_AudioFormat format);
type AudioFilter *func(cvt *AudioCVT, format AudioFormat)
