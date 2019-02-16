package BISMW

import (
	"fmt"
	"weather-dump/src/CCSDS/Frames"
	"weather-dump/src/Meteor"
)

type Data struct {
	channelData map[uint16]*Channel
	spacecraft  Meteor.SpacecraftParameters
}

func NewData(scid uint8) *Data {
	e := Data{}
	e.channelData = make(map[uint16]*Channel)
	e.spacecraft = Meteor.Spacecrafts[scid]
	return &e
}

func (e *Data) Process() {
	fmt.Println("[BISMW] Processing BISMW channels data...")
	for _, channel := range e.channelData {
		channel.Fix(e.spacecraft)
	}
}

func (e *Data) SaveAllChannels(outputFolder string) {
	fmt.Println("[BISMW] Exporting BISMW channels products...")
	for _, channel := range e.channelData {
		channel.Export(outputFolder)
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
