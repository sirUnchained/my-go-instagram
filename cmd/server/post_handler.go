package main

import (
	"net/http"

	"github.com/sirUnchained/my-go-instagram/internal/payloads"
	"github.com/sirUnchained/my-go-instagram/internal/scripts"
)

func (s *server) createPostHandler(w http.ResponseWriter, r *http.Request) {
	// var postP payloads.CreatePostPayload
	var fileP []payloads.CreateFilePayload
	user := scripts.GetUserFromContext(r)

	status, err := scripts.ReadFormFiles(w, r, user.Id, &fileP)
	switch status {
	case http.StatusBadRequest:
		s.badRequestResponse(w, r, err)
		return
	case http.StatusInternalServerError:
		s.internalServerErrorResponse(w, r, err)
		return
	}

	if err := scripts.JsonResponse(w, http.StatusCreated, fileP); err != nil {
		s.internalServerErrorResponse(w, r, err)
		return
	}
}
