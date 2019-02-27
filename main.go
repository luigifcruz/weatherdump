package main

import (
	"fmt"
	remoteHandler "weather-dump/src/handlers/remote"
	terminalHandler "weather-dump/src/handlers/terminal"

	kingpin "gopkg.in/alecthomas/kingpin.v2"
)

const startMessage = `
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

var (
	output = kingpin.Flag("output", "Custom output folder. Default is the current input file folder.").Default("").String()

	hrd       = kingpin.Command("hrd", "Activate workflow for the HRD protocol (NPOESS & NPP).")
	hrdFile   = hrd.Arg("file", "input file path").Required().ExistingFile()
	hrdFormat = hrd.Flag("decoded", "input file format").Short('d').Default("false").Bool()

	lrpt       = kingpin.Command("lrpt", "Activate workflow for the LRPT protocol (Meteor-MN2).")
	lrptFile   = lrpt.Arg("file", "input file path").Required().ExistingFile()
	lrptFormat = lrpt.Flag("decoded", "input file format").Short('d').Default("false").Bool()

	remote = kingpin.Command("remote", "Activate the remote API for the GUI.")
)

func main() {
	fmt.Println(startMessage)

	kingpin.CommandLine.HelpFlag.Short('h')
	kingpin.Version("Beta 1")

	switch kingpin.Parse() {
	case "hrd":
		terminalHandler.HandleInput(*hrdFile, *hrdFormat, *output, "hrd")
	case "lrpt":
		terminalHandler.HandleInput(*lrptFile, *lrptFormat, *output, "lrpt")
	case "remote":
		remoteHandler.New().Listen()
	}
}
