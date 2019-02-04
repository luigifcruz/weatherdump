package BISMW

type Line struct {
	segments map[uint8]*Segment
}

func NewLine() *Line {
	e := Line{}
	e.segments = make(map[uint8]*Segment)
	return &e
}
