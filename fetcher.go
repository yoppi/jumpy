package jumpy

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/moovweb/gokogiri"
	html "github.com/moovweb/gokogiri/html"
)

type Fetcher interface {
	Fetch(string) (*Page, error)
}

type fetcher struct {
	Fetcher
	linkCh         chan string
	pageCh         chan *Page
	callbackPageCh chan *Page
	Bucket         *Bucket
}

func NewFetcher(bucket *Bucket, linkCh chan string, pageCh chan *Page, callbackPageCh chan *Page) *fetcher {
	f := &fetcher{
		Bucket:         bucket,
		linkCh:         linkCh,
		pageCh:         pageCh,
		callbackPageCh: callbackPageCh,
	}

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

func (f *fetcher) Fetch(link string) (*Page, error) {
	resp, err := http.Get(link)
	if err != nil {
		fmt.Printf("%v\n", err)
		return nil, err
	}
	defer resp.Body.Close()

	doc, err := f.parseBody(resp)
	if err != nil {
		return nil, err
	}

	linkUrl, _ := url.Parse(link)

	return NewPage(linkUrl, doc), nil
}

func (f *fetcher) parseBody(resp *http.Response) (*html.HtmlDocument, error) {
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("%v\n", err)
		return nil, err
	}

	doc, err := gokogiri.ParseHtml(body)
	if err != nil {
		fmt.Printf("%v\n", err)
		return nil, err
	}

	return doc, nil
}
