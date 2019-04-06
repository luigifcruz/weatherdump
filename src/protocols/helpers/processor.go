package helpers

import (
	"github.com/gorilla/websocket"
	"sort"
)

type Manifest map[uint16]*struct {
	Name        string
	Description string
	Filename    string
	Activated   bool
	Finished    bool
}

type ProcessingManifest struct {
	Parser   Manifest
	Composer Manifest
}

func (e ProcessingManifest) ParserCount() int {
	return len(e.Parser)
}

func (e ProcessingManifest) ComposerCount() int {
	return len(e.Composer)
}

func (e *Manifest) File(i uint16, filename string) {
	(*e)[i].Filename = filename
}

func (e *Manifest) Completed(i uint16, socket *websocket.Conn) {
	(*e)[i].Finished = true
	if socket != nil {

		socket.WriteJSON(e)
	}
}

func (e Manifest) Ordered() []uint16 {
	keys := make([]int, 0, len(e))
	for k := range e {
		keys = append(keys, int(k))
	}
	sort.Ints(keys)
	res := make([]uint16, 0, len(e))
	for _, k := range keys {
		res = append(res, uint16(k))
	}
	return res
}
