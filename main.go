package main

import (
	"log"
	"os"

	"github.com/ayush6624/go-chatgpt"
	frezhchatgpt "github.com/jrh3k5/frezh/internal/chatgpt"
	frezhhttp "github.com/jrh3k5/frezh/internal/http"
	"github.com/jrh3k5/frezh/internal/ocr"
)

func main() {
	log.Print("Starting frezh...")

	chatgptKey := os.Getenv("CHATGPT_KEY")
	if chatgptKey == "" {
		log.Fatal("CHATGPT_KEY is not set")
	}

	chatgptClient, err := chatgpt.NewClient(chatgptKey)
	if err != nil {
		log.Fatalf("failed to create ChatGPT client: %v", err)
	}

	chatgptService := frezhchatgpt.NewAyushService(chatgptClient)
	ocrProcessor := &ocr.Gosseract{}

	if err := frezhhttp.StartServer(chatgptService, ocrProcessor); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}
