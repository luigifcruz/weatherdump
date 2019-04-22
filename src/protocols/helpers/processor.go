package helpers

import (
	"sort"

	"github.com/fatih/color"
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

type ManifestList map[uint16]*Manifest

type ProcessingManifest struct {
	Parser   ManifestList
	Composer ManifestList
	SocketConnection
	Progress
}

func (e *ProcessingManifest) Start() {
	if !e.IsRegistred() {
		e.Progress.Start(len(e.Parser), len(e.Composer))
	}
	e.WaitForClient(nil)
}

func (e *ProcessingManifest) Stop(re string) {
	e.Progress.Stop()
	color.Magenta(re)
}

func (e *ProcessingManifest) Update() {
	e.SocketConnection.SendJSON(e)
}

func (e *ProcessingManifest) ParserCompleted(key uint16) {
	e.Update()
	e.Parser[key].Completed()
	e.Progress.IncrementParser()
}

func (e *ProcessingManifest) ComposerCompleted(key uint16) {
	e.Update()
	e.Composer[key].Completed()
	e.Progress.IncrementComposer()
}

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
