package bismw

import (
	"encoding/binary"
	"fmt"
	"weather-dump/src/protocols/lrpt"
)

const segmentDataMinimum = 13

type Segment struct {
	time  lrpt.Time
	MCUN  uint8
	QT    uint8
	DC    uint8
	AC    uint8
	QFM   uint16
	QF    uint8
	valid bool
	mcus  [14][]int64
}

func NewSegment(buf []byte) *Segment {
	e := Segment{}
	e.FromBinary(buf)
	return &e
}

func (e *Segment) FromBinary(dat []byte) {
	if len(dat) < segmentDataMinimum {
		return
	}

	e.time.FromBinary(dat[0:])
	e.MCUN = uint8(dat[8])
	e.QT = uint8(dat[9])
	e.DC = uint8(dat[10]) & 0xF0 >> 4
	e.AC = uint8(dat[10]) & 0x0F
	e.QFM = binary.BigEndian.Uint16(dat[11:])
	e.QF = uint8(dat[13])
	e.valid = true

	e.huffmanDecode(dat[14:])
}

func (e Segment) GetMCUNumber() uint8 {
	return e.MCUN
}

func (e Segment) GetDate() lrpt.Time {
	return e.time
}

func (e Segment) Print() {
	fmt.Println("### LRPT Segment Frame")
	fmt.Printf("MCU Number: %d\n", e.MCUN)
	fmt.Printf("Quantization Table: %08b\n", e.QT)
	fmt.Printf("Huffman Table DC: %04b\n", e.DC)
	fmt.Printf("Huffman Table AC: %04b\n", e.AC)
	fmt.Printf("Quality Factor Marker: %16b\n", e.QFM)
	fmt.Printf("Quality Factor: %08b\n", e.QF)
	fmt.Println()
	e.time.Print()
}

func (e *Segment) huffmanDecode(data []byte) {
	buf := convertToArray(data)
	lastDC := int64(0)

	for i := 0; i < 14; i++ {
		val := findDC(buf)
		if val == cfc[0] {
			e.valid = false
		}
		e.mcus[i] = []int64{val + lastDC}
		lastDC = e.mcus[i][0]

		for j := 0; j < 63; {
			vals := findAC(buf)
			j += len(vals)

			if vals[0] == cfc[0] {
				e.valid = false
			}
			if vals[0] != eob[0] {
				e.mcus[i] = append(e.mcus[i], vals...)
			} else {
				break
			}
		}

		if len(e.mcus[i]) > 64 {
			e.mcus[i] = e.mcus[i][:64]
		}

		e.mcus[i] = append(e.mcus[i], make([]int64, 64-len(e.mcus[i]))...)
	}
}

func (e Segment) RenderSegment(buf *[64 * 14]byte) {
	if !e.valid {
		return
	}

	quantizationTable := getQuantizationTable(float64(e.QF))
	output := [14][64]uint8{}

	for y := 0; y < 14; y++ {
		var buf [64]int64
		for x := 0; x < 64; x++ {
			buf[x] = e.mcus[y][zigzag[x]] * quantizationTable[x]
		}

		idct(&buf)

		for x := 0; x < 64; x++ {
			normalizedPixel := buf[x] + 128

			if normalizedPixel > 255 {
				normalizedPixel = 255
			}
			if normalizedPixel < 0 {
				normalizedPixel = 0
			}

			output[y][x] = uint8(normalizedPixel)
		}
	}

	o := 0
	for y := 0; y < 8; y++ {
		for x := 0; x < 112; x++ {
			(*buf)[o] = output[x/8][(y*8)+(x%8)]
			o++
		}
	}
}
