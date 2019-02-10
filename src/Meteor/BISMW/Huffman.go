package BISMW

import (
	"reflect"
)

var eob = []int{-999}
var cfc = []int{-998}

func getValue(dat []bool) int {
	if len(dat) == 0 {
		return 0
	}

	var result int
	for i := 1; i < len(dat); i++ {
		if dat[i] {
			result = result | 0x01<<uint(len(dat)-1-i)
		}
	}
	result += 0x01 << uint(len(dat)-1)
	if !dat[0] {
		result *= -1
	}

	return result
}

func findDC(dat *[]bool) int {
	buf := *dat
	for _, m := range dcCategories {
		klen := len(m.code)
		if len(buf) < klen {
			continue
		}

		if reflect.DeepEqual(buf[:klen], m.code) {
			if len(buf) < klen+m.len {
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

func findAC(dat *[]bool) []int {
	buf := *dat
	for _, m := range acCategories {
		klen := len(m.code)
		if len(buf) < klen {
			continue
		}

		if reflect.DeepEqual(buf[:klen], m.code) {
			if m.clen == 0 && m.zlen == 0 {
				*dat = buf[klen:]
				return eob
			}
			vals := make([]int, m.zlen+1)
			if !(m.zlen == 15 && m.clen == 0) {
				if len(buf) < klen+m.clen {
					break
				}
				vals[m.zlen] = getValue(buf[klen : klen+m.clen])
			}
			*dat = buf[klen+m.clen:]
			return vals
		}
	}

	*dat = nil
	return cfc
}

func convertToArray(buf []byte) *[]bool {
	var soft = make([]bool, len(buf)*8)
	for i, m := range buf {
		soft[0+8*i] = m>>7&0x01 == 0x01
		soft[1+8*i] = m>>6&0x01 == 0x01
		soft[2+8*i] = m>>5&0x01 == 0x01
		soft[3+8*i] = m>>4&0x01 == 0x01
		soft[4+8*i] = m>>3&0x01 == 0x01
		soft[5+8*i] = m>>2&0x01 == 0x01
		soft[6+8*i] = m>>1&0x01 == 0x01
		soft[7+8*i] = m>>0&0x01 == 0x01
	}
	return &soft
}
