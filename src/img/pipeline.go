package img

import (
	"reflect"
	"sort"
	"strings"
)

type Pipeline struct {
	taskCount  int
	target     *Img
	exceptions map[string]int
	currTasks  map[string]int
}

func NewPipeline() Pipeline {
	return Pipeline{
		currTasks:  make(map[string]int),
		exceptions: make(map[string]int),
	}
}

func (e *Pipeline) AddException(method string, enabled bool) {
	if !enabled {
		e.exceptions[method] = -1
	}
}

func (e *Pipeline) ResetExceptions() {
	e.exceptions = make(map[string]int)
}

func (e *Pipeline) AddPipe(method string, enabled bool) {
	if enabled {
		e.currTasks[method] = e.taskCount
		e.taskCount++
	}
}

func (e *Pipeline) Target(img Img) *Pipeline {
	e.target = &img
	return e
}

func (e *Pipeline) Process() *Pipeline {
	for _, task := range getKeys(e.currTasks) {
		if !strings.Contains(task, "Export") && e.exceptions[task] == 0 {
			reflect.ValueOf(*e.target).MethodByName(task).Call([]reflect.Value{})
		}
	}
	return e
}

func (e *Pipeline) Export(args ...interface{}) *Pipeline {
	inputs := make([]reflect.Value, len(args))
	for i := range args {
		inputs[i] = reflect.ValueOf(args[i])
	}

	for _, task := range getKeys(e.currTasks) {
		if strings.Contains(task, "Export") && e.exceptions[task] == 0 {
			reflect.ValueOf(*e.target).MethodByName(task).Call(inputs)
		}
	}
	return e
}

func getKeys(tasks map[string]int) []string {
	keys := make([]string, 0, len(tasks))
	for k := range tasks {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}
