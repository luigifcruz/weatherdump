package BISMW

import (
	"fmt"
	"image"
	"image/jpeg"
	"os"
	"weather-dump/src/Meteor"
)

const maxFrameCount = 8192

type Channel struct {
	apid       uint16
	parameters ChannelParameters
	fileName   string
	height     uint32
	width      uint32
	startTime  Meteor.Time
	endTime    Meteor.Time
	lines      map[uint32]*Line
	count      uint32
	lastFrame  uint32
}

func NewChannel(apid uint16) *Channel {
	e := Channel{}
	e.apid = apid
	e.lines = make(map[uint32]*Line)
	e.count = 0
	return &e
}

func (e *Channel) Fix() {
	e.parameters = ChannelsParameters[e.apid]

	//e.startTime = e.segments[e.start].header.GetDate()
	//e.endTime = e.segments[e.end].header.GetDate()
	//e.fileName = fmt.Sprintf("%s_%s_BISMW_%s_%s", scft.Filename, scft.SignalName, e.parameters.ChannelName, e.startTime.GetZulu())
	e.height = e.count * uint32(e.parameters.SegmentHeight) / 14
	e.width = e.parameters.FinalProductWidth

	var final []byte
	for i := uint32(0); i < e.count; i++ {
		buf := e.lines[i].ExportLine()
		final = append(final, buf[:]...)
	}
	name := fmt.Sprintf("./out.jpeg")
	output, err := os.Create(name)
	if err != nil {
		fmt.Println(err)
	}
	defer output.Close()
	fmt.Println(len(final))
	s := image.NewGray(image.Rect(0, 0, 1568, len(final)/1568/14))
	s.Pix = final
	jpeg.Encode(output, s, nil)
}
