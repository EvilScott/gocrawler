package crawl

import (
    "io/ioutil"
    "net/url"
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
        <a href="/yes">yes</a>
        <a href="/yes/yes">yes2</a>
        <a href="#">here</a>
        <a href="/">here</a>
        <a href="">here</a>
        <a>no</a>
    </body>
    </html>
    `
    base, _ := url.Parse("http://example.com")
    body := ioutil.NopCloser(strings.NewReader(html))
    expected := []string{"/yes", "/yes/yes", "", "/"}
    util.AssertEqualSlice(t, expected, GrabLinks(base, body))
}
