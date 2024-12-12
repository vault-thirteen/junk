package mm

// void (SDLCALL *channel_finished)(int channel)
type ChannelFinished *func(channel int)
