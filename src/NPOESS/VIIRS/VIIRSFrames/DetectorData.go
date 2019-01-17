package VIIRSFrames

import (
	"encoding/binary"
	"fmt"
	"math"
	"unsafe"
)

/*
#include <stdlib.h>
#include <stdint.h>
#include <libaec.h>
#include <string.h>
#cgo LDFLAGS: -laec

void decompress(char *input, char *output, int inputLen, int outputLen) {
	struct aec_stream strm;

	strm.bits_per_sample = 15;
	strm.block_size = 8;
	strm.rsi = 128;
	strm.flags = AEC_DATA_MSB | AEC_DATA_PREPROCESS;
	strm.next_in = input;
	strm.avail_in = inputLen;
	strm.next_out = output;
	strm.avail_out = outputLen * sizeof(char);

	aec_decode_init(&strm);
	aec_decode(&strm, AEC_FLUSH);
	aec_decode_end(&strm);
}
*/
import "C"

func Decompress(data []byte, inputLen int, outputLen int) []byte {
	var slice = make([]byte, outputLen)
	C.decompress((*C.char)(unsafe.Pointer(&data[0])), (*C.char)(unsafe.Pointer(&slice[0])), C.int(inputLen), C.int(outputLen))
	return slice
}

type DetectorData struct {
	fillData       uint8
	checksumOffset uint16
	aggregator     []byte
	checksum       uint32
	syncWord       uint32
	diffBuf        []byte
}

func NewDetectorData() *DetectorData {
	return &DetectorData{}
}

func (e *DetectorData) FromBinary(buf *[]byte) {
	dat := *buf

	if len(dat) < 4 {
		return
	}

	e.fillData = uint8(dat[0])
	e.checksumOffset = binary.BigEndian.Uint16(dat[2:])

	cso := e.checksumOffset

	if int(cso) >= len(dat)+4 || cso == 0 {
		return
	}

	e.aggregator = dat[4:cso]
	bitSlicer(&e.aggregator, int(e.fillData))

	if (len(dat) - int(cso)) > 8 {
		e.checksum = binary.BigEndian.Uint32(dat[cso:])
		e.syncWord = binary.BigEndian.Uint32(dat[cso+4:])
		*buf = (*buf)[cso+8:]
	} else {
		e.syncWord = 0xC000FFEE
	}
}

func (e DetectorData) Print() {
	fmt.Println("### VIIRS Aggregator")
	fmt.Printf("Fill Data: %08b\n", e.fillData)
	fmt.Printf("Checksum Offset: %016b\n", e.checksumOffset)
	fmt.Printf("Data Size: %d\n", len(e.aggregator))
	fmt.Printf("Checksum: %032b\n", e.checksum)
	fmt.Printf("Sync Word: %032b\n", e.syncWord)
	fmt.Println()
}

// Struct Validation
// Struct Get
func (e DetectorData) GetChecksum() uint16 {
	return e.checksumOffset
}

func (e DetectorData) GetData(syncwork uint32, width int, oversample int) []byte {
	if len(e.diffBuf) > 0 {
		return e.diffBuf
	}

	if len(e.aggregator) > 8 && (syncwork == e.syncWord || e.syncWord == 0xC000FFEE) {
		var buf []byte
		size := width * 2 * oversample // 16-bits pixels * oversample
		dat := Decompress(e.aggregator, len(e.aggregator), size)

		if oversample == 1 {
			return dat
		}

		for x := 0; x < size; x += oversample * 2 {
			var val uint16

			if oversample > 1 {
				val += binary.BigEndian.Uint16(dat[x : x+2])
				val += binary.BigEndian.Uint16(dat[x+2 : x+4])
			}

			if oversample > 2 {
				val += binary.BigEndian.Uint16(dat[x+4 : x+6])
			}

			val /= uint16(oversample)

			b := make([]byte, 2)
			binary.BigEndian.PutUint16(b, val)
			buf = append(buf, b...)
		}

		return buf
	}

	return make([]byte, width*2)
}

// Struct Set
func (e *DetectorData) SetData(dat *[]byte) {
	e.diffBuf = make([]byte, len(*dat))
	copy(e.diffBuf, *dat)
}

// Struct Tools
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

	if len(*dat) < len(buf) {
		return
	}

	buf[len(*dat)-1] = uint8(buf[len(*dat)-1]) & ^(uint8(math.Pow(2, float64(bits))) - 1)
}
