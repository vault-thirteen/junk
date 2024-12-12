package m

// CleanupAudioStreamResamplerFunc
// typedef void (*SDL_CleanupAudioStreamResamplerFunc)(SDL_AudioStream *stream);
type CleanupAudioStreamResamplerFunc *func(stream *AudioStream)
