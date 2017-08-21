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
    // take worker count from args
    var workers int
    flag.IntVar(&workers, "w", 3, "number of workers")

    // take user-agent from args
    var userAgent string
    flag.StringVar(&userAgent, "ua", "gocrawler/1.0", "user agent string")

    // parse args
    flag.Parse()

    // parse initial url
    start, err := url.Parse(flag.Arg(0))
    if err != nil {
        fmt.Print(err.Error())
        os.Exit(1)
    }

    // url set
    results := types.NewURLSet()

    // channels
    todos := make(chan string, 1000)
    found := make(chan []string, workers)

    // worker status
    wg := sync.WaitGroup{}

    // crawler workers
    for i := 1; i <= workers; i++ {
        go crawl.Worker(i, start.Scheme, start.Hostname(), userAgent, todos, found, &wg)
    }

    // listening for crawler results
    fmt.Printf("Starting crawl with %d workers ...\n", workers)
    todos <- start.Path
    results.AddURL(start.Path)

    // main loop
    go func() {
        for links := range found {
            wg.Add(1)
            for _, link := range links {
                if !results.AddURL(link) {
                    todos <- link
                }
            }
            wg.Done()
        }
    }()

    // wait for all workers to finish
    time.Sleep(time.Second * 5)
    wg.Wait()
    fmt.Printf("Finished! Links found:%s\n", results.String())
}
