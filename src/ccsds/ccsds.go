package ccsds

import (
	"weather-dump/src/ccsds/frames"
	"weather-dump/src/ccsds/parameters"
)

var Version = parameters.Version

type Worker struct {
	spacePackets []frames.SpacePacketFrame
	tmpPacket    *frames.SpacePacketFrame
	buffer       []byte
}

func New() *Worker {
	return &Worker{}
}

func (e Worker) GetSpacePackets() []frames.SpacePacketFrame {
	return e.spacePackets
}

func (e *Worker) CloseFrame() {
	if e.tmpPacket != nil {
		e.spacePackets = append(e.spacePackets, *e.tmpPacket)
	}

	e.buffer = make([]byte, 0)
	e.tmpPacket = nil
}

func (e *Worker) CreatePacket(buf []byte) {
	if len(buf) == 0 {
		return
	}

	if e.tmpPacket != nil {
		e.CloseFrame()
	}

	if len(buf) > 6 {
		e.tmpPacket = &frames.SpacePacketFrame{}
		e.tmpPacket.FromBinary(buf[:6])
		buf = buf[6:]
	} else {
		e.buffer = buf
	}

	if e.tmpPacket != nil {
		buf = e.tmpPacket.FeedData(buf)

		if e.tmpPacket.IsValid() {
			e.CloseFrame()
		}

		if len(buf) > 0 {
			e.CreatePacket(buf)
		}
	}
}

func (e *Worker) ParseMPDU(MPDU frames.MultiplexingFrame) {
	if !MPDU.IsValid() {
		//fmt.Println("[CCSDS] Not Valid MPDU frame, skipping...")
		return
	}

	dat := MPDU.GetPacketZone()
	fhp := MPDU.GetFirstHeaderPointer()

	if MPDU.HaveNewPackage() && fhp > uint16(len(dat)) {
		//fmt.Println("[CCSDS] First header pointer bigger than buffer, skipping...")
		return
	}

	if MPDU.HaveNewPackage() {
		if e.tmpPacket == nil && len(e.buffer) > 0 {
			buf := append(e.buffer, dat[:fhp]...)
			e.CreatePacket(buf)
		} else if e.tmpPacket != nil {
			e.tmpPacket.FeedData(dat[:fhp])
			e.CloseFrame()
		}

		e.CreatePacket(dat[fhp:])
	} else {
		if len(e.buffer) > 0 && e.tmpPacket == nil {
			buf := append(e.buffer, dat...)
			e.CreatePacket(buf)
		} else if e.tmpPacket == nil {
			// IGNORE
		} else {
			e.tmpPacket.FeedData(dat)
		}
	}

	if e.tmpPacket != nil && e.tmpPacket.IsValid() {
		e.CloseFrame()
	}
}
