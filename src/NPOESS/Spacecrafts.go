package NPOESS

type SpacecraftParameters struct {
	Filename          string
	FullName          string
	SignalName        string
	TrueColorChannels [3]uint16
}

var Spacecrafts = map[uint8]SpacecraftParameters{
	159: {
		Filename:          "NOAA20",
		FullName:          "Joint Polar Satellite System",
		SignalName:        "HRD",
		TrueColorChannels: [3]uint16{800, 801, 802},
	},
	157: {
		Filename:          "NPP1",
		FullName:          "Suomi National Polar-orbiting Partnership",
		SignalName:        "HRD",
		TrueColorChannels: [3]uint16{800, 801, 802},
	},
}
