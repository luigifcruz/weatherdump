package NPOESS

import (
	"fmt"
	"weather-dump/src/CCSDS/Frames"
)

type Telemetry struct {
	data []byte
}

func (e Telemetry) Parse(packet Frames.SpacePacketFrame) {
	t := Time{}
	t.FromBinary(packet.GetData())
	t.Print()
	fmt.Printf("Remain Data: %d\n", len(packet.GetData()[8:]))
	fmt.Printf("%b\n", packet.GetData()[8:])
}
