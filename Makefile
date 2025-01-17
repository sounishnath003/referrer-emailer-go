
.PHONY: install
install:
	go mod tidy
	go mod download
	go mod verify

.PHONY: build
build:
	CGO_ENABLED=0 GOOS=linux GO_ARCH=amd64 go build -o ./tmp/main cmd/*.go

.PHONY: run
run: build
	./tmp/main

.PHONY: send-email
send-email:
	curl http://localhost:3000/api/send-email -H "Content-Type: application/json" -d '{"from": "flock.sinasini@gmail.com", "to": ["flock.sinasini@gmail.com", "almyhle.johnshon@gmail.com"], "subject": "my-first-subject", "body": "My emailer custom go is running. <strong>my content of email</strong>. I am doing super good."}' | jq