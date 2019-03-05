package bismw

import (
	"fmt"
	"path/filepath"
	"weather-dump/src/protocols/lrpt"
	"weather-dump/src/tools/img"
)

// Channel struct.
type Channel struct {
	apid       uint16
	parameters ChannelParameters
	fileName   string
	height     uint32
	width      uint32
	startTime  lrpt.Time
	endTime    lrpt.Time
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

	i := img.NewGray(&buf, int(e.width), int(e.height)).Equalize()
	if e.parameters.Inversion {
		i.Invert()
	}

	outputName, _ := filepath.Abs(fmt.Sprintf("%s/%s", outputFolder, e.fileName))
	i.ExportPNG(outputName)
}

// Fix the channel metadata.
func (e *Channel) Fix(scft lrpt.SpacecraftParameters) {
	e.parameters = ChannelsParameters[e.apid]
	e.startTime = e.lines[0].GetDate()
	e.endTime = e.lines[e.count/14].GetDate()
	e.fileName = fmt.Sprintf("%s_%s_BISMW_%s_%d", scft.Filename, scft.SignalName, e.parameters.ChannelName, e.startTime.GetMilliseconds())
	e.height = e.count * uint32(e.parameters.BlockDim) / 14
	e.width = e.parameters.FinalWidth
}
