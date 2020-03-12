.PHONY: build
build:
	docker build -t kenanbek/flash:latest .

.PHONY: run
run:
	docker run --rm -p 8080:8080 kenanbek/flash

.PHONY: push
push:
	docker push kenanbek/flash

.PHONY: build/%
build/%:
	docker build -t kenanbek/flash:$* .

.PHONY: run/%
run/%:
	docker run --rm -p 8080:8080 kenanbek/flash:$*

.PHONY: push/%
push/%:
	docker push kenanbek/flash:$*
