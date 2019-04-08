package remote

import (
	"github.com/nfnt/resize"
	"image"
	"image/jpeg"
	"net/http"
	"os"
)

func (s *Remote) thumbnailHandler(w http.ResponseWriter, r *http.Request) {
	filePath := r.FormValue("filepath")

	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		ResError(w, "INPUT_FILE_NOT_FOUND", "")
		return
	}

	file, err := os.Open(filePath)
	img, _, err := image.Decode(file)

	if err != nil {
		ResError(w, "INVALID_FILE", "")
		return
	}
	file.Close()

	m := resize.Resize(350, 0, img, resize.Lanczos3)

	w.Header().Set("Content-Type", "image/jpeg")
	jpeg.Encode(w, m, nil)
}
