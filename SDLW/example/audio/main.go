package main

import (
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/vault-thirteen/SDLW/SDL"
	m "github.com/vault-thirteen/SDLW/SDL/model"
	sdlm "github.com/vault-thirteen/SDLW/SDL_Mixer"
	mm "github.com/vault-thirteen/SDLW/SDL_Mixer/model"
	"github.com/vault-thirteen/SDLW/helper/audio"
	"github.com/vault-thirteen/SDLW/helper/cpu"
	"github.com/vault-thirteen/SDLW/win32"
)

const ErrEndianness = "endianness error"

func main() {
	var err error
	err = win32.LoadLibrary()
	mustBeNoError(err)

	err = sdl.LoadLibrary(sdl.SdlDll)
	mustBeNoError(err)

	sdlMixerExternalFunctions := sdlm.ExternalFunctions{
		ExtFnGetError: sdl.GetFnGetError(),
	}
	err = sdlm.LoadLibrary(sdlm.SdlMixDll, &sdlMixerExternalFunctions)
	mustBeNoError(err)

	err = work()
	mustBeNoError(err)
}

func mustBeNoError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func work() (err error) {
	// Show versions.
	fmt.Println(fmt.Sprintf("SDL Version: %v.", sdl.GetVersion().Text()))
	fmt.Println(fmt.Sprintf("SDL Mixer Version: %v.", sdlm.Linked_Version().Text()))

	// Initialization of SDL.
	if sdl.Init(sdl.INIT_AUDIO) != 0 {
		return sdl.GetError()
	}
	defer sdl.Quit()
	sdlInitStatus := sdl.WasInit(0)
	fmt.Println(fmt.Sprintf("SDL Initialization Status: %v.", sdlInitStatus))

	// Initialization of SDL Mixer.
	sdlMixerInitFlags := sdlm.Init(m.Int(mm.INIT_MP3))
	if sdlMixerInitFlags == 0 {
		return sdlm.GetError()
	}
	defer sdlm.Quit()
	fmt.Println(fmt.Sprintf("SDL Mixer Initialization Status: %v.", sdlMixerInitFlags))

	// Audio Info.
	var audioInfo *audio.Info
	audioInfo, err = audio.GetInfo()
	if err != nil {
		return err
	}
	device := audioInfo.PlaybackDevices[0]
	fmt.Println(fmt.Sprintf("Using device: %s.", device.Name))

	err = playSound(device)
	if err != nil {
		return err
	}

	err = playMusic(device)
	if err != nil {
		return err
	}

	return nil
}

func playSound(device *audio.Device) (err error) {
	// Load a WAV file.
	filePath := `C:\Windows\Media\chimes.wav`
	var wavSpec m.AudioSpec
	var wavBuffer *m.Uint8
	var wavLength m.Uint32
	_ = sdl.LoadWAV(filePath, &wavSpec, &wavBuffer, &wavLength)
	defer func() {
		sdl.FreeWAV(wavBuffer)
	}()

	// Open the device.
	var s2 m.AudioSpec
	deviceId := sdl.OpenAudioDevice(device.Name, 0, &wavSpec, &s2, 0)
	if deviceId == 0 {
		return sdl.GetError()
	}
	defer sdl.CloseAudioDevice(deviceId)
	fmt.Println(fmt.Sprintf("Device ID: %d.", deviceId))

	// Play audio.
	ret := sdl.QueueAudio(deviceId, wavBuffer, wavLength)
	if ret != 0 {
		return sdl.GetError()
	}
	sdl.PauseAudioDevice(deviceId, 0)

	duration := audio.DurationOfAudio(wavSpec, wavLength)
	time.Sleep(duration)

	return nil
}

func playMusic(device *audio.Device) (err error) {
	var format m.Uint16
	switch cpu.HardwareEndianness() {
	case cpu.EndiannessLittleEndian:
		format = sdl.AUDIO_S16LSB
	case cpu.EndiannessBigEndian:
		format = sdl.AUDIO_S16MSB
	default:
		return errors.New(ErrEndianness)
	}

	if sdlm.OpenAudioDevice(48000, format, 2, 2048, device.Name, sdl.AUDIO_ALLOW_ANY_CHANGE) != 0 {
		return sdlm.GetError()
	}
	defer sdl.CloseAudio()

	// Load an MP3 file.
	// External file. Copy it before starting the example.
	filePath := `Music.mp3`
	music := sdlm.LoadMUS(filePath)
	if music == nil {
		return sdlm.GetError()
	}
	defer sdlm.FreeMusic(music)

	if sdlm.PlayMusic(music, 1) != 0 {
		return sdlm.GetError()
	}

	duration := audio.DurationOfMusic(music)
	time.Sleep(duration)

	sdlm.Test(filePath) // This shit does not work.

	return nil
}
