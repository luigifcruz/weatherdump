package terminal

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	interfaces "weather-dump/src/interfaces"
	meteorDecoder "weather-dump/src/meteor/decoder"
	meteorProcessor "weather-dump/src/meteor/processor"
	npoessDecoder "weather-dump/src/npoess/decoder"
	npoessProcessor "weather-dump/src/npoess/processor"
)

func HandleInput(inputFile string, inputFormat string, outputPath string, datalink string) {
	fmt.Println("[CLI] Decoding started...", outputPath)

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

	var processorMakers = map[string]func(string) interfaces.Processor{
		"lrpt": meteorProcessor.NewProcessor,
		"hrd":  npoessProcessor.NewProcessor,
	}

	var decoderMakers = map[string]func(string) interfaces.Decoder{
		"lrpt": meteorDecoder.NewDecoder,
		"hrd":  npoessDecoder.NewDecoder,
	}

	if inputFormat == "grcout" {
		decodedFile := fmt.Sprintf("%s/decoded_%s.bin", outputPath, strings.ToLower(inputFileName))
		decoderMakers[datalink]("").Work(inputFile, decodedFile)
		inputFile = decodedFile
	}

	if _, err := os.Stat(inputFile); os.IsNotExist(err) || inputFile == "" {
		fmt.Println("[CLI] Input file not found. Exiting...\nError:", err)
		os.Exit(0)
	}

	processor := processorMakers[datalink]("")
	processor.Work(inputFile)
	processor.ExportAll(outputPath)
}
