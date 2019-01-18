package VIIRS

import (
	"encoding/binary"
	"sort"
)

func MinIntSlice(v []int) int {
	sort.Ints(v)
	return v[0]
}

func MaxIntSlice(v []int) int {
	sort.Ints(v)
	return v[len(v)-1]
}

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
