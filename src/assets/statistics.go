package assets

// Statistics struct for the decoder.
type Statistics struct {
	VCID                      uint8
	PacketNumber              uint32
	FrameBits                 uint16
	SignalQuality             uint8
	SyncCorrelation           uint8
	LostPackets               uint64
	AverageVitCorrections     uint16
	AverageRSCorrections      [4]uint8
	DroppedPackets            uint64
	ReceivedPacketsPerChannel [256]int64
	LostPacketsPerChannel     [256]int64
	TotalPackets              uint64
	TotalBytesRead            uint64
	TotalBytes                uint64
	SyncWord                  [4]uint8
	FrameLock                 bool
	Finished                  bool
	TaskName                  string
}

// Finish decoding process.
func (e *Statistics) Finish() {
	e.Finished = true
	e.TotalBytesRead = e.TotalBytes
}
