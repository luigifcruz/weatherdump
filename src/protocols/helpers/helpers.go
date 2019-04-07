package helpers

func ShiftWithConstantSize(arr *[]byte, pos int, length int) {
	for i := 0; i < length-pos; i++ {
		(*arr)[i] = (*arr)[pos+i]
	}
}

func WatchFor(signal chan bool, method func() bool) {
	for {
		select {
		case <-signal:
			return
		default:
			if method() {
				return
			}
		}
	}
}

func MaxIntSlice(v []int) int {
	index := 0
	max := 0
	for i, e := range v {
		if e > max {
			index = i
			max = e
		}
	}
	return index
}
