package request

import (
	"strconv"

	"github.com/gin-gonic/gin"
)

func Pagination(c *gin.Context) (start, limit int) {
	start, err := strconv.Atoi(c.Query("start"))
	if err != nil || start < 1 {
		start = 1
	}
	limit, err = strconv.Atoi(c.Query("limit"))
	if err != nil || limit < 0 || limit > 200 {
		limit = 20
	}
	return (start - 1) * limit, limit
}
