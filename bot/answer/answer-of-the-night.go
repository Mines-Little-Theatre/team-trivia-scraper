package answer

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"golang.org/x/net/html"
)

const FreeAnswerURL = "https://teamtrivia.com/free/"

type AnswerData struct {
	Title, Blurb, Date, Answer string
}

func Fetch(ctx context.Context) (*AnswerData, error) {
	doc, err := fetchDocument(ctx)
	if err != nil {
		return nil, fmt.Errorf("fetch document: %w", err)
	}

	return extractData(doc)
}

func fetchDocument(ctx context.Context) (*html.Node, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", FreeAnswerURL, nil)
	if err != nil {
		return nil, err
	}

	req.AddCookie(&http.Cookie{Name: "new_site", Value: "Y"})
	if regionID := os.Getenv("TRIVIA_BOT_REGION_ID"); regionID != "" {
		req.AddCookie(&http.Cookie{Name: "region_ID", Value: regionID})
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
