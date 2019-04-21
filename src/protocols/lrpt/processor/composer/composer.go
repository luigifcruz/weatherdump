package composer

import (
	"fmt"
	"path/filepath"
	"sort"
	"weatherdump/src/img"
	"weatherdump/src/protocols/lrpt"
	"weatherdump/src/protocols/lrpt/processor/parser"
)

type Composer struct {
	pipeline         img.Pipeline
	scft             lrpt.SpacecraftParameters
	Equalize         bool
	ShortName        string
	FileName         string
	RequiredChannels []uint16
}

func (e *Composer) Register(pipeline img.Pipeline, scft lrpt.SpacecraftParameters) *Composer {
	e.pipeline = pipeline
	e.scft = scft
	return e
}

func (e Composer) Render(ch parser.List, outputFolder string) string {
	ch01 := ch[e.RequiredChannels[0]]
	ch02 := ch[e.RequiredChannels[1]]
	ch03 := ch[e.RequiredChannels[2]]

	// Check if required channels exist.
	if !ch01.HasData || !ch02.HasData || !ch03.HasData {
		//fmt.Println("[COM] Can't export component channel. Not all required channels are available.")
		return ""
	}

	outputName, _ := filepath.Abs(fmt.Sprintf("%s/%s_%s_COMP_%s_LRPT",
		outputFolder, e.scft.Filename, e.scft.SignalName, e.FileName))

	// Synchronize all channels scans.
	firstScan := make([]int, 3)
	lastScan := make([]int, 3)
	offset := make([]int, 3)

	firstScan[0], lastScan[0], offset[0] = ch01.GetBounds()
	firstScan[1], lastScan[1], offset[1] = ch02.GetBounds()
	firstScan[2], lastScan[2], offset[2] = ch03.GetBounds()

	ch01.SetBounds(MaxIntSlice(firstScan), MinIntSlice(lastScan), MaxIntSlice(offset))
	ch02.SetBounds(MaxIntSlice(firstScan), MinIntSlice(lastScan), MaxIntSlice(offset))
	ch03.SetBounds(MaxIntSlice(firstScan), MinIntSlice(lastScan), MaxIntSlice(offset))

	ch01.Process(e.scft)
	ch02.Process(e.scft)
	ch03.Process(e.scft)

	// Create output image struct.
	w, h := ch01.GetDimensions()
	finalBuf := make([]byte, w*h*4)

	for p := 3; p < len(finalBuf); p += 4 {
		finalBuf[p] = 0xFF
	}

	// Compose images and fill buffer.
	var buf []byte
	e.pipeline.Target(img.NewGray(&buf, w, h))
	e.pipeline.AddException("Invert", false)
	e.pipeline.AddException("Equalize", e.Equalize)

	ch01.Export(&buf, e.scft)
	e.pipeline.Process()
	for p := 1; p < len(finalBuf); p += 4 {
		finalBuf[p] = buf[p/4]
	}

	ch02.Export(&buf, e.scft)
	e.pipeline.Process()
	for p := 0; p < len(finalBuf); p += 4 {
		finalBuf[p] = buf[p/4]
	}

	ch03.Export(&buf, e.scft)
	e.pipeline.Process()
	for p := 2; p < len(finalBuf); p += 4 {
		finalBuf[p] = buf[p/4]
	}

	// Render and save the true-color image.
	e.pipeline.Target(img.NewRGBA(&finalBuf, w, h)).Export(outputName, 100)
	e.pipeline.ResetExceptions()
	return outputName
}

func MinIntSlice(v []int) int {
	sort.Ints(v)
	return v[0]
}

func MaxIntSlice(v []int) int {
	sort.Ints(v)
	return v[len(v)-1]
}
