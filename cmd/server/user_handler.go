package main

import (
	"errors"
	"net/http"

	"github.com/go-playground/validator/v10"
	error_messages "github.com/sirUnchained/my-go-instagram/internal/errors"
	"github.com/sirUnchained/my-go-instagram/internal/payloads"
	"github.com/sirUnchained/my-go-instagram/internal/scripts"
)

func (s *server) createUserHandler(w http.ResponseWriter, r *http.Request) {
	var userP payloads.UserPayload
	if err := scripts.ReadJson(w, r, &userP); err != nil {
		s.badRequestResponse(w, r, err)
		return
	}

	v := validator.New(validator.WithRequiredStructEnabled())
	if err := v.Struct(userP); err != nil {
		s.badRequestResponse(w, r, err)
		return
	}

	ctx := r.Context()
	if err := s.postgreStorage.UserStore.Create(ctx, &userP); err != nil {
		switch {
		case errors.Is(err, error_messages.USERNAME_DUP) || errors.Is(err, error_messages.EMAIL_DUP):
			s.badRequestResponse(w, r, err)
		default:
			s.internalServerErrorResponse(w, r, err)
		}
		return
	}

	if err := scripts.ErrorResponse(w, http.StatusCreated, "user created"); err != nil {
		s.internalServerErrorResponse(w, r, err)
		return
	}

}
