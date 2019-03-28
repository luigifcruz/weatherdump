package parser

import (
	"fmt"
	"weather-dump/src/ccsds/frames"
	"weather-dump/src/protocols/lrpt"
	"weather-dump/src/protocols/lrpt/processor/parser/segment"
)

const maxFrameCount = 8192 * 3

// Channel struct.
type Channel struct {
	APID         uint16
	ChannelName  string
	BlockDim     int
	Invert       bool
	FinalWidth   uint32
	FileName     string
	Height       uint32
	Width        uint32
	StartTime    lrpt.Time
	EndTime      lrpt.Time
	HasData      bool
	SegmentCount uint32

	FirstSegment uint32
	LastSegment  uint32

	segments map[uint32]*segment.Data
	rollover uint32
	last     uint32
	offset   uint32
}

// NewChannel instance.
func (e *Channel) init() {
	e.HasData = true
	e.LastSegment = 0x00000000
	e.FirstSegment = 0xFFFFFFFF

	e.segments = make(map[uint32]*segment.Data)
}

func (e Channel) GetBounds() (int, int) {
	return int(e.FirstSegment) / 14, int(e.LastSegment) / 14
}

func (e *Channel) SetBounds(first, last int) {
	e.FirstSegment = uint32(first * 14)
	e.LastSegment = uint32(last * 14)
}

func (e Channel) GetDimensions() (int, int) {
	return int(e.Width), int(e.Height)
}

// Fix the channel metadata.
func (e *Channel) Process(scft lrpt.SpacecraftParameters) {
	f := e.FirstSegment % 14
	for i := uint32(0); i < f; i++ {
		e.segments[i] = segment.NewFiller()
		e.FirstSegment--
	}

	for i := e.FirstSegment; i <= e.LastSegment; i++ {
		if e.segments[i] == nil {
			e.segments[i] = segment.NewFiller()
		}
	}

	e.FileName = fmt.Sprintf("%s_%s_BISMW_%s_%d", scft.Filename, scft.SignalName, e.ChannelName, e.StartTime.GetMilliseconds())
	e.Height = (e.LastSegment - e.FirstSegment + 28) * uint32(e.BlockDim) / 14
	e.Width = e.FinalWidth
}

func (e *Channel) Parse(packet frames.SpacePacketFrame) {
	if !packet.IsValid() {
		return
	}

	if new := segment.New(packet.GetData()); new.IsValid() {
		if !e.HasData {
			e.init()
		}

		if e.last > uint32(packet.GetSequenceCount()) && e.last > 16000 {
			e.rollover += 16384
		}
		e.last = uint32(packet.GetSequenceCount())

		if uint32(new.GetMCUNumber())/14 == 0 {
			e.offset = 43 - (uint32(packet.GetSequenceCount())+e.rollover)%43
		}

		t := uint32(packet.GetSequenceCount()) + e.rollover + e.offset
		id := t/43*14 + uint32(new.GetMCUNumber())/14

		if e.LastSegment < id {
			e.LastSegment = id
			e.EndTime = new.GetDate()
		}

		if e.FirstSegment > id {
			e.FirstSegment = id
			e.StartTime = new.GetDate()
		}

		e.segments[id] = new
		e.SegmentCount++
	}
}
