package m

/*
typedef enum

	{
	    SDL_ENOMEM,
	    SDL_EFREAD,
	    SDL_EFWRITE,
	    SDL_EFSEEK,
	    SDL_UNSUPPORTED,
	    SDL_LASTERROR
	} SDL_errorcode;
*/
type Errorcode = Int

const (
	SDL_ENOMEM      = Errorcode(0)
	SDL_EFREAD      = Errorcode(1)
	SDL_EFWRITE     = Errorcode(2)
	SDL_EFSEEK      = Errorcode(3)
	SDL_UNSUPPORTED = Errorcode(4)
	SDL_LASTERROR   = Errorcode(5)
)
