package img

import (
	"image"
	"image/jpeg"
	"image/png"
	"os"
)

type RGBA struct {
	buf    *[]byte
	width  int
	height int
}

func NewRGBA(buf *[]byte, width, height int) Img {
	return &RGBA{buf, width, height}
}

func (e *RGBA) Flop() Img {
	return e
}

func (e *RGBA) Equalize() Img {
	return e
}

func (e *RGBA) Invert() Img {
	return e
}

func (e *RGBA) ExportPNG(outputFile string, quality int) Img {
	o, _ := os.Create(outputFile + ".png")
	defer o.Close()

	img := image.NewRGBA(image.Rect(0, 0, e.width, e.height))
	img.Pix = *e.buf

	enc := &png.Encoder{
		CompressionLevel: png.DefaultCompression,
	}
	enc.Encode(o, img)
	return e
}

func (e *RGBA) ExportJPEG(outputFile string, quality int) Img {
	o, _ := os.Create(outputFile + ".jpeg")
	defer o.Close()

	img := image.NewRGBA(image.Rect(0, 0, e.width, e.height))
	img.Pix = *e.buf

	var opt jpeg.Options
	opt.Quality = quality
	jpeg.Encode(o, img, &opt)
	return e
}
