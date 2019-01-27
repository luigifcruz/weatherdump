package Sensor

type maskDC struct {
	code []bool
	len  int
}

var dcCategories = [12]maskDC{
	maskDC{[]bool{false, false}, 0},
	maskDC{[]bool{false, true, false}, 1},
	maskDC{[]bool{false, true, true}, 2},
	maskDC{[]bool{true, false, false}, 3},
	maskDC{[]bool{true, false, true}, 4},
	maskDC{[]bool{true, true, false}, 5},
	maskDC{[]bool{true, true, true, false}, 6},
	maskDC{[]bool{true, true, true, true, false}, 7},
	maskDC{[]bool{true, true, true, true, true, false}, 8},
	maskDC{[]bool{true, true, true, true, true, true, false}, 9},
	maskDC{[]bool{true, true, true, true, true, true, true, false}, 10},
	maskDC{[]bool{true, true, true, true, true, true, true, true, false}, 11},
}
