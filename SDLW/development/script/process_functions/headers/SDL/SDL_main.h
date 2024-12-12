typedef int (*SDL_main_func)(int argc, char *argv[]);
extern DECLSPEC void SDLCALL SDL_SetMainReady(void);
extern DECLSPEC int SDLCALL SDL_RegisterApp(const char *name, Uint32 style, void *hInst);
extern DECLSPEC void SDLCALL SDL_UnregisterApp(void);
