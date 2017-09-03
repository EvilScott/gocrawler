package crawl

import (
    "errors"
    "fmt"
    "io"
    "net/http"
    "sync"
    "time"

    "github.com/evilscott/gocrawler/robots"

    "golang.org/x/net/html"
)

// Config keeps track of pertinent settings for the crawler.
type Config struct {
    Exclusions robots.Exclusion
    RedirectCount int
    UserAgent string
}

// GrabLinks returns a slice of found links.
func GrabLinks(body io.Reader) []string {
    // Keep track of found links.
    var links []string

    // Tokenize the body for walking.
    z := html.NewTokenizer(body)

    // Walk page elements looking for links.
    for tt := z.Next(); tt != html.ErrorToken; tt = z.Next() {
        // Skip any tags that are not a starting tag.
        if tt != html.StartTagToken {
            continue
        }

        // Walk element attributes looking for href.
        tag, attr := z.TagName()
        if attr == true && string(tag) == "a" {
            walk := true
            for walk {
                name, value, more := z.TagAttr()
                if string(name) == "href" {
                    links = append(links, string(value))
                    walk = false
                } else if more == false {
                    walk = false
                }
            }
        }
    }

    // Return slice of found links.
    return links
}

// Worker grabs URLs from a given channel and crawls them for links.
func Worker(id int, c Config, todos <-chan string, found chan<- []string, wg *sync.WaitGroup) {
    // Create reuseable HTTP client.
    client := &http.Client{}

    // Handle redirects.
    client.CheckRedirect = func(req *http.Request, via []*http.Request) error {
        if len(via) >= c.RedirectCount {
            return errors.New(fmt.Sprintf("stopped after %d redirects", c.RedirectCount))
        }
        return nil
    }

    // Listen to todos channel.
    for target := range todos {
        // Notify WaitGroup that Worker is busy.
        wg.Add(1)

        // Create request for target URL.
        req, err := http.NewRequest("GET", target, nil)
        if err != nil {
            fmt.Println(err.Error())
            wg.Done()
            continue
        }

        // Set User-Agent for request.
        req.Header.Set("User-Agent", c.UserAgent)

        // Send the request.
        fmt.Printf("Crawler #%d %s\n", id, target)
        resp, err := client.Do(req)
        if err != nil {
            fmt.Println(err.Error())
            wg.Done()
            continue
        }

        // Grab links and send them to the found channel for processing.
        found <- GrabLinks(resp.Body)

        // Close the response body and notify the WaitGroup that the Worker is not busy.
        resp.Body.Close()
        wg.Done()

        // Throttle requests if specified by config.
        time.Sleep(c.Exclusions.CrawlDelay * time.Second)
    }
}
