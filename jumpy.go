package jumpy

import (
	"net/url"
	"runtime"
)

var options map[string]string

func setup() {
	numCpu := runtime.NumCPU()
	if runtime.GOMAXPROCS(0) < numCpu {
		runtime.GOMAXPROCS(numCpu)
	}
}

func Crawl(rootUrl string, options map[string]string, callback func(*Page)) {
	root, err := url.Parse(rootUrl)
	if err != nil {
		panic(err)
	}

	setup()
	// TODO: Handle options
	options = options

	callbackPageCh := make(chan *Page, 10000)
	stop := make(chan int)
	crawler := NewCrawler(root, callbackPageCh, stop)

	crawler.crawl()

	for {
		select {
		case page := <-callbackPageCh:
			callback(page)
			crawler.pageCh <- page
		case <-stop:
			return
		}
	}
}
