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

	userRouter := r.PathPrefix("/user").Subrouter()
	blogRouter := r.PathPrefix("/blog").Subrouter()
	commentRouter := r.PathPrefix("/comment").Subrouter()

	api.RegisterUserRoutes(userRouter)
	api.RegisterBlogRoutes(blogRouter)
	api.RegisterCommentRoutes(commentRouter)

	godotenv.Load()
	database.InitDB()

	err := http.ListenAndServe("0.0.0.0:4000", r)

	if err == nil {
		fmt.Println("Server running successfully!")
	} else {
		panic(err)
	}
}
