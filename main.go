package main

import (
	"os"
	"fmt"
	"errors"
	"encoding/json"
	"time"
	"log"
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
//		panic(err)
		log.Panicf("Can't Open file(%s), err(%s)\n", filename, err)
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
		log.Printf("Can't move file pointer by listheader size, err(%s)\n", err)
		return
	}

	chunkHeader := Header{}
	if nil != chunkHeader.Get(f) {
		return
	}

	gprmcs := GetGPRMC(f, chunkHeader)
	if j, err := json.Marshal(gprmcs); err != nil {
		log.Printf("Can't make json output, err(%s)\n", err)
	} else {
		fmt.Printf("%s \n", j)
	}
}

func main() {
	date, mon, day := time.Now().Date()
	d := fmt.Sprintf("%04d%02d%02d", date, mon, day )
	fpLog, err := os.OpenFile("/home/ub1st/gprmc/log/"+"gprmc_"+d+".log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		panic(err)
	}
	defer fpLog.Close()
	log.SetOutput(fpLog)

	log.Printf("-------- Start GPRMC Parser. ------------\n")
	if len(os.Args) == 2 {
		log.Printf("target file is %s\n", os.Args[1])
		readFile(os.Args[1])
	} else {
		log.Printf("invalid args count(%d). \n", len(os.Args))
	}
	log.Printf("-------- End GPRMC Parser.   ------------\n")
}
