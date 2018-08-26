package main

import (
	"fmt"
	"os"
	"log"
	"io/ioutil"
	"osp-noaa-dump/VIIRS/ScienceData"
	"osp-noaa-dump/CCSDS"
	"osp-noaa-dump/CCSDS/Frames"
	"github.com/urfave/cli"
)

const frameSize = 892

func runDecoder(fileName string, outputPath string) {
	fmt.Println("Satellite Helper App - NPOESS Edition")

	if _, err := os.Stat(outputPath); os.IsNotExist(err) {
		os.Mkdir(outputPath, os.ModePerm)
	}

	file, _ := ioutil.ReadFile(fileName)

	d := CCSDS.CCSDS{}
	t := VIIRS.ScienceData{}
	scid := uint8(0)

	bytesCount := 0
	bytesNumber := len(file)

	fmt.Println("Decoding CCSDS packets...")

	for bytesCount < bytesNumber {
		s := Frames.TransferFrame{}
		s.FromBinary(file[bytesCount:])
		scid = s.GetSCID()

		if s.GetVCID() == 16 {
			p := Frames.MultiplexingFrame{}
			p.FromBinary(s.GetMPDU())

			CCSDS.ParseMPDU(&d, p)
		}

		bytesCount += frameSize
	}

	fmt.Println("Decoding science packets...")

	for _, packet := range d.GetSpacePackets() {
		if packet.GetAPID() >= 800 && packet.GetAPID() <= 823 {
			t.Parse(packet)
		}
	}

	t.SaveAllChannels(scid)

	fmt.Println("Done! Products saved.")
}

func main() {
	var fileType string
	var outputPath string

	app := cli.NewApp()

	app.Name = "polardump"
	app.UsageText = "polardump [DATALINK] [FILE_PATH] [OPTIONS]..."
	app.Author = "Luigi Cruz (@luigifcruz) for Open Satellite Project"
	app.Usage = "OSP's universal decoder for polar orbiting satellites."
	app.Version = "0.1.0"

	app.Flags = []cli.Flag {
		cli.StringFlag {
		  Name: "type",
		  Value: "decoded",
		  Usage: "file encoding type",
		  Destination: &fileType,
		},
		cli.StringFlag {
			Name: "output",
			Value: "./output",
			Usage: "path where the data decoded will be saved",
			Destination: &outputPath,
		},
	}

	app.Commands = []cli.Command {
		{
			Name: "hrd",
			Usage: "decoder for X-Band High Rate Data (HRD) signal (Suomi & NOAA-20)",
			Category: "DATALINK",
			Action: func(c *cli.Context) error {
				if fileType == "cadu" {
					fmt.Println("[ERROR] The CADU type input not supported yet.")
					os.Exit(0)
				}

				if len(c.Args().First()) == 0 {
					fmt.Println("[ERROR] Missing file_path.")
					os.Exit(0)
				}

				runDecoder(c.Args().First(), outputPath)
				return nil
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
