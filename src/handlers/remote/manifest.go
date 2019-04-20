package remote

import (
	"encoding/json"
	"net/http"
	"weatherdump/src/handlers"
	"weatherdump/src/protocols/helpers"
)

type manifestRequest struct {
	Datalink string                     `schema:"datalink,required"`
	Manifest helpers.ProcessingManifest `schema:"-"`
}

func (s *Remote) manifestHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	var req manifestRequest
	if err := decoder.Decode(&req, r.PostForm); err != nil {
		ResError(w, "INVALID_REQUEST", err.Error())
		return
	}

	processor := handlers.AvailableProcessors[req.Datalink]("", nil)
	req.Manifest = processor.GetProductsManifest()

	request, _ := json.Marshal(req)
	ResSuccess(w, "MANIFEST", string(request))
}
