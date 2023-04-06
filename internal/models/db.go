package models

import (
	"apify-service/internal/driver"
	"log"
	"sync"

	"gorm.io/gorm"
)

var lock = &sync.Mutex{}

type DB struct {
	SQL *gorm.DB
}

var singleInstance *DB

func GetDB() *gorm.DB {
	if singleInstance == nil {
		lock.Lock()
		defer lock.Unlock()
		if singleInstance == nil {
			db, err := driver.ConnectSQL()
			if err != nil {
				log.Fatal(err)
			}
			singleInstance = &DB{
				SQL: db,
			}
		}
	}

	return singleInstance.SQL
}
