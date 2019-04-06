package decoder

import (
	"encoding/binary"
	"fmt"
	"io"
	"log"
	"os"
	"weather-dump/src/handlers/interfaces"
	"weather-dump/src/protocols/helpers"

	"github.com/fatih/color"
	"github.com/gosuri/uiprogress"
	SatHelper "github.com/luigifreitas/libsathelper"
)

// Decoder ASM
// Synchronized + Post-Viterbi + RS Corrected + Unscrambled

type AsmDecoder struct {
	hardData     []byte
	rsWorkBuffer []byte
	reedSolomon  SatHelper.ReedSolomon
	Statistics   helpers.Statistics
}

func NewAsmDecoder(uuid string) interfaces.Decoder {
	e := AsmDecoder{}

	e.Statistics.Register("hrd", uuid)

	e.hardData = make([]byte, datalink[id].FrameSize)
	e.rsWorkBuffer = make([]byte, 255)
	e.reedSolomon = SatHelper.NewReedSolomon()
	e.reedSolomon.SetCopyParityToOutput(true)
	e.Statistics.DroppedPackets = 0
	e.Statistics.TotalPackets = 1

	return &e
}

func (e *AsmDecoder) Work(inputPath string, outputPath string, signal chan bool) {
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

	if !e.Statistics.IsRegistred() {
		progress.Start()
	}

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

	e.Statistics.WaitForClient(signal)

	helpers.WatchFor(signal, func() bool {
		n, err := input.Read(e.hardData)
		if datalink[id].FrameSize != n || err != nil {
			if err != io.EOF && err != nil {
				log.Fatal(err)
			}
			return true
		}

		e.Statistics.TotalBytesRead += uint64(n)
		bar.Set(int(e.Statistics.TotalBytesRead))

		if e.Statistics.TotalPackets%averageLastNSamples == 0 {
			e.Statistics.AverageRSCorrections = [4]int{}
		}

		helpers.ShiftWithConstantSize(&e.hardData, datalink[id].SyncWordSize, datalink[id].FrameSize-datalink[id].SyncWordSize)
		e.Statistics.TotalPackets++

		e.Statistics.VCID = e.hardData[1] & 0x3F
		e.Statistics.FrameBits = uint16(datalink[id].FrameBits)
		e.Statistics.PacketNumber = binary.BigEndian.Uint32(e.hardData[2:]) & 0xFFFFFF00 >> 8

		e.Statistics.ReceivedPacketsPerChannel[e.Statistics.VCID]++
		dat := e.hardData[:datalink[id].FrameSize-datalink[id].RsParityBlockSize-datalink[id].SyncWordSize]
		output.Write(dat)

		if e.Statistics.TotalPackets%512 == 0 {
			e.Statistics.Update()
		}

		return false
	})

	os.Remove(outputPath + ".buf")

	color.Green("[DEC] Decoding finished! File saved in the same folder.\n")

	e.Statistics.Finish()

	if !e.Statistics.IsRegistred() {
		progress.Stop()
	}
}
