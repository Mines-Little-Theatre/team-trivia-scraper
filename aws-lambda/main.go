package main

import (
	"context"
	"log"
	"os"

	"github.com/Mines-Little-Theatre/team-trivia-scraper/bot"
	"github.com/aws/aws-lambda-go/lambda"
)

func main() {
	config := bot.Config{
		WebhookID:    readEnv("TRIVIA_WEBHOOK_ID"),
		WebhookToken: readEnv("TRIVIA_WEBHOOK_TOKEN"),
		Message:      readEnv("TRIVIA_MESSAGE"),
	}

	lambda.Start(func(context.Context, any) ([]byte, error) {
		bot.Run(config)
		return nil, nil
	})
}

func readEnv(key string) string {
	var result string
	var ok bool
	if result, ok = os.LookupEnv(key); !ok {
		log.Fatalf("please set the %s environment variable", key)
	}
	return result
}
