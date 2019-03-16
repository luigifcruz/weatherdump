package parser

import (
	"fmt"
	"runtime"
	"sync"
	"weather-dump/src/protocols/hrd"
)

func (e *Channel) Export(buf *[]byte, ch List, scft hrd.SpacecraftParameters) bool {
	fmt.Printf("[SEN] Rendering Channel %s.\n", e.ChannelName)

	if !e.HasData {
		return false
	}

	e.Process(scft)
	*buf = make([]byte, e.Width*e.Height*2)

	decimation := 1
	var codedChannel *Channel

	if e.ReconstructionBand != 000 {
		codedChannel = ch[e.ReconstructionBand]

		if !codedChannel.HasData {
			return false
		}

		if e.ChannelName[0] != codedChannel.ChannelName[0] {
			decimation = 2
		}

		codedChannel.Process(scft)
		if codedChannel.ReconstructionBand != 000 {
			if !codedChannel.Export(&[]byte{}, ch, scft) {
				return false
			}
		}
	}

	offset := 0
	threads := runtime.NumCPU() * 2
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
				for i := 0; i < e.AggregationZoneHeight; i++ {
					for j := range e.AggregationZoneWidth {

						if e.ReconstructionBand != 000 {
							if codedChannel.segments[x] == nil {
								continue
							}

							differentialData := codedChannel.segments[x].Body[i/decimation].Detector[j].GetData()
							e.segments[x].Body[i].Detector[j].Integrate(differentialData, decimation)
						}

						if len(*buf) > 0 {
							data := e.segments[x].Body[i].Detector[j].GetData()
							copy((*buf)[index:], *data)
							index += len(*data)
						}
					}
				}
			}
		}(&wg, s+int(e.FirstSegment), f+int(e.FirstSegment), offset)

		offset += ((e.AggregationZoneHeight * 2 * int(e.Width)) * (f - s))
	}

	wg.Wait()
	e.ReconstructionBand = 000
	return true
}
