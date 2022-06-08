package main

import (
	"log"
	"os"

	"github.com/Mines-Little-Theatre/team-trivia-scraper/bot"
)

func main() {
	bot.Run(bot.Config{
		WebhookURL:          readEnv("WEBHOOK_URL"),
		FreeAnswerMessage:   readEnv("FREE_ANSWER_MESSAGE"),
		NoFreeAnswerMessage: readEnv("NO_FREE_ANSWER_MESSAGE"),
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
