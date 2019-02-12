package BISMW

// Made by Artyom Litvinovich
// Ported to GoLang by Luigi Cruz

import "math"

var cosine [8][8]float64
var alpha [8]float64

func initCos() {
	for y := 0; y < 8; y++ {
		for x := 0; x < 8; x++ {
			cosine[y][x] = math.Cos(math.Pi / 16.0 * (2.0*float64(y) + 1.0) * float64(x))
		}
	}

	alpha[0] = 1.0 / math.Sqrt(2.0)
	for x := 1; x < 8; x++ {
		alpha[x] = 1.0
	}
}

func calculateIdct(res, inpt *[64]float64) {
	initCos()
	for y := 0; y < 8; y++ {
		for x := 0; x < 8; x++ {
			var s float64
			for u := 0; u < 8; u++ {
				coeff := inpt[0*8+u] * alpha[0] * cosine[y][0]
				coeff += inpt[1*8+u] * alpha[1] * cosine[y][1]
				coeff += inpt[2*8+u] * alpha[2] * cosine[y][2]
				coeff += inpt[3*8+u] * alpha[3] * cosine[y][3]
				coeff += inpt[4*8+u] * alpha[4] * cosine[y][4]
				coeff += inpt[5*8+u] * alpha[5] * cosine[y][5]
				coeff += inpt[6*8+u] * alpha[6] * cosine[y][6]
				coeff += inpt[7*8+u] * alpha[7] * cosine[y][7]
				s += alpha[u] * cosine[x][u] * coeff
			}
			res[y*8+x] = s / 4.0
		}
	}
}
