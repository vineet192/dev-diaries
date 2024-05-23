package api

import (
	"inventory/api/blog"
	"inventory/api/comment"
	"inventory/api/user"

	"github.com/gorilla/mux"
)

func RegisterUserRoutes(router *mux.Router) {
	router.HandleFunc("/", user.CreateUser).Methods("POST")
	router.HandleFunc("/", user.EditUser).Methods("PUT")
	router.HandleFunc("/{id}", user.DeleteUserByID).Methods("DELETE")
}

func RegisterBlogRoutes(router *mux.Router) {
	router.HandleFunc("/", blog.PostBlog).Methods("POST")
	router.HandleFunc("/{id}/comment", blog.PostComment).Methods("POST")
	router.HandleFunc("/{id}", blog.EditBlog).Methods("PUT")
	router.HandleFunc("/{id}", blog.DeleteBlogByID).Methods("DELETE")
}

func RegisterCommentRoutes(router *mux.Router) {
	router.HandleFunc("/{id}", comment.DeleteCommentByID).Methods("DELETE")
}
