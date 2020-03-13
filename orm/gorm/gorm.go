package gorm

import (
	"fmt"
	"time"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"

	"github.com/CharlesBases/common/log"
)

var (
	orm *gorm.DB

	debug           = false
	maxIdleConns    = 2000
	maxOpenConns    = 1000
	connMaxLifetime = 10
)

func init() {
	db, err := gorm.Open("mysql", addr())
	if err != nil {
		log.Error(fmt.Sprintf(" - db connect(%s) error - %s", addr(), err.Error()))
		panic("user stop run")
	}

	db.DB().SetMaxIdleConns(maxIdleConns)
	db.DB().SetMaxOpenConns(maxOpenConns)
	db.DB().SetConnMaxLifetime(time.Duration(connMaxLifetime) * time.Second)

	db.LogMode(debug)
	db.SetLogger(new(Logger))
	db.SingularTable(true)
	db.BlockGlobalUpdate(true)
	db.Set("gorm:association_autoupdate", false).Set("gorm:association_autocreate", false)

	// db.Exec("set sql_mode=(select replace(@@sql_mode,'ONLY_FULL_GROUP_BY',''))")

	if orm != nil {
		orm.Close()
	}

	orm = db
}

func Gorm() *gorm.DB {
	return orm
}

func addr() string {
	return "root:password@tcp(127.0.0.1:3306)/mysql"
}

type Logger struct{}

func (l *Logger) Print(v ...interface{}) {
	if len(v) != 0 {
		switch v[0] {
		case "log":
			log.Errorf("MySQL >>> \npath: %v\nerr : %v", v[1], v[2])
		case "sql":
			log.Debugf("MySQL >>> \npath: %v\nsql : %v | %v | %v | %v", v[1], v[3], v[4], v[5], v[2])
		case "info":
			log.Info("SQL - ", v)
		case "warn":
			log.Warn("SQL - ", v)
		}
	}
}
