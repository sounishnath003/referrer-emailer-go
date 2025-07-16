
.PHONY: install
install:
	go mod tidy
	go mod download
	go mod verify

.PHONY: build
build:
	CGO_ENABLED=0 GO_ARCH=amd64 go build -ldflags "-s -w" -o ./tmp/main ./cmd/*.go

.PHONY: run
run: build
	./tmp/main

.PHONY: web
web:
	cd web && npm start

.PHONY: pdf-service
pdf-service:
	cd pdf-service && npm run dev

.PHONY: all
all: build
	cd pdf-service && npm run dev &
	source .env && make run &
	cd web && npm start

.PHONY: compose-up
compose-up:
	docker-compose down
	# docker rmi referrer-emailer
	docker-compose up --build

.PHONY: docker-build
docker-build:
	docker rmi -f $$(docker images referrer-emailer -qa)
	docker build -t referrer-emailer -f Dockerfile .
	docker images

.PHONY: docker-run
docker-run:
	docker images
	source .env;
	docker run -ti -e MAIL_ADDR=$MAIL_ADDR -e MAIL_SECRET=$MAIL_SECRET -e MONGO_DB_URI=$MONGO_DB_URI -v ./storage:/home/nonroot/storage -p 3000:3000 referrer-emailer:latest