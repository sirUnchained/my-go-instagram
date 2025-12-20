package scripts

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func WriteJson(w http.ResponseWriter, status int, data any) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	return json.NewEncoder(w).Encode(data)
}

func ErrorResponse(w http.ResponseWriter, status int, data any) error {
	var errorRes struct {
		Error any `json:"error"`
	}

	errorRes.Error = data

	return WriteJson(w, status, errorRes)
}

func JsonResponse(w http.ResponseWriter, status int, data any) error {
	var response struct {
		Data any `json:"error"`
	}

	response.Data = data

	return WriteJson(w, status, response)
}

func ReadJson(w http.ResponseWriter, r *http.Request, playload any) error {
	maxBytes := 1_048_578
	r.Body = http.MaxBytesReader(w, r.Body, int64(maxBytes))

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	return decoder.Decode(playload)
}

func ReadForm(w http.ResponseWriter, r *http.Request, playload any) error {
	// set 20 mb size limit
	maxSize := 2_0971_560
	if err := r.ParseMultipartForm(int64(maxSize)); err != nil {
		return err
	}

	// get all files and check counts
	files := r.MultipartForm.File["files"]
	if len(files) > 5 {
		// todo create error types
		return fmt.Errorf("file limit reached! oly 5 files allowed")
	}

	// process each file
	for i, fileHeader := range files {
		file, err := fileHeader.Open()
		if err != nil {
			return err
		}
		defer file.Close()

		fileBytes, err := io.ReadAll(file)
		if err != nil {
			return err
		}

		fmt.Printf("File %d: %s (%d bytes)\n", i+1, fileHeader.Filename, len(fileBytes))
	}

	return nil
}
