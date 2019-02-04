package BISMW

import (
	"fmt"
	"weather-dump/src/CCSDS/Frames"
)

type Data struct {
	channelData map[uint16]*Channel
}

func NewData(scid uint8) *Data {
	e := Data{}
	e.channelData = make(map[uint16]*Channel)
	return &e
}

func (e *Data) Process() {
	fmt.Println("[BISMW] Processing BISMW channels data...")
	for _, channel := range e.channelData {
		channel.Fix()
	}
}

func (e *Data) Parse(packet Frames.SpacePacketFrame) {
	ch := e.channelData
	apid := packet.GetAPID()

	if packet.IsValid() {
		frameCount := uint32(packet.GetSequenceCount())

		if ch[apid] == nil {
			ch[apid] = NewChannel(apid)
			ch[apid].lastFrame = frameCount - 30
		}

		for {
			if frameCount-ch[apid].lastFrame > 30 && frameCount-ch[apid].lastFrame < 16350 {
				if ch[apid].lines[ch[apid].count] == nil {
					ch[apid].lines[ch[apid].count] = NewLine()
				}
				ch[apid].lastFrame += 14
				ch[apid].count++
				//fmt.Println(frameCount, ch[apid].lastFrame, frameCount-ch[apid].lastFrame, "Add filler...")
			} else {
				//fmt.Println(frameCount, ch[apid].lastFrame, frameCount-ch[apid].lastFrame, "Pass.")
				break
			}
		}

		if ch[apid].lines[ch[apid].count] == nil {
			ch[apid].lines[ch[apid].count] = NewLine()
		}

		segment := NewSegment(packet.GetData())
		ch[apid].lines[ch[apid].count].segments[segment.GetMCUNumber()] = segment

		ch[apid].lastFrame = frameCount
		ch[apid].count++
		return
	}
}
