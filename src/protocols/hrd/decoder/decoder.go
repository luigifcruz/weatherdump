package decoder

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"weather-dump/src/handlers/interfaces"

	b64 "encoding/base64"
	"encoding/json"

	SatHelper "github.com/OpenSatelliteProject/libsathelper"
	"github.com/gorilla/websocket"
)

const defaultFlywheelRecheck = 1024
const averageLastNSamples = 8192
const lastFrameDataBits = 64
const lastFrameData = lastFrameDataBits / 8
const uselastFrameData = true
const id = "HRD"

var upgrader = websocket.Upgrader{}

type Worker struct {
	viterbiData     []byte
	decodedData     []byte
	lastFrameEnd    []byte
	codedData       []byte
	rsCorrectedData []byte
	rsWorkBuffer    []byte
	syncWord        []byte
	viterbi         SatHelper.Viterbi27
	reedSolomon     SatHelper.ReedSolomon
	correlator      SatHelper.Correlator
	packetFixer     SatHelper.PacketFixer
	Statistics      Statistics
	constSock       *websocket.Conn
	statsSock       *websocket.Conn
}

func NewDecoder(uuid string) interfaces.Decoder {
	e := Worker{}

	if uuid != "" {
		http.HandleFunc(fmt.Sprintf("/hrd/%s/constellation", uuid), e.constellation)
		http.HandleFunc(fmt.Sprintf("/hrd/%s/statistics", uuid), e.statistics)
	}

	if uselastFrameData {
		e.viterbiData = make([]byte, Datalink[id].CodedFrameSize+lastFrameDataBits)
		e.decodedData = make([]byte, Datalink[id].FrameSize+lastFrameData)
		e.lastFrameEnd = make([]byte, lastFrameDataBits)

		e.viterbi = SatHelper.NewViterbi27(Datalink[id].FrameBits + lastFrameDataBits)

		for i := 0; i < lastFrameDataBits; i++ {
			e.lastFrameEnd[i] = 128
		}
	} else {
		e.viterbiData = make([]byte, Datalink[id].CodedFrameSize)
		e.decodedData = make([]byte, Datalink[id].FrameSize)

		e.viterbi = SatHelper.NewViterbi27(Datalink[id].FrameBits)
	}

	e.codedData = make([]byte, Datalink[id].CodedFrameSize)
	e.rsCorrectedData = make([]byte, Datalink[id].FrameSize)
	e.rsWorkBuffer = make([]byte, 255)

	e.reedSolomon = SatHelper.NewReedSolomon()
	e.correlator = SatHelper.NewCorrelator()
	e.packetFixer = SatHelper.NewPacketFixer()

	e.syncWord = make([]byte, 4)

	e.reedSolomon.SetCopyParityToOutput(true)

	e.correlator.AddWord(Datalink[id].SyncWords[0])
	e.correlator.AddWord(Datalink[id].SyncWords[1])
	e.correlator.AddWord(Datalink[id].SyncWords[2])
	e.correlator.AddWord(Datalink[id].SyncWords[3])

	e.correlator.AddWord(Datalink[id].SyncWords[4])
	e.correlator.AddWord(Datalink[id].SyncWords[5])
	e.correlator.AddWord(Datalink[id].SyncWords[6])
	e.correlator.AddWord(Datalink[id].SyncWords[7])

	return &e
}

func (e *Worker) Work(inputPath string, outputPath string, g *bool) {
	var isCorrupted bool
	lastFrameOk := false

	fmt.Printf("[DEC] Initializing decoding process...\n")

	var averageRSCorrections float32
	var averageVitCorrections float32
	var lostPacketsPerChannel [256]int64
	var lastPacketCount [256]int64
	var receivedPacketsPerChannel [256]int64
	var phaseShift SatHelper.SatHelperPhaseShift
	var flywheelCount = 0

	input, err := os.Open(inputPath)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer input.Close()

	output, err := os.Create(outputPath)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer output.Close()

	fi, _ := os.Stat(inputPath)

	fmt.Printf("[DEC] Starting decoding the signal. This might take a while...\n")
	for *g {
		n, err := input.Read(e.codedData)
		if Datalink[id].CodedFrameSize != n {
			break
		}

		if err == nil {
			e.Statistics.TotalBytesRead += uint64(n)
			if (e.Statistics.TotalPackets%32 == 0) && e.constSock != nil {
				e.constSock.WriteMessage(1, []byte(b64.StdEncoding.EncodeToString(e.codedData[:256])))
			}

			if e.Statistics.TotalPackets%averageLastNSamples == 0 {
				averageRSCorrections = 0
				averageVitCorrections = 0
			}

			if flywheelCount == defaultFlywheelRecheck {
				lastFrameOk = false
				flywheelCount = 0
			}

			if !lastFrameOk {
				e.correlator.Correlate(&e.codedData[0], uint(Datalink[id].CodedFrameSize))
			} else {
				e.correlator.Correlate(&e.codedData[0], uint(Datalink[id].CodedFrameSize)/64)
				if e.correlator.GetHighestCorrelationPosition() != 0 {
					e.correlator.Correlate(&e.codedData[0], uint(Datalink[id].CodedFrameSize))
					flywheelCount = 0
				}
			}
			flywheelCount++

			word := e.correlator.GetCorrelationWordNumber()
			pos := e.correlator.GetHighestCorrelationPosition()
			corr := e.correlator.GetHighestCorrelation()

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

			if corr < Datalink[id].MinCorrelationBits {
				//fmt.Printf("[DEC] Not enough correlations %d/%d. Skipping...\n", corr, Datalink[id].MinCorrelationBits)
				continue
			}

			if pos != 0 {
				shiftWithConstantSize(&e.codedData, int(pos), Datalink[id].CodedFrameSize)
				offset := Datalink[id].CodedFrameSize - int(pos)

				buffer := make([]byte, int(pos))
				n, err = input.Read(buffer)

				e.Statistics.TotalBytesRead += uint64(n)
				if err != nil {
					fmt.Println(err)
					break
				}

				for i := offset; i < Datalink[id].CodedFrameSize; i++ {
					e.codedData[i] = buffer[i-offset]
				}
			}

			e.packetFixer.FixPacket(&e.codedData[0], uint(Datalink[id].CodedFrameSize), phaseShift, iqInv)

			if uselastFrameData {
				for i := 0; i < lastFrameDataBits; i++ {
					e.viterbiData[i] = e.lastFrameEnd[i]
				}
				for i := lastFrameDataBits; i < Datalink[id].CodedFrameSize+lastFrameDataBits; i++ {
					e.viterbiData[i] = e.codedData[i-lastFrameDataBits]
				}
			} else {
				for i := 0; i < Datalink[id].CodedFrameSize; i++ {
					e.viterbiData[i] = e.codedData[i]
				}
			}

			e.viterbi.Decode(&e.viterbiData[0], &e.decodedData[0])

			nrzmDecodeSize := Datalink[id].FrameSize

			if uselastFrameData {
				nrzmDecodeSize += lastFrameData
			}

			SatHelper.DifferentialEncodingNrzmDecode(&e.decodedData[0], nrzmDecodeSize)

			signalErrors := float32(e.viterbi.GetPercentBER())
			signalErrors = 100 - (signalErrors * 10)
			signalQuality := uint8(signalErrors)

			averageVitCorrections += float32(e.viterbi.GetBER())

			if uselastFrameData {
				shiftWithConstantSize(&e.decodedData, lastFrameData/2, Datalink[id].FrameSize+lastFrameData/2)
				for i := 0; i < lastFrameDataBits; i++ {
					e.lastFrameEnd[i] = e.viterbiData[Datalink[id].CodedFrameSize+i]
				}
			}

			for i := 0; i < Datalink[id].SyncWordSize; i++ {
				e.syncWord[i] = e.decodedData[i]
				e.Statistics.SyncWord[i] = e.decodedData[i]
			}

			shiftWithConstantSize(&e.decodedData, Datalink[id].SyncWordSize, Datalink[id].FrameSize-Datalink[id].SyncWordSize)

			e.Statistics.AverageVitCorrections += uint16(e.viterbi.GetBER())
			e.Statistics.TotalPackets++

			SatHelper.DeRandomizerDeRandomize(&e.decodedData[0], Datalink[id].FrameSize-Datalink[id].SyncWordSize)

			derrors := make([]int32, Datalink[id].RsBlocks)

			for i := 0; i < Datalink[id].RsBlocks; i++ {
				e.reedSolomon.Deinterleave(&e.decodedData[0], &e.rsWorkBuffer[0], byte(i), byte(Datalink[id].RsBlocks))
				derrors[i] = int32(int8(e.reedSolomon.Decode_ccsds(&e.rsWorkBuffer[0])))
				e.reedSolomon.Interleave(&e.rsWorkBuffer[0], &e.rsCorrectedData[0], byte(i), byte(Datalink[id].RsBlocks))
				if derrors[i] != -1 {
					averageRSCorrections += float32(derrors[i])
				}
				e.Statistics.RsErrors[i] = derrors[i]
			}

			if derrors[0] == -1 && derrors[1] == -1 && derrors[2] == -1 && derrors[3] == -1 {
				isCorrupted = true
				lastFrameOk = false
				e.Statistics.DroppedPackets++
			} else {
				isCorrupted = false
				lastFrameOk = true
			}

			scid := ((e.rsCorrectedData[0] & 0x3F) << 2) | (e.rsCorrectedData[1]&0xC0)>>6
			vcid := e.rsCorrectedData[1] & 0x3F
			counter := uint(e.rsCorrectedData[2])
			counter = SatHelper.ToolsSwapEndianess(counter)
			counter &= 0xFFFFFF00
			counter >>= 8

			e.Statistics.SCID = scid
			e.Statistics.VCID = vcid

			vitBitErr := e.viterbi.GetBER()

			vitBitErr -= lastFrameDataBits / 2

			if vitBitErr < 0 {
				vitBitErr = 0
			}

			if signalQuality > 100 || isCorrupted {
				signalQuality = 0
			}

			e.Statistics.PacketNumber = uint64(counter)
			e.Statistics.VitErrors = uint16(vitBitErr)
			e.Statistics.FrameBits = uint16(Datalink[id].FrameBits)
			e.Statistics.SignalQuality = signalQuality
			e.Statistics.SyncCorrelation = uint8(corr)
			e.Statistics.TotalBytes = uint64(fi.Size())

			if !isCorrupted {
				if lastPacketCount[vcid]+1 != int64(counter) && lastPacketCount[vcid] > -1 {
					lostCount := int(int64(counter) - lastPacketCount[vcid] - 1)
					e.Statistics.LostPackets += uint64(lostCount)
					lostPacketsPerChannel[vcid] += int64(lostCount)
				}
				lastPacketCount[vcid] = int64(counter)
				if receivedPacketsPerChannel[vcid] == -1 {
					receivedPacketsPerChannel[vcid] = 1
				} else {
					receivedPacketsPerChannel[vcid] = receivedPacketsPerChannel[vcid] + 1
				}

				if e.Statistics.TotalPackets%averageLastNSamples == 0 {
					e.Statistics.AverageRSCorrections = uint8(averageRSCorrections / 4)
					e.Statistics.AverageVitCorrections = uint16(averageVitCorrections)
				} else {
					e.Statistics.AverageRSCorrections = uint8(averageRSCorrections / float32(4*(e.Statistics.TotalPackets%averageLastNSamples)))
					e.Statistics.AverageVitCorrections = uint16(averageVitCorrections / float32(e.Statistics.TotalPackets%averageLastNSamples))
				}
				e.Statistics.FrameLock = 1
				for i := 0; i < 256; i++ {
					e.Statistics.ReceivedPacketsPerChannel[i] = receivedPacketsPerChannel[i]
					e.Statistics.LostPacketsPerChannel[i] = lostPacketsPerChannel[i]
				}

				dat := e.rsCorrectedData[:Datalink[id].FrameSize-Datalink[id].RsParityBlockSize-Datalink[id].SyncWordSize]
				output.Write(dat)
			} else {
				e.Statistics.FrameLock = 0
			}

			if e.Statistics.TotalPackets%32 == 0 && e.statsSock != nil {
				json, err := json.Marshal(e.Statistics)
				if err == nil {
					e.statsSock.WriteMessage(1, []byte(json))
				}
			}

			if e.Statistics.TotalPackets%512 == 0 {
				fmt.Printf("[DEC] (%04.1f%%) Dropped Packets: %4.1f%%; Signal Quality: %2d%%\n",
					float32(e.Statistics.TotalBytesRead)/float32(fi.Size())*100,
					float32(e.Statistics.DroppedPackets)/float32(e.Statistics.TotalPackets)*100,
					e.Statistics.SignalQuality)
			}

		} else {
			if err != io.EOF {
				fmt.Println(err)
			}
			break
		}
	}

	fmt.Printf("[DEC] Decoding finished! File saved as %s\n", outputPath)
}

func (e *Worker) constellation(w http.ResponseWriter, r *http.Request) {
	upgrader.CheckOrigin = func(r *http.Request) bool { return true }
	e.constSock, _ = upgrader.Upgrade(w, r, nil)
}

func (e *Worker) statistics(w http.ResponseWriter, r *http.Request) {
	upgrader.CheckOrigin = func(r *http.Request) bool { return true }
	e.statsSock, _ = upgrader.Upgrade(w, r, nil)
}

func shiftWithConstantSize(arr *[]byte, pos int, length int) {
	for i := 0; i < length-pos; i++ {
		(*arr)[i] = (*arr)[pos+i]
	}
}
