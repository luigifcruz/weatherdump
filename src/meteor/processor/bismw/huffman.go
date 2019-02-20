package bismw

var eob = []float64{-99999}
var cfc = []float64{-99998}

func getValue(dat []bool) float64 {
	var result int
	for i := 0; i < len(dat); i++ {
		if dat[i] {
			result = result | 0x01<<uint(len(dat)-1-i)
		}
	}
	if !dat[0] {
		result -= (1 << uint(len(dat))) - 1
	}
	return float64(result)
}

func fastEqual(a, b []bool) bool {
	for i := 0; i < len(b); i++ {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

func findDC(dat *[]bool) float64 {
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
	return cfc[0]
}

func findAC(dat *[]bool) []float64 {
	bufl := len(*dat)
	for _, m := range acCategories {
		klen := len(m.code)
		if bufl < klen {
			continue
		}

		if fastEqual((*dat)[:klen], m.code) {
			if m.clen == 0 && m.zlen == 0 {
				*dat = (*dat)[klen:]
				return eob
			}
			vals := make([]float64, m.zlen+1)
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
	return cfc
}

func convertToArray(buf []byte) *[]bool {
	var soft = make([]bool, len(buf)*8)
	for i := 0; i < len(buf); i++ {
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
