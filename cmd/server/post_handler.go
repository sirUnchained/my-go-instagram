package main

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/sirUnchained/my-go-instagram/internal/helpers"
	"github.com/sirUnchained/my-go-instagram/internal/payloads"
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
//	@Success		201		{object}	helpers.DataRes{data=models.PostModel}
//	@Failure		400		{object}	helpers.ErrorRes
//	@Failure		401		{object}	helpers.ErrorRes
//	@Failure		403		{object}	helpers.ErrorRes
//	@Failure		413		{object}	helpers.ErrorRes
//	@Failure		500		{object}	helpers.ErrorRes
//	@Security		ApiKeyAuth
//	@Router			/posts/new [post]
func (s *server) createPostHandler(w http.ResponseWriter, r *http.Request) {
	var postP payloads.CreatePostPayload
	user := helpers.GetUserFromContext(r)

	status, err := helpers.ReadFormFiles(w, r, user.Id, &postP)
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

	if err := helpers.JsonResponse(w, http.StatusCreated, post); err != nil {
		s.internalServerErrorResponse(w, r, err)
		return
	}
}

// GetPost godoc
//
//	@Summary		get single post
//	@Description	get one post by ID
//	@Tags			posts
//	@Accept			json
//	@Produce		json
//	@Param			postid	path		int	true	"post ID"
//	@Success		200	{object}	helpers.DataRes{data=models.PostModel}
//	@Failure		400	{object}	helpers.ErrorRes
//	@Failure		404	{object}	helpers.ErrorRes
//	@Failure		500	{object}	helpers.ErrorRes
//	@Security		ApiKeyAuth
//	@Router			/posts/{postid} [get]
func (s *server) getPostHandler(w http.ResponseWriter, r *http.Request) {
	postid, err := strconv.ParseInt(chi.URLParam(r, "postid"), 10, 64)
	if err != nil {
		s.badRequestResponse(w, r, fmt.Errorf("invalid postid"))
		return
	}

	ctx := r.Context()
	post, err := s.postgreStorage.PostStore.GetById(ctx, postid)
	if err != nil {
		s.internalServerErrorResponse(w, r, err)
		return
	}

	if err := helpers.JsonResponse(w, http.StatusOK, post); err != nil {
		s.internalServerErrorResponse(w, r, err)
		return
	}
}
