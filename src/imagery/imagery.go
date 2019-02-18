package imagery

import (
	"sync"
)

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
	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()

		for p := 0; p < len(*img)/2; p++ {
			(*img)[p] = 255 - (*img)[p]
		}
	}()

	go func() {
		defer wg.Done()

		for p := len(*img) / 2; p < len(*img); p++ {
			(*img)[p] = 255 - (*img)[p]
		}
	}()

	wg.Wait()
}
