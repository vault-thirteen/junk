package audio

import (
	"time"

	"github.com/vault-thirteen/SDLW/SDL"
	m "github.com/vault-thirteen/SDLW/SDL/model"
	sdlm "github.com/vault-thirteen/SDLW/SDL_Mixer"
	mm "github.com/vault-thirteen/SDLW/SDL_Mixer/model"
)

// GetInfo collects information about audio subsystem.
// SDL audio must be already initialized before using this function.
func GetInfo() (info *Info, err error) {
	info = &Info{}

	// Drivers.
	info.DriversCount = sdl.GetNumAudioDrivers()
	info.Drivers = make([]*Driver, 0, info.DriversCount)
	for i := m.Int(0); i < info.DriversCount; i++ {
		driverName := sdl.GetAudioDriver(i)
		info.Drivers = append(info.Drivers,
			&Driver{Index: i, Name: driverName})
	}

	// Playback Devices.
	info.PlaybackDevices, err = getDevices(0)
	if err != nil {
		return nil, err
	}
	info.PlaybackDevicesCount = m.Int(len(info.PlaybackDevices))

	// Recording Devices.
	info.RecorderDevices, err = getDevices(1)
	if err != nil {
		return nil, err
	}
	info.RecorderDevicesCount = m.Int(len(info.RecorderDevices))

	return info, nil
}

// getDevices lists devices for the specified mode.
func getDevices(devMode m.Int) (devices []*Device, err error) {
	n := sdl.GetNumAudioDevices(devMode)
	devices = make([]*Device, 0, n)

	var device *Device
	for i := m.Int(0); i < n; i++ {
		device, err = getDeviceInfo(i, devMode)
		if err != nil {
			return nil, err
		}

		devices = append(devices, device)
	}

	return devices, nil
}

// getDeviceInfo gets device information.
// 'devMode' is an 'isCapture' argument for 'GetAudioDeviceName' function.
func getDeviceInfo(index m.Int, devMode m.Int) (dev *Device, err error) {
	var deviceName = sdl.GetAudioDeviceName(index, devMode)

	var deviceSpec m.AudioSpec
	ret := sdl.GetAudioDeviceSpec(index, devMode, &deviceSpec)
	if ret != 0 {
		return nil, sdl.GetError()
	}

	return &Device{
		Index: index,
		Name:  deviceName,
		Spec:  &deviceSpec,
	}, nil
}

// DurationOfAudio returns duration of audio.
func DurationOfAudio(as m.AudioSpec, len m.Uint32) time.Duration {
	var sampleSize m.Uint32 = m.Uint32(sdl.AUDIO_BITSIZE(m.Uint16(as.Format))) / 8
	var sampleCount m.Uint32 = len / sampleSize
	var sampleLen m.Uint32
	if as.Channels > 0 {
		sampleLen = sampleCount / m.Uint32(as.Channels)
	} else {
		sampleLen = sampleCount
	}
	var seconds = m.Double(sampleLen) / m.Double(as.Freq)
	return time.Duration(m.Double(time.Second) * seconds)
}

// DurationOfMusic returns duration of music.
func DurationOfMusic(music *mm.Music) time.Duration {
	// SDL library returns -1 length in C code and 0 in Go code.
	// What does it mean ?
	// 1. SDL library does not know how to get duration of MPEG Layer 3 format.
	// 2. Golang ist komplete Scheiße. Das stimmt.
	// See the 'Test' function to get more information.
	dur := sdlm.MusicDuration(music)
	return time.Duration(m.Double(time.Second) * dur)
}
