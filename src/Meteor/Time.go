package Meteor

import (
	"encoding/binary"
	"fmt"
)

type Time struct {
	day          uint16
	milliseconds uint32
	microseconds uint16
}

func (e *Time) FromBinary(dat []byte) {
	e.day = binary.BigEndian.Uint16(dat[0:])
	e.milliseconds = binary.BigEndian.Uint32(dat[2:])
	e.microseconds = binary.BigEndian.Uint16(dat[6:])
}

func (e Time) Print() {
	fmt.Println("### Time Frame Segment")
	fmt.Printf("Day: %d\n", e.day)
	fmt.Printf("Milliseconds: %d\n", e.milliseconds)
	fmt.Printf("Microseconds: %d\n", e.microseconds)
	fmt.Println()
}

func (e Time) GetMilliseconds() uint32 {
	return e.milliseconds
}
