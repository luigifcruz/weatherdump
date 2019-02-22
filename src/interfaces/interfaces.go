package interfaces

type Processor interface {
	Work(string)
	ExportAll(string)
}

type Decoder interface {
	Work(string, string)
}
