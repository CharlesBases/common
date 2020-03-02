package gorm

import (
	"fmt"
	"strings"

	"github.com/jinzhu/gorm"
)

const (
	query   = "select"
	orderby = "orderby"
	groupby = "groupby"
	offset  = "offset"
	limit   = "limit"
)

type (
	Parameter map[string]interface{}
)

type GormDB struct {
	DB *gorm.DB
}

type bind struct {
	query   string
	where   []string
	orderby string
	groupby string
	limit   [2]int
}

func NewGormDB(table string) *GormDB {
	return &GormDB{
		DB: DB,
	}
}

func (gormDB *GormDB) Search(result interface{}, parameter Parameter) error {
	return gormDB.build(parameter).Find(result).Error
}

func (gormDB *GormDB) SearchOne(result interface{}, parameter Parameter) error {
	return gormDB.build(parameter).First(result).Error
}

func (gormDB *GormDB) build(parameter Parameter) *gorm.DB {
	DB := gormDB.DB
	if DB == nil {
		panic("user stop run!")
	}

	bind := new(bind)
	bind.where = make([]string, 0)

	// parse
	for key, val := range parameter {
		field := strings.ToLower(strings.TrimSpace(key))

		switch strings.ToLower(field) {
		case query:
			if value, ok := val.([]string); ok {
				bind.query = strings.Join(value, ", ")
			}
		case orderby:
			if value, ok := val.(string); ok {
				bind.orderby = value
			}
		case groupby:
			if value, ok := val.(string); ok {
				bind.groupby = value
			}
		case offset:
			if value, ok := val.(int); ok {
				bind.limit[0] = value
			}
		case limit:
			if value, ok := val.(int); ok {
				bind.limit[1] = value
			}
		default:
			switch val.(type) {
			case string:
				if strings.Contains(val.(string), " ") {
					bind.where = append(bind.where, (fmt.Sprintf(`(%s %s)`, field, val.(string))))
					continue
				}
				bind.where = append(bind.where, fmt.Sprintf(`(%s = "%s")`, field, val.(string)))
			default:
				bind.where = append(bind.where, fmt.Sprintf(`(%s = %v)`, field, val))
			}
		}
	}

	// bind
	if len(bind.query) != 0 {
		DB = DB.Select(bind.query)
	}
	if len(bind.where) != 0 {
		DB = DB.Where(strings.Join(bind.where, " AND "))
	}
	if len(bind.orderby) != 0 {
		DB = DB.Order(bind.orderby)
	}
	if len(bind.groupby) != 0 {
		DB = DB.Group(bind.groupby)
	}
	if bind.limit[0] != 0 || bind.limit[1] != 0 {
		DB = DB.Offset(bind.limit[0]).Limit(bind.limit[1])
	}

	return DB
}
