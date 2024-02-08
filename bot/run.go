package bot

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"image/jpeg"
	"image/png"
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
			err = cryForHelp(ctx, session, webhookID, webhookToken)
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
			pngData, err := dalle.GenerateImage(ctx, answerData.Answer)
			if err != nil && !errors.Is(err, dalle.ErrNoToken) {
				log.Println("generate image:", err)
				embed.Footer = &discordgo.MessageEmbedFooter{
					Text: "Image generation failed: " + err.Error(),
				}
			} else if err == nil {
				image, err := png.Decode(bytes.NewReader(pngData))
				if err != nil {
					log.Println("decode image:", err)
					embed.Footer = &discordgo.MessageEmbedFooter{
						Text: "Image generation failed (decoding): " + err.Error(),
					}
				} else {
					jpegBuf := new(bytes.Buffer)
					err := jpeg.Encode(jpegBuf, image, nil)
					if err != nil {
						log.Println("encode image:", err)
						embed.Footer = &discordgo.MessageEmbedFooter{
							Text: "Image generation failed (encoding): " + err.Error(),
						}
					} else {
						log.Printf("JPEG encoding reduced image size from %d to %d (%+.1f%%)", len(pngData), jpegBuf.Len(), 100.0*float64(jpegBuf.Len()-len(pngData))/float64(len(pngData)))
						webhookMessage.Files = []*discordgo.File{{
							Name:        "image.jpg",
							ContentType: "image/jpeg",
							Reader:      jpegBuf,
						}}
						embed.Image = &discordgo.MessageEmbedImage{
							URL: "attachment://image.jpg",
						}
						embed.Footer = &discordgo.MessageEmbedFooter{
							Text: "Image is AI-generated (DALL·E 3)",
						}
					}
				}
			}
		}
	}

	_, err = session.WebhookExecute(webhookID, webhookToken, true, webhookMessage, discordgo.WithContext(ctx))
	if err != nil {
		log.Println("execute webhook:", err)
		return cryForHelp(ctx, session, webhookID, webhookToken)
	}

	log.Println("finished posting")
	return nil
}

func cryForHelp(ctx context.Context, session *discordgo.Session, webhookID, webhookToken string) error {
	distressMessage, err := readEnv("TRIVIA_BOT_CRY_FOR_HELP")
	if err != nil {
		return err
	}
	_, err = session.WebhookExecute(webhookID, webhookToken, true, &discordgo.WebhookParams{
		Content: distressMessage,
	}, discordgo.WithContext(ctx))
	return err
}
