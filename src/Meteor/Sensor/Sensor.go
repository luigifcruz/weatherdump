package Sensor

import (
	"encoding/binary"
	"fmt"
)

const sensorDataMinimum = 13

type Sensor struct {
	day  uint16
	msec uint32
	usec uint16
	MCUN uint8
	QT   uint8
	DC   uint8
	AC   uint8
	QFM  uint16
	QF   uint8
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
}
