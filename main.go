package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"weather-dump/src/CCSDS"
	"weather-dump/src/CCSDS/Frames"
	VIIRS "weather-dump/src/VIIRS/ScienceData"

	"github.com/urfave/cli"
	"gopkg.in/gographics/imagick.v2/imagick"
)

const frameSize = 892

func runHRDDecoder(fileName string, outputPath string) {
	fmt.Println("Decoding started...")

	if _, err := os.Stat(outputPath); os.IsNotExist(err) {
		os.Mkdir(outputPath, os.ModePerm)
	}

	file, _ := ioutil.ReadFile(fileName)

	d := CCSDS.CCSDS{}
	viirs := VIIRS.ScienceData{}
	scid := uint8(0)

	bytesCount := 0
	bytesNumber := len(file)

	fmt.Println("Decoding CCSDS packets...")

	for bytesCount < bytesNumber {
		s := Frames.TransferFrame{}
		s.FromBinary(file[bytesCount:])
		scid = s.GetSCID()

		if s.GetVCID() == 16 || s.GetVCID() == 0 {
			p := Frames.MultiplexingFrame{}
			p.FromBinary(s.GetMPDU())

			CCSDS.ParseMPDU(&d, p)
		}

		bytesCount += frameSize
	}

	fmt.Println("Decoding science packets...")
	skippedPackets := 0

	for _, packet := range d.GetSpacePackets() {
		if !packet.IsValid() {
			skippedPackets += 1
			continue
		}

		if packet.GetAPID() == 8 {
			packet.Print()
			viirs.ParseTime(packet)
		}

		if packet.GetAPID() >= 800 && packet.GetAPID() <= 823 {
			//packet.Print()
			//viirs.Parse(packet)
		}
	}

	fmt.Printf("Found %d invalid packets...\n", skippedPackets)

	viirs.SetOutputFolder(outputPath)
	viirs.SaveAllChannels(scid)
	viirs.ExportTrueColor(scid)

	fmt.Println("Done! Products saved.")
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
					fmt.Println("[ERROR] The CADU type input not supported yet.")
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
