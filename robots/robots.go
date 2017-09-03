package robots

import (
    "io"
    "io/ioutil"
    "fmt"
    "regexp"
    "strconv"
    "strings"
)

type Exclusion struct {
    Allow []string
    Disallow []string
    CrawlDelay int
}

// parseSection parses a single section of robots.txt for use in Parse.
func parseSection(section string) Exclusion {
    var err error

    allow := []string{}
    allowMatch := regexp.MustCompile("Allow:\\s?(.+)\\s?").FindAllStringSubmatch(section, -1)
    for _, match := range allowMatch {
        allow = append(allow, match[1])
    }

    disallow := []string{}
    disallowMatch := regexp.MustCompile("Disallow:\\s?(.+)\\s?").FindAllStringSubmatch(section, -1)
    for _, match := range disallowMatch {
        disallow = append(disallow, match[1])
    }

    var crawlDelay int64
    crawlDelayMatch := regexp.MustCompile("Crawl-delay:\\s?(\\d+)\\s?").FindStringSubmatch(section)
    if len(crawlDelayMatch) == 2 {
        crawlDelay, err = strconv.ParseInt(crawlDelayMatch[1], 10, 0)
        if err != nil {
            crawlDelay = 0
        }
    } else {
        crawlDelay = 0
    }

    return Exclusion{
        Allow: allow,
        Disallow: disallow,
        CrawlDelay: int(crawlDelay),
    }
}

// Parse processes a robots.txt input and outputs an Exclusion struct to guide the crawler for a given User-Agent
func Parse(ua string, txt io.Reader) Exclusion {
    b, err := ioutil.ReadAll(txt)
    if err != nil {
        fmt.Println(err.Error())
        return Exclusion{}
    }
    body := string(b)

    // Remove comments.
    re := regexp.MustCompile("#[^\n]+")
    body = re.ReplaceAllLiteralString(body, "")

    // Break apart User-agent sections.
    re = regexp.MustCompile("((User-agent:[^\n]+\n)+)")
    body = re.ReplaceAllString(body, "!!!\n$1")
    sections := strings.Split(body, "!!!\n")
    // Check for at least one User-agent section.
    if len(sections) == 0 {
        return Exclusion{}
    } else {
        sections = sections[1:]
    }

    // Case insensitive check for passed User-agent ua and fallback to *.
    uaRegex, err := regexp.Compile(fmt.Sprintf("(?i)user-agent:\\s*%s\\s*", ua))
    if err != nil {
        fmt.Println(err.Error())
        return Exclusion{}
    }
    wildRegex := regexp.MustCompile("(?i)user-agent:\\s*\\*\\s*")
    var ex Exclusion
    for _, section := range sections {
        switch {
        case uaRegex.FindString(section) != "":
            ex = parseSection(section)
        case wildRegex.FindString(section) != "":
            if ex.Blank() {
                ex = parseSection(section)
            }
        }
    }

    // Return our best matching Exclusion (or a blank one).
    return ex
}

// Blank returns true if an Exclusion has no data in it.
func (e Exclusion) Blank() bool {
    return len(e.Allow) == 0 && len(e.Disallow) == 0 && e.CrawlDelay == 0
}

// Allowed uses a robots Exclusion struct determine if a url is allowed to be crawled
func (e Exclusion) Allowed(url string) bool {
    for _, allow := range e.Allow {
        if strings.Index(url, allow) == 0 || allow == "*"{
            return true
        }
    }
    for _, disallow := range e.Disallow {
        if strings.Index(url, disallow) == 0 || disallow == "*" {
            return false
        }
    }
    return true
}
