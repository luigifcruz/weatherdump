package decoder

import (
	"fmt"
	"io"
	"log"
	"os"
	"weather-dump/src/handlers/interfaces"
	"weather-dump/src/protocols/helpers"

	"encoding/binary"

	"github.com/fatih/color"
	"github.com/gosuri/uiprogress"
	SatHelper "github.com/luigifreitas/libsathelper"
)

type SoftSymbolDecoder struct {
	viterbiData  []byte
	decodedData  []byte
	lastFrameEnd []byte
	codedData    []byte
	rsWorkBuffer []byte
	viterbi      SatHelper.Viterbi27
	reedSolomon  SatHelper.ReedSolomon
	correlator   SatHelper.Correlator
	packetFixer  SatHelper.PacketFixer
	Statistics   helpers.Statistics
}

func NewSoftSymbolDecoder(uuid string) interfaces.Decoder {
	e := SoftSymbolDecoder{}

	e.Statistics.Register("hrd", uuid)

	if uselastFrameData {
		e.viterbiData = make([]byte, datalink[id].CodedFrameSize+lastFrameDataBits)
		e.decodedData = make([]byte, datalink[id].FrameSize+lastFrameData)
		e.lastFrameEnd = make([]byte, lastFrameDataBits)
		e.viterbi = SatHelper.NewViterbi27(datalink[id].FrameBits + lastFrameDataBits)

		for i := 0; i < lastFrameDataBits; i++ {
			e.lastFrameEnd[i] = 128
		}
	} else {
		e.viterbiData = make([]byte, datalink[id].CodedFrameSize)
		e.decodedData = make([]byte, datalink[id].FrameSize)
		e.viterbi = SatHelper.NewViterbi27(datalink[id].FrameBits)
	}

	e.codedData = make([]byte, datalink[id].CodedFrameSize)
	e.rsWorkBuffer = make([]byte, 255)

	e.reedSolomon = SatHelper.NewReedSolomon()
	e.correlator = SatHelper.NewCorrelator()
	e.packetFixer = SatHelper.NewPacketFixer()

	e.reedSolomon.SetCopyParityToOutput(true)

	e.correlator.AddWord(datalink[id].SyncWords[0])
	e.correlator.AddWord(datalink[id].SyncWords[1])
	e.correlator.AddWord(datalink[id].SyncWords[2])
	e.correlator.AddWord(datalink[id].SyncWords[3])

	e.correlator.AddWord(datalink[id].SyncWords[4])
	e.correlator.AddWord(datalink[id].SyncWords[5])
	e.correlator.AddWord(datalink[id].SyncWords[6])
	e.correlator.AddWord(datalink[id].SyncWords[7])

	return &e
}

func (e *SoftSymbolDecoder) Work(inputPath string, outputPath string, signal chan bool) {
	var phaseShift SatHelper.SatHelperPhaseShift
	flywheelCount := 0

	fi, err := os.Stat(inputPath)
	output, err := os.Create(outputPath)
	input, err := os.Open(inputPath)
	if err != nil {
		log.Fatal(err)
	}

	defer input.Close()
	defer output.Close()

	e.Statistics.TotalBytes = uint64(fi.Size())
	e.Statistics.TaskName = "Decoding soft-symbol file"

	progress := uiprogress.New()

	if !e.Statistics.IsRegistred() {
		progress.Start()
	}

	bar := progress.AddBar(int(fi.Size())).AppendCompleted()
	bar.PrependFunc(func(b *uiprogress.Bar) string {
		return "[DEC] Decoding soft-symbol file	"
	})
	bar.AppendFunc(func(b *uiprogress.Bar) string {
		s := e.Statistics
		return fmt.Sprintf("\n[DEC] Decoder Statistics	 [VCID: %2d] [VIT: %5d] [QUAL: %2d%%] [RS: %2d %2d %2d %2d] [DROPPED: %4.1f%%]",
			s.VCID, s.AverageVitCorrections, s.SignalQuality,
			s.AverageRSCorrections[0], s.AverageRSCorrections[1], s.AverageRSCorrections[2], s.AverageRSCorrections[3],
			float32(s.DroppedPackets)/float32(s.TotalPackets)*100)
	})

	e.Statistics.WaitForClient(signal)

	helpers.WatchFor(signal, func() bool {
		n, err := input.Read(e.codedData)
		if datalink[id].CodedFrameSize != n || err != nil {
			if err != io.EOF {
				log.Fatal(err)
			}
			return true
		}

		e.Statistics.TotalBytesRead += uint64(n)
		bar.Set(int(e.Statistics.TotalBytesRead))

		if e.Statistics.TotalPackets%32 == 0 {
			e.Statistics.Constellation = e.codedData[:200]
			e.Statistics.Update()
		}

		if e.Statistics.TotalPackets%averageLastNSamples == 0 {
			e.Statistics.AverageRSCorrections = [4]int{}
			e.Statistics.AverageVitCorrections = 0
		}

		if flywheelCount == defaultFlywheelRecheck {
			e.Statistics.FrameLock = false
			flywheelCount = 0
		}

		if !e.Statistics.FrameLock {
			e.correlator.Correlate(&e.codedData[0], uint(datalink[id].CodedFrameSize))
		} else {
			e.correlator.Correlate(&e.codedData[0], uint(datalink[id].CodedFrameSize)/64)
			if e.correlator.GetHighestCorrelationPosition() != 0 {
				e.correlator.Correlate(&e.codedData[0], uint(datalink[id].CodedFrameSize))
				flywheelCount = 0
			}
		}
		flywheelCount++

		cor := e.correlator.GetHighestCorrelation()
		word := e.correlator.GetCorrelationWordNumber()
		pos := e.correlator.GetHighestCorrelationPosition()

		if cor > datalink[id].MinCorrelationBits {
			iqInv := (word / 4) > 0
			switch word % 4 {
			case 0:
				phaseShift = SatHelper.DEG_0
			case 1:
				phaseShift = SatHelper.DEG_90
			case 2:
				phaseShift = SatHelper.DEG_180
			case 3:
				phaseShift = SatHelper.DEG_270
			}

			if pos != 0 {
				helpers.ShiftWithConstantSize(&e.codedData, int(pos), datalink[id].CodedFrameSize)
				offset := datalink[id].CodedFrameSize - int(pos)

				buffer := make([]byte, int(pos))
				n, err = input.Read(buffer)

				e.Statistics.TotalBytesRead += uint64(n)
				bar.Set(int(e.Statistics.TotalBytesRead))
				if err != nil {
					fmt.Println(err)
					return true
				}

				for i := offset; i < datalink[id].CodedFrameSize; i++ {
					e.codedData[i] = buffer[i-offset]
				}
			}

			e.packetFixer.FixPacket(&e.codedData[0], uint(datalink[id].CodedFrameSize), phaseShift, iqInv)

			if uselastFrameData {
				for i := 0; i < lastFrameDataBits; i++ {
					e.viterbiData[i] = e.lastFrameEnd[i]
				}
				for i := lastFrameDataBits; i < datalink[id].CodedFrameSize+lastFrameDataBits; i++ {
					e.viterbiData[i] = e.codedData[i-lastFrameDataBits]
				}
			} else {
				for i := 0; i < datalink[id].CodedFrameSize; i++ {
					e.viterbiData[i] = e.codedData[i]
				}
			}

			e.viterbi.Decode(&e.viterbiData[0], &e.decodedData[0])

			nrzmDecodeSize := datalink[id].FrameSize

			if uselastFrameData {
				nrzmDecodeSize += lastFrameData
			}

			SatHelper.DifferentialEncodingNrzmDecode(&e.decodedData[0], nrzmDecodeSize)

			signalErrors := float32(e.viterbi.GetPercentBER())
			signalErrors = 100 - (signalErrors * 10)

			if uselastFrameData {
				helpers.ShiftWithConstantSize(&e.decodedData, lastFrameData/2, datalink[id].FrameSize+lastFrameData/2)
				for i := 0; i < lastFrameDataBits; i++ {
					e.lastFrameEnd[i] = e.viterbiData[datalink[id].CodedFrameSize+i]
				}
			}

			for i := 0; i < datalink[id].SyncWordSize; i++ {
				e.Statistics.SyncWord[i] = e.decodedData[i]
			}

			helpers.ShiftWithConstantSize(&e.decodedData, datalink[id].SyncWordSize, datalink[id].FrameSize-datalink[id].SyncWordSize)

			e.Statistics.TotalPackets++

			SatHelper.DeRandomizerDeRandomize(&e.decodedData[0], datalink[id].FrameSize-datalink[id].SyncWordSize)

			var derrors [4]int
			for i := 0; i < datalink[id].RsBlocks; i++ {
				e.reedSolomon.Deinterleave(&e.decodedData[0], &e.rsWorkBuffer[0], byte(i), byte(datalink[id].RsBlocks))
				derrors[i] = int(int8(e.reedSolomon.Decode_ccsds(&e.rsWorkBuffer[0])))
				e.reedSolomon.Interleave(&e.rsWorkBuffer[0], &e.decodedData[0], byte(i), byte(datalink[id].RsBlocks))
				if derrors[i] != -1 {
					e.Statistics.AverageRSCorrections[i] = (e.Statistics.AverageRSCorrections[i] + derrors[i]) / 2
				}
			}

			if derrors[0] == -1 && derrors[1] == -1 && derrors[2] == -1 && derrors[3] == -1 {
				e.Statistics.AverageRSCorrections = [4]int{-1, -1, -1, -1}
				e.Statistics.FrameLock = false
				e.Statistics.DroppedPackets++
			} else {
				e.Statistics.FrameLock = true
			}

			e.Statistics.VCID = e.decodedData[1] & 0x3F
			e.Statistics.FrameBits = uint16(datalink[id].FrameBits)
			e.Statistics.PacketNumber = binary.BigEndian.Uint32(e.decodedData[2:]) & 0xFFFFFF00 >> 8
			e.Statistics.SignalQuality = uint8(signalErrors)
			e.Statistics.SyncCorrelation = uint8(cor)

			if e.Statistics.SignalQuality > 100 || !e.Statistics.FrameLock {
				e.Statistics.SignalQuality = 0
			}

			e.Statistics.AverageVitCorrections += e.viterbi.GetBER()
			e.Statistics.AverageVitCorrections /= 2

			if e.Statistics.FrameLock {
				e.Statistics.ReceivedPacketsPerChannel[e.Statistics.VCID]++
				dat := e.decodedData[:datalink[id].FrameSize-datalink[id].RsParityBlockSize-datalink[id].SyncWordSize]
				output.Write(dat)
			}
		}

		return false
	})

	color.Green("[DEC] Decoding finished! File saved in the same folder.\n")

	e.Statistics.Finish()

	if !e.Statistics.IsRegistred() {
		progress.Stop()
	}
}
