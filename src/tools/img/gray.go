package img

import (
	"image"
	"image/jpeg"
	"image/png"
	"math"
	"os"
	"runtime"
	"sync"
	"weather-dump/src/tools/parallel"
)

type Gray struct {
	buf     *[]byte
	width   int
	height  int
	threads int
}

func NewGray(buf *[]byte, width, height int) Img {
	return &Gray{buf, width, height, runtime.NumCPU()}
}

func (e *Gray) Flop() Img {
	var wg sync.WaitGroup
	wg.Add(e.threads)

	for s, f := range parallel.SerialRange(0, len(*e.buf)/(e.width), e.threads) {
		go func(wg *sync.WaitGroup, start, finish int) {
			defer wg.Done()

			for l := start; l < finish; l++ {
				for p := 0; p < e.width/2; p++ {
					fp := p + (l * e.width)
					lp := e.width - p + (l * e.width) - 1

					f := (*e.buf)[lp]
					l := (*e.buf)[fp]

					(*e.buf)[fp] = f
					(*e.buf)[lp] = l
				}
			}
		}(&wg, s, f)
	}

	wg.Wait()
	return e
}

func (e *Gray) Equalize() Img {
	var hist [256]int
	var nlvl [256]uint8

	totalPixels := len(*e.buf)

	for p := 0; p < totalPixels; p++ {
		hist[(*e.buf)[p]]++
	}

	firstNonZero := 0
	for hist[firstNonZero] == 0 {
		firstNonZero++
	}

	if hist[firstNonZero] == totalPixels {
		for p := 0; p < totalPixels; p++ {
			(*e.buf)[p] = uint8(totalPixels)
		}
		return e
	}

	pixelScale := float64(len(hist)-1) / float64(totalPixels-hist[firstNonZero])
	firstNonZero++

	frequencyCount := 0
	for ; firstNonZero < len(hist); firstNonZero++ {
		frequencyCount += hist[firstNonZero]
		nlvl[firstNonZero] = uint8(math.Max(0, math.Min(float64(frequencyCount)*pixelScale, 255)))
	}

	var wg sync.WaitGroup
	wg.Add(e.threads)

	for s, f := range parallel.SerialRange(0, totalPixels, e.threads) {
		go func(wg *sync.WaitGroup, start, finish int) {
			defer wg.Done()
			for p := start; p < finish; p++ {
				(*e.buf)[p] = nlvl[(*e.buf)[p]]
			}
		}(&wg, s, f)
	}

	wg.Wait()
	return e
}

func (e *Gray) Invert() Img {
	var wg sync.WaitGroup
	wg.Add(e.threads)

	for s, f := range parallel.SerialRange(0, len(*e.buf), e.threads) {
		go func(wg *sync.WaitGroup, start, finish int) {
			defer wg.Done()

			for p := start; p < finish; p++ {
				(*e.buf)[p] = 255 - (*e.buf)[p]
			}
		}(&wg, s, f)
	}

	wg.Wait()
	return e
}

func (e *Gray) ExportPNG(outputFile string, quality int) Img {
	o, _ := os.Create(outputFile + ".png")
	defer o.Close()

	img := image.NewGray(image.Rect(0, 0, e.width, e.height))
	img.Pix = *e.buf

	enc := &png.Encoder{
		CompressionLevel: png.DefaultCompression,
	}
	enc.Encode(o, img)
	return e
}

func (e *Gray) ExportJPEG(outputFile string, quality int) Img {
	o, _ := os.Create(outputFile + ".jpeg")
	defer o.Close()

	img := image.NewGray(image.Rect(0, 0, e.width, e.height))
	img.Pix = *e.buf

	var opt jpeg.Options
	opt.Quality = quality
	jpeg.Encode(o, img, &opt)
	return e
}
