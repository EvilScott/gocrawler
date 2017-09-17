package crawl

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/evilscott/gocrawler/types"

	"golang.org/x/net/html"
)

// GrabLinks returns a slice of found links.
func GrabLinks(body io.Reader) []string {
	// Keep track of found links.
	var links []string

	// Tokenize the body for walking.
	z := html.NewTokenizer(body)

	// Walk page elements looking for links.
	for tt := z.Next(); tt != html.ErrorToken; tt = z.Next() {
		// Skip any tags that are not a starting tag.
		if tt != html.StartTagToken && tt != html.SelfClosingTagToken {
			continue
		}

		// Walk element attributes looking for href.
		tag, attr := z.TagName()
		switch {
		// Skip crawling pages with meta robots.
		case attr && string(tag) == "meta":
			var metaName, metaContent string
			walk := true
			for walk {
				name, value, more := z.TagAttr()
				if string(name) == "name" {
					metaName = string(value)
				} else if string(name) == "contents" {
					metaContent = string(value)
				}
				if more == false {
					walk = false
				}
			}
			if strings.ToLower(metaName) == "robots" {
				re := regexp.MustCompile(",\\s?")
				for _, content := range re.Split(metaContent, -1) {
					if content == "nofollow" {
						return []string{}
					}
				}
			}
		// Record all links.
		case attr && string(tag) == "a":
			walk := true
			for walk {
				var link string
				name, value, more := z.TagAttr()
				if string(name) == "href" {
					link = string(value)
				} else if string(name) == "rel" && string(value) == "nofollow" {
					walk = false
				}
				if more == false {
					if link != "" {
						links = append(links, link)
					}
					walk = false
				}
			}
		}
	}

	// Return slice of found links.
	return links
}

// Worker grabs URLs from a given channel and crawls them for links.
func Worker(id int, c types.Config, chans types.ChannelGroup, wg *sync.WaitGroup) {
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
	for target := range chans.TODOs {
		// Create request for target URL.
		req, err := http.NewRequest("GET", target, nil)
		if err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
			time.Sleep(time.Second * time.Duration(c.Exclusions.CrawlDelay))
			wg.Done()
			continue
		}

		// Set User-Agent for request.
		req.Header.Set("User-Agent", c.UserAgent)

		// Send the request.
		if c.QuietMode == false {
			fmt.Printf("Crawler #%d %s :: crawling ...\n", id, target)
		}
		start := time.Now()
		resp, err := client.Do(req)
		dur := int(time.Now().Sub(start) / time.Millisecond)
		if err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
			time.Sleep(time.Second * time.Duration(c.Exclusions.CrawlDelay))
			wg.Done()
			continue
		}
		if c.VerboseMode {
			fmt.Printf("Crawler #%d %s :: %s (%dms)\n", id, target, resp.Status, dur)
		}

		// Log non 2xx/3xx responses.
		if resp.StatusCode >= 400 {
			fmt.Fprintf(os.Stderr, "%s :: %s", target, resp.Status)
		}

		// Record extra data for responses.
		wg.Add(1)
		chans.Responses <- types.ResponseData{
			URL:    target,
			Status: resp.Status,
			Code:   resp.StatusCode,
			Time:   dur,
		}

		// Grab links and send them to the found channel for processing.
		links := GrabLinks(resp.Body)
		if c.VerboseMode {
			fmt.Printf("Crawler #%d %s :: %d links found\n", id, target, len(links))
		}
		wg.Add(1)
		chans.Found <- links

		// Throttle requests if specified by config.
		time.Sleep(time.Second * time.Duration(c.Exclusions.CrawlDelay))

		// Close the response body and notify the WaitGroup that the Worker is not busy.
		resp.Body.Close()
		wg.Done()
	}
}
