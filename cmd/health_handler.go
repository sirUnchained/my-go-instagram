package main

import (
	"net/http"

	"github.com/sirUnchained/my-go-instagram/internal/scripts"
)

func (s *server) checkHealthHandler(w http.ResponseWriter, r *http.Request) {
	scripts.WriteJson(w, http.StatusOK, map[string]string{"msg": "all right"})
}
