package VIIRS

import (
	"fmt"
	"image"
	"image/draw"
	"image/png"
	"math"
	"os"

	"github.com/nfnt/resize"
)

type BowTie struct {
	x0 int
	y0 int
	x1 int
	y1 int
}

func IsBlackOrWhite(r, g, b uint32) bool {
	return (r == 0 || r == 65532) && (g == 0 || g == 65532) && (b == 0 || b == 65532)
}

func LinearInterp(r float32, q0, q1 uint32) uint16 {
	return uint16(float32(q0)*r + float32(q1)*(1-r))
}

func CosineInterp(r float32, q0, q1 uint32) uint16 {
	var mu = float64(r)
	var mu2 = (1 - math.Cos(mu*math.Pi)) / 2

	return uint16(float64(q0)*(1-mu2) + float64(q1)*(mu2))
}

func ProcessBowTie(m *image.Gray16, bt BowTie) {
	var gm = m
	var w = bt.x1 - bt.x0
	var h = bt.y1 - bt.y0

	var btImg = image.NewGray16(image.Rect(0, 0, w, 2))
	// Slice Two lines
	for x := 0; x < w; x++ {
		btImg.Set(x, 0, gm.At(bt.x0+x, bt.y0-1))
		btImg.Set(x, 1, gm.At(bt.x0+x, bt.y1+1))
	}

	btImgRes := resize.Resize(uint(float32(w)), uint(h), btImg, resize.Lanczos3)

	for x := 0; x < w; x++ {
		for y := 0; y < h; y++ {
			gm.Set(bt.x0+x, bt.y0+y, btImgRes.At(x, y))
		}
	}
}

func PerformInterpolation(img *image.Gray16, cs ChannelParameters) {
	var bowTies []BowTie

	bounds := img.Bounds()

	fmt.Printf("[INTERPOLATOR] Interpolating Channel %s\n", cs.ChannelName)

	var x = bounds.Min.X
	for z := 0; z < 6; z++ {
		if z == 2 || z == 3 {
			x += cs.AggregationZoneWidth[z]
			continue
		}
		var csw = cs.AggregationZoneWidth[z]
		var csh = cs.BowTieHeight[z]
		for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
			r, g, b, _ := img.At(x, y).RGBA()
			if IsBlackOrWhite(r, g, b) {
				var nextY = y + csh
				if nextY > bounds.Max.Y || nextY+1 > bounds.Max.Y {
					break
				}
				nr, ng, nb, _ := img.At(x, nextY).RGBA()
				n2r, n2g, n2b, _ := img.At(x, nextY+1).RGBA()
				if IsBlackOrWhite(nr, ng, nb) && !IsBlackOrWhite(n2r, n2g, n2b) {
					bt := BowTie{
						x0: x,
						y0: y,
						x1: x + csw + 1,
						y1: y + csh + 1,
					}
					bowTies = append(bowTies, bt)
					y += cs.BowTieHeight[0]
				}
			}
		}
		x += cs.AggregationZoneWidth[z]
	}

	outputFile, _ := os.Create(fmt.Sprintf("mask_%s.png", cs.ChannelName))
	imas := image.NewGray(image.Rect(0, 0, img.Bounds().Dx(), img.Bounds().Dy()))
	draw.Draw(imas, imas.Bounds(), image.Black, image.ZP, draw.Src)
	for _, bowTie := range bowTies {
		var rect = image.Rect(bowTie.x0, bowTie.y0, bowTie.x1, bowTie.y1)
		draw.Draw(imas, rect.Bounds(), image.White, image.ZP, draw.Src)
	}
	encoder := png.Encoder{CompressionLevel: png.NoCompression}
	encoder.Encode(outputFile, imas)
	outputFile.Close()

	for _, bowTie := range bowTies {
		ProcessBowTie(img, bowTie)
	}
}
