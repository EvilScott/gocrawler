package main

import (
	"fmt"
	"golang.org/x/net/html"
	"net/http"
	"net/url"
)

const START_URL = "http://www.stedwards.edu/"

func crawl(targetURL string) []string {
	baseURL, _ := url.Parse(targetURL)
	linkMap := NewURLSet()

	// hit URL
	resp, _ := http.Get(targetURL)
	defer resp.Body.Close()

	// walk elements
	z := html.NewTokenizer(resp.Body)
	for tt := z.Next(); tt != html.ErrorToken; tt = z.Next() {
		// skip anything that isn't a starting tag
		if tt != html.StartTagToken {
			continue
		}

		// walk element attrs
		tagName, hasAttr := z.TagName()
		if hasAttr == true && string(tagName) == "a" {
			attrWalk := true
			for attrWalk {
				name, value, more := z.TagAttr()
				if string(name) == "href" {
					currentUrl, _ := url.Parse(string(value))
					resolvedUrl := baseURL.ResolveReference(currentUrl)
					if linkMap.Add(resolvedUrl.String()) {
						fmt.Printf("Link found: %s\n", resolvedUrl)
					}
				}
				if more == false {
					attrWalk = false
				}
			}
		}
	}
	return linkMap.Slice()
}

func main() {
	fmt.Println("Starting crawl...")
	links := crawl(START_URL)
	fmt.Printf("Links: %s\n", links)
	fmt.Println("Finished!")
}
