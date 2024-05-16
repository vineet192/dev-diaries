package main

import (
	"fmt"
	"inventory/api"
	"inventory/database"
	"net/http"

	"github.com/lpernett/godotenv"

	"github.com/gorilla/mux"
)

func main() {

	r := mux.NewRouter()
	r.HandleFunc("/user", api.CreateUser).Methods("POST")
	godotenv.Load()
	database.InitDB()

	err := http.ListenAndServe("0.0.0.0:3000", r)

	if err == nil {
		fmt.Println("Server running successfully!")
	} else {
		panic(err)
	}
}
