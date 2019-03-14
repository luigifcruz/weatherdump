package parallel

func SerialRange(start, finish, threads int) map[int]int {
	r := make(map[int]int)

	for t := 0; t < threads; t++ {
		s := finish / threads * t
		f := finish / threads * (t + 1)

		if t == threads-1 {
			f = finish
		}

		r[s] = f
	}

	return r
}
