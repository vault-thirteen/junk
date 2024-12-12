package m

// AudioStream
/*
struct _SDL_AudioStream

	{
		SDL_AudioCVT cvt_before_resampling;
		SDL_AudioCVT cvt_after_resampling;
		SDL_DataQueue *queue;
		SDL_bool first_run;
		Uint8 *staging_buffer;
		int staging_buffer_size;
		int staging_buffer_filled;
		Uint8 *work_buffer_base;
		int work_buffer_len;
		int src_sample_frame_size;
		SDL_AudioFormat src_format;
		Uint8 src_channels;
		int src_rate;
		int dst_sample_frame_size;
		SDL_AudioFormat dst_format;
		Uint8 dst_channels;
		int dst_rate;
		double rate_incr;
		Uint8 pre_resample_channels;
		int packetlen;
		int resampler_padding_samples;
		float *resampler_padding;
		void *resampler_state;
		SDL_ResampleAudioStreamFunc resampler_func;
		SDL_ResetAudioStreamResamplerFunc reset_resampler_func;
		SDL_CleanupAudioStreamResamplerFunc cleanup_resampler_func;
	};
*/
// typedef struct _SDL_AudioStream SDL_AudioStream;
type AudioStream struct {
	CvtBeforeResampling     AudioCVT
	CvtAfterResampling      AudioCVT
	Queue                   *DataQueue
	FirstRun                Bool
	StagingBuffer           *Uint8
	StagingBufferSize       Int
	StagingBufferFilled     Int
	WorkBufferBase          *Uint8
	WorkBufferLen           Int
	SrcSampleFrameSize      Int
	SrcFormat               AudioFormat
	SrcChannels             Uint8
	SrcRate                 Int
	DstSampleFrameSize      Int
	DstFormat               AudioFormat
	DstChannels             Uint8
	DstRate                 Int
	RateIncr                Double
	PreResampleChannels     Uint8
	Packetlen               Int
	ResamplerPaddingSamples Int
	ResamplerPadding        *Float
	ResamplerState          *Void
	ResamplerFunc           ResampleAudioStreamFunc
	ResetResamplerFunc      ResetAudioStreamResamplerFunc
	CleanupResamplerFunc    CleanupAudioStreamResamplerFunc
}
