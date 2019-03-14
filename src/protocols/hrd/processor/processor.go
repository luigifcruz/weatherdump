package processor

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"path/filepath"
	"sort"
	"weather-dump/src/ccsds"
	"weather-dump/src/ccsds/frames"
	"weather-dump/src/handlers/interfaces"
	"weather-dump/src/protocols/hrd"
	"weather-dump/src/protocols/hrd/processor/composer"
	"weather-dump/src/protocols/hrd/processor/parser"
	"weather-dump/src/tools/img"

	"github.com/gorilla/websocket"
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

	http.HandleFunc(fmt.Sprintf("/npoess/%s/statistics", uuid), e.statistics)

	return &e
}

func (e *Worker) Work(inputFile string) {
	fmt.Println("[PRC] WARNING! This processor is currently in BETA development state.")

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

	fmt.Printf("[PRC] Found %d packets from VCID 16.\n", len(e.ccsds.GetSpacePackets()))

	for _, packet := range e.ccsds.GetSpacePackets() {
		if packet.GetAPID() >= 800 && packet.GetAPID() <= 823 {
			channels[packet.GetAPID()].Parse(packet)
		}
	}

	fmt.Println("[PRC] Finished decoding all packets...")
}

func (e *Worker) Export(outputPath string, wf img.Pipeline) {
	fmt.Printf("[PRC] Exporting VIIRS science products to %s...\n", outputPath)

	for _, apid := range getKeys(channels) {
		ch := channels[uint16(apid)]

		if !ch.HasData {
			continue
		}

		ch.Fix(hrd.Spacecrafts[e.scid])

		w, h := ch.GetDimensions()
		buf := make([]byte, w*h*2)

		reconChannel := ch.ReconstructionBand
		if reconChannel == 000 {
			ch.ExportUncoded(&buf)
		} else {
			if !channels[reconChannel].HasData {
				continue
			}
			ch.ExportCoded(&buf, channels[reconChannel])
		}

		outputName, _ := filepath.Abs(fmt.Sprintf("%s/%s", outputPath, ch.FileName))
		wf.AddException("Invert", ch.Invert)
		wf.Target(img.NewGray16(&buf, w, h)).Process().Export(outputName, 100)
		wf.ResetExceptions()
	}

	c := composer.Composers["True-Color"]
	c.Register(wf, hrd.Spacecrafts[e.scid]).Render(channels, outputPath)

	fmt.Println("[PRC] Done! Products saved.")
}

func (e *Worker) statistics(w http.ResponseWriter, r *http.Request) {
	upgrader.CheckOrigin = func(r *http.Request) bool { return true }
	e.statsSock, _ = upgrader.Upgrade(w, r, nil)
}

func getKeys(tasks parser.ChannelList) []int {
	keys := make([]int, 0, len(tasks))
	for k := range tasks {
		keys = append(keys, int(k))
	}
	sort.Ints(keys)
	return keys
}
