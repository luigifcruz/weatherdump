package main

import (
	"fmt"
	"log"
	"os"

	"weather-dump/src/handlers/remote"
	"weather-dump/src/handlers/terminal"

	"github.com/urfave/cli"
)

const welcome = `
======================= Open Satellite Project =======================
__          __        _   _               _____                        
\ \        / /       | | | |             |  __ \                       
 \ \  /\  / /__  __ _| |_| |__   ___ _ __| |  | |_   _ _ __ ___  _ __  
  \ \/  \/ / _ \/ _' | __| '_ \ / _ \ '__| |  | | | | | '_ ' _ \| '_ \ 
   \  /\  /  __/ (_| | |_| | | |  __/ |  | |__| | |_| | | | | | | |_) |
    \/  \/ \___|\__,_|\__|_| |_|\___|_|  |_____/ \__,_|_| |_| |_| .__/ 
                                                                | |    
								|_|
========================= CLI Version Beta 1 =========================    
`

func main() {
	fmt.Println(welcome)

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

				terminal.HandleInput(c.Args().First(), inputFormat, outputFolder, "hrd")
				return nil
			},
		}, {
			Name:     "lrpt",
			Usage:    "decoder for VHF band Low Rate Picture Transfer (LRPT) signal (MeteorM-N2)",
			Category: "DATALINK",
			Action: func(c *cli.Context) error {
				if len(c.Args().First()) == 0 {
					fmt.Println("[ERROR] Missing file_path.")
					os.Exit(0)
				}

				terminal.HandleInput(c.Args().First(), inputFormat, outputFolder, "lrpt")
				return nil
			},
		}, {
			Name:     "remote",
			Usage:    "listen to network commands",
			Category: "DATALINK",
			Action: func(c *cli.Context) error {
				remote.New().Listen()
				return nil
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
