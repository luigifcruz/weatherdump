package main

import (
	"fmt"
	"io/ioutil"
	"osp-noaa-dump/VIIRS/ScienceData"
	"osp-noaa-dump/CCSDS"
	"osp-noaa-dump/CCSDS/Frames"
)

const frameSize = 892

func main() {
	fmt.Println("Satellite Helper App - NPOESS Edition")

	file, _ := ioutil.ReadFile("./npp_q.raw")

	d := CCSDS.CCSDS{}
	t := VIIRS.ScienceData{}
	scid := uint8(0)

	bytesCount := 0
	bytesNumber := len(file)

	fmt.Println("Decoding packets...")

	for bytesCount < bytesNumber {
		s := Frames.TransferFrame{}
		s.FromBinary(file[bytesCount:])
		scid = s.GetSCID()

		if s.GetVCID() == 16 {
			p := Frames.MultiplexingFrame{}
			p.FromBinary(s.GetMPDU())

			CCSDS.ParseMPDU(&d, p)
		}

		bytesCount += frameSize
	}

	for _, packet := range d.GetSpacePackets() {
		if packet.GetAPID() >= 800 && packet.GetAPID() <= 823 {
			t.Parse(packet)
		}
	}

	t.SaveAllChannels(scid)

	fmt.Println("Done! Products saved.")
}