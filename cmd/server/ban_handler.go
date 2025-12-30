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

// BanUser godoc
//
//	@Summary		ban the user
//	@Description	with incoming id we ban the user (only adins can do it).
//	@Tags			users
//	@Accept			json
//	@Produce		json
//	@Param			payload	body		payloads.CreateBanPayload	true	"User credentials"
//	@Param			userid	path		int	true	"User ID"
//	@Success		200	{object}	nil
//	@Failure		400	{object}	helpers.ErrorRes
//	@Failure		404	{object}	helpers.ErrorRes
//	@Failure		500	{object}	helpers.ErrorRes
//	@Security		ApiKeyAuth
//	@Router			/users/ban/{userid} [delete]
func (s *server) banUserHandler(w http.ResponseWriter, r *http.Request) {
	banP := &payloads.CreateBanPayload{}
	if err := helpers.ReadJson(w, r, banP); err != nil {
		s.badRequestResponse(w, r, err)
		return
	}

	userid, err := strconv.ParseInt(chi.URLParam(r, "userid"), 10, 64)
	if err != nil {
		s.badRequestResponse(w, r, fmt.Errorf("invalid userid"))
		return
	}

	ctx := r.Context()
	user, err := s.postgreStorage.UserStore.GetById(ctx, userid)
	if err != nil {
		switch {
		case errors.Is(err, global_varables.USERNAME_DUP):
			s.badRequestResponse(w, r, fmt.Errorf("you are not allowed to use this username"))
			return
		case errors.Is(err, global_varables.EMAIL_DUP):
			s.badRequestResponse(w, r, fmt.Errorf("you are not allowed to use this email"))
			return
		case errors.Is(err, global_varables.NOT_FOUND_ROW):
			s.notFoundResponse(w, r, fmt.Errorf("no user with this email found"))
			return
		default:
			s.internalServerErrorResponse(w, r, err)
			return
		}
	}

	if user.Role.Name == global_varables.ADMIN_ROLE {
		s.badRequestResponse(w, r, fmt.Errorf("admins cannot be banned"))
		return
	}

	if err := s.postgreStorage.BanStore.Create(ctx, user, banP); err != nil {
		switch {
		case errors.Is(err, global_varables.NOT_FOUND_ROW):
			s.notFoundResponse(w, r, fmt.Errorf("no user with this email found"))
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

// BanUser godoc
//
//	@Summary		ban the user
//	@Description	with incoming id we ban the user (only adins can do it).
//	@Tags			users
//	@Accept			json
//	@Produce		json
//	@Param			payload	body		payloads.UnbanPayload	true	"User credentials"
//	@Success		200	{object}	nil
//	@Failure		400	{object}	helpers.ErrorRes
//	@Failure		404	{object}	helpers.ErrorRes
//	@Failure		500	{object}	helpers.ErrorRes
//	@Security		ApiKeyAuth
//	@Router			/users/unban/{userid} [delete]
func (s *server) unbanUserHandler(w http.ResponseWriter, r *http.Request) {
	unbanP := &payloads.UnbanPayload{}
	if err := helpers.ReadJson(w, r, unbanP); err != nil {
		s.badRequestResponse(w, r, err)
		return
	}

	ctx := r.Context()
	if err := s.postgreStorage.BanStore.Delete(ctx, unbanP.Email); err != nil {
		switch {
		case errors.Is(err, global_varables.NOT_FOUND_ROW):
			s.notFoundResponse(w, r, fmt.Errorf("no email found in bans"))
		default:
			s.internalServerErrorResponse(w, r, err)
			return
		}
		return
	}

	if err := helpers.JsonResponse(w, http.StatusOK, nil); err != nil {
		s.internalServerErrorResponse(w, r, err)
		return
	}

}
