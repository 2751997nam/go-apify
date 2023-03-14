package main

import (
	"fmt"
	"log"
	"os"
	"product-service/internal/driver"

	"gorm.io/gorm"
)

var webPort = os.Getenv("PORT")

func main() {
	if webPort == "" {
		webPort = "8080"
	}

	db, err := run()
	if err != nil {
		log.Fatal("cannot open connection with db")
	}
	d, err := db.DB()
	if err != nil {
		log.Fatal("cannot open connection with db")
	}
	defer d.Close()

	router := routes()

	router.Run(fmt.Sprintf(":%s", webPort))

}

func run() (*gorm.DB, error) {
	db, err := driver.ConnectSQL()
	if err != nil {
		log.Fatal("cannot open connection with db")
	}

	return db, err
}
