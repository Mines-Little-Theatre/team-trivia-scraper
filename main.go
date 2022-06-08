package main

import (
	"io"
	"log"
	"mime/multipart"
	"os"
)

type Config struct {
	WebhookURL          string
	FreeAnswerMessage   string
	NoFreeAnswerMessage string
}

func main() {
	config := Config{
		WebhookURL:          readEnv("WEBHOOK_URL"),
		FreeAnswerMessage:   readEnv("FREE_ANSWER_MESSAGE"),
		NoFreeAnswerMessage: readEnv("NO_FREE_ANSWER_MESSAGE"),
	}

	body, bodyWriter := io.Pipe()
	mpartWriter := multipart.NewWriter(bodyWriter)
	mpartBoundary := mpartWriter.Boundary()

	go writeMessage(&config, bodyWriter, mpartWriter)
	executeWebhook(&config, body, mpartBoundary)
}

func readEnv(key string) string {
	var result string
	var ok bool
	if result, ok = os.LookupEnv(key); !ok {
		log.Fatalf("please set the %s environment variable", key)
	}
	return result
}
