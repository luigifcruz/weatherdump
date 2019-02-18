package BISMW

import (
	"fmt"
	"image"
	"image/png"
	"os"
	"weather-dump/src/Meteor"
	"weather-dump/src/imagery"
)

const maxFrameCount = 8192

// Channel struct.
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

// NewChannel instance.
func NewChannel(apid uint16) *Channel {
	e := Channel{}
	e.apid = apid
	e.lines = make(map[uint32]*Line)
	e.count = 0
	return &e
}

// Export the image of the respective channel.
func (e *Channel) Export(outputFolder string) {
	var buf []uint8
	for i := uint32(0); i < e.count; i++ {
		line := e.lines[i].RenderLine()
		buf = append(buf, line[:]...)
	}

	imagery.HistogramEqualizationU8(&buf)
	if e.parameters.Inversion {
		imagery.PixelInversionU8(&buf)
	}

	output, _ := os.Create(fmt.Sprintf("%s/%s.png", outputFolder, e.fileName))
	defer output.Close()
	s := image.NewGray(image.Rect(0, 0, int(e.width), int(e.height)))
	s.Pix = buf
	png.Encode(output, s)
}

// Fix the channel metadata.
func (e *Channel) Fix(scft Meteor.SpacecraftParameters) {
	e.parameters = ChannelsParameters[e.apid]
	e.startTime = e.lines[0].GetDate()
	e.endTime = e.lines[e.count/14].GetDate()
	e.fileName = fmt.Sprintf("%s_%s_BISMW_%s_%d", scft.Filename, scft.SignalName, e.parameters.ChannelName, e.startTime.GetMilliseconds())
	e.height = e.count * uint32(e.parameters.BlockDim) / 14
	e.width = e.parameters.FinalWidth
}
