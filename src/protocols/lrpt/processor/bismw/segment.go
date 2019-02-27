package bismw

import (
	"encoding/binary"
	"fmt"
	"weather-dump/src/protocols/lrpt"
)

const segmentDataMinimum = 13

type Segment struct {
	time    lrpt.Time
	MCUN    uint8
	QT      uint8
	DC      uint8
	AC      uint8
	QFM     uint16
	QF      uint8
	payload []byte
	mcus    [14][]float64
	export  [14][64]byte
}

func NewSegment(buf []byte) *Segment {
	e := Segment{}
	e.FromBinary(buf)
	valid := e.huffmanDecode()
	if valid {
		e.dequantize()
	}
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

func (e Segment) RenderSegment(buf *[64 * 14]byte) {
	o := 0
	for y := 0; y < 8; y++ {
		for x := 0; x < 112; x++ {
			(*buf)[o] = e.export[x/8][(y*8)+(x%8)]
			o++
		}
	}
}

func (e *Segment) huffmanDecode() bool {
	buf := convertToArray(e.payload)
	lastDC := 0.0

	for i := 0; i < 14; i++ {
		val := findDC(buf)
		if val == cfc[0] {
			//fmt.Println("[JPEG] Invalid DC value, frame can't be restored.")
			return false
		}
		e.mcus[i] = []float64{val + lastDC}
		lastDC = e.mcus[i][0]

		for j := 0; j < 63; {
			vals := findAC(buf)
			j += len(vals)

			if vals[0] == cfc[0] {
				//fmt.Println("[JPEG] Invalid AC value, frame can't be restored.")
				return false
			}
			if vals[0] != eob[0] {
				e.mcus[i] = append(e.mcus[i], vals...)
			} else {
				//fmt.Printf("EOB! Chunks: %02d MCU#: %02d LEN: %08d DC: %02f %02f\n", j+1, i, len(*buf), e.mcus[i][0], val)
				break
			}
		}

		if len(e.mcus[i]) > 64 {
			//fmt.Println("[JPEG] Invalid number of blocks. Cropping...")
			e.mcus[i] = e.mcus[i][:64]
		}

		e.mcus[i] = append(e.mcus[i], make([]float64, 64-len(e.mcus[i]))...)
	}

	return true
}

func (e *Segment) dequantize() {
	quantizationTable := getQuantizationTable(float64(e.QF))

	for y := 0; y < 14; y++ {
		var buf [64]int64
		for x := 0; x < 64; x++ {
			buf[x] = int64((e.mcus[y][zigzag[x]] * float64(quantizationTable[x])) + 0.5)
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

			e.export[y][x] = uint8(normalizedPixel)
		}
		e.mcus[y] = nil
	}
}
