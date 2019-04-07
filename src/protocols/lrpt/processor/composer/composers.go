package composer

import (
	"weather-dump/src/protocols/helpers"
)

type List map[uint16]Composer

var Composers = List{
	000: {
		FileName:         "FALSECOLOR",
		RequiredChannels: []uint16{65, 68, 64},
		Equalize:         false,
	},
	001: {
		FileName:         "TRUECOLOR",
		RequiredChannels: []uint16{65, 66, 64},
		Equalize:         true,
	},
}

var Manifest = helpers.ManifestList{
	000: {
		Name:        "False Color",
		Description: "False Color RGB Composite",
		Activated:   true,
	},
	001: {
		Name:        "True-Color",
		Description: "True Color RGB Composite",
		Activated:   true,
	},
}
