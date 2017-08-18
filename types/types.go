package types

import (
    "fmt"
    "sync"
)

type URLSet struct {
    set map[string]int
    m sync.RWMutex
}

func (set *URLSet) AddURL(url string) bool {
    set.m.Lock()
    defer set.m.Unlock()

    count := set.set[url]
    set.set[url] = count + 1
    return count > 0
}

func (set *URLSet) AddURLs(urls []string) {
    set.m.Lock()
    for _, url := range urls {
        set.AddURL(url)
    }
    set.m.Unlock()
}

func (set *URLSet) Contains(url string) bool {
    set.m.RLock()
    defer set.m.RUnlock()

    return set.set[url] > 0
}

func (set *URLSet) Length() int {
    set.m.RLock()
    defer set.m.RUnlock()

    return len(set.Slice())
}

func (set *URLSet) Slice() []string {
    set.m.RLock()
    defer set.m.RUnlock()

    var links []string
    for l := range set.set {
        links = append(links, l)
    }
    return links
}

func (set *URLSet) String() string {
    set.m.RLock()
    defer set.m.RUnlock()

    var out string
    for _, link := range set.Slice() {
        out = fmt.Sprintf("%s\n%s :: %d", out, link, set.set[link])
    }
    return out
}

func NewURLSet() URLSet {
    return URLSet{make(map[string]int), sync.RWMutex{}}
}
