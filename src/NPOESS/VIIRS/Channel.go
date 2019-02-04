package VIIRS

import (
	"fmt"
	"weather-dump/src/NPOESS"
	"weather-dump/src/NPOESS/VIIRS/viirsframes"
)

const MaxFrameCount = 8192

type Channel struct {
	apid       uint16
	parameters ChannelParameters
	fileName   string
	height     uint32
	width      uint32
	startTime  NPOESS.Time
	endTime    NPOESS.Time
	segments   map[uint32]*Segment
	start      uint32
	end        uint32
	count      uint32

	scanCount  uint32
	exctdCount uint32
}

func NewChannel(apid uint16) *Channel {
	e := Channel{}
	e.apid = apid
	e.segments = make(map[uint32]*Segment)
	e.end = 0x00000000
	e.start = 0xFFFFFFFF
	e.count = 0
	return &e
}

type Segment struct {
	header *viirsframes.FrameHeader
	body   [32]viirsframes.FrameBody
}

func NewSegment(header *viirsframes.FrameHeader) *Segment {
	e := Segment{}
	e.header = header
	return &e
}

func NewFillSegment(scanNumber uint32) *Segment {
	fillFrame := Segment{}
	fillFrame.header = viirsframes.NewFillFrameHeader(scanNumber)
	for i := 0; i < 32; i++ {
		fillFrame.body[i] = *viirsframes.NewFillFrameBody()
	}
	return &fillFrame
}

func (e *Channel) Fix(scft NPOESS.SpacecraftParameters) {
	e.parameters = ChannelsParameters[e.apid]

	if e.end-e.start > MaxFrameCount {
		fmt.Printf("[VIIRS] Potentially invalid channel %s was found.\n", e.parameters.ChannelName)
		fmt.Println("	It's too long for the round earth, trying to correct...")

		if (e.end - e.end - e.count) < MaxFrameCount {
			e.start = e.end - e.count
		}

		if (e.start + e.count - e.start) < MaxFrameCount {
			e.end = e.start + e.count
		}

		if e.end-e.start > MaxFrameCount {
			fmt.Println("	Cannot find any valid number, skipping channel.")
			return
		}

		fmt.Println("	Found a valid number. Channel can still be damaged.")
	}

	for i := e.end; i >= e.start; i-- {
		if e.segments[i] == nil {
			e.segments[i] = NewFillSegment(i)
		}
	}

	e.startTime = e.segments[e.start].header.GetDate()
	e.endTime = e.segments[e.end].header.GetDate()
	e.fileName = fmt.Sprintf("%s_%s_VIIRS_%s_%s", scft.Filename, scft.SignalName, e.parameters.ChannelName, e.startTime.GetZulu())
	e.height = (e.end - e.start) * uint32(e.parameters.AggregationZoneHeight)
	e.width = e.parameters.FinalProductWidth
}

func (e Channel) ComposeUncoded(outputFolder string) {
	var buf []byte

	fmt.Printf("[VIIRS] Rendering Uncoded Channel %s\n", e.parameters.ChannelName)

	if len(e.segments) > 0 {
		for x := e.end; x >= e.start; x-- {
			packet := e.segments[x]
			for i := 0; i < e.parameters.AggregationZoneHeight; i++ {
				for j, segment := range e.parameters.AggregationZoneWidth {
					oversampleSize := e.parameters.OversampleZone[j]
					buf = append(buf, packet.body[i].GetData(j, segment, oversampleSize, false)...)
				}
			}
		}

		ExportGrayscale(buf, e, outputFolder)
	}
}

func (e *Channel) ComposeCoded(outputFolder string, r *Channel) {
	var buf []byte

	decFactor := map[bool]int{false: 2, true: 1}
	bandComp := []rune(e.parameters.ChannelName)[0] == []rune(ChannelsParameters[e.parameters.ReconstructionBand].ChannelName)[0]

	fmt.Printf("[VIIRS] Rendering Coded Channel %s with reconstruction channel %s\n",
		e.parameters.ChannelName, ChannelsParameters[e.parameters.ReconstructionBand].ChannelName)

	if len(e.segments) > 0 && len(r.segments) > 0 {
		for x := e.end; x >= e.start; x-- {
			packet := e.segments[x]
			for i := 0; i < e.parameters.AggregationZoneHeight; i++ {
				for j, segment := range e.parameters.AggregationZoneWidth {
					if r.segments[x] != nil && !packet.body[i].IsFillerFrame() && !r.segments[x].body[i/decFactor[bandComp]].IsFillerFrame() {
						var image []uint16

						baseData := packet.body[i].GetData(j, segment, e.parameters.OversampleZone[j], false)
						reconData := r.segments[x].body[i/decFactor[bandComp]].GetData(j, segment, e.parameters.OversampleZone[j], true)
						reconPixel := ConvertToU16(reconData)

						for y, basePixel := range ConvertToU16(baseData) {
							pixel := int16(basePixel) + int16(reconPixel[y/decFactor[bandComp]]) - int16(16383)
							image = append(image, uint16(pixel))
						}

						diffImage := ConvertToByte(image)
						e.segments[x].body[i].SetData(j, &diffImage)
						buf = append(buf, diffImage...)
					} else {
						buf = append(buf, make([]byte, segment*2)...)
					}
				}
			}
		}

		ExportGrayscale(buf, *e, outputFolder)
	}
}
