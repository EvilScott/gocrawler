package robots

import (
    "io"
    "io/ioutil"
    "fmt"
    "strings"
)

type Exclusion struct {
    allow []string
    disallow []string
    crawlDelay int
}

// Parse processes a robots.txt input and outputs an Exclusion struct to guide the crawler for a given User-Agent
func Parse(ua string, txt io.Reader) Exclusion {
    ex := Exclusion{}
    body, err := ioutil.ReadAll(txt)
    if err != nil {
        fmt.Println(err.Error())
        return ex
    }
    for _, line := range strings.Split(string(body), "\n") {
        switch {
        case strings.Index(line, "User-agent:") == 0:
            continue // TODO
        case strings.Index(line, "Disallow:") == 0:
            continue // TODO
        case strings.Index(line, "Allow:") == 0:
            continue // TODO
        case strings.Index(line, "Crawl-delay:") == 0:
            continue // TODO
        default:
            continue
        }
    }
    return Exclusion{}
}

// Allowed uses a robots Exclusion struct determine if a url is allowed to be crawled
func (e Exclusion) Allowed(url string) bool {
    for _, allow := range e.allow {
        if strings.Index(url, allow) == 0 {
            return true
        }
    }
    for _, disallow := range e.disallow {
        if strings.Index(url, disallow) == 0 {
            return false
        }
    }
    return true
}
