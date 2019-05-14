package request

import (
	"net/http"
	"strconv"
)

func Pagination(r *http.Request) (start, limit int) {
	start, err := strconv.Atoi(r.FormValue("start"))
	if err != nil || start < 1 {
		start = 1
	}
	limit, err = strconv.Atoi(r.FormValue("limit"))
	if err != nil || limit < 0 || limit > 200 {
		limit = 20
	}
	return (start - 1) * limit, limit
}

func Page(limit int, count int) int {
	return (count + limit - 1) / limit
}
