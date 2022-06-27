# API demo, 目的
- 使用 Golang 的框架 gorilla/mux framework
- 使用 swagger 來產生 API 文件
- 使用 RS256 來產生 keypair, 給 JWT 進行認證使用

# 版本
- golang 1.16.5
- gorilla/mux, v1.8.0
- gorm, v1.21.11
- postman, collection v2.1, newman, 5.3.2
- Docker version 20.10.6, build 370c289
- docker-compose version 1.29.1, build c34c88b2

# Serivce securities
- HTTPS
- CSRF
- JWT Authorization
- XSS protection
- HSTS Protection
- Forbidden to show nginx's version
- SQL injection

# 本地執行 API Service
1. 產生 RS256 keypair
    1. `cd keypair`
    2. `ssh-keygen -t rsa -b 4096 -m PEM -f jwt_rs256.key`, kepp the passphrase is empty.
    3. `openssl rsa -in jwt_rs256.key -pubout -outform PEM -out jwt_rs256.key.pub`  
2. `docker-compose up -d`
3. `cd app && go run main.go`, port is `:8080`

# websocket
1. 執行本地 API service
2. Run the client, `cd websocket_client && go run main.go`

# run api unit test
1. 執行本地 API service
2. `docker exec app go test ./...`

# postman, newman, 執行自動測試
1. 執行本地 API service
2. import the file: `./postman/...` into your postman app.

# swagger, 產生文件
- swagger info
    - version: v0.27.0
    - commit: 43c2774170504d87b104e3e4d68626aac2cd447d
    - [github](https://github.com/go-swagger/go-swagger)
    - [install](https://goswagger.io/install.html)
- `swagger generate spec -w ./swaggerdoc -o ./swagger.json`
- 產生文件: `swagger serve swagger.json`
- 建立文件網站: `swagger serve -F swagger swagger.json`