package parser

import "weather-dump/src/assets"

type List map[uint16]*Channel

var Channels = List{
	64: {APID: 64, ChannelName: "CH64", BlockDim: 8, Invert: false, FinalWidth: 1568},
	65: {APID: 65, ChannelName: "CH65", BlockDim: 8, Invert: false, FinalWidth: 1568},
	66: {APID: 66, ChannelName: "CH66", BlockDim: 8, Invert: false, FinalWidth: 1568},
	67: {APID: 67, ChannelName: "CH67", BlockDim: 8, Invert: false, FinalWidth: 1568},
	68: {APID: 68, ChannelName: "CH68", BlockDim: 8, Invert: true, FinalWidth: 1568},
	69: {APID: 69, ChannelName: "CH69", BlockDim: 8, Invert: false, FinalWidth: 1568},
}

var Manifest = assets.Manifest{
	64: {
		Name: "Channel 64",
	},
	65: {
		Name: "Channel 64",
	},
	66: {
		Name: "Channel 64",
	},
	67: {
		Name: "Channel 64",
	},
	68: {
		Name: "Channel 64",
	},
	69: {
		Name: "Channel 64",
	},
}
