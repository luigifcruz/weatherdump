package BISMW

type ChannelParameters struct {
	APID              uint16
	ChannelName       string
	SegmentWidth      int
	SegmentHeight     int
	FinalProductWidth uint32
}

var ChannelsIndex = [24]uint16{64, 65, 66, 67, 68, 69}

var ChannelsParameters = map[uint16]ChannelParameters{
	64: {APID: 64, ChannelName: "CH64", SegmentWidth: 8, SegmentHeight: 8, FinalProductWidth: 1578},
	65: {APID: 65, ChannelName: "CH65", SegmentWidth: 8, SegmentHeight: 8, FinalProductWidth: 1578},
	66: {APID: 66, ChannelName: "CH66", SegmentWidth: 8, SegmentHeight: 8, FinalProductWidth: 1578},
	67: {APID: 67, ChannelName: "CH67", SegmentWidth: 8, SegmentHeight: 8, FinalProductWidth: 1578},
	68: {APID: 68, ChannelName: "CH68", SegmentWidth: 8, SegmentHeight: 8, FinalProductWidth: 1578},
	69: {APID: 69, ChannelName: "CH69", SegmentWidth: 8, SegmentHeight: 8, FinalProductWidth: 1578},
}
