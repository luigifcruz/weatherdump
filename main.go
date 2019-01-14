package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"time"
	"weather-dump/src/CCSDS"
	"weather-dump/src/CCSDS/Frames"
	"weather-dump/src/NPOESS/Decoder"
	"weather-dump/src/NPOESS/VIIRS"

	"github.com/urfave/cli"
	"gopkg.in/gographics/imagick.v2/imagick"
)

const frameSize = 892

func runHRDDecoder(inputPath string, inputFormat string, outputFolder string) {
	fmt.Println("[HRD] Decoding started...")

	if _, err := os.Stat(outputFolder); os.IsNotExist(err) {
		os.Mkdir(outputFolder, os.ModePerm)
	}

	outputFolder = fmt.Sprintf("%s/NPOESS-HRD-%s", outputFolder, time.Now().Format(time.RFC3339))
	os.Mkdir(outputFolder, os.ModePerm)

	if inputFormat == "grcout" {
		dec := Decoder.NewDecoder()
		outputFile := fmt.Sprintf("%s/decoded.bin", outputFolder)
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
	viirs := VIIRS.Data{}
	scid := uint8(0)

	fmt.Println("[HRD] Decoding CCSDS packets...")

	for i := len(file); i > 0; i -= frameSize {
		s := Frames.NewTransferFrame(file[(len(file) - i):])
		scid = s.GetSCID()

		if s.IsReplay() {
			p := Frames.NewMultiplexingFrame(s.GetMPDU())

			switch s.GetVCID() {
			case 16:
				ch16.ParseMPDU(*p) // VCID 16 Parser (VIIRS)
			}
		}
	}

	fmt.Println("[HRD] Decoding VCID 16 packets...")

	skippedPackets := 0
	for _, packet := range ch16.GetSpacePackets() {
		if !packet.IsValid() {
			skippedPackets += 1
			continue
		}

		if packet.GetAPID() >= 800 && packet.GetAPID() <= 823 {
			viirs.Parse(packet)
		}
	}

	fmt.Printf("[HRD] Found %d/%d invalid packets in VCID 16...\n", skippedPackets, len(ch16.GetSpacePackets()))
	fmt.Printf("[HRD] Exporting VIIRS science products to %s...\n", outputFolder)

	viirs.SetOutputFolder(outputFolder)
	viirs.SaveAllChannels(scid)
	viirs.ExportTrueColor(scid)

	fmt.Println("[HRD] Done! Products saved.")
}

func settingsPrint(inputFormat string, outputPath string, datalinkName string) {
	fmt.Println("============== WeatherDump ==============")
	fmt.Println("============= CONFIGURATION =============")
	fmt.Println("Datalink:", datalinkName)
	fmt.Println("Input Format:", inputFormat)
	fmt.Println("Output Folder:", outputPath)
	fmt.Println("=========================================")
}

func main() {
	imagick.Initialize()
	defer imagick.Terminate()

	var outputFolder string
	var inputFormat string

	app := cli.NewApp()

	app.Name = "weatherdump"
	app.UsageText = "weatherdump [OPTIONS] [DATALINK] [FILE_PATH]"
	app.Author = "Luigi Cruz (@luigifcruz) for Open Satellite Project"
	app.Usage = "OSP's universal decoder for polar orbiting satellites."
	app.Version = "1.0.0"

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:        "format",
			Value:       "grcout",
			Usage:       "input format [grcout or decoded]",
			Destination: &inputFormat,
		},
		cli.StringFlag{
			Name:        "output",
			Value:       "./output",
			Usage:       "folder where the products will be saved",
			Destination: &outputFolder,
		},
	}

	app.Commands = []cli.Command{
		{
			Name:     "hrd",
			Usage:    "decoder for X-Band High Rate Data (HRD) signal (Suomi & NOAA-20)",
			Category: "DATALINK",
			Action: func(c *cli.Context) error {
				if len(c.Args().First()) == 0 {
					fmt.Println("[ERROR] Missing file_path.")
					os.Exit(0)
				}

				settingsPrint(outputFolder, outputFolder, "HRD")
				runHRDDecoder(c.Args().First(), inputFormat, outputFolder)
				return nil
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
