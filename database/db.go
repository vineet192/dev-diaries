package database

import (
	"devdiaries/models"
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
	DB.AutoMigrate(
		&models.User{},
		&models.Blog{},
		&models.BlogReaction{},
		&models.CommentReaction{},
		&models.Tag{})

	if err == nil {
		fmt.Println("Database connected successfully")
	} else {
		panic(err)
	}
}
