package aotn

import (
	"fmt"
	"log"
	"net/http"

	"github.com/Mines-Little-Theatre/team-trivia-scraper/bot/suppliers"
	"github.com/bwmarrin/discordgo"
	"golang.org/x/net/html"
)

const freeAnswerURL = "https://teamtrivia.com/free/"

type freeAnswerData struct {
	title, blurb, date, answer string
}

type AnswerOfTheNight struct{}

func init() {
	suppliers.RegisterSupplier("aotn", AnswerOfTheNight{})
}

func (AnswerOfTheNight) SupplyData(context *suppliers.SupplierContext) error {
	doc, err := fetchDocument(context.Config("REGION_ID"))
	if err != nil {
		return err
	}

	data := extractData(doc)
	embed := createEmbed(data)

	if data.answer != "" {
		openaiToken := context.Config("OPENAI_TOKEN")
		if openaiToken != "" {
			imageUrl, err := generateImage(data.answer, openaiToken)
			if err != nil {
				log.Println("aotn generate image:", err)
				embed.Footer = &discordgo.MessageEmbedFooter{
					Text: "Image generation failed: " + err.Error(),
				}
			} else {
				embed.Image = &discordgo.MessageEmbedImage{
					URL:    imageUrl,
					Width:  256,
					Height: 256,
				}
				embed.Footer = &discordgo.MessageEmbedFooter{
					Text: "Image is AI-generated",
				}
			}
		}
	}

	context.AddEmbed("answer", embed)
	return nil
}

func fetchDocument(regionID string) (*html.Node, error) {
	req, err := http.NewRequest("GET", freeAnswerURL, nil)
	if err != nil {
		return nil, err
	}

	if regionID != "" {
		req.AddCookie(&http.Cookie{Name: "region_ID", Value: regionID})
	}
	req.AddCookie(&http.Cookie{Name: "new_site", Value: "Y"})

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	} else if resp.StatusCode != 200 {
		resp.Body.Close()
		return nil, fmt.Errorf("status %d from free answer page", resp.StatusCode)
	}

	doc, err := html.Parse(resp.Body)
	resp.Body.Close()
	return doc, err
}

func createEmbed(data freeAnswerData) *discordgo.MessageEmbed {
	result := &discordgo.MessageEmbed{
		Title: data.title,
		// Description: data.blurb,
		URL:   freeAnswerURL,
		Color: 0x0069b5,
	}

	if data.date != "" || data.answer != "" {
		result.Fields = []*discordgo.MessageEmbedField{{
			Name:  data.date,
			Value: data.answer,
		}}
	}

	return result
}
