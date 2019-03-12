package bismw

import (
	"fmt"
	"runtime"
	"sync"
	"weather-dump/src/protocols/lrpt"
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

// Export the image of the respective channel.
func (e *Channel) Compose() *[]byte {
	buf := make([]byte, e.count*8*e.width)

	threads := runtime.NumCPU()
	var wg sync.WaitGroup
	wg.Add(threads)

	for t := 0; t < threads; t++ {
		end := int(e.count) / threads * (t + 1)

		if t == threads-1 {
			end = int(e.count)
		}

		go func(wg *sync.WaitGroup, start, finish int) {
			defer wg.Done()
			for i := start; i < finish; i++ {
				e.lines[uint32(i)].RenderLine(&buf, int(i)*8*int(e.width))
			}
		}(&wg, int(e.count)/threads*t, end)
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
