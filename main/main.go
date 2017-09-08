package main

import (
    "flag"
    "fmt"
    "net/http"
    "net/url"
    "os"
    "sync"
    "time"

    "github.com/evilscott/gocrawler/crawl"
    "github.com/evilscott/gocrawler/robots"
    "github.com/evilscott/gocrawler/types"
)

func main() {
    // Take Worker count from args.
    var workers int
    flag.IntVar(&workers, "w", 3, "number of workers")

    // Take User-Agent from args.
    var userAgent string
    flag.StringVar(&userAgent, "ua", "gocrawler/1.0", "user agent string")

    // Take CrawlDelay from args (overridden by higher value in robots.txt).
    var crawlDelay int
    flag.IntVar(&crawlDelay, "cd", 0, "crawl delay")

    // Take RedirectCount from args.
    var redirectCount int
    flag.IntVar(&redirectCount, "r", 10, "redirect count")

    // Take QuietMode from args.
    var quietMode bool
    flag.BoolVar(&quietMode, "q", false, "quiet mode")

    // Take VerboseMode from args.
    var verboseMode bool
    flag.BoolVar(&verboseMode, "v", false, "verbose mode")

    // Verbose mode overrides Quiet mode.
    if verboseMode && quietMode {
        quietMode = false
    }

    // Parse arguments.
    flag.Parse()

    // Parse base URL.
    base, err := url.Parse(flag.Arg(0))
    if err != nil {
        fmt.Fprint(os.Stderr, err.Error())
        os.Exit(1)
    }

    // Create exclusions from robots.txt file.
    ex := robots.Exclusion{}
    res, err := http.Get(fmt.Sprintf("%s://%s/robots.txt", base.Scheme, base.Host))
    if err == nil {
        ex = robots.Parse(userAgent, res.Body)
    }
    if crawlDelay > ex.CrawlDelay {
        ex.CrawlDelay = crawlDelay
    }
    ex.VerboseMode = verboseMode

    // Create common worker config.
    c := crawl.Config{
        Exclusions: ex,
        QuietMode: quietMode,
        RedirectCount: redirectCount,
        UserAgent: userAgent,
        VerboseMode: verboseMode,
    }

    // Keep track of results.
    results := types.NewResultSet(*base, ex)

    // Create channels.
    todos := make(chan string, 1000)
    found := make(chan []string, workers)
    badURLs := make(chan [2]string, workers)

    // Keep track of Worker status.
    wg := sync.WaitGroup{}

    // Create crawl Workers.
    for i := 1; i <= workers; i++ {
        go crawl.Worker(i, c, todos, found, badURLs, &wg)
    }

    // Start crawl with base URL.
    if c.QuietMode == false {
        fmt.Printf("Starting crawl with %d workers ...\n", workers)
    }
    found <- []string{base.Path}

    // Routine to process found URLs.
    go func() {
        var todoCount int
        for links := range found {
            reasons := make(map[string]int)
            wg.Add(1)
            todoCount = 0
            for _, link := range links {
                shouldCrawl, reason, crawlURL := results.Add(link)
                if shouldCrawl {
                    todoCount++
                    todos <- crawlURL
                } else {
                    reasons[reason]++
                }
            }
            if c.VerboseMode == true {
                fmt.Printf("%d links processed; %d new links to crawl\n", len(links), todoCount)
                for r, c := range reasons {
                    fmt.Printf("%s %d\n", r, c)
                }
            }
            wg.Done()
        }
    }()

    // Routine to process error URLs.
    go func() {
        for bad := range badURLs {
            wg.Add(1)
            results.AddError(bad[0], bad[1])
            wg.Done()
        }
    }()

    // Wait for all workers to finish.
    time.Sleep(time.Second * 5)
    wg.Wait()
    if c.QuietMode == false {
        fmt.Println("Finished! Links found:")
    }
    fmt.Println(results.String())
}
