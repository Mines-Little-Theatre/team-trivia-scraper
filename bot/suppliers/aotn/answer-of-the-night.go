package aotn

import (
	"fmt"
	"net/http"
	"time"

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

func (AnswerOfTheNight) SupplyData(context *suppliers.SupplierContext) error {
	doc, err := fetchDocument(context.Config("REGION_ID"))
	if err != nil {
		return err
	}

	data := extractData(doc)
	today := time.Now()
	data.imageURL = fmt.Sprintf("https://www.triviocity.com/game/inserts/%04[1]d/%02[2]d/%04[1]d-%02[2]d-%02[3]d/1/r4_1.jpg", today.Year(), today.Month(), today.Day())
	context.AddEmbed("answer", createEmbed(data))
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

	if data.imageURL != "" {
		result.Image = &discordgo.MessageEmbedImage{
			URL: data.imageURL,
		}
	}

	return result
}
