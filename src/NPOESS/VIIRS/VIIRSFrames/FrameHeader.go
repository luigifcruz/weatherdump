package VIIRSFrames

import (
	"encoding/binary"
	"fmt"
	"weather-dump/src/NPOESS"
)

const frameHeaderMinimum = 52

type FrameHeader struct {
	time             NPOESS.Time
	numberOfSegments uint8
	sequenceCount    uint32
	packetTime       NPOESS.Time
	formatVersion    uint8
	instrumentNumber uint8
	hamSide          uint8
	scanSync         uint8
	selfTestPattern  uint8
	scanNumber       uint32
	scanTerminus     NPOESS.Time
	sensorMode       uint8
	viirsModel       uint8
	fswVersion       uint16
	bandControlWorld uint32
	partialStart     uint16
	numberOfSamples  uint16
	sampleDelay      uint16
	fillFrame        bool
}

func NewFillFrameHeader(scanNumber uint32) *FrameHeader {
	e := FrameHeader{}
	e.scanNumber = scanNumber
	e.fillFrame = true
	return &e
}

func NewFrameHeader(buf []byte) *FrameHeader {
	e := FrameHeader{}
	e.FromBinary(buf)
	return &e
}

func (e *FrameHeader) FromBinary(dat []byte) {
	if len(dat) < frameHeaderMinimum {
		return
	}

	e.time.FromBinary(dat[0:8])
	e.numberOfSegments = dat[8]
	// Spare 8 bits
	e.sequenceCount = binary.BigEndian.Uint32(dat[10:])
	e.packetTime.FromBinary(dat[14:22])
	e.formatVersion = uint8(dat[22])
	e.instrumentNumber = uint8(dat[23])
	// Spare 16 bits
	e.hamSide = uint8(dat[26]) >> 7
	e.scanSync = uint8(dat[26]&0x40) >> 6
	e.selfTestPattern = uint8(dat[26]&0x3C) >> 2
	// Spare 10 bits
	e.scanNumber = binary.BigEndian.Uint32(dat[28:])
	e.scanTerminus.FromBinary(dat[32:40])
	e.sensorMode = uint8(dat[40])
	e.viirsModel = uint8(dat[41])
	e.fswVersion = binary.BigEndian.Uint16(dat[42:])
	e.bandControlWorld = binary.BigEndian.Uint32(dat[44:])
	e.partialStart = binary.BigEndian.Uint16(dat[48:])
	e.numberOfSamples = binary.BigEndian.Uint16(dat[50:])
	e.sampleDelay = binary.BigEndian.Uint16(dat[52:])
	e.fillFrame = false
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

// Struct Validation
func (e FrameHeader) IsValid() bool {
	return !e.fillFrame
}

// Struct Get
func (e FrameHeader) GetDateString() string {
	return e.time.GetZulu()
}

func (e FrameHeader) GetDate() NPOESS.Time {
	return e.time
}

func (e FrameHeader) GetNumberOfSegments() uint8 {
	return e.numberOfSegments
}

func (e FrameHeader) GetSequenceCount() uint32 {
	return e.sequenceCount
}

func (e FrameHeader) GetScanNumber() uint32 {
	return e.scanNumber
}

// Struct Set
