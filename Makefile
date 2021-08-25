SHELL:=/bin/bash
.PHONY: db table

db:
	@\
	docker run -d --name db -p 5432:5432 -e POSTGRES_DB=ui_test -e POSTGRES_USER=ui_test -e POSTGRES_PASSWORD=ui_test postgres:13.3;

table: 
	@\
	docker exec -it db psql -U ui_test -w ui_test -c "CREATE TABLE users(acct varchar(20) NOT NULL, pwd text NOT NULL, fullname varchar(20) NOT NULL, created_at timestamp NOT NULL DEFAULT NOW(), updated_at timestamp NOT NULL DEFAULT NOW(), PRIMARY KEY(acct), UNIQUE(acct));";