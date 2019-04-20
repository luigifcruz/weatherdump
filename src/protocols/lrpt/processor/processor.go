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
	"weatherdump/src/protocols/lrpt"
	"weatherdump/src/protocols/lrpt/processor/composer"
	"weatherdump/src/protocols/lrpt/processor/parser"

	"github.com/fatih/color"
	"github.com/gorilla/websocket"
	"github.com/gosuri/uiprogress"
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

	e.manifest.Register("lrpt", uuid)

	return &e
}

func (e *Worker) Work(inputFile string) {
	color.Yellow("[PRC] WARNING! This processor is currently in ALPHA development state.")
	scidStat := [256]int{}

	file, _ := ioutil.ReadFile(inputFile)
	for i := len(file); i > 0; i -= frameSize {
		f := frames.NewTransferFrame(file[(len(file) - i):])
		p := frames.NewMultiplexingFrame(ccsds.Version["LRPT"], f.GetMPDU())

		if !f.IsReplay() && p.IsValid() {
			scidStat[f.GetSCID()]++
			switch f.GetVCID() {
			case 5:
				e.ccsds.ParseMPDU(*p) // VCID 5 Parser
			}
		}
	}

	for _, packet := range e.ccsds.GetSpacePackets() {
		if packet.GetAPID() >= 64 && packet.GetAPID() <= 69 {
			e.channels[packet.GetAPID()].Parse(packet)
		}
	}

	e.scid = uint8(helpers.MaxIntSlice(scidStat[:]))
	fmt.Printf("[PRC] Decoded %d packets from VCID 16.\n", len(e.ccsds.GetSpacePackets()))
}

func (e *Worker) Export(outputPath string, wf img.Pipeline) {
	fmt.Printf("[PRC] Exporting BISMW science products.\n")
	var currentParser, currentComposer uint16

	progress := uiprogress.New()

	if !e.manifest.IsRegistred() {
		progress.Start()
	}

	bar1 := progress.AddBar(e.manifest.ParserCount()).AppendCompleted()
	bar2 := progress.AddBar(e.manifest.ComposerCount()).AppendCompleted()

	bar1.PrependFunc(func(b *uiprogress.Bar) string {
		switch currentParser {
		case 0:
			return fmt.Sprintf("[DEC] Starting render		")
		case 9999:
			return fmt.Sprintf("[DEC] Processing completed ")
		default:
			return fmt.Sprintf("[DEC] Rendering %s	", e.manifest.Parser[currentParser].Name)
		}
	})

	bar2.PrependFunc(func(b *uiprogress.Bar) string {
		switch currentComposer {
		case 0:
			return fmt.Sprintf("[DEC] Waiting for render")
		case 9999:
			return fmt.Sprintf("[DEC] Components completed ")
		default:
			return fmt.Sprintf("[DEC] Rendering %s	", e.manifest.Composer[currentComposer].Name)
		}
	})

	e.manifest.WaitForClient(nil)

	for _, apid := range e.manifest.Parser.Parse() {
		currentParser = apid
		ch := e.channels[apid]

		var buf []byte
		if ch.Export(&buf, lrpt.Spacecrafts[e.scid]) {
			w, h := ch.GetDimensions()
			outputName, _ := filepath.Abs(fmt.Sprintf("%s/%s", outputPath, ch.FileName))

			wf.AddException("Invert", ch.Invert)
			wf.Target(img.NewGray(&buf, w, h)).Process().Export(outputName, 100)
			wf.ResetExceptions()

			e.manifest.Parser[apid].FileName(outputName)
		}

		e.manifest.Parser[apid].Completed()
		e.manifest.Update()
		bar1.Incr()
	}

	currentParser = 9999

	for _, code := range e.manifest.Composer.Parse() {
		currentComposer = code
		c := composer.Composers[code]
		outputName := c.Register(wf, lrpt.Spacecrafts[e.scid]).Render(e.channels, outputPath)
		e.manifest.Composer[code].FileName(outputName)
		e.manifest.Composer[code].Completed()
		e.manifest.Update()
		bar2.Incr()
	}

	currentComposer = 9999

	if !e.manifest.IsRegistred() {
		progress.Stop()
	}

	e.channels = make(parser.List)
	e.ccsds = nil
	color.Green("[PRC] Done! All products and components were saved.")
}

func (e Worker) GetProductsManifest() helpers.ProcessingManifest {
	return helpers.ProcessingManifest{
		Parser:   parser.Manifest,
		Composer: composer.Manifest,
	}
}
