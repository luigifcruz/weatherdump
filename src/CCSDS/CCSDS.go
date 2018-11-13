package CCSDS

import (
	"fmt"
	"weather-dump/src/CCSDS/Frames"
)

const invalidAPID = uint16(65535)

type CCSDS struct {
	lastAPID uint16
	spacePackets []Frames.SpacePacketFrame
	pendingPackets [2047]Frames.SpacePacketFrame
	buffer []byte
}

func (e CCSDS) GetSpacePackets() []Frames.SpacePacketFrame {
	return e.spacePackets
}

func CreatePacket(e *CCSDS) (uint16, []byte) {
	dat := e.buffer
	apid := invalidAPID

	for {
		if len(dat) < 6 {
			return apid, dat
		}

		s := Frames.SpacePacketFrame{}
		s.FromBinary(dat)
		apid = s.GetAPID()
		dat = dat[6:]

		if s.GetAPID() != 2047 {
			e.pendingPackets[apid] = s
		} else {
			apid = invalidAPID
		}

		if s.GetPacketLength() + 1 < uint16(len(dat)) {
			if apid != invalidAPID {
				e.spacePackets = append(e.spacePackets, e.pendingPackets[apid])
				e.pendingPackets[apid] = Frames.SpacePacketFrame{}
			}

			dat = dat[s.GetPacketLength()+1:]
			apid = invalidAPID
		} else {
			break
		}
	}

	return apid, []byte{}
}

func ParseMPDU(e *CCSDS, MPDU Frames.MultiplexingFrame) {

		dat := MPDU.GetPacketZone()
		fhp := MPDU.GetFirstHeaderPointer()

		if MPDU.HaveNewPackage() {
			if e.lastAPID == invalidAPID && len(e.buffer) > 0 {
				if fhp > 0 {
					e.buffer = append(e.buffer, dat[:fhp]...)
				}

				e.lastAPID, dat = CreatePacket(e)
				if e.lastAPID == invalidAPID {
					e.buffer = dat
				} else {
					e.buffer = []byte{}
				}
			}

			// Finishing another Space Packet!
			if e.lastAPID != invalidAPID && len(dat) > 0 {
				if e.lastAPID > 0 {
					e.pendingPackets[e.lastAPID].FeedData(dat[:fhp])
				}
				e.spacePackets = append(e.spacePackets, e.pendingPackets[e.lastAPID])
				e.pendingPackets[e.lastAPID] = Frames.SpacePacketFrame{}
				e.lastAPID = invalidAPID
			}

			// Try to create a new packet!
			if len(dat) > int(fhp) {
				e.buffer = append(e.buffer, dat[fhp:]...)
				e.lastAPID, dat = CreatePacket(e)
				if e.lastAPID == invalidAPID {
					e.buffer = dat
				} else {
					e.buffer = []byte{}
				}
			}
		} else {
			if len(e.buffer) > 0 && e.lastAPID == invalidAPID {
				e.buffer = append(e.buffer, dat...)
				e.lastAPID, dat = CreatePacket(e)
				if e.lastAPID == invalidAPID {
					e.buffer = dat
				} else {
					e.buffer = []byte{}
				}
			} else if len(e.buffer) > 0 {
				fmt.Errorf("problem with continuation package")
			} else if e.lastAPID == invalidAPID {
				e.buffer = append(e.buffer, dat...)
				e.lastAPID, dat = CreatePacket(e)
				if e.lastAPID == invalidAPID {
					e.buffer = dat
				} else {
					e.buffer = []byte{}
				}
			} else {
				e.pendingPackets[e.lastAPID].FeedData(dat)
			}
		}
}