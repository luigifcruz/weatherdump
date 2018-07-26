package VIIRS

type ChannelParameters struct {
	APID uint16
	ChannelName string
	AggregationZoneWidth [6]int
	AggregationZoneHeight int
	BowTieHeight [6]int
	FinalProductWidth int
}

var ChannelsParameters = map[uint16]ChannelParameters{
	801: {
		APID: 801,
		ChannelName: "M04",
		AggregationZoneWidth: [6]int{640, 368, 592, 592, 368, 640},
		AggregationZoneHeight: 15,
		BowTieHeight: [6]int{3, 1, 0, 0, 1, 3},
		FinalProductWidth: 3200,
	},
	805: {
		APID: 805,
		ChannelName: "M06",
		AggregationZoneWidth: [6]int{640, 368, 592, 592, 368, 640},
		AggregationZoneHeight: 15,
		BowTieHeight: [6]int{3, 1, 0, 0, 1, 3},
		FinalProductWidth: 3200,
	},
	807: {
		APID: 807,
		ChannelName: "M09",
		AggregationZoneWidth: [6]int{640, 368, 592, 592, 368, 640},
		AggregationZoneHeight: 15,
		BowTieHeight: [6]int{3, 1, 0, 0, 1, 3},
		FinalProductWidth: 3200,
	},
	808: {
		APID: 808,
		ChannelName: "M10",
		AggregationZoneWidth: [6]int{640, 368, 592, 592, 368, 640},
		AggregationZoneHeight: 15,
		BowTieHeight: [6]int{3, 1, 0, 0, 1, 3},
		FinalProductWidth: 3200,
	},
	812: {
		APID: 812,
		ChannelName: "M12",
		AggregationZoneWidth: [6]int{640, 368, 592, 592, 368, 640},
		AggregationZoneHeight: 15,
		BowTieHeight: [6]int{3, 1, 0, 0, 1, 3},
		FinalProductWidth: 3200,
	},
	814: {
		APID: 814,
		ChannelName: "M16",
		AggregationZoneWidth: [6]int{640, 368, 592, 592, 368, 640},
		AggregationZoneHeight: 15,
		BowTieHeight: [6]int{3, 1, 0, 0, 1, 3},
		FinalProductWidth: 3200,
	},
	815: {
		APID: 815,
		ChannelName: "M15",
		AggregationZoneWidth: [6]int{640, 368, 592, 592, 368, 640},
		AggregationZoneHeight: 15,
		BowTieHeight: [6]int{3, 1, 0, 0, 1, 3},
		FinalProductWidth: 3200,
	},
	816: {
		APID: 816,
		ChannelName: "M14",
		AggregationZoneWidth: [6]int{640, 368, 592, 592, 368, 640},
		AggregationZoneHeight: 15,
		BowTieHeight: [6]int{3, 1, 0, 0, 1, 3},
		FinalProductWidth: 3200,
	},
	818: {
		APID: 818,
		ChannelName: "I01",
		AggregationZoneWidth: [6]int{1280, 736, 1184, 1184, 736, 1280},
		AggregationZoneHeight: 31,
		BowTieHeight: [6]int{6, 2, 0, 0, 2, 6},
		FinalProductWidth: 6400,
	},
}