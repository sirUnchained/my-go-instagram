package main

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt/v5"
	global_varables "github.com/sirUnchained/my-go-instagram/internal/global"
	"github.com/sirUnchained/my-go-instagram/internal/payloads"
	"github.com/sirUnchained/my-go-instagram/internal/scripts"
)

// CreateUser godoc
//
//	@Summary		create user
//	@Description	we can register new user and return a token
//	@Tags			auth
//	@Accept			json
//	@Produce		json
//	@Param			payload	body		payloads.CreateUserPayload	true	"User credentials"
//	@Success		200	{object}	payloads.UserWithToken
//	@Failure		400	{object}	error
//	@Failure		500	{object}	error
//	@Router			/auth/register [post]
func (s *server) registerUserHandler(w http.ResponseWriter, r *http.Request) {
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

// CreateUser godoc
//
//	@Summary		login user
//	@Description	we can login user and return a token
//	@Tags			auth
//	@Accept			json
//	@Produce		json
//	@Param			payload	body		payloads.LoginUserPayload	true	"User credentials"
//	@Success		200	{object}	payloads.UserWithToken
//	@Failure		400	{object}	error
//	@Failure		500	{object}	error
//	@Router			/auth/login [post]
func (s *server) loginUserHandler(w http.ResponseWriter, r *http.Request) {
	var loginP payloads.LoginUserPayload
	if err := scripts.ReadJson(w, r, &loginP); err != nil {
		s.badRequestResponse(w, r, err)
		return
	}

	v := validator.New(validator.WithRequiredStructEnabled())
	if err := v.Struct(loginP); err != nil {
		s.badRequestResponse(w, r, err)
		return
	}

	ctx := r.Context()
	user, err := s.postgreStorage.UserStore.GetByEmail(ctx, loginP.Email)
	if err != nil {
		switch {
		case errors.Is(err, global_varables.NOT_FOUND_ROW):
			s.badRequestResponse(w, r, fmt.Errorf("invalid email or password"))
			return
		default:
			s.badRequestResponse(w, r, err)
			return
		}
	}

	if err := user.Password.Compare(loginP.Password); err != nil {
		s.badRequestResponse(w, r, fmt.Errorf("invalid email or password"))
		return
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

	var result payloads.UserWithToken
	result.Token = token
	result.User = *user

	if err := scripts.JsonResponse(w, http.StatusOK, result); err != nil {
		s.internalServerErrorResponse(w, r, err)
		return
	}

}
