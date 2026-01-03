package helpers

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

func GetLimitOffset(r *http.Request) (int64, int64) {
	limit, err := strconv.ParseInt(chi.URLParam(r, "limit"), 10, 64)
	if err != nil || limit < 0 || limit > 100 {
		limit = 20
	}

	offset, err := strconv.ParseInt(chi.URLParam(r, "offset"), 10, 64)
	if err != nil || offset < 0 {
		offset = 0
	}

	return limit, offset
}
