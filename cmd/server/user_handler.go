package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	global_varables "github.com/sirUnchained/my-go-instagram/internal/global"
	"github.com/sirUnchained/my-go-instagram/internal/helpers"
	"github.com/sirUnchained/my-go-instagram/internal/payloads"
)

// GetUser godoc
//
//	@Summary		get single user
//	@Description	get one user by ID
//	@Tags			users
//	@Accept			json
//	@Produce		json
//	@Param			userid	path		int	true	"User ID"
//	@Success		200	{object}	helpers.DataRes{data=models.UserModel}
//	@Failure		400	{object}	helpers.ErrorRes
//	@Failure		403	{object}	helpers.ErrorRes
//	@Failure		404	{object}	helpers.ErrorRes
//	@Failure		500	{object}	helpers.ErrorRes
//	@Security		ApiKeyAuth
//	@Router			/users/{userid} [get]
func (s *server) getUserHandler(w http.ResponseWriter, r *http.Request) {
	targetUserId, err := strconv.ParseInt(chi.URLParam(r, "userid"), 10, 64)
	if err != nil {
		s.badRequestResponse(w, r, fmt.Errorf("invalid id"))
		return
	}

	ctx := r.Context()
	targetUser, err := s.postgreStorage.UserStore.GetById(ctx, targetUserId)
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

	if err := helpers.JsonResponse(w, http.StatusOK, targetUser); err != nil {
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
//	@Success		200	{object}	helpers.DataRes{data=models.UserModel}
//	@Failure		400	{object}	helpers.ErrorRes
//	@Failure		404	{object}	helpers.ErrorRes
//	@Failure		500	{object}	helpers.ErrorRes
//	@Security		ApiKeyAuth
//	@Router			/users/me [get]
func (s *server) getMeHandler(w http.ResponseWriter, r *http.Request) {
	user := helpers.GetUserFromContext(r)

	if err := helpers.JsonResponse(w, http.StatusOK, user); err != nil {
		s.internalServerErrorResponse(w, r, err)
		return
	}
}

// UpdateUser godoc
//
//	@Summary		update the user in the token
//	@Description	with client token we'll update the user which use token
//	@Tags			users
//	@Accept			json
//	@Produce		json
//	@Param			payload	body		payloads.CreateUserPayload	true	"User credentials"
//	@Success		200	{object}	helpers.DataRes{data=models.UserModel}
//	@Failure		400	{object}	helpers.ErrorRes
//	@Failure		404	{object}	helpers.ErrorRes
//	@Failure		500	{object}	helpers.ErrorRes
//	@Security		ApiKeyAuth
//	@Router			/users/update [put]
func (s *server) updateMeHandler(w http.ResponseWriter, r *http.Request) {
	user := helpers.GetUserFromContext(r)
	var userP payloads.CreateUserPayload
	if err := helpers.ReadJson(w, r, &userP); err != nil {
		s.badRequestResponse(w, r, err)
		return
	}

	v := validator.New(validator.WithRequiredStructEnabled())
	if err := v.Struct(userP); err != nil {
		s.badRequestResponse(w, r, err)
		return
	}

	ctx := r.Context()
	user, err := s.postgreStorage.UserStore.UpdateData(ctx, user, &userP)
	if err != nil {
		switch {
		case errors.Is(err, global_varables.USERNAME_DUP):
			s.badRequestResponse(w, r, fmt.Errorf("you are not allowed to use this username"))
			return
		case errors.Is(err, global_varables.EMAIL_DUP):
			s.badRequestResponse(w, r, fmt.Errorf("you are not allowed to use this email"))
			return
		default:
			s.internalServerErrorResponse(w, r, err)
			return
		}
	}

	if err := helpers.JsonResponse(w, http.StatusOK, user); err != nil {
		s.internalServerErrorResponse(w, r, err)
		return
	}
}
