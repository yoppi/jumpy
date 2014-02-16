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
	nodes, searchErr := p.Doc.Search("//a[@href]")
	if searchErr != nil {
		fmt.Errorf("%v", searchErr)
	}

	var links []string
	for _, node := range nodes {
		link := node.Attribute("href").String()
		if normalized, ok := p.Normalize(link); ok {
			links = append(links, normalized)
		}
	}

	return links
}

func (p *Page) Normalize(link string) (string, bool) {
	linkUrl, err := url.Parse(link)
	if err != nil {
		log.Printf("parse error url[%v]\n", link)
		return "", false
	}

	if linkUrl.Host == "" {
		if linkUrl.Path == "/" {
			//log.Printf("root url\n")
			return "", false
		}
		if linkUrl.Path == "" {
			//log.Printf("invalid url[%v]", link)
			return "", false
		}
		if string(linkUrl.Path[0]) == "/" {
			return p.Url.Scheme + "://" + p.Url.Host + linkUrl.Path, true
		} else {
			pageUrl := p.Url.String()
			if string(pageUrl[len(pageUrl)-1]) != "/" {
				pageUrl += "/"
			}
			return pageUrl + link, true
		}
	}

	if linkUrl.Host != p.Url.Host {
		//log.Printf("outside domain[%v]\n", link)
		return "", false
	}

	if linkUrl.Scheme == "javascript" {
		//log.Printf("invalid url[%v]", link)
		return "", false
	}

	return linkUrl.Scheme +"://" + linkUrl.Host + linkUrl.Path, true
}

