package BISMW

import (
	"encoding/binary"
	"fmt"
	"weather-dump/src/Meteor"
)

const segmentDataMinimum = 13

type Segment struct {
	time    Meteor.Time
	MCUN    uint8
	QT      uint8
	DC      uint8
	AC      uint8
	QFM     uint16
	QF      uint8
	payload []byte
	mcus    [14][64]float64
	export  [14][64]byte
}

func NewSegment(buf []byte) *Segment {
	e := Segment{}
	e.FromBinary(buf)
	e.HuffmanDecode()
	e.Dequantize()
	e.RenderSegment()
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

	e.payload = dat[14:]
}

func (e Segment) GetMCUNumber() uint8 {
	return e.MCUN
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

func (e *Segment) HuffmanDecode() {
	buf := convertToArray(e.payload)

	for i := 0; i < 14; i++ {
		val := findDC(buf)
		if val == cfc[0] {
			fmt.Println("[JPEG] Invalid DC value, frame can't be restored.")
			return
		}

		tmp := []float64{val}
		if i != 0 {
			tmp[0] += e.mcus[i-1][0]
		}

		for j := 0; j < 63; {
			vals := findAC(buf)
			j += len(vals)

			if vals[0] == cfc[0] {
				fmt.Println("[JPEG] Invalid AC value, frame can't be restored.")
				return
			}
			if vals[0] == eob[0] {
				//fmt.Printf("EOB! Chunks: %02d MCU#: %02d LEN: %08d DC: %d %d\n", j+1, i, len(*buf), tmp[0], val)
				break
			} else {
				tmp = append(tmp, vals...)
			}
		}

		if len(tmp) > 64 {
			fmt.Println("[JPEG] Invalid number of blocks.")
			return
		}

		tmp = append(tmp, make([]float64, 64-len(tmp))...)
		copy(e.mcus[i][:], tmp[:])
	}
}

func (e *Segment) Dequantize() {
	quantizationTable := getQuantizationTable(float64(e.QF))

	for y := 0; y < 14; y++ {
		var input, output [64]float64
		for x := 0; x < 64; x++ {
			input[x] = e.mcus[y][zigzag[x]] * quantizationTable[x]
		}

		calculateIdct(&output, &input)

		for x := 0; x < 64; x++ {
			normalizedPixel := output[x] + 128
			e.export[y][x] = uint8(normalizedPixel + 0.5)

			if normalizedPixel > 255 {
				e.export[y][x] = 255
			}
			if normalizedPixel < 0 {
				e.export[y][x] = 0
			}
		}
	}
}

func (e Segment) RenderSegment() []byte {
	var buf = make([]byte, 64*14)
	o := 0
	for y := 0; y < 8; y++ {
		for x := 0; x < 112; x++ {
			buf[o] = byte(e.export[x/8][y*8+x-(x/8*8)])
			o++
		}
	}

	return buf
}

func (e Segment) GetDate() Meteor.Time {
	return e.time
}
