package ScienceFrames

import (
	"osp-noaa-dump/VIIRS/Common"
	"fmt"
)

type FrameHeader struct {
	time VIIRS.Time
	numberOfSegments uint8

	sequenceCount uint32
	packetTime VIIRS.Time
	formatVersion uint8
	instrumentNumber uint8

	hamSide uint8
	scanSync uint8
	selfTestPattern uint8

	scanNumber uint32
	scanTerminus VIIRS.Time
	sensorMode uint8
	viirsModel uint8
	fswVersion uint16
	bandControlWorld uint32
	partialStart uint16
	numberOfSamples uint16
	sampleDelay uint16
}

func NewHeader() *FrameHeader {
	return &FrameHeader{}
}

func (e FrameHeader) GetDate() string {
	return e.time.GetZulu()
}

func (e FrameHeader) GetNumberOfSegments() uint8 {
	return e.numberOfSegments
}

func (e FrameHeader) GetSequenceCount() uint32 {
	return e.sequenceCount
}

func (e *FrameHeader) FromBinary(dat []byte) {
	e.time.FromBinary(dat[0:8])
	e.numberOfSegments = dat[8]
	// Spare 8 bits
	e.sequenceCount = uint32(dat[10]) << 24 | uint32(dat[11]) << 16 | uint32(dat[12]) << 8 | uint32(dat[13])
	e.packetTime.FromBinary(dat[14:22])
	e.formatVersion = uint8(dat[22])
	e.instrumentNumber = uint8(dat[23])
	// Spare 16 bits
	e.hamSide = uint8(dat[26]) >> 7
	e.scanSync = uint8(dat[26] & 0x40) >> 6
	e.selfTestPattern = uint8(dat[26] & 0x3C) >> 2
	// Spare 10 bits
	e.scanNumber = uint32(dat[28]) << 24 | uint32(dat[29]) << 16 | uint32(dat[30]) << 8 | uint32(dat[31])
	e.scanTerminus.FromBinary(dat[32:40])
	e.sensorMode = uint8(dat[40])
	e.viirsModel = uint8(dat[41])
	e.fswVersion = uint16(dat[42]) << 8 | uint16(dat[43])
	e.bandControlWorld = uint32(dat[44]) << 24 | uint32(dat[45]) << 16 | uint32(dat[46]) << 8 | uint32(dat[47])
	e.partialStart = uint16(dat[48]) << 8 | uint16(dat[49])
	e.numberOfSamples = uint16(dat[50]) << 8 | uint16(dat[51])
	e.sampleDelay = uint16(dat[52]) << 8 | uint16(dat[53])
}

func (e FrameHeader) Print() {
	fmt.Println("### VIIRS Science Header")
	fmt.Printf("Day Time: %s\n", e.time.GetZulu())
	fmt.Printf("Number of Segments %08b\n", e.numberOfSegments)
	fmt.Println()
	fmt.Printf("Sequence Count: %032b\n", e.sequenceCount)
	fmt.Printf("Packet Time: %s\n", e.packetTime.GetZulu())
	fmt.Printf("Format Version: %08b\n", e.formatVersion)
	fmt.Printf("Instrument Number: %08b\n", e.instrumentNumber)
	fmt.Println()
	fmt.Printf("HAM Side: %01b\n", e.hamSide)
	fmt.Printf("Scan Synch: %01b\n", e.scanSync)
	fmt.Printf("Self Test Data Patter: %04b\n", e.selfTestPattern)
	fmt.Printf("Scan Number: %032b\n", e.scanNumber)
	fmt.Printf("Scan Terminus: %s\n", e.scanTerminus.GetZulu())
	fmt.Printf("Sensor Mode: %08b\n", e.sensorMode)
	fmt.Printf("VIIRS Model: %08b\n", e.viirsModel)
	fmt.Printf("FSW Version: %016b\n", e.fswVersion)
	fmt.Printf("Band Controll Word: %032b\n", e.bandControlWorld)
	fmt.Printf("Partial Start: %016b\n", e.partialStart)
	fmt.Printf("Number of Samples: %016b\n", e.numberOfSamples)
	fmt.Printf("Sample Delay: %016b\n", e.sampleDelay)
	fmt.Println()
}