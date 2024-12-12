#include <stdlib.h>
#include <SDL.h>
#include "SDL_mixer.h"
#include "console.h"


void playSound(int devIndex, const char *devName, const char *file) {
    SDL_AudioSpec wavSpec;
    Uint8 *wavBuffer;
    Uint32 wavLength;
    SDL_LoadWAV(file, &wavSpec, &wavBuffer, &wavLength);

    SDL_AudioSpec s2;
    SDL_AudioDeviceID deviceId = SDL_OpenAudioDevice(devName, 0, &wavSpec, &s2, 0);
    if (deviceId == 0) {
        console_print("SDL_OpenAudioDevice: %s.\n", SDL_GetError());
        return;
    }
    console_print("Device ID: %d.", deviceId);

    if (SDL_QueueAudio(deviceId, wavBuffer, wavLength) != 0) {
        console_print("SDL_QueueAudio: %s.\n", SDL_GetError());
    }
    SDL_PauseAudioDevice(deviceId, 0);
    Sleep(5 * 1000);

    SDL_CloseAudioDevice(deviceId);
    SDL_FreeWAV(wavBuffer);
}


void playMusic(int devIndex, const char *devName, const char *file) {
    if (Mix_OpenAudioDevice(48000, AUDIO_S16, 2, 2048, devName, SDL_AUDIO_ALLOW_ANY_CHANGE) != 0) {
        console_print("SDL_QueueAudio: %s.\n", SDL_GetError());
        return;
    }

    Mix_Music *music = Mix_LoadMUS(file);
    double duration = Mix_MusicDuration(music);
    console_print("duration: %f.\n", duration);

    Mix_FreeMusic(music);
    Mix_CloseAudio();
}


int APIENTRY WinMain(HINSTANCE hInstance, HINSTANCE hPrevInstance, LPSTR lpCmdLine, int nShowCmd) {
    int exitCode = EXIT_FAILURE;
    AllocConsole();
    if (SDL_Init(SDL_INIT_AUDIO) != 0) {
        console_print("SDL_Init: %s.\n", SDL_GetError());
        goto failure;
    }
    if (Mix_Init(MIX_INIT_MP3) == 0) {
        console_print("Mix_Init: %s.\n", Mix_GetError());
        goto failure;
    }

    SDL_version ver;
    SDL_GetVersion(&ver);
    console_print("SDL Version: %d.%d.%d.\n", ver.major, ver.minor, ver.patch);

    int devMode = 0;
    int devicesCount = SDL_GetNumAudioDevices(devMode);
    int devIndex = 0;
    const char *devName = SDL_GetAudioDeviceName(devIndex, devMode);
    console_print("Device Name: %s.\n", devName);
    SDL_AudioSpec deviceSpec;
    if (SDL_GetAudioDeviceSpec(devIndex, devMode, &deviceSpec) != 0) {
        console_print("SDL_GetAudioDeviceSpec: %s.\n", SDL_GetError());
        goto failure;
    }

    playSound(devIndex, devName, "C:\\Windows\\Media\\chimes.wav");
    playMusic(devIndex, devName, "D:\\Temp\\1\\music.mp3");
    goto success;

    failure:
    goto fin;

    success:
    exitCode = EXIT_SUCCESS;
    goto fin;

    fin:
    Mix_Quit();
    SDL_Quit();
    FreeConsole();
    return exitCode;
}
