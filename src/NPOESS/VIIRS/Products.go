package VIIRS

import (
	"sort"
	"weather-dump/src/NPOESS/VIIRS/VIIRSFrames"
)

type Data struct {
	tempSegments [2047]Segment
	dataSegments []Segment
	outputFolder string
}

type Segment struct {
	APID   uint16
	header VIIRSFrames.FrameHeader
	body   [32]VIIRSFrames.FrameBody
}

func NewFillSegment(scanNumber uint32) *Segment {
	fillFrame := Segment{}
	fillFrame.header = *VIIRSFrames.NewFillFrameHeader(scanNumber)
	for i := 0; i < 32; i += 1 {
		fillFrame.body[i] = *VIIRSFrames.NewFillFrameBody()
	}
	return &fillFrame
}

func (e Data) GetChannelPackets(chAPID uint16) (uint32, uint32, map[uint32]*Segment) {
	var pkts []*Segment

	for i, segment := range e.dataSegments {
		if segment.APID == chAPID {
			pkts = append(pkts, &e.dataSegments[i])
		}
	}

	sort.SliceStable(pkts, func(i, j int) bool {
		return pkts[j].header.GetSequenceCount() < pkts[i].header.GetSequenceCount()
	})

	if len(pkts) == 0 {
		return 0, 0, nil
	}

	lastN := pkts[0].header.GetScanNumber()
	frstN := pkts[len(pkts)-1].header.GetScanNumber()
	realN := lastN - frstN + 1

	usedN := 0
	z := lastN

	buf := make(map[uint32]*Segment, realN)
	for i := uint32(0); i < realN; i += 1 {
		if lastN != pkts[usedN].header.GetScanNumber() {
			buf[lastN] = NewFillSegment(lastN)
		} else {
			buf[lastN] = pkts[usedN]
			usedN += 1
		}
		lastN -= 1
	}

	return frstN, z, buf
}

func (e Data) GetTimestamp(chAPID uint16) string {
	a, _, segments := e.GetChannelPackets(chAPID)
	for _, seg := range segments {
		if seg.header.IsValid() {
			return seg.header.GetDate()
		}
	}
	return segments[a].header.GetDate()
}
