package VIIRS

type SpacecraftParameters struct {
	Filename string
	FullName string
	SignalName string
}

var Spacecrafts = map[uint8]SpacecraftParameters{
	159: {
		Filename: "NOAA20",
		FullName: "Joint Polar Satellite System",
		SignalName: "HRD",
	},
	157: {
		Filename: "NPP1",
		FullName: "Suomi National Polar-orbiting Partnership",
		SignalName: "HRD",
	},
}