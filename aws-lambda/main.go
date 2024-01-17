package main

import (
	"log"

	"github.com/Mines-Little-Theatre/team-trivia-scraper/bot"
	"github.com/aws/aws-lambda-go/lambda"
)

func main() {
	log.SetFlags(log.Lshortfile)
	lambda.Start(bot.Run)
}
