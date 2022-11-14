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
	CryForHelp   string // optional
}

func Run(config *Config) error {
	session, err := discordgo.New("")
	if err != nil {
		// reading the source code, it appears that this will never actually return an error
		log.Println(err)
		return err
	}

	data := new(discordgo.WebhookParams)
	data.Content = config.Message

	data.Embeds = embeds.CollectEmbeds(
		embeds.AnswerOfTheNight,
	)

	log.Println("finished collecting embeds, posting")

	_, err = session.WebhookExecute(config.WebhookID, config.WebhookToken, false, data)
	if err != nil {
		log.Println(err)
		// attempt to cry for help
		if config.CryForHelp != "" {
			log.Println("crying for help")
			_, err = session.WebhookExecute(config.WebhookID, config.WebhookToken, false, &discordgo.WebhookParams{
				Content: config.CryForHelp,
			})
			if err != nil {
				return err
			}

			log.Println("finished crying for help")
			return nil
		} else {
			return err
		}
	}

	log.Println("finished posting")
	return nil
}
