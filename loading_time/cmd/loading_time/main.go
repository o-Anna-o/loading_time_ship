package main

import (
	"loading_time/internal/api"
	"log"
)

func main() {
	log.Println("Application started!")
	api.StartServer()
	log.Println("Application terminated!")
}
