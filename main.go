package main

import (
	"log"

	"github.com/jrh3k5/frezh/http"
)

func main() {
	log.Print("Starting frezh...")

	if err := http.StartServer(); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}
