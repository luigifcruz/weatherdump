package remote

import (
	"encoding/json"
	"fmt"
	"github.com/fatih/color"
	"net/http"
	"os"
	"strings"
	"weatherdump/src/handlers"
)

type decoderRequest struct {
	InputFile  string `schema:"inputFile,required"`
	Datalink   string `schema:"datalink,required"`
	Decoder    string `schema:"decoder,required"`
	OutputPath string `schema:"outputPath"`
}

func (s *Remote) decoderHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	var req decoderRequest
	if err := decoder.Decode(&req, r.PostForm); err != nil {
		ResError(w, "INVALID_REQUEST", err.Error())
		return
	}

	if _, err := os.Stat(req.InputFile); os.IsNotExist(err) {
		ResError(w, "INPUT_FILE_NOT_FOUND", "")
		return
	}

	id := s.register()
	workingPath, fileName := handlers.GenerateDirectories(req.InputFile, req.OutputPath)
	decodedFile := fmt.Sprintf("%s/decoded_%s.bin", workingPath, strings.ToLower(fileName))

	go func() {
		s.routines[id] = make(chan bool)
		handlers.AvailableDecoders[req.Datalink][req.Decoder](id.String()).Work(req.InputFile, decodedFile, s.routines[id])
		delete(s.routines, id)
		color.Magenta("[RMT] Decoder %s exited.\n", id.String())
	}()

	req.OutputPath = decodedFile
	request, _ := json.Marshal(req)
	ResSuccess(w, id.String(), string(request))
}
