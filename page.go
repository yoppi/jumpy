package jumpy

import (
	"fmt"
	"log"
	"net/url"
	"path"

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

	nodes, searchErr := p.Doc.Search("//a[@href]")
	if searchErr != nil {
		fmt.Printf("%v", searchErr)
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
	linkUrl, err := url.Parse(link)
	if err != nil {
		log.Printf("parse error url[%v]\n", link)
		return "", false
	}

	if linkUrl.Host == "" {
		if linkUrl.Path == "/" {
			return "", false
		}
		if linkUrl.Path == "" {
			return "", false
		}
		if string(linkUrl.Path[0]) == "/" {
			return p.Url.Scheme + "://" + p.Url.Host + p.NormalizePath(linkUrl.Path), true
		} else {
			pageUrl := p.Url.String()
			if string(pageUrl[len(pageUrl)-1]) != "/" {
				pageUrl += "/"
			}
			return pageUrl + link, true
		}
	}

	if linkUrl.Host != p.Url.Host {
		return "", false
	}

	if linkUrl.Scheme == "javascript" {
		return "", false
	}

	return linkUrl.Scheme +"://" + linkUrl.Host + p.NormalizePath(linkUrl.Path), true
}

func (p *Page) NormalizePath(_path string) string {
	ret := path.Clean(_path)

	if string(_path[len(_path) - 1]) == "/" {
		ret += "/"
	}

	return ret
}
