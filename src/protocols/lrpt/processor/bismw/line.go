package bismw

import "weather-dump/src/protocols/lrpt"

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
	var line [64 * 14 * 14]byte
	var segments [14][64 * 14]byte

	for i := 0; i < 14; i++ {
		if e.segments[uint8(i)] == nil {
			return line[:]
		}
		e.segments[uint8(i)].RenderSegment(&segments[i])
	}

	o := 0
	for y := 0; y < 8; y++ {
		for x := 0; x < 1568; x++ {
			line[o] = segments[uint8(x/112)][(y*112)+(x%112)]
			o++
		}
	}

	return line[:]
}

func (e Line) GetDate() lrpt.Time {
	return lrpt.Time{}
	//return e.segments[0].GetDate()
}
