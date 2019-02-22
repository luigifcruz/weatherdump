package npoess

import (
	"encoding/binary"
	"fmt"
	"time"
)

// NASA's Timestamp Epoch is 1st January 1858
// Probably because Explorer 1 Launch Year (?)

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
	fmt.Printf("Days since 1958: %d\n", e.day)
	fmt.Printf("Milliseconds: %d\n", e.milliseconds)
	fmt.Printf("Microseconds: %d\n", e.microseconds)
	fmt.Printf("RFC3339: %s\n", e.GetZulu())
	fmt.Println()
}

func (e Time) GetZulu() string {
	return e.GetDate().UTC().Format(time.RFC3339)
}

func (e Time) GetDate() time.Time {
	// Start from epoch
	millis := int64(0)

	// Add spacecraft epoch count
	millis += int64(e.day) * 24 * 60 * 60 * 1000
	millis += int64(e.milliseconds)

	// Subtract days from January 1, 1958 to January 1, 1970
	millis -= 4383 * 24 * 60 * 60 * 1000

	// Convert to Normal Date
	nanos := millis * int64(time.Millisecond)
	return time.Unix(0, nanos)
}
