package helpers

import (
	"bytes"
	"io"
	"log"
	"os"
	"sync"
)

func ShiftWithConstantSize(arr *[]byte, pos int, length int) {
	for i := 0; i < length-pos; i++ {
		(*arr)[i] = (*arr)[pos+i]
	}
}

func WatchFor(signal chan bool, method func() bool) {
	for {
		select {
		case <-signal:
			return
		default:
			if method() {
				return
			}
		}
	}
}

func MaxIntSlice(v []int) int {
	index := 0
	max := 0
	for i, e := range v {
		if e > max {
			index = i
			max = e
		}
	}
	return index
}

func CaptureOutput(f func()) string {
	reader, writer, err := os.Pipe()
	if err != nil {
		panic(err)
	}
	stdout := os.Stdout
	stderr := os.Stderr
	defer func() {
		os.Stdout = stdout
		os.Stderr = stderr
		log.SetOutput(os.Stderr)
	}()
	os.Stdout = writer
	os.Stderr = writer
	log.SetOutput(writer)
	out := make(chan string)
	wg := new(sync.WaitGroup)
	wg.Add(1)
	go func() {
		var buf bytes.Buffer
		wg.Done()
		io.Copy(&buf, reader)
		out <- buf.String()
	}()
	wg.Wait()
	f()
	writer.Close()
	return <-out
}
