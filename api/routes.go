package api

import (
	"devdiaries/api/auth"
	"devdiaries/api/blog"
	"devdiaries/api/comment"
	"devdiaries/api/middleware"
	"devdiaries/api/user"

	"github.com/gorilla/mux"
)

func RegisterUserRoutes(router *mux.Router) {

	router.HandleFunc("/follower/{follower_id}", user.AddFollower).Methods("POST")
	router.HandleFunc("/{id}", user.EditUser).Methods("PUT")
	router.HandleFunc("/{id}", user.GetUserByID).Methods("GET")
	router.HandleFunc("/", user.GetUser).Methods("GET")
	router.HandleFunc("/{id}/blog", user.GetBlogs).Methods("GET")
	router.HandleFunc("/{id}/blog_feed", user.GetBlogFeed).Methods("GET")

	//Logged in user can only delete their own account
	deleteUserRoute := router.Path("/{id}").Subrouter()
	deleteUserRoute.Use(middleware.ValidateUserID)
	deleteUserRoute.HandleFunc("", user.DeleteUserByID).Methods("DELETE")

	router.HandleFunc("/follower/{follower_id}", user.RemoveFollower).Methods("DELETE")
}

func RegisterBlogRoutes(router *mux.Router) {
	router.HandleFunc("/", blog.PostBlog).Methods("POST")
	router.HandleFunc("/{id}/comment", blog.PostComment).Methods("POST")
	router.HandleFunc("/{id}/reaction", blog.PostReaction).Methods("POST")
	router.HandleFunc("/{blog_id}/reaction/{user_id}", blog.DeleteReaction).Methods("DELETE")
	router.HandleFunc("/{id}", blog.EditBlog).Methods("PUT")
	router.HandleFunc("/{id}", blog.DeleteBlogByID).Methods("DELETE")
}

func RegisterCommentRoutes(router *mux.Router) {
	router.HandleFunc("/{id}", comment.DeleteCommentByID).Methods("DELETE")
	router.HandleFunc("/{id}/reaction", comment.PostReaction).Methods("POST")
	router.HandleFunc("/{comment_id}/reaction/{user_id}", comment.DeleteReaction).Methods("DELETE")
}

func RegisterAuthRoutes(router *mux.Router) {
	router.HandleFunc("/signup", auth.SignUp).Methods("POST")
	router.HandleFunc("/login", auth.Login).Methods("POST")
}
