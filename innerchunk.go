package main

import (
	"os"
	"bytes"
	"encoding/binary"
	"fmt"
)

type InnerChunk struct {
	SID  [4]byte
	Size uint32
}

func (c *InnerChunk) Get(f *os.File) error{
	b := make([]byte, 8)
	readsize, err := f.Read(b)
	if ( len(b) != readsize || err != nil) {
		panic(err)
	}

	buf := bytes.NewBuffer(b)
	err = binary.Read(buf, binary.LittleEndian, c)
	if err != nil {
		fmt.Println("header Read failed:", err)
		return err
	}

	return nil
}

