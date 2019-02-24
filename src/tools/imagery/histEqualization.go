package imagery

import "math"

// HistogramEqualizationU16 of a 16-bit grayscale image.
func HistogramEqualizationU16(img *[]uint16) {
	var hist, nlvl [65536]int
	totalPixels := len(*img)

	for p := 0; p < totalPixels; p++ {
		hist[(*img)[p]]++
	}

	firstNonZero := 0
	for hist[firstNonZero] == 0 {
		firstNonZero++
	}

	if hist[firstNonZero] == totalPixels {
		for p := 0; p < totalPixels; p++ {
			(*img)[p] = uint16(totalPixels)
		}
		return
	}

	pixelScale := float64(len(hist)-1) / float64(totalPixels-hist[firstNonZero])
	firstNonZero++

	frequencyCount := 0
	for ; firstNonZero < len(hist); firstNonZero++ {
		frequencyCount += hist[firstNonZero]
		nlvl[firstNonZero] = int(math.Max(0, math.Min(float64(frequencyCount)*pixelScale, 65535)))
	}

	for p := 0; p < totalPixels; p++ {
		(*img)[p] = uint16(nlvl[(*img)[p]])
	}
}

// HistogramEqualizationU8 of a 16-bit grayscale image.
func HistogramEqualizationU8(img *[]byte) {
	var hist, nlvl [256]int
	totalPixels := len(*img)

	for p := 0; p < totalPixels; p++ {
		hist[(*img)[p]]++
	}

	firstNonZero := 0
	for hist[firstNonZero] == 0 {
		firstNonZero++
	}

	if hist[firstNonZero] == totalPixels {
		for p := 0; p < totalPixels; p++ {
			(*img)[p] = uint8(totalPixels)
		}
		return
	}

	pixelScale := float64(len(hist)-1) / float64(totalPixels-hist[firstNonZero])
	firstNonZero++

	frequencyCount := 0
	for ; firstNonZero < len(hist); firstNonZero++ {
		frequencyCount += hist[firstNonZero]
		nlvl[firstNonZero] = int(math.Max(0, math.Min(float64(frequencyCount)*pixelScale, 255)))
	}

	for p := 0; p < totalPixels; p++ {
		(*img)[p] = uint8(nlvl[(*img)[p]])
	}
}
