package chip8

import (
	"github.com/veandco/go-sdl2/sdl"
	"log"
)

type APU struct {
	DeviceID sdl.AudioDeviceID
	Specs    sdl.AudioSpec
}

func APUInit() *APU {
	audioSpecs := sdl.AudioSpec{
		Freq:     44100,
		Format:   sdl.AUDIO_F32LSB,
		Channels: 1,
		Samples:  2048,
	}

	specs := sdl.AudioSpec{}
	if sdl.GetNumAudioDevices(false) <= 0 {
		log.Fatal("No Audio devices found!")
	}

	audioDevice, err := sdl.OpenAudioDevice("", false, &audioSpecs, &specs, sdl.AUDIO_ALLOW_ANY_CHANGE)
	if err != nil {
		log.Fatal("err")
	}

	apu := APU{
		DeviceID: audioDevice,
		Specs:    specs,
	}
	sdl.PauseAudioDevice(apu.DeviceID, false)
	return &apu
}
