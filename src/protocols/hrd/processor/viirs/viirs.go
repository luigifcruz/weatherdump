package viirs

import (
	"weather-dump/src/ccsds/frames"
	viirsFrames "weather-dump/src/protocols/hrd/processor/viirs/frames"
)

type Worker struct {
	channelData map[uint16]*Channel
}

func New() *Worker {
	return &Worker{make(map[uint16]*Channel)}
}

func (e *Worker) Channel(apid uint16) *Channel {
	return e.channelData[apid]
}

func (e *Worker) Parse(packet frames.SpacePacketFrame) {
	ch := e.channelData
	apid := packet.GetAPID()

	if packet.GetSequenceFlags() == 1 && packet.IsValid() {
		if ch[apid] == nil {
			ch[apid] = NewChannel(apid)
		}

		frameHeader := viirsFrames.NewFrameHeader(packet.GetData())
		ch[apid].scanCount = frameHeader.GetScanNumber()
		ch[apid].exctdCount = frameHeader.GetSequenceCount() + uint32(frameHeader.GetNumberOfSegments()) + 2
		ch[apid].segments[ch[apid].scanCount] = NewSegment(frameHeader)

		if ch[apid].end < ch[apid].scanCount {
			ch[apid].end = ch[apid].scanCount
		}

		if ch[apid].start > ch[apid].scanCount {
			ch[apid].start = ch[apid].scanCount
		}

		ch[apid].count++
		return
	}

	if ch[apid] != nil {
		frameBody := viirsFrames.NewFrameBody(packet.GetData())
		if frameBody.GetSequenceCount() <= ch[apid].exctdCount && frameBody.GetDetectorNumber() < 32 {
			ch[apid].segments[ch[apid].scanCount].body[frameBody.GetDetectorNumber()] = *frameBody
		}
	}
}
