package crawl

import (
	"strings"
	"testing"

	"github.com/evilscott/gocrawler/util"
)

func TestGrabLinks(t *testing.T) {
	html := `
    <html>
    <head>
        <link href="/no" />
    </head>
    <body>
        <a href="/yes">yes</a>
        <a href="#">here</a>
        <a href="/">here</a>
        <a>no</a>
        <a foo="bar">foo</a>
        <a href="/nope" rel="nofollow">bar</a>
    </body>
    </html>
    `
	body := strings.NewReader(html)
	expected := []string{"/yes", "#", "/"}
	util.AssertEqualSlice(t, expected, GrabLinks(body), "GrabLinks")
}

func TestGrabLinksMeta(t *testing.T) {
	html := `
    <html>
    <head>
        <meta name="robots" contents="noindex, nofollow" />
    </head>
    <body>
        <a href="/foo">foo</a>
        <a href="/bar">bar</a>
        <a href="/baz">baz</a>
    </body>
    </html>
    `
	body := strings.NewReader(html)
	util.AssertEqualSlice(t, []string{}, GrabLinks(body), "GrabLinksMeta")
}
