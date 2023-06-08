package utils

import (
	"log"
	"os"
	"strings"

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

// ReadEnvList reads a comma-separated list from an environment variable, fatally logging an error message if it is not set
func ReadEnvList(key string) []string {
	var value = ReadEnv(key)
	if value == "" {
		return []string{}
	} else {
		return strings.Split(value, ",")
	}
}

// ReadConfig reads a bot.Config from standard environment variable names
func ReadConfig() *bot.Config {
	return &bot.Config{
		WebhookID:    ReadEnv("TRIVIA_WEBHOOK_ID"),
		WebhookToken: ReadEnv("TRIVIA_WEBHOOK_TOKEN"),
		Suppliers:    ReadEnvList("TRIVIA_SUPPLIERS"),
		Message:      ReadEnv("TRIVIA_MESSAGE"),
		Embeds:       ReadEnvList("TRIVIA_EMBEDS"),
		CryForHelp:   os.Getenv("TRIVIA_CRY_FOR_HELP"),
	}
}
