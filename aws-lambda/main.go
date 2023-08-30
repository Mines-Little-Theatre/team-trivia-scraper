package main

import (
	"context"
	"log"

	"github.com/Mines-Little-Theatre/team-trivia-scraper/bot"
	"github.com/Mines-Little-Theatre/team-trivia-scraper/utils"
	"github.com/aws/aws-lambda-go/lambda"
)

func main() {
	log.SetFlags(log.Lshortfile)

	config := utils.ReadConfig()

	lambda.Start(func(ctx context.Context) error {
		return bot.Run(ctx, config)
	})
}
