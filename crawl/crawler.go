package crawl

import (
	"fmt"
	"net/http"
	"net/url"
	"sync"

	"golang.org/x/net/html"

	"github.com/evilscott/gocrawler/types"
)

func grabLinks(baseURL *url.URL, resp *http.Response) []string {
	// keep track of urls on page
	linkMap := types.NewURLSet()

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

func Crawler(id int, scheme, domain string, todoURLs <-chan string, foundURLs chan<- string, wg sync.WaitGroup) {
	for path := range todoURLs {
		// signal that work is happening
		wg.Add(1)

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

		// close response body and signal work is done
		resp.Body.Close()
		wg.Done()
	}
}
