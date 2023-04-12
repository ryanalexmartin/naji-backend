package main

import (
	"backend/server"
)

func main() {
	addr := "localhost:8080"
	server.StartServer(addr)
}
