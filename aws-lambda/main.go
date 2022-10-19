package main

import (
	"context"

	"github.com/Mines-Little-Theatre/team-trivia-scraper/bot"
	"github.com/Mines-Little-Theatre/team-trivia-scraper/utils"
	"github.com/aws/aws-lambda-go/lambda"
)

func main() {
	config := bot.Config{
		WebhookID:    utils.ReadEnv("TRIVIA_WEBHOOK_ID"),
		WebhookToken: utils.ReadEnv("TRIVIA_WEBHOOK_TOKEN"),
		Message:      utils.ReadEnv("TRIVIA_MESSAGE"),
	}

	lambda.Start(func(context.Context, any) ([]byte, error) {
		bot.Run(config)
		return nil, nil
	})
}
