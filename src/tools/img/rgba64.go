package img

import (
	"image"
	"image/jpeg"
	"image/png"
	"os"
	"runtime"
)

type RGBA64 struct {
	buf     *[]byte
	width   int
	height  int
	threads int
}

func NewRGBA64(buf *[]byte, width, height int) Img {
	return &RGBA64{buf, width, height, runtime.NumCPU()}
}

func (e *RGBA64) Flop() Img {
	return e
}

func (e *RGBA64) Equalize() Img {
	return e
}

func (e *RGBA64) Invert() Img {
	return e
}

func (e *RGBA64) ExportPNG(outputFile string, quality int) Img {
	o, _ := os.Create(outputFile + ".png")
	defer o.Close()

	img := image.NewRGBA64(image.Rect(0, 0, e.width, e.height))
	img.Pix = *e.buf

	enc := &png.Encoder{
		CompressionLevel: png.DefaultCompression,
	}
	enc.Encode(o, img)
	return e
}

func (e *RGBA64) ExportJPEG(outputFile string, quality int) Img {
	o, _ := os.Create(outputFile + ".jpeg")
	defer o.Close()

	img := image.NewRGBA64(image.Rect(0, 0, e.width, e.height))
	img.Pix = *e.buf

	var opt jpeg.Options
	opt.Quality = quality
	jpeg.Encode(o, img, &opt)
	return e
}
