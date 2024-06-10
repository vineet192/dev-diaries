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
### [USER](models/user.go)
```
GET /user/{id}
GET /user/{id}/blog
GET /user/{id}/blog_feed
```

```
POST /user/
POST /user/{id}/follower/{follower_id}
```
```
PUT /user/{id}
```
```
DELETE /user/{id}
DELETE user/{id}/follower/{follower_id}
```

### [BLOG](models/blog.go)
```
POST /blog/
POST /blog/{id}/comment
POST /blog/{id}/reaction
```
```
PUT /blog/{id}
```
```
DELETE /blog/{id}
DELETE /blog/{blog_id}/reaction/{user_id}
```

### [COMMENT](models/comment.go)
```
POST /comment/{id}/reaction
```
```
DELETE /comment/{id}
DELETE /comment/{comment_id}/reaction/{user_id}
```