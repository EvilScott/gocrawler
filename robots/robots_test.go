package robots

import (
    "strings"
    "testing"

    "github.com/evilscott/gocrawler/util"
)

func TestParse(t *testing.T) {
    var txt string = `# ignore comment lines
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
Allow: *
Disallow: /admin
`

    fooEx := Parse("foocrawler", strings.NewReader(txt))
    util.AssertEqualSlice(t, []string{"/foo/bar"}, fooEx.allow, "Parse fooEx.allow")
    util.AssertEqualSlice(t, []string{"/foo", "/bar"}, fooEx.disallow, "Parse fooEx.disallow")
    util.AssertEquals(t, 10, fooEx.crawlDelay, "Parse fooEx.crawlDelay")

    barEx := Parse("barcrawler", strings.NewReader(txt))
    util.AssertEqualSlice(t, []string{},  barEx.allow, "Parse barEx.allow")
    util.AssertEqualSlice(t, []string{"/"}, barEx.disallow, "Parse barEx.disallow")
    util.AssertEquals(t, 0, barEx.crawlDelay, "Parse barEx.crawlDelay")

    bazEx := Parse("BaZcRaWlEr", strings.NewReader(txt))
    util.AssertEqualSlice(t, barEx.allow, bazEx.allow, "Parse bazEx.allow")
    util.AssertEqualSlice(t, barEx.disallow, bazEx.disallow, "Parse bazEx.disallow")
    util.AssertEquals(t, barEx.crawlDelay, bazEx.crawlDelay, "Parse bazEx.crawlDelay")

    otherEx := Parse("othercrawler", strings.NewReader(txt))
    util.AssertEqualSlice(t, []string{"*"}, otherEx.allow, "Parse otherEx.allow")
    util.AssertEqualSlice(t, []string{"/admin"}, otherEx.disallow, "Parse otherEx.disallow")
    util.AssertEquals(t, 0, otherEx.crawlDelay, "Parse otherEx.crawlDelay")
}

func TestExclusion_Blank(t *testing.T) {
    blankEx := Exclusion{
        allow: []string{},
        disallow: []string{},
        crawlDelay: 0,
    }
    util.AssertEquals(t, true, blankEx.Blank(), "Blank")

    notBlankEx := Exclusion{
        allow: []string{"foo"},
        disallow: []string{"bar"},
        crawlDelay: 3,
    }
    util.AssertEquals(t, false, notBlankEx.Blank(), "Blank")

}

func TestExclusion_Allowed(t *testing.T) {
    ex := Exclusion{
        allow: []string{"/foo/bar"},
        disallow: []string{"/foo", "/bar"},
        crawlDelay: 10,
    }
    util.AssertEquals(t, true, ex.Allowed("/foo/bar"), "Allowed")
    util.AssertEquals(t, false, ex.Allowed("/foo"), "Allowed")
    util.AssertEquals(t, false, ex.Allowed("/bar"), "Allowed")
    util.AssertEquals(t, true, ex.Allowed("/"), "Allowed")
    util.AssertEquals(t, true, ex.Allowed("/foo/bar/baz"), "Allowed")
    util.AssertEquals(t, true, ex.Allowed("/baz"), "Allowed")
    util.AssertEquals(t, false, ex.Allowed("/foo/baz"), "Allowed")
}

func TestExclusion_AllowedWildcards(t *testing.T) {
    ex := Exclusion{
        allow: []string{"/allowed"},
        disallow: []string{"*"},
    }
    util.AssertEquals(t, true, ex.Allowed("/allowed"), "AllowedWildcard")
    util.AssertEquals(t, false, ex.Allowed("/not/allowed"), "AllowedWildcard")

    ex = Exclusion{
        allow: []string{"*"},
        disallow: []string{"/disallowed"},
    }
    util.AssertEquals(t, true, ex.Allowed("/disallowed"), "AllowedWildcard")
    util.AssertEquals(t, true, ex.Allowed("/allowed"), "AllowedWildcard")
}
