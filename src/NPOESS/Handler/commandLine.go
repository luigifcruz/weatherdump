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
	"weather-dump/src/NPOESS/Decoder"
	"weather-dump/src/NPOESS/VIIRS"
)

const frameSize = 892

func CommandLine(inputPath string, inputFormat string, outputFolder string) {
	fmt.Println("[HRD] Decoding started...")

	go http.ListenAndServe(":3000", nil)

	if _, err := os.Stat(outputFolder); os.IsNotExist(err) {
		os.Mkdir(outputFolder, os.ModePerm)
	}

	fileName := time.Now().Format(time.RFC3339)
	r, err := regexp.Compile("^.*\\/(.*)\\.\\w+$")

	if err == nil {
		fileName = r.FindStringSubmatch(inputPath)[1]
	}

	outputFolder = fmt.Sprintf("%s/NPOESS-HRD-%s", outputFolder, strings.ToUpper(fileName))
	os.Mkdir(outputFolder, os.ModePerm)

	if inputFormat == "grcout" {
		dec := Decoder.NewDecoder()
		outputFile := fmt.Sprintf("%s/decoded-%s.bin", outputFolder, strings.ToLower(fileName))
		dec.DecodeFile(inputPath, outputFile)
		inputPath = outputFile
	}

	file, err := ioutil.ReadFile(inputPath)
	if err != nil {
		fmt.Println("[HRD] Input file not found. Exiting...")
		fmt.Println(err)
		os.Exit(0)
	}

	ch16 := CCSDS.CCSDS{}
	scid := uint8(0)

	fmt.Println("[HRD] Decoding CCSDS packets...")

	for i := len(file); i > 0; i -= frameSize {
		s := Frames.NewTransferFrame(file[(len(file) - i):])
		scid = s.GetSCID()

		if s.IsReplay() {
			p := Frames.NewMultiplexingFrame(CCSDS.Version["HRD"], s.GetMPDU())

			switch s.GetVCID() {
			case 16:
				ch16.ParseMPDU(*p) // VCID 16 Parser (VIIRS)
			}
		}
	}

	fmt.Printf("[HRD] Decoding %d VCID 16 packets...\n", len(ch16.GetSpacePackets()))

	viirs := VIIRS.NewData(scid)
	for _, packet := range ch16.GetSpacePackets() {
		if packet.GetAPID() >= 800 && packet.GetAPID() <= 823 {
			viirs.Parse(packet)
		}
	}

	fmt.Printf("[HRD] Exporting VIIRS science products to %s...\n", outputFolder)

	viirs.Process()
	viirs.SaveAllChannels(outputFolder)
	viirs.SaveTrueColorChannel(outputFolder)

	fmt.Println("[HRD] Done! Products saved.")
}
