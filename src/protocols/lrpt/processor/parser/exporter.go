package parser

import (
	"fmt"
	"weather-dump/src/protocols/lrpt"
)

func (e *Channel) Export(buf *[]byte, scft lrpt.SpacecraftParameters) bool {
	if !e.HasData {
		return false
	}

	e.Process(scft)
	*buf = make([]byte, e.Height*e.Width)

	// 14*14*8

	fmt.Println(e.LastSegment-e.FirstSegment, e.SegmentCount-1)

	index := (14 * 64) * (e.FirstSegment % 14)
	fmt.Println(index)
	for x := e.FirstSegment; x < e.LastSegment; x += 14 {
		for i := uint32(0); i < 8; i++ {
			for j := uint32(0); j < 14; j++ {
				if segment := e.segments[x+j]; segment != nil && segment.IsValid() {
					//fmt.Println(index, i, e.FirstSegment, e.LastSegment, x+j, e.Height*e.Width)
					copy((*buf)[index:], segment.Lines[i][:])
				}
				index += 8 * 14
			}
		}
	}

	fmt.Println(len(*buf), index)

	/*

		for x := uint32(0); x < e.SegmentCount/14; x++ {
			for i := uint32(0); i < 8; i++ {
				for j := uint32(0); j < 14; j++ {
					if segment := e.blocks[x].Segments[uint8(j)]; segment != nil {
						copy((*buf)[index:], segment.Lines[i][:])
					}
					index += 8 * 14
				}
			}
		}*/
	return true
}
