package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt/v5"
	global_varables "github.com/sirUnchained/my-go-instagram/internal/global"
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
	user, err := s.postgreStorage.UserStore.GetById(ctx, userid)
	if err != nil {
		switch {
		case errors.Is(err, global_varables.NOT_FOUND_ROW):
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

func (s *server) getMeHandler(w http.ResponseWriter, r *http.Request) {
	user := scripts.GetUserFromContext(r)

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
	user, err := s.postgreStorage.UserStore.Create(ctx, &userP)
	if err != nil {
		switch {
		case errors.Is(err, global_varables.USERNAME_DUP) || errors.Is(err, global_varables.EMAIL_DUP):
			s.badRequestResponse(w, r, fmt.Errorf("user already exists"))
			return
		default:
			s.internalServerErrorResponse(w, r, err)
			return
		}
	}

	now := time.Now()
	claims := jwt.RegisteredClaims{
		Subject:   fmt.Sprintf("%d", user.Id),
		ExpiresAt: jwt.NewNumericDate(now.Add(s.serverConfigs.auth.expMin)),
		IssuedAt:  jwt.NewNumericDate(now),
		NotBefore: jwt.NewNumericDate(now),
		Issuer:    s.serverConfigs.auth.iss,
		Audience:  jwt.ClaimStrings{s.serverConfigs.auth.aud},
	}

	token, err := s.auth.GenerateToken(claims)
	if err != nil {
		s.internalServerErrorResponse(w, r, err)
		return
	}

	if err := scripts.JsonResponse(w, http.StatusCreated, token); err != nil {
		s.internalServerErrorResponse(w, r, err)
		return
	}

}
