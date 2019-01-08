package Frames

import (
	"encoding/binary"
	"fmt"
)

const frameSize = 892

type TransferFrame struct {
	versionNumber       uint8
	SCID                uint8
	VCID                uint8
	virtualChannelCount uint32
	replayFlag          uint8
	MPDU                []byte
}

func (e *TransferFrame) FromBinary(dat []byte) {
	e.versionNumber = dat[0] >> 6
	e.SCID = (dat[0]&0x3F)<<2 | (dat[1]&0xC0)>>6
	e.VCID = dat[1] & 0x3F
	e.virtualChannelCount = binary.BigEndian.Uint32(dat[2:]) >> 8
	e.replayFlag = dat[5] >> 7
	e.MPDU = dat[6:892]
}

func (e TransferFrame) GetMPDU() []byte {
	return e.MPDU
}

func (e TransferFrame) GetVCID() uint8 {
	return e.VCID
}

func (e TransferFrame) GetSCID() uint8 {
	return e.SCID
}

func (e TransferFrame) Print() {
	fmt.Println("### Transfer Frame Primary Header")
	fmt.Printf("Version Number: %02b\n", e.versionNumber)
	fmt.Printf("Spacecraft ID: %08b\n", e.SCID)
	fmt.Printf("Virtual Channel ID: %06b\n", e.VCID)
	fmt.Printf("Virtual Channel Count: %024b\n", e.virtualChannelCount)
	fmt.Printf("Replay Flag: %01b\n", e.replayFlag)
	fmt.Println()
}
