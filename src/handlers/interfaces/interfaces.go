package interfaces

import (
	"weatherdump/src/img"
	"weatherdump/src/protocols/helpers"
)

type ProcessorMakers map[string]func(string, *helpers.ProcessingManifest) Processor
type Processor interface {
	Work(string)
	Export(string, img.Pipeline)
	GetProductsManifest() helpers.ProcessingManifest
}

type DecoderMakers map[string]map[string]func(string) Decoder
type Decoder interface {
	Work(string, string, chan bool)
}
