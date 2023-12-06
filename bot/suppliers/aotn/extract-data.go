package aotn

import (
	"strings"

	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

type unit struct{}

func extractData(doc *html.Node) (data freeAnswerData) {
	main := findChildElementNamed(atom.Main, doc)
	section := findChildElementNamed(atom.Section, main)
	flexContainer := findChildElementNamed(atom.Div, section)
	adSpaceDiv := findChildElementNamed(atom.Div, flexContainer)

	titleDiv := findNextSiblingElementNamed(atom.Div, adSpaceDiv)
	data.title = formattedContent(titleDiv)

	blurbDiv := findNextSiblingElementNamed(atom.Div, titleDiv)
	data.blurb = formattedContent(blurbDiv)

	dateDiv := findNextSiblingElementNamed(atom.Div, blurbDiv)
	data.date = formattedContent(dateDiv)

	answerDiv := findNextSiblingElementNamed(atom.Div, dateDiv)
	data.answer = formattedContent(answerDiv)

	return
}

func findChildElementNamed(tagName atom.Atom, n *html.Node) *html.Node {
	for at := advanceElementSearch(n, n); at != nil; at = advanceElementSearch(at, n) {
		if at.Type == html.ElementNode && at.DataAtom == tagName {
			return at
		}
	}

	return nil
}

func advanceElementSearch(n *html.Node, topmost *html.Node) *html.Node {
	if n.FirstChild != nil {
		return n.FirstChild
	} else if n != topmost {
		if n.NextSibling != nil {
			return n.NextSibling
		} else if n != topmost {
			for n := n.Parent; n != topmost && n != nil; n = n.Parent {
				if n.NextSibling != nil {
					return n.NextSibling
				}
			}
		}
	}

	return nil
}

func findNextSiblingElementNamed(tagName atom.Atom, n *html.Node) *html.Node {
	for at := n.NextSibling; at != nil; at = at.NextSibling {
		if at.Type == html.ElementNode && at.DataAtom == tagName {
			return at
		}
	}

	return nil
}

var formatTokens = map[atom.Atom]string{
	atom.Em:     "_",
	atom.Strong: "**",
}

func formattedContent(n *html.Node) string {
	var out strings.Builder
	currentFormats := make(map[string]unit, len(formatTokens))
	formatContentRecursive(&out, currentFormats, n)
	return strings.TrimSpace(out.String())
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
