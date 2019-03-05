package img

import (
	"image"
	"image/jpeg"
	"image/png"
	"math"
	"os"
	"runtime"
	"sync"
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

	lines := len(*e.buf) / (e.width * 2)
	for t := 0; t < e.threads; t++ {
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
		}(&wg, lines/e.threads*t, lines/e.threads*(t+1))
	}

	wg.Wait()
	return e
}

func (e *Gray) Equalize() Img {
	var hist, nlvl [256]int
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
		nlvl[firstNonZero] = int(math.Max(0, math.Min(float64(frequencyCount)*pixelScale, 255)))
	}

	var wg sync.WaitGroup
	wg.Add(e.threads)

	for t := 0; t < e.threads; t++ {
		go func(wg *sync.WaitGroup, start, finish int) {
			defer wg.Done()
			for p := start; p < finish; p++ {
				(*e.buf)[p] = uint8(nlvl[(*e.buf)[p]])
			}
		}(&wg, totalPixels/e.threads*t, totalPixels/e.threads*(t+1))
	}

	wg.Wait()
	return e
}

func (e *Gray) Invert() Img {
	var wg sync.WaitGroup
	wg.Add(e.threads)

	pixels := len(*e.buf)
	for t := 0; t < e.threads; t++ {
		go func(wg *sync.WaitGroup, start, finish int) {
			defer wg.Done()

			for p := start; p < finish; p++ {
				(*e.buf)[p] = 255 - (*e.buf)[p]
			}
		}(&wg, pixels/e.threads*t, pixels/e.threads*(t+1))
	}

	wg.Wait()
	return e
}

func (e *Gray) ExportPNG(outputFile string) {
	o, _ := os.Create(outputFile + ".png")
	defer o.Close()

	img := image.NewGray(image.Rect(0, 0, e.width, e.height))
	img.Pix = *e.buf

	png.Encode(o, img)
}

func (e *Gray) ExportJPEG(outputFile string, quality int) {
	o, _ := os.Create(outputFile + ".jpeg")
	defer o.Close()

	img := image.NewGray(image.Rect(0, 0, e.width, e.height))
	img.Pix = *e.buf

	var opt jpeg.Options
	opt.Quality = quality
	jpeg.Encode(o, img, &opt)
}
