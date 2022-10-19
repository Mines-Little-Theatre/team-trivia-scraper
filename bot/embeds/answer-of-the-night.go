package embeds

import (
	"bytes"
	"log"
	"net/http"

	"github.com/bwmarrin/discordgo"
	"golang.org/x/net/html"
)

const freeAnswerURL = "https://teamtrivia.com/free/"

// AnswerOfTheNight is an EmbedProvider that fetches the contents of https://teamtrivia.com/free/
func AnswerOfTheNight(out chan<- *discordgo.MessageEmbed) {
	defer close(out)

	resp, err := http.Get(freeAnswerURL)
	if err != nil {
		log.Println("free answer: failed to create request:", err)
		return
	}
	defer resp.Body.Close()

	result := new(discordgo.MessageEmbed)
	result.Type = discordgo.EmbedTypeRich
	result.URL = freeAnswerURL
	result.Color = 0x00cccc
	result.Image = new(discordgo.MessageEmbedImage)
	field := new(discordgo.MessageEmbedField)
	result.Fields = []*discordgo.MessageEmbedField{field}

	// title is in an h1.title
	// description is in an h6
	// field name is in an h3>strong
	// field value is in an h2
	// (the terrible header structure is very convenient for this specific purpose)
	// and the image url is in an img.img-fluid[src]

	z := html.NewTokenizer(resp.Body)

htmlLoop:
	for anyEmptyStrings(result.Title, field.Name, field.Value, result.Image.URL) {
		switch z.Next() {
		case html.ErrorToken:
			break htmlLoop
		case html.StartTagToken:
			tagName, hasAttr := z.TagName()
			switch {
			case bytes.Equal(tagName, []byte("h1")):
				// don't need to check class, there is only one h1
				if z.Next() != html.TextToken {
					log.Println("free answer: no text in the h1")
				} else {
					result.Title = string(z.Text())
				}
			case bytes.Equal(tagName, []byte("h3")):
				if z.Next() != html.StartTagToken {
					log.Println("free answer: no strong in the h3")
				} else if z.Next() != html.TextToken {
					log.Println("free answer: no text in the strong in the h3")
				} else {
					field.Name = string(z.Text())
				}
			case bytes.Equal(tagName, []byte("h2")):
				if z.Next() != html.TextToken {
					log.Println("free answer: no text in the h2")
				} else {
					field.Value = string(z.Text())
				}
			case bytes.Equal(tagName, []byte("img")):
				// need to find src and check class
				var key, val []byte
				var src string
				var hasFluidClass bool
			attrLoop:
				for hasAttr && (src == "" || !hasFluidClass) {
					key, val, hasAttr = z.TagAttr()
					switch {
					case bytes.Equal(key, []byte("src")):
						src = string(val)
					case bytes.Equal(key, []byte("class")):
						classes := bytes.Split(val, []byte(" "))
						for _, class := range classes {
							if bytes.Equal(class, []byte("img-fluid")) {
								hasFluidClass = true
								continue attrLoop
							}
						}
					}
				}
				if hasFluidClass {
					result.Image.URL = src
				}
			}
		}
	}

	if field.Value != "" {
		if field.Name == "" {
			// field name is required, and we definitely want the value we have
			// so put in some placeholder value that's obviously an error
			log.Println("free answer: no name for field")
			field.Name = "undefined"
		}

		log.Println("free answer: the free answer is", field.Value)
		out <- result
	} else {
		log.Println("free answer: no free answer")
	}
}

func anyEmptyStrings(strings ...string) bool {
	for _, s := range strings {
		if s == "" {
			return true
		}
	}
	return false
}
