package main

import (
	"os"
	"fmt"
	"errors"
	"encoding/json"
)

func getChunkData(f *os.File, bufSize uint32) (string, error) {
	b := make([]byte, bufSize)
	readsize, err := f.Read(b)
	if ( len(b) != readsize || err != nil) {
		panic(err)
	}

	if '$' != b[0] {
		return "", errors.New("not a gprmc data")
	}

	return string(b[1:readsize]), nil
}

func readFile(filename string) {
	f, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	riffHeader := Header{}
	if nil != riffHeader.Get(f) {
		return
	}

	listHeader := Header{}
	if nil != listHeader.Get(f) {
		return
	}
	_, err = f.Seek(int64(listHeader.Size - 12 + 8), os.SEEK_CUR)
	if nil != err {
		return
	}

	chunkHeader := Header{}
	if nil != chunkHeader.Get(f) {
		return
	}

	gprmcs := GetGPRMC(f, chunkHeader)
	j, _ := json.Marshal(gprmcs)
	fmt.Printf("%s \n", j)
}

func main() {
	if len(os.Args) == 2 {
		readFile(os.Args[1])
	}
}
