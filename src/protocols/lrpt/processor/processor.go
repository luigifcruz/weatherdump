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
	"weather-dump/src/protocols/lrpt"
	"weather-dump/src/protocols/lrpt/processor/parser"
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
	e := Worker{
		ccsds: ccsds.New(),
	}

	if uuid != "" {
		http.HandleFunc(fmt.Sprintf("/lrpt/%s/statistics", uuid), e.statistics)
	}

	return &e
}

func (e *Worker) Work(inputFile string) {
	color.Yellow("[PRC] WARNING! This processor is currently in ALPHA development state.")

	file, _ := ioutil.ReadFile(inputFile)
	for i := len(file); i > 0; i -= frameSize {
		f := frames.NewTransferFrame(file[(len(file) - i):])
		e.scid = f.GetSCID()

		if !f.IsReplay() {
			p := frames.NewMultiplexingFrame(ccsds.Version["LRPT"], f.GetMPDU())
			switch f.GetVCID() {
			case 5:
				e.ccsds.ParseMPDU(*p) // VCID 5 Parser
			}
		}
	}

	for _, packet := range e.ccsds.GetSpacePackets() {
		if packet.GetAPID() >= 64 && packet.GetAPID() <= 69 {
			channels[packet.GetAPID()].Parse(packet)
		}
	}

	fmt.Printf("[PRC] Decoded %d packets from VCID 16.\n", len(e.ccsds.GetSpacePackets()))
}

func (e *Worker) Export(outputPath string, wf img.Pipeline, manifest assets.ProcessingManifest) {
	fmt.Printf("[PRC] Exporting BISMW science products.\n")
	var currentParser uint16

	progress := uiprogress.New()
	progress.Start()
	bar := progress.AddBar(manifest.ParserCount()).AppendCompleted()

	bar.PrependFunc(func(b *uiprogress.Bar) string {
		switch currentParser {
		case 0:
			return fmt.Sprintf("[DEC] Starting decoder		")
		case 9999:
			return fmt.Sprintf("[DEC] Processing completed 	")
		default:
			return fmt.Sprintf("[DEC] Rendering channel %s	", manifest.Parser[currentParser].Name)
		}
	})

	for _, apid := range manifest.Parser.Ordered() {
		currentParser = apid
		ch := channels[apid]

		var buf []byte
		if ch.Export(&buf, lrpt.Spacecrafts[e.scid]) {
			w, h := ch.GetDimensions()
			outputName, _ := filepath.Abs(fmt.Sprintf("%s/%s", outputPath, ch.FileName))

			wf.AddException("Invert", ch.Invert)
			wf.Target(img.NewGray(&buf, w, h)).Process().Export(outputName, 100)
			wf.ResetExceptions()
		}

		manifest.Parser.Completed(apid, e.statsSock)
		bar.Incr()
	}

	currentParser = 9999
	progress.Stop()
	color.Green("[PRC] Done! All products and components were saved.")
}

func (e Worker) GetProductsManifest() assets.ProcessingManifest {
	return assets.ProcessingManifest{
		Parser:   parser.Manifest,
		Composer: assets.Manifest{},
	}
}

func (e *Worker) statistics(w http.ResponseWriter, r *http.Request) {
	upgrader.CheckOrigin = func(r *http.Request) bool { return true }
	e.statsSock, _ = upgrader.Upgrade(w, r, nil)
}
