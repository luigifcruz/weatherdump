package jpeg

// EOB indicates the End Of Block of each MCU.
var EOB = []int64{-99999}

// CFC indicates that no match was found inside the Huffman LUT.
var CFC = []int64{-99998}

// GetQuantizationTable returns the standard quantization table
// with the quality factor correction.
func GetQuantizationTable(qf float64) []int64 {
	var table [64]int64

	if (qf > 20) && (qf < 50) {
		qf = 5000 / qf
	} else {
		qf = 200 - (2 * qf)
	}

	for x := 0; x < 64; x++ {
		table[x] = int64((qf / 100 * qTable[x]) + 0.5)
		if table[x] < 1 {
			table[x] = 1
		}
	}
	return table[:]
}

// FindDC decodes and return the next DC coefficient by
// applying Huffman to the bool slice received.
func FindDC(dat *[]bool) int64 {
	buf := *dat
	bufl := len(*dat)
	for _, m := range dcCategories {
		klen := len(m.code)
		if bufl < klen {
			continue
		}

		if fastEqual(buf[:klen], m.code) {
			if bufl < klen+m.len {
				break
			}
			*dat = buf[klen+m.len:]
			if m.len == 0 {
				return 0
			}
			return getValue(buf[klen : klen+m.len])
		}
	}

	*dat = nil
	return CFC[0]
}

// FindAC decodes and return the AC coefficient by applying Huffman.
func FindAC(dat *[]bool) []int64 {
	bufl := len(*dat)
	for _, m := range acCategories {
		klen := len(m.code)
		if bufl < klen {
			continue
		}

		if fastEqual((*dat)[:klen], m.code) {
			if m.clen == 0 && m.zlen == 0 {
				*dat = (*dat)[klen:]
				return EOB
			}
			vals := make([]int64, m.zlen+1)
			if !(m.zlen == 15 && m.clen == 0) {
				if bufl < klen+m.clen {
					break
				}
				vals[m.zlen] = getValue((*dat)[klen : klen+m.clen])
			}
			*dat = (*dat)[klen+m.clen:]
			return vals
		}
	}

	*dat = nil
	return CFC
}

// ConvertToArray receives the byte slice and convert
// each bit to a boolean slice that will be returned
// as a pointer.
func ConvertToArray(buf []byte, len int) *[]bool {
	var soft = make([]bool, len*8)
	for i := 0; i < len; i++ {
		soft[0+8*i] = buf[i]>>7&0x01 == 0x01
		soft[1+8*i] = buf[i]>>6&0x01 == 0x01
		soft[2+8*i] = buf[i]>>5&0x01 == 0x01
		soft[3+8*i] = buf[i]>>4&0x01 == 0x01
		soft[4+8*i] = buf[i]>>3&0x01 == 0x01
		soft[5+8*i] = buf[i]>>2&0x01 == 0x01
		soft[6+8*i] = buf[i]>>1&0x01 == 0x01
		soft[7+8*i] = buf[i]>>0&0x01 == 0x01
	}
	return &soft
}

func getValue(dat []bool) int64 {
	var result int64
	for i := 0; i < len(dat); i++ {
		if dat[i] {
			result = result | 0x01<<uint(len(dat)-1-i)
		}
	}
	if !dat[0] {
		result -= (1 << uint(len(dat))) - 1
	}
	return result
}

func fastEqual(a, b []bool) bool {
	for i := 0; i < len(b); i++ {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}
