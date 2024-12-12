#include "SDL.h"
#include "SDL_mixer.h"
#include "music.h"

// Why the fuck is this type declared in a code file instead of header file ?!
// Looks like developers of SDL library were drunk when they wrote all this ...
// Source: src/music.c.
struct _Mix_Music {
    Mix_MusicInterface *interface;
    void *context;

    SDL_bool playing;
    Mix_Fading fading;
    int fade_step;
    int fade_steps;

    char filename[1024];
};
