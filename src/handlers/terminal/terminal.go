package terminal

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"weather-dump/src/handlers"
)

func HandleInput(inputFile string, inputFormat string, outputPath string, datalink string) {
	fmt.Printf("[CLI] Activating %s workflow...\n", strings.ToUpper(datalink))

	heartbeat := true
	workingPath := filepath.Dir(inputFile)
	inputFileName := filepath.Base(inputFile)
	inputFileName = strings.TrimSuffix(inputFileName, filepath.Ext(inputFile))

	if outputPath != "" {
		workingPath = outputPath
	}

	outputPath = fmt.Sprintf("%s/OUTPUT_%s", workingPath, strings.ToUpper(inputFileName))
	if _, err := os.Stat(outputPath); os.IsNotExist(err) {
		os.Mkdir(outputPath, os.ModePerm)
	}

	if inputFormat == "grcout" {
		decodedFile := fmt.Sprintf("%s/decoded_%s.bin", outputPath, strings.ToLower(inputFileName))
		handlers.AvailableDecoders[datalink]("").Work(inputFile, decodedFile, &heartbeat)
		inputFile = decodedFile
	}

	if _, err := os.Stat(inputFile); os.IsNotExist(err) || inputFile == "" {
		fmt.Println("[CLI] Input file not found. Exiting...\nError:", err)
		os.Exit(0)
	}

	processor := handlers.AvailableProcessors[datalink]("")
	processor.Work(inputFile)
	processor.ExportAll(outputPath)
}
