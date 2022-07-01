SHELL:=/bin/bash
.PHONY: test build dev

test:
	@\
	newman run "./postman/postman_collection.json" -e "./postman/env.json" -r cli,junit --bail --disable-unicode;

build:
	@\
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o cmd/docker_app ./app;

key:
	@\
	cd app/keypair;\
	ssh-keygen -t rsa -b 4096 -m PEM -f jwt_rs256.key;\
	openssl rsa -in jwt_rs256.key -pubout -outform PEM -out jwt_rs256.key.pub;