package composer

import (
	"fmt"
	"path/filepath"
	"sort"
	"weather-dump/src/img"
	"weather-dump/src/protocols/hrd"
	"weather-dump/src/protocols/hrd/processor/parser"
)

type Composer struct {
	pipeline         img.Pipeline
	scft             hrd.SpacecraftParameters
	ShortName        string
	FileName         string
	RequiredChannels []uint16
}

func (e *Composer) Register(pipeline img.Pipeline, scft hrd.SpacecraftParameters) *Composer {
	e.pipeline = pipeline
	e.scft = scft
	return e
}

func (e Composer) Render(ch parser.List, outputFolder string) {
	ch01 := ch[e.RequiredChannels[0]]
	ch02 := ch[e.RequiredChannels[1]]
	ch03 := ch[e.RequiredChannels[2]]

	// Check if required channels exist.
	if !ch01.HasData || !ch02.HasData || !ch03.HasData {
		//fmt.Println("[COM] Can't export component channel. Not all required channels are available.")
		return
	}

	outputName, _ := filepath.Abs(fmt.Sprintf("%s/%s_%s_COMP_%s_VIIRS_%s",
		outputFolder, e.scft.Filename, e.scft.SignalName, e.FileName, ch01.StartTime.GetZuluSafe()))

	// Synchronize all channels scans.
	firstScan := make([]int, 3)
	lastScan := make([]int, 3)

	firstScan[0], lastScan[0] = ch01.GetBounds()
	firstScan[1], lastScan[1] = ch02.GetBounds()
	firstScan[2], lastScan[2] = ch03.GetBounds()

	ch01.SetBounds(MaxIntSlice(firstScan), MinIntSlice(lastScan))
	ch02.SetBounds(MaxIntSlice(firstScan), MinIntSlice(lastScan))
	ch03.SetBounds(MaxIntSlice(firstScan), MinIntSlice(lastScan))

	ch01.Process(e.scft)
	ch02.Process(e.scft)
	ch03.Process(e.scft)

	// Create output image struct.
	w, h := ch01.GetDimensions()
	finalBuf := make([]byte, w*h*8)

	for p := 6; p < len(finalBuf); p += 8 {
		finalBuf[p+0] = 0xFF
		finalBuf[p+1] = 0xFF
	}

	// Compose images and fill buffer.
	var buf []byte
	e.pipeline.Target(img.NewGray16(&buf, w, h))
	e.pipeline.AddException("Invert", false)

	ch01.Export(&buf, ch, e.scft)
	e.pipeline.Process()
	for p := 2; p < len(finalBuf); p += 8 {
		finalBuf[p+0] = buf[(p/4)+0]
		finalBuf[p+1] = buf[(p/4)+1]
	}

	ch02.Export(&buf, ch, e.scft)
	e.pipeline.Process()
	for p := 0; p < len(finalBuf); p += 8 {
		finalBuf[p+0] = buf[(p/4)+0]
		finalBuf[p+1] = buf[(p/4)+1]
	}

	ch03.Export(&buf, ch, e.scft)
	e.pipeline.Process()
	for p := 4; p < len(finalBuf); p += 8 {
		finalBuf[p+0] = buf[(p/4)-1]
		finalBuf[p+1] = buf[(p/4)+0]
	}

	// Render and save the true-color image.
	e.pipeline.Target(img.NewRGBA64(&finalBuf, w, h)).Export(outputName, 100)
	e.pipeline.ResetExceptions()
}

func MinIntSlice(v []int) int {
	sort.Ints(v)
	return v[0]
}

func MaxIntSlice(v []int) int {
	sort.Ints(v)
	return v[len(v)-1]
}
