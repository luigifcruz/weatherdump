package block

import (
	"weather-dump/src/protocols/lrpt"
)

type Data struct {
	Segments map[uint8]*Segment
}

func New() *Data {
	return &Data{make(map[uint8]*Segment)}
}

func (e *Data) AddMCU(dat []byte) {
	segment := NewSegment(dat)
	e.Segments[segment.GetMCUNumber()/14] = segment
}

func (e Data) GetDate() lrpt.Time {
	for _, segment := range e.Segments {
		return segment.GetDate()
	}
	return lrpt.Time{}
}
