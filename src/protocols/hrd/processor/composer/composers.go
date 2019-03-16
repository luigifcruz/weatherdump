package composer

import (
	"weather-dump/src/assets"
)

type List map[uint16]Composer

var Composers = List{
	000: {
		FileName:         "TRUECOLOR_M_CH",
		RequiredChannels: []uint16{800, 801, 802},
	},
	001: {
		FileName:         "NAT_COLOR_I_CH",
		RequiredChannels: []uint16{820, 819, 818},
	},
	002: {
		FileName:         "NAT_COLOR_M_CH",
		RequiredChannels: []uint16{808, 806, 801},
	},
}

var Manifest = assets.Manifest{
	000: {
		Name:        "True-Color",
		Description: "Moderate Resolution Channels True-Color Component",
	},
	001: {
		Name:        "Imagery Natural-Color",
		Description: "Imagery Channels Natural-Color Component",
	},
	002: {
		Name:        "Moderate Natural-Color",
		Description: "Moderate Resolution Natural-Color Component",
	},
}
