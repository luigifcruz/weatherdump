package remote

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"weather-dump/src/assets"
	"weather-dump/src/handlers"
	"weather-dump/src/img"
)

type processorRequest struct {
	InputFile  string `schema:"inputPath,required"`
	Datalink   string `schema:"datalink,required"`
	Pipeline   string `schema:"pipeline,required"`
	Manifest   string `schema:"manifest,required"`
	OutputPath string `schema:"outputPath"`
}

func (s *Remote) processorHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	var req processorRequest
	if err := decoder.Decode(&req, r.PostForm); err != nil {
		ResError(w, "INVALID_REQUEST", err.Error())
		return
	}

	if _, err := os.Stat(req.InputFile); os.IsNotExist(err) {
		ResError(w, "INPUT_FILE_NOT_FOUND", "")
		return
	}

	id := s.register()
	wf := img.NewPipeline()
	req.OutputPath, _ = handlers.GenerateDirectories(req.InputFile, req.OutputPath)

	var p map[string]struct {
		Name      string
		Activated bool
	}
	json.Unmarshal([]byte(req.Pipeline), &p)

	for key, task := range p {
		wf.AddPipe(key, task.Activated)
	}

	go func() {
		processor := handlers.AvailableProcessors[req.Datalink](id.String())
		processor.Work(req.InputFile)

		var m assets.ProcessingManifest
		json.Unmarshal([]byte(req.Manifest), &m)
		processor.Export(req.OutputPath, wf, m)

		fmt.Printf("[RMT] Decoder %s exited.\n", id.String())
		runtime.Free()
	}()

	request, _ := json.Marshal(req)
	ResSuccess(w, id.String(), string(request))
}
