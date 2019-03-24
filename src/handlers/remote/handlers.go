package remote

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
	"weather-dump/src/assets"
	"weather-dump/src/handlers"
	"weather-dump/src/tools/img"

	uuid "github.com/satori/go.uuid"
)

func (s *Remote) decoderStart(w http.ResponseWriter, r *http.Request, vars map[string]string, id uuid.UUID) {
	inputFile := r.FormValue("inputFile")

	if handlers.AvailableDecoders[vars["datalink"]][vars["decoder"]] == nil {
		ResError(w, "INVALID_DECODER_FILE_DESCRIPTOR", "")
		return
	}

	if _, err := os.Stat(inputFile); os.IsNotExist(err) || inputFile == "" {
		ResError(w, "INPUT_FILE_NOT_FOUND", "")
		return
	}

	workingPath, fileName := handlers.GenerateDirectories(inputFile, r.FormValue("outputPath"))
	decodedFile := fmt.Sprintf("%s/decoded_%s.bin", workingPath, strings.ToLower(fileName))

	go func() {
		handlers.AvailableDecoders[vars["datalink"]][vars["decoder"]](id.String()).Work(inputFile, decodedFile, &s.processes[id].heartbeart)
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

	if r.FormValue("manifest") == "" {
		ResError(w, "INVALID_MANIFEST", "")
		return
	}

	wf := img.NewPipeline()

	var pipeline map[string]struct {
		Name      string
		Activated bool
	}
	json.Unmarshal([]byte(r.FormValue("pipeline")), &pipeline)

	for key, task := range pipeline {
		wf.AddPipe(key, task.Activated)
	}

	go func() {
		defer s.terminate(id)

		processor := handlers.AvailableProcessors[vars["datalink"]](id.String())
		processor.Work(inputFile)

		var manifest assets.ProcessingManifest
		json.Unmarshal([]byte(r.FormValue("manifest")), &manifest)
		processor.Export(workingPath, wf, manifest)
	}()

	ResSuccess(w, id.String(), workingPath)
}
