package jumpy

import (
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/moovweb/gokogiri"
)

type Fetcher interface {
	Fetch(string) (*Page, error)
}

type fetcher struct {
	Fetcher
	linkCh chan string
	pageCh chan *Page
	callbackPageCh chan *Page
	Bucket *Bucket
}

func NewFetcher(bucket *Bucket, linkCh chan string, pageCh chan *Page, callbackPageCh chan *Page) *fetcher {
	f := &fetcher{Bucket: bucket, linkCh: linkCh, pageCh: pageCh, callbackPageCh: callbackPageCh}

	go func() {
		for {
			select {
			case link := <-linkCh:
				if !f.Bucket.Exist(link) {
					page, err := f.Fetch(link)
					if err != nil {
						return
					}
					f.Bucket.Add(link, page)
					f.callbackPageCh <- page
				}
			}
		}
	}()

	return f
}

func (f *fetcher) Fetch(url string) (*Page, error) {
	resp, err := http.Get(url)
	if err != nil {
		fmt.Printf("%v\n", err)
		return nil, err
	}
	defer resp.Body.Close()

	return f.parseBody(resp)
}

func (f *fetcher) parseBody(resp *http.Response) (*Page, error) {
	body, readErr := ioutil.ReadAll(resp.Body)
	if readErr != nil {
		fmt.Printf("%v\n", readErr)
		return nil, readErr
	}

	doc, parseErr := gokogiri.ParseHtml(body)
	if parseErr != nil {
		fmt.Printf("%v\n", parseErr)
		return nil, parseErr
	}

	return NewPage(resp.Request.URL, doc), nil
}

