package remote

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
	"weather-dump/src/handlers"
	"weather-dump/src/tools/img"

	uuid "github.com/satori/go.uuid"
)

func (s *Remote) decoderStart(w http.ResponseWriter, r *http.Request, vars map[string]string, id uuid.UUID) {
	inputFile := r.FormValue("inputFile")

	if _, err := os.Stat(inputFile); os.IsNotExist(err) || inputFile == "" {
		ResError(w, "INPUT_FILE_NOT_FOUND", "")
		return
	}

	workingPath, fileName := handlers.GenerateDirectories(inputFile, r.FormValue("outputPath"))
	decodedFile := fmt.Sprintf("%s/decoded_%s.bin", workingPath, strings.ToLower(fileName))

	go func() {
		handlers.AvailableDecoders[vars["datalink"]](id.String()).Work(inputFile, decodedFile, &s.processes[id].heartbeart)
		s.terminate(id)
	}()

	ResSuccess(w, id.String(), decodedFile)
}

func (s *Remote) processorStart(w http.ResponseWriter, r *http.Request, vars map[string]string, id uuid.UUID) {
	inputFile := r.FormValue("inputFile")
	workingPath, _ := handlers.GenerateDirectories(inputFile, r.FormValue("outputPath"))

	if _, err := os.Stat(inputFile); os.IsNotExist(err) || inputFile == "" {
		ResError(w, "INPUT_FILE_NOT_FOUND", "")
		return
	}

	wf := img.NewPipeline()

	var pipeline map[string]bool
	json.Unmarshal([]byte(r.FormValue("pipeline")), &pipeline)

	for task, enabled := range pipeline {
		wf.AddPipe(task, enabled)
	}

	go func() {
		processor := handlers.AvailableProcessors[vars["datalink"]](id.String())
		processor.Work(inputFile)
		processor.Export(workingPath, wf)
		s.terminate(id)
	}()

	ResSuccess(w, "PROCESSOR_STARTED", workingPath)
}
