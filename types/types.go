package types

import (
	"fmt"
	"net/url"

	"github.com/evilscott/gocrawler/robots"
)

// Constant reasons for not crawling.
const (
	CANNOT_PARSE = "cannot parse"
	DISALLOWED   = "disallowed"
	DUPLICATE    = "duplicate"
	EXTERNAL     = "external"
)

// Config keeps track of pertinent settings for the crawler.
type Config struct {
	BufferSize    int
	Exclusions    robots.Exclusion
	QuietMode     bool
	RedirectCount int
	UserAgent     string
	VerboseMode   bool
}

// URLData keeps track of extra data associated with URLs.
type ResponseData struct {
	URL    string
	Status string
	Code   int
	Time   int
}

// ChannelGroup keeps track of all the channels a crawl routine needs.
type ChannelGroup struct {
	TODOs     chan string
	Found     chan []string
	Responses chan ResponseData
}

// Result handles info for a given link found during the crawl.
type Result struct {
	external       bool
	redirect       bool
	redirectedFrom url.URL
	url            url.URL
}

// ResultSet keeps track of found links during the crawl and associated data.
type ResultSet struct {
	base        url.URL
	set         map[string][]Result
	responseSet map[string]ResponseData
	errorSet    map[string]string
	ex          robots.Exclusion
}

// NewResultSet creates a new ResultSet struct.
func NewResultSet(base url.URL, ex robots.Exclusion) ResultSet {
	return ResultSet{
		base:        base,
		set:         make(map[string][]Result),
		responseSet: make(map[string]ResponseData),
		errorSet:    make(map[string]string),
		ex:          ex,
	}
}

// Add adds a new link to the ResultSet and returns if it should be crawled, the reason for exclusion (if applicable),
// and the full URL.
func (rs ResultSet) Add(link string) (shouldCrawl bool, reason, crawlURL string) {
	// Parse the new link and resolve it against the base URL.
	parsed, err := url.Parse(string(link))
	if err != nil {
		return false, CANNOT_PARSE, link
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
	switch {
	case !rs.ex.Allowed(resolved.Path):
		shouldCrawl = false
		reason = DISALLOWED
	case len(rs.set[key]) > 0:
		shouldCrawl = false
		reason = DUPLICATE
	case external:
		shouldCrawl = false
		reason = EXTERNAL
	default:
		shouldCrawl = true
	}

	// Add the link to the ResultSet.
	rs.set[key] = append(rs.set[key], Result{
		external:       external,
		redirect:       false,
		redirectedFrom: url.URL{},
		url:            *resolved,
	})

	// Return whether whether it should be crawled and the link itself.
	return shouldCrawl, reason, fmt.Sprintf("%s://%s%s", resolved.Scheme, resolved.Host, resolved.Path)
}

// AddResponse adds the response (and error if appropriate) to the ResultSet.
func (rs ResultSet) AddResponse(resp ResponseData) {
	rs.responseSet[resp.URL] = resp
	if resp.Code >= 400 {
		rs.errorSet[resp.URL] = resp.Status
	}
}

// String returns a string representation of the ResultSet.
func (rs ResultSet) String() string {
	// Iterate over ResultSet and build string representation.
	var out string
	for link, resp := range rs.responseSet {
		out += fmt.Sprintf("%s %d (%dms)\n", link, resp.Code, resp.Time)
	}
	for link, results := range rs.set {
		if results[0].external {
			out += fmt.Sprintf("%s\n", link)
		}
	}
	for link := range rs.errorSet {
		out += fmt.Sprintf("%s %s\n", link, rs.errorSet[link])
	}
	return out
}
