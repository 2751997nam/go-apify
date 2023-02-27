package config

import (
	"fmt"
	"os"
)

type DatabaseConfig struct {
	Port     string
	Host     string
	User     string
	Password string
	Database string
}

func (dc *DatabaseConfig) toString() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", dc.User, dc.Password, dc.Host, dc.Port, dc.Database)
}

func GetDsn() string {
	dc := DatabaseConfig{
		Port:     os.Getenv("DB_PORT"),
		Host:     os.Getenv("DB_HOST"),
		User:     os.Getenv("DB_USER"),
		Password: os.Getenv("DB_PASSWORD"),
		Database: os.Getenv("DB_DATABASE"),
	}

	return dc.toString()
}
