package VIIRS

import (
	"bytes"
	"fmt"
	"image"
	"image/png"
	"path/filepath"
	"sort"
	"weather-dump/src/CCSDS/Frames"
	VIIRS "weather-dump/src/VIIRS/Common"
	"weather-dump/src/VIIRS/ScienceData/ScienceFrames"

	"gopkg.in/gographics/imagick.v2/imagick"
)

const firstPacket = 1
const lastPacket = 2

type ScienceData struct {
	tempSegments [2047]Segment
	dataSegments []Segment
	outputFolder string
}

type Segment struct {
	APID   uint16
	header ScienceFrames.FrameHeader
	body   [32]ScienceFrames.FrameBody
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

	fmt.Printf("[RENDER] Rendering Channel %s\n", cs.ChannelName)
	fmt.Printf("[RENDER] Coded Channel: %s <= Reconstruction Channel: %s\n", cs.ChannelName, ChannelsParameters[cs.ReconstructionBand].ChannelName)

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

	fmt.Printf("[RENDER] Rendering Channel %s\n", cs.ChannelName)

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

	w := cs.FinalProductWidth
	h := len(packets) * cs.AggregationZoneHeight

	img := image.NewGray16(image.Rect(0, 0, w, h))
	img.Pix = buf

	png_img := new(bytes.Buffer)
	png.Encode(png_img, img)

	mw := imagick.NewMagickWand()

	mw.ReadImageBlob(png_img.Bytes())
	mw.EqualizeImage()
	mw.FlopImage()
	mw.WriteImage(outputName)
}

func (e ScienceData) ExportTrueColor(SCID uint8) {
	sc := Spacecrafts[SCID]
	packets := e.GetChannelPackets(sc.TrueColorChannels[0])

	fmt.Println("[RENDER] Exporting true color image.")

	R, _ := filepath.Abs(fmt.Sprintf("%s/%s_%s_VIIRS_%s_%s.png",
		e.outputFolder,
		sc.Filename,
		sc.SignalName,
		ChannelsParameters[sc.TrueColorChannels[1]].ChannelName,
		packets[0].header.GetDate()))

	G, _ := filepath.Abs(fmt.Sprintf("%s/%s_%s_VIIRS_%s_%s.png",
		e.outputFolder,
		sc.Filename,
		sc.SignalName,
		ChannelsParameters[sc.TrueColorChannels[0]].ChannelName,
		packets[0].header.GetDate()))

	B, _ := filepath.Abs(fmt.Sprintf("%s/%s_%s_VIIRS_%s_%s.png",
		e.outputFolder,
		sc.Filename,
		sc.SignalName,
		ChannelsParameters[sc.TrueColorChannels[2]].ChannelName,
		packets[0].header.GetDate()))

	RGB, _ := filepath.Abs(fmt.Sprintf("%s/%s_%s_VIIRS_%s_%s.png",
		e.outputFolder,
		sc.Filename,
		sc.SignalName,
		"TRUECOLOR",
		packets[0].header.GetDate()))

	mwR := imagick.NewMagickWand()
	mwG := imagick.NewMagickWand()
	mwB := imagick.NewMagickWand()

	mwR.ReadImage(R)
	mwG.ReadImage(G)
	mwB.ReadImage(B)

	mwRGB := imagick.NewMagickWand()
	mwRGB.AddImage(mwR)
	mwRGB.AddImage(mwG)
	mwRGB.AddImage(mwB)
	mwRGB.ResetIterator()

	mwRGB = mwRGB.CombineImages(imagick.CHANNEL_RED | imagick.CHANNEL_GREEN | imagick.CHANNEL_BLUE)
	mwRGB.WriteImage(RGB)
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
