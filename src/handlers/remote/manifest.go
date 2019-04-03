package remote

import (
	"encoding/json"
	"net/http"
	"weather-dump/src/assets"
	"weather-dump/src/handlers"
)

type manifestRequest struct {
	Datalink string                    `schema:"datalink,required"`
	Manifest assets.ProcessingManifest `schema:"-"`
}

func (s *Remote) manifestHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	var req manifestRequest
	if err := decoder.Decode(&req, r.PostForm); err != nil {
		ResError(w, "INVALID_REQUEST", err.Error())
		return
	}

	processor := handlers.AvailableProcessors[req.Datalink]("")
	req.Manifest = processor.GetProductsManifest()

	request, _ := json.Marshal(req)
	ResSuccess(w, "MANIFEST", string(request))
}
