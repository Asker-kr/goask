package main

import (
	"os"
	"fmt"
	"errors"
	"encoding/json"
	"time"
	"log"
	"flag"
	"path/filepath"
	"strings"
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

func readFile(filename string, c bool) {
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

	_, err = f.Seek(int64(listHeader.Size-12+8), os.SEEK_CUR)
	if nil != err {
		log.Printf("Can't move file pointer by listheader size, err(%s)\n", err)
		return
	}

	chunkHeader := Header{}
	if nil != chunkHeader.Get(f) {
		return
	}

	gprmcs := GetGPRMC(f, chunkHeader)
	if c {
		log.Printf("%d\t%s\n", len(gprmcs), filename)
	} else {
		if j, err := json.Marshal(gprmcs); err != nil {
			log.Printf("Can't make json output, err(%s)\n", err)
		} else {
			fmt.Printf("%s \n", j)
		}
	}
}

func visit(path string, f os.FileInfo, err error) error {
	if !f.IsDir() && filepath.Ext(path) == ".avi" && ( strings.Contains( filepath.Base(path), "E1.") || strings.Contains( filepath.Base(path), "E2.")) {
//		log.Printf("now file is %s\n", path)
		readFile(path, true)
	}
	return nil
}

func main() {
	c := flag.Bool("check", false, "check data counts")
	flag.Parse()

	date, mon, day := time.Now().Date()
	d := fmt.Sprintf("%04d%02d%02d", date, mon, day)
	fpLog, err := os.OpenFile("./log/"+"gprmc_"+d+".log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		panic(err)
	}
	defer fpLog.Close()
	log.SetOutput(fpLog)

	log.Printf("-------- Start GPRMC Parser. ------------\n")

	if len(flag.Args()) == 1 {
		log.Printf("target file is %s\n", flag.Args()[0])
		if *c {
			log.Printf("Check is true.\n")
			err := filepath.Walk(flag.Args()[0], visit)
			log.Printf("filepath.Walk() returned %v\n", err)
		} else {
			readFile(flag.Args()[0], false)
		}
	} else {
		log.Printf("target file name not exist.\n")
	}
	log.Printf("-------- End GPRMC Parser.   ------------\n")
}
