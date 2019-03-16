package terminal

import (
	"fmt"
	"os"
	"strings"
	"weather-dump/src/handlers"
	"weather-dump/src/tools/img"
)

func HandleInput(datalink, inputFile, outputPath string, inputDecoded bool, wf img.Pipeline) {
	fmt.Printf("[CLI] Activating %s workflow...\n", strings.ToUpper(datalink))

	heartbeat := true
	workingPath, fileName := handlers.GenerateDirectories(inputFile, outputPath)

	if !inputDecoded {
		decodedFile := fmt.Sprintf("%s/decoded_%s.bin", workingPath, strings.ToLower(fileName))
		handlers.AvailableDecoders[datalink]("").Work(inputFile, decodedFile, &heartbeat)
		inputFile = decodedFile
	}

	if _, err := os.Stat(inputFile); os.IsNotExist(err) || inputFile == "" {
		fmt.Println("[CLI] Input file not found. Exiting...\nError:", err)
		os.Exit(0)
	}

	processor := handlers.AvailableProcessors[datalink]("")
	processor.Work(inputFile)

	manifest := processor.GetProductsManifest()
	processor.Export(workingPath, wf, manifest)
}
