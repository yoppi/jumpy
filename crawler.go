package jumpy

import (
	"log"
	"net/url"
	"runtime"
)

type Crawler struct {
	RootUrl *url.URL
	Fetchers []Fetcher
	Parser *Parser
	callbackPageCh chan *Page
	pageCh chan *Page
	linkCh chan string
	stop chan int
}

func NewCrawler(rootUrl *url.URL, callbackPageCh chan *Page, stop chan int) *Crawler {
	pageCh := make(chan *Page, 10000)
	linkCh := make(chan string, 10000)
	backet := NewBucket()
	fetchers := make([]Fetcher, runtime.GOMAXPROCS(0))
	for i := 0; i < runtime.GOMAXPROCS(0); i++ {
		fetchers = append(fetchers, NewFetcher(backet, linkCh, pageCh, callbackPageCh))
	}
	parser := NewParser(backet, pageCh, linkCh)

	return &Crawler{rootUrl, fetchers, parser, callbackPageCh, pageCh, linkCh, stop}
}

func (c *Crawler) crawl() {
	log.Printf("Crawl starting at [%v]\n", c.RootUrl.String())
	c.linkCh <- c.RootUrl.String()
}
