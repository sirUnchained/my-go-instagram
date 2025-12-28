package main

import (
	"net/http"

	"github.com/sirUnchained/my-go-instagram/internal/payloads"
	"github.com/sirUnchained/my-go-instagram/internal/scripts"
	"github.com/sirUnchained/my-go-instagram/internal/storage/models"
)

// CreatePost godoc
//
//	@Summary		create a post with files
//	@Description	create a post with multiple files (images/documents) using form-data
//	@Tags			posts
//	@Accept			multipart/form-data
//	@Produce		json
//	@Param			description	formData	string	false	"Post description"
//	@Param			tags		formData	string	false	"Tags as JSON array e.g. [\"tag1\",\"tag2\"]"
//	@Param			files		formData	[]file	true	"Files to upload (1-5 files, max 10MB each)"
//	@Success		201		{object}	models.PostModel
//	@Failure		400		{object}	map[string]interface{}
//	@Failure		401		{object}	map[string]interface{}
//	@Failure		413		{object}	map[string]interface{}
//	@Failure		500		{object}	map[string]interface{}
//	@Security		ApiKeyAuth
//	@Router			/posts/new [post]
func (s *server) createPostHandler(w http.ResponseWriter, r *http.Request) {
	var postP payloads.CreatePostPayload
	user := scripts.GetUserFromContext(r)

	status, err := scripts.ReadFormFiles(w, r, user.Id, &postP)
	switch status {
	case http.StatusBadRequest:
		s.badRequestResponse(w, r, err)
		return
	case http.StatusInternalServerError:
		s.internalServerErrorResponse(w, r, err)
		return
	}

	// set files to db
	ctx := r.Context()
	files, err := s.postgreStorage.FileStore.Create(ctx, user.Id, postP.Files)
	if err != nil {
		s.internalServerErrorResponse(w, r, err)
		return
	}

	// if we have tags, save them in db
	var tags []models.TagModel
	if len(postP.Tags) > 0 {
		ctx := r.Context()
		tags, err = s.postgreStorage.TagStore.Create(ctx, user.Id, postP.Tags)
		if err != nil {
			s.internalServerErrorResponse(w, r, err)
			return
		}
	}

	// save post
	ctx = r.Context()
	post, err := s.postgreStorage.PostStore.Create(ctx, &postP, &files, &tags, user)
	if err != nil {
		s.internalServerErrorResponse(w, r, err)
		return
	}

	post.Files = files
	post.Tags = tags

	if err := scripts.JsonResponse(w, http.StatusCreated, post); err != nil {
		s.internalServerErrorResponse(w, r, err)
		return
	}
}
