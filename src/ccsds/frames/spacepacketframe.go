package frames

import (
	"encoding/binary"
	"fmt"
)

const spacePacketFrameMinimum = 6

// SpacePacketFrame data structure.
type SpacePacketFrame struct {
	versionNumber       uint8
	typeIndicator       uint8
	secondaryHeaderFlag uint8
	APID                uint16
	sequenceFlags       uint8
	packetSeqCount      uint16
	packetDataLength    uint16
	packetData          []byte

	dataLength int
}

// FromBinary parses the binary data into the dectector struct.
func (e *SpacePacketFrame) FromBinary(dat []byte) {
	if len(dat) < spacePacketFrameMinimum {
		return
	}

	e.dataLength += len(dat)
	e.versionNumber = dat[0] >> 5
	e.typeIndicator = (dat[0] & 0x1F) >> 4
	e.secondaryHeaderFlag = (dat[0] & 0x0F) >> 3
	e.APID = (uint16(dat[0]&0x07) << 8) | uint16(dat[1])
	e.sequenceFlags = uint8(binary.BigEndian.Uint16(dat[2:]) & 0xC000 >> 14)
	e.packetSeqCount = binary.BigEndian.Uint16(dat[2:]) & 0x3FFF
	e.packetDataLength = binary.BigEndian.Uint16(dat[4:])
	e.packetData = dat[6:]

	if uint16(len(e.packetData)) > e.packetDataLength+1 {
		e.packetData = e.packetData[:e.packetDataLength+1]
	}
}

// FeedData receives chunks of data and append it to the current data.
func (e *SpacePacketFrame) FeedData(dat []byte) []byte {
	currentData := (e.packetDataLength + 1)
	dataLeft := currentData - uint16(len(e.packetData))

	if currentData < 6 || dataLeft > uint16(len(dat)) {
		e.dataLength += len(dat)
		e.packetData = append(e.packetData, dat...)
		return nil
	}

	e.dataLength += int(dataLeft)
	e.packetData = append(e.packetData, dat[:dataLeft]...)
	return dat[dataLeft:]
}

// GetAPID return the current frame APID value.
func (e SpacePacketFrame) GetAPID() uint16 {
	return e.APID
}

// GetSequenceCount returns the current frame sequence count.
func (e SpacePacketFrame) GetSequenceCount() uint16 {
	return e.packetSeqCount
}

// GetPacketLength returns the current frame data length.
func (e SpacePacketFrame) GetPacketLength() uint16 {
	return e.packetDataLength
}

// GetData returns a slice containing the frame data.
func (e SpacePacketFrame) GetData() []byte {
	return e.packetData
}

// GetSequenceFlags returns the sequence flags of the current frame.
func (e SpacePacketFrame) GetSequenceFlags() uint8 {
	return e.sequenceFlags
}

// IsValid checks if the current packet is valid by comparing the data size.
// This is helpful to identify corrupted packets.
func (e SpacePacketFrame) IsValid() bool {
	return (e.packetDataLength + 1) == uint16(len(e.packetData))
}

// Print all exported variables from the current class into the terminal.
func (e SpacePacketFrame) Print() {
	fmt.Println("### Space Packet Primary Header")
	fmt.Printf("Version Number: %03b\n", e.versionNumber)
	fmt.Printf("Type Indicator: %01b\n", e.typeIndicator)
	fmt.Printf("Secondary Header: %01b\n", e.secondaryHeaderFlag)
	fmt.Printf("APID: %011b\n", e.APID)
	fmt.Printf("Sequence Flag: %02b\n", e.sequenceFlags)
	fmt.Printf("Sequence Count: %014b\n", e.packetSeqCount)
	fmt.Printf("Packet Length/Current Length: %d/%d\n", e.packetDataLength, len(e.packetData)-1)
	fmt.Println()
}
