# Dev Diaries
Blog API written entirely in [Go](https://go.dev/). This repository contains the backend code written in Go with MySQL as the database.

## Steps to run
- This project uses [Infisical](https://infisical.com/) as the secrets vault. Create an account for their free cloud service or have a local instance running
- Add the following secrets to Infisical:
    - `DB_PASSWORD=<your-MySQL-user-password>`
    - `DB_USER=<your-MySQL-username>`
    - `JWT_SECRET=<256-bit-secret>`
- Add a .env file in the project root with the following variables
    - `DB_URL=@tcp(<your-db-host>:<your-db-port>)/blog_db?charset=utf8mb4&parseTime=True&loc=Local`
    - `TOKEN_VALIDITY_MINUTES=<duration-of-JWT-Token-in-minutes>`
    - `INFISCAL_SECRET=<infisical-secret>`
    - `INFISCAL_CLIENT_ID=<infisical-client-id>`
    - `INFISCAL_HOST=<infisical-host-url>`
    - `INFISCAL_ENVIRONMENT=<dev/staging/prduction>`
    - `INFISCAL_WORKSPACE_ID=<infisical-workspace-id>`
- Run `go mod download` to install all dependencies
- Run `go run main.go` to run the server on PORT 4000

## [API](https://app.swaggerhub.com/apis/vineet192/DevDiaries/1.0.0)