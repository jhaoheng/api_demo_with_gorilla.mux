SHELL:=/bin/bash
.PHONY: test dev

test:
	@\
	newman run "./postman/postman_collection.json" -e "./postman/env.json" -r cli,junit --bail --disable-unicode;


dev:
	@\
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o cmd/docker_app ./app;\
	chmod 777 cmd/docker_app;\
	docker-compose down && docker-compose up -d;\