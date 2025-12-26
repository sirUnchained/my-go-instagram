package main

import (
	"net/http"

	"github.com/sirUnchained/my-go-instagram/internal/payloads"
	"github.com/sirUnchained/my-go-instagram/internal/scripts"
)

func (s *server) createPostHandler(w http.ResponseWriter, r *http.Request) {
	var postP payloads.CreatePostPayload

	status, err := scripts.ReadFormFiles(w, r, postP)
	switch status {
	case http.StatusBadRequest:
		s.badRequestResponse(w, r, err)
		return
	case http.StatusInternalServerError:
		s.internalServerErrorResponse(w, r, err)
	}
}
