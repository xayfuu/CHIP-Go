package chip8

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

type ROM struct {
	Size         int64
	Name         string
	Reader       *os.File
	Filepath     string
	InitialState []byte
}

func ReadROM(fp string) (*ROM, error) {
	r := ROM{
		Name:     filepath.Base(fp),
		Filepath: fp,
	}

	fInfo, err := os.Stat(fp)
	if err != nil {
		return nil, err
	}
	r.Size = fInfo.Size()

	if r.Reader, err = os.Open(fp); err != nil {
		return nil, err
	}

	return &r, nil
}

func (r *ROM) CheckMemoryOverflow(memSize int64) error {
	if r.Size > memSize {
		return errors.New(fmt.Sprintf("Program too large to fit in memory. (%d > %d)", r.Size, memSize))
	}
	return nil
}
