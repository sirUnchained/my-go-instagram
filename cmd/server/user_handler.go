package main

import (
	"net/http"

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
	targetUser := scripts.GetUserByIdFromContext(r)

	if err := scripts.JsonResponse(w, http.StatusOK, targetUser); err != nil {
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
