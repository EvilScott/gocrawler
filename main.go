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

	// variables
	var link string
	var newLinks []string
	todoLinks := []string{startURL.Path}
	visitedLinks := NewURLSet()
	foundLinks := NewURLSet()

	// main loop
	fmt.Println("Starting crawl ...")
	for len(todoLinks) > 0 {
		// pop a link from the front of todos
		link, todoLinks = todoLinks[0], todoLinks[1:]

		// mark it as visited and crawl it
		// TODO do this in a goroutine
		newLinks = crawl(startURL.Scheme, startURL.Hostname(), link)
		// TODO check here for successful crawl
		visitedLinks.AddURL(link)

		// record new links
		foundLinks.AddURLs(newLinks)

		// record unvisited links
		for _, newLink := range newLinks {
			if !visitedLinks.Contains(newLink) && !contains(newLink, todoLinks) {
				todoLinks = append(todoLinks, newLink)
			}
		}
	}
	fmt.Printf("Finished! Links found:%s\n", foundLinks.String())
}
