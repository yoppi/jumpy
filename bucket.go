package jumpy

import (
	"sync"
)

type Bucket struct {
	Storage Storage
}

type Storage interface {
	Exist(string) bool
	Add(string, *Page)
	Keys() []string
}

type OnMemory struct {
	sync.Mutex
	Storage
	Pages map[string]*Page
}

func NewOnMemory() *OnMemory {
	return &OnMemory{Pages: map[string]*Page{}}
}

func (s *OnMemory) Exist(url string) bool {
	_, ok := s.Pages[url]
	return ok
}

func (s *OnMemory) Add(url string, page *Page) {
	s.Lock()
	s.Pages[url] = page
	s.Unlock()
}

func (s *OnMemory) Keys() []string {
	keys := make([]string, len(s.Pages))
	for k, _ := range s.Pages {
		keys = append(keys, k)
	}
	return keys
}

func NewBucket() *Bucket {
	storage := NewOnMemory()
	return &Bucket{storage}
}

func (b *Bucket) Exist(url string) bool{
	return b.Storage.Exist(url)
}

func (b *Bucket) Add(url string, page *Page) {
	b.Storage.Add(url, page)
}


