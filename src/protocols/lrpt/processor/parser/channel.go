package parser

import (
	"fmt"
	"weather-dump/src/ccsds/frames"
	"weather-dump/src/protocols/lrpt"
	"weather-dump/src/protocols/lrpt/processor/parser/block"
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

	segments  map[uint32]*block.Segment
	rollover  uint16
	lastCount uint16
}

// NewChannel instance.
func (e *Channel) init() {
	e.HasData = true
	e.LastSegment = 0x00000000
	e.FirstSegment = 0xFFFFFFFF

	e.segments = make(map[uint32]*block.Segment)
}

func (e Channel) GetDimensions() (int, int) {
	return int(e.Width), int(e.Height)
}

// Fix the channel metadata.
func (e *Channel) Process(scft lrpt.SpacecraftParameters) {
	for i := e.FirstSegment; i <= e.LastSegment; i++ {
		if e.segments[i] == nil {
			e.segments[i] = block.NewFillSegment()
			e.SegmentCount++
		}
	}

	e.StartTime = e.segments[e.FirstSegment].GetDate()
	e.EndTime = e.segments[e.LastSegment].GetDate()
	e.FileName = fmt.Sprintf("%s_%s_BISMW_%s_%d", scft.Filename, scft.SignalName, e.ChannelName, e.StartTime.GetMilliseconds())
	e.Height = (e.SegmentCount + 28) * uint32(e.BlockDim) / 14
	e.Width = e.FinalWidth
}

func (e *Channel) Parse(packet frames.SpacePacketFrame) {
	if !packet.IsValid() {
		return
	}

	if !e.HasData {
		e.init()
	}

	if e.lastCount-e.rollover > packet.GetSequenceCount() {
		fmt.Println("rollover", e.rollover)
		e.rollover += 16383
	}
	e.lastCount = packet.GetSequenceCount() + e.rollover

	new := block.NewSegment(packet.GetData())
	id := uint32(e.lastCount/43*14) + uint32(new.GetMCUNumber()/14)
	//fmt.Println(e.lastCount, new.GetID()/88, new.GetMCUNumber()/14, id)

	e.segments[id] = new

	if e.LastSegment < id {
		e.LastSegment = id
	}

	if e.FirstSegment > id {
		e.FirstSegment = id
	}

	e.SegmentCount++
}
