package img

import (
	"encoding/binary"
	"image"
	"image/jpeg"
	"image/png"
	"math"
	"os"

	"github.com/luigifreitas/gofast"
)

type Gray16 struct {
	buf    *[]byte
	width  int
	height int
}

func NewGray16(buf *[]byte, width, height int) Img {
	return &Gray16{buf, width, height}
}

func (e *Gray16) Flop() Img {
	finish := len(*e.buf) / (e.width * 2)
	dw := e.width * 2

	gofast.For(0, finish, 1, func(i int) {
		for p := 0; p < e.width; p += 2 {
			fp := p + (i * dw)
			lp := dw - p + (i * dw) - 1

			l1 := (*e.buf)[fp]
			l2 := (*e.buf)[fp+1]
			f1 := (*e.buf)[lp-1]
			f2 := (*e.buf)[lp]

			(*e.buf)[fp] = f1
			(*e.buf)[fp+1] = f2
			(*e.buf)[lp-1] = l1
			(*e.buf)[lp] = l2
		}
	})
	return e
}

func (e *Gray16) Equalize() Img {
	var hist [65536]int
	var nlvl [65536]uint16

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
		nlvl[firstNonZero] = uint16(math.Max(0, math.Min(float64(frequencyCount)*pixelScale, 65535)))
	}

	gofast.For(0, totalPixels, 2, func(i int) {
		binary.BigEndian.PutUint16((*e.buf)[i:], nlvl[binary.BigEndian.Uint16((*e.buf)[i:])])
	})
	return e
}

func (e *Gray16) Invert() Img {
	gofast.For(0, len(*e.buf), 1, func(i int) {
		(*e.buf)[i] = 255 - (*e.buf)[i]
	})
	return e
}

func (e *Gray16) ExportPNG(outputFile string, quality int) Img {
	o, _ := os.Create(outputFile + ".png")
	defer o.Close()

	img := image.NewGray16(image.Rect(0, 0, e.width, e.height))
	img.Pix = *e.buf

	enc := &png.Encoder{
		CompressionLevel: png.DefaultCompression,
	}
	enc.Encode(o, img)
	return e
}

func (e *Gray16) ExportJPEG(outputFile string, quality int) Img {
	o, _ := os.Create(outputFile + ".jpeg")
	defer o.Close()

	img := image.NewGray16(image.Rect(0, 0, e.width, e.height))
	img.Pix = *e.buf

	var opt jpeg.Options
	opt.Quality = quality
	jpeg.Encode(o, img, &opt)
	return e
}
