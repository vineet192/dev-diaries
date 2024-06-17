package main

import (
	"devdiaries/api"
	"devdiaries/database"
	"fmt"
	"net/http"

	"github.com/lpernett/godotenv"

	"github.com/gorilla/mux"
)

func main() {
	r := mux.NewRouter()

	userRouter := r.PathPrefix("/user").Subrouter()
	blogRouter := r.PathPrefix("/blog").Subrouter()
	commentRouter := r.PathPrefix("/comment").Subrouter()
	authRouter := r.PathPrefix("/").Subrouter()

	api.RegisterUserRoutes(userRouter)
	api.RegisterBlogRoutes(blogRouter)
	api.RegisterCommentRoutes(commentRouter)
	api.RegisterAuthRoutes(authRouter)

	godotenv.Load()
	database.InitDB()

	err := http.ListenAndServe("0.0.0.0:4000", r)

	if err == nil {
		fmt.Println("Server running successfully!")
	} else {
		panic(err)
	}
}
