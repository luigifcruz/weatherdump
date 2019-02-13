package BISMW

import "weather-dump/src/Meteor"

type Line struct {
	segments map[uint8]*Segment
}

func NewLine() *Line {
	e := Line{}
	e.segments = make(map[uint8]*Segment)
	return &e
}

func (e *Line) AddMCU(dat []byte) {
	segment := NewSegment(dat)
	e.segments[segment.GetMCUNumber()/14] = segment
}

func (e Line) RenderLine() []byte {
	var buf = make([]byte, 64*14*14)
	fillerSegment := [64 * 14]byte{}

	o := 0
	for y := 0; y < 8; y++ {
		for x := 0; x < 1568; x++ {
			var segment []byte

			if e.segments[uint8(x/112)] == nil {
				segment = fillerSegment[:]
				return buf
			} else {
				segment = e.segments[uint8(x/112)].RenderSegment()
			}

			buf[o] = segment[(y*112)+(x%112)]
			o++
		}
	}

	return buf
}

func (e Line) GetDate() Meteor.Time {
	return Meteor.Time{}
	//return e.segments[0].GetDate()
}
