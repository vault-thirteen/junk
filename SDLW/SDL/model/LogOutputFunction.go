package m

// LogOutputFunction
// void (SDLCALL *SDL_LogOutputFunction)(void *userdata, int category, SDL_LogPriority priority, const char *message);
// SDL_log.h
type LogOutputFunction *func(userdata *Void, category Int, priority LogPriority, message *Char)
