package helpers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/websocket"
)

var (
	upgrader            = websocket.Upgrader{}
	averageLastNSamples = 8192
)

type SocketConnection struct {
	sockets   *websocket.Conn
	registred bool
}

// Statistics struct for the decoder.
type Statistics struct {
	VCID                      uint8
	PacketNumber              uint32
	FrameBits                 uint16
	SignalQuality             uint8
	SyncCorrelation           uint8
	LostPackets               uint64
	AverageVitCorrections     int
	AverageRSCorrections      [4]int
	DroppedPackets            uint64
	ReceivedPacketsPerChannel [256]int64
	LostPacketsPerChannel     [256]int64
	TotalPackets              uint64
	TotalBytesRead            uint64
	TotalBytes                uint64
	SyncWord                  [4]uint8
	FrameLock                 bool
	Finished                  bool
	TaskName                  string
	Constellation             []byte
	SocketConnection
}

func (e *Statistics) Register(datalink, uuid string) {
	if uuid == "" {
		return
	}

	e.registred = true
	http.HandleFunc(fmt.Sprintf("/socket/%s/%s", datalink, uuid), func(w http.ResponseWriter, r *http.Request) {
		upgrader.CheckOrigin = func(r *http.Request) bool { return true }
		e.sockets, _ = upgrader.Upgrade(w, r, nil)
	})
}

func (e Statistics) IsRegistred() bool {
	return e.registred
}

func (e *Statistics) WaitForClient(signal chan bool) {
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

func (e *Statistics) Update() {
	if !e.IsRegistred() {
		return
	}

	if json, err := json.Marshal(e); err == nil {
		e.sockets.WriteMessage(1, []byte(json))
	}
}

// Finish decoding process.
func (e *Statistics) Finish() {
	e.Finished = true
	e.TotalBytesRead = e.TotalBytes
	e.Update()
}
