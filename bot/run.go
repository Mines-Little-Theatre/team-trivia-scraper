package bot

import (
	"context"
	"log"

	"github.com/Mines-Little-Theatre/team-trivia-scraper/bot/suppliers"
	_ "github.com/Mines-Little-Theatre/team-trivia-scraper/bot/suppliers/aotn"
	_ "github.com/Mines-Little-Theatre/team-trivia-scraper/bot/suppliers/errsupply"
	"github.com/bwmarrin/discordgo"
)

type Config struct {
	WebhookID    string
	WebhookToken string
	Suppliers    []string
	Message      string
	Embeds       []string
	CryForHelp   string // optional
}

func Run(ctx context.Context, config *Config) error {
	session, err := discordgo.New("")
	if err != nil {
		// reading the source code, it appears that this will never actually return an error
		log.Println(err)
		return err
	}

	supplierResults := suppliers.RunSuppliers(config.Suppliers)

	data := new(discordgo.WebhookParams)
	data.Content = config.Message
	for _, embedName := range config.Embeds {
		embed, ok := supplierResults.Embeds[embedName]
		if ok {
			data.Embeds = append(data.Embeds, embed)
		}
	}

	log.Println("finished collecting supplier results, posting")

	_, err = session.WebhookExecute(config.WebhookID, config.WebhookToken, true, data, discordgo.WithContext(ctx))
	if err != nil {
		log.Println(err)
		// attempt to cry for help
		if config.CryForHelp != "" {
			log.Println("crying for help")
			_, err = session.WebhookExecute(config.WebhookID, config.WebhookToken, true, &discordgo.WebhookParams{
				Content: config.CryForHelp,
			}, discordgo.WithContext(ctx))
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
