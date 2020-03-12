package gorm

import (
	"fmt"
	"testing"
)

func TestBind(t *testing.T) {
	param := map[string]interface{}{
		"select":  []string{"id", "name"},
		"id":      10086,
		"name":    `like "%中国%"`,
		"status":  "<> 3",
		"orderby": "id DESC",
		"groupby": "name",
		"limit":   1,
		"offset":  1,
	}
	summary := Bind(param)
	fmt.Println(summary)
}
