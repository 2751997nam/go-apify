package driver

import (
	"fmt"
	"product-service/internal/config"

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
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		CreateBatchSize: 1000,
	})
	if err != nil {
		return nil, err
	}

	return db, nil
}
