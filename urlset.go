package main

import "fmt"

type URLSet struct {
	set map[string]int
}

func (set *URLSet) AddURL(url string) bool {
	count, found := set.set[url]
	set.set[url] = count + 1
	return !found
}

func (set *URLSet) AddURLs(urls []string) {
	for _, url := range urls {
		set.AddURL(url)
	}
}

func (set *URLSet) Slice() []string {
	var links []string
	for l := range set.set {
		links = append(links, l)
	}
	return links
}

func (set *URLSet) String() string {
	var out string
	for _, link := range set.Slice() {
		out = fmt.Sprintf("%s\n%s :: %d", out, link, set.set[link])
	}
	return out
}

func NewURLSet() URLSet {
	return URLSet{set: make(map[string]int)}
}
