package crawl

import (
    "fmt"
    "io"
    "net/http"
    "net/url"
    "sync"

    "golang.org/x/net/html"

    "github.com/evilscott/gocrawler/types"
)

func GrabLinks(base *url.URL, body io.ReadCloser) []string {
    // keep track of urls on page
    links := types.NewURLSet()

    // walk elements
    z := html.NewTokenizer(body)
    for tt := z.Next(); tt != html.ErrorToken; tt = z.Next() {
        // skip anything that isn't a starting tag
        if tt != html.StartTagToken {
            continue
        }

        // walk element attrs
        tag, attr := z.TagName()
        if attr == true && string(tag) == "a" {
            walk := true
            for walk {
                name, value, more := z.TagAttr()
                if string(name) == "href" {
                    current, _ := url.Parse(string(value))
                    resolved := base.ResolveReference(current)
                    if base.Scheme == resolved.Scheme && base.Hostname() == resolved.Hostname() {
                        links.AddURL(resolved.Path)
                    }
                    walk = false
                }
                if more == false {
                    walk = false
                }
            }
        }
    }
    return links.Slice()
}

func Worker(id int, scheme, domain string, todos <-chan string, found chan<- []string, wg *sync.WaitGroup) {
    for path := range todos {
        // mark working
        wg.Add(1)

        // create url from domain and path
        target := fmt.Sprintf("%s://%s%s", scheme, domain, path)
        base, err := url.Parse(target)
        if err != nil {
            fmt.Println(err.Error())
            continue
        }

        // hit URL
        fmt.Printf("Crawler #%d %s\n", id, target)
        resp, err := http.Get(target)
        if err != nil {
            fmt.Println(err.Error())
            continue
        }

        // grab links and send them to channel
        found <- GrabLinks(base, resp.Body)

        // close response body and signal work is done
        resp.Body.Close()
        wg.Done()
    }
}
