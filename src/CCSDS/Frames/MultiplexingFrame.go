package Frames

import (
	"encoding/binary"
	"fmt"
	"weather-dump/src/CCSDS/Parameters"
)

const multiplexingFrameMinimum = 2

type MultiplexingFrame struct {
	firstHeaderPointer uint16
	packetZone         []byte
	CCSDS              int
}

func NewMultiplexingFrame(version int, dat []byte) *MultiplexingFrame {
	e := MultiplexingFrame{}
	e.CCSDS = version
	e.FromBinary(dat)
	return &e
}

func (e *MultiplexingFrame) FromBinary(dat []byte) {
	if len(dat) < multiplexingFrameMinimum {
		return
	}

	switch e.CCSDS {
	case Parameters.Version["LRPT"]:
		e.firstHeaderPointer = binary.BigEndian.Uint16(dat[2:]) & 0x7FF
		e.packetZone = dat[4:]
	case Parameters.Version["HRD"]:
		e.firstHeaderPointer = binary.BigEndian.Uint16(dat[0:]) & 0x7FF
		e.packetZone = dat[2:]
	}
}

func (e MultiplexingFrame) Print() {
	fmt.Println("### Multiplexing Frame Header")
	fmt.Printf("First Header Pointer: %011b\n", e.firstHeaderPointer)
	fmt.Println()
}

func (e MultiplexingFrame) IsValid() bool {
	switch e.CCSDS {
	case Parameters.Version["LRPT"]:
		return len(e.packetZone) == (886 - 4)
	case Parameters.Version["HRD"]:
		return len(e.packetZone) == (886 - 2)
	}
	return false
}

func (e MultiplexingFrame) GetPacketZone() []byte {
	return e.packetZone
}

func (e MultiplexingFrame) GetFirstHeaderPointer() uint16 {
	return e.firstHeaderPointer
}

func (e MultiplexingFrame) HaveNewPackage() bool {
	return e.firstHeaderPointer != 2047
}
