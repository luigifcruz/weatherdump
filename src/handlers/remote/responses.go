package remote

import (
	"encoding/json"
	"net/http"
)

type Response struct {
	Res         bool
	Code        string
	Description string
}

func ResError(w http.ResponseWriter, code, description string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusBadRequest)
	json.NewEncoder(w).Encode(Response{false, code, description})
}

func ResSuccess(w http.ResponseWriter, code, description string) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(Response{true, code, description})
}
