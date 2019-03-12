package main

import (
	"fmt"
	"time"
	remoteHandler "weather-dump/src/handlers/remote"
	terminalHandler "weather-dump/src/handlers/terminal"
	"weather-dump/src/tools/img"

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

	inputFormat = kingpin.Flag("decoded", "input file format").Short('d').Default("false").Bool()

	exportPNG  = kingpin.Flag("png", "export pictures as PNG").Default("false").Bool()
	exportJPEG = kingpin.Flag("jpeg", "export pictures as JPEG").Default("false").Bool()
	equalize   = kingpin.Flag("equalize", "histogram equalize output pictures (--no-equalize turn-off)").Short('e').Default("true").Bool()
	invert     = kingpin.Flag("invert", "invert infrared output pictures (--no-invert turn-off)").Short('i').Default("true").Bool()
	flop       = kingpin.Flag("flop", "flop output pictures (--no-flop turn-off)").Short('f').Default("true").Bool()

	hrd          = kingpin.Command("hrd", "Activate workflow for the HRD protocol (NPOESS & NPP).")
	hrdInputFile = hrd.Arg("file", "input file path").Required().ExistingFile()

	lrpt          = kingpin.Command("lrpt", "Activate workflow for the LRPT protocol (Meteor-MN2).")
	lrptInputFile = lrpt.Arg("file", "input file path").Required().ExistingFile()

	remote     = kingpin.Command("remote", "Activate the remote API for the GUI.")
	remotePort = remote.Arg("port", "server listen port").Default("3000").String()
)

func main() {
	fmt.Println(startMessage)

	kingpin.CommandLine.HelpFlag.Short('h')
	kingpin.Version("Beta 1")

	datalink := kingpin.Parse()

	if datalink == "remote" {
		remoteHandler.New().Listen(*remotePort)
	}

	wf := img.NewPipeline()

	wf.AddPipe("Equalize", *equalize)
	wf.AddPipe("Flop", *flop)
	wf.AddPipe("Invert", *invert)
	wf.AddPipe("ExportPNG", *exportPNG)
	wf.AddPipe("ExportJPEG", *exportJPEG)

	start := time.Now()
	terminalHandler.HandleInput(datalink, *lrptInputFile+*hrdInputFile, *output, *inputFormat, wf)
	fmt.Printf("Tasks finished in %s\n", time.Since(start))
}
