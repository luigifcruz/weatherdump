package processor

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"path/filepath"
	"weather-dump/src/assets"
	"weather-dump/src/ccsds"
	"weather-dump/src/ccsds/frames"
	"weather-dump/src/handlers/interfaces"
	"weather-dump/src/protocols/hrd"
	"weather-dump/src/protocols/hrd/processor/composer"
	"weather-dump/src/protocols/hrd/processor/parser"
	"weather-dump/src/tools/img"

	"github.com/fatih/color"
	"github.com/gorilla/websocket"
	"github.com/gosuri/uiprogress"
)

var upgrader = websocket.Upgrader{}
var channels = parser.Channels

const frameSize = 892

type Worker struct {
	ccsds     *ccsds.Worker
	scid      uint8
	statsSock *websocket.Conn
}

func NewProcessor(uuid string) interfaces.Processor {
	e := Worker{}
	e.ccsds = ccsds.New()

	if uuid != "" {
		http.HandleFunc(fmt.Sprintf("/hrd/%s/statistics", uuid), e.statistics)
	}

	return &e
}

func (e *Worker) Work(inputFile string) {
	color.Yellow("[PRC] WARNING! This processor is currently in BETA development state.")

	file, _ := ioutil.ReadFile(inputFile)
	for i := len(file); i > 0; i -= frameSize {
		f := frames.NewTransferFrame(file[(len(file) - i):])
		e.scid = f.GetSCID()

		if f.IsReplay() {
			p := frames.NewMultiplexingFrame(ccsds.Version["HRD"], f.GetMPDU())
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

	fmt.Printf("[PRC] Decoded %d packets from VCID 16.\n", len(e.ccsds.GetSpacePackets()))
}

func (e *Worker) Export(outputPath string, wf img.Pipeline, manifest assets.ProcessingManifest) {
	fmt.Printf("[PRC] Exporting VIIRS science products.\n")
	var currentComposer, currentParser uint16

	progress := uiprogress.New()
	progress.Start()

	bar1 := progress.AddBar(manifest.ParserCount()).AppendCompleted()
	bar2 := progress.AddBar(manifest.ComposerCount()).AppendCompleted()

	bar1.PrependFunc(func(b *uiprogress.Bar) string {
		switch currentParser {
		case 0:
			return fmt.Sprintf("[DEC] Starting render		")
		case 9999:
			return fmt.Sprintf("[DEC] Processing completed 	")
		default:
			return fmt.Sprintf("[DEC] Rendering channel %s	", manifest.Parser[currentParser].Name)
		}
	})

	bar2.PrependFunc(func(b *uiprogress.Bar) string {
		switch currentComposer {
		case 0:
			return fmt.Sprintf("[DEC] Waiting for render	")
		case 9999:
			return fmt.Sprintf("[DEC] Components completed	")
		default:
			return fmt.Sprintf("[DEC] Rendering %s	", manifest.Composer[currentComposer].Name)
		}
	})

	for _, apid := range manifest.Parser.Ordered() {
		currentParser = apid
		ch := channels[apid]

		var buf []byte
		if ch.Export(&buf, channels, hrd.Spacecrafts[e.scid]) {
			w, h := ch.GetDimensions()
			outputName, _ := filepath.Abs(fmt.Sprintf("%s/%s", outputPath, ch.FileName))

			wf.AddException("Invert", ch.Invert)
			wf.Target(img.NewGray16(&buf, w, h)).Process().Export(outputName, 100)
			wf.ResetExceptions()
		}

		manifest.Parser.Completed(apid, e.statsSock)
		bar1.Incr()
	}

	currentParser = 9999

	for code := range manifest.Composer {
		currentComposer = code
		c := composer.Composers[uint16(code)]
		c.Register(wf, hrd.Spacecrafts[e.scid]).Render(channels, outputPath)
		manifest.Composer.Completed(code, e.statsSock)
		bar2.Incr()
	}

	currentComposer = 9999
	progress.Stop()
	color.Green("[PRC] Done! All products and components were saved.")
}

func (e *Worker) statistics(w http.ResponseWriter, r *http.Request) {
	upgrader.CheckOrigin = func(r *http.Request) bool { return true }
	e.statsSock, _ = upgrader.Upgrade(w, r, nil)
}

func (e Worker) GetProductsManifest() assets.ProcessingManifest {
	return assets.ProcessingManifest{
		Parser:   parser.Manifest,
		Composer: composer.Manifest,
	}
}
