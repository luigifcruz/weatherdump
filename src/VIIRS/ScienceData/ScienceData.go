package VIIRS

import (
	"fmt"
	"image"
	"image/png"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"weather-dump/src/CCSDS/Frames"
	"weather-dump/src/VIIRS/Common"
	"weather-dump/src/VIIRS/ScienceData/ScienceFrames"
)

const firstPacket = 1
const lastPacket = 2

type ScienceData struct {
	tempSegments [2047]Segment
	dataSegments []Segment
	outputFolder string
}

type Segment struct {
	APID uint16
	header ScienceFrames.FrameHeader
	body [32]ScienceFrames.FrameBody
}

func (e ScienceData) GetChannelPackets(channelAPID uint16) []*Segment {
	var list []*Segment

	for i, segment := range e.dataSegments {
		if segment.APID == channelAPID {
			list = append(list, &e.dataSegments[i])
		}
	}

	sort.SliceStable(list, func(i, j int) bool {
		return list[j].header.GetSequenceCount() < list[i].header.GetSequenceCount()
	})

	return list
}

func (e ScienceData) SaveAllChannels(SCID uint8) {
	var codedChannels []int

	for _, channel := range ChannelsParameters {
		if channel.ReconstructionBand == 000 {
			e.SaveUncodedChannel(channel.APID, SCID)
		} else {
			codedChannels = append(codedChannels, int(channel.APID))
		}
	}

	sort.Ints(codedChannels)

	for _, channel := range codedChannels {
		e.SaveCodedChannel(uint16(channel), SCID)
	}
}

func (e ScienceData) SaveCodedChannel(channelAPID uint16, SCID uint8) {
	var buf []byte
	cs := ChannelsParameters[channelAPID]

	basePackets := e.GetChannelPackets(channelAPID)
	reconPackets := e.GetChannelPackets(cs.ReconstructionBand)

	fmt.Printf("Coded Channel: %s <= Reconstruction Channel: %s\n", cs.ChannelName, ChannelsParameters[cs.ReconstructionBand].ChannelName)

	if len(basePackets) > 0 && len(reconPackets) > 0 {
		for x, packet := range basePackets {
			for i := 0; i < cs.AggregationZoneHeight; i++ {
				for j, segment := range cs.AggregationZoneWidth {
					if packet.body[i].IsValid() {
						var image []uint16

						baseData := packet.body[i].GetData(j, segment, cs.OversampleZone[j])
						reconData := reconPackets[x].body[i].GetData(j, segment, cs.OversampleZone[j])

						reconPixel := VIIRS.ConvertToU16(reconData)
						for y, basePixel := range VIIRS.ConvertToU16(baseData) {
							pixel := int16(basePixel) + int16(reconPixel[y]) - int16(16383)
							image = append(image, uint16(pixel))
						}


						diffImage := VIIRS.ConvertToByte(image)
						basePackets[x].body[i].SetData(j, &diffImage)
						buf = append(buf, diffImage...)
					} else {
						buf = append(buf, make([]byte, segment*2)...)
					}
				}
			}
		}

		e.ProcessBuf(buf, cs, basePackets, SCID)
	}
}

func (e ScienceData) SaveUncodedChannel(channelAPID uint16, SCID uint8) {
	var buf []byte
	packets := e.GetChannelPackets(channelAPID)
	cs := ChannelsParameters[channelAPID]

	if len(packets) > 0 {
		for _, packet := range packets {
			for i := 0; i < cs.AggregationZoneHeight; i++ {
				for j, segment := range cs.AggregationZoneWidth {
					if packet.body[i].IsValid() {
						oversampleSize := cs.OversampleZone[j]
						buf = append(buf, packet.body[i].GetData(j, segment, oversampleSize)...)
					} else {
						buf = append(buf, make([]byte, segment*2)...)
					}
				}
			}
		}

		e.ProcessBuf(buf, cs, packets, SCID)
	}
}

func (e ScienceData) ProcessBuf(buf []byte, cs ChannelParameters, packets []*Segment, SCID uint8) {
	sc := Spacecrafts[SCID]
	outputName, _ := filepath.Abs(fmt.Sprintf("%s/%s_%s_VIIRS_%s_%s.png", e.outputFolder, sc.Filename, sc.SignalName, cs.ChannelName, packets[0].header.GetDate()))
	outputFile, _ := os.Create(outputName)

	img := image.NewGray16(image.Rect(0, 0, cs.FinalProductWidth, len(packets)*cs.AggregationZoneHeight))
	img.Pix = buf

	if strings.ContainsAny(cs.ChannelName, "I") {
		PerformInterpolation(img, cs)
	}

	encoder := png.Encoder{ CompressionLevel: png.NoCompression }
	encoder.Encode(outputFile, img)

	outputFile.Close()

	VIIRS.HistEqualization(outputName)
}

func (e *ScienceData) Parse(packet Frames.SpacePacketFrame) {
	ts := &e.tempSegments[packet.GetAPID()]

	if packet.GetSequenceFlags() == firstPacket {
		ts.header.FromBinary(packet.GetData())
		ts.APID = packet.GetAPID()
	} else {
		frameBody := &ScienceFrames.FrameBody{}
		frameBody.FromBinary(packet.GetData())
		if frameBody.IsValid() {
			ts.body[frameBody.GetDetectorNumber()].FromBinary(packet.GetData())
		}
	}

	if packet.GetSequenceFlags() == lastPacket && ts.header.GetNumberOfSegments() > 0 {
		e.dataSegments = append(e.dataSegments, *ts)
	}
}

func (e *ScienceData) SetOutputFolder(path string) {
	e.outputFolder = path
}
