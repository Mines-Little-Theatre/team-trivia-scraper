package aotn

import (
	"fmt"
	"net/http"

	"github.com/bwmarrin/discordgo"
	"golang.org/x/net/html"
)

const freeAnswerURL = "https://teamtrivia.com/free/"

type freeAnswerData struct {
	title, blurb, date, answer, imageURL string
}

// AnswerOfTheNight is an EmbedProvider that fetches the contents of https://teamtrivia.com/free/
func AnswerOfTheNight() (*discordgo.MessageEmbed, error) {
	doc, err := fetchDocument()
	if err != nil {
		return nil, err
	}

	data := extractData(doc)

	return createEmbed(data), nil
}

func fetchDocument() (*html.Node, error) {
	resp, err := http.Get(freeAnswerURL)
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
	result := new(discordgo.MessageEmbed)
	result.Type = discordgo.EmbedTypeRich

	result.Title = data.title
	// result.Description = data.blurb
	result.URL = freeAnswerURL
	result.Color = 0x00cccc

	if data.date != "" || data.answer != "" {
		result.Fields = []*discordgo.MessageEmbedField{{
			Name:  data.date,
			Value: data.answer,
		}}
	}

	if data.imageURL != "" {
		result.Image = &discordgo.MessageEmbedImage{
			URL: data.imageURL,
		}
	}

	return result
}
