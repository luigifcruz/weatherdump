package parser

import (
	"fmt"
	"weather-dump/src/ccsds/frames"
	"weather-dump/src/protocols/lrpt"
	"weather-dump/src/protocols/lrpt/processor/parser/block"
)

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
	LastFrame    uint32

	blocks map[uint32]*block.Data
}

// NewChannel instance.
func (e *Channel) init() {
	e.HasData = true
	e.blocks = make(map[uint32]*block.Data)
}

func (e Channel) GetDimensions() (int, int) {
	return int(e.Width), int(e.Height)
}

// Fix the channel metadata.
func (e *Channel) Process(scft lrpt.SpacecraftParameters) {
	e.StartTime = e.blocks[0].GetDate()
	e.EndTime = e.blocks[e.SegmentCount/14].GetDate()
	e.FileName = fmt.Sprintf("%s_%s_BISMW_%s_%d", scft.Filename, scft.SignalName, e.ChannelName, e.StartTime.GetMilliseconds())
	e.Height = (e.SegmentCount + 1) * uint32(e.BlockDim) / 14
	e.Width = e.FinalWidth
}

func (e *Channel) Parse(packet frames.SpacePacketFrame) {
	if !packet.IsValid() {
		return
	}

	frameCount := uint32(packet.GetSequenceCount())

	if !e.HasData {
		e.init()
		e.LastFrame = frameCount - 30
	}

	for {
		if frameCount-e.LastFrame > 30 && frameCount-e.LastFrame < 16350 {
			e.blocks[e.SegmentCount] = block.New()
			e.LastFrame += 14
			e.SegmentCount++
		} else {
			break
		}
	}

	if e.blocks[e.SegmentCount] == nil {
		e.blocks[e.SegmentCount] = block.New()
	}

	e.blocks[e.SegmentCount/14].AddMCU(packet.GetData())
	e.LastFrame = frameCount
	e.SegmentCount++
}
