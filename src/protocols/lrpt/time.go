package lrpt

import (
	"encoding/binary"
	"fmt"
)

type Time struct {
	day          uint16
	milliseconds uint32
	microseconds uint16
}

// FromBinary parses the binary data into the dectector struct.
func (e *Time) FromBinary(dat []byte) {
	e.day = binary.BigEndian.Uint16(dat[0:])
	e.milliseconds = binary.BigEndian.Uint32(dat[2:])
	e.microseconds = binary.BigEndian.Uint16(dat[6:])
}

// Print all exported variables from the current class into the terminal.
func (e Time) Print() {
	fmt.Println("### Time Frame Segment")
	fmt.Printf("Day: %d\n", e.day)
	fmt.Printf("Milliseconds: %d\n", e.milliseconds)
	fmt.Printf("Microseconds: %d\n", e.microseconds)
	fmt.Println()
}

// IsValid checks if the current time is valid.
// This is helpful to identify corrupted segments.
func (e Time) IsValid() bool {
	return e.day == 0 && e.microseconds == 0
}

func (e Time) GetZuluSafe() string {
	return string(e.GetMilliseconds())
}

func (e Time) GetMilliseconds() uint32 {
	return e.milliseconds
}
