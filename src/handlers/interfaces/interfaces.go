package interfaces

import "weather-dump/src/tools/img"

type ProcessorMakers map[string]func(string) Processor
type Processor interface {
	Work(string)
	Export(string, img.Pipeline)
}

type DecoderMakers map[string]func(string) Decoder
type Decoder interface {
	Work(string, string, *bool)
}
