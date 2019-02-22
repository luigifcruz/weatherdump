package processor

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"weather-dump/src/ccsds"
	"weather-dump/src/ccsds/frames"
	"weather-dump/src/interfaces"
	"weather-dump/src/meteor/processor/bismw"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{}

const frameSize = 892

type Worker struct {
	ccsds     *ccsds.Worker
	bismw     *bismw.Worker
	scid      uint8
	statsSock *websocket.Conn
}

func NewProcessor(uuid string) interfaces.Processor {
	e := Worker{}
	e.ccsds = ccsds.New()
	e.bismw = bismw.New()

	http.HandleFunc(fmt.Sprintf("/meteor/%s/statistics", uuid), e.statistics)
	return &e
}

func (e *Worker) Work(inputFile string) {
	fmt.Println("[PRC] Decoding CCSDS packets...")

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

	fmt.Printf("[PRC] Decoding %d VCID 5 packets...\n", len(e.ccsds.GetSpacePackets()))
	for _, packet := range e.ccsds.GetSpacePackets() {
		if packet.GetAPID() >= 64 && packet.GetAPID() <= 69 {
			e.bismw.Parse(packet)
		}
	}

	fmt.Println("[PRC] Processing channels...")
	e.bismw.Process(e.scid)

	fmt.Println("[PRC] Finished decoding packets...")
}

func (e *Worker) ExportAll(outputPath string) {
	fmt.Printf("[PRC] Exporting BISMW science products to %s...\n", outputPath)
	e.bismw.SaveAllChannels(outputPath)
	fmt.Println("[PRC] Done! Products saved.")
}

func (e *Worker) statistics(w http.ResponseWriter, r *http.Request) {
	upgrader.CheckOrigin = func(r *http.Request) bool { return true }
	e.statsSock, _ = upgrader.Upgrade(w, r, nil)
}
