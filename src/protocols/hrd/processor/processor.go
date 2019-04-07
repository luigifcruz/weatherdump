package processor

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"weather-dump/src/ccsds"
	"weather-dump/src/ccsds/frames"
	"weather-dump/src/handlers/interfaces"
	"weather-dump/src/img"
	"weather-dump/src/protocols/helpers"
	"weather-dump/src/protocols/hrd"
	"weather-dump/src/protocols/hrd/processor/composer"
	"weather-dump/src/protocols/hrd/processor/parser"

	"github.com/fatih/color"
	"github.com/gorilla/websocket"
	"github.com/gosuri/uiprogress"
)

var upgrader = websocket.Upgrader{}
var channels = parser.Channels

const frameSize = 892

type Worker struct {
	ccsds    *ccsds.Worker
	scid     uint8
	manifest helpers.ProcessingManifest
}

func NewProcessor(uuid string, manifest *helpers.ProcessingManifest) interfaces.Processor {
	e := Worker{
		ccsds: ccsds.New(),
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
			channels[packet.GetAPID()].Parse(packet)
		}
	}

	e.scid = uint8(helpers.MaxIntSlice(scidStat[:]))
	fmt.Printf("[PRC] Decoded %d packets from VCID 16.\n", len(e.ccsds.GetSpacePackets()))
}

func (e *Worker) Export(outputPath string, wf img.Pipeline) {
	fmt.Printf("[PRC] Exporting VIIRS science products.\n")
	var currentComposer, currentParser uint16

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
			return fmt.Sprintf("[DEC] Processing completed 	")
		default:
			return fmt.Sprintf("[DEC] Rendering channel %s	", e.manifest.Parser[currentParser].Name)
		}
	})

	bar2.PrependFunc(func(b *uiprogress.Bar) string {
		switch currentComposer {
		case 0:
			return fmt.Sprintf("[DEC] Waiting for render	")
		case 9999:
			return fmt.Sprintf("[DEC] Components completed	")
		default:
			return fmt.Sprintf("[DEC] Rendering %s	", e.manifest.Composer[currentComposer].Name)
		}
	})

	e.manifest.WaitForClient(nil)

	for _, apid := range e.manifest.Parser.Parse() {
		currentParser = apid
		ch := channels[apid]

		var buf []byte
		if ch.Export(&buf, channels, hrd.Spacecrafts[e.scid]) {
			w, h := ch.GetDimensions()
			outputName, _ := filepath.Abs(fmt.Sprintf("%s/%s", outputPath, ch.FileName))

			wf.AddException("Invert", ch.Invert)
			wf.Target(img.NewGray16(&buf, w, h)).Process().Export(outputName, 100)
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
		c := composer.Composers[uint16(code)]
		outputName := c.Register(wf, hrd.Spacecrafts[e.scid]).Render(channels, outputPath)
		e.manifest.Composer[code].FileName(outputName)
		e.manifest.Composer[code].Completed()
		e.manifest.Update()
		bar2.Incr()
	}

	currentComposer = 9999

	if !e.manifest.IsRegistred() {
		progress.Stop()
	}

	color.Green("[PRC] Done! All products and components were saved.")
}

func (e Worker) GetProductsManifest() helpers.ProcessingManifest {
	return helpers.ProcessingManifest{
		Parser:   parser.Manifest,
		Composer: composer.Manifest,
	}
}
