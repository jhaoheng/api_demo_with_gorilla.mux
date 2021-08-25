# readme
- This is a demo for show how to build APIs with gorilla/mux framework.
- And use swagger to build the docs.
- This repo also use RS256 to generate the keypair for JWT authentication.

# run web api server by golang
- `cd app && go run main.go`, port is `:8080`
- golang 1.16.5
- gorilla/mux, v1.8.0
- gorm, v1.21.11
- securities
  - XSS
  - CSRF
  - JWT Authorization

# run docker-compose, localhost api server
1. Generate RS256 keypair for JWT at localhost
  1. `cd keypair`
  2. `ssh-keygen -t rsa -b 4096 -m PEM -f jwtRS256.key`, kepp the passphrase is empty.
  3. `openssl rsa -in jwtRS256.key -pubout -outform PEM -out jwtRS256.key.pub`  
2. `docker-compose up -d`
  - Docker version 20.10.6, build 370c289
  - docker-compose version 1.29.1, build c34c88b2
  - securities
    - HTTPS
    - XSS protection
    - HSTS Protection
    - Forbidden to show nginx's version

# websocket
1. Run api server
2. Run the client, `cd websocket_client && go run main.go`

# run api unit test
1. Run api server
2. `docker exec app go test`

# postman
1. Run api server
2. import the file: `./postman/...` into your postman app.

# swagger, generate doc
- swagger info
  - version: v0.27.0
  - commit: 43c2774170504d87b104e3e4d68626aac2cd447d
  - [github](https://github.com/go-swagger/go-swagger)
  - [install](https://goswagger.io/install.html)
- `swagger generate spec -w ./swaggerdoc -o ./swagger.json`
- start
  - `swagger serve swagger.json`
  - `swagger serve -F swagger swagger.json`

# deploy to AWS ECS
1. Build golang app
2. Build docker image
  - add the go app in it.
3. Publish the image to AWS ECR.
4. Build two ECS Task
  - One is set the parameter `env=stg` in the task.
  - Another is set `env=prod`
5. Run the ECS Service with the ECS Task which env you need to work.