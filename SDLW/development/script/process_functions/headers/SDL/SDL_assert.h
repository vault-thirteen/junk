typedef enum { SDL_ASSERTION_RETRY, SDL_ASSERTION_BREAK, SDL_ASSERTION_ABORT, SDL_ASSERTION_IGNORE, SDL_ASSERTION_ALWAYS_IGNORE } SDL_AssertState;
typedef struct SDL_AssertData { int always_ignore; unsigned int trigger_count; const char *condition; const char *filename; int linenum; const char *function; const struct SDL_AssertData *next; } SDL_AssertData;
extern DECLSPEC SDL_AssertState SDLCALL SDL_ReportAssertion(SDL_AssertData *, const char *, const char *, int)
typedef SDL_AssertState (SDLCALL *SDL_AssertionHandler)(const SDL_AssertData* data, void* userdata);
extern DECLSPEC void SDLCALL SDL_SetAssertionHandler(SDL_AssertionHandler handler, void *userdata);
extern DECLSPEC SDL_AssertionHandler SDLCALL SDL_GetDefaultAssertionHandler(void);
extern DECLSPEC SDL_AssertionHandler SDLCALL SDL_GetAssertionHandler(void **puserdata);
extern DECLSPEC const SDL_AssertData * SDLCALL SDL_GetAssertionReport(void);
extern DECLSPEC void SDLCALL SDL_ResetAssertionReport(void);
