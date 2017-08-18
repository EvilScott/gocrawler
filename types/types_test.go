package types

import (
    "testing"

    "github.com/evilscott/gocrawler/util"
)

func TestURLSet_AddURL(t *testing.T) {
    urlSet := NewURLSet()
    util.Assert(t, false, urlSet.AddURL("foo"))
    util.Assert(t, true, urlSet.AddURL("foo"))
    util.Assert(t, false, urlSet.AddURL("bar"))
    util.Assert(t, 2, urlSet.set["foo"])
    util.Assert(t, 1, urlSet.set["bar"])
}

func TestURLSet_AddURLs(t *testing.T) {
    urlSet := NewURLSet()
    urlSet.AddURLs([]string{"foo", "foo", "bar"})
    util.Assert(t, 2, urlSet.set["foo"])
    util.Assert(t, 1, urlSet.set["bar"])
}

func TestURLSet_Contains(t *testing.T) {
    urlSet := NewURLSet()
    util.Assert(t, false, urlSet.Contains("foo"))
    urlSet.AddURL("foo")
    util.Assert(t, true, urlSet.Contains("foo"))
}

func TestURLSet_Length(t *testing.T) {
    urlSet := NewURLSet()
    urlSet.AddURL("foo")
    util.Assert(t, 1, urlSet.Length())
    urlSet.AddURL("foo")
    urlSet.AddURL("bar")
    util.Assert(t, 2, urlSet.Length())
}
