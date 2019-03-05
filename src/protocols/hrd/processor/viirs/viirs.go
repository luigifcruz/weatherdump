package viirs

import (
	"fmt"
	"image"
	"image/png"
	"os"
	"path/filepath"
	"weather-dump/src/assets"
	"weather-dump/src/ccsds/frames"
	"weather-dump/src/protocols/hrd"
	viirsFrames "weather-dump/src/protocols/hrd/processor/viirs/frames"
	"weather-dump/src/tools/img"
)

const firstPacket = 1

type Worker struct {
	channelData map[uint16]*Channel
}

func New() *Worker {
	e := Worker{}
	e.channelData = make(map[uint16]*Channel)
	return &e
}

func (e Worker) Export(delegate *assets.ExportDelegate, outputFolder string) {
	for _, i := range ChannelsIndex {
		if e.channelData[i] == nil {
			continue
		}

		var buf []byte
		reconChannel := e.channelData[i].parameters.ReconstructionBand
		if reconChannel == 000 {
			e.channelData[i].ComposeUncoded(&buf)
		} else {
			if e.channelData[reconChannel] != nil {
				e.channelData[i].ComposeCoded(&buf, e.channelData[reconChannel])
			}
		}
		if buf != nil {
			outputName, _ := filepath.Abs(fmt.Sprintf("%s/%s", outputFolder, e.channelData[i].fileName))
			i := img.NewGray16(&buf, int(e.channelData[i].width), int(e.channelData[i].height))

			if delegate.Equalize {
				i.Equalize()
			}

			if delegate.Flip {
				i.Flop()
			}

			if delegate.ExportPNG {
				i.ExportPNG(outputName)
			}

			if delegate.QualityJPEG > 0 {
				i.ExportJPEG(outputName, delegate.QualityJPEG)
			}
		}
	}
}

func (e Worker) SaveTrueColorChannel(scid uint8, outputFolder string) {
	colorChannels := hrd.Spacecrafts[scid].TrueColorChannels
	fmt.Println("[SEN] Exporting true color channel.")

	ch01 := e.channelData[colorChannels[0]]
	ch02 := e.channelData[colorChannels[1]]
	ch03 := e.channelData[colorChannels[2]]

	// Check if required channels exist.
	if ch01 == nil || ch02 == nil || ch03 == nil {
		fmt.Println("[SEN] Can't export true color channel. Not all required channels are available.")
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
	e.Process(scid)

	// Create output image struct.
	i := image.NewRGBA64(image.Rect(0, 0, int(ch01.width), int(ch01.height)))
	bufferSize := int(ch01.width*ch01.height) * 8
	finalImage := make([]byte, bufferSize)

	for p := 6; p < bufferSize; p += 8 {
		finalImage[p+0] = 0xFF
		finalImage[p+1] = 0xFF
	}

	// Compose images and fill buffer.
	var buf []byte
	ref := 0

	ch01.ComposeUncoded(&buf)
	img.NewGray16(&buf, int(ch01.width), int(ch01.height)).Flop().Equalize()

	for p := 2; p < bufferSize; p += 8 {
		finalImage[p+0] = buf[ref]
		finalImage[p+1] = buf[ref]
		ref += 2
	}

	buf = nil
	ref = 0
	e.channelData[colorChannels[1]].ComposeCoded(&buf, ch01)
	img.NewGray16(&buf, int(ch02.width), int(ch02.height)).Flop().Equalize()

	for p := 0; p < bufferSize; p += 8 {
		finalImage[p+0] = buf[ref]
		finalImage[p+1] = buf[ref]
		ref += 2
	}

	buf = nil
	ref = 0
	e.channelData[colorChannels[2]].ComposeCoded(&buf, ch01)
	img.NewGray16(&buf, int(ch03.width), int(ch03.height)).Flop().Equalize()

	for p := 4; p < bufferSize; p += 8 {
		finalImage[p+0] = buf[ref]
		finalImage[p+1] = buf[ref]
		ref += 2
	}

	// Render and save the true-color image.
	i.Pix = finalImage
	outputName, _ := filepath.Abs(fmt.Sprintf("%s/TRUECOLOR_VIIRS_%s.png", outputFolder, ch01.endTime.GetZuluSafe()))
	outputFile, err := os.Create(outputName)
	if err != nil {
		fmt.Println("[EXPORT] Error saving final image...")
	}
	png.Encode(outputFile, i)
	outputFile.Close()
}

func (e *Worker) Process(scid uint8) {
	for _, channel := range e.channelData {
		channel.Fix(hrd.Spacecrafts[scid])
	}
}

func (e *Worker) Parse(packet frames.SpacePacketFrame) {
	ch := e.channelData
	apid := packet.GetAPID()

	if packet.GetSequenceFlags() == firstPacket && packet.IsValid() {
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
