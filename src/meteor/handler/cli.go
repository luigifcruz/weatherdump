package meteor

import (
	"fmt"
	"net/http"
	"os"
	"regexp"
	"strings"
	"time"
	"weather-dump/src/meteor/decoder"
	"weather-dump/src/meteor/processor"
)

func CommandLine(inputFile string, inputFormat string, outputPath string) {
	fmt.Println("[LRPT] Decoding started...")

	go http.ListenAndServe(":3000", nil)

	if _, err := os.Stat(outputPath); os.IsNotExist(err) {
		os.Mkdir(outputPath, os.ModePerm)
	}

	fileName := time.Now().Format(time.RFC3339)
	r, err := regexp.Compile("^.*\\/(.*)\\.\\w+$")

	if err == nil {
		fileName = r.FindStringSubmatch(inputFile)[1]
	}

	outputPath = fmt.Sprintf("%s/METEOR-LRPT-%s", outputPath, strings.ToUpper(fileName))
	os.Mkdir(outputPath, os.ModePerm)

	if inputFormat == "grcout" {
		outputFile := fmt.Sprintf("%s/decoded-%s.bin", outputPath, strings.ToLower(fileName))
		decoder := decoder.NewDecoder("")
		decoder.Work(inputFile, outputFile)
		inputFile = outputFile
	}

	if _, err := os.Stat(inputFile); os.IsNotExist(err) || inputFile == "" {
		fmt.Println("[HRD] Input file not found. Exiting...", err)
		os.Exit(0)
	}

	processor := processor.NewProcessor("")
	processor.Work(inputFile)
	processor.ExportAll(outputPath)
}
