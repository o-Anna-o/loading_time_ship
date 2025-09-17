package main

// go run cmd/loading_time/main.go

import (
	"loading_time/internal/api"
	"log"
)

func main() {
	log.Println("Application started!")
	api.StartServer()
	log.Println("Application terminated!")
}
