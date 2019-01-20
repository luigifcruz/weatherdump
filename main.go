package main

import (
	"fmt"
	"log"
	"os"
	"weather-dump/src/NPOESS/Handler"

	"github.com/urfave/cli"
	"gopkg.in/gographics/imagick.v2/imagick"
)

const frameSize = 892

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
	app.Usage = "OSP's universal decoder for sun-synchronous satellites."
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
				Handler.CommandLine(c.Args().First(), inputFormat, outputFolder)
				return nil
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
