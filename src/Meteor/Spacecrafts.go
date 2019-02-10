package Meteor

type SpacecraftParameters struct {
	Filename   string
	FullName   string
	SignalName string
}

var Spacecrafts = map[uint8]SpacecraftParameters{
	000: {
		Filename:   "METEOR_MN2",
		FullName:   "Meteor",
		SignalName: "LRPT",
	},
}
