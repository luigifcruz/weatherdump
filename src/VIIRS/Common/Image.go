package VIIRS

import (
	"encoding/binary"
)

func ConvertToU16(data []byte) []uint16 {
	var buf []uint16
	for i := 0; i < len(data); i += 2 {
		buf = append(buf, binary.BigEndian.Uint16(data[i:]))
	}
	return buf
}

func ConvertToByte(data []uint16) []byte {
	var buf []byte
	bb := make([]byte, 2)

	for _, d := range data {
		binary.BigEndian.PutUint16(bb, d)
		buf = append(buf, bb...)
	}

	return buf
}

func FindColorDepth(dat []uint16) (uint16, uint16) {
	max := uint16(0)
	min := uint16(65535)

	for _, e := range dat {
		if e > max {
			max = e
		}
	}
	for _, e := range dat {
		if e < min {
			min = e
		}
	}

	return min, max
}
