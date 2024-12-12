package mm

import "github.com/vault-thirteen/SDLW/SDL/model"

// Music
/*
struct _Mix_Music {
    Mix_MusicInterface *interface;
    void *context;

    SDL_bool playing;
    Mix_Fading fading;
    int fade_step;
    int fade_steps;

    char filename[1024];
};
typedef struct _Mix_Music Mix_Music;
*/
type Music struct {
	Interface *MusicInterface
	Context   *m.Void

	Playing   m.Bool
	Fading    Fading
	FadeStep  m.Int
	FadeSteps m.Int

	Filename [1024]m.Char
}
