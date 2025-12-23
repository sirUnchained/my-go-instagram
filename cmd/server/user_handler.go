package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	error_messages "github.com/sirUnchained/my-go-instagram/internal/errors"
	"github.com/sirUnchained/my-go-instagram/internal/payloads"
	"github.com/sirUnchained/my-go-instagram/internal/scripts"
)

func (s *server) getUserHandler(w http.ResponseWriter, r *http.Request) {
	userid, err := strconv.ParseInt(chi.URLParam(r, "userid"), 10, 64)
	if err != nil {
		s.badRequestResponse(w, r, fmt.Errorf("invalid id"))
		return
	}

	ctx := r.Context()
	user, err := s.postgreStorage.UserStore.Get(ctx, userid)
	if err != nil {
		switch {
		case errors.Is(err, error_messages.NOT_FOUND_ROW):
			s.notFoundResponse(w, r, err)
			return
		default:
			s.internalServerErrorResponse(w, r, err)
			return
		}
	}

	if err := scripts.JsonResponse(w, http.StatusOK, user); err != nil {
		s.internalServerErrorResponse(w, r, err)
		return
	}
}

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
