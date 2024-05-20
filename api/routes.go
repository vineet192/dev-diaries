package api

import (
	"inventory/api/blog"
	"inventory/api/user"

	"github.com/gorilla/mux"
)

func RegisterUserRoutes(router *mux.Router) {
	router.HandleFunc("/", user.CreateUser).Methods("POST")
}

func RegisterBlogRoutes(router *mux.Router) {
	router.HandleFunc("/", blog.PostBlog).Methods("POST")
}
