package VIIRS

import (
	"osp-noaa-dump/CCSDS/Frames"
	"osp-noaa-dump/VIIRS/ScienceData/ScienceFrames"
	"unsafe"
	"sort"
	"image"
	"os"
	"image/png"
	"fmt"
	"osp-noaa-dump/VIIRS/Common"
	"strings"
)

/*
#include <stdlib.h>
#include <stdint.h>
#include <libaec.h>
#include <string.h>
#cgo LDFLAGS: -laec

void decompress(char *input, char *output, int inputLen, int outputLen) {
	struct aec_stream strm;

	strm.bits_per_sample = 15;
	strm.block_size = 8;
	strm.rsi = 128;
	strm.flags = AEC_DATA_MSB | AEC_DATA_PREPROCESS;
	strm.next_in = input;
	strm.avail_in = inputLen;
	strm.next_out = output;
	strm.avail_out = outputLen * sizeof(char);

	aec_decode_init(&strm);
	aec_decode(&strm, AEC_FLUSH);
	aec_decode_end(&strm);
}
*/
import "C"

func Decompress(data []byte, inputLen int, outputLen int) []byte {
	var slice = make([]byte, outputLen)
	C.decompress((*C.char)(unsafe.Pointer(&data[0])), (*C.char)(unsafe.Pointer(&slice[0])), C.int(inputLen), C.int(outputLen))
	return slice
}

const firstPacket = 1
const lastPacket = 2

type ScienceData struct {
	tempSegments [2047]Segment
	dataSegments []Segment
}

type Segment struct {
	APID uint16
	header ScienceFrames.FrameHeader
	body [32]ScienceFrames.FrameBody
}

func (e ScienceData) GetChannelPackets(channelAPID uint16) []Segment {
	var list []Segment

	for _, segment := range e.dataSegments {
		if segment.APID == channelAPID {
			list = append(list, segment)
		}
	}

	sort.SliceStable(list, func(i, j int) bool {
		return list[j].header.GetSequenceCount() < list[i].header.GetSequenceCount()
	})

	return list
}

func (e ScienceData) SaveAllChannels(SCID uint8) {
	for _, channel := range ChannelsParameters {
		e.SaveChannel(channel.APID, SCID)
	}
}

func (e ScienceData) SaveChannel(channelAPID uint16, SCID uint8) {
	var buf []byte
	packets := e.GetChannelPackets(channelAPID)
	cs := ChannelsParameters[channelAPID]

	if len(packets) > 0 {
		for _, packet := range packets {
			for i := 0; i < cs.AggregationZoneHeight; i++ {
				for j, segment := range cs.AggregationZoneWidth {
					if packet.body[i].IsValid() {
						dat := packet.body[i].GetData(j)
						buf = append(buf, Decompress(dat, len(dat), segment*2)...)
					} else {
						buf = append(buf, make([]byte, segment*2)...)
					}
				}
			}
		}

		VIIRS.NormalizeImage(&buf)

		sc := Spacecrafts[SCID]
		outputName := fmt.Sprintf("output/%s_%s_VIIRS_%s_%s.png", sc.Filename, sc.SignalName, cs.ChannelName, packets[0].header.GetDate())
		outputFile, _ := os.Create(outputName)

		img := image.NewGray16(image.Rect(0, 0, cs.FinalProductWidth, len(packets)*cs.AggregationZoneHeight))
		img.Pix = buf

		if strings.ContainsAny(cs.ChannelName, "I") {
			PerformInterpolation(img, cs)
		}

		encoder := png.Encoder{ CompressionLevel: png.NoCompression }
		encoder.Encode(outputFile, img)
		outputFile.Close()
	}
}

func (e *ScienceData) Parse(packet Frames.SpacePacketFrame) {
	ts := &e.tempSegments[packet.GetAPID()]

	if packet.GetSequenceFlags() == firstPacket {
		ts.header.FromBinary(packet.GetData())
		ts.APID = packet.GetAPID()
	} else {
		frameBody := &ScienceFrames.FrameBody{}
		frameBody.FromBinary(packet.GetData())
		ts.body[frameBody.GetDetectorNumber()].FromBinary(packet.GetData())
	}

	if packet.GetSequenceFlags() == lastPacket && ts.header.GetNumberOfSegments() > 0 {
		e.dataSegments = append(e.dataSegments, *ts)
	}
}
