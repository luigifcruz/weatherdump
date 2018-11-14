package VIIRS

import (
	"encoding/binary"
	"fmt"
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

	for _, e := range dat { if e > max { max = e } }
	for _, e := range dat { if e < min { min = e } }

	return min, max
}

func NormalizeImage(data *[]byte) {
	u16 := ConvertToU16(*data)

	min, max := FindColorDepth(u16)

	fmt.Println("Min:", min, "Max:", max)

	for i, j := range u16 {
		u16[i] = uint16((((float64(j) - float64(min)) * float64(65535 - 0)) / (float64(max) - float64(min))) + 0.00)
	}

	*data = ConvertToByte(u16)
}