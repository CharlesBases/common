package gorm

import (
	"fmt"
	"time"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"

	"github.com/CharlesBases/common/log"
)

var (
	DB *gorm.DB

	debug           = false
	maxIdleConns    = 2000
	maxOpenConns    = 1000
	connMaxLifetime = 10
)

func InitGorm(database string) {
	initMySql(database)
}

func initMySql(database string) {
	db, err := gorm.Open("mysql", database)
	if err != nil {
		log.Error(fmt.Sprintf(" - db dsn(%s) error - ", database), err.Error())
		return
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

	if DB != nil {
		DB.Close()
	}

	DB = db
}

type Logger struct {
}

func (l *Logger) Print(v ...interface{}) {
	log.Debug("SQL - ", v)
}
