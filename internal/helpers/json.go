package helpers

import (
	"encoding/json"
	"net/http"
)

type ErrorRes struct {
	Error any `json:"error"`
}

type DataRes struct {
	Data any `json:"data"`
}

func WriteJson(w http.ResponseWriter, status int, data any) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	return json.NewEncoder(w).Encode(data)
}

func ErrorResponse(w http.ResponseWriter, status int, data any) error {
	var errorRes ErrorRes

	errorRes.Error = data

	return WriteJson(w, status, &errorRes)
}

func JsonResponse(w http.ResponseWriter, status int, data any) error {
	var response DataRes

	response.Data = data

	return WriteJson(w, status, &response)
}

func ReadJson(w http.ResponseWriter, r *http.Request, playload any) error {
	maxBytes := 1_048_578
	r.Body = http.MaxBytesReader(w, r.Body, int64(maxBytes))

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	return decoder.Decode(playload)
}
