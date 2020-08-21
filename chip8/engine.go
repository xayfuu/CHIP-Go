package chip8

import (
	"github.com/veandco/go-sdl2/sdl"
	"log"
)

const (
	SDL_WINDOW_FLAGS   = sdl.WINDOW_OPENGL
	SDL_SCREEN_FORMAT  = sdl.PIXELFORMAT_RGB888
	SDL_SCREEN_HEIGHT  = H * 0xA
	SDL_SCREEN_WIDTH   = W * 0xA
	SDL_TEXTURE_WIDTH  = H * 0x2
	SDL_TEXTURE_HEIGHT = W * 0x2
	SDL_TEXTURE_ACCESS = sdl.TEXTUREACCESS_TARGET
)

type SDL struct {
	Window   *sdl.Window
	Renderer *sdl.Renderer
	Screen   *sdl.Texture
	Keymap   map[sdl.Scancode]uint16
}

var keymapping = map[sdl.Scancode]uint16{
	sdl.SCANCODE_X: 0x0,
	sdl.SCANCODE_1: 0x1,
	sdl.SCANCODE_2: 0x2,
	sdl.SCANCODE_3: 0x3,
	sdl.SCANCODE_Q: 0x4,
	sdl.SCANCODE_W: 0x5,
	sdl.SCANCODE_E: 0x6,
	sdl.SCANCODE_A: 0x7,
	sdl.SCANCODE_S: 0x8,
	sdl.SCANCODE_D: 0x9,
	sdl.SCANCODE_Y: 0xA, // QWERTZ keyboard, change to sdl.SCANCODE_Z if QWERTY keyboard
	sdl.SCANCODE_C: 0xB,
	sdl.SCANCODE_4: 0xC,
	sdl.SCANCODE_R: 0xD,
	sdl.SCANCODE_F: 0xE,
	sdl.SCANCODE_V: 0xF,
}

func SDLInit() *SDL {
	err := sdl.Init(sdl.INIT_VIDEO | sdl.INIT_AUDIO)
	if err != nil {
		log.Fatal(err)
	}
	window, renderer, err := sdl.CreateWindowAndRenderer(SDL_SCREEN_WIDTH, SDL_SCREEN_HEIGHT, uint32(SDL_WINDOW_FLAGS))
	if err != nil {
		log.Fatal(err)
	}
	window.SetTitle("CHIP-8 Emulator")

	screen, err := renderer.CreateTexture(uint32(SDL_SCREEN_FORMAT), SDL_TEXTURE_ACCESS, SDL_TEXTURE_WIDTH, SDL_TEXTURE_HEIGHT)
	if err != nil {
		log.Fatal(err)
	}
	_sdl := SDL{
		Window:   window,
		Renderer: renderer,
		Screen:   screen,
		Keymap:   keymapping,
	}
	return &_sdl
}
