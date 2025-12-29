package main

import (
	"net/http"

	"github.com/sirUnchained/my-go-instagram/internal/helpers"
)

func (s *server) internalServerErrorResponse(w http.ResponseWriter, r *http.Request, err error) {
	s.logger.Errorln("internal seerver error", "method", r.Method, "path", r.URL.Path, "error", err.Error())
	helpers.ErrorResponse(w, http.StatusInternalServerError, "server ran into a problem, we will fix this as soon as we can")
}

func (s *server) badRequestResponse(w http.ResponseWriter, r *http.Request, err error) {
	s.logger.Errorln("bad request error", "method", r.Method, "path", r.URL.Path, "error", err.Error())
	helpers.ErrorResponse(w, http.StatusBadRequest, err.Error())
}

func (s *server) notFoundResponse(w http.ResponseWriter, r *http.Request, err error) {
	s.logger.Errorln("not found error", "method", r.Method, "path", r.URL.Path, "error", err.Error())
	helpers.ErrorResponse(w, http.StatusNotFound, "nothing found")
}

func (s *server) unauthorizedResponse(w http.ResponseWriter, r *http.Request, err error) {
	s.logger.Errorln("user is not authorized", "method", r.Method, "path", r.URL.Path, "error", err.Error())
	helpers.ErrorResponse(w, http.StatusUnauthorized, "unauthorized")
}

func (s *server) forbiddenResponse(w http.ResponseWriter, r *http.Request, err error) {
	s.logger.Errorln("forbidden", "method", r.Method, "path", r.URL.Path, "error", err.Error())
	helpers.ErrorResponse(w, http.StatusForbidden, err.Error())
}
