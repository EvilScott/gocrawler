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
	util.AssertEqualSlice(t, []string{"/foo/bar"}, fooEx.Allow, "Parse fooEx.Allow")
	util.AssertEqualSlice(t, []string{"/foo", "/bar"}, fooEx.Disallow, "Parse fooEx.Disallow")
	util.AssertEquals(t, 10, fooEx.CrawlDelay, "Parse fooEx.CrawlDelay")

	barEx := Parse("barcrawler", strings.NewReader(txt))
	util.AssertEqualSlice(t, []string{}, barEx.Allow, "Parse barEx.Allow")
	util.AssertEqualSlice(t, []string{"/"}, barEx.Disallow, "Parse barEx.Disallow")
	util.AssertEquals(t, 0, barEx.CrawlDelay, "Parse barEx.CrawlDelay")

	bazEx := Parse("BaZcRaWlEr", strings.NewReader(txt))
	util.AssertEqualSlice(t, barEx.Allow, bazEx.Allow, "Parse bazEx.Allow")
	util.AssertEqualSlice(t, barEx.Disallow, bazEx.Disallow, "Parse bazEx.Disallow")
	util.AssertEquals(t, barEx.CrawlDelay, bazEx.CrawlDelay, "Parse bazEx.CrawlDelay")

	otherEx := Parse("othercrawler", strings.NewReader(txt))
	util.AssertEqualSlice(t, []string{"*"}, otherEx.Allow, "Parse otherEx.Allow")
	util.AssertEqualSlice(t, []string{"/admin"}, otherEx.Disallow, "Parse otherEx.Disallow")
	util.AssertEquals(t, 0, otherEx.CrawlDelay, "Parse otherEx.CrawlDelay")
}

func TestExclusion_Blank(t *testing.T) {
	blankEx := Exclusion{
		Allow:      []string{},
		Disallow:   []string{},
		CrawlDelay: 0,
	}
	util.AssertEquals(t, true, blankEx.Blank(), "Blank")

	notBlankEx := Exclusion{
		Allow:      []string{"foo"},
		Disallow:   []string{"bar"},
		CrawlDelay: 3,
	}
	util.AssertEquals(t, false, notBlankEx.Blank(), "Blank")

}

func TestExclusion_Allowed(t *testing.T) {
	ex := Exclusion{
		Allow:      []string{"/foo/bar"},
		Disallow:   []string{"/foo", "/bar"},
		CrawlDelay: 10,
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
		Allow:    []string{"/allowed"},
		Disallow: []string{"*"},
	}
	util.AssertEquals(t, true, ex.Allowed("/allowed"), "AllowedWildcard")
	util.AssertEquals(t, false, ex.Allowed("/not/allowed"), "AllowedWildcard")

	ex = Exclusion{
		Allow:    []string{"*"},
		Disallow: []string{"/Disallowed"},
	}
	util.AssertEquals(t, true, ex.Allowed("/Disallowed"), "AllowedWildcard")
	util.AssertEquals(t, true, ex.Allowed("/allowed"), "AllowedWildcard")
}
