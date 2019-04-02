package decoder

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"weather-dump/src/assets"
	"weather-dump/src/handlers/interfaces"

	"github.com/fatih/color"
	"github.com/gorilla/websocket"
	"github.com/gosuri/uiprogress"
	SatHelper "github.com/luigifreitas/libsathelper"
)

// Decoder ASM
// Synchronized + Post-Viterbi + RS Corrected + Unscrambled

type AsmDecoder struct {
	hardData     []byte
	rsWorkBuffer []byte
	reedSolomon  SatHelper.ReedSolomon
	Statistics   assets.Statistics
	constSock    *websocket.Conn
	statsSock    *websocket.Conn
}

func NewAsmDecoder(uuid string) interfaces.Decoder {
	e := AsmDecoder{}

	if uuid != "" {
		http.HandleFunc(fmt.Sprintf("/hrd/%s/constellation", uuid), e.constellation)
		http.HandleFunc(fmt.Sprintf("/hrd/%s/statistics", uuid), e.statistics)
	}

	e.hardData = make([]byte, datalink[id].FrameSize)
	e.rsWorkBuffer = make([]byte, 255)
	e.reedSolomon = SatHelper.NewReedSolomon()
	e.reedSolomon.SetCopyParityToOutput(true)
	e.Statistics.DroppedPackets = 0
	e.Statistics.TotalPackets = 1

	return &e
}

func (e *AsmDecoder) Work(inputPath string, outputPath string, g *bool) {
	color.Yellow("[DEC] WARNING! This decoder is currently in BETA development state.")

	fi, err := os.Stat(inputPath)
	input, err := os.Open(inputPath)
	output, err := os.Create(outputPath)
	if err != nil {
		log.Fatal(err)
	}

	defer input.Close()
	defer output.Close()

	e.Statistics.TotalBytes = uint64(fi.Size())
	e.Statistics.TaskName = "Decoding CADU file	"

	progress := uiprogress.New()
	progress.Start()

	bar := progress.AddBar(int(fi.Size())).AppendCompleted()

	bar.PrependFunc(func(b *uiprogress.Bar) string {
		return "[DEC] Decoding CADU file	"
	})

	bar.AppendFunc(func(b *uiprogress.Bar) string {
		s := e.Statistics
		return fmt.Sprintf("\n[DEC] Decoder Statistics	 [VCID: %2d] [#%8d]", s.VCID, s.PacketNumber)
	})

	e.Statistics.TotalBytesRead = 0
	e.Statistics.TotalBytes = uint64(fi.Size())
	e.Statistics.TaskName = "Decoding soft-symbol file"

	for *g {
		n, err := input.Read(e.hardData)
		if datalink[id].FrameSize != n {
			break
		}

		if err == nil {
			e.Statistics.TotalBytesRead += uint64(n)
			bar.Set(int(e.Statistics.TotalBytesRead))

			if e.Statistics.TotalPackets%averageLastNSamples == 0 {
				e.Statistics.AverageRSCorrections = [4]int{}
			}

			shiftWithConstantSize(&e.hardData, datalink[id].SyncWordSize, datalink[id].FrameSize-datalink[id].SyncWordSize)
			e.Statistics.TotalPackets++

			e.Statistics.VCID = e.hardData[1] & 0x3F
			e.Statistics.FrameBits = uint16(datalink[id].FrameBits)
			e.Statistics.PacketNumber = binary.BigEndian.Uint32(e.hardData[2:]) & 0xFFFFFF00 >> 8

			e.Statistics.ReceivedPacketsPerChannel[e.Statistics.VCID]++
			dat := e.hardData[:datalink[id].FrameSize-datalink[id].RsParityBlockSize-datalink[id].SyncWordSize]
			output.Write(dat)

			if e.Statistics.TotalPackets%32 == 0 && e.statsSock != nil {
				e.updateStatistics(e.Statistics)
			}
		} else {
			if err != io.EOF {
				log.Fatal(err)
			}
			break
		}
	}

	progress.Stop()
	os.Remove(outputPath + ".buf")

	if e.statsSock != nil {
		e.Statistics.Finish()
		e.updateStatistics(e.Statistics)
	}

	color.Green("[DEC] Decoding finished! File saved in the same folder.\n")
}

func (e *AsmDecoder) updateStatistics(s assets.Statistics) {
	json, err := json.Marshal(s)
	if err == nil {
		e.statsSock.WriteMessage(1, []byte(json))
	}
}

func (e *AsmDecoder) constellation(w http.ResponseWriter, r *http.Request) {
	upgrader.CheckOrigin = func(r *http.Request) bool { return true }
	e.constSock, _ = upgrader.Upgrade(w, r, nil)
}

func (e *AsmDecoder) statistics(w http.ResponseWriter, r *http.Request) {
	upgrader.CheckOrigin = func(r *http.Request) bool { return true }
	e.statsSock, _ = upgrader.Upgrade(w, r, nil)
}

func shiftWithConstantSize(arr *[]byte, pos int, length int) {
	for i := 0; i < length-pos; i++ {
		(*arr)[i] = (*arr)[pos+i]
	}
}
