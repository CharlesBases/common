package gorm

import (
	"fmt"
	"strings"
)

type Parameter map[string]interface{}

type summary struct {
	query   string
	where   string
	orderby string
	groupby string
	limit   int
	offset  int
}

func Bind(parameter Parameter) *summary {
	var (
		result = new(summary)
		wheres = make([]string, 0)
	)

	for key, val := range parameter {
		key = strings.ToLower(strings.TrimSpace(key))
		switch key {
		case "select":
			if _, ok := val.([]string); ok {
				result.query = strings.Join(val.([]string), ", ")
			}
		case "orderby":
			if _, ok := val.(string); ok {
				result.orderby = val.(string)
			}
		case "groupby":
			if _, ok := val.(string); ok {
				result.groupby = val.(string)
			}
		case "limit":
			if _, ok := val.(int); ok {
				result.limit = val.(int)
			}
		case "offset":
			if _, ok := val.(int); ok {
				result.offset = val.(int)
			}
		default:
			switch val.(type) {
			case int64, float64:
				wheres = append(wheres, fmt.Sprintf(`%s = %v`, key, val))
			case string:
				if strings.Contains(val.(string), " ") {
					wheres = append(wheres, fmt.Sprintf(`%s %v`, key, val))
					continue
				}
				wheres = append(wheres, fmt.Sprintf(`%s = "%v"`, key, val))
			}
		}
	}
	result.where = strings.Join(wheres, " AND ")

	return result
}
