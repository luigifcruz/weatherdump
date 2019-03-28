package frames

import (
	"encoding/binary"
	"fmt"
)

const frameSize = 892
const transferFrameMinimum = frameSize

// TransferFrame data structure.
type TransferFrame struct {
	versionNumber       uint8
	SCID                uint8
	VCID                uint8
	virtualChannelCount uint32
	replayFlag          uint8
	MPDU                []byte
}

// NewTransferFrame returns a new TransferFrame pointer
// populated with the binary data passed to it.
func NewTransferFrame(dat []byte) *TransferFrame {
	e := TransferFrame{}
	e.FromBinary(dat)
	return &e
}

// FromBinary parses the binary data into the dectector struct.
func (e *TransferFrame) FromBinary(dat []byte) {
	if len(dat) < transferFrameMinimum {
		return
	}

	e.versionNumber = dat[0] >> 6
	e.SCID = (dat[0]&0x3F)<<2 | (dat[1]&0xC0)>>6
	e.VCID = dat[1] & 0x3F
	e.virtualChannelCount = binary.BigEndian.Uint32(dat[2:]) >> 8
	e.replayFlag = dat[5] >> 7
	e.MPDU = dat[6:892]
}

// IsReplay returns if the current frame is replay.
func (e TransferFrame) IsReplay() bool {
	return e.replayFlag == 0x01
}

// GetMPDU returns the MPDU of the current frame.
func (e TransferFrame) GetMPDU() []byte {
	return e.MPDU
}

// GetVCID returns the VCID of the current frame.
func (e TransferFrame) GetVCID() uint8 {
	return e.VCID
}

// GetSCID returns the SCID of the current frame.
func (e TransferFrame) GetSCID() uint8 {
	return e.SCID
}

// Print all exported variables from the current class into the terminal.
func (e TransferFrame) Print() {
	fmt.Println("### Transfer Frame Primary Header")
	fmt.Printf("Version Number: %02b\n", e.versionNumber)
	fmt.Printf("Spacecraft ID: %08b\n", e.SCID)
	fmt.Printf("Virtual Channel ID: %06b\n", e.VCID)
	fmt.Printf("Virtual Channel Count: %024b\n", e.virtualChannelCount)
	fmt.Printf("Replay Flag: %01b\n", e.replayFlag)
	fmt.Println()
}
