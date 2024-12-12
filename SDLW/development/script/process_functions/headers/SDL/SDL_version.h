typedef struct SDL_version { Uint8 major; Uint8 minor; Uint8 patch; } SDL_version;
extern DECLSPEC void SDLCALL SDL_GetVersion(SDL_version * ver);
extern DECLSPEC const char *SDLCALL SDL_GetRevision(void);
extern SDL_DEPRECATED DECLSPEC int SDLCALL SDL_GetRevisionNumber(void);
