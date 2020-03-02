package gorm

import "testing"

func TestBind(t *testing.T) {
	parameter := Parameter{
		"select":   []string{"id", "name", "phone"},
		"id":       1,
		"sex":      "in (1, 2)",
		"registry": `in ("anhui", "shanghai")`,
		"birthday": `between "1990-01-01" and "2000-01-01"`,
		"weight":   110.5,
		"name":     `like "%王%"`,
		"phone":    "15695655353",
	}
	new(GormDB).build(parameter)
}
