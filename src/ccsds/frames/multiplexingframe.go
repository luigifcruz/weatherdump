package frames

import (
	"encoding/binary"
	"fmt"
	"weatherdump/src/ccsds/parameters"
)

const multiplexingFrameMinimum = 2

// MultiplexingFrame data structure.
type MultiplexingFrame struct {
	firstHeaderPointer uint16
	packetZone         []byte
	CCSDS              int
}

// NewMultiplexingFrame returns a new MultiplexingFrame pointer
// populated with the binary data passed to it.
func NewMultiplexingFrame(version int, dat []byte) *MultiplexingFrame {
	e := MultiplexingFrame{}
	e.CCSDS = version
	e.FromBinary(dat)
	return &e
}

// FromBinary parses the binary data into the dectector struct.
func (e *MultiplexingFrame) FromBinary(dat []byte) {
	if len(dat) < multiplexingFrameMinimum {
		return
	}

	switch e.CCSDS {
	case parameters.Version["LRPT"]:
		e.firstHeaderPointer = binary.BigEndian.Uint16(dat[2:]) & 0x7FF
		e.packetZone = dat[4:]
	case parameters.Version["HRD"]:
		e.firstHeaderPointer = binary.BigEndian.Uint16(dat[0:]) & 0x7FF
		e.packetZone = dat[2:]
	}
}

// Print all exported variables from the current class into the terminal.
func (e MultiplexingFrame) Print() {
	fmt.Println("### Multiplexing Frame Header")
	fmt.Printf("First Header Pointer: %011b\n", e.firstHeaderPointer)
	fmt.Println()
}

// IsValid checks if the current frame is valid by comparing the data size.
// This is helpful to identify corrupted packets.
func (e MultiplexingFrame) IsValid() bool {
	switch e.CCSDS {
	case parameters.Version["LRPT"]:
		return len(e.packetZone) == (886 - 4)
	case parameters.Version["HRD"]:
		return len(e.packetZone) == (886 - 2)
	}
	return false
}

// GetPacketZone returns the current packet zone.
func (e MultiplexingFrame) GetPacketZone() []byte {
	return e.packetZone
}

// GetFirstHeaderPointer returns the current FHP value.
func (e MultiplexingFrame) GetFirstHeaderPointer() uint16 {
	return e.firstHeaderPointer
}

// HaveNewPackage indicates if the frame contains a new package.
func (e MultiplexingFrame) HaveNewPackage() bool {
	return e.firstHeaderPointer != 2047
}
