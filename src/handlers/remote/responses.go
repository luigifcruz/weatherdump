package remote

import (
	"encoding/json"
	"net/http"
)

// Response standard response structure.
type Response struct {
	Res  bool
	Code string
	Data string
}

// ResError ends the current request with Code 400 (Bad Request).
// Will return standard response structure.
func ResError(w http.ResponseWriter, code, data string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusBadRequest)
	json.NewEncoder(w).Encode(Response{false, code, data})
}

// ResSuccess ends the current request with Code 200 (Success).
// Will return standard response structure.
func ResSuccess(w http.ResponseWriter, code, data string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(Response{true, code, data})
}
