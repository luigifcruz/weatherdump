package interfaces

import (
	"weather-dump/src/assets"
)

type ProcessorMakers map[string]func(string) Processor
type Processor interface {
	Work(string)
	Export(*assets.ExportDelegate, string)
}

type DecoderMakers map[string]func(string) Decoder
type Decoder interface {
	Work(string, string, *bool)
}
