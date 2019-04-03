package remote

import (
	"fmt"
	"log"
	"net/http"

	httpHandlers "github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/gorilla/schema"
	uuid "github.com/satori/go.uuid"
)

var decoder = schema.NewDecoder()

type Remote struct {
	routines map[uuid.UUID](chan bool)
}

func New() *Remote {
	return &Remote{make(map[uuid.UUID](chan bool))}
}

func (s *Remote) Listen(port string) {
	fmt.Println("[RMT] Starting to listen requests from port " + port + "...")

	r := mux.NewRouter()
	r.HandleFunc("/start/processor", s.processorHandler)
	r.HandleFunc("/start/decoder", s.decoderHandler)
	r.HandleFunc("/abort/{id}", s.abortHandler)
	r.HandleFunc("/get/manifest", s.manifestHandler)

	origins := httpHandlers.AllowedOrigins([]string{"http://localhost:3002"})
	headers := httpHandlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type"})
	http.Handle("/", httpHandlers.CORS(origins, headers)(r))
	log.Fatal(http.ListenAndServe("127.0.0.1:"+port, nil))
}

func (s *Remote) register() uuid.UUID {
	id := uuid.Must(uuid.NewV4(), nil)
	fmt.Printf("[RMT] Process registered: %s\n", id.String())
	return id
}

func (s *Remote) terminate(id uuid.UUID) {
	s.routines[id] <- true
	delete(s.routines, id)
	fmt.Printf("[RMT] Process terminated: %s\n", id.String())
}

func (s *Remote) abortHandler(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.FromString(mux.Vars(r)["id"])

	if err != nil || s.routines[id] == nil {
		ResError(w, "INVALID_ID", "Invalid ID or process already exited.")
		return
	}

	s.terminate(id)
	ResSuccess(w, "PROCESS_TERMINATED", "")
}
