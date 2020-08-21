package chip8

import (
	"errors"
	"io"
	"log"
	"math/rand"
	"time"
)

const (
	NANO       = 0x3B9ACA00
	PC_INITIAL = 0x200
	H          = 0x20
	W          = 0x40
)

type CPU struct {
	Memory  [4096]byte
	PC      uint16
	Display [W * H]byte
	V       [16]byte
	I       uint16
	Keys    [16]bool
	DT      byte
	ST      byte
	Stack   [16]uint16
	SP      uint16
	WKey    *byte
	Cycle   int64
	Clock   int64
	Ticks   int64
}

var fontset = []byte{
	0xF0, 0x90, 0x90, 0x90, 0xF0, // 0
	0x20, 0x60, 0x20, 0x20, 0x70, // 1
	0xF0, 0x10, 0xF0, 0x80, 0xF0, // 2
	0xF0, 0x10, 0xF0, 0x10, 0xF0, // 3
	0x90, 0x90, 0xF0, 0x10, 0x10, // 4
	0xF0, 0x80, 0xF0, 0x10, 0xF0, // 5
	0xF0, 0x80, 0xF0, 0x90, 0xF0, // 6
	0xF0, 0x10, 0x20, 0x40, 0x40, // 7
	0xF0, 0x90, 0xF0, 0x90, 0xF0, // 8
	0xF0, 0x90, 0xF0, 0x10, 0xF0, // 9
	0xF0, 0x90, 0xF0, 0x90, 0x90, // A
	0xE0, 0x90, 0xE0, 0x90, 0xE0, // B
	0xF0, 0x80, 0x80, 0x80, 0xF0, // C
	0xE0, 0x90, 0x90, 0x90, 0xE0, // D
	0xF0, 0x80, 0xF0, 0x80, 0xF0, // E
	0xF0, 0x80, 0xF0, 0x80, 0x80, // F
}

func CPUInit() *CPU {
	c := CPU{
		I:     0,
		SP:    0,
		Cycle: 0,
		WKey:  nil,
		DT:    0,
		ST:    0,
		PC:    PC_INITIAL,
		Ticks: 1000,
		Clock: time.Now().UnixNano(),
	}
	c.InitialState()
	c.LoadFontset()
	return &c
}

func (c *CPU) InitialState() *CPU {
	c.ClearDisplay()
	c.ClearStack()
	c.ClearV()
	c.ClearMemory()
	return c
}

func (c *CPU) ClearDisplay() *CPU {
	for x := 0; x < W*H; x++ {
		c.Display[x] = 0x0
	}
	return c
}

func (c *CPU) ClearStack() *CPU {
	for i, _ := range c.Stack {
		c.Stack[i] = 0x0
	}
	return c
}

func (c *CPU) ClearV() *CPU {
	for i, _ := range c.V {
		c.V[i] = 0x0
	}
	return c
}

func (c *CPU) ClearMemory() *CPU {
	for i, _ := range c.Memory {
		c.Memory[i] = 0x0
	}
	return c
}

func (c *CPU) MemorySize() int64 {
	return 0x1000 - PC_INITIAL
}

func (c *CPU) GetEmptyMemoryArr() []byte {
	return make([]byte, c.MemorySize())
}

func (c *CPU) IsInInitialState() error {
	if c.PC != PC_INITIAL {
		return errors.New("CPU is not in initial state.")
	}
	return nil
}

func (c *CPU) LoadFontset() *CPU {
	copy(c.Memory[0x0:], fontset)
	return c
}

func (c *CPU) LoadGame(fp string) (*ROM, error) {
	if err := c.IsInInitialState(); err != nil {
		return nil, err
	}
	rom, err := ReadROM(fp)

	if err := rom.CheckMemoryOverflow(c.MemorySize()); err != nil {
		return nil, err
	}
	memLoader := c.GetEmptyMemoryArr()

	if _, err = rom.Reader.Read(memLoader); err != nil && err != io.EOF {
		return nil, err
	}
	defer rom.Reader.Close()
	copy(c.Memory[PC_INITIAL:], memLoader)

	rom.InitialState = memLoader
	return rom, nil
}

func (c *CPU) NextInstruction() uint16 {
	pc := c.PC
	c.PC += 2
	return uint16(c.Memory[pc])<<8 | uint16(c.Memory[pc+1])
}

func (c *CPU) NextTick() error {
	now := time.Now().UnixNano()
	cycles := (now - c.Clock) * c.Ticks / NANO

	for c.Cycle < cycles {
		err := c.RunCycle()
		if err != nil {
			return err
		}
		if c.WKey != nil {
			c.Cycle = cycles
		}
	}

	return nil
}

// https://en.wikipedia.org/wiki/CHIP-8#Opcode_table
// http://devernay.free.fr/hacks/chip8/C8TECH10.HTM
func (c *CPU) RunCycle() error {
	if c.WKey != nil {
		log.Printf("\tCPU[%v]:\t\tWaiting for Key input...\n", c.Cycle)
		return nil
	}

	c.Cycle += 1
	op := c.NextInstruction()

	opset := op & 0xF000

	addr := op & 0x0FFF
	kk := op & 0x00FF
	nibble := op & 0x000F

	x := (op & 0x0F00) >> 8
	y := (op & 0x00F0) >> 4

	switch opset {
	case 0x0000:
		switch kk {
		case 0x0000:
		case 0x00E0:
			c.ClearDisplay()
		case 0x00EE:
			c.SP -= 1
			c.PC = c.Stack[c.SP]
		}
	case 0x1000:
		c.PC = addr
	case 0x2000:
		c.Stack[c.SP] = c.PC
		c.SP += 1
		c.PC = addr
	case 0x3000:
		if c.V[x] == byte(kk) {
			c.PC += 2
		}
	case 0x4000:
		if c.V[x] != byte(kk) {
			c.PC += 2
		}
	case 0x5000:
		if c.V[x] == c.V[y] {
			c.PC += 2
		}
	case 0x6000:
		c.V[x] = byte(kk)
	case 0x7000:
		c.V[x] += byte(kk)
	case 0x8000:
		switch nibble {
		case 0x0000:
			c.V[x] = c.V[y]
		case 0x0001:
			c.V[x] = c.V[x] | c.V[y]
		case 0x0002:
			c.V[x] = c.V[x] & c.V[y]
		case 0x0003:
			c.V[x] = c.V[x] ^ c.V[y]
		case 0x0004:
			if c.V[x] < c.V[y] {
				c.V[0x000F] = 1
			} else {
				c.V[0x000F] = 0
			}
			c.V[x] += c.V[y]
		case 0x0005:
			if c.V[x] >= c.V[y] {
				c.V[0x000F] = 1
			} else {
				c.V[0x000F] = 0
			}
			c.V[x] -= c.V[y]
		case 0x0006:
			if c.V[x]&0x0001 == 1 {
				c.V[0x000F] = 1
			} else {
				c.V[0x000F] = 0
			}
			c.V[x] >>= 1
		case 0x0007:
			if c.V[y] >= c.V[x] {
				c.V[0x000F] = 1
			} else {
				c.V[0x000F] = 0
			}
			c.V[x] = c.V[y] - c.V[x]
		case 0x000E:
			c.V[0xF] = c.V[x] >> 7
			c.V[x] <<= 1
		}
	case 0x9000:
		if c.V[x] != c.V[y] {
			c.PC += 2
		}
	case 0xA000:
		c.I = addr
	case 0xB000:
		c.PC = addr + uint16(c.V[0x0000])
	case 0xC000:
		c.V[x] = byte(rand.Intn(256)) & byte(kk)
	case 0xD000:

		collision := false
		sm := c.Memory[c.I:]
		x := c.V[x]
		y := c.V[y]
		for iy := uint16(0); iy < nibble; iy++ {
			for ix := uint16(0); ix < 8; ix++ {
				tx := int(x) + int(ix)
				ty := int(y) + int(iy)
				if tx >= W || ty >= H {
					continue
				}
				s := c.Display[ty*W+tx]
				d := (sm[iy] >> (7 - ix)) & 0x01
				c.Display[ty*W+tx] ^= byte(d)
				if s == 1 && d == 1 {
					collision = true
				}
			}
		}
		if collision {
			c.V[0xf] = 1
		} else {
			c.V[0xf] = 0
		}
	case 0xE000:
		switch kk {
		case 0x00A1:
			if c.Keys[c.V[x]] {
				c.PC += 2
			}
		case 0x00E9:
			if c.Keys[c.V[x]] {
				c.PC += 2
			}
		}
	case 0xF000:
		switch kk {
		case 0x0007:
			c.V[x] = c.DT
		case 0x000A:
			c.WKey = &c.V[x]
		case 0x0015:
			c.DT = c.V[x]
		case 0x0018:
			c.ST = c.V[x]
		case 0x001E:
			c.I += uint16(c.V[x])
		case 0x0029:
			c.I = uint16(c.V[x] * 0x5)
		case 0x0033:
			c.Memory[c.I] = c.V[x] / 100
			c.Memory[c.I+1] = (c.V[x] / 10) % 10
			c.Memory[c.I+2] = (c.V[x] % 100) % 10
		case 0x0055:
			for i := uint16(0x0000); i <= x; i++ {
				c.Memory[c.I+i] = c.V[i]
			}
		case 0x0065:
			for i := uint16(0x0000); i <= x; i++ {
				c.V[i] = c.Memory[c.I+i]
			}
		}
	default:
		log.Fatalf("Unknown opcode: %x", op)
	}
	if c.DT > 0 {
		c.DT = c.DT - 1
	}
	if c.ST > 0 {
		c.ST = c.ST - 1
	}
	return nil
}

func (c *CPU) KeyRelease(key uint16) {
	if int(key) < len(c.Keys) {
		c.Keys[key] = false
	}
}

func (c *CPU) KeyPress(key uint16) {
	if int(key) < len(c.Keys) {
		c.Keys[key] = true
		if c.WKey != nil {
			*c.WKey = byte(key)
			c.WKey = nil
		}
	}
}
