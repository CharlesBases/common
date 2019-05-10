package db

import (
	"os"

	"common/log"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

var (
	DB *gorm.DB

	debug        = false
	maxIdleConns = 2000
	maxOpenConns = 1000
)

func InitGorm(database string) *gorm.DB {
	openMySql(database)
	return DB
}

func openMySql(database string) {
	db, err := gorm.Open("mysql", database)
	if err != nil {
		log.Error(" - 数据库连接失败 - ", err.Error())
		os.Exit(0)
	}

	if db.DB() == nil {
		log.Error(" - 数据库连接出现未知错误 - ", err.Error())
		os.Exit(0)
	}

	err = db.DB().Ping()
	if err != nil {
		log.Error(" - 数据库Ping不通 - ", err.Error())
		os.Exit(0)
	}

	db.DB().SetMaxIdleConns(maxIdleConns)
	db.DB().SetMaxOpenConns(maxOpenConns)

	db.LogMode(debug)
	db.SetLogger(new(Logger))
	db.SingularTable(true)
	db.BlockGlobalUpdate(true)
	db.Set("gorm:association_autoupdate", false).Set("gorm:association_autocreate", false)

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
