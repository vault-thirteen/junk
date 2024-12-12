extern DECLSPEC int SDLCALL SDL_SetError(SDL_PRINTF_FORMAT_STRING const char *fmt, ...) SDL_PRINTF_VARARG_FUNC(1);
extern DECLSPEC const char *SDLCALL SDL_GetError(void);
extern DECLSPEC char * SDLCALL SDL_GetErrorMsg(char *errstr, int maxlen);
extern DECLSPEC void SDLCALL SDL_ClearError(void);
typedef enum { SDL_ENOMEM, SDL_EFREAD, SDL_EFWRITE, SDL_EFSEEK, SDL_UNSUPPORTED, SDL_LASTERROR } SDL_errorcode;
extern DECLSPEC int SDLCALL SDL_Error(SDL_errorcode code);
