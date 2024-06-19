# Dev Diaries
Dev diaries is a blog service. This repository contains the backend code written in Go with MySQL as the database.

## Dependencies
- [Go](https://go.dev/doc/install)
- [Gorilla/Mux](https://github.com/gorilla/mux)
- [GORM](https://gorm.io/index.html)
- [MySQL](https://dev.mysql.com/downloads/mysql/)

## Steps to run
- Set up a local or remote MySQL instance and have it running
- Add a .env file in the project root with the following variables
    - `DB_URL=<mysql_connection_string>`
- Run `go run main.go` to run the server on PORT 4000

## API
### AUTH

`POST /signup`
[Body](payload/request/signupBody.go)

`POST /login`
[Body](payload/request/loginBody.go)

### [USER](models/user.go)
`GET /user/{id}`

`GET /user/{id}/blog`
[Query](payload/request/blogQuery.go)

`GET /user/{id}/blog_feed`

`POST /user/{id}/follower/{follower_id}`
[Body](models/user.go)

`PUT /user/{id}`
[Body](models/user.go)

`DELETE /user/{id}`

`DELETE user/{id}/follower/{follower_id}`

### [BLOG](models/blog.go)
`POST /blog/`
[Body](models/blog.go)

`POST /blog/{id}/comment`
[Body](models/comment.go)

`POST /blog/{id}/reaction`
[Body](models/blog_reaction.go)

`PUT /blog/{id}`
[Body](models/blog.go)

`DELETE /blog/{id}`

`DELETE /blog/{blog_id}/reaction/{user_id}`

### [COMMENT](models/comment.go)
`POST /comment/{id}/reaction`
[Body](models/comment_reaction.go)

`DELETE /comment/{id}`

`DELETE /comment/{comment_id}/reaction/{user_id}`