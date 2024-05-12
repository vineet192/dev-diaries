package main

import (
	"fmt"
	"inventory/utilities"
	"net/http"

	"github.com/gorilla/mux"
)

// func initDB() {
// 	dsn := "user:pass@tcp(127.0.0.1:3306)/dbname?charset=utf8mb4&parseTime=True&loc=Local"
// 	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
// }

func main() {

	r := mux.NewRouter()

	utilities.LoadEnv()

	err := http.ListenAndServe("0.0.0.0:3000", r)

	if err != nil {
		fmt.Println(err, "An error has occured while starting the server")
	}
	fmt.Println("Welcome to Vineet's inventory API, written in GO!")
}
