package main

import (
	"flag"
	"fmt"

	"github.com/yoppi/jumpy"
)

var (
	command  string
	rootUrl  string
	commands map[string]func(*jumpy.Page)
)

func init() {
	flag.StringVar(&command, "command", "url", "Affect crawled page.")
	flag.StringVar(&rootUrl, "root", "", "Crawling target root url.")
	commands = map[string]func(*jumpy.Page){
		"url": func(page *jumpy.Page) {
			fmt.Println(page.Url)
		},
	}
}

func main() {
	flag.Parse()
	if len(rootUrl) == 0 {
		panic("Must specify root url with --root option.")
	}

	if f, ok := commands[command]; ok {
		jumpy.Crawl(rootUrl, map[string]string{}, f)
	}
}
