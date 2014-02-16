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
}

func NewFetcher() *fetcher {
	return &fetcher{}
}

func (f *fetcher) Fetch(url string) (*Page, error) {
	resp, err := http.Get(url)
	if err != nil {
		fmt.Errorf("%v", err)
	}
	defer resp.Body.Close()

	return f.parseBody(resp)
}

func (f *fetcher) parseBody(resp *http.Response) (*Page, error) {
	body, readErr := ioutil.ReadAll(resp.Body)
	if readErr != nil {
		fmt.Errorf("%v", readErr)
		return nil, readErr
	}

	doc, parseErr := gokogiri.ParseHtml(body)
	if parseErr != nil {
		fmt.Errorf("%v", parseErr)
		return nil, parseErr
	}

	return NewPage(resp.Request.URL, doc), nil
}

