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

	Recon bool
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

// GetData from the current detector.
func (e Detector) GetData(syncWord uint32, width int, oversample int) []byte {
	if e.Recon {
		return e.data
	}

	if len(e.data) < 8 && (syncWord != e.syncWord || e.syncWord != 0xC000FFEE) {
		return make([]byte, width*2)
	}

	size := width * 2 * oversample
	dat, _ := Decompress(e.data, len(e.data), size)

	if oversample == 1 {
		return dat
	}

	buf := make([]byte, width*2)
	for x := 0; x < size; x += oversample * 2 {
		var val uint16

		switch oversample {
		case 2:
			val += binary.BigEndian.Uint16(dat[x:])
			val += binary.BigEndian.Uint16(dat[x+2:])
		case 3:
			val += binary.BigEndian.Uint16(dat[x:])
			val += binary.BigEndian.Uint16(dat[x+2:])
			val += binary.BigEndian.Uint16(dat[x+4:])
		}

		val /= uint16(oversample)
		binary.BigEndian.PutUint16(buf[x/oversample:], val)
	}

	return buf
}

// SetData updates the data inside the detector.
func (e *Detector) SetData(dat []byte) {
	e.data = make([]byte, len(dat))
	copy(e.data, dat)
	e.Recon = true
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
