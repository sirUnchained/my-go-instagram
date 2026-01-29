package main

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	global_varables "github.com/sirUnchained/my-go-instagram/internal/global"
	"github.com/sirUnchained/my-go-instagram/internal/helpers"
)

// follow godoc
//
//	@Summary		a user followings
//	@Description	you can see who a user is following
//	@Tags			follows
//	@Accept			json
//	@Produce		json
//	@Param			userid	path		int		true	"post id"
//	@Param			limit	query		int		false	"number of comments to return (default: 20, max: 100)"
//	@Param			offset	query		int		false	"number of comments to skip (default: 0)"
//	@Success		200		{object}	helpers.DataRes{Data=[]models.UserModel}
//	@Failure		400		{object}	helpers.ErrorRes
//	@Failure		403		{object}	helpers.ErrorRes
//	@Failure		404		{object}	helpers.ErrorRes
//	@Failure		500		{object}	helpers.ErrorRes
//	@Security		ApiKeyAuth
//	@Router			/users/{userid}/followings [get]
func (s *server) getFollowingsHandler(w http.ResponseWriter, r *http.Request) {
	userid, err := strconv.ParseInt(chi.URLParam(r, "userid"), 10, 64)
	if err != nil {
		s.badRequestResponse(w, r, err)
		return
	}

	limit, offset := helpers.GetLimitOffset(r)

	ctx := r.Context()
	users, err := s.postgreStorage.FollowStore.GetFollowings(ctx, userid, limit, offset)
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

	if err := helpers.JsonResponse(w, http.StatusOK, users); err != nil {
		s.internalServerErrorResponse(w, r, err)
		return
	}
}

// follow godoc
//
//	@Summary		a user followers
//	@Description	you can see who follows a user
//	@Tags			follows
//	@Accept			json
//	@Produce		json
//	@Param			userid	path		int		true	"post id"
//	@Param			limit	query		int		false	"number of comments to return (default: 20, max: 100)"
//	@Param			offset	query		int		false	"number of comments to skip (default: 0)"
//	@Success		200		{object}	helpers.DataRes{Data=[]models.UserModel}
//	@Failure		400		{object}	helpers.ErrorRes
//	@Failure		403		{object}	helpers.ErrorRes
//	@Failure		404		{object}	helpers.ErrorRes
//	@Failure		500		{object}	helpers.ErrorRes
//	@Security		ApiKeyAuth
//	@Router			/users/{userid}/followers [get]
func (s *server) getFollowersHandler(w http.ResponseWriter, r *http.Request) {
	userid, err := strconv.ParseInt(chi.URLParam(r, "userid"), 10, 64)
	if err != nil {
		s.badRequestResponse(w, r, err)
		return
	}

	limit, offset := helpers.GetLimitOffset(r)

	ctx := r.Context()
	users, err := s.postgreStorage.FollowStore.GetFollowers(ctx, userid, limit, offset)
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

	if err := helpers.JsonResponse(w, http.StatusOK, users); err != nil {
		s.internalServerErrorResponse(w, r, err)
		return
	}
}

// follow godoc
//
//	@Summary		follow a user
//	@Description	follow a user
//	@Tags			follows
//	@Accept			json
//	@Produce		json
//	@Param			userid	path		int	true	"Target User ID"
//	@Success		201	{object}	helpers.DataRes{data=nil}
//	@Failure		400	{object}	helpers.ErrorRes
//	@Failure		403	{object}	helpers.ErrorRes
//	@Failure		404	{object}	helpers.ErrorRes
//	@Failure		500	{object}	helpers.ErrorRes
//	@Security		ApiKeyAuth
//	@Router			/users/{userid}/follow [post]
func (s *server) followUserHandler(w http.ResponseWriter, r *http.Request) {
	user := helpers.GetUserFromContext(r)

	targetUserId, err := strconv.ParseInt(chi.URLParam(r, "targetuserid"), 10, 64)
	if err != nil {
		s.badRequestResponse(w, r, err)
		return
	}

	ctx := r.Context()
	if err := s.postgreStorage.FollowStore.Create(ctx, user.Id, targetUserId); err != nil {
		s.internalServerErrorResponse(w, r, err)
		return
	}

	if err := helpers.JsonResponse(w, http.StatusCreated, nil); err != nil {
		s.internalServerErrorResponse(w, r, err)
		return
	}
}

// follow godoc
//
//	@Summary		unfollow a user
//	@Description	unfollow a user
//	@Tags			follows
//	@Accept			json
//	@Produce		json
//	@Param			userid	path		int	true	"Target User ID"
//	@Success		200	{object}	helpers.DataRes{data=nil}
//	@Failure		400	{object}	helpers.ErrorRes
//	@Failure		403	{object}	helpers.ErrorRes
//	@Failure		404	{object}	helpers.ErrorRes
//	@Failure		500	{object}	helpers.ErrorRes
//	@Security		ApiKeyAuth
//	@Router			/users/{userid}/unfollow [post]
func (s *server) unfollowUserHandler(w http.ResponseWriter, r *http.Request) {
	user := helpers.GetUserFromContext(r)

	targetUserId, err := strconv.ParseInt(chi.URLParam(r, "targetuserid"), 10, 64)
	if err != nil {
		s.badRequestResponse(w, r, err)
		return
	}

	ctx := r.Context()
	if err := s.postgreStorage.FollowStore.Delete(ctx, user.Id, targetUserId); err != nil {
		s.internalServerErrorResponse(w, r, err)
		return
	}

	if err := helpers.JsonResponse(w, http.StatusOK, nil); err != nil {
		s.internalServerErrorResponse(w, r, err)
		return
	}
}
