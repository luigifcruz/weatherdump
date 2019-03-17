package parser

import (
	"weather-dump/src/protocols/lrpt"
)

func (e *Channel) Export(buf *[]byte, scft lrpt.SpacecraftParameters) bool {
	if !e.HasData {
		return false
	}

	e.Process(scft)
	*buf = make([]byte, e.Height*e.Width)

	index := 0
	for x := uint32(0); x < e.SegmentCount/14; x++ {
		for i := uint32(0); i < 8; i++ {
			for j := uint32(0); j < 14; j++ {
				if segment := e.blocks[x].Segments[uint8(j)]; segment != nil {
					copy((*buf)[index:], segment.Lines[i][:])
				}
				index += 8 * 14
			}
		}
	}

	return true
}
