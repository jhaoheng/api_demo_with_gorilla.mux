SHELL:=/bin/bash
.PHONY: test

test:
	@\
	newman run "./postman/postman_collection.json" -e "./postman/env.json" -r cli,junit --bail --disable-unicode;