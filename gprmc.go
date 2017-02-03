package main

import (
	"strings"
	"strconv"
	"os"
	"reflect"
	"math"
)

type Gprmc struct {
	Idx      int     `json:"seq"`
	DateTime string  `json:"time"`
	Valid    bool    `json:"valid"`
	Y        float64 `json:"wgs84_y"`
	X        float64 `json:"wgs84_x"`
	Speed    float64 `json:"speed"`
	Dir      float64 `json:"direction"`
}

func makeDateTime(tockens []string) string {
	t := tockens[1][0:6]
	d := "20" + tockens[9][4:6] + tockens[9][0:2] + tockens[9][2:4]
	return d + t
}

func (g *Gprmc) SetData( idx int, s string) {
	g.Idx = idx

	tockens := strings.Split(s, ",")
	g.DateTime = makeDateTime(tockens)
	g.Valid = tockens[2] == "A"

	if y, err := strconv.ParseFloat(tockens[3], 64); err == nil {
		y = math.Floor(y/100) + (y - math.Floor(y/100)*100)/60
		g.Y = y
	}

	if x, err := strconv.ParseFloat(tockens[5], 64); err == nil {
		x = math.Floor(x/100) + (x - math.Floor(x/100)*100)/60
		g.X = x
	}

	if sp, err := strconv.ParseFloat(tockens[7], 64); err == nil {
		g.Speed = sp * 1.852
	}

	if dir, err := strconv.ParseFloat(tockens[8], 64); err == nil {
		g.Dir = dir
	}
}

func GetGPRMC(f *os.File, h Header) []Gprmc {
	gprmcs := []Gprmc{}

	readed := uint32(0)
	headerSize := uint32(reflect.TypeOf(h).Size())
	i := 0
	for n := uint32(0); n < h.Size-(headerSize-8); n += readed {
		chunk := InnerChunk{}
		if nil != chunk.Get(f) {
			return gprmcs
		}

		if 0 != chunk.Size%2 {
			chunk.Size++
		}

		readed = 8 + chunk.Size

		if 0 == chunk.Size {
			continue
		}

		if chunk.SID[2] == 't' && chunk.SID[3] == 'x' {
			s, err := getChunkData(f, chunk.Size)
			if nil != err {
				continue
			}

			i++
			g := Gprmc{}
			g.SetData( i, s)
			gprmcs = append(gprmcs, g)
		} else {
			_, err := f.Seek(int64(chunk.Size), os.SEEK_CUR)
			if nil != err {
				return gprmcs
			}
		}
	}

	return gprmcs
}
