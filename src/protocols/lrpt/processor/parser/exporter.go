package parser

import (
	"weatherdump/src/protocols/lrpt"
)

// Export the assets data inside the current LRPT channel.
// Data allocation with the current bounds occurs inside this function.
func (e *Channel) Export(buf *[]byte, scft lrpt.SpacecraftParameters) bool {
	if !e.HasData {
		return false
	}

	e.Process(scft)
	*buf = make([]byte, e.Height*e.Width)

	index := 0
	for x := e.FirstSegment; x < e.LastSegment; x += 14 {
		for i := uint32(0); i < 8; i++ {
			for j := uint32(0); j < 14; j++ {
				if s := e.segments[x+j]; s != nil && s.IsValid() {
					copy((*buf)[index:], s.Lines[i][:])
				}
				index += 8 * 14
			}
		}
	}

	return true
}
