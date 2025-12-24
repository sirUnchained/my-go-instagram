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

// GetUser godoc
//
//	@Summary		get single user
//	@Description	get one user by ID
//	@Tags			users
//	@Accept			json
//	@Produce		json
//	@Param			userid	path		int	true	"User ID"
//	@Success		200	{object}	models.UserModel
//	@Failure		400	{object}	error
//	@Failure		404	{object}	error
//	@Failure		500	{object}	error
//	@Security		ApiKeyAuth
//	@Router			/users/{userid} [get]
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

// GetMe godoc
//
//	@Summary		get the user in the token
//	@Description	usign the client token we'll give the user which use token
//	@Tags			users
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	models.UserModel
//	@Failure		400	{object}	error
//	@Failure		404	{object}	error
//	@Failure		500	{object}	error
//	@Security		ApiKeyAuth
//	@Router			/users/me [get]
func (s *server) getMeHandler(w http.ResponseWriter, r *http.Request) {
	user := scripts.GetUserFromContext(r)

	if err := scripts.JsonResponse(w, http.StatusOK, user); err != nil {
		s.internalServerErrorResponse(w, r, err)
		return
	}
}

// CreateUser godoc
//
//	@Summary		create user
//	@Description	we can register new user and return a token
//	@Tags			users
//	@Accept			json
//	@Produce		json
//	@Param			payload	body		payloads.CreateUserPayload	true	"User credentials"
//	@Success		200	{object}	payloads.UserWithToken
//	@Failure		400	{object}	error
//	@Failure		500	{object}	error
//	@Router			/users/new [post]
func (s *server) createUserHandler(w http.ResponseWriter, r *http.Request) {
	var userP payloads.CreateUserPayload
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

	var userToken payloads.UserWithToken = payloads.UserWithToken{User: *user, Token: token}

	if err := scripts.JsonResponse(w, http.StatusCreated, userToken); err != nil {
		s.internalServerErrorResponse(w, r, err)
		return
	}

}

// Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJzZXJ2ZXIiLCJzdWIiOiIxMyIsImF1ZCI6WyJzaXJ1bmNoYWluZWQtaW5zdGFncmFtIl0sImV4cCI6MTc2NjU3Mjc0MCwibmJmIjoxNzY2NDg2MzQwLCJpYXQiOjE3NjY0ODYzNDB9.HQtCTzDdkCHaGhddDXbVB2wQ5WLlgQ_zLDkxEuOQ2Ek
