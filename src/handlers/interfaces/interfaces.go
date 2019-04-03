package interfaces

import (
	"weather-dump/src/assets"
	"weather-dump/src/img"
)

type ProcessorMakers map[string]func(string) Processor
type Processor interface {
	Work(string)
	Export(string, img.Pipeline, assets.ProcessingManifest)
	GetProductsManifest() assets.ProcessingManifest
}

type DecoderMakers map[string]map[string]func(string) Decoder
type Decoder interface {
	Work(string, string, chan bool)
}

func WatchFor(signal chan bool, method func() bool) {
	for {
		select {
		case <-signal:
			return
		default:
			if method() {
				return
			}
		}
	}
}
