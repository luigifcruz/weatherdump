package handlers

import (
	"weather-dump/src/handlers/interfaces"
	npoessDecoder "weather-dump/src/protocols/hrd/decoder"
	npoessProcessor "weather-dump/src/protocols/hrd/processor"
	meteorDecoder "weather-dump/src/protocols/lrpt/decoder"
	meteorProcessor "weather-dump/src/protocols/lrpt/processor"
)

// AvailableDecoders shows the currently available decoders for this build.
var AvailableDecoders = interfaces.DecoderMakers{
	"lrpt": meteorDecoder.NewDecoder,
	"hrd":  npoessDecoder.NewDecoder,
}

// AvailableProcessors shows the currently available processors for this build.
var AvailableProcessors = interfaces.ProcessorMakers{
	"lrpt": meteorProcessor.NewProcessor,
	"hrd":  npoessProcessor.NewProcessor,
}
