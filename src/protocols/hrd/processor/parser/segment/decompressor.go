package segment

/*
#include <stdlib.h>
#include <stdint.h>
#include <libaec.h>
#include <string.h>
#cgo LDFLAGS: -laec

int decompress(char *input, char *output, int inputLen, int outputLen) {
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

	if (aec_decode_init(&strm) != AEC_OK)
        return 1;
    if (aec_decode(&strm, AEC_FLUSH) != AEC_OK)
		return 1;

	aec_decode_end(&strm);
	return 0;
}
*/
import "C"
import "unsafe"

// Decompress the RICE encoded data into the final 16-bits pixel slice.
func Decompress(data []byte, inputLen int, outputLen int) ([]byte, int) {
	var slice = make([]byte, outputLen)
	err := C.decompress((*C.char)(unsafe.Pointer(&data[0])), (*C.char)(unsafe.Pointer(&slice[0])), C.int(inputLen), C.int(outputLen))
	return slice, int(err)
}
