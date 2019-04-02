package terminal

import (
	"fmt"
	"os"
	"strings"
	"weather-dump/src/handlers"
	"weather-dump/src/img"

	"github.com/fatih/color"
)

// HandleInput with user defined functions gathered by the CLI tool.
func HandleInput(datalink, inputFile, outputPath, decoderType string, wf img.Pipeline) {
	fmt.Printf("[CLI] Activating %s workflow.\n", strings.ToUpper(datalink))

	heartbeat := true
	workingPath, fileName := handlers.GenerateDirectories(inputFile, outputPath)

	if decoderType != "none" {
		if handlers.AvailableDecoders[datalink][decoderType] == nil {
			color.Yellow("[CLI] Invalid decoder input. Try 'weatherdump %s -h' for more information.", datalink)
			return
		}

		decodedFile := fmt.Sprintf("%s/decoded_%s.bin", workingPath, strings.ToLower(fileName))
		handlers.AvailableDecoders[datalink][decoderType]("").Work(inputFile, decodedFile, &heartbeat)
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
