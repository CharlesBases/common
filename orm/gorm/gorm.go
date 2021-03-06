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
	return "root:password@tcp(127.0.0.1:3306)/mysql?charset=utf8mb4&parseTime=True&loc=Local"
}

type Logger struct{}

func (l *Logger) Print(v ...interface{}) {
	if len(v) != 0 {
		switch v[0] {
		case "sql":
			if len(v) == 6 {
				log.Debugf("MySQL >>>\npath: %v\nsql : %v | %v | %v | %v", v[1], v[3], v[4], v[5], v[2])
				return
			}
		case "log":
			if len(v) == 3 {
				log.Errorf("MySQL >>>\npath: %v\nerr : %v", v[1], v[2])
				return
			}
		case "info":
			if len(v) == 2 {
				log.Infof("MySQL >>>\ninfo: %v", v[1])
				return
			}
		case "warn":
			if len(v) == 2 {
				log.Warnf("MySQL >>>\nwarn: %v", v[1])
				return
			}
		}

		log.Debug("MySQL >>>", v)
	}
}
