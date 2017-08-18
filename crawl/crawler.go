package crawl

import (
	"fmt"
	"io"
	"net/http"
	"net/url"

	"golang.org/x/net/html"

	"github.com/evilscott/gocrawler/types"
)

func grabLinks(baseURL *url.URL, body io.ReadCloser) []string {
	// keep track of urls on page
	linkMap := types.NewURLSet()

	// walk elements
	z := html.NewTokenizer(body)
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

func Crawler(id int, scheme, domain string, todos <-chan string, found chan<- string) {
	for {
		select {
		case path := <- todos:
			// create url from domain and path
			target := fmt.Sprintf("%s://%s%s", scheme, domain, path)
			base, err := url.Parse(target)
			if err != nil {
				fmt.Println(err.Error())
				continue
			}

			// hit URL
			fmt.Printf("Crawler #%d crawling %s ...\n", id, target)
			resp, err := http.Get(target)
			if err != nil {
				fmt.Println(err.Error())
				continue
			}

			// grab links and send them to channel
			for _, link := range grabLinks(base, resp.Body) {
				found <- link
			}

			// close response body and signal work is done
			resp.Body.Close()

		default: // noop
		}
	}
}
