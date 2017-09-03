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
        <a href="">here</a>
        <a>no</a>
    </body>
    </html>
    `
    body := strings.NewReader(html)
    expected := []string{"/yes", "#", "/", ""}
    util.AssertEqualSlice(t, expected, GrabLinks(body), "GrabLinks")
}
