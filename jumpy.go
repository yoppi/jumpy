package jumpy

import (
	"net/url"
)

var options map[string]string

// Start crawling at rootUrl
func Crawl(rootUrl string, options map[string]string, callback func(*Page)) {
	root, err := url.Parse(rootUrl)
	if err != nil {
		panic(err)
	}

	// TODO: Handle options
	options = options

	callbackPageCh := make(chan *Page, 10000)
	stop := make(chan int)
	crawler := NewCrawler(root, callbackPageCh, stop)

	crawler.crawl()

	for {
		select {
		case page := <-callbackPageCh:
			go callback(page)
			crawler.pageCh <- page
		case <-stop:
			return
		}
	}
}
