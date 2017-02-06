package main

import (
	"os"
	"bytes"
	"encoding/binary"
	"log"
	"errors"
)

type InnerChunk struct {
	SID  [4]byte
	Size uint32
}

func (c *InnerChunk) Get(f *os.File) error{
	b := make([]byte, 8)
	readSize, err := f.Read(b)
	if err != nil {
		log.Printf("Can't read header data readsize(%d), err(%s)\n", readSize, err)
		return err
	}

	if len(b) != readSize  {
		log.Printf("Read header size invalid readsize(%d)\n", readSize)
		return errors.New("Invalid read Size.")
	}

	buf := bytes.NewBuffer(b)
	err = binary.Read(buf, binary.LittleEndian, c)
	if err != nil {
		log.Printf("Can't conver chunk data to binary, err(%s)\n", err)
		return err
	}

	return nil
}

