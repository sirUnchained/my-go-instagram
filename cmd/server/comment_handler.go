package main

import (
	"net/http"

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
		s.internalServerErrorResponse(w, r, err)
		return
	}

	if err := helpers.JsonResponse(w, http.StatusCreated, nil); err != nil {
		s.internalServerErrorResponse(w, r, err)
		return
	}
}
