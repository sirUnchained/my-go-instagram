package main

import (
	"net/http"

	"github.com/sirUnchained/my-go-instagram/internal/helpers"
)

func (s *server) checkHealthHandler(w http.ResponseWriter, r *http.Request) {
	helpers.WriteJson(w, http.StatusOK, map[string]string{"msg": "all right"})
}
