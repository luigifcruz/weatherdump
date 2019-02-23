package remote

import (
	"fmt"
	"net/http"
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
	decoderMakers     interfaces.DecoderMakers
	processorMakers   interfaces.ProcessorMakers
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
	r.HandleFunc("/api/{datalink}/terminate", s.terminate)
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

	s.states[u1].processorMakers = interfaces.ProcessorMakers{
		"lrpt": meteorProcessor.NewProcessor,
		"hrd":  npoessProcessor.NewProcessor,
	}

	s.states[u1].decoderMakers = interfaces.DecoderMakers{
		"lrpt": meteorDecoder.NewDecoder,
		"hrd":  npoessDecoder.NewDecoder,
	}

	ResSuccess(w, u1.String(), "")
}

func (e *Remote) terminate(w http.ResponseWriter, r *http.Request) {

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
