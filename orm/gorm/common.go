package gorm

import (
	"fmt"
	"strings"

	"github.com/jinzhu/gorm"
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

type GormDao struct {
	*gorm.DB
}

func (dao *GormDao) Search(resultsPoint interface{}, parameter Parameter) error {
	return dao.bind(parameter).Find(resultsPoint).Error
}

func (dao *GormDao) bind(parameter Parameter) *gorm.DB {
	var (
		result = new(summary)
		wheres = make([]string, 0)
	)

	for key, val := range parameter {
		key = strings.ToLower(strings.TrimSpace(key))
		switch key {
		case "select":
			if _, ok := val.([]string); ok {
				result.query = join(val.([]string))
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
			case uint, uint32, uint64, int, int32, int64, float32, float64:
				wheres = append(wheres, fmt.Sprintf(`(%s = %v)`, key, val))
			case string:
				if strings.Contains(val.(string), " ") {
					wheres = append(wheres, fmt.Sprintf(`(%s %v)`, key, val))
					continue
				}
				wheres = append(wheres, fmt.Sprintf(`(%s = "%v")`, key, val))
			}
		}
	}
	result.where = strings.Join(wheres, " AND ")

	db := dao.DB

	if len(result.query) != 0 {
		db = db.Select(result.query)
	}
	if len(result.where) != 0 {
		db = db.Where(result.where)
	}
	if result.limit != 0 {
		if result.offset != 0 {
			db = db.Offset(result.offset)
		}
		db = db.Limit(result.limit)
	}
	if len(result.groupby) != 0 {
		db = db.Group(result.groupby)
	}
	if len(result.orderby) != 0 {
		db = db.Order(result.orderby)
	}

	return db
}

func join(slice []string) string {
	s := make([]string, 0)
	for _, val := range slice {
		if len(val) != 0 {
			s = append(s, val)
		}
	}
	return strings.Join(s, ", ")
}
