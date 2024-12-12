#include <SDL.h>
#include <stdio.h>
#include <stdlib.h>

int WinMain(int argc, char *argv[]) {
    if (SDL_Init(SDL_INIT_EVERYTHING) != 0) {
        fprintf(stderr, "SDL_Init Error: %s\n", SDL_GetError());
        return EXIT_FAILURE;
    }

    SDL_Window *win = SDL_CreateWindow("Hello World!", SDL_WINDOWPOS_CENTERED, SDL_WINDOWPOS_CENTERED, 640, 480, SDL_WINDOW_SHOWN);
    if (win == NULL) {
        fprintf(stderr, "SDL_CreateWindow Error: %s\n", SDL_GetError());
        return EXIT_FAILURE;
    }

    SDL_Surface *surf = SDL_GetWindowSurface(win);
    if (surf == NULL) {
        fprintf(stderr, "SDL_GetWindowSurface Error: %s\n", SDL_GetError());
        SDL_DestroyWindow(win);
        SDL_Quit();
        return EXIT_FAILURE;
    }

    // Windows O.S. has pixel format SDL_PIXELFORMAT_RGB888.
    Uint32 wpf = SDL_GetWindowPixelFormat(win);
    if (wpf == SDL_PIXELFORMAT_UNKNOWN) {
        fprintf(stderr, "SDL_GetWindowPixelFormat Error: %s\n", SDL_GetError());
        SDL_DestroyWindow(win);
        SDL_Quit();
        return EXIT_FAILURE;
    }

    SDL_Color col = {255, 255, 255, 255};
    Uint32 pixel = SDL_MapRGBA(surf->format, col.r, col.g, col.b, col.a);

    errno_t err = SDL_FillRect(surf, NULL, pixel);
    if (err != 0) {
        fprintf(stderr, "SDL_FillRect Error: %s\n", SDL_GetError());
        SDL_DestroyWindow(win);
        SDL_Quit();
        return EXIT_FAILURE;
    }

    err = SDL_UpdateWindowSurface(win);
    if (err != 0) {
        fprintf(stderr, "SDL_UpdateWindowSurface Error: %s\n", SDL_GetError());
        SDL_DestroyWindow(win);
        SDL_Quit();
        return EXIT_FAILURE;
    }

    char running = 1;
    SDL_Event event;
    while (running) {
        //while (SDL_PollEvent(&event)) {} // For heavy games.
        if (SDL_WaitEvent(&event) != 1) {
            fprintf(stderr, "SDL_WaitEvent Error: %s\n", SDL_GetError());
            SDL_Delay(100);
        }

        if (event.type == SDL_QUIT) {
            fprintf(stderr, "SDL_QUIT Event, TS: %d\n", event.quit.timestamp);
            running = 0;
            break;
        }
    }

    SDL_DestroyWindow(win);
    SDL_Quit();

    return EXIT_SUCCESS;
}
