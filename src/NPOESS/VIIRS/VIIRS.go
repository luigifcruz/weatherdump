package VIIRS

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"image"
	"image/png"
	"os"
	"path/filepath"
	"sort"
	"weather-dump/src/CCSDS/Frames"
	"weather-dump/src/NPOESS"
	"weather-dump/src/NPOESS/VIIRS/VIIRSFrames"

	"gopkg.in/gographics/imagick.v2/imagick"
)

const firstPacket = 1
const lastPacket = 2

func ConvertToU16(data []byte) []uint16 {
	var buf []uint16
	for i := 0; i < len(data); i += 2 {
		buf = append(buf, binary.BigEndian.Uint16(data[i:]))
	}
	return buf
}

func ConvertToByte(data []uint16) []byte {
	var buf []byte
	bb := make([]byte, 2)

	for _, d := range data {
		binary.BigEndian.PutUint16(bb, d)
		buf = append(buf, bb...)
	}

	return buf
}

func (e Data) SaveAllChannels(SCID uint8) {
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

func (e Data) SaveCodedChannel(channelAPID uint16, SCID uint8) {
	var buf []byte
	cs := ChannelsParameters[channelAPID]

	a, z, basePackets := e.GetChannelPackets(channelAPID)
	_, _, reconPackets := e.GetChannelPackets(cs.ReconstructionBand)

	decFactor := map[bool]int{false: 2, true: 1}
	bandComp := []rune(cs.ChannelName)[0] == []rune(ChannelsParameters[cs.ReconstructionBand].ChannelName)[0]

	fmt.Printf("[VIIRS] Rendering Channel %s\n", cs.ChannelName)
	fmt.Printf("[VIIRS] (Decimation Factor: %d) Coded Channel: %s <= Reconstruction Channel: %s\n",
		decFactor[bandComp],
		cs.ChannelName,
		ChannelsParameters[cs.ReconstructionBand].ChannelName)

	if len(basePackets) > 0 && len(reconPackets) > 0 {
		for x := z; x >= a; x -= 1 {
			packet := basePackets[x]
			for i := 0; i < cs.AggregationZoneHeight; i++ {
				for j, segment := range cs.AggregationZoneWidth {
					if reconPackets[x] != nil && packet.body[i].IsValid() && reconPackets[x].body[i/decFactor[bandComp]].IsValid() {
						var image []uint16

						baseData := packet.body[i].GetData(j, segment, cs.OversampleZone[j])
						reconData := reconPackets[x].body[i/decFactor[bandComp]].GetData(j, segment, cs.OversampleZone[j])
						reconPixel := ConvertToU16(reconData)

						for y, basePixel := range ConvertToU16(baseData) {
							pixel := int16(basePixel) + int16(reconPixel[y/decFactor[bandComp]]) - int16(16383)
							image = append(image, uint16(pixel))
						}

						diffImage := ConvertToByte(image)
						basePackets[x].body[i].SetData(j, &diffImage)
						buf = append(buf, diffImage...)
					} else {
						buf = append(buf, make([]byte, segment*2)...)
					}
				}
			}
		}

		e.ProcessBuf(buf, cs, len(basePackets), e.GetTimestamp(channelAPID), SCID)
	}
}

func (e Data) SaveUncodedChannel(channelAPID uint16, SCID uint8) {
	var buf []byte
	a, z, packets := e.GetChannelPackets(channelAPID)
	cs := ChannelsParameters[channelAPID]

	fmt.Printf("[VIIRS] Rendering Uncoded Channel %s\n", cs.ChannelName)

	if len(packets) > 0 {
		for x := z; x >= a; x -= 1 {
			packet := packets[x]
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

		e.ProcessBuf(buf, cs, len(packets), e.GetTimestamp(channelAPID), SCID)
	}
}

func (e Data) ProcessBuf(buf []byte, cs ChannelParameters, len int, date string, SCID uint8) {
	sc := NPOESS.Spacecrafts[SCID]
	outputName, _ := filepath.Abs(fmt.Sprintf("%s/%s_%s_VIIRS_%s_%s.png", e.outputFolder, sc.Filename, sc.SignalName, cs.ChannelName, date))

	w := cs.FinalProductWidth
	h := len * cs.AggregationZoneHeight

	img := image.NewGray16(image.Rect(0, 0, w, h))
	img.Pix = buf

	png_img := new(bytes.Buffer)
	png.Encode(png_img, img)

	mw := imagick.NewMagickWand()
	defer mw.Destroy()

	mw.ReadImageBlob(png_img.Bytes())
	mw.EqualizeImage()
	mw.FlopImage()
	mw.WriteImage(outputName)
}

func (e Data) ExportTrueColor(SCID uint8) {
	sc := NPOESS.Spacecrafts[SCID]

	fmt.Println("[VIIRS] Exporting true color image.")

	R, _ := filepath.Abs(fmt.Sprintf("%s/%s_%s_VIIRS_%s_%s.png",
		e.outputFolder,
		sc.Filename,
		sc.SignalName,
		ChannelsParameters[sc.TrueColorChannels[1]].ChannelName,
		e.GetTimestamp(sc.TrueColorChannels[1])))

	G, _ := filepath.Abs(fmt.Sprintf("%s/%s_%s_VIIRS_%s_%s.png",
		e.outputFolder,
		sc.Filename,
		sc.SignalName,
		ChannelsParameters[sc.TrueColorChannels[0]].ChannelName,
		e.GetTimestamp(sc.TrueColorChannels[0])))

	B, _ := filepath.Abs(fmt.Sprintf("%s/%s_%s_VIIRS_%s_%s.png",
		e.outputFolder,
		sc.Filename,
		sc.SignalName,
		ChannelsParameters[sc.TrueColorChannels[2]].ChannelName,
		e.GetTimestamp(sc.TrueColorChannels[2])))

	RGB, _ := filepath.Abs(fmt.Sprintf("%s/%s_%s_VIIRS_%s_%s.png",
		e.outputFolder,
		sc.Filename,
		sc.SignalName,
		"TRUECOLOR",
		e.GetTimestamp(sc.TrueColorChannels[0])))

	if _, err := os.Stat(R); os.IsNotExist(err) {
		fmt.Println("[VIIRS] Red channel doesn't exists. Can't create true-color product.")
		return
	}

	if _, err := os.Stat(G); os.IsNotExist(err) {
		fmt.Println("[VIIRS] Green channel doesn't exists. Can't create true-color product.")
		return
	}

	if _, err := os.Stat(B); os.IsNotExist(err) {
		fmt.Println("[VIIRS] Blue channel doesn't exists. Can't create true-color product.")
		return
	}

	// Load all channels for RGB.
	mwR := imagick.NewMagickWand()
	defer mwR.Destroy()

	mwG := imagick.NewMagickWand()
	defer mwG.Destroy()

	mwB := imagick.NewMagickWand()
	defer mwB.Destroy()

	mwR.ReadImage(R)
	mwG.ReadImage(G)
	mwB.ReadImage(B)

	// Merge them togheter to create True Color Image.
	mwRGB := imagick.NewMagickWand()
	defer mwRGB.Destroy()

	mwRGB.AddImage(mwR)
	mwRGB.AddImage(mwG)
	mwRGB.AddImage(mwB)
	mwRGB.ResetIterator()

	mwRGB = mwRGB.CombineImages(imagick.CHANNEL_RED | imagick.CHANNEL_GREEN | imagick.CHANNEL_BLUE)
	mwRGB.WriteImage(RGB)
}

func (e *Data) Parse(packet Frames.SpacePacketFrame) {
	ts := &e.tempSegments[packet.GetAPID()]

	if packet.GetSequenceFlags() == firstPacket {
		ts.header.FromBinary(packet.GetData())
		ts.APID = packet.GetAPID()
	} else {
		frameBody := &VIIRSFrames.FrameBody{}
		frameBody.FromBinary(packet.GetData())
		if frameBody.IsValid() {
			ts.body[frameBody.GetDetectorNumber()].FromBinary(packet.GetData())
		}
	}

	if packet.GetSequenceFlags() == lastPacket && ts.header.GetNumberOfSegments() > 0 {
		e.dataSegments = append(e.dataSegments, *ts)
	}
}

func (e *Data) SetOutputFolder(path string) {
	e.outputFolder = path
}
