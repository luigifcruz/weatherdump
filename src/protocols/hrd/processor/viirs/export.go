package viirs

import (
	"fmt"
	"image"
	"image/png"
	"os"
	"path/filepath"
	"weather-dump/src/tools/imagery"
)

func ProcessImage(buf *[]byte, e Channel) {
	array := ConvertToU16(*buf)
	imagery.HistogramEqualizationU16(&array)
	imagery.FlopU16(&array, int(e.width))
	*buf = ConvertToByte(array)
}

func ExportGrayscale(buf *[]byte, e Channel, outputFolder string) {
	outputName, _ := filepath.Abs(fmt.Sprintf("%s/%s.png", outputFolder, e.fileName))

	img := image.NewGray16(image.Rect(0, 0, int(e.width), int(e.height)))
	img.Pix = *buf

	outputFile, err := os.Create(outputName)
	if err != nil {
		fmt.Println("[EXP] Error saving final image...", err)
	}
	png.Encode(outputFile, img)
	outputFile.Close()
}
