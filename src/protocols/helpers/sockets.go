package helpers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{}

type SocketConnection struct {
	sockets   *websocket.Conn
	registred bool
}

func (e *SocketConnection) Register(datalink, uuid string) {
	if uuid == "" {
		return
	}

	e.registred = true
	http.HandleFunc(fmt.Sprintf("/socket/%s/%s", datalink, uuid), func(w http.ResponseWriter, r *http.Request) {
		upgrader.CheckOrigin = func(r *http.Request) bool { return true }
		e.sockets, _ = upgrader.Upgrade(w, r, nil)
	})
}

func (e *SocketConnection) IsRegistred() bool {
	return e.registred
}

func (e *SocketConnection) SendJSON(msg interface{}) {
	if json, err := json.Marshal(msg); err == nil && e.IsRegistred() {
		e.sockets.WriteMessage(1, json)
	}
}

func (e *SocketConnection) WaitForClient(signal chan bool) {
	if !e.IsRegistred() {
		return
	}

	WatchFor(signal, func() bool {
		for e.sockets == nil {
			return false
		}
		return true
	})
}
