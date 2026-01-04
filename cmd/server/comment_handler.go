package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	global_varables "github.com/sirUnchained/my-go-instagram/internal/global"
	"github.com/sirUnchained/my-go-instagram/internal/helpers"
	"github.com/sirUnchained/my-go-instagram/internal/payloads"
)

// CreateComment godoc
//
//	@Summary		create a new comment
//	@Description	wyou can create comment for a post or just reply to one
//	@Tags			comments
//	@Accept			json
//	@Produce		json
//	@Param			payload	body		payloads.CreateCommentPayload	true	"comments credentials"
//	@Success		200	{object}	nil
//	@Failure		400	{object}	helpers.ErrorRes
//	@Failure		404	{object}	helpers.ErrorRes
//	@Failure		500	{object}	helpers.ErrorRes
//	@Security		ApiKeyAuth
//	@Router			/comments/new [post]
func (s *server) createCommentHandler(w http.ResponseWriter, r *http.Request) {
	commentP := payloads.CreateCommentPayload{}
	if err := helpers.ReadJson(w, r, &commentP); err != nil {
		s.badRequestResponse(w, r, err)
		return
	}

	user := helpers.GetUserFromContext(r)
	commentP.CreatorID = user.Id

	ctx := r.Context()
	if err := s.postgreStorage.CommentStore.Create(ctx, user.Id, &commentP); err != nil {
		switch {
		case errors.Is(err, global_varables.NOT_FOUND_ROW):
			s.notFoundResponse(w, r, fmt.Errorf("no such a post exists or maybe parrent comment not found"))
			return
		default:
			s.internalServerErrorResponse(w, r, err)
			return
		}
	}

	if err := helpers.JsonResponse(w, http.StatusCreated, nil); err != nil {
		s.internalServerErrorResponse(w, r, err)
		return
	}
}

// GetComment godoc
//
//	@Summary		get post comments
//	@Description	you can get comments for a post with pagination
//	@Tags			comments
//	@Accept			json
//	@Produce		json
//	@Param			postid	path		int		true	"post id"
//	@Param			limit	query		int		false	"number of comments to return (default: 20, max: 100)"
//	@Param			offset	query		int		false	"number of comments to skip (default: 0)"
//	@Success		200		{object}	helpers.DataRes{Data=[]models.CommentModel}
//	@Failure		400		{object}	helpers.ErrorRes
//	@Failure		404		{object}	helpers.ErrorRes
//	@Failure		500		{object}	helpers.ErrorRes
//	@Security		ApiKeyAuth
//	@Router			/comments/posts/{postid} [get]
func (s *server) getCommentsHandler(w http.ResponseWriter, r *http.Request) {
	postid, err := strconv.ParseInt(chi.URLParam(r, "postid"), 10, 64)
	if err != nil {
		s.badRequestResponse(w, r, err)
		return
	}

	limit, offset := helpers.GetLimitOffset(r)

	ctx := r.Context()
	comments, err := s.postgreStorage.CommentStore.GetPostComments(ctx, postid, limit, offset)
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

	if err := helpers.JsonResponse(w, http.StatusOK, comments); err != nil {
		s.internalServerErrorResponse(w, r, err)
		return
	}

}

// GetReplyComment godoc
//
//	@Summary		get replies to a comments
//	@Description	you can get replies to a comments with a pagination
//	@Tags			comments
//	@Accept			json
//	@Produce		json
//	@Param			commentid	path		int		true	"comment id"
//	@Param			limit		query		int		false	"number of comments to return (default: 20, max: 100)"
//	@Param			offset		query		int		false	"number of comments to skip (default: 0)"
//	@Success		200			{object}	helpers.DataRes{Data=[]models.CommentModel}
//	@Failure		400			{object}	helpers.ErrorRes
//	@Failure		404			{object}	helpers.ErrorRes
//	@Failure		500			{object}	helpers.ErrorRes
//	@Security		ApiKeyAuth
//	@Router			/comments/{commentid}/replies [get]
func (s *server) getReplyCommentsHandler(w http.ResponseWriter, r *http.Request) {
	commentid, err := strconv.ParseInt(chi.URLParam(r, "commentid"), 10, 64)
	if err != nil {
		s.badRequestResponse(w, r, err)
		return
	}

	limit, offset := helpers.GetLimitOffset(r)

	ctx := r.Context()
	comments, err := s.postgreStorage.CommentStore.GetRepliedComments(ctx, commentid, limit, offset)
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

	if err := helpers.JsonResponse(w, http.StatusOK, comments); err != nil {
		s.internalServerErrorResponse(w, r, err)
		return
	}

}

// DeleteComment godoc
//
//	@Summary		delete a comment
//	@Description	you can get dalete a comment or a replied comment
//	@Tags			comments
//	@Accept			json
//	@Produce		json
//	@Param			commentid	path		int		true	"comment id"
//	@Success		200			{object}	nil
//	@Failure		400			{object}	helpers.ErrorRes
//	@Failure		404			{object}	helpers.ErrorRes
//	@Failure		500			{object}	helpers.ErrorRes
//	@Security		ApiKeyAuth
//	@Router			/comments/{commentid} [delete]
func (s *server) deleteCommentHandler(w http.ResponseWriter, r *http.Request) {
	commentid, err := strconv.ParseInt(chi.URLParam(r, "commentid"), 10, 64)
	if err != nil {
		s.badRequestResponse(w, r, err)
		return
	}

	ctx := r.Context()
	err = s.postgreStorage.CommentStore.Delete(ctx, commentid)
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

	if err := helpers.JsonResponse(w, http.StatusOK, nil); err != nil {
		s.internalServerErrorResponse(w, r, err)
		return
	}
}
