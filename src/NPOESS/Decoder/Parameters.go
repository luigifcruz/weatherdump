package Decoder

type Parameters struct {
	FrameSize          int
	FrameBits          int
	CodedFrameSize     int
	MinCorrelationBits uint
	SyncWordSize       int
	RsParityBlockSize  int
	RsBlocks           int
	HrdUw0             uint64
	HrdUw1             uint64
	HrdUw2             uint64
	HrdUw3             uint64
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
		HrdUw0:             0xfc4ef4fd0cc2df89,
		HrdUw1:             0x56275254a66b45ec,
		HrdUw2:             0x03b10b02f33d2076,
		HrdUw3:             0xa9d8adab5994ba89,
	},
}
