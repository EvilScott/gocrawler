package robots

import (
    "strings"
    "testing"

    "github.com/evilscott/gocrawler/util"
)

func TestParse(t *testing.T) {
    var txt string = `
    # ignore comment lines
    User-agent: foocrawler # ignore inline comments
    Disallow: /foo
    Disallow: /bar
    Allow: /foo/bar
    Crawl-delay: 10

    User-agent: barcrawler
    User-agent: bazcrawler
    Disallow: /
    Crawl-delay: invalid

    User-agent: *
    Allow: /
    Disallow: /admin
    `

    fooEx := Parse("foocrawler", strings.NewReader(txt))
    util.AssertEqualSlice(t, []string{"/foo/bar"}, fooEx.allow)
    util.AssertEqualSlice(t, []string{"/foo", "/bar"}, fooEx.disallow)
    util.AssertEquals(t, 10, fooEx.crawlDelay)

    barEx := Parse("barcrawler", strings.NewReader(txt))
    util.AssertEqualSlice(t, []string{}, fooEx.allow)
    util.AssertEqualSlice(t, []string{"/"}, fooEx.disallow)
    util.AssertEquals(t, 0, fooEx.crawlDelay)

    bazEx := Parse("bazcrawler", strings.NewReader(txt))
    util.AssertEqualSlice(t, barEx.allow, bazEx.allow)
    util.AssertEqualSlice(t, barEx.disallow, bazEx.disallow)
    util.AssertEquals(t, barEx.crawlDelay, bazEx.crawlDelay)

    otherEx := Parse("othercrawler", strings.NewReader(txt))
    util.AssertEqualSlice(t, []string{"/"}, otherEx.allow)
    util.AssertEqualSlice(t, []string{}, otherEx.disallow)
    util.AssertEquals(t, 0, otherEx.crawlDelay)
}

func TestExclusion_Allowed(t *testing.T) {
    ex := Exclusion{
        allow: []string{"/foo/bar"},
        disallow: []string{"/foo", "/bar"},
        crawlDelay: 10,
    }
    util.AssertEquals(t, true, ex.Allowed("/foo/bar"))
    util.AssertEquals(t, false, ex.Allowed("/foo"))
    util.AssertEquals(t, false, ex.Allowed("/bar"))
    util.AssertEquals(t, true, ex.Allowed("/"))
    util.AssertEquals(t, true, ex.Allowed("/foo/bar/baz"))
    util.AssertEquals(t, true, ex.Allowed("/baz"))
    util.AssertEquals(t, false, ex.Allowed("/foo/baz"))
}
