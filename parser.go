package jumpy

import(
	. "github.com/ahmetalpbalkan/go-linq"
)

type Parser struct {
	pageCh chan *Page
	linkCh chan string
}

func NewParser(pageCh chan *Page, linkCh chan string) *Parser {
	parser := &Parser{pageCh, linkCh}

	go func() {
		for {
			select {
			case page := <-pageCh:
				links, err := From(page.Links()).Distinct().Results()
				if err != nil {
					continue
				}
				for _, link := range links {
					parser.linkCh <- link.(string)
				}
			}
		}
	}()

	return parser
}
