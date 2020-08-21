package main

import "C"

import (
	"fmt"
	"log"
	"math/rand"
	"runtime"
	"time"

	"./chip8"
)

var (
	c8 *chip8.Chip8
)

var romPath = "./roms/"

func init() {
	runtime.LockOSThread()
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	rand.Seed(time.Now().UTC().UnixNano())
	c8 = chip8.Init()
}

func main() {
	c8.CPU.LoadGame(fmt.Sprintf("%s%s", romPath, "PONG"))

	cpuTick := time.NewTicker(time.Millisecond)
	frameTick := time.NewTicker(time.Second / 60)
	soundTick := time.NewTicker(time.Second / 60)

	for c8.IsAlive() {
		select {

		case <-frameTick.C:
			c8.UpdateDisplay()

		case <-soundTick.C:
			//apu.Beep()

		case <-cpuTick.C:
			c8.Beep()
			err := c8.CPU.NextTick()
			if err != nil {
				log.Fatal(err)
			}
		}
	}
}
