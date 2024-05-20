package database

import (
	"fmt"
	"os"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB
var err error

func InitDB() {
	dsn := os.Getenv("DB_URL")
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{TranslateError: true})

	if err == nil {
		fmt.Println("Database connected successfully")
	} else {
		panic(err)
	}
}
