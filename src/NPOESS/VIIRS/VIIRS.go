package VIIRS

import (
	"fmt"
	"os"
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
				e.channelData[i].ComposeCoded(outputFolder, e.channelData[reconChannel])
			}
		}
	}
}

func (e Data) SaveTrueColorChannel(outputFolder string) {
	colorChannels := e.spacecraft.TrueColorChannels
	fmt.Println("[VIIRS] Exporting True Color Channel.")

	// Synchronize all channels scans.
	firstScan := make([]int, 3)
	lastScan := make([]int, 3)

	firstScan[0] = int(e.channelData[colorChannels[0]].start)
	firstScan[1] = int(e.channelData[colorChannels[1]].start)
	firstScan[2] = int(e.channelData[colorChannels[2]].start)

	lastScan[0] = int(e.channelData[colorChannels[0]].end)
	lastScan[1] = int(e.channelData[colorChannels[1]].end)
	lastScan[2] = int(e.channelData[colorChannels[2]].end)

	e.channelData[colorChannels[0]].end = uint32(MinIntSlice(lastScan))
	e.channelData[colorChannels[1]].end = uint32(MinIntSlice(lastScan))
	e.channelData[colorChannels[2]].end = uint32(MinIntSlice(lastScan))

	e.channelData[colorChannels[0]].start = uint32(MaxIntSlice(firstScan))
	e.channelData[colorChannels[1]].start = uint32(MaxIntSlice(firstScan))
	e.channelData[colorChannels[2]].start = uint32(MaxIntSlice(firstScan))

	// Fix channel parameters.
	e.Process()

	// Decode all channels.
	e.channelData[colorChannels[0]].ComposeUncoded("/tmp")
	e.channelData[colorChannels[1]].ComposeCoded("/tmp", e.channelData[colorChannels[0]])
	e.channelData[colorChannels[2]].ComposeCoded("/tmp", e.channelData[colorChannels[0]])

	// Generate the true color image.
	ExportTrueColor(outputFolder, e.channelData[colorChannels[1]], e.channelData[colorChannels[0]], e.channelData[colorChannels[2]])

	// Cleaning up our garbage.
	os.Remove(fmt.Sprintf("/tmp/%s.png", e.channelData[colorChannels[0]].fileName))
	os.Remove(fmt.Sprintf("/tmp/%s.png", e.channelData[colorChannels[1]].fileName))
	os.Remove(fmt.Sprintf("/tmp/%s.png", e.channelData[colorChannels[2]].fileName))
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
		frameBody := VIIRSFrames.NewFrameBody(packet.GetData())
		if frameBody.GetSequenceCount() <= ch[apid].exctdCount && frameBody.GetDetectorNumber() < 32 {
			ch[apid].segments[ch[apid].scanCount].body[frameBody.GetDetectorNumber()] = *frameBody
		}
	}
}
