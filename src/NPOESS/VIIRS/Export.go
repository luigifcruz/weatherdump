package VIIRS

import (
	"bytes"
	"fmt"
	"image"
	"image/png"
	"os"
	"path/filepath"

	"gopkg.in/gographics/imagick.v2/imagick"
)

func ExportGrayscale(buf []byte, e Channel, outputFolder string) {
	outputName, _ := filepath.Abs(fmt.Sprintf("%s/%s.png", outputFolder, e.fileName))

	img := image.NewGray16(image.Rect(0, 0, int(e.width), int(e.height)))
	img.Pix = buf

	pngImg := new(bytes.Buffer)
	png.Encode(pngImg, img)

	mw := imagick.NewMagickWand()
	defer mw.Destroy()

	mw.ReadImageBlob(pngImg.Bytes())
	mw.EqualizeImage()
	mw.FlopImage()
	mw.WriteImage(outputName)
}

func ExportTrueColor(outputFolder string, r, g, b *Channel) {
	fmt.Println("[VIIRS] Saving true color image.")

	R, _ := filepath.Abs(fmt.Sprintf("/tmp/%s.png", r.fileName))
	G, _ := filepath.Abs(fmt.Sprintf("/tmp/%s.png", g.fileName))
	B, _ := filepath.Abs(fmt.Sprintf("/tmp/%s.png", b.fileName))

	RGB, _ := filepath.Abs(fmt.Sprintf("%s/TRUECOLOR_VIIRS_%s.png", outputFolder, r.endTime.GetZulu()))

	if _, err := os.Stat(R); os.IsNotExist(err) {
		fmt.Println("[VIIRS] Red channel doesn't exists. Can't create true-color product.")
		return
	}

	if _, err := os.Stat(G); os.IsNotExist(err) {
		fmt.Println("[VIIRS] Green channel doesn't exists. Can't create true-color product.")
		fmt.Println(G)
		return
	}

	if _, err := os.Stat(B); os.IsNotExist(err) {
		fmt.Println("[VIIRS] Blue channel doesn't exists. Can't create true-color product.")
		return
	}

	// Load all channels for RGB.
	mwR := imagick.NewMagickWand()
	defer mwR.Destroy()

	mwG := imagick.NewMagickWand()
	defer mwG.Destroy()

	mwB := imagick.NewMagickWand()
	defer mwB.Destroy()

	mwR.ReadImage(R)
	mwG.ReadImage(G)
	mwB.ReadImage(B)

	// Merge them togheter to create True Color Image.
	mwRGB := imagick.NewMagickWand()
	defer mwRGB.Destroy()

	mwRGB.AddImage(mwR)
	mwRGB.AddImage(mwG)
	mwRGB.AddImage(mwB)
	mwRGB.ResetIterator()

	mwRGB = mwRGB.CombineImages(imagick.CHANNEL_RED | imagick.CHANNEL_GREEN | imagick.CHANNEL_BLUE)
	mwRGB.WriteImage(RGB)
}
