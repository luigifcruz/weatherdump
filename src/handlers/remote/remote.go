package remote

import (
	"fmt"
	"net/http"
	"weather-dump/src/handlers"

	httpHandlers "github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	uuid "github.com/satori/go.uuid"
)

type handler map[string]func(http.ResponseWriter, *http.Request, map[string]string, uuid.UUID)

type process struct {
	heartbeart bool
}

type Remote struct {
	processes     map[uuid.UUID]*process
	startHandlers handler
}

func New() *Remote {
	e := Remote{}
	e.processes = make(map[uuid.UUID]*process)
	e.startHandlers = handler{
		"decoder":   e.decoderStart,
		"processor": e.processorStart,
	}
	return &e
}

func (s *Remote) Listen() {
	origins := httpHandlers.AllowedOrigins([]string{"http://localhost:3002"})
	headers := httpHandlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type"})

	r := mux.NewRouter()
	r.HandleFunc("/{datalink}/{cmd}/{handler}", s.router)
	http.Handle("/", httpHandlers.CORS(origins, headers)(r))

	fmt.Println("[RMT] Starting to listen requests...")
	http.ListenAndServe("127.0.0.1:3000", nil)
}

func (s *Remote) register() uuid.UUID {
	id := uuid.Must(uuid.NewV4())
	s.processes[id] = &process{true}
	fmt.Printf("[RMT] Process registered: %s\n", id.String())
	return id
}

func (s *Remote) terminate(id uuid.UUID) {
	s.processes[id].heartbeart = false
	delete(s.processes, id)
	fmt.Printf("[RMT] Process terminated: %s\n", id.String())
}

func (s *Remote) router(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	if handlers.AvailableDecoders[vars["datalink"]] == nil {
		ResError(w, "INVALID_DATALINK", "Datalink not supported.")
		return
	}

	if s.startHandlers[vars["handler"]] == nil {
		ResError(w, "INVALID_HANDLER", "Handler not supported.")
		return
	}

	switch vars["cmd"] {
	case "abort":
		id, err := uuid.FromString(r.FormValue("id"))

		if err != nil || s.processes[id] == nil {
			ResError(w, "INVALID_ID", "Invalid ID or process already exited.")
			return
		}

		s.terminate(id)
		ResSuccess(w, "PROCESS_TERMINATED", "")
	case "start":
		s.startHandlers[vars["handler"]](w, r, vars, s.register())
	default:
		ResError(w, "INVALID_COMMAND", "Invalid command.")
	}
}
