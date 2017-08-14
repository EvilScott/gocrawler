package main

import (
	"fmt"
	"net/url"
	"os"
)

func main() {
	// parse initial url
	startURL, err := url.Parse(os.Args[1])
	if err != nil {
		fmt.Print(err.Error())
		os.Exit(1)
	}

	// url set
	foundLinks := NewURLSet()

	// channels
	todoURLs := make(chan string)
	foundURLs := make(chan []string)

	// crawler workers
	// TODO configurable number of workers via args
	for i := 1; i <= 3; i++ {
		go crawler(i, startURL.Scheme, startURL.Hostname(), todoURLs, foundURLs)
	}

	// waiting on workers
	fmt.Println("Starting crawl ...")
	todoURLs <- startURL.Path
	for founds := range foundURLs {
		for _, found := range founds {
			if !foundLinks.AddURL(found) {
				todoURLs <- found
			}
		}
	}
	fmt.Printf("Finished! Links found:%s\n", foundLinks.String())
}
