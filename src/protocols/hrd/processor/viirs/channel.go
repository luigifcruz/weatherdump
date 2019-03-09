package viirs

import (
	"encoding/binary"
	"fmt"
	"runtime"
	"sync"
	"weather-dump/src/protocols/hrd"
	"weather-dump/src/protocols/hrd/processor/viirs/frames"
)

const maxFrameCount = 8192

type Channel struct {
	apid       uint16
	parameters ChannelParameters
	fileName   string
	height     uint32
	width      uint32
	startTime  hrd.Time
	endTime    hrd.Time
	segments   map[uint32]*Segment
	start      uint32
	end        uint32
	count      uint32

	scanCount  uint32
	exctdCount uint32
}

func NewChannel(apid uint16) *Channel {
	return &Channel{
		apid:     apid,
		segments: make(map[uint32]*Segment),
		end:      0x00000000,
		start:    0xFFFFFFFF,
		count:    0x00000000,
	}
}

type Segment struct {
	header *frames.FrameHeader
	body   [32]frames.FrameBody
}

func NewSegment(header *frames.FrameHeader) *Segment {
	return &Segment{
		header: header,
	}
}

func NewFillSegment(scanNumber uint32) *Segment {
	fillFrame := Segment{header: frames.NewFillFrameHeader(scanNumber)}
	for i := 0; i < 32; i++ {
		fillFrame.body[i] = *frames.NewFillFrameBody()
	}
	return &fillFrame
}

func (e Channel) GetReconstructionBand() uint16 {
	return e.parameters.ReconstructionBand
}

func (e Channel) GetFileName() string {
	return e.fileName
}

func (e Channel) GetDimensions() (int, int) {
	return int(e.width), int(e.height)
}

func (e *Channel) Fix(scft hrd.SpacecraftParameters) {
	e.parameters = ChannelsParameters[e.apid]

	if e.end-e.start > maxFrameCount {
		fmt.Printf("[SEN] Potentially invalid channel %s was found.\n", e.parameters.ChannelName)
		fmt.Println("	It's too long for the round earth, trying to correct...")

		if (e.end - e.end - e.count) < maxFrameCount {
			e.start = e.end - e.count
		}

		if (e.start + e.count - e.start) < maxFrameCount {
			e.end = e.start + e.count
		}

		if e.end-e.start > maxFrameCount {
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
	e.fileName = fmt.Sprintf("%s_%s_VIIRS_%s_%s", scft.Filename, scft.SignalName, e.parameters.ChannelName, e.startTime.GetZuluSafe())
	e.height = (e.end - e.start + 1) * uint32(e.parameters.AggregationZoneHeight)
	e.width = e.parameters.FinalProductWidth
}

func (e Channel) ComposeUncoded(buf *[]byte) {
	fmt.Printf("[SEN] Rendering Uncoded Channel %s\n", e.parameters.ChannelName)

	if !(len(e.segments) > 0) {
		buf = nil
		return
	}

	index := 0
	for x := e.end; x >= e.start; x-- {
		for i := 0; i < e.parameters.AggregationZoneHeight; i++ {
			for j, segment := range e.parameters.AggregationZoneWidth {
				data := e.segments[x].body[i].GetData(j, segment, e.parameters.OversampleZone[j], false)
				copy((*buf)[index:], data)
				index += len(data)
			}
		}
	}
}

func (e *Channel) ComposeCoded(buf *[]byte, r *Channel) {
	decFactor := map[bool]int{false: 2, true: 1}
	bandComp := []rune(e.parameters.ChannelName)[0] == []rune(ChannelsParameters[e.parameters.ReconstructionBand].ChannelName)[0]

	fmt.Printf("[SEN] Rendering Coded Channel %s with reconstruction channel %s\n",
		e.parameters.ChannelName, ChannelsParameters[e.parameters.ReconstructionBand].ChannelName)

	if !(len(e.segments) > 0 && len(r.segments) > 0) {
		buf = nil
		return
	}

	offset := 0
	threads := runtime.NumCPU()
	segments := e.end - e.start

	var wg sync.WaitGroup
	wg.Add(int(threads))

	for t := 0; t < threads; t++ {
		f := (int(segments) / threads * (threads - t))
		s := (int(segments) / threads * (threads - t - 1))

		if t == 0 {
			f = int(segments)
		}

		go func(wg *sync.WaitGroup, start, finish, index int) {
			defer wg.Done()
			for x := uint32(finish); x >= uint32(start); x-- {
				packet := e.segments[x]
				for i := 0; i < e.parameters.AggregationZoneHeight; i++ {
					for j, segment := range e.parameters.AggregationZoneWidth {
						if r.segments[x] == nil || packet.body[i].IsFillerFrame() && r.segments[x].body[i/decFactor[bandComp]].IsFillerFrame() {
							continue
						}

						baseData := packet.body[i].GetData(j, segment, e.parameters.OversampleZone[j], false)
						reconData := r.segments[x].body[i/decFactor[bandComp]].GetData(j, segment, e.parameters.OversampleZone[j], true)

						for b := 0; b < len(baseData); b += 2 {
							newPixel := binary.BigEndian.Uint16(baseData[b:]) + binary.BigEndian.Uint16(reconData[b/decFactor[bandComp]/2*2:]) - 16383
							binary.BigEndian.PutUint16((*buf)[index+b:], newPixel)
						}

						e.segments[x].body[i].SetData(j, (*buf)[index:index+len(baseData)])
						index += len(baseData)
					}
				}
			}
		}(&wg, s+int(e.start), f+int(e.start), offset)

		offset += ((e.parameters.AggregationZoneHeight * 2 * int(e.width)) * (f - s))
	}

	wg.Wait()
}
