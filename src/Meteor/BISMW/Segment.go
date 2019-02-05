package BISMW

import (
	"encoding/binary"
	"fmt"
	"reflect"
	"weather-dump/src/Meteor"
)

var eob = []int{-999}
var cfc = []int{-998}

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

	mcus  [14][]int
	final [14][64]byte
	valid bool
}

func NewSegment(buf []byte) *Segment {
	e := Segment{}
	e.FromBinary(buf)
	e.Parse()
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

func getValue(dat []bool) int {
	if len(dat) == 0 {
		fmt.Println("Got invalid value...")
		return 0
	}

	result := 0x00
	for i := len(dat) - 1; i > 0; i-- {
		if dat[i] {
			result = result | 0x0001<<uint(i-2)
		}
	}
	result += 0x01 << uint(len(dat)-1)
	if !dat[0] {
		result *= -1
	}
	return result
}

func findDC(dat *[]bool) int {
	buf := *dat
	for _, m := range dcCategories {
		klen := len(m.code)
		if len(buf) < klen {
			continue
		}

		if reflect.DeepEqual(buf[:klen], m.code) {
			if len(buf) < klen+m.len {
				break
			}
			*dat = buf[klen+m.len:]
			if m.len == 0 {
				return 0
			}
			return getValue(buf[klen : klen+m.len])
		}
	}
	*dat = nil
	return cfc[0]
}

func findAC(dat *[]bool) []int {
	buf := *dat
	for _, m := range acCategories {
		klen := len(m.code)
		if len(buf) < klen {
			continue
		}

		if reflect.DeepEqual(buf[:klen], m.code) {
			if m.clen == 0 && m.zlen == 0 {
				*dat = buf[klen:]
				return eob
			}
			vals := make([]int, m.zlen+1)
			if !(m.zlen == 15 && m.clen == 0) {
				if len(buf) < klen+m.clen {
					break
				}
				vals[m.zlen] = getValue(buf[klen : klen+m.clen])
			}
			*dat = buf[klen+m.clen:]
			return vals
		}
	}

	*dat = nil
	return cfc
}

func convertToArray(buf []byte) *[]bool {
	var soft = make([]bool, len(buf)*8)
	for i, m := range buf {
		soft[0+8*i] = m>>7&0x01 == 0x01
		soft[1+8*i] = m>>6&0x01 == 0x01
		soft[2+8*i] = m>>5&0x01 == 0x01
		soft[3+8*i] = m>>4&0x01 == 0x01
		soft[4+8*i] = m>>3&0x01 == 0x01
		soft[5+8*i] = m>>2&0x01 == 0x01
		soft[6+8*i] = m>>1&0x01 == 0x01
		soft[7+8*i] = m>>0&0x01 == 0x01
	}
	return &soft
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

var qTable = [64]int{
	16, 11, 10, 16, 24, 40, 51, 61,
	12, 12, 14, 19, 26, 58, 60, 55,
	14, 13, 16, 24, 40, 57, 69, 56,
	14, 17, 22, 29, 51, 87, 80, 62,
	18, 22, 37, 56, 68, 109, 103, 77,
	24, 35, 55, 64, 81, 104, 113, 92,
	49, 64, 78, 87, 103, 121, 120, 101,
	72, 92, 95, 98, 112, 100, 103, 99,
}

var zigzag = [64]int{
	0, 1, 5, 6, 14, 15, 27, 28,
	2, 4, 7, 13, 16, 26, 29, 42,
	3, 8, 12, 17, 25, 30, 41, 43,
	9, 11, 18, 24, 31, 40, 44, 53,
	10, 19, 23, 32, 39, 45, 52, 54,
	20, 22, 33, 38, 46, 51, 55, 60,
	21, 34, 37, 47, 50, 56, 59, 61,
	35, 36, 48, 49, 57, 58, 62, 63,
}

func (e *Segment) Parse() {
	buf := convertToArray(e.payload)
	for i := 0; i < 14; i++ {
		val := findDC(buf)
		if val == cfc[0] {
			fmt.Println("[JPEG] Invalid DC value, frame can't be restored.")
			return
		}

		if i == 0 {
			e.mcus[i] = []int{val}
		} else {
			e.mcus[i] = []int{val + e.mcus[i-1][0]}
		}

		for j := 0; j < 63; {
			vals := findAC(buf)
			j += len(vals)

			if vals[0] == cfc[0] {
				fmt.Println("[JPEG] Invalid AC value, frame can't be restored.")
				return
			}
			if vals[0] == eob[0] {
				//fmt.Printf("EOB! Chunks: %02d MCU#: %02d LEN: %08d DC: %d %d\n", j+1, i, len(*buf), e.mcus[i][0], val)
				break
			} else {
				e.mcus[i] = append(e.mcus[i], vals...)
			}
		}

		if len(e.mcus[i]) > 64 {
			fmt.Println("[JPEG] Invalid number of blocks.")
			return
		}

		e.mcus[i] = append(e.mcus[i], make([]int, 64-len(e.mcus[i]))...)
	}

	if len(*buf) > 16 {
		fmt.Println("[JPEG] Invalid number of remaining bits.")
	}

	for y := 0; y < 14; y++ {
		for x := 0; x < 64; x++ {
			e.mcus[y][x] = e.mcus[y][zigzag[x]]

			f := (200 - 2*float32(e.QF)) / 100.0
			if (e.QF > 20) && (e.QF < 50) {
				f = (5000 / float32(e.QF)) / 100.0
			}

			quantizationValue := int(float32(qTable[x]) * f)
			if quantizationValue != 0 {
				e.mcus[y][x] *= quantizationValue
			}
		}

		buf := [64]int{}
		copy(buf[:], e.mcus[y])
		Meteor.Idct(&buf)
		copy(e.mcus[y], buf[:])

		for x := 0; x < 64; x++ {
			e.final[y][x] = byte(e.mcus[y][x] + 128)
		}
	}
	/*
		name := fmt.Sprintf("./out_%d.jpeg", t)
		output, err := os.Create(name)
		if err != nil {
			fmt.Println(err)
			return
		}
		defer output.Close()

		img := [64 * 14]byte{}
		o := 0
		for y := 0; y < 8; y++ {
			for x := 0; x < 112; x++ {
				//fmt.Println(o, x/8, y*8+x-((x/8)*8))
				img[o] = e.final[x/8][y*8+x-((x/8)*8)]
				o++
			}
		}

		s := image.NewGray(image.Rect(0, 0, 112, 8))
		s.Pix = img[:]
		jpeg.Encode(output, s, nil)
	*/
}

func (e Segment) ExportSegment() [64 * 14]byte {
	img := [64 * 14]byte{}
	o := 0
	for y := 0; y < 8; y++ {
		for x := 0; x < 112; x++ {
			//fmt.Println(o, x/8, y*8+x-((x/8)*8))
			img[o] = e.final[x/8][y*8+x-((x/8)*8)]
			o++
		}
	}
	return img
}
