package htmlparser

import (
	"strings"

	"golang.org/x/net/html"
)

type HtmlParser struct {
	data string
}

func NewParser(data string) *HtmlParser {
	return &HtmlParser{
		data: data,
	}
}

func parse(text string) *tokensList {
	tkn := html.NewTokenizer(strings.NewReader(text))
	tokens := NewTokensList()
	for {
		tt := tkn.Next()
		switch {
		case tt == html.ErrorToken:
			return tokens
		case tt == html.StartTagToken:
			tokens.addToken(tkn.Token().String())
		case tt == html.EndTagToken:
			tokens.closeToken(tkn.Token().String())
		case tt == html.TextToken:
			tokens.text(tkn.Token().Data)
		}
	}
}

func (p *HtmlParser) Body() string {
	tokens := parse(p.data)
	return tokens.body()
}

func (p *HtmlParser) Head() string {
	tokens := parse(p.data)
	return tokens.head()
}

func (p *HtmlParser) Add(comm string) string {
	splited := strings.Split(p.data, "</body>")
	splited[0] += "\n" + comm + "\n"
	return splited[0] + "</body>" + splited[1]
}
