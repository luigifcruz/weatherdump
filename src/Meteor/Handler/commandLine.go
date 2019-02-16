package Handler

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"
	"strings"
	"time"
	"weather-dump/src/CCSDS"
	"weather-dump/src/CCSDS/Frames"
	"weather-dump/src/Meteor/BISMW"
	"weather-dump/src/Meteor/Decoder"
)

const frameSize = 892

func CommandLine(inputPath string, inputFormat string, outputFolder string) {
	fmt.Println("[LRPT] Decoding started...")

	go http.ListenAndServe(":3000", nil)

	if _, err := os.Stat(outputFolder); os.IsNotExist(err) {
		os.Mkdir(outputFolder, os.ModePerm)
	}

	fileName := time.Now().Format(time.RFC3339)
	r, err := regexp.Compile("^.*\\/(.*)\\.\\w+$")

	if err == nil {
		fileName = r.FindStringSubmatch(inputPath)[1]
	}

	outputFolder = fmt.Sprintf("%s/METEOR-LRPT-%s", outputFolder, strings.ToUpper(fileName))
	os.Mkdir(outputFolder, os.ModePerm)

	if inputFormat == "grcout" {
		dec := Decoder.NewDecoder()
		outputFile := fmt.Sprintf("%s/decoded-%s.bin", outputFolder, strings.ToLower(fileName))
		dec.DecodeFile(inputPath, outputFile)
		inputPath = outputFile
	}

	file, err := ioutil.ReadFile(inputPath)
	if err != nil {
		fmt.Println("[LRPT] Input file not found. Exiting...")
		fmt.Println(err)
		os.Exit(0)
	}

	ch05 := CCSDS.CCSDS{}
	scid := uint8(0)

	fmt.Println("[LRPT] Decoding CCSDS packets...")

	for i := len(file); i > 0; i -= frameSize {
		s := Frames.NewTransferFrame(file[(len(file) - i):])
		scid = s.GetSCID()

		if !s.IsReplay() {
			p := Frames.NewMultiplexingFrame(CCSDS.Version["LRPT"], s.GetMPDU())

			switch s.GetVCID() {
			case 5:
				ch05.ParseMPDU(*p) // VCID 5 Parser
			}
		}
	}

	bismw := BISMW.NewData(scid)
	for _, packet := range ch05.GetSpacePackets() {
		if packet.GetAPID() >= 64 && packet.GetAPID() <= 69 {
			bismw.Parse(packet)
		}
	}

	bismw.Process()
	bismw.SaveAllChannels(outputFolder)
}
