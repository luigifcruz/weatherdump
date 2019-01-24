package Handler

import (
	"fmt"
	"net/http"
	"os"
	"regexp"
	"strings"
	"time"
	"weather-dump/src/Meteor/Decoder"
)

func CommandLine(inputPath string, inputFormat string, outputFolder string) {
	fmt.Println("[LRPT] Decoding started...")

	go http.ListenAndServe(":3000", nil)

	if _, err := os.Stat(outputFolder); os.IsNotExist(err) {
		os.Mkdir(outputFolder, os.ModePerm)
	}

	fileName := time.Now().Format(time.RFC3339)
	r, err := regexp.Compile("^.*\\/(.*)\\.\\w+$")

	if err == nil {
		fileName = r.FindStringSubmatch(inputPath)[1]
	}

	outputFolder = fmt.Sprintf("%s/METEOR-LRPT-%s", outputFolder, strings.ToUpper(fileName))
	os.Mkdir(outputFolder, os.ModePerm)

	if inputFormat == "grcout" {
		dec := Decoder.NewDecoder()
		outputFile := fmt.Sprintf("%s/decoded-%s.bin", outputFolder, strings.ToLower(fileName))
		dec.DecodeFile(inputPath, outputFile)
		inputPath = outputFile
	}
}
