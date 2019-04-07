package helpers

import (
	"sort"
)

type Manifest struct {
	Name        string
	Description string
	Filename    string
	Activated   bool
	Finished    bool
}

func (e *Manifest) FileName(filename string) {
	e.Filename = filename
}

func (e *Manifest) Completed() {
	e.Finished = true
}

type ProcessingManifest struct {
	Parser   ManifestList
	Composer ManifestList
	SocketConnection
}

func (e ProcessingManifest) ParserCount() int {
	return len(e.Parser)
}

func (e ProcessingManifest) ComposerCount() int {
	return len(e.Composer)
}

func (e *ProcessingManifest) Update() {
	e.SocketConnection.SendJSON(e)
}

type ManifestList map[uint16]*Manifest

func (e ManifestList) Parse() []uint16 {
	keys := make([]int, 0, len(e))
	for k := range e {
		if (*e[k]).Activated {
			keys = append(keys, int(k))
		}
	}
	sort.Ints(keys)
	res := make([]uint16, 0, len(e))
	for _, k := range keys {
		res = append(res, uint16(k))
	}
	return res
}
