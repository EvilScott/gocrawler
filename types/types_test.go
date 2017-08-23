package types

import (
    "net/url"
    "testing"

    "github.com/evilscott/gocrawler/util"
)

func TestResultSet_Add(t *testing.T) {
    var shouldCrawl bool
    var crawlURL string

    base, _ := url.Parse("http://www.test.com/")
    rs := NewResultSet(*base)

    shouldCrawl, crawlURL = rs.Add("/foo")
    util.AssertEquals(t, true, shouldCrawl)
    util.AssertEquals(t, "http://www.test.com/foo", crawlURL)

    shouldCrawl, crawlURL = rs.Add("/foo")
    util.AssertEquals(t, false, shouldCrawl)
    util.AssertEquals(t, "http://www.test.com/foo", crawlURL)

    shouldCrawl, crawlURL = rs.Add("/bar")
    util.AssertEquals(t, true, shouldCrawl)
    util.AssertEquals(t, "http://www.test.com/bar", crawlURL)

    shouldCrawl, crawlURL = rs.Add("https://otherdomain.com/baz")
    util.AssertEquals(t, false, shouldCrawl)
    util.AssertEquals(t, "https://otherdomain.com/baz", crawlURL)

    shouldCrawl, crawlURL = rs.Add("https://www.test.com/bbq")
    util.AssertEquals(t, false, shouldCrawl)
    util.AssertEquals(t, "https://www.test.com/bbq", crawlURL)

    util.AssertEquals(t, 2, len(rs.set["/foo"]))
    util.AssertEquals(t, 1, len(rs.set["/bar"]))
}
