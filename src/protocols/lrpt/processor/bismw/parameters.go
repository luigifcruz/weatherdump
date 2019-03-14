package bismw

type ChannelParameters struct {
	APID        uint16
	ChannelName string
	BlockDim    int
	Invert      bool
	FinalWidth  uint32
}

var ChannelsIndex = [24]uint16{64, 65, 66, 67, 68, 69}

var ChannelsParameters = map[uint16]ChannelParameters{
	64: {APID: 64, ChannelName: "CH64", BlockDim: 8, Invert: false, FinalWidth: 1568},
	65: {APID: 65, ChannelName: "CH65", BlockDim: 8, Invert: false, FinalWidth: 1568},
	66: {APID: 66, ChannelName: "CH66", BlockDim: 8, Invert: false, FinalWidth: 1568},
	67: {APID: 67, ChannelName: "CH67", BlockDim: 8, Invert: false, FinalWidth: 1568},
	68: {APID: 68, ChannelName: "CH68", BlockDim: 8, Invert: true, FinalWidth: 1568},
	69: {APID: 69, ChannelName: "CH69", BlockDim: 8, Invert: false, FinalWidth: 1568},
}
