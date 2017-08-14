package main

import (
	"fmt"
	"golang.org/x/net/html"
	"net/http"
	"net/url"
)

func grabLinks(baseURL *url.URL, resp *http.Response) []string {
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
					if baseURL.Scheme == resolvedUrl.Scheme && baseURL.Hostname() == resolvedUrl.Hostname() {
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

func crawler(id int, scheme, domain string, todoURLs <-chan string, foundURLs chan<- string) {
	for path := range todoURLs {
		// create url from domain and path
		targetURL := fmt.Sprintf("%s://%s%s", scheme, domain, path)
		baseURL, err := url.Parse(targetURL)
		if err != nil {
			fmt.Println(err.Error())
			continue
		}

		// hit URL
		fmt.Printf("Crawler #%d crawling %s ...\n", id, targetURL)
		resp, err := http.Get(targetURL)
		if err != nil {
			fmt.Println(err.Error())
			continue
		}

		// grab links and send them to channel
		for _, found := range grabLinks(baseURL, resp) {
			foundURLs <- found
		}

		// close response body
		resp.Body.Close()
	}
}
