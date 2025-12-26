package scripts

import (
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/sirUnchained/my-go-instagram/internal/payloads"
)

const (
	MAX_UPLOAD_SIZE = 10 * 1024 * 1024 // 10MB maximum for each file
	MAX_FILES       = 5                // maximum 5 files
	UPLOADS_DIR     = "./public/uploads/"
)

var allowedMimeTypes = map[string]bool{
	// image formats
	"image/jpeg":    true,
	"image/png":     true,
	"image/gif":     true,
	"image/webp":    true,
	"image/svg+xml": true,

	// document formats
	"application/pdf":    true,
	"text/plain":         true,
	"application/msword": true,
	"application/vnd.openxmlformats-officedocument.wordprocessingml.document": true,

	// video formats
	"video/mp4":        true,
	"video/mpeg":       true,
	"video/ogg":        true,
	"video/webm":       true,
	"video/quicktime":  true,
	"video/x-msvideo":  true,
	"video/x-matroska": true,
	"video/x-flv":      true,

	// audio formats
	"audio/mpeg": true,
	"audio/wav":  true,
	"audio/ogg":  true,
	"audio/webm": true,
}

func ReadFormFiles(w http.ResponseWriter, r *http.Request, userid int64, playload *[]payloads.CreateFilePayload) (int, error) {
	w.Header().Set("Content-Type", "application/json")

	userDir := filepath.Join(UPLOADS_DIR, fmt.Sprintf("%d", userid))
	if _, err := os.Stat(userDir); os.IsNotExist(err) {
		if err := os.MkdirAll(userDir, 0775); err != nil {
			return http.StatusInternalServerError, err
		}
	}

	r.Body = http.MaxBytesReader(w, r.Body, MAX_UPLOAD_SIZE*MAX_FILES)

	if err := r.ParseMultipartForm(MAX_UPLOAD_SIZE * MAX_FILES); err != nil {
		return http.StatusBadRequest, fmt.Errorf("maximum upload size should be 50MB and files count shoud be between 1 to 5")
	}

	files := r.MultipartForm.File["files"]
	err := validateFilesCount(files)
	if err != nil {
		return http.StatusBadRequest, err
	}

	// processing files
	for _, fileHeader := range files {
		// open file
		file, err := fileHeader.Open()
		if err != nil {
			return http.StatusInternalServerError, err
		}
		defer file.Close()

		// check the file size
		if fileHeader.Size > MAX_UPLOAD_SIZE {
			return http.StatusBadRequest, fmt.Errorf("one of the files are larger than %d", MAX_UPLOAD_SIZE)
		}

		// check file type
		switch code, err := validateFileType(file); code {
		case http.StatusInternalServerError:
			return http.StatusInternalServerError, err

		case http.StatusBadRequest:
			return http.StatusInternalServerError, err
		}

		// generate a unique name + the file itself
		fileExt := filepath.Ext(fileHeader.Filename)
		uniqueName := fmt.Sprintf("%d-%s-%s", time.Now().UnixNano(), generateRandomString(16), fileExt)
		path, err := generateFile(file, uniqueName, userid)
		if err != nil {
			return http.StatusInternalServerError, err
		}

		newFile := payloads.CreateFilePayload{
			Filename:  uniqueName,
			Filepath:  path,
			SizeBytes: int(fileHeader.Size),
			Creator:   userid,
		}

		*playload = append(*playload, newFile)
	}

	return http.StatusCreated, nil
}

func validateFilesCount(files []*multipart.FileHeader) error {
	if len(files) == 0 {
		return fmt.Errorf("no files uploaded, minimum is '1' file")
	}

	if len(files) > MAX_FILES {
		return fmt.Errorf("too many files, amximum is '%d' files", MAX_FILES)
	}

	return nil
}

func validateFileType(file multipart.File) (int, error) {
	buffer := make([]byte, 512)
	_, err := file.Read(buffer)
	if err != nil && err != io.EOF {
		// if we faile here it means somthing went wrong in server
		return http.StatusInternalServerError, err
	}

	// reset file pointer to begining
	file.Seek(0, 0)

	mimeType := http.DetectContentType(buffer)
	if !allowedMimeTypes[mimeType] {
		return http.StatusBadRequest, fmt.Errorf("mimeType not allowed")
	}

	return http.StatusOK, nil
}

func generateRandomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[time.Now().UnixNano()%int64(len(charset))]
	}
	return string(b)
}

func generateFile(file multipart.File, name string, userid int64) (string, error) {
	// create destination file
	dstPath := filepath.Join(fmt.Sprint(UPLOADS_DIR, "/", userid), name)
	dst, err := os.Create(dstPath)
	if err != nil {
		return "", err
	}
	defer dst.Close()

	// copy file to destination
	if _, err := io.Copy(dst, file); err != nil {
		return "", err
	}

	return dstPath, nil
}
