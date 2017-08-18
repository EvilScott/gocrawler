package types

import (
    "testing"

    "github.com/evilscott/gocrawler/util"
)

func TestURLSet_AddURL(t *testing.T) {
    urlSet := NewURLSet()
    util.AssertEquals(t, false, urlSet.AddURL("foo"))
    util.AssertEquals(t, true, urlSet.AddURL("foo"))
    util.AssertEquals(t, false, urlSet.AddURL("bar"))
    util.AssertEquals(t, 2, urlSet.set["foo"])
    util.AssertEquals(t, 1, urlSet.set["bar"])
}

func TestURLSet_AddURLs(t *testing.T) {
    urlSet := NewURLSet()
    urlSet.AddURLs([]string{"foo", "foo", "bar"})
    util.AssertEquals(t, 2, urlSet.set["foo"])
    util.AssertEquals(t, 1, urlSet.set["bar"])
}

func TestURLSet_Contains(t *testing.T) {
    urlSet := NewURLSet()
    util.AssertEquals(t, false, urlSet.Contains("foo"))
    urlSet.AddURL("foo")
    util.AssertEquals(t, true, urlSet.Contains("foo"))
}

func TestURLSet_Length(t *testing.T) {
    urlSet := NewURLSet()
    urlSet.AddURL("foo")
    util.AssertEquals(t, 1, urlSet.Length())
    urlSet.AddURL("foo")
    urlSet.AddURL("bar")
    util.AssertEquals(t, 2, urlSet.Length())
}
