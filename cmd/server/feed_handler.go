package main

import (
	"net/http"

	"github.com/sirUnchained/my-go-instagram/internal/helpers"
)

// Feed godoc
//
//	@Summary		Get feed
//	@Description	just get feed
//	@Tags			posts
//	@Accept			json
//	@Produce		json
//	@Param			limit	query		int		false	"number of comments to return (default: 20, max: 100)"
//	@Param			offset	query		int		false	"number of comments to skip (default: 0)"
//	@Success		200	{object}	helpers.DataRes{data=[]models.PostModel}
//	@Failure		401	{object}	helpers.ErrorRes
//	@Failure		500	{object}	helpers.ErrorRes
//	@Security		ApiKeyAuth
//	@Router			/posts/feed [get]
func (s *server) getFeedHandler(w http.ResponseWriter, r *http.Request) {
	limit, offset := helpers.GetLimitOffset(r)
	user := helpers.GetUserFromContext(r)

	ctx := r.Context()
	feed, err := s.postgreStorage.PostStore.GetFeed(ctx, limit, offset, user.Id)
	if err != nil {
		s.internalServerErrorResponse(w, r, err)
		return
	}

	if err := helpers.JsonResponse(w, http.StatusOK, feed); err != nil {
		s.internalServerErrorResponse(w, r, err)
		return
	}
}
