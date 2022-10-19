package main

import (
	"github.com/Mines-Little-Theatre/team-trivia-scraper/bot"
	"github.com/Mines-Little-Theatre/team-trivia-scraper/utils"
)

func main() {
	bot.Run(bot.Config{
		WebhookID:    utils.ReadEnv("TRIVIA_WEBHOOK_ID"),
		WebhookToken: utils.ReadEnv("TRIVIA_WEBHOOK_TOKEN"),
		Message:      utils.ReadEnv("TRIVIA_MESSAGE"),
	})
}
