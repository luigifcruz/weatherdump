package remote

import (
	"fmt"
	"log"
	"net/http"
	"runtime/debug"

	"github.com/fatih/color"
	"github.com/google/uuid"
	httpHandlers "github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/gorilla/schema"
)

var decoder = schema.NewDecoder()

// Remote data structure.
type Remote struct {
	routines map[uuid.UUID](chan bool)
}

// New returns a pointer containing an initialized Remote data struct.
func New() *Remote {
	return &Remote{make(map[uuid.UUID](chan bool))}
}

// Listen to client requests for HTTP and WebSockets API.
func (s *Remote) Listen(serverPort, clientPort string) {
	fmt.Println("[RMT] Server listening to port " + serverPort)
	fmt.Println("[RMT] Accepting requests from client port " + clientPort)

	r := mux.NewRouter()
	r.HandleFunc("/start/processor", s.processorHandler)
	r.HandleFunc("/start/decoder", s.decoderHandler)
	r.HandleFunc("/abort/{id}", s.abortHandler)
	r.HandleFunc("/get/manifest", s.manifestHandler)
	r.HandleFunc("/get/thumbnail", s.thumbnailHandler)

	origins := httpHandlers.AllowedOrigins([]string{"http://localhost:" + clientPort})
	headers := httpHandlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type"})
	http.Handle("/", httpHandlers.CORS(origins, headers)(r))
	log.Fatal(http.ListenAndServe("127.0.0.1:"+serverPort, nil))
}

func (s *Remote) register() uuid.UUID {
	id := uuid.Must(uuid.NewRandom())
	color.Magenta("[RMT] Process registered: %s\n", id.String())
	return id
}

func (s *Remote) terminate(id uuid.UUID) {
	s.routines[id] <- true
	delete(s.routines, id)
	color.Magenta("[RMT] Process terminated: %s\n", id.String())
	debug.FreeOSMemory()
}

func (s *Remote) abortHandler(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(mux.Vars(r)["id"])

	if err != nil || s.routines[id] == nil {
		ResError(w, "INVALID_ID", "Invalid ID or process already exited.")
		return
	}

	s.terminate(id)
	ResSuccess(w, "PROCESS_TERMINATED", "")
}
