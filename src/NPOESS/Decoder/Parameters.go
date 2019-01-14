package Decoder

type Parameters struct {
	FrameSize          int
	FrameBits          int
	CodedFrameSize     int
	MinCorrelationBits uint
	SyncWordSize       int
	RsParityBlockSize  int
	RsBlocks           int
	HritUw0            uint64
	HritUw2            uint64
}

var Datalink = map[string]Parameters{
	"HRD": {
		FrameSize:          1024,
		FrameBits:          (1024 * 8),
		CodedFrameSize:     ((1024 * 8) * 2),
		MinCorrelationBits: 46,
		SyncWordSize:       4,
		RsParityBlockSize:  (32 * 4),
		RsBlocks:           4,
		HritUw0:            0xfc4ef4fd0cc2df89,
		HritUw2:            0x25010b02f33d2076,
	},
}
