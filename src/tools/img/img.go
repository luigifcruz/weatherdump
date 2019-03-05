package img

type Img interface {
	Invert() Img
	Flop() Img
	Equalize() Img
	ExportPNG(string)
	ExportJPEG(string, int)
}
