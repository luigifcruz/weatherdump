package parser

import (
	"fmt"
	"weather-dump/src/ccsds/frames"
	"weather-dump/src/protocols/hrd"
	"weather-dump/src/protocols/hrd/processor/parser/segment"
)

const maxFrameCount = 8192

type Channel struct {
	APID                  uint16
	ChannelName           string
	AggregationZoneWidth  [6]int
	AggregationZoneHeight int
	BowTieHeight          [6]int
	OversampleZone        [6]int
	Width                 uint32
	Height                uint32
	ReconstructionBand    uint16
	Invert                bool
	FileName              string
	FirstSegment          uint32
	LastSegment           uint32
	SegmentCount          uint32
	HasData               bool
	StartTime             hrd.Time
	EndTime               hrd.Time

	decoded    bool
	segments   map[uint32]*segment.Data
	scanCount  uint32
	exctdCount uint32
}

func (e *Channel) init() {
	e.segments = make(map[uint32]*segment.Data)

	e.HasData = true
	e.LastSegment = 0x00000000
	e.FirstSegment = 0xFFFFFFFF
	e.SegmentCount = 0x00000000
}

func (e Channel) GetDimensions() (int, int) {
	return int(e.Width), int(e.Height)
}

func (e Channel) GetBounds() (int, int) {
	return int(e.FirstSegment), int(e.LastSegment)
}

func (e *Channel) SetBounds(first, last int) {
	e.FirstSegment = uint32(first)
	e.LastSegment = uint32(last)
}

func (e *Channel) Fix(scft hrd.SpacecraftParameters) {
	if e.LastSegment-e.FirstSegment > maxFrameCount {
		fmt.Printf("[SEN] Potentially invalid channel %s was found.\n", e.ChannelName)
		fmt.Println("	It's too long for the round earth, trying to correct...")

		if (e.LastSegment - e.LastSegment - e.SegmentCount) < maxFrameCount {
			e.FirstSegment = e.LastSegment - e.SegmentCount
		}

		if (e.FirstSegment + e.SegmentCount - e.FirstSegment) < maxFrameCount {
			e.LastSegment = e.FirstSegment + e.SegmentCount
		}

		if e.LastSegment-e.FirstSegment > maxFrameCount {
			fmt.Println("	Cannot find any valid number, skipping channel.")
			return
		}

		fmt.Println("	Found a valid number. Channel can still be damaged.")
	}

	for i := e.LastSegment; i >= e.FirstSegment; i-- {
		if e.segments[i] == nil {
			e.segments[i] = segment.NewFillSegment(i)
		}
	}

	e.StartTime = e.segments[e.FirstSegment].Header.GetDate()
	e.EndTime = e.segments[e.LastSegment].Header.GetDate()
	e.FileName = fmt.Sprintf("%s_%s_VIIRS_%s_%s", scft.Filename, scft.SignalName, e.ChannelName, e.StartTime.GetZuluSafe())
	e.Height = (e.LastSegment - e.FirstSegment + 1) * uint32(e.AggregationZoneHeight)
}

func (e *Channel) Parse(packet frames.SpacePacketFrame) {
	if packet.GetSequenceFlags() == 1 && packet.IsValid() {
		if !e.HasData {
			e.init()
		}

		frameHeader := segment.NewFrameHeader(packet.GetData())
		e.scanCount = frameHeader.GetScanNumber()
		e.exctdCount = frameHeader.GetSequenceCount() + uint32(frameHeader.GetNumberOfSegments()) + 2
		e.segments[e.scanCount] = segment.NewSegment(frameHeader)

		if e.LastSegment < e.scanCount {
			e.LastSegment = e.scanCount
		}

		if e.FirstSegment > e.scanCount {
			e.FirstSegment = e.scanCount
		}

		e.SegmentCount++
		return
	}

	if e.HasData {
		body := segment.NewBody(packet.GetData())
		if body.GetSequenceCount() <= e.exctdCount && body.GetDetectorNumber() < 32 {
			e.segments[e.scanCount].Body[body.GetDetectorNumber()] = *body
		}
	}
}
