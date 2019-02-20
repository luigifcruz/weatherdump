package remote

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"regexp"
	"strings"
	"time"
	meteorDecoder "weather-dump/src/meteor/decoder"
	meteorProcessor "weather-dump/src/meteor/processor"
	npoessDecoder "weather-dump/src/npoess/decoder"
	npoessProcessor "weather-dump/src/npoess/processor"

	"github.com/gorilla/mux"
	uuid "github.com/satori/go.uuid"
)

type State struct {
	locked            bool
	workingPath       string
	activatedDatalink string
	npoessDecoder     *npoessDecoder.Worker
	npoessProcessor   *npoessProcessor.Worker
	meteorDecoder     *meteorDecoder.Worker
	meteorProcessor   *meteorProcessor.Worker
}

type Remote struct {
	states map[uuid.UUID]*State
}

func New() *Remote {
	e := Remote{}
	e.states = make(map[uuid.UUID]*State)
	return &e
}

func (s *Remote) Listen() {
	r := mux.NewRouter()
	r.HandleFunc("/{datalink}/register", s.register)
	r.HandleFunc("/{datalink}/{id}/{process}/{cmd}", s.router)
	http.Handle("/", r)

	fmt.Println("[REMOTE] Starting to listen requests...")
	http.ListenAndServe(":3000", nil)
}

func (s *Remote) register(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	u1 := uuid.Must(uuid.NewV4())

	s.states[u1] = &State{}
	s.states[u1].locked = false
	s.states[u1].activatedDatalink = vars["datalink"]

	switch vars["datalink"] {
	case "meteor":
		s.states[u1].workingPath = fmt.Sprintf("%s/METEOR-LRPT-%s", r.FormValue("workingPath"), time.Now().Format(time.RFC3339))
	case "npoess":
		s.states[u1].workingPath = fmt.Sprintf("%s/NPOESS-HRD-%s", r.FormValue("workingPath"), time.Now().Format(time.RFC3339))
	}
	os.Mkdir(s.states[u1].workingPath, os.ModePerm)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(SuccessResponse{true, u1.String(), ""})
}

func (s *Remote) router(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	vars := mux.Vars(r)
	id, err := uuid.FromString(vars["id"])

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorResponse{false, "INVALID_ID", err.Error()})
		return
	}

	if s.states[id] == nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorResponse{false, "ID_NOT_REGISTERED", ""})
		return
	}

	if vars["cmd"] != "abort" && s.states[id].locked {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorResponse{false, "ID_LOCKED", ""})
		return
	}

	if vars["datalink"] != s.states[id].activatedDatalink {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorResponse{false, "WRONG_DATALINK", ""})
		return
	}

	switch vars["process"] {
	case "decoder":
		s.decoderHandler(w, r, vars, id)
	case "processor":
		s.processorHandler(w, r, vars, id)
	case "exporter":
		s.exporterHandler(w, r, vars, id)
	default:
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorResponse{false, "INVALID_PROCESS", ""})
	}
}

func (s *Remote) decoderHandler(w http.ResponseWriter, r *http.Request, vars map[string]string, id uuid.UUID) {
	inputFile := r.FormValue("inputFile")
	if _, err := os.Stat(inputFile); os.IsNotExist(err) || inputFile == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorResponse{false, "INPUT_FILE_NOT_FOUND", ""})
		return
	}

	var fileName string
	m, err := regexp.Compile("^.*\\/(.*)\\.\\w+$")
	if err == nil {
		fileName = m.FindStringSubmatch(inputFile)[1]
	}
	outputFile := fmt.Sprintf("%s/decoded-%s.bin", s.states[id].workingPath, strings.ToLower(fileName))

	go func() {
		fmt.Println("[REMOTE] New goroutine spawned! ID:", id.String())
		s.states[id].locked = true

		switch vars["datalink"] {
		case "meteor":
			s.states[id].meteorDecoder = meteorDecoder.NewDecoder(id.String())
			s.states[id].meteorDecoder.Work(inputFile, outputFile)
		case "npoess":
			s.states[id].npoessDecoder = npoessDecoder.NewDecoder(id.String())
			s.states[id].npoessDecoder.Work(inputFile, outputFile)
		}

		s.states[id].locked = false
	}()

	json.NewEncoder(w).Encode(SuccessResponse{true, "DECODER_STARTED", outputFile})
}

func (s *Remote) processorHandler(w http.ResponseWriter, r *http.Request, vars map[string]string, id uuid.UUID) {
	inputFile := r.FormValue("inputFile")
	if _, err := os.Stat(inputFile); os.IsNotExist(err) || inputFile == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorResponse{false, "INPUT_FILE_NOT_FOUND", ""})
		return
	}

	go func() {
		fmt.Println("[REMOTE] New goroutine spawned! ID:", id.String())
		s.states[id].locked = true

		switch vars["datalink"] {
		case "meteor":
			s.states[id].meteorProcessor = meteorProcessor.NewProcessor(id.String())
			s.states[id].meteorProcessor.Work(inputFile)
		case "npoess":
			s.states[id].npoessProcessor = npoessProcessor.NewProcessor(id.String())
			s.states[id].npoessProcessor.Work(inputFile)
		}

		s.states[id].locked = false
	}()

	json.NewEncoder(w).Encode(SuccessResponse{true, "PROCESSOR_STARTED", ""})
}

func (s *Remote) exporterHandler(w http.ResponseWriter, r *http.Request, vars map[string]string, id uuid.UUID) {
	go func() {
		fmt.Println("[REMOTE] New goroutine spawned! ID:", id.String())

		switch vars["datalink"] {
		case "meteor":
			s.states[id].meteorProcessor.ExportAll(s.states[id].workingPath)
		case "npoess":
			s.states[id].npoessProcessor.ExportAll(s.states[id].workingPath)
		}
	}()

	json.NewEncoder(w).Encode(SuccessResponse{true, "EXPORTER_STARTED", ""})
}
