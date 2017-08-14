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
	foundURLs := make(chan string)

	// crawler workers
	// TODO configurable number of workers via args
	for i := 1; i <= 3; i++ {
		go crawler(i, startURL.Scheme, startURL.Hostname(), todoURLs, foundURLs)
	}

	// listening for crawler results
	fmt.Println("Starting crawl ...")
	todoURLs <- startURL.Path
	for found := range foundURLs {
		if !foundLinks.AddURL(found) {
			todoURLs <- found
		}
	}
	fmt.Printf("Finished! Links found:%s\n", foundLinks.String())
}
