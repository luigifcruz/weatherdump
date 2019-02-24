package bismw

import (
	"fmt"
	"weather-dump/src/ccsds/frames"
	"weather-dump/src/protocols/lrpt"
)

type Worker struct {
	channelData map[uint16]*Channel
}

func New() *Worker {
	e := Worker{}
	e.channelData = make(map[uint16]*Channel)
	return &e
}

func (e *Worker) Process(scid uint8) {
	for _, channel := range e.channelData {
		channel.Fix(lrpt.Spacecrafts[scid])
	}
}

func (e *Worker) SaveAllChannels(outputFolder string) {
	fmt.Println("[SEN] Exporting BISMW channels products...")
	for _, channel := range e.channelData {
		channel.Export(outputFolder)
	}
}

func (e *Worker) Parse(packet frames.SpacePacketFrame) {
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
				ch[apid].lines[ch[apid].count] = NewLine()
				ch[apid].lastFrame += 14
				ch[apid].count++
			} else {
				break
			}
		}

		if ch[apid].lines[ch[apid].count] == nil {
			ch[apid].lines[ch[apid].count] = NewLine()
		}

		ch[apid].lines[ch[apid].count/14].AddMCU(packet.GetData())
		ch[apid].lastFrame = frameCount
		ch[apid].count++
		return
	}
}
