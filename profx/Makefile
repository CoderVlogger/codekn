.PHONY: build
build:
	docker build -t kenanbek/profx:latest .

.PHONY: run
run:
	docker run --rm kenanbek/profx

.PHONY: push
push:
	docker push kenanbek/profx

.PHONY: build/%
build/%:
	docker build -t kenanbek/profx:$* .

.PHONY: run/%
run/%:
	docker run --rm kenanbek/profx:$*

.PHONY: push/%
push/%:
	docker push kenanbek/profx:$*

.PHONY: local/db
local/db/start:
	docker run -d --name profx -p 3306:3306 -e MYSQL_USER=profx -e MYSQL_PASSWORD=profx -e MYSQL_DATABASE=profx -e MYSQL_ROOT_PASSWORD=profx mysql:5.7.25

.PHONY: local/db/stop
local/db/stop:
	docker stop profx

.PHONY: local/db/destroy
local/db/destroy: local/db/stop
	docker rm profx

.PHONY: local/db/recreate
local/db/recreate: local/db/destroy local/db/start
	echo "recreated"
