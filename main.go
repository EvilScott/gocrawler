package main

import (
	"net/http"
	"golang.org/x/net/html"
	"fmt"
)

const START_URL = "http://www.stedwards.edu/"

func crawl(url string) []string {
	var links []string
	resp, _ := http.Get(url)
	defer resp.Body.Close()
	z := html.NewTokenizer(resp.Body)
	tokenWalk := true
	for tokenWalk {
		tt := z.Next()
		switch tt {
		case html.StartTagToken:
			tagName, hasAttr := z.TagName()
			if hasAttr == true && string(tagName) == "a" {
				attrWalk := true
				for attrWalk {
					name, value, more := z.TagAttr()
					if string(name) == "href" {
						fmt.Printf("Link found: %s\n", value)
						links = append(links, string(value))
					}
					if more == false {
						attrWalk = false
					}
				}
			}
		case html.ErrorToken:
			tokenWalk = false
		}
	}
	return links
}

func main() {
	fmt.Println("Starting crawl...")
	links := crawl(START_URL)
	fmt.Printf("Links: %s", links)
	fmt.Println("Finished!")
}
