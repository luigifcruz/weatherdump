package Sensor

type maskAC struct {
	code []bool
	clen int
	zlen int
}

var acCategories = [162]maskAC{
	maskAC{[]bool{true, false, true, false}, 0, 0},                                                                              // 0/0
	maskAC{[]bool{false, false}, 1, 0},                                                                                          // 0/1
	maskAC{[]bool{false, true}, 2, 0},                                                                                           // 0/2
	maskAC{[]bool{true, false, false}, 3, 0},                                                                                    // 0/3
	maskAC{[]bool{true, false, true, true}, 4, 0},                                                                               // 0/4
	maskAC{[]bool{true, true, false, true, false}, 5, 0},                                                                        // 0/5
	maskAC{[]bool{true, true, true, true, false, false, false}, 6, 0},                                                           // 0/6
	maskAC{[]bool{true, true, true, true, true, false, false, false}, 7, 0},                                                     // 0/7
	maskAC{[]bool{true, true, true, true, true, true, false, true, true, false}, 8, 0},                                          // 0/8
	maskAC{[]bool{true, true, true, true, true, true, true, true, true, false, false, false, false, false, true, false}, 9, 0},  // 0/9
	maskAC{[]bool{true, true, true, true, true, true, true, true, true, false, false, false, false, false, true, true}, 10, 0},  // 0/10
	maskAC{[]bool{true, true, false, false}, 1, 1},                                                                              // 1/1
	maskAC{[]bool{true, true, false, true, true}, 2, 1},                                                                         // 1/2
	maskAC{[]bool{true, true, true, true, false, false, true}, 3, 1},                                                            // 1/2
	maskAC{[]bool{true, true, true, true, true, false, true, true, false}, 4, 1},                                                // 1/3
	maskAC{[]bool{true, true, true, true, true, true, true, false, true, true, false}, 5, 1},                                    // 1/4
	maskAC{[]bool{true, true, true, true, true, true, true, true, true, false, false, false, false, true, false, false}, 6, 1},  // 1/6
	maskAC{[]bool{true, true, true, true, true, true, true, true, true, false, false, false, false, true, false, true}, 7, 1},   // 1/7
	maskAC{[]bool{true, true, true, true, true, true, true, true, true, false, false, false, false, true, true, false}, 8, 1},   // 1/8
	maskAC{[]bool{true, true, true, true, true, true, true, true, true, false, false, false, false, true, true, true}, 9, 1},    // 1/9
	maskAC{[]bool{true, true, true, true, true, true, true, true, true, false, false, false, true, false, false, false}, 10, 1}, // 1/10
	maskAC{[]bool{true, true, true, false, false}, 1, 2},                                                                        // 2/1
	maskAC{[]bool{true, true, true, true, true, false, false, true}, 2, 2},                                                      // 2/2
	maskAC{[]bool{true, true, true, true, true, true, false, true, true, true}, 3, 2},                                           // 2/3
	maskAC{[]bool{true, true, true, true, true, true, true, true, false, true, false, false}, 4, 2},                             // 2/4
	maskAC{[]bool{true, true, true, true, true, true, true, true, true, false, false, false, true, false, false, true}, 5, 2},   // 2/5
	maskAC{[]bool{true, true, true, true, true, true, true, true, true, false, false, false, true, false, true, false}, 6, 2},   // 2/6
	maskAC{[]bool{true, true, true, true, true, true, true, true, true, false, false, false, true, false, true, true}, 7, 2},    // 2/7
	maskAC{[]bool{true, true, true, true, true, true, true, true, true, false, false, false, true, true, false, false}, 8, 2},   // 2/8
	maskAC{[]bool{true, true, true, true, true, true, true, true, true, false, false, false, true, true, false, true}, 9, 2},    // 2/9
	maskAC{[]bool{true, true, true, true, true, true, true, true, true, false, false, false, true, true, true, false}, 10, 2},   // 2/10
	maskAC{[]bool{true, true, true, false, true, false}, 1, 3},                                                                  // 3/1
	maskAC{[]bool{true, true, true, true, true, false, true, true, true}, 2, 3},                                                 // 3/2
	maskAC{[]bool{true, true, true, true, true, true, true, true, false, true, false, true}, 3, 3},                              // 3/3
	maskAC{[]bool{true, true, true, true, true, true, true, true, true, false, false, false, true, true, true, true}, 4, 3},     // 3/4
	maskAC{[]bool{true, true, true, true, true, true, true, true, true, false, false, true, false, false, false, false}, 5, 3},  // 3/5
	maskAC{[]bool{true, true, true, true, true, true, true, true, true, false, false, true, false, false, false, true}, 6, 3},   // 3/6
	maskAC{[]bool{true, true, true, true, true, true, true, true, true, false, false, true, false, false, true, false}, 7, 3},   // 3/7
	maskAC{[]bool{true, true, true, true, true, true, true, true, true, false, false, true, false, false, true, true}, 8, 3},    // 3/8
	maskAC{[]bool{true, true, true, true, true, true, true, true, true, false, false, true, false, true, false, false}, 9, 3},   // 3/9
	maskAC{[]bool{true, true, true, true, true, true, true, true, true, false, false, true, false, true, false, true}, 10, 3},   // 3/10
	maskAC{[]bool{true, true, true, false, true, true}, 1, 4},
	maskAC{[]bool{true, true, true, true, true, true, true, false, false, false}, 2, 4},
	maskAC{[]bool{true, true, true, true, true, true, true, true, true, false, false, true, false, true, true, false}, 3, 4},
	maskAC{[]bool{true, true, true, true, true, true, true, true, true, false, false, true, false, true, true, true}, 4, 4},
	maskAC{[]bool{true, true, true, true, true, true, true, true, true, false, false, true, true, false, false, false}, 5, 4},
	maskAC{[]bool{true, true, true, true, true, true, true, true, true, false, false, true, true, false, false, true}, 6, 4},
	maskAC{[]bool{true, true, true, true, true, true, true, true, true, false, false, true, true, false, true, false}, 7, 4},
	maskAC{[]bool{true, true, true, true, true, true, true, true, true, false, false, true, true, false, true, true}, 8, 4},
	maskAC{[]bool{true, true, true, true, true, true, true, true, true, false, false, true, true, true, false, false}, 9, 4},
	maskAC{[]bool{true, true, true, true, true, true, true, true, true, false, false, true, true, true, false, true}, 10, 4},
	maskAC{[]bool{true, true, true, true, false, true, false}, 1, 5},
	maskAC{[]bool{true, true, true, true, true, true, true, false, true, true, true}, 2, 5},
	maskAC{[]bool{true, true, true, true, true, true, true, true, true, false, false, true, true, true, true, false}, 3, 5},
	maskAC{[]bool{true, true, true, true, true, true, true, true, true, false, false, true, true, true, true, true}, 4, 5},
	maskAC{[]bool{true, true, true, true, true, true, true, true, true, false, true, false, false, false, false, false}, 5, 5},
	maskAC{[]bool{true, true, true, true, true, true, true, true, true, false, true, false, false, false, false, true}, 6, 5},
	maskAC{[]bool{true, true, true, true, true, true, true, true, true, false, true, false, false, false, true, false}, 7, 5},
	maskAC{[]bool{true, true, true, true, true, true, true, true, true, false, true, false, false, false, true, true}, 8, 5},
	maskAC{[]bool{true, true, true, true, true, true, true, true, true, false, true, false, false, true, false, false}, 9, 5},
	maskAC{[]bool{true, true, true, true, true, true, true, true, true, false, true, false, false, true, false, true}, 10, 5},
	maskAC{[]bool{true, true, true, true, false, true, true}, 1, 6},
	maskAC{[]bool{true, true, true, true, true, true, true, true, false, true, true, false}, 2, 6},
	maskAC{[]bool{true, true, true, true, true, true, true, true, true, false, true, false, false, true, true, false}, 3, 6},
	maskAC{[]bool{true, true, true, true, true, true, true, true, true, false, true, false, false, true, true, true}, 4, 6},
	maskAC{[]bool{true, true, true, true, true, true, true, true, true, false, true, false, true, false, false, false}, 5, 6},
	maskAC{[]bool{true, true, true, true, true, true, true, true, true, false, true, false, true, false, false, true}, 6, 6},
	maskAC{[]bool{true, true, true, true, true, true, true, true, true, false, true, false, true, false, true, false}, 7, 6},
	maskAC{[]bool{true, true, true, true, true, true, true, true, true, false, true, false, true, false, true, true}, 8, 6},
	maskAC{[]bool{true, true, true, true, true, true, true, true, true, false, true, false, true, true, false, false}, 9, 6},
	maskAC{[]bool{true, true, true, true, true, true, true, true, true, false, true, false, true, true, false, true}, 10, 6},
	maskAC{[]bool{true, true, true, true, true, false, true, false}, 1, 7},
	maskAC{[]bool{true, true, true, true, true, true, true, true, false, true, true, true}, 2, 7},
	maskAC{[]bool{true, true, true, true, true, true, true, true, true, false, true, false, true, true, true, false}, 3, 7},
	maskAC{[]bool{true, true, true, true, true, true, true, true, true, false, true, false, true, true, true, true}, 4, 7},
	maskAC{[]bool{true, true, true, true, true, true, true, true, true, false, true, true, false, false, false, false}, 5, 7},
	maskAC{[]bool{true, true, true, true, true, true, true, true, true, false, true, true, false, false, false, true}, 6, 7},
	maskAC{[]bool{true, true, true, true, true, true, true, true, true, false, true, true, false, false, true, false}, 7, 7},
	maskAC{[]bool{true, true, true, true, true, true, true, true, true, false, true, true, false, false, true, true}, 8, 7},
	maskAC{[]bool{true, true, true, true, true, true, true, true, true, false, true, true, false, true, false, false}, 9, 7},
	maskAC{[]bool{true, true, true, true, true, true, true, true, true, false, true, true, false, true, false, true}, 10, 7},
	maskAC{[]bool{true, true, true, true, true, true, false, false, false}, 1, 8},
	maskAC{[]bool{true, true, true, true, true, true, true, true, true, false, false, false, false, false, false}, 2, 8},
	maskAC{[]bool{true, true, true, true, true, true, true, true, true, false, true, true, false, true, true, false}, 3, 8},
	maskAC{[]bool{true, true, true, true, true, true, true, true, true, false, true, true, false, true, true, true}, 4, 8},
	maskAC{[]bool{true, true, true, true, true, true, true, true, true, false, true, true, true, false, false, false}, 5, 8},
	maskAC{[]bool{true, true, true, true, true, true, true, true, true, false, true, true, true, false, false, true}, 6, 8},
	maskAC{[]bool{true, true, true, true, true, true, true, true, true, false, true, true, true, false, true, false}, 7, 8},
	maskAC{[]bool{true, true, true, true, true, true, true, true, true, false, true, true, true, false, true, true}, 8, 8},
	maskAC{[]bool{true, true, true, true, true, true, true, true, true, false, true, true, true, true, false, false}, 9, 8},
	maskAC{[]bool{true, true, true, true, true, true, true, true, true, false, true, true, true, true, false, true}, 10, 8},
	maskAC{[]bool{true, true, true, true, true, true, false, false, true}, 1, 9},
	maskAC{[]bool{true, true, true, true, true, true, true, true, true, false, true, true, true, true, true, false}, 2, 9},
	maskAC{[]bool{true, true, true, true, true, true, true, true, true, false, true, true, true, true, true, true}, 3, 9},
	maskAC{[]bool{true, true, true, true, true, true, true, true, true, true, false, false, false, false, false, false}, 4, 9},
	maskAC{[]bool{true, true, true, true, true, true, true, true, true, true, false, false, false, false, false, true}, 5, 9},
	maskAC{[]bool{true, true, true, true, true, true, true, true, true, true, false, false, false, false, true, false}, 6, 9},
	maskAC{[]bool{true, true, true, true, true, true, true, true, true, true, false, false, false, false, true, true}, 7, 9},
	maskAC{[]bool{true, true, true, true, true, true, true, true, true, true, false, false, false, true, false, false}, 8, 9},
	maskAC{[]bool{true, true, true, true, true, true, true, true, true, true, false, false, false, true, false, true}, 9, 9},
	maskAC{[]bool{true, true, true, true, true, true, true, true, true, true, false, false, false, true, true, false}, 10, 9},
	maskAC{[]bool{true, true, true, true, true, true, false, true, false}, 1, 10},
	maskAC{[]bool{true, true, true, true, true, true, true, true, true, true, false, false, false, true, true, true}, 2, 10},
	maskAC{[]bool{true, true, true, true, true, true, true, true, true, true, false, false, true, false, false, false}, 3, 10},
	maskAC{[]bool{true, true, true, true, true, true, true, true, true, true, false, false, true, false, false, true}, 4, 10},
	maskAC{[]bool{true, true, true, true, true, true, true, true, true, true, false, false, true, false, true, false}, 5, 10},
	maskAC{[]bool{true, true, true, true, true, true, true, true, true, true, false, false, true, false, true, true}, 6, 10},
	maskAC{[]bool{true, true, true, true, true, true, true, true, true, true, false, false, true, true, false, false}, 7, 10},
	maskAC{[]bool{true, true, true, true, true, true, true, true, true, true, false, false, true, true, false, true}, 8, 10},
	maskAC{[]bool{true, true, true, true, true, true, true, true, true, true, false, false, true, true, true, false}, 9, 10},
	maskAC{[]bool{true, true, true, true, true, true, true, true, true, true, false, false, true, true, true, true}, 10, 10},
	maskAC{[]bool{true, true, true, true, true, true, true, false, false, true}, 1, 11},
	maskAC{[]bool{true, true, true, true, true, true, true, true, true, true, false, true, false, false, false, false}, 2, 11},
	maskAC{[]bool{true, true, true, true, true, true, true, true, true, true, false, true, false, false, false, true}, 3, 11},
	maskAC{[]bool{true, true, true, true, true, true, true, true, true, true, false, true, false, false, true, false}, 4, 11},
	maskAC{[]bool{true, true, true, true, true, true, true, true, true, true, false, true, false, false, true, true}, 5, 11},
	maskAC{[]bool{true, true, true, true, true, true, true, true, true, true, false, true, false, true, false, false}, 6, 11},
	maskAC{[]bool{true, true, true, true, true, true, true, true, true, true, false, true, false, true, false, true}, 7, 11},
	maskAC{[]bool{true, true, true, true, true, true, true, true, true, true, false, true, false, true, true, false}, 8, 11},
	maskAC{[]bool{true, true, true, true, true, true, true, true, true, true, false, true, false, true, true, true}, 9, 11},
	maskAC{[]bool{true, true, true, true, true, true, true, true, true, true, false, true, true, false, false, false}, 10, 11},
	maskAC{[]bool{true, true, true, true, true, true, true, false, true, false}, 1, 12},
	maskAC{[]bool{true, true, true, true, true, true, true, true, true, true, false, true, true, false, false, true}, 2, 12},
	maskAC{[]bool{true, true, true, true, true, true, true, true, true, true, false, true, true, false, true, false}, 3, 12},
	maskAC{[]bool{true, true, true, true, true, true, true, true, true, true, false, true, true, false, true, true}, 4, 12},
	maskAC{[]bool{true, true, true, true, true, true, true, true, true, true, false, true, true, true, false, false}, 5, 12},
	maskAC{[]bool{true, true, true, true, true, true, true, true, true, true, false, true, true, true, false, true}, 6, 12},
	maskAC{[]bool{true, true, true, true, true, true, true, true, true, true, false, true, true, true, true, false}, 7, 12},
	maskAC{[]bool{true, true, true, true, true, true, true, true, true, true, false, true, true, true, true, true}, 8, 12},
	maskAC{[]bool{true, true, true, true, true, true, true, true, true, true, true, false, false, false, false, false}, 9, 12},
	maskAC{[]bool{true, true, true, true, true, true, true, true, true, true, true, false, false, false, false, true}, 10, 12},
	maskAC{[]bool{true, true, true, true, true, true, true, true, false, false, false}, 1, 13},
	maskAC{[]bool{true, true, true, true, true, true, true, true, true, true, true, false, false, false, true, false}, 2, 13},
	maskAC{[]bool{true, true, true, true, true, true, true, true, true, true, true, false, false, false, true, true}, 3, 13},
	maskAC{[]bool{true, true, true, true, true, true, true, true, true, true, true, false, false, true, false, false}, 4, 13},
	maskAC{[]bool{true, true, true, true, true, true, true, true, true, true, true, false, false, true, false, true}, 5, 13},
	maskAC{[]bool{true, true, true, true, true, true, true, true, true, true, true, false, false, true, true, false}, 6, 13},
	maskAC{[]bool{true, true, true, true, true, true, true, true, true, true, true, false, false, true, true, true}, 7, 13},
	maskAC{[]bool{true, true, true, true, true, true, true, true, true, true, true, false, true, false, false, false}, 8, 13},
	maskAC{[]bool{true, true, true, true, true, true, true, true, true, true, true, false, true, false, false, true}, 9, 13},
	maskAC{[]bool{true, true, true, true, true, true, true, true, true, true, true, false, true, false, true, false}, 10, 13},
	maskAC{[]bool{true, true, true, true, true, true, true, true, true, true, true, false, true, false, true, true}, 1, 14},
	maskAC{[]bool{true, true, true, true, true, true, true, true, true, true, true, false, true, true, false, false}, 2, 14},
	maskAC{[]bool{true, true, true, true, true, true, true, true, true, true, true, false, true, true, false, true}, 3, 14},
	maskAC{[]bool{true, true, true, true, true, true, true, true, true, true, true, false, true, true, true, false}, 4, 14},
	maskAC{[]bool{true, true, true, true, true, true, true, true, true, true, true, false, true, true, true, true}, 5, 14},
	maskAC{[]bool{true, true, true, true, true, true, true, true, true, true, true, true, false, false, false, false}, 6, 14},
	maskAC{[]bool{true, true, true, true, true, true, true, true, true, true, true, true, false, false, false, true}, 7, 14},
	maskAC{[]bool{true, true, true, true, true, true, true, true, true, true, true, true, false, false, true, false}, 8, 14},
	maskAC{[]bool{true, true, true, true, true, true, true, true, true, true, true, true, false, false, true, true}, 9, 14},
	maskAC{[]bool{true, true, true, true, true, true, true, true, true, true, true, true, false, true, false, false}, 10, 14},
	maskAC{[]bool{true, true, true, true, true, true, true, true, false, false, true}, 0, 15},
	maskAC{[]bool{true, true, true, true, true, true, true, true, true, true, true, true, false, true, false, true}, 2, 15},
	maskAC{[]bool{true, true, true, true, true, true, true, true, true, true, true, true, false, true, true, false}, 2, 15},
	maskAC{[]bool{true, true, true, true, true, true, true, true, true, true, true, true, false, true, true, true}, 3, 15},
	maskAC{[]bool{true, true, true, true, true, true, true, true, true, true, true, true, true, false, false, false}, 4, 15},
	maskAC{[]bool{true, true, true, true, true, true, true, true, true, true, true, true, true, false, false, true}, 5, 15},
	maskAC{[]bool{true, true, true, true, true, true, true, true, true, true, true, true, true, false, true, false}, 6, 15},
	maskAC{[]bool{true, true, true, true, true, true, true, true, true, true, true, true, true, false, true, true}, 7, 15},
	maskAC{[]bool{true, true, true, true, true, true, true, true, true, true, true, true, true, true, false, false}, 8, 15},
	maskAC{[]bool{true, true, true, true, true, true, true, true, true, true, true, true, true, true, false, true}, 9, 15},
	maskAC{[]bool{true, true, true, true, true, true, true, true, true, true, true, true, true, true, true, false}, 10, 15},
}
