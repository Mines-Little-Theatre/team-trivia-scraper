package bot

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/Mines-Little-Theatre/team-trivia-scraper/bot/answer"
	"github.com/Mines-Little-Theatre/team-trivia-scraper/bot/dalle"
	"github.com/bwmarrin/discordgo"
)

func readEnv(key string) (string, error) {
	value, ok := os.LookupEnv(key)
	if !ok {
		return "", fmt.Errorf("%s environment variable not set", key)
	}
	return value, nil
}

func Run(ctx context.Context) (err error) {
	session, err := discordgo.New("")
	if err != nil {
		return fmt.Errorf("create discord session: %w", err)
	}
	defer session.Close()

	webhookID, err := readEnv("TRIVIA_BOT_WEBHOOK_ID")
	if err != nil {
		return err
	}
	webhookToken, err := readEnv("TRIVIA_BOT_WEBHOOK_TOKEN")
	if err != nil {
		return err
	}

	defer func() {
		if panicValue := recover(); panicValue != nil {
			log.Println("panicked:", panicValue)
			distressMessage, ok := os.LookupEnv("TRIVIA_BOT_CRY_FOR_HELP")
			if ok {
				log.Println("crying for help")
			}
			_, err = session.WebhookExecute(webhookID, webhookToken, true, &discordgo.WebhookParams{
				Content: distressMessage,
			}, discordgo.WithContext(ctx))
		}
	}()

	webhookMessage := new(discordgo.WebhookParams)
	webhookMessage.Content = os.Getenv("TRIVIA_BOT_MESSAGE")

	answerData, err := answer.Fetch(ctx)
	if err != nil {
		webhookMessage.Embeds = []*discordgo.MessageEmbed{{
			Description: "Failed to retrieve the free answer: " + err.Error(),
			Color:       0xffcc00,
		}}
	} else {
		embed := &discordgo.MessageEmbed{
			Title: answerData.Title,
			URL:   answer.FreeAnswerURL,
			Color: 0x0069b5,
		}
		webhookMessage.Embeds = []*discordgo.MessageEmbed{embed}

		if answerData.Date != "" || answerData.Answer != "" {
			embed.Fields = append(embed.Fields, &discordgo.MessageEmbedField{
				Name:  answerData.Date,
				Value: answerData.Answer,
			})
		}

		if answerData.Answer != "" {
			imageURL, err := dalle.GenerateImage(ctx, answerData.Answer)
			if err != nil && !errors.Is(err, dalle.ErrNoToken) {
				log.Println("generate image:", err)
				embed.Footer = &discordgo.MessageEmbedFooter{
					Text: "Image generation failed: " + err.Error(),
				}
			} else if err == nil {
				embed.Image = &discordgo.MessageEmbedImage{
					URL:    imageURL,
					Width:  256,
					Height: 256,
				}
				embed.Footer = &discordgo.MessageEmbedFooter{
					Text: "Image is AI-generated",
				}
			}
		}
	}

	_, err = session.WebhookExecute(webhookID, webhookToken, true, webhookMessage, discordgo.WithContext(ctx))
	if err != nil {
		panic(err)
	}

	log.Println("finished posting")
	return nil
}
