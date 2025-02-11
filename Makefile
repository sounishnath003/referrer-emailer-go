
.PHONY: install
install:
	go mod tidy
	go mod download
	go mod verify

.PHONY: build
build:
	CGO_ENABLED=0 GO_ARCH=amd64 go build -o ./tmp/main ./cmd/*.go

.PHONY: run
run: build
	./tmp/main

.PHONY: web
web:
	cd web && npm start

.PHONY: all
all: build
	source .env && make run &
	cd web && npm start

.PHONY: docker-build
docker-build:
	docker build -t referrer-emailer -f Dockerfile .
	docker images