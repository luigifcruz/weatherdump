package Frames

import "fmt"

type MultiplexingFrame struct {
	firstHeaderPointer uint16
	packetZone []byte
}

func (e MultiplexingFrame) GetPacketZone() []byte {
	return e.packetZone
}

func (e MultiplexingFrame) GetFirstHeaderPointer() uint16 {
	return e.firstHeaderPointer
}

func (e *MultiplexingFrame) FromBinary(dat []byte) {
	e.firstHeaderPointer = (uint16(dat[0]) << 8 | uint16(dat[1])) & 0x7FF
	e.packetZone = dat[2:886]
}

func (e MultiplexingFrame) Print() {
	fmt.Println("### Multiplexing Frame Header")
	fmt.Printf("First Header Pointer: %011b\n", e.firstHeaderPointer)
	fmt.Println()
}

func (e MultiplexingFrame) HaveNewPackage() bool {
	return e.firstHeaderPointer != 2047
}