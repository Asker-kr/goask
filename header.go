package main

import (
	"os"
	"bytes"
	"encoding/binary"
	"fmt"
)

type Header struct {
	SID   [4]byte
	Size  uint32
	Stype [4]byte
}

func (h *Header) Get(f *os.File) error {
	b := make([]byte, 12)
	readsize, err := f.Read(b)
	if ( len(b) != readsize || err != nil) {
		panic(err)
	}

	buf := bytes.NewBuffer(b)
	err = binary.Read(buf, binary.LittleEndian, h)
	if err != nil {
		fmt.Println("header Read failed:", err)
		return err
	}

	return nil
}
