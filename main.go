package main

import (
	"log"
	"os"

	"github.com/Mines-Little-Theatre/team-trivia-scraper/bot"
)

func main() {
	bot.Run(bot.Config{
		WebhookID:    readEnv("TRIVIA_WEBHOOK_ID"),
		WebhookToken: readEnv("TRIVIA_WEBHOOK_TOKEN"),
		Message:      readEnv("TRIVIA_MESSAGE"),
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
