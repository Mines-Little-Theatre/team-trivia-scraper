package aotn

import (
	"log"
	"net/url"
	"strings"

	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

type unit struct{}

func extractData(doc *html.Node) (data freeAnswerData) {
	// title is the content of the first (and only) h1
	h1 := findNextElementNamed(atom.H1, doc)
	if h1 == nil {
		return
	}
	data.title = formattedContent(h1)

	// blurb is the content of the subsequent h6
	h6 := findNextElementNamed(atom.H6, h1)
	if h6 == nil {
		return
	}
	data.blurb = formattedContent(h6)

	// date is the content of the subsequent h3
	h3 := findNextElementNamed(atom.H3, h6)
	if h3 == nil {
		return
	}
	data.date = formattedContent(h3)

	// answer is the content of the subsequent h2
	h2 := findNextElementNamed(atom.H2, h3)
	if h2 == nil {
		return
	}
	data.answer = formattedContent(h2)

	// (the terrible header structure is very convenient for this)
	// and the image url is the src of the subsequent img if it (only) has the img-fluid class
	img := findNextElementNamed(atom.Img, h2)
	if img == nil {
		return
	}
	hasFluidClass := false
	src := ""
	for _, attr := range img.Attr {
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
			log.Println("aotn: could not make sense of src:", src)
		} else {
			data.imageURL = src
		}
	}

	return
}

func findNextElementNamed(tagName atom.Atom, n *html.Node) *html.Node {
	for n := advanceElementSearch(n); n != nil; n = advanceElementSearch(n) {
		if n.Type == html.ElementNode && n.DataAtom == tagName {
			return n
		}
	}

	return nil
}

func advanceElementSearch(n *html.Node) *html.Node {
	if n.FirstChild != nil {
		return n.FirstChild
	} else if n.NextSibling != nil {
		return n.NextSibling
	} else {
		for n := n.Parent; n != nil; n = n.Parent {
			if n.NextSibling != nil {
				return n.NextSibling
			}
		}
		return nil
	}
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
