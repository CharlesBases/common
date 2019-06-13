package db

import (
	"fmt"
	"sync"
	"time"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"

	"github.com/CharlesBases/common/log"
)

var (
	DB *gorm.DB

	debug        = false
	maxIdleConns = 2000
	maxOpenConns = 1000

	Sync sync.RWMutex
)

func InitGorm(database string) {
	initMySql(database)

	go func() {
		ticker := time.NewTicker(time.Second * 10)
		for {
			<-ticker.C
			if !ping() {
				log.Error(fmt.Sprintf(" - db ping error, connect again. - "))
				initMySql(database)
			}
		}
	}()
}

func initMySql(database string) {
	db, err := gorm.Open("mysql", database)
	if err != nil {
		log.Error(fmt.Sprintf(" - db dsn(%s) error - ", database), err.Error())
		return
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

func ping() bool {
	Sync.RLock()
	status := true
	if err := DB.DB().Ping(); err != nil {
		DB.Close()
		DB = nil
		status = false
	}
	Sync.RUnlock()
	return status
}

type Logger struct {
}

func (l *Logger) Print(v ...interface{}) {
	log.Debug("SQL - ", v)
}
