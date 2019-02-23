package interfaces

type ProcessorMakers map[string]func(string) Processor
type Processor interface {
	Work(string)
	ExportAll(string)
}

type DecoderMakers map[string]func(string) Decoder
type Decoder interface {
	Work(string, string)
}
