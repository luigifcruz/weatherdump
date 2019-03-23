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

	SatHelper "github.com/OpenSatelliteProject/libsathelper"
	"github.com/fatih/color"
	"github.com/gorilla/websocket"
	"github.com/gosuri/uiprogress"
)

type CaduDecoder struct {
	hardData     []byte
	softData     []byte
	rsWorkBuffer []byte
	correlator   SatHelper.Correlator
	reedSolomon  SatHelper.ReedSolomon
	Statistics   assets.Statistics
	constSock    *websocket.Conn
	statsSock    *websocket.Conn
}

func NewCaduDecoder(uuid string) interfaces.Decoder {
	e := CaduDecoder{}

	if uuid != "" {
		http.HandleFunc(fmt.Sprintf("/hrd/%s/constellation", uuid), e.constellation)
		http.HandleFunc(fmt.Sprintf("/hrd/%s/statistics", uuid), e.statistics)
	}

	e.softData = make([]byte, datalink[id].FrameBits)
	e.hardData = make([]byte, datalink[id].FrameSize)
	e.rsWorkBuffer = make([]byte, 255)
	e.correlator = SatHelper.NewCorrelator()
	e.reedSolomon = SatHelper.NewReedSolomon()
	e.reedSolomon.SetCopyParityToOutput(true)
	e.Statistics.DroppedPackets = 0
	e.Statistics.TotalPackets = 1

	e.correlator.AddWord(uint(0x1ACFFC1D))
	e.correlator.AddWord(uint(0xE53003E2))

	return &e
}

func (e *CaduDecoder) Work(inputPath string, outputPath string, g *bool) {
	color.Red("[DEC] WARNING! This decoder is currently in ALPHA development state.")
	flywheelCount := 0

	fi, err := os.Stat(inputPath)
	input, err := os.Open(inputPath)
	outputBuf, err := os.Create(outputPath + ".buf")
	if err != nil {
		log.Fatal(err)
	}

	e.Statistics.TotalBytes = uint64(fi.Size())
	e.Statistics.TaskName = "Converting CADU file"

	progress := uiprogress.New()
	progress.Start()

	bar1 := progress.AddBar(int(fi.Size())).AppendCompleted()
	bar2 := progress.AddBar(int(fi.Size()) * 8).AppendCompleted()

	bar1.PrependFunc(func(b *uiprogress.Bar) string {
		return "[DEC] Converting CADU file	"
	})

	bar2.PrependFunc(func(b *uiprogress.Bar) string {
		return "[DEC] Decoding soft-symbol file	"
	})

	bar2.AppendFunc(func(b *uiprogress.Bar) string {
		s := e.Statistics
		return fmt.Sprintf("\n[DEC] Decoder Statistics	 [VCID: %2d] [#%8d]  [CORR: %2d] [RS: %2d %2d %2d %2d]  [DROPPED: %4.1f%%]",
			s.VCID, s.PacketNumber, s.SyncCorrelation,
			s.AverageRSCorrections[0], s.AverageRSCorrections[1], s.AverageRSCorrections[2], s.AverageRSCorrections[3],
			float32(s.DroppedPackets)/float32(s.TotalPackets)*100)
	})

	for *g {
		n, err := input.Read(e.hardData)
		if datalink[id].FrameSize != n {
			break
		}

		if err == nil {
			e.Statistics.TotalBytesRead += uint64(n)
			bar1.Set(int(e.Statistics.TotalBytesRead))

			convertToArray(e.hardData, &e.softData, datalink[id].FrameSize)
			outputBuf.Write(e.softData)

			if e.Statistics.TotalBytesRead%1e4 == 0 && e.statsSock != nil {
				e.updateStatistics(e.Statistics)
			}
		} else {
			if err != io.EOF {
				log.Fatal(err)
			}
			break
		}
	}

	input.Close()
	outputBuf.Close()

	fi, err = os.Stat(outputPath + ".buf")
	output, err := os.Create(outputPath)
	inputBuf, err := os.Open(outputPath + ".buf")
	if err != nil {
		log.Fatal(err)
	}

	e.Statistics.TotalBytesRead = 0
	e.Statistics.TotalBytes = uint64(fi.Size())
	e.Statistics.TaskName = "Decoding soft-symbol file"

	for *g {
		n, err := inputBuf.Read(e.softData)
		if datalink[id].FrameBits != n {
			break
		}

		if err == nil {
			e.Statistics.TotalBytesRead += uint64(n)
			bar2.Set(int(e.Statistics.TotalBytesRead))

			if e.Statistics.TotalPackets%averageLastNSamples == 0 {
				e.Statistics.AverageRSCorrections = [4]int{}
			}

			if flywheelCount == defaultFlywheelRecheck*8 {
				e.Statistics.FrameLock = false
				flywheelCount = 0
			}

			if !e.Statistics.FrameLock {
				e.correlator.Correlate(&e.softData[0], uint(datalink[id].FrameBits))
			} else {
				e.correlator.Correlate(&e.softData[0], uint(datalink[id].FrameBits)/128)
				if e.correlator.GetHighestCorrelationPosition() != 0 {
					e.correlator.Correlate(&e.softData[0], uint(datalink[id].FrameBits))
					flywheelCount = 0
				}
			}
			flywheelCount++

			pos := e.correlator.GetHighestCorrelationPosition()
			cor := e.correlator.GetHighestCorrelation()

			if cor < datalink[id].MinCorrelationBits/2 {
				//fmt.Printf("[DEC] Not enough correlations %d/%d. Skipping...\n", corr, Datalink[id].MinCorrelationBits)
				continue
			}

			if pos != 0 {
				shiftWithConstantSize(&e.softData, int(pos), datalink[id].FrameBits)
				offset := datalink[id].FrameBits - int(pos)

				buffer := make([]byte, int(pos))
				n, err = inputBuf.Read(buffer)

				e.Statistics.TotalBytesRead += uint64(n)
				bar2.Set(int(e.Statistics.TotalBytesRead))
				if err != nil {
					fmt.Println(err)
					break
				}

				for i := offset; i < datalink[id].FrameBits; i++ {
					e.softData[i] = buffer[i-offset]
				}
			}

			for i := 0; i < datalink[id].FrameBits; i += 8 {
				b := byte(0x00)
				for j := i; j < i+8 && j < datalink[id].FrameBits; j++ {
					v := byte(0x00)
					if e.softData[j] > 128 {
						v = byte(0x01)
					}
					b = (b << 1) | v
				}
				e.hardData[i/8] = b
			}

			shiftWithConstantSize(&e.hardData, datalink[id].SyncWordSize, datalink[id].FrameSize-datalink[id].SyncWordSize)
			SatHelper.DeRandomizerDeRandomize(&e.hardData[0], datalink[id].FrameSize-datalink[id].SyncWordSize)
			e.Statistics.TotalPackets++

			var derrors [4]int
			for i := 0; i < datalink[id].RsBlocks; i++ {
				e.reedSolomon.Deinterleave(&e.hardData[0], &e.rsWorkBuffer[0], byte(i), byte(datalink[id].RsBlocks))
				derrors[i] = int(int8(e.reedSolomon.Decode_ccsds(&e.rsWorkBuffer[0])))
				e.reedSolomon.Interleave(&e.rsWorkBuffer[0], &e.hardData[0], byte(i), byte(datalink[id].RsBlocks))
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

			e.Statistics.SyncCorrelation = uint8(cor)
			e.Statistics.VCID = e.hardData[1] & 0x3F
			e.Statistics.FrameBits = uint16(datalink[id].FrameBits)
			e.Statistics.PacketNumber = binary.BigEndian.Uint32(e.hardData[2:]) & 0xFFFFFF00 >> 8

			if e.Statistics.FrameLock {
				e.Statistics.ReceivedPacketsPerChannel[e.Statistics.VCID]++
				dat := e.hardData[:datalink[id].FrameSize-datalink[id].RsParityBlockSize-datalink[id].SyncWordSize]
				output.Write(dat)
			}

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

	output.Close()
	inputBuf.Close()
	progress.Stop()
	os.Remove(outputPath + ".buf")

	if e.statsSock != nil {
		e.Statistics.Finish()
		e.updateStatistics(e.Statistics)
	}

	color.Green("[DEC] Decoding finished! File saved in the same folder.\n")
}

func convertToArray(hard []byte, soft *[]byte, len int) {
	var buf = make([]bool, len*8)
	for i := 0; i < len; i++ {
		buf[0+8*i] = hard[i]>>7&0x01 == 0x01
		buf[1+8*i] = hard[i]>>6&0x01 == 0x01
		buf[2+8*i] = hard[i]>>5&0x01 == 0x01
		buf[3+8*i] = hard[i]>>4&0x01 == 0x01
		buf[4+8*i] = hard[i]>>3&0x01 == 0x01
		buf[5+8*i] = hard[i]>>2&0x01 == 0x01
		buf[6+8*i] = hard[i]>>1&0x01 == 0x01
		buf[7+8*i] = hard[i]>>0&0x01 == 0x01
	}
	for i := 0; i < len*8; i++ {
		if buf[i] == true {
			(*soft)[i] = 0xFF
		} else {
			(*soft)[i] = 0x00
		}
	}
}

func (e *CaduDecoder) updateStatistics(s assets.Statistics) {
	json, err := json.Marshal(s)
	if err == nil {
		e.statsSock.WriteMessage(1, []byte(json))
	}
}

func (e *CaduDecoder) constellation(w http.ResponseWriter, r *http.Request) {
	upgrader.CheckOrigin = func(r *http.Request) bool { return true }
	e.constSock, _ = upgrader.Upgrade(w, r, nil)
}

func (e *CaduDecoder) statistics(w http.ResponseWriter, r *http.Request) {
	upgrader.CheckOrigin = func(r *http.Request) bool { return true }
	e.statsSock, _ = upgrader.Upgrade(w, r, nil)
}

func shiftWithConstantSize(arr *[]byte, pos int, length int) {
	for i := 0; i < length-pos; i++ {
		(*arr)[i] = (*arr)[pos+i]
	}
}
