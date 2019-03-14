package parser

import (
	"encoding/binary"
	"fmt"
	"runtime"
	"sync"
)

func (e Channel) Export(nuf *[]byte, ch ChannelList) {

}

func (e Channel) ExportUncoded(buf *[]byte) {
	fmt.Printf("[SEN] Rendering Uncoded Channel %s\n", e.ChannelName)

	index := 0
	for x := e.LastSegment; x >= e.FirstSegment; x-- {
		for i := 0; i < e.AggregationZoneHeight; i++ {
			for j, segment := range e.AggregationZoneWidth {
				data := e.segments[x].Body[i].GetData(j, segment, e.OversampleZone[j])
				copy((*buf)[index:], data)
				index += len(data)
			}
		}
	}
}

func (e *Channel) ExportCoded(buf *[]byte, r *Channel) {
	decFactor := map[bool]int{false: 2, true: 1}
	bandComp := []rune(e.ChannelName)[0] == []rune(Channels[e.ReconstructionBand].ChannelName)[0]

	fmt.Printf("[SEN] Rendering Coded Channel %s with reconstruction channel %s\n",
		e.ChannelName, Channels[e.ReconstructionBand].ChannelName)

	offset := 0
	threads := runtime.NumCPU()
	segments := e.LastSegment - e.FirstSegment

	var wg sync.WaitGroup
	wg.Add(int(threads))

	for t := 0; t < threads; t++ {
		f := (int(segments) / threads * (threads - t))
		s := (int(segments) / threads * (threads - t - 1))

		if t == 0 {
			f = int(segments)
		}

		if t == threads-1 {
			s--
		}

		go func(wg *sync.WaitGroup, start, finish, index int) {
			defer wg.Done()
			for x := uint32(finish); x > uint32(start); x-- {
				packet := e.segments[x]
				for i := 0; i < e.AggregationZoneHeight; i++ {
					for j, segment := range e.AggregationZoneWidth {
						if r.segments[x] == nil {
							continue
						}

						baseData := packet.Body[i].GetData(j, segment, e.OversampleZone[j])
						reconData := r.segments[x].Body[i/decFactor[bandComp]].GetData(j, segment, e.OversampleZone[j])

						for b := 0; b < len(baseData); b += 2 {
							newPixel := binary.BigEndian.Uint16(baseData[b:]) + binary.BigEndian.Uint16(reconData[b/decFactor[bandComp]/2*2:]) - 16383
							binary.BigEndian.PutUint16((*buf)[index+b:], newPixel)
						}

						if !e.decoded {
							e.segments[x].Body[i].SetData(j, (*buf)[index:index+len(baseData)])
						}
						index += len(baseData)
					}
				}
			}
		}(&wg, s+int(e.FirstSegment), f+int(e.FirstSegment), offset)

		offset += ((e.AggregationZoneHeight * 2 * int(e.Width)) * (f - s))
	}

	wg.Wait()
	e.decoded = true
}
