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

	lambda.Start(func(context.Context, any) ([]byte, error) {
		return nil, bot.Run(config)
	})
}
