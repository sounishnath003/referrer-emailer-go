
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
	curl http://localhost:3000/api/send-email -H "Content-Type: application/json" -d '{"from": "flock.sinasini@gmail.com", "to": ["flock.sinasini@gmail.com", "almyhle.johnshon@gmail.com"], "subject": "my-first-subject", "body": "<div><h3>Welcome Sounish!</h3><p>To the whole new world of building softwares referrals</p><p>Lorem ipsum dolor sit amet consectetur adipisicing elit. Quod ea, et in autem sint, ex sapiente consequunturassumenda magni est debitis voluptas nemo praesentium optio, itaque nisi minus totam quo.</p><p> Thanks and regards</p><p>Customgo-emailer-service</p></div>"}' | jq