package request

import (
	"net/http"
	"strconv"
)

func Pagination(r *http.Request) (int, int) {
	offset, err := strconv.Atoi(r.FormValue("offset"))
	if err != nil || offset < 1 {
		offset = 1
	}
	limit, err := strconv.Atoi(r.FormValue("limit"))
	if err != nil || limit < 0 || limit > 200 {
		limit = 20
	}
	return (offset - 1) * limit, limit
}

func PageTotal(limit int, count int) int {
	return (count + limit - 1) / limit
}
