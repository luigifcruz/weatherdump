package Decoder

type Parameters struct {
	FrameSize          int
	FrameBits          int
	CodedFrameSize     int
	MinCorrelationBits uint
	SyncWordSize       int
	RsParityBlockSize  int
	RsBlocks           int
	SyncWords          [8]uint64
}

var Datalink = map[string]Parameters{
	"LRPT": {
		FrameSize:          1024,
		FrameBits:          (1024 * 8),
		CodedFrameSize:     ((1024 * 8) * 2),
		MinCorrelationBits: 46,
		SyncWordSize:       4,
		RsParityBlockSize:  (32 * 4),
		RsBlocks:           4,
		SyncWords: [8]uint64{
			0xfca2b63db00d9794,
			0x56fbd394daa4c1c2,
			0x035d49c24ff2686b,
			0xa9042c6b255b3e3d,
			0xfc51793e700e6b68,
			0xa9f7e368e558c2c1,
			0x03ae86c18ff19497,
			0x56081c971aa73d3e,
		},
	},
}
