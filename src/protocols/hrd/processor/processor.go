package processor

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"weatherdump/src/ccsds"
	"weatherdump/src/ccsds/frames"
	"weatherdump/src/handlers/interfaces"
	"weatherdump/src/img"
	"weatherdump/src/protocols/helpers"
	"weatherdump/src/protocols/hrd"
	"weatherdump/src/protocols/hrd/processor/composer"
	"weatherdump/src/protocols/hrd/processor/parser"

	"github.com/fatih/color"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{}

const frameSize = 892

type Worker struct {
	ccsds    *ccsds.Worker
	scid     uint8
	manifest helpers.ProcessingManifest
	channels parser.List
}

func NewProcessor(uuid string, manifest *helpers.ProcessingManifest) interfaces.Processor {
	e := Worker{
		ccsds:    ccsds.New(),
		channels: parser.New(),
	}

	if manifest == nil {
		e.manifest = e.GetProductsManifest()
	} else {
		e.manifest = *manifest
	}

	e.manifest.Register("hrd", uuid)

	return &e
}

func (e *Worker) Work(inputFile string) {
	color.Yellow("[PRC] WARNING! This processor is currently in BETA development state.")
	scidStat := [256]int{}

	file, _ := ioutil.ReadFile(inputFile)
	for i := len(file); i > 0; i -= frameSize {
		f := frames.NewTransferFrame(file[(len(file) - i):])
		p := frames.NewMultiplexingFrame(ccsds.Version["HRD"], f.GetMPDU())

		if f.IsReplay() && p.IsValid() {
			scidStat[f.GetSCID()]++
			switch f.GetVCID() {
			case 16:
				e.ccsds.ParseMPDU(*p) // VCID 16 Parser (VIIRS)
			}
		}
	}

	for _, packet := range e.ccsds.GetSpacePackets() {
		if packet.GetAPID() >= 800 && packet.GetAPID() <= 823 {
			e.channels[packet.GetAPID()].Parse(packet)
		}
	}

	e.scid = uint8(helpers.MaxIntSlice(scidStat[:]))
	fmt.Printf("[PRC] Decoded %d packets from VCID 16.\n", len(e.ccsds.GetSpacePackets()))
}

func (e *Worker) Export(outputPath string, wf img.Pipeline) {
	fmt.Printf("[PRC] Exporting VIIRS science products.\n")
	e.manifest.Start()

	re := helpers.CaptureOutput(func() {
		for _, apid := range e.manifest.Parser.Parse() {
			ch := e.channels[apid]

			var buf []byte
			if ch.Export(&buf, e.channels, hrd.Spacecrafts[e.scid]) {
				w, h := ch.GetDimensions()
				outputName, _ := filepath.Abs(fmt.Sprintf("%s/%s", outputPath, ch.FileName))

				wf.AddException("Invert", ch.Invert)
				wf.Target(img.NewGray16(&buf, w, h)).Process().Export(outputName, 100)
				wf.ResetExceptions()

				e.manifest.Parser[apid].FileName(outputName)
			}

			e.manifest.ParserCompleted(apid)
		}

		for _, code := range e.manifest.Composer.Parse() {
			c := composer.Composers[uint16(code)]
			outputName := c.Register(wf, hrd.Spacecrafts[e.scid]).Render(e.channels, outputPath)
			e.manifest.Composer[code].FileName(outputName)
			e.manifest.ComposerCompleted(code)
		}
	})

	e.channels = make(parser.List)
	e.ccsds = nil

	e.manifest.Stop(re)
	color.Green("[PRC] Done! All products and components were saved.")
}

func (e Worker) GetProductsManifest() helpers.ProcessingManifest {
	return helpers.ProcessingManifest{
		Parser:   parser.Manifest,
		Composer: composer.Manifest,
	}
}
