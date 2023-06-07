package bot

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/Mines-Little-Theatre/team-trivia-scraper/bot/suppliers"
	_ "github.com/Mines-Little-Theatre/team-trivia-scraper/bot/suppliers/aotn"
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

	supplierResults := suppliers.RunSuppliers([]string{"aotn"})

	data := new(discordgo.WebhookParams)
	data.Content = config.Message

	embed, ok := supplierResults.Embeds["aotn:answer"]
	if ok {
		embedData, _ := json.Marshal(embed)
		fmt.Println(string(embedData))
		data.Embeds = append(data.Embeds, embed)
	}

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
