package processor

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"path/filepath"
	"weather-dump/src/ccsds"
	"weather-dump/src/ccsds/frames"
	"weather-dump/src/handlers/interfaces"
	"weather-dump/src/protocols/lrpt"
	"weather-dump/src/protocols/lrpt/processor/bismw"
	"weather-dump/src/tools/img"

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
	e := Worker{
		ccsds: ccsds.New(),
		bismw: bismw.New(),
	}

	http.HandleFunc(fmt.Sprintf("/meteor/%s/statistics", uuid), e.statistics)
	return &e
}

func (e *Worker) Work(inputFile string) {
	fmt.Println("[PRC] WARNING! This processor is currently in ALPHA development state.")

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

	fmt.Printf("[PRC] Found %d packets from VCID 16.\n", len(e.ccsds.GetSpacePackets()))
	for _, packet := range e.ccsds.GetSpacePackets() {
		if packet.GetAPID() >= 64 && packet.GetAPID() <= 69 {
			e.bismw.Parse(packet)
		}
	}

	fmt.Println("[PRC] Finished decoding all packets...")
}

func (e *Worker) Export(outputPath string, wf img.Pipeline) {
	fmt.Printf("[PRC] Exporting BISMW science products to %s...\n", outputPath)

	for _, apid := range bismw.ChannelsIndex {
		channel := e.bismw.Channel(apid)

		if channel == nil {
			continue
		}

		channel.Fix(lrpt.Spacecrafts[e.scid])
		buf := channel.Compose()
		w, h := channel.GetDimensions()

		outputName, _ := filepath.Abs(fmt.Sprintf("%s/%s", outputPath, channel.GetFileName()))
		wf.Marshal("Invert", bismw.ChannelsParameters[apid].Inversion)
		wf.Target(img.NewGray(buf, w, h)).Process().Export(outputName, 100)
		wf.Unmarshal()
	}

	fmt.Println("[PRC] Done! Products saved.")
}

func (e *Worker) statistics(w http.ResponseWriter, r *http.Request) {
	upgrader.CheckOrigin = func(r *http.Request) bool { return true }
	e.statsSock, _ = upgrader.Upgrade(w, r, nil)
}
