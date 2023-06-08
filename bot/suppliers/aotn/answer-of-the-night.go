package aotn

import (
	"fmt"
	"net/http"

	"github.com/Mines-Little-Theatre/team-trivia-scraper/bot/suppliers"
	"github.com/bwmarrin/discordgo"
	"golang.org/x/net/html"
)

const freeAnswerURL = "https://teamtrivia.com/free/"

type freeAnswerData struct {
	title, blurb, date, answer, imageURL string
}

type AnswerOfTheNight struct{}

func init() {
	suppliers.RegisterSupplier("aotn", AnswerOfTheNight{})
}

func (a AnswerOfTheNight) SupplyData(context *suppliers.SupplierContext) error {
	doc, err := fetchDocument(context.Config("REGION_ID"))
	if err != nil {
		return err
	}

	data := extractData(doc)
	context.AddEmbed("answer", createEmbed(data))
	return nil
}

func fetchDocument(regionID string) (*html.Node, error) {
	req, err := http.NewRequest("GET", freeAnswerURL, nil)
	if err != nil {
		return nil, err
	}

	if regionID != "" {
		req.Header.Add("Cookie", "region_ID="+regionID)
	}

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
