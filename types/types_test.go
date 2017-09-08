package types

import (
    "net/url"
    "testing"

    "github.com/evilscott/gocrawler/robots"
    "github.com/evilscott/gocrawler/util"
)

func TestResultSet_Add(t *testing.T) {
    var shouldCrawl bool
    var reason, crawlURL string
    ex := robots.Exclusion{}

    base, _ := url.Parse("http://www.test.com/")
    rs := NewResultSet(*base, ex)

    shouldCrawl, reason, crawlURL = rs.Add("/foo")
    util.AssertEquals(t, true, shouldCrawl, "Add /foo - shouldCrawl")
    util.AssertEquals(t, "", reason, "Add /foo - reason")
    util.AssertEquals(t, "http://www.test.com/foo", crawlURL, "Add /foo - crawlURL")

    shouldCrawl, reason, crawlURL = rs.Add("/foo")
    util.AssertEquals(t, false, shouldCrawl, "Add /foo (duplicate) - shouldCrawl")
    util.AssertEquals(t, DUPLICATE, reason, "Add /foo (duplicate - reason)")
    util.AssertEquals(t, "http://www.test.com/foo", crawlURL, "Add /foo (duplicate) - crawlURL")

    shouldCrawl, reason, crawlURL = rs.Add("/bar")
    util.AssertEquals(t, true, shouldCrawl, "Add /bar - shouldCrawl")
    util.AssertEquals(t, "", reason, "Add /bar - reason")
    util.AssertEquals(t, "http://www.test.com/bar", crawlURL, "Add /bar - crawlURL")

    shouldCrawl, reason, crawlURL = rs.Add("https://otherdomain.com/baz")
    util.AssertEquals(t, false, shouldCrawl, "Add /baz (eternal) - shouldCrawl")
    util.AssertEquals(t, EXTERNAL, reason, "Add /baz (eternal) - reason")
    util.AssertEquals(t, "https://otherdomain.com/baz", crawlURL, "Add /baz (eternal) - crawlURL")

    shouldCrawl, reason, crawlURL = rs.Add("https://www.test.com/bbq")
    util.AssertEquals(t, false, shouldCrawl, "Add /bbq (external) - shouldCrawl")
    util.AssertEquals(t, EXTERNAL, reason, "Add /bbq (external) - reason")
    util.AssertEquals(t, "https://www.test.com/bbq", crawlURL, "Add /bbq (external) - crawlURL")

    util.AssertEquals(t, 2, len(rs.set["/foo"]), "Add count /foo")
    util.AssertEquals(t, 1, len(rs.set["/bar"]), "Add count /bar")
}
