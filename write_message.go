package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"strings"

	"golang.org/x/net/html"
)

const url = "https://teamtrivia.com/free/"

func writeMessage(config *Config, bodyWriter *io.PipeWriter, mpartWriter *multipart.Writer) {
	resp, err := http.Get(url)
	if err != nil {
		log.Fatalln("failed to create scraping request:", err)
	}
	defer resp.Body.Close()

	z := html.NewTokenizer(resp.Body)

	content, err := mpartWriter.CreateFormField("content")
	if err != nil {
		log.Fatalln("failed to create content part:", err)
	}

	// find first/only h2, that has the answer or is empty
	for {
		switch z.Next() {
		case html.ErrorToken:
			log.Println("did not find an h2")
			_, err := fmt.Fprint(content, config.NoFreeAnswerMessage)
			if err != nil {
				log.Fatalln("failed to write content:", err)
			}
			err = mpartWriter.Close()
			if err != nil {
				log.Fatalln("failed to close mpartWriter:", err)
			}
			return
		case html.StartTagToken:
			if tagName, _ := z.TagName(); bytes.Equal(tagName, []byte("h2")) {
				// get the text; if the next token is not text, there isn't a free answer (yet?)
				if z.Next() != html.TextToken {
					log.Println("no text in the h2")
					_, err := fmt.Fprint(content, config.NoFreeAnswerMessage)
					if err != nil {
						log.Fatalln("failed to write content:", err)
					}
					err = mpartWriter.Close()
					if err != nil {
						log.Fatalln("failed to close mpartWriter:", err)
					}
				} else {
					_, err := fmt.Fprintf(content, config.FreeAnswerMessage, z.Text())
					if err != nil {
						log.Fatalln("failed to write content:", err)
					}

					// from here on, if anything goes wrong, just log the error and continue
					// we've written the answer, we can at least try to close mpartWriter and bodyWriter

					// find the image?
					err = findImage(z, mpartWriter)
					if err != nil {
						log.Println(err)
					}

					resp.Body.Close()
					err = mpartWriter.Close()
					if err != nil {
						log.Fatalln("could not close mpartWriter:", err)
					}
					err = bodyWriter.Close()
					if err != nil {
						log.Fatalln("could not close bodyWriter:", err)
					}
				}

				return
			}
		}
	}
}

func findImage(z *html.Tokenizer, mpartWriter *multipart.Writer) error {
	for {
		switch z.Next() {
		case html.ErrorToken:
			// no image, oh well
			return nil
		case html.EndTagToken:
			// we might be exiting the div where the image would be
			if tagName, _ := z.TagName(); bytes.Equal(tagName, []byte("div")) {
				// no image, oh well
				return nil
			}
		case html.StartTagToken, html.SelfClosingTagToken:
			// this is maybe the image?
			if tagName, hasAttr := z.TagName(); bytes.Equal(tagName, []byte("img")) {
				// found it!
				for hasAttr {
					var key, val []byte // have to declare them here so hasAttr doesn't get redeclared
					key, val, hasAttr = z.TagAttr()
					if bytes.Equal(key, []byte("src")) {
						// found the image URL!
						url := string(val)
						lastSlash := strings.LastIndexByte(url, '/')
						filename := url[lastSlash+1:]

						resp, err := http.Get(url)
						if err != nil {
							return fmt.Errorf("failed to request image: %w", err)
						}
						defer resp.Body.Close()

						if resp.StatusCode != 200 {
							return fmt.Errorf("image request failed with status %d", resp.StatusCode)
						}

						file, err := mpartWriter.CreateFormFile("files[0]", filename)
						if err != nil {
							return fmt.Errorf("failed to create file part: %w", err)
						}

						_, err = io.Copy(file, resp.Body)
						if err != nil {
							return fmt.Errorf("failed to copy image file: %w", err)
						}
					}
				}
				return nil
			}
		}
	}
}
