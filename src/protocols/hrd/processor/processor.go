package processor

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"path/filepath"
	"weather-dump/src/ccsds"
	"weather-dump/src/ccsds/frames"
	"weather-dump/src/handlers/interfaces"
	"weather-dump/src/protocols/hrd"
	"weather-dump/src/protocols/hrd/processor/viirs"
	"weather-dump/src/tools/img"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{}

const frameSize = 892

type Worker struct {
	ccsds     *ccsds.Worker
	viirs     *viirs.Worker
	scid      uint8
	statsSock *websocket.Conn
}

func NewProcessor(uuid string) interfaces.Processor {
	e := Worker{}
	e.ccsds = ccsds.New()
	e.viirs = viirs.New()

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
			e.viirs.Parse(packet)
		}
	}

	fmt.Println("[PRC] Finished decoding all packets...")
}

func (e *Worker) Export(outputPath string, wf img.Pipeline) {
	fmt.Printf("[PRC] Exporting VIIRS science products to %s...\n", outputPath)

	for _, apid := range viirs.ChannelsIndex {
		channel := e.viirs.Channel(apid)

		if channel == nil {
			continue
		}

		channel.Fix(hrd.Spacecrafts[e.scid])

		w, h := channel.GetDimensions()
		buf := make([]byte, w*h*2)

		reconChannel := channel.GetReconstructionBand()
		if reconChannel == 000 {
			channel.ComposeUncoded(&buf)
		} else {
			if e.viirs.Channel(reconChannel) == nil {
				continue
			}
			channel.ComposeCoded(&buf, e.viirs.Channel(reconChannel))
		}

		outputName, _ := filepath.Abs(fmt.Sprintf("%s/%s", outputPath, channel.GetFileName()))
		wf.Target(img.NewGray16(&buf, w, h)).Process().Export(outputName, 100)
	}

	fmt.Println("[PRC] Done! Products saved.")
}

func (e *Worker) statistics(w http.ResponseWriter, r *http.Request) {
	upgrader.CheckOrigin = func(r *http.Request) bool { return true }
	e.statsSock, _ = upgrader.Upgrade(w, r, nil)
}
