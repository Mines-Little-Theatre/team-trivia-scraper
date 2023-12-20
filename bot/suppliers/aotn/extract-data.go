package aotn

import (
	"strings"

	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

type unit struct{}

func extractData(doc *html.Node) (data freeAnswerData) {
	main := findChildElementNamed(doc, atom.Main)

	titleDiv := findChildNodeWhere(main, func(n *html.Node) bool {
		return n.Type == html.ElementNode && n.DataAtom == atom.Div && findChildNodeWhere(n, func(n *html.Node) bool {
			return n.Type == html.TextNode && strings.ToLower(n.Data) == "answer of the night"
		}) != nil
	})
	data.title = formattedContent(titleDiv)

	blurbDiv := findNextSiblingElementNamed(titleDiv, atom.Div)
	data.blurb = formattedContent(blurbDiv)

	dateDiv := findNextSiblingElementNamed(blurbDiv, atom.Div)
	data.date = formattedContent(dateDiv)

	answerDiv := findNextSiblingElementNamed(dateDiv, atom.Div)
	data.answer = formattedContent(answerDiv)

	return
}

func findChildNodeWhere(n *html.Node, predicate func(*html.Node) bool) *html.Node {
	for at := advanceNodeSearch(n, n); at != nil; at = advanceNodeSearch(at, n) {
		if predicate(at) {
			return at
		}
	}

	return nil
}

func advanceNodeSearch(n *html.Node, topmost *html.Node) *html.Node {
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

func findChildElementNamed(n *html.Node, tagName atom.Atom) *html.Node {
	return findChildNodeWhere(n, func(n *html.Node) bool { return n.Type == html.ElementNode && n.DataAtom == tagName })
}

func findNextSiblingElementNamed(n *html.Node, tagName atom.Atom) *html.Node {
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
