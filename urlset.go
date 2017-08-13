package main

type URLSet struct {
	set map[string]bool
}

func (set *URLSet) Add(url string) bool {
	_, found := set.set[url]
	set.set[url] = true
	return !found
}

func (set *URLSet) Slice() []string {
	var links []string
	for l := range set.set {
		links = append(links, l)
	}
	return links
}

func NewURLSet() URLSet {
	return URLSet{set: make(map[string]bool)}
}
