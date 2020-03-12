package gorm

import (
	"testing"
)

func TestBind(t *testing.T) {
	param := map[string]interface{}{
		"select":  []string{"id", "name"},
		"id":      "in (1, 2, 3)",
		"name":    `like "%张%"`,
		"sex":     "<> 3",
		"height":  175,
		"groupby": "name",
		"orderby": "id DESC",
	}
	new(GormDao).bind(param)
}
