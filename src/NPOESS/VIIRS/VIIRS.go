package VIIRS

import (
	"fmt"
	"weather-dump/src/CCSDS/Frames"
	"weather-dump/src/NPOESS"
	"weather-dump/src/NPOESS/VIIRS/VIIRSFrames"
)

const firstPacket = 1
const lastPacket = 2

type Data struct {
	tempSegments [2047]Segment
	dataSegments []Segment
	outputFolder string

	spacecraft  NPOESS.SpacecraftParameters
	channelData map[uint16]*Channel
}

func NewData(scid uint8) *Data {
	e := Data{}
	e.spacecraft = NPOESS.Spacecrafts[scid]
	e.channelData = make(map[uint16]*Channel)
	return &e
}

func (e Data) SaveAllChannels(outputFolder string) {
	for _, i := range ChannelsIndex {
		if e.channelData[i] == nil {
			continue
		}

		reconChannel := e.channelData[i].parameters.ReconstructionBand
		if reconChannel == 000 {
			e.channelData[i].ComposeUncoded(outputFolder)
		} else {
			if e.channelData[reconChannel] != nil {
				//e.channelData[i].ComposeCoded(outputFolder, e.channelData[reconChannel])
			}
		}
	}

	//colorChannels := e.spacecraft.TrueColorChannels
	//ExportTrueColor(outputFolder, e.channelData[colorChannels[1]], e.channelData[colorChannels[0]], e.channelData[colorChannels[2]])
}

func (e *Data) Process() {
	fmt.Println("[VIIRS] Processing VIIRS channels data...")
	for _, channel := range e.channelData {
		channel.Fix(e.spacecraft)
	}
}

func (e *Data) Parse(packet Frames.SpacePacketFrame) {
	ch := e.channelData
	apid := packet.GetAPID()

	if packet.GetSequenceFlags() == firstPacket && packet.IsValid() {
		if ch[apid] == nil {
			ch[apid] = NewChannel(apid)
		}

		frameHeader := VIIRSFrames.NewFrameHeader(packet.GetData())
		ch[apid].scanCount = frameHeader.GetScanNumber()
		ch[apid].exctdCount = frameHeader.GetSequenceCount() + uint32(frameHeader.GetNumberOfSegments()) + 2
		ch[apid].segments[ch[apid].scanCount] = NewSegment(frameHeader)
		//fmt.Printf("%d %d %d==%d\n", apid, ch[apid].scanCount, ch[apid].exctdCount, frameHeader.GetSequenceCount())

		if ch[apid].end < ch[apid].scanCount {
			ch[apid].end = ch[apid].scanCount
		}

		if ch[apid].start > ch[apid].scanCount {
			ch[apid].start = ch[apid].scanCount
		}

		return
	}

	if ch[apid] != nil {
		frameBody := VIIRSFrames.NewFrameBody(packet.GetData())
		//fmt.Printf("%d %d %d>%d %d (%d)\n", apid, ch[apid].scanCount, ch[apid].exctdCount, frameBody.GetSequenceCount(), frameBody.GetDetectorNumber(), frameBody.GetSequenceCount() <= ch[apid].exctdCount)
		if frameBody.GetSequenceCount() <= ch[apid].exctdCount && frameBody.GetDetectorNumber() < 32 {
			ch[apid].segments[ch[apid].scanCount].body[frameBody.GetDetectorNumber()] = *frameBody
		}
	}
}
