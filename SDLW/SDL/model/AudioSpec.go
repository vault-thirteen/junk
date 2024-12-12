package m

// AudioSpec
/*
typedef struct SDL_AudioSpec
{
	int freq;
	SDL_AudioFormat format;
	Uint8 channels;
	Uint8 silence;
	Uint16 samples;
	Uint16 padding;
	Uint32 size;
	SDL_AudioCallback callback;
	void *userdata;
} SDL_AudioSpec;
*/
type AudioSpec struct {
	Freq     Int
	Format   AudioFormat
	Channels Uint8
	Silence  Uint8
	Samples  Uint16
	Padding  Uint16
	Size     Uint32
	Callback AudioCallback
	Userdata *Void
}
