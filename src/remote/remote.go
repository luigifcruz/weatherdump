package remote

import (
	"fmt"
	"net/http"
	"os"
	"regexp"
	"strings"
	"weather-dump/src/interfaces"
	meteorDecoder "weather-dump/src/meteor/decoder"
	meteorProcessor "weather-dump/src/meteor/processor"
	npoessDecoder "weather-dump/src/npoess/decoder"
	npoessProcessor "weather-dump/src/npoess/processor"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	uuid "github.com/satori/go.uuid"
)

type State struct {
	locked            bool
	workingPath       string
	activatedDatalink string
	decoderMakers     map[string]func(string) interfaces.Decoder
	processorMakers   map[string]func(string) interfaces.Processor
	processor         interfaces.Processor
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
	origins := handlers.AllowedOrigins([]string{"http://localhost:3002"})
	headers := handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type"})

	r := mux.NewRouter()
	r.HandleFunc("/api/{datalink}/register", s.register)
	r.HandleFunc("/api/{datalink}/{id}/{process}/{cmd}", s.router)
	http.Handle("/", handlers.CORS(origins, headers)(r))

	fmt.Println("[RMT] Starting to listen requests...")
	http.ListenAndServe(":3000", nil)
}

func (s *Remote) register(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	u1 := uuid.Must(uuid.NewV4())

	s.states[u1] = &State{}
	s.states[u1].locked = false
	s.states[u1].activatedDatalink = vars["datalink"]

	s.states[u1].processorMakers = map[string]func(string) interfaces.Processor{
		"lrpt": meteorProcessor.NewProcessor,
		"hrd":  npoessProcessor.NewProcessor,
	}

	s.states[u1].decoderMakers = map[string]func(string) interfaces.Decoder{
		"lrpt": meteorDecoder.NewDecoder,
		"hrd":  npoessDecoder.NewDecoder,
	}

	ResSuccess(w, u1.String(), "")
}

func (s *Remote) router(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := uuid.FromString(vars["id"])

	if err != nil || s.states[id] == nil {
		ResError(w, "INVALID_ID", "Invalid or not registed ID.")
		return
	}

	if s.states[id].decoderMakers[vars["datalink"]] == nil || vars["datalink"] != s.states[id].activatedDatalink {
		ResError(w, "INVALID_DATALINK", "Datalink not supported or differs from the activated one.")
		return
	}

	if vars["cmd"] != "abort" && s.states[id].locked {
		ResError(w, "ID_LOCKED", "Current ID is working on something. Be patient or create a new one.")
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
		ResError(w, "INVALID_PROCESS", "")
	}
}

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
