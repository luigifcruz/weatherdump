package remote

import (
	"fmt"
	"net/http"
	"os"
	"regexp"
	"strings"

	uuid "github.com/satori/go.uuid"
)

func (s *Remote) decoderHandler(w http.ResponseWriter, r *http.Request, vars map[string]string, id uuid.UUID) {
	inputFile := r.FormValue("inputFile")

	if _, err := os.Stat(inputFile); os.IsNotExist(err) || inputFile == "" {
		ResError(w, "INPUT_FILE_NOT_FOUND", "")
		return
	}

	var fileName string
	m, err := regexp.Compile("^.*\\/(.*)\\.\\w+$")
	if err == nil {
		fileName = m.FindStringSubmatch(inputFile)[1]
	}
	outputFile := fmt.Sprintf("%s/decoded-%s.bin", s.states[id].workingPath, strings.ToLower(fileName))

	go func() {
		s.states[id].locked = true
		s.states[id].decoderMakers[vars["datalink"]](id.String()).Work(inputFile, outputFile)
		s.states[id].locked = false
	}()

	ResSuccess(w, "DECODER_STARTED", outputFile)
}

func (s *Remote) processorHandler(w http.ResponseWriter, r *http.Request, vars map[string]string, id uuid.UUID) {
	inputFile := r.FormValue("inputFile")

	if _, err := os.Stat(inputFile); os.IsNotExist(err) || inputFile == "" {
		ResError(w, "INPUT_FILE_NOT_FOUND", "")
		return
	}

	go func() {
		s.states[id].locked = true
		s.states[id].processor = s.states[id].processorMakers[vars["datalink"]](id.String())
		s.states[id].processor.Work(inputFile)
		s.states[id].locked = false
	}()

	ResSuccess(w, "PROCESSOR_STARTED", "")
}

func (s *Remote) exporterHandler(w http.ResponseWriter, r *http.Request, vars map[string]string, id uuid.UUID) {
	if s.states[id].processor == nil {
		ResError(w, "PROCESSOR_NOT_LOADED", "")
		return
	}

	go func() {
		s.states[id].locked = true
		go s.states[id].processor.ExportAll(s.states[id].workingPath)
		s.states[id].locked = false
	}()

	ResSuccess(w, "EXPORTER_STARTED", "")
}
