package segment

type Data struct {
	Header *Header
	Body   [32]Body
}

func NewSegment(header *Header) *Data {
	return &Data{
		Header: header,
	}
}

func NewFillSegment(scanNumber uint32) *Data {
	fillFrame := Data{Header: NewFillHeader(scanNumber)}
	for i := 0; i < 32; i++ {
		fillFrame.Body[i] = *NewFillBody()
	}
	return &fillFrame
}
