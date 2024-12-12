package audio

import m "github.com/vault-thirteen/SDLW/SDL/model"

type Info struct {
	DriversCount m.Int
	Drivers      []*Driver

	PlaybackDevicesCount m.Int
	PlaybackDevices      []*Device

	RecorderDevicesCount m.Int
	RecorderDevices      []*Device
}
