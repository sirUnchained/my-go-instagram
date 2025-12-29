package helpers

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/sirUnchained/my-go-instagram/internal/payloads"
)

const (
	MAX_POST_UPLOAD_SIZE = 50 * 1024 * 1024
	MAX_POST_FILES       = 5
	MAX_CONTENT_SIZE     = 1 * 1024 * 1024
	MAX_AVATARS          = 1
	MAX_AVATAR_SIZE      = 1 * 512
	MAX_BIO_MUSICS       = 1
	UPLOADS_DIR          = "./public/uploads/"
)

var allowedMimeTypes = map[string]bool{
	// image formats
	"image/jpeg": true,
	"image/png":  true,
	"image/webp": true,

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

func ReadFormFiles(w http.ResponseWriter, r *http.Request, userid int64, playload any) (int, error) {
	userDir := filepath.Join(UPLOADS_DIR, fmt.Sprintf("%d", userid))
	if _, err := os.Stat(userDir); os.IsNotExist(err) {
		if err := os.MkdirAll(userDir, 0775); err != nil {
			return http.StatusInternalServerError, err
		}
	}

	switch p := playload.(type) {
	case *payloads.CreatePostPayload:
		maxSize := int64(MAX_POST_UPLOAD_SIZE*MAX_POST_FILES + MAX_CONTENT_SIZE)
		r.Body = http.MaxBytesReader(w, r.Body, maxSize)
		if err := r.ParseMultipartForm(maxSize); err != nil {
			if strings.Contains(err.Error(), "request body too large") {
				return http.StatusBadRequest, fmt.Errorf("maximum upload size is %dMB", maxSize/(1024*1024))
			}
			return http.StatusBadRequest, fmt.Errorf("failed to parse form: %w", err)
		}

		p.Creator = userid
		p.Description = r.FormValue("description")
		if len(p.Description) > 2048 {
			return http.StatusBadRequest, fmt.Errorf("description too long, maximum 2048 characters")
		}
		tagsStr := r.FormValue("tags")
		if tagsStr == "" {
			p.Tags = []string{}
		} else {
			if err := json.Unmarshal([]byte(tagsStr), &p.Tags); err != nil {
				log.Printf("Error unmarshaling tags: %v, tags string: %s", err, tagsStr)
				return http.StatusBadRequest, fmt.Errorf("invalid tags format: %w", err)
			}
		}

		files := r.MultipartForm.File["files"]
		err := validateFilesCount(files)
		if err != nil {
			return http.StatusBadRequest, err
		}

		// processing files
		var fileList []payloads.CreateFilePayload
		for _, fileHeader := range files {
			// open file
			file, err := fileHeader.Open()
			if err != nil {
				return http.StatusInternalServerError, err
			}
			defer file.Close()

			// check the file size
			if fileHeader.Size > MAX_POST_UPLOAD_SIZE {
				return http.StatusBadRequest, fmt.Errorf("one of the files are larger than %dMB", MAX_POST_UPLOAD_SIZE)
			}

			// check file type
			switch code, err := validateFileType(file); code {
			case http.StatusInternalServerError:
				return http.StatusInternalServerError, err

			case http.StatusBadRequest:
				return http.StatusBadRequest, err
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

			fileList = append(fileList, newFile)
		}
		p.Files = fileList
	case *payloads.CreateUserPayload:
		maxSize := int64(MAX_POST_UPLOAD_SIZE*MAX_POST_FILES + MAX_CONTENT_SIZE)
		r.Body = http.MaxBytesReader(w, r.Body, maxSize)
		if err := r.ParseMultipartForm(maxSize); err != nil {
			if strings.Contains(err.Error(), "request body too large") {
				return http.StatusBadRequest, fmt.Errorf("maximum upload size is %dMB", maxSize/(1024*1024))
			}
			return http.StatusBadRequest, fmt.Errorf("failed to parse form: %w", err)
		}

		p.Username = r.FormValue("username")
		p.Fullname = r.FormValue("fullname")
		p.Email = r.FormValue("email")
		p.Password = r.FormValue("password")
		p.Bio = r.FormValue("bio")
		avatarFile := r.MultipartForm.File["avatar"]
		if len(avatarFile) > 1 || len(avatarFile) < 1 {
			p.Avatar = payloads.CreateFilePayload{Creator: userid,
				Filename:  "deafult",
				Filepath:  "./public/deafult.jpeg",
				SizeBytes: 0,
			}
			return http.StatusCreated, nil
		}

		v := validator.New(validator.WithRequiredStructEnabled())
		if err := v.Struct(p); err != nil {
			return http.StatusBadRequest, err
		}

		for _, fileHeader := range avatarFile {
			// open file
			file, err := fileHeader.Open()
			if err != nil {
				return http.StatusInternalServerError, err
			}
			defer file.Close()

			// check the file size
			if fileHeader.Size > MAX_AVATAR_SIZE {
				return http.StatusBadRequest, fmt.Errorf("one of the files are larger than %dMB", MAX_POST_UPLOAD_SIZE)
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

			p.Avatar = payloads.CreateFilePayload{
				Filename:  uniqueName,
				Filepath:  path,
				SizeBytes: int(fileHeader.Size),
				Creator:   userid,
			}
		}
	}
	return http.StatusCreated, nil
}

func validateFilesCount(files []*multipart.FileHeader) error {
	if len(files) == 0 {
		return fmt.Errorf("no files uploaded, minimum is '1' file")
	}

	if len(files) > MAX_POST_FILES {
		return fmt.Errorf("too many files, maximum is '%d' files", MAX_POST_FILES)
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
