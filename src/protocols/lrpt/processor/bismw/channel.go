package bismw

import (
	"fmt"
	"runtime"
	"sync"
	"weather-dump/src/protocols/lrpt"
	"weather-dump/src/tools/parallel"
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
	return &Channel{
		apid:  apid,
		lines: make(map[uint32]*Line),
		count: 0,
	}
}

// Compose the image of the respective channel.
func (e *Channel) Compose() *[]byte {
	buf := make([]byte, e.count*8*e.width)

	threads := runtime.NumCPU()
	var wg sync.WaitGroup
	wg.Add(threads)

	for s, f := range parallel.SerialRange(0, int(e.count)/14, threads) {
		go func(wg *sync.WaitGroup, start, finish int) {
			defer wg.Done()
			for i := start; i < finish; i++ {
				e.lines[uint32(i)].RenderLine(&buf, int(i)*8*int(e.width))
			}
		}(&wg, s, f)
	}

	wg.Wait()
	return &buf
}

func (e Channel) GetFileName() string {
	return e.fileName
}

func (e Channel) GetDimensions() (int, int) {
	return int(e.width), int(e.height)
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
