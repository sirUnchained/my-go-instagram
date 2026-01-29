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

// report godoc
//
//	@Summary		getting reports
//	@Description	admin can get reports
//	@Tags			reports
//	@Accept			json
//	@Produce		json
//	@Param			reportid	path		int		true	"post id"
//	@Param			limit		query		int		false	"number of reports to return (default: 20, max: 100)"
//	@Param			offset		query		int		false	"number of reports to skip (default: 0)"
//	@Success		200		{object}	helpers.DataRes{Data=[]models.ReportModel}
//	@Failure		400		{object}	helpers.ErrorRes
//	@Failure		403		{object}	helpers.ErrorRes
//	@Failure		500		{object}	helpers.ErrorRes
//	@Security		ApiKeyAuth
//	@Router			/reports/{reportid} [get]
func (s *server) getReportsHandler(w http.ResponseWriter, r *http.Request) {
	limit, offset := helpers.GetLimitOffset(r)

	ctx := r.Context()
	reports, err := s.postgreStorage.ReportStore.GetReports(ctx, limit, offset)
	if err != nil {
		s.internalServerErrorResponse(w, r, err)
		return
	}

	if helpers.JsonResponse(w, http.StatusOK, reports); err != nil {
		s.internalServerErrorResponse(w, r, err)
		return
	}
}

// report godoc
//
//	@Summary		create report
//	@Description	just create a report
//	@Tags			reports
//	@Accept			json
//	@Produce		json
//	@Param			payload	body		payloads.CreateReportPayload	true	"report credentials"
//	@Success		200	{object}		helpers.DataRes{Data=nil}
//	@Failure		400	{object}		helpers.ErrorRes
//	@Failure		500	{object}		helpers.ErrorRes
//	@Router			/reports/{reportid} [post]
func (s *server) createReportsHandler(w http.ResponseWriter, r *http.Request) {
	var reportP payloads.CreateReportPayload
	user := helpers.GetUserFromContext(r)

	if err := helpers.ReadJson(w, r, &reportP); err != nil {
		s.badRequestResponse(w, r, err)
		return
	}

	v := validator.New(validator.WithRequiredStructEnabled())
	if err := v.Struct(reportP); err != nil {
		s.badRequestResponse(w, r, err)
		return
	}

	ctx := r.Context()
	if err := s.postgreStorage.ReportStore.Create(ctx, user.Id, reportP); err != nil {
		s.internalServerErrorResponse(w, r, err)
		return
	}

	if err := helpers.JsonResponse(w, http.StatusCreated, nil); err != nil {
		s.internalServerErrorResponse(w, r, err)
		return
	}
}

// report godoc
//
//	@Summary		delete a report
//	@Description	delete a report by ID
//	@Tags			reports
//	@Accept			json
//	@Produce		json
//	@Param			reportid	path		int	true	"post ID"
//	@Success		200			{object}	helpers.DataRes{data=nil}
//	@Failure		404			{object}	helpers.ErrorRes
//	@Failure		500			{object}	helpers.ErrorRes
//	@Security		ApiKeyAuth
//	@Router			/reports/{reportid} [delete]
func (s *server) deleteReportsHandler(w http.ResponseWriter, r *http.Request) {
	reportid, err := strconv.ParseInt(chi.URLParam(r, "reportid"), 10, 64)
	if err != nil {
		s.badRequestResponse(w, r, err)
		return
	}

	ctx := r.Context()
	err = s.postgreStorage.ReportStore.Delete(ctx, reportid)
	if err != nil {
		switch {
		case errors.Is(err, global_varables.NOT_FOUND_ROW):
			s.notFoundResponse(w, r, fmt.Errorf("no report with this id exists"))
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
