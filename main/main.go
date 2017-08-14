package main

import (
	"fmt"
	"net/url"
	"os"
	"sync"

	"github.com/evilscott/gocrawler/crawl"
	"github.com/evilscott/gocrawler/types"
)

func main() {
	// parse initial url
	startURL, err := url.Parse(os.Args[1])
	if err != nil {
		fmt.Print(err.Error())
		os.Exit(1)
	}

	// url set
	foundLinks := types.NewURLSet()

	// channels
	todoURLs := make(chan string)
	foundURLs := make(chan string)

	// wait group
	var wg sync.WaitGroup

	// crawler workers
	// TODO configurable number of workers via args
	for i := 1; i <= 3; i++ {
		go crawl.Crawler(i, startURL.Scheme, startURL.Hostname(), todoURLs, foundURLs, wg)
	}

	// listening for crawler results
	fmt.Println("Starting crawl ...")
	todoURLs <- startURL.Path
	for found := range foundURLs {
		wg.Add(1)
		if !foundLinks.AddURL(found) {
			todoURLs <- found
		}
		wg.Done()
	}

	// wait for all workers to finish
	wg.Wait()
	fmt.Printf("Finished! Links found:%s\n", foundLinks.String())
}
