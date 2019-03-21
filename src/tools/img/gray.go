package img

import (
	"image"
	"image/jpeg"
	"image/png"
	"math"
	"os"

	"github.com/luigifreitas/gofast"
)

type Gray struct {
	buf    *[]byte
	width  int
	height int
}

func NewGray(buf *[]byte, width, height int) Img {
	return &Gray{buf, width, height}
}

func (e *Gray) Flop() Img {
	finish := len(*e.buf) / e.width
	gofast.For(0, finish, 1, func(i int) {
		for p := 0; p < e.width/2; p++ {
			fp := p + (i * e.width)
			lp := e.width - p + (i * e.width) - 1

			f := (*e.buf)[lp]
			l := (*e.buf)[fp]

			(*e.buf)[fp] = f
			(*e.buf)[lp] = l
		}
	})
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

	gofast.For(0, totalPixels, 1, func(i int) {
		(*e.buf)[i] = nlvl[(*e.buf)[i]]
	})
	return e
}

func (e *Gray) Invert() Img {
	gofast.For(0, len(*e.buf), 1, func(i int) {
		(*e.buf)[i] = 255 - (*e.buf)[i]
	})
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
