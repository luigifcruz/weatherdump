package img

import (
	"encoding/binary"
	"image"
	"image/jpeg"
	"image/png"
	"math"
	"os"
	"runtime"
	"sync"
)

type Gray16 struct {
	buf     *[]byte
	width   int
	height  int
	threads int
}

func NewGray16(buf *[]byte, width, height int) Img {
	return &Gray16{buf, width, height, runtime.NumCPU()}
}

func (e *Gray16) Flop() Img {
	var wg sync.WaitGroup
	wg.Add(e.threads)

	lines := len(*e.buf) / (e.width * 2)
	for t := 0; t < e.threads; t++ {
		go func(wg *sync.WaitGroup, start, finish int) {
			defer wg.Done()

			dw := e.width * 2
			for l := start; l < finish; l++ {
				for p := 0; p < e.width; p += 2 {
					fp := p + (l * dw)
					lp := dw - p + (l * dw) - 1

					l1 := (*e.buf)[fp]
					l2 := (*e.buf)[fp+1]
					f1 := (*e.buf)[lp-1]
					f2 := (*e.buf)[lp]

					(*e.buf)[fp] = f1
					(*e.buf)[fp+1] = f2
					(*e.buf)[lp-1] = l1
					(*e.buf)[lp] = l2
				}
			}
		}(&wg, lines/e.threads*t, lines/e.threads*(t+1))
	}

	wg.Wait()
	return e
}

func (e *Gray16) Equalize() Img {
	var hist, nlvl [65536]int
	totalPixels := len(*e.buf)

	for p := 0; p < totalPixels; p += 2 {
		hist[binary.BigEndian.Uint16((*e.buf)[p:])]++
	}

	firstNonZero := 0
	for hist[firstNonZero] == 0 {
		firstNonZero++
	}

	if hist[firstNonZero] == totalPixels/2 {
		for p := 0; p < totalPixels; p += 2 {
			binary.BigEndian.PutUint16((*e.buf)[p:], uint16(totalPixels/2))
		}
		return e
	}

	pixelScale := float64(len(hist)-1) / float64((totalPixels/2)-hist[firstNonZero])
	firstNonZero++

	frequencyCount := 0
	for ; firstNonZero < len(hist); firstNonZero++ {
		frequencyCount += hist[firstNonZero]
		nlvl[firstNonZero] = int(math.Max(0, math.Min(float64(frequencyCount)*pixelScale, 65535)))
	}

	var wg sync.WaitGroup
	wg.Add(e.threads)

	for t := 0; t < e.threads; t++ {
		go func(wg *sync.WaitGroup, start, finish int) {
			defer wg.Done()
			for p := start; p < finish; p += 2 {
				binary.BigEndian.PutUint16((*e.buf)[p:], uint16(nlvl[binary.BigEndian.Uint16((*e.buf)[p:])]))
			}
		}(&wg, totalPixels/e.threads*t, totalPixels/e.threads*(t+1))
	}

	wg.Wait()
	return e
}

func (e *Gray16) Invert() Img {
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

func (e *Gray16) ExportPNG(outputFile string) {
	o, _ := os.Create(outputFile + ".png")
	defer o.Close()

	img := image.NewGray16(image.Rect(0, 0, e.width, e.height))
	img.Pix = *e.buf

	png.Encode(o, img)
}

func (e *Gray16) ExportJPEG(outputFile string, quality int) {
	o, _ := os.Create(outputFile + ".jpeg")
	defer o.Close()

	img := image.NewGray16(image.Rect(0, 0, e.width, e.height))
	img.Pix = *e.buf

	var opt jpeg.Options
	opt.Quality = quality
	jpeg.Encode(o, img, &opt)
}
