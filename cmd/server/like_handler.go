package main

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	global_varables "github.com/sirUnchained/my-go-instagram/internal/global"
	"github.com/sirUnchained/my-go-instagram/internal/helpers"
)

// LikePost godoc
//
//	@Summary		like a post
//	@Description	you can like a post
//	@Tags			likes
//	@Accept			json
//	@Produce		json
//	@Param			postid	path		int	true	"post ID"
//	@Success		201	{object}	helpers.DataRes{data=nil}
//	@Failure		400	{object}	helpers.ErrorRes
//	@Failure		404	{object}	helpers.ErrorRes
//	@Failure		500	{object}	helpers.ErrorRes
//	@Security		ApiKeyAuth
//	@Router			/posts/{postid}/like [post]
func (s *server) likePostHandler(w http.ResponseWriter, r *http.Request) {
	user := helpers.GetUserFromContext(r)

	postid, err := strconv.ParseInt(chi.URLParam(r, "postid"), 10, 64)
	if err != nil {
		s.badRequestResponse(w, r, err)
		return
	}

	ctx := r.Context()
	if err := s.postgreStorage.LikeStore.Create(ctx, postid, user.Id); err != nil {
		if !errors.Is(err, global_varables.DUP_ITEM) {
			s.internalServerErrorResponse(w, r, err)
			return
		}
	}

	if err := helpers.JsonResponse(w, http.StatusCreated, nil); err != nil {
		s.internalServerErrorResponse(w, r, err)
		return
	}
}

// LikePost godoc
//
//	@Summary		dislike a post
//	@Description	you can dislike a post
//	@Tags			likes
//	@Accept			json
//	@Produce		json
//	@Param			postid	path		int	true	"post ID"
//	@Success		200	{object}	helpers.DataRes{data=nil}
//	@Failure		400	{object}	helpers.ErrorRes
//	@Failure		404	{object}	helpers.ErrorRes
//	@Failure		500	{object}	helpers.ErrorRes
//	@Security		ApiKeyAuth
//	@Router			/posts/{postid}/dislike [post]
func (s *server) dislikePostHandler(w http.ResponseWriter, r *http.Request) {
	user := helpers.GetUserFromContext(r)

	postid, err := strconv.ParseInt(chi.URLParam(r, "postid"), 10, 64)
	if err != nil {
		s.badRequestResponse(w, r, err)
		return
	}

	ctx := r.Context()
	if err := s.postgreStorage.LikeStore.Delete(ctx, postid, user.Id); err != nil {
		s.internalServerErrorResponse(w, r, err)
		return
	}

	if err := helpers.JsonResponse(w, http.StatusOK, nil); err != nil {
		s.internalServerErrorResponse(w, r, err)
		return
	}
}
