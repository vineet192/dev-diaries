package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {

	r := mux.NewRouter()

	err := http.ListenAndServe("0.0.0.0:3000", r)

	if err != nil {
		fmt.Println(err, "An error has occured while starting the server")
	}
	fmt.Println("Welcome to Vineet's inventory API, written in GO!")
}
