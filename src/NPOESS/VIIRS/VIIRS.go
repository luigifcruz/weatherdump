package VIIRS

import (
	"fmt"
	"os"
	"weather-dump/src/CCSDS/Frames"
	"weather-dump/src/NPOESS"
	"weather-dump/src/NPOESS/VIIRS/viirsframes"
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
	fmt.Println("[VIIRS] Exporting true color channel.")

	ch01 := e.channelData[colorChannels[0]]
	ch02 := e.channelData[colorChannels[1]]
	ch03 := e.channelData[colorChannels[2]]

	// Check if required channels exist.
	if ch01 == nil || ch02 == nil || ch03 == nil {
		fmt.Println("[VIIRS] Can't export true color channel. Not all required channels are available.")
		return
	}

	// Synchronize all channels scans.
	firstScan := make([]int, 3)
	lastScan := make([]int, 3)

	firstScan[0] = int(ch01.start)
	firstScan[1] = int(ch02.start)
	firstScan[2] = int(ch03.start)

	lastScan[0] = int(ch01.end)
	lastScan[1] = int(ch02.end)
	lastScan[2] = int(ch03.end)

	ch01.end = uint32(MinIntSlice(lastScan))
	ch02.end = uint32(MinIntSlice(lastScan))
	ch03.end = uint32(MinIntSlice(lastScan))

	ch01.start = uint32(MaxIntSlice(firstScan))
	ch02.start = uint32(MaxIntSlice(firstScan))
	ch03.start = uint32(MaxIntSlice(firstScan))

	// Fix channel parameters.
	e.Process()

	// Decode all channels.
	ch01.ComposeUncoded("/tmp")
	e.channelData[colorChannels[1]].ComposeCoded("/tmp", ch01)
	e.channelData[colorChannels[2]].ComposeCoded("/tmp", ch01)

	// Generate the true color image.
	ExportTrueColor(outputFolder, ch02, ch01, ch03)

	// Cleaning up our garbage.
	os.Remove(fmt.Sprintf("/tmp/%s.png", ch01.fileName))
	os.Remove(fmt.Sprintf("/tmp/%s.png", ch02.fileName))
	os.Remove(fmt.Sprintf("/tmp/%s.png", ch03.fileName))
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

		frameHeader := viirsframes.NewFrameHeader(packet.GetData())
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
		frameBody := viirsframes.NewFrameBody(packet.GetData())
		if frameBody.GetSequenceCount() <= ch[apid].exctdCount && frameBody.GetDetectorNumber() < 32 {
			ch[apid].segments[ch[apid].scanCount].body[frameBody.GetDetectorNumber()] = *frameBody
		}
	}
}
