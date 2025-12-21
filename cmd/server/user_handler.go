package main

import (
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/sirUnchained/my-go-instagram/internal/payloads"
	"github.com/sirUnchained/my-go-instagram/internal/scripts"
)

func (s *server) createUserHandler(w http.ResponseWriter, r *http.Request) {
	var userP payloads.UserPayload
	if err := scripts.ReadJson(w, r, &userP); err != nil {
		scripts.ErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	v := validator.New(validator.WithRequiredStructEnabled())
	if err := v.Struct(userP); err != nil {
		scripts.ErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	ctx := r.Context()
	if err := s.postgreStorage.UserStore.Create(ctx, &userP); err != nil {
		s.logger.Errorln(err)
		scripts.ErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	if err := scripts.ErrorResponse(w, http.StatusCreated, "user created"); err != nil {
		scripts.ErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

}
