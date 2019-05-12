package parser

import (
	"fmt"
	"weatherdump/src/ccsds/frames"
	"weatherdump/src/protocols/lrpt"
	"weatherdump/src/protocols/lrpt/processor/parser/segment"
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
func (e Channel) GetBounds() (int, int, int) {
	return int(e.FirstSegment / 14), int(e.LastSegment / 14), int(e.offset)
}

// SetBounds for the passed values.
// After calling this function the Process() also should be called.
func (e *Channel) SetBounds(first, last, offset int) {
	e.FirstSegment = uint32((first + (offset - int(e.offset))) * 14)
	e.LastSegment = uint32((last + (offset - int(e.offset))) * 14)
}

// GetDimensions returns the width and height of the current channel.
// Should be called after the Process().
func (e Channel) GetDimensions() (int, int) {
	return int(e.Width), int(e.Height)
}

// GetTime returns the time of the first and last valid frames.
// Should be called after the Process().
func (e Channel) GetTime() (int, int) {
	return int(e.StartTime.GetMilliseconds()), int(e.EndTime.GetMilliseconds())
}

// Process corrects the current channel metadata.
// Should be called every time SetBounds() is called.
func (e *Channel) Process(scft lrpt.SpacecraftParameters) {
	e.FirstSegment -= e.FirstSegment % 14
	e.LastSegment -= e.LastSegment % 14

	for i := e.FirstSegment; i <= e.LastSegment; i++ {
		if e.segments[i] == nil {
			e.segments[i] = segment.NewFiller()
		}
	}

	for i := e.FirstSegment; i <= e.LastSegment; i++ {
		if e.segments[i].GetDate().GetMilliseconds() != uint32(0) {
			e.StartTime = e.segments[i].GetDate()
			break
		}
	}

	for i := e.LastSegment; i >= e.FirstSegment; i-- {
		if e.segments[i].GetDate().GetMilliseconds() != uint32(0) {
			e.EndTime = e.segments[i].GetDate()
			break
		}
	}

	e.FileName = fmt.Sprintf("%s_%s_BISMW_%s_%d", scft.Filename, scft.SignalName, e.ChannelName, e.StartTime.GetMilliseconds())
	e.Height = ((e.LastSegment - e.FirstSegment) / 14) * 8
	e.Width = e.FinalWidth

	if e.Height*e.Width < 100 {
		e.HasData = false
	}
}

// Parse the current Space Packet Frame into each LRPT protocol channel structure.
func (e *Channel) Parse(packet frames.SpacePacketFrame) {
	if new := segment.New(packet.GetData()); new.IsValid() && packet.IsValid() {
		if !e.HasData {
			e.init()
		}

		sequence := uint32(packet.GetSequenceCount())
		mcuNumber := uint32(new.GetMCUNumber()) / 14

		if e.lastSeq > sequence && e.lastSeq > 16000 && sequence < 1000 {
			e.rollover += 16384
		}

		if mcuNumber == 0 && e.offset == 0 {
			e.offset = (sequence + e.rollover) % 43 % 14
			fmt.Println(e.offset)
		}

		id := ((sequence + e.rollover - e.offset) / 43 * 14) + mcuNumber

		if e.LastSegment < id {
			e.LastSegment = id
		}

		if e.FirstSegment > id {
			e.FirstSegment = id
		}

		e.lastSeq = sequence
		e.segments[id] = new
		e.SegmentCount++
	}
}
