package main

import (
	"github.com/sounishnath003/customgo-mailer-service/internal/core"
	"github.com/sounishnath003/customgo-mailer-service/internal/server"
)

func main() {

	co := core.NewCore()

	server := server.NewServer(co)
	panic(server.Start())
}
