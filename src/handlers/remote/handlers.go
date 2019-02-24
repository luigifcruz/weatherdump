package remote

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"weather-dump/src/handlers"

	uuid "github.com/satori/go.uuid"
)

func generateDirectories(inputFile string) (string, string) {
	inputFileName := filepath.Base(inputFile)
	inputFileName = strings.TrimSuffix(inputFileName, filepath.Ext(inputFile))
	workingPath := filepath.Dir(inputFile)

	if !strings.Contains(inputFile, "/OUTPUT_") {
		workingPath = fmt.Sprintf("%s/OUTPUT_%s", filepath.Dir(inputFile), strings.ToUpper(inputFileName))
		if _, err := os.Stat(workingPath); os.IsNotExist(err) {
			os.Mkdir(workingPath, os.ModePerm)
		}
	}

	return workingPath, inputFileName
}

func (s *Remote) decoderStart(w http.ResponseWriter, r *http.Request, vars map[string]string, id uuid.UUID) {
	inputFile := r.FormValue("inputFile")

	if _, err := os.Stat(inputFile); os.IsNotExist(err) || inputFile == "" {
		ResError(w, "INPUT_FILE_NOT_FOUND", "")
		return
	}

	workingPath, fileName := generateDirectories(inputFile)
	decodedFile := fmt.Sprintf("%s/decoded_%s.bin", workingPath, strings.ToLower(fileName))

	go func() {
		handlers.AvailableDecoders[vars["datalink"]](id.String()).Work(inputFile, decodedFile, &s.processes[id].heartbeart)
	}()

	ResSuccess(w, id.String(), decodedFile)
}

func (s *Remote) processorStart(w http.ResponseWriter, r *http.Request, vars map[string]string, id uuid.UUID) {
	inputFile := r.FormValue("inputFile")
	workingPath, _ := generateDirectories(inputFile)

	if _, err := os.Stat(inputFile); os.IsNotExist(err) || inputFile == "" {
		ResError(w, "INPUT_FILE_NOT_FOUND", "")
		return
	}

	go func() {
		processor := handlers.AvailableProcessors[vars["datalink"]](id.String())
		processor.Work(inputFile)
		processor.ExportAll(workingPath)
	}()

	ResSuccess(w, "PROCESSOR_STARTED", "")
}
