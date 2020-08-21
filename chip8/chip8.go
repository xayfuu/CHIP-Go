package chip8

import (
	"encoding/binary"
	"github.com/veandco/go-sdl2/sdl"
	"math"
)

type Chip8 struct {
	Engine *SDL
	CPU    *CPU
	APU    *APU
}

func Init() *Chip8 {
	return &Chip8{
		Engine: SDLInit(),
		CPU:    CPUInit(),
		APU:    APUInit(),
	}
}

func (c8 *Chip8) UpdateDisplay() error {
	err := c8.Engine.Renderer.SetRenderTarget(c8.Engine.Screen)
	if err != nil {
		return err
	}

	c8.Engine.Renderer.SetDrawColor(38, 38, 38, 255)
	c8.Engine.Renderer.Clear()

	c8.Engine.Renderer.SetDrawColor(242, 242, 242, 255)

	for py := int32(0); py < H; py++ {
		for px := int32(0); px < W; px++ {
			if c8.CPU.Display[py*W+px] == 0x1 {
				c8.Engine.Renderer.DrawPoint(px, py)
			}
		}
	}
	c8.Engine.Renderer.SetRenderTarget(nil)
	displayRect := sdl.Rect{
		W: int32(W),
		H: int32(H),
	}
	c8.Engine.Renderer.Copy(c8.Engine.Screen, &displayRect, &sdl.Rect{W: SDL_SCREEN_WIDTH, H: SDL_SCREEN_HEIGHT})
	c8.Engine.Renderer.Present()

	return nil
}

func (c8 *Chip8) IsAlive() bool {
	for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
		switch e := event.(type) {
		case *sdl.QuitEvent:
			return false
		case *sdl.KeyboardEvent:
			key, keyExists := c8.Engine.Keymap[e.Keysym.Scancode]
			if e.GetType() == sdl.KEYUP {
				if keyExists {
					c8.CPU.KeyRelease(key)
				}
			} else {
				if keyExists {
					c8.CPU.KeyPress(key)
				}
			}
		}
	}
	return true
}

func (c8 *Chip8) Beep() {
	if c8.APU.DeviceID != 0 && c8.CPU.ST > 0 {
		sample := make([]byte, 4)
		binary.LittleEndian.PutUint32(sample, math.Float32bits(1.0))
		n := int(c8.APU.Specs.Channels) * int(c8.APU.Specs.Samples) * 4
		data := make([]byte, n)

		// 128 samples per 1/60 of a second
		for i := 0; i < n; i += 4 {
			copy(data[i:], sample)
		}

		if err := sdl.QueueAudio(c8.APU.DeviceID, data); err != nil {
			println(err)
		}
	}
}
