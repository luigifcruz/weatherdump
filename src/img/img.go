package img

type Img interface {
	Invert() Img
	Flop() Img
	Equalize() Img
	ExportPNG(string, int) Img
	ExportJPEG(string, int) Img
}
