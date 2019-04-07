package composer

import (
	"weather-dump/src/protocols/helpers"
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

var Manifest = helpers.ManifestList{
	000: {
		Name:        "True-Color",
		Description: "Moderate RGB Composite",
		Activated:   true,
	},
	001: {
		Name:        "Ima. Nat-Color",
		Description: "Imagery Natural Composite",
		Activated:   true,
	},
	002: {
		Name:        "Mod. Nat-Color",
		Description: "Moderate Natural Composite",
		Activated:   true,
	},
}
