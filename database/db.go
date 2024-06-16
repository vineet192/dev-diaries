package database

import (
	"devdiaries/api/utilities"
	"devdiaries/models"
	"errors"
	"fmt"
	"os"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB
var err error

func InitDB() {

	dbUser, DBUserErr := utilities.GetSecret("DB_USER")
	dbPass, DBPassErr := utilities.GetSecret("DB_PASSWORD")

	if errors.Join(DBUserErr, DBPassErr) != nil {
		panic(errors.Join(DBUserErr, DBPassErr))
	}

	dsn := fmt.Sprintf("%s:%s%s", dbUser, dbPass, os.Getenv("DB_URL"))
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
