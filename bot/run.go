package bot

import (
	"log"

	"github.com/Mines-Little-Theatre/team-trivia-scraper/bot/embeds"
	"github.com/bwmarrin/discordgo"
)

type Config struct {
	WebhookID    string
	WebhookToken string
	Message      string
}

func Run(config Config) {
	session, err := discordgo.New("")
	if err != nil {
		// reading the source code, it appears that this will never actually return an error
		log.Fatalln(err)
	}

	data := new(discordgo.WebhookParams)
	data.Content = config.Message

	data.Embeds = embeds.CollectEmbeds(
		embeds.AnswerOfTheNight,
	)

	_, err = session.WebhookExecute(config.WebhookID, config.WebhookToken, false, data)
	if err != nil {
		log.Fatalln(err)
	}
}
