package aotn

import (
	"log"
	"net/url"
	"strings"

	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

type unit struct{}

func extractData(doc *html.Node) freeAnswerData {
	// title is the content of the only h1
	// blurb is the content of the only h6
	// date is the content of the only h3
	// answer is the content of the only h2
	// (the terrible header structure is very convenient for this)
	// and the image url is the src of the only img.img-fluid, which only has the one class

	var data freeAnswerData
	extractDataRecursive(&data, doc)
	return data
}

// return true if done
func extractDataRecursive(data *freeAnswerData, n *html.Node) bool {
	if n.Type == html.ElementNode {
		switch n.DataAtom {
		case atom.H1:
			data.title = formattedContent(n)
			return !anyEmptyStrings(data.blurb, data.date, data.answer, data.imageURL)
		case atom.H6:
			data.blurb = formattedContent(n)
			return !anyEmptyStrings(data.title, data.date, data.answer, data.imageURL)
		case atom.H3:
			data.date = formattedContent(n)
			return !anyEmptyStrings(data.title, data.blurb, data.answer, data.imageURL)
		case atom.H2:
			data.answer = formattedContent(n)
			return !anyEmptyStrings(data.title, data.blurb, data.date, data.imageURL)
		case atom.Img:
			hasFluidClass := false
			src := ""
			for _, attr := range n.Attr {
				switch attr.Key {
				case "class":
					hasFluidClass = attr.Val == "img-fluid"
				case "src":
					src = attr.Val
				}
			}

			if hasFluidClass {
				src, err := url.JoinPath(freeAnswerURL, src)
				if err != nil {
					log.Println("free answer: could not make sense of src", src)
				} else {
					data.imageURL = src
					return !anyEmptyStrings(data.title, data.blurb, data.date, data.answer)
				}
			}
		}
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		if extractDataRecursive(data, c) {
			return true
		}
	}

	return false
}

var formatTokens = map[atom.Atom]string{
	atom.Em:     "_",
	atom.Strong: "**",
}

func formattedContent(n *html.Node) string {
	var out strings.Builder
	currentFormats := make(map[string]unit, len(formatTokens))
	formatContentRecursive(&out, currentFormats, n)
	return out.String()
}

func formatContentRecursive(out *strings.Builder, currentFormats map[string]unit, n *html.Node) {
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		switch c.Type {
		case html.TextNode:
			out.WriteString(c.Data)
		case html.ElementNode:
			format, ok := formatTokens[c.DataAtom]
			if ok {
				_, already := currentFormats[format]
				ok = !already
			}

			if ok {
				currentFormats[format] = unit{}
				out.WriteString(format)
			}
			for c := n.FirstChild; c != nil; c = c.NextSibling {
				formatContentRecursive(out, currentFormats, c)
			}
			if ok {
				delete(currentFormats, format)
				out.WriteString(format)
			}
		}
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
