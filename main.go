package main

import (
	"fmt"
	"golang.org/x/net/html"
	"net/http"
	"net/url"
	"os"
)

func crawl(scheme string, domain string, path string) []string {
	// create url from domain and path
	targetURL := fmt.Sprintf("%s://%s%s", scheme, domain, path)
	baseURL, err := url.Parse(targetURL)
	if err != nil {
		fmt.Print(err.Error())
		return []string{}
	}

	// hit URL
	fmt.Printf("Crawling %s ...", targetURL)
	resp, err := http.Get(targetURL)
	if err != nil {
		fmt.Print(err.Error())
		return []string{}
	}
	defer resp.Body.Close()

	// keep track of urls on page
	linkMap := NewURLSet()

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
					if scheme == resolvedUrl.Scheme && domain == resolvedUrl.Hostname() {
						linkMap.AddURL(resolvedUrl.Path)
					}
					attrWalk = false
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
	startURL, err := url.Parse(os.Args[1])
	if err != nil {
		fmt.Print(err.Error())
		os.Exit(1)
	}
	var visitedLinks []string
	foundLinks := NewURLSet()
	fmt.Print("Starting crawl ...\n")
	foundLinks.AddURLs(crawl(startURL.Scheme, startURL.Hostname(), startURL.Path))
	visitedLinks = append(visitedLinks, startURL.Path)
	fmt.Printf("Finished! Links found:\n\n%s\n", foundLinks.String())
}
