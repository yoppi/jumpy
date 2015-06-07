package jumpy

import (
	"fmt"
	"log"
	"net/url"

	html "github.com/moovweb/gokogiri/html"
)

type Page struct {
	Url *url.URL
	Doc *html.HtmlDocument
}

func NewPage(url *url.URL, doc *html.HtmlDocument) *Page {
	return &Page{url, doc}
}

func (p *Page) Links() []string {
	var links []string

	nodes, err := p.Doc.Search("//a[@href]")
	if err != nil {
		fmt.Printf("%v", err)
		return links
	}

	for _, node := range nodes {
		link := node.Attribute("href").String()
		if normalized, ok := p.Normalize(link); ok {
			links = append(links, normalized)
		}
	}

	return links
}

func (p *Page) Normalize(link string) (string, bool) {
	normalized, err := p.Url.Parse(link)
	if err != nil {
		log.Printf("parse error url[%v]\n", link)
		return "", false
	}

	if normalized.Host != p.Url.Host {
		return "", false
	}

	return normalized.String(), true
}
