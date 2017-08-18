package main

import (
    "flag"
    "fmt"
    "net/url"
    "os"

    "github.com/evilscott/gocrawler/crawl"
    "github.com/evilscott/gocrawler/types"
)

func main() {
    // take worker count from args
    var workers int
    flag.IntVar(&workers, "workers", 3, "number of workers")
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
    found := make(chan string, 100)

    // crawler workers
    for i := 1; i <= workers; i++ {
        go crawl.Worker(i, start.Scheme, start.Hostname(), todos, found)
    }

    // listening for crawler results
    fmt.Printf("Starting crawl with %d workers ...\n", workers)
    todos <- start.Path
    results.AddURL(start.Path)

    // main loop
    for {
        select {
        case link := <-found:
            if !results.AddURL(link) {
                todos <- link
            }
        default:
            // TODO check for end condition
        }
    }

    // wait for all workers to finish
    fmt.Printf("Finished! Links found:%s\n", results.String())
}
