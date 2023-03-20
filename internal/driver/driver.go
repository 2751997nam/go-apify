package driver

import (
	"fmt"
	"log"
	"product-service/internal/config"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func ConnectSQL() (*gorm.DB, error) {
	dsn := config.GetDsn()
	fmt.Println(dsn)
	db, err := NewDatabase(dsn)
	if err != nil {
		panic(err)
	}

	err = testDB(db)
	if err != nil {
		return nil, err
	}

	return db, nil
}

func testDB(db *gorm.DB) error {
	d, err := db.DB()
	if err != nil {
		return err
	}

	err = d.Ping()
	if err != nil {
		return err
	}

	return nil
}

func NewDatabase(dsn string) (*gorm.DB, error) {

	counts := 0
	for {
		db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
			CreateBatchSize: 1000,
		})
		if err != nil {
			log.Println("Mysql not yet ready ...")
			counts++
		} else {
			log.Println("Connected to Mysql!")
			return db, nil
		}

		if counts > 20 {
			log.Println(err)
			return nil, err
		}

		log.Println("Backing off for two seconds....")
		time.Sleep(2 * time.Second)
	}
}
