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
	mcus    [14][]int
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

func (e *Segment) Parse() {
	fmt.Printf("[JPEG] Packet size %d\n", len(e.payload))
	buf := convertToArray(e.payload)

	for i := 0; i < 14; i++ {
		val := findDC(buf)
		if val == cfc[0] {
			fmt.Println("[JPEG] Invalid DC value, frame can't be restored.")
			return
		}

		e.mcus[i] = []int{val}

		for j := 0; j < 62; {
			vals := findAC(buf)
			j += len(vals)

			if vals[0] == cfc[0] {
				fmt.Println("[JPEG] Invalid AC value, frame can't be restored.")
				return
			}
			if vals[0] == eob[0] {
				//fmt.Printf("EOB! Chunks: %02d MCU#: %02d LEN: %08d DC: %d\n", j+1, i, len(*buf), val)
				break
			} else {
				e.mcus[i] = append(e.mcus[i], vals...)
			}
		}

		if len(e.mcus[i]) > 64 {
			fmt.Println("WTF = ", len(*buf))
			return
		}

		e.mcus[i] = append(e.mcus[i], make([]int, 64-len(e.mcus[i]))...)
	}

	fmt.Println(len(*buf))
	//os.Exit(0)
}
