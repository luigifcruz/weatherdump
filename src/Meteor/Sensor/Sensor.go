package Sensor

import (
	"encoding/binary"
	"fmt"
	"reflect"
)

var eob = []int{-999}
var cfc = []int{-998}

const sensorDataMinimum = 13

type mcu struct {
	blocks [64]int
}

type Sensor struct {
	day     uint16
	msec    uint32
	usec    uint16
	MCUN    uint8
	QT      uint8
	DC      uint8
	AC      uint8
	QFM     uint16
	QF      uint8
	payload []byte
	units   [14]mcu
}

func NewSensor(buf []byte) *Sensor {
	e := Sensor{}
	e.FromBinary(buf)
	return &e
}

func (e *Sensor) FromBinary(dat []byte) {
	if len(dat) < sensorDataMinimum {
		return
	}

	e.day = binary.BigEndian.Uint16(dat[0:])
	e.msec = binary.BigEndian.Uint32(dat[2:])
	e.usec = binary.BigEndian.Uint16(dat[6:])

	e.MCUN = uint8(dat[8])
	e.QT = uint8(dat[9])
	e.DC = uint8(dat[10]) & 0xF0 >> 4
	e.AC = uint8(dat[10]) & 0x0F
	e.QFM = binary.BigEndian.Uint16(dat[11:])
	e.QF = uint8(dat[13])

	e.payload = dat[14:]
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

func (e Sensor) Print() {
	fmt.Println("### LRPT Sensor Frame")
	fmt.Printf("Days: %d\n", e.day)
	fmt.Printf("Milliseconds: %d\n", e.msec)
	fmt.Printf("Microseconds: %d\n", e.usec)

	fmt.Printf("First MCU Number: %08b\n", e.MCUN)
	fmt.Printf("Quantization Table: %08b\n", e.QT)
	fmt.Printf("Huffman Table DC: %04b\n", e.DC)
	fmt.Printf("Huffman Table AC: %04b\n", e.AC)
	fmt.Printf("Quality Factor Marker: %16b\n", e.QFM)
	fmt.Printf("Quality Factor: %08b\n", e.QF)
}

func (e *Sensor) Parse() {
	fmt.Printf("[JPEG] Packet size %d\n", len(e.payload))
	buf := convertToArray(e.payload)

	for mculen := 0; mculen < 14; mculen++ {
		val := findDC(buf)
		if val == cfc[0] {
			fmt.Println("[JPEG] Invalid DC value, frame can't be restored.")
			return
		}

		for chunks := 0; chunks < 62; {
			vals := findAC(buf)
			chunks += len(vals)

			if vals[0] == cfc[0] {
				fmt.Println("[JPEG] Invalid AC value, frame can't be restored.")
				return
			}
			if vals[0] == eob[0] {
				fmt.Printf("EOB! Chunks: %02d MCU#: %02d LEN: %08d DC: %d\n", chunks+1, mculen, len(*buf), val)
				break
			}
		}
	}

	fmt.Println(len(*buf))
	//os.Exit(0)
}
