package bot

import (
	"io"
	"mime/multipart"
)

type Config struct {
	WebhookURL          string
	FreeAnswerMessage   string
	NoFreeAnswerMessage string
}

func Run(config Config) {
	body, bodyWriter := io.Pipe()
	mpartWriter := multipart.NewWriter(bodyWriter)
	mpartBoundary := mpartWriter.Boundary()

	go writeMessage(&config, bodyWriter, mpartWriter)
	executeWebhook(&config, body, mpartBoundary)
}
