package utils

import (
	"log"
	"os"

	"github.com/Mines-Little-Theatre/team-trivia-scraper/bot"
)

// ReadEnv reads an environment variable, fatally logging an error message if it is not set
func ReadEnv(key string) string {
	var result string
	var ok bool
	if result, ok = os.LookupEnv(key); !ok {
		log.Fatalf("please set the %s environment variable", key)
	}
	return result
}

// ReadConfig reads a bot.Config from standard environment variable names
func ReadConfig() *bot.Config {
	return &bot.Config{
		WebhookID:    ReadEnv("TRIVIA_WEBHOOK_ID"),
		WebhookToken: ReadEnv("TRIVIA_WEBHOOK_TOKEN"),
		Message:      ReadEnv("TRIVIA_MESSAGE"),
		CryForHelp:   os.Getenv("TRIVIA_CRY_FOR_HELP"),
	}
}
