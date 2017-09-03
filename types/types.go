package types

import (
    "fmt"
    "net/url"
    "sync"

    "github.com/evilscott/gocrawler/robots"
)

// Result handles info for a given link found during the crawl.
type Result struct {
    external bool
    redirect bool
    redirectedFrom url.URL
    url url.URL
}

// ResultSet keeps track of found links during the crawl and associated data.
type ResultSet struct {
    base url.URL
    set map[string][]Result
    m sync.RWMutex
    ex robots.Exclusion
}

// NewResultSet creates a new ResultSet struct.
func NewResultSet(base url.URL, ex robots.Exclusion) ResultSet {
    return ResultSet{
        base: base,
        set: make(map[string][]Result),
        m: sync.RWMutex{},
        ex: ex,
    }
}

// Add adds a new link to the ResultSet and returns if it should be crawled and the full URL
func (rs ResultSet) Add(link string) (bool, string) {
    // Utilize RWLock for safe map access.
    rs.m.Lock()
    defer rs.m.Unlock()

    // Parse the new link and resolve it against the base URL.
    parsed, err := url.Parse(string(link))
    if err != nil {
        return false, link
    }
    resolved := rs.base.ResolveReference(parsed)

    // Test if the added if link is internal or external.
    external := true
    key := resolved.String()
    if rs.base.Scheme == resolved.Scheme && rs.base.Host == resolved.Host {
        external = false
        key = resolved.Path
    }

    // Test if the added link should be crawled or not.
    found := len(rs.set[key]) == 0 && external == false && rs.ex.Allowed(resolved.Path)

    // Add the link to the ResultSet.
    rs.set[key] = append(rs.set[key], Result{
        external: external,
        redirect: false,
        redirectedFrom: url.URL{},
        url: *resolved,
    })

    // Return whether whether it should be crawled and the link itself.
    return found, fmt.Sprintf("%s://%s%s", resolved.Scheme, resolved.Host, resolved.Path)
}

// String returns a string representation of the ResultSet.
func (rs ResultSet) String() string {
    rs.m.Lock()
    defer rs.m.Unlock()

    // Iterate over ResultSet and build string representation.
    var out string
    for link := range rs.set {
        out = fmt.Sprintf("%s\n%s", out, link)
    }
    return out
}
