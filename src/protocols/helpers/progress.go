package helpers

import (
	"github.com/gosuri/uiprogress"
)

type Progress struct {
	activated   bool
	progress    *uiprogress.Progress
	parserBar   *uiprogress.Bar
	composerBar *uiprogress.Bar
}

func (e *Progress) Start(parserCount, composerCount int) {
	e.progress = uiprogress.New()
	e.parserBar = e.progress.AddBar(parserCount).AppendCompleted()
	e.composerBar = e.progress.AddBar(composerCount).AppendCompleted()

	e.parserBar.PrependFunc(func(b *uiprogress.Bar) string {
		return "[DEC] Rendering channels	"
	})

	e.composerBar.PrependFunc(func(b *uiprogress.Bar) string {
		return "[DEC] Rendering composites	"
	})

	e.progress.Start()
}

func (e *Progress) Stop() {
	e.progress.Stop()
}

func (e *Progress) IncrementParser() {
	e.parserBar.Incr()
}

func (e *Progress) IncrementComposer() {
	e.composerBar.Incr()
}
