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
	lastSeq  uint32
	lastTime uint32
	offset   uint32
}

// NewChannel instance.
func (e *Channel) init() {
	e.HasData = true
	e.LastSegment = 0x00000000
	e.FirstSegment = 0xFFFFFFFF

	e.segments = make(map[uint32]*segment.Data)
}

// GetBounds returns the number of the first and last segment
// of the current channel. This should be called after Process().
func (e Channel) GetBounds() (int, int) {
	return int(e.FirstSegment) / 14, int(e.LastSegment) / 14
}

// SetBounds for the passed values.
// After calling this function the Process() also should be called.
func (e *Channel) SetBounds(first, last int) {
	e.FirstSegment = uint32(first * 14)
	e.LastSegment = uint32(last * 14)
}

// GetDimensions returns the width and height of the current channel.
// Should be called after the Process().
func (e Channel) GetDimensions() (int, int) {
	return int(e.Width), int(e.Height)
}

// Process corrects the current channel metadata.
// Should be called every time SetBounds() is called.
func (e *Channel) Process(scft lrpt.SpacecraftParameters) {
	e.FirstSegment -= e.FirstSegment % 14
	e.LastSegment += e.LastSegment % 14

	for i := e.FirstSegment; i <= e.LastSegment; i++ {
		if e.segments[i] == nil {
			e.segments[i] = segment.NewFiller()
		}
	}

	e.FileName = fmt.Sprintf("%s_%s_BISMW_%s_%d", scft.Filename, scft.SignalName, e.ChannelName, e.StartTime.GetMilliseconds())
	e.Height = (e.LastSegment - e.FirstSegment + 1) * uint32(e.BlockDim) / 14
	e.Width = e.FinalWidth
}

// Parse the current Space Packet Frame into each LRPT protocol channel structure.
func (e *Channel) Parse(packet frames.SpacePacketFrame) {
	if new := segment.New(packet.GetData()); new.IsValid() && packet.IsValid() {
		if !e.HasData {
			e.init()
		}

		sequence := uint32(packet.GetSequenceCount())
		mcuNumber := uint32(new.GetMCUNumber()) / 14

		if e.lastSeq > sequence && e.lastSeq > 16000 && e.lastTime < new.GetDate().GetMilliseconds() {
			e.rollover += 16384
		}

		if mcuNumber == 0 {
			e.offset = 43 - (sequence+e.rollover)%43
		}

		id := ((sequence + e.rollover + e.offset) / 43 * 14) + mcuNumber

		if e.LastSegment < id {
			e.LastSegment = id
			e.EndTime = new.GetDate()
		}

		if e.FirstSegment > id {
			e.FirstSegment = id
			e.StartTime = new.GetDate()
		}

		e.lastSeq = sequence
		e.lastTime = new.GetDate().GetMilliseconds()
		e.segments[id] = new
		e.SegmentCount++
	}
}
