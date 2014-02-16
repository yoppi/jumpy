package jumpy

import (
	"log"
	"net/url"
)

type Crawler struct {
	RootUrl *url.URL
	Fetcher Fetcher
	Bucket *Bucket
	callbackPageCh chan *Page
	pageCh chan *Page
	linkCh chan string
	stop chan int
}

func NewCrawler(rootUrl *url.URL, callbackPageCh chan *Page, stop chan int) *Crawler {
	pageCh := make(chan *Page, 10000)
	linkCh := make(chan string, 10000)
	crawler := &Crawler{rootUrl, NewFetcher(), NewBucket(), callbackPageCh, pageCh, linkCh, stop}

	go func() {
		for {
			select {
			case page := <-crawler.pageCh:
				for _, link := range page.Links() {
					if !crawler.Bucket.Exist(link) {
						crawler.linkCh <- link
					}
				}
			case link := <-crawler.linkCh:
				if page, ok := crawler.fetch(link); ok {
					crawler.Bucket.Add(link, page)
					crawler.callbackPageCh <- page
				}
			}
		}
	}()

	return crawler
}

func (c *Crawler) crawl() {
	log.Printf("Crawl starting at [%v]\n", c.RootUrl.String())
	c.linkCh <- c.RootUrl.String()
}

func (c *Crawler) fetch(link string) (*Page, bool) {
	page, err := c.Fetcher.Fetch(link)
	if err != nil {
		return nil, false
	}
	return page, true
}

