package Sensor

import (
	"encoding/binary"
	"fmt"
	"os"
	"reflect"
)

const sensorDataMinimum = 13

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

type Mask struct {
	code []bool
	len  int
}

var dcCategories = [12]Mask{
	Mask{[]bool{false, false}, 3},                                           // 00
	Mask{[]bool{true, false, true}, 4},                                      // 101
	Mask{[]bool{true, true, false}, 5},                                      // 110
	Mask{[]bool{true, false, false}, 2},                                     // 100
	Mask{[]bool{false, true, true}, 1},                                      // 011
	Mask{[]bool{false, true, false}, 0},                                     // 010
	Mask{[]bool{true, true, true, false}, 6},                                // 1110
	Mask{[]bool{true, true, true, true, false}, 7},                          // 11110
	Mask{[]bool{true, true, true, true, true, false}, 8},                    // 111110
	Mask{[]bool{true, true, true, true, true, true, false}, 9},              // 1111110
	Mask{[]bool{true, true, true, true, true, true, true, false}, 10},       // 11111110
	Mask{[]bool{true, true, true, true, true, true, true, true, false}, 11}, // 111111110
}

func findCategory(dat []bool) []bool {
	for _, m := range dcCategories {
		if reflect.DeepEqual(dat[:len(m.code)], m.code) {
			fmt.Println(m.len)
			return dat[len(m.code)+m.len:]
		}
	}
	return nil
}

func convertToArray(buf []byte) []bool {
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
	return soft
}

func (e Sensor) Print() {
	fmt.Println("# LRPT Sensor Frame")
	fmt.Printf("Days: %d\n", e.day)
	fmt.Printf("Milliseconds: %d\n", e.msec)
	fmt.Printf("Microseconds: %d\n", e.usec)

	fmt.Printf("First MCU Number: %08b\n", e.MCUN)
	fmt.Printf("Quantization Table: %08b\n", e.QT)
	fmt.Printf("Huffman Table DC: %04b\n", e.DC)
	fmt.Printf("Huffman Table AC: %04b\n", e.AC)
	fmt.Printf("Quality Factor Marker: %16b\n", e.QFM)
	fmt.Printf("Quality Factor: %08b\n", e.QF)
	fmt.Println()

	fmt.Printf("%08b %08b\n", e.payload[0], e.payload[1])

	g := convertToArray(e.payload[0:2])
	g = findCategory(g)

	os.Exit(0)
}
