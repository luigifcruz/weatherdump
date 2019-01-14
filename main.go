package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"weather-dump/src/CCSDS"
	"weather-dump/src/CCSDS/Frames"
	"weather-dump/src/NPOESS/Decoder"
	"weather-dump/src/NPOESS/VIIRS"

	"github.com/urfave/cli"
	"gopkg.in/gographics/imagick.v2/imagick"
)

const frameSize = 892

func runHRDDecoder(fileName string, outputPath string) {
	fmt.Println("[HRD] Decoding started...")

	dec := Decoder.NewDecoder()
	dec.DecodeFile()

	if _, err := os.Stat(outputPath); os.IsNotExist(err) {
		os.Mkdir(outputPath, os.ModePerm)
	}

	file, _ := ioutil.ReadFile(fileName)

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
	fmt.Printf("[HRD] Exporting VIIRS science products to %s...\n", outputPath)

	viirs.SetOutputFolder(outputPath)
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

	var inputFormat string
	var outputPath string

	app := cli.NewApp()

	app.Name = "weatherdump"
	app.UsageText = "weatherdump [OPTIONS] [DATALINK] [FILE_PATH]"
	app.Author = "Luigi Cruz (@luigifcruz) for Open Satellite Project"
	app.Usage = "OSP's universal decoder for polar orbiting satellites."
	app.Version = "0.1.0"

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:        "format",
			Value:       "decoded",
			Usage:       "input format [decoded or cadu]",
			Destination: &inputFormat,
		},
		cli.StringFlag{
			Name:        "output",
			Value:       "./output",
			Usage:       "path where the data decoded will be saved",
			Destination: &outputPath,
		},
	}

	app.Commands = []cli.Command{
		{
			Name:     "hrd",
			Usage:    "decoder for X-Band High Rate Data (HRD) signal (Suomi & NOAA-20)",
			Category: "DATALINK",
			Action: func(c *cli.Context) error {
				if inputFormat == "cadu" {
					fmt.Println("[ERROR] The CADU type input isn't supported yet.")
					os.Exit(0)
				}

				if len(c.Args().First()) == 0 {
					fmt.Println("[ERROR] Missing file_path.")
					os.Exit(0)
				}

				settingsPrint(inputFormat, outputPath, "HRD")
				runHRDDecoder(c.Args().First(), outputPath)
				return nil
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
