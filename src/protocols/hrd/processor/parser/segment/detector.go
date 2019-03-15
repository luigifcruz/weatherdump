package segment

import (
	"encoding/binary"
	"fmt"
	"math"
)

const detectorMinimum = 88

// Detector is the final data structure from the VIIRS pictures products.
type Detector struct {
	fillData       uint8
	checksumOffset uint16
	checksum       uint32
	syncWord       uint32
	data           []byte
}

// NewDetector returns a pointer of a new Detector.
func NewDetector() *Detector {
	return &Detector{}
}

// FromBinary parses the binary data into the dectector struct.
func (e *Detector) FromBinary(buf *[]byte) {
	if len(*buf) < detectorMinimum {
		return
	}

	dat := *buf
	e.fillData = uint8(dat[0])
	e.checksumOffset = binary.BigEndian.Uint16(dat[2:])

	cso := e.checksumOffset

	if int(cso) >= len(dat)+4 || cso < 4 {
		return
	}

	e.data = dat[4:cso]
	bitSlicer(&e.data, int(e.fillData))

	if (len(dat) - int(cso)) > 8 {
		e.checksum = binary.BigEndian.Uint32(dat[cso:])
		e.syncWord = binary.BigEndian.Uint32(dat[cso+4:])
		*buf = (*buf)[cso+8:]
	} else {
		e.syncWord = 0xC000FFEE
	}
}

// Print the values contained inside the detector struct into stdout.
func (e Detector) Print() {
	fmt.Println("### VIIRS Aggregator")
	fmt.Printf("Fill Data: %08b\n", e.fillData)
	fmt.Printf("Checksum Offset: %016b\n", e.checksumOffset)
	fmt.Printf("Data Size: %d\n", len(e.data))
	fmt.Printf("Checksum: %032b\n", e.checksum)
	fmt.Printf("Sync Word: %032b\n", e.syncWord)
	fmt.Println()
}

// GetChecksum value from the detector struct.
func (e Detector) GetChecksum() uint16 {
	return e.checksumOffset
}

func (e Detector) IsValid(syncWord uint32) bool {
	return len(e.data) > 8 && (syncWord == e.syncWord || e.syncWord == 0xC000FFEE)
}

// GetData from the current detector.
func (e Detector) GetData() *[]byte {
	return &e.data
}

func (e *Detector) Pad(width int) {
	e.data = make([]byte, width*2)
}

func (e *Detector) Decompress(width, oversample int) {
	e.data, _ = Decompress(e.data, len(e.data), width*2*oversample)
}

func (e *Detector) Decimate(width, oversample int) {
	if oversample == 1 {
		return
	}

	for x := 0; x < len(e.data); x += oversample * 2 {
		var val uint16

		switch oversample {
		case 2:
			val += binary.BigEndian.Uint16(e.data[x:])
			val += binary.BigEndian.Uint16(e.data[x+2:])
		case 3:
			val += binary.BigEndian.Uint16(e.data[x:])
			val += binary.BigEndian.Uint16(e.data[x+2:])
			val += binary.BigEndian.Uint16(e.data[x+4:])
		}

		val /= uint16(oversample)
		binary.BigEndian.PutUint16(e.data[x/oversample:], val)
	}

	e.data = append([]byte(nil), e.data[:width*2]...)
}

// SetData updates the data inside the detector.
func (e *Detector) Integrate(diff *[]byte, decimation int) {
	for i := 0; i < len(e.data); i += 2 {
		var base, differential uint16
		base = binary.BigEndian.Uint16(e.data[i:])
		if (i/decimation/2*2)+1 < len(*diff) {
			differential = binary.BigEndian.Uint16((*diff)[i/decimation/2*2:])
		}
		binary.BigEndian.PutUint16(e.data[i:], base+differential-16383)
	}
}

func bitSlicer(dat *[]byte, size int) {
	buf := *dat
	bits, bytes := 0, 0

	for size%8 != 0 {
		bits++
		size--
	}

	bytes = len(*dat) - (size / 8)

	if bytes > len(*dat) || bytes < 0 {
		return
	}

	*dat = (*dat)[:bytes]

	if len(*dat)-1 < len(buf) {
		return
	}

	buf[len(*dat)-1] = uint8(buf[len(*dat)-1]) & ^(uint8(math.Pow(2, float64(bits))) - 1)
}
