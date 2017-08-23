package main

import (
    "flag"
    "fmt"
    "net/url"
    "os"
    "sync"
    "time"

    "github.com/evilscott/gocrawler/crawl"
    "github.com/evilscott/gocrawler/types"
)

func main() {
    // Take Worker count from args.
    var workers int
    flag.IntVar(&workers, "w", 3, "number of workers")

    // Take User-Agent from args.
    var userAgent string
    flag.StringVar(&userAgent, "ua", "gocrawler/1.0", "user agent string")

    // Parse arguments.
    flag.Parse()

    // Parse base URL.
    base, err := url.Parse(flag.Arg(0))
    if err != nil {
        fmt.Print(err.Error())
        os.Exit(1)
    }

    // Keep track of results.
    results := types.NewResultSet(*base)

    // Create channels.
    todos := make(chan string, 1000)
    found := make(chan []string, workers)

    // Keep track of Worker status.
    wg := sync.WaitGroup{}

    // Create crawl Workers.
    for i := 1; i <= workers; i++ {
        go crawl.Worker(i, userAgent, todos, found, &wg)
    }

    // Start crawl with base URL.
    fmt.Printf("Starting crawl with %d workers ...\n", workers)
    todos <- base.String()

    // Routine to process found URLs.
    go func() {
        for links := range found {
            wg.Add(1)
            for _, link := range links {
                shouldCrawl, crawlURL := results.Add(link)
                if shouldCrawl {
                    todos <- crawlURL
                }
            }
            wg.Done()
        }
    }()

    // Wait for all workers to finish.
    time.Sleep(time.Second * 5)
    wg.Wait()
    fmt.Printf("Finished! Links found:%s\n", results.String())
}
