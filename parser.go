package jumpy

type Parser struct {
	pageCh chan *Page
	linkCh chan string
	Bucket *Bucket
}

func NewParser(bucket *Bucket, pageCh chan *Page, linkCh chan string) *Parser {
	parser := &Parser{pageCh, linkCh, bucket}

	go func() {
		for {
			select {
			case page := <-pageCh:
				// TODO: need uniquify?
				for _, link := range page.Links() {
					parser.linkCh <- link
				}
			}
		}
	}()

	return parser
}
