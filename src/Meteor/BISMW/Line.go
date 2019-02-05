package BISMW

type Line struct {
	segments map[uint8]*Segment
}

func NewLine() *Line {
	e := Line{}
	e.segments = make(map[uint8]*Segment)
	return &e
}

func (e Line) ExportLine() [64 * 14 * 14]byte {
	img := [64 * 14 * 14]byte{}

	o := 0
	for y := 0; y < 8; y++ {
		for x := 0; x < 1568; x++ {
			if e.segments[uint8(x/112)] == nil {
				return [64 * 14 * 14]byte{}
			}
			//fmt.Println(uint8(x/112), o, y*112+x-((x/112)*112))
			segment := e.segments[uint8(x/112)].ExportSegment()
			img[o] = segment[y*112+x-((x/112)*112)]
			o++
		}
	}

	return img
}
