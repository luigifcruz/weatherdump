package imagery

// FlopU16 a 16-bits grayscale image.
func FlopU16(img *[]uint16, w int) {
	buf := make([]uint16, len(*img))
	for p := 0; p < len(*img)-w; p++ {
		buf[p] = (*img)[(w-(p%w))+((p/w)*w)]
	}
	*img = buf
}

// PixelInversionU8 of a 8-bits grayscale image.
func PixelInversionU8(img *[]byte) {
	for p := 0; p < len(*img); p++ {
		(*img)[p] = 255 - (*img)[p]
	}
}
