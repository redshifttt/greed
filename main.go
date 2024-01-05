package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	// "math"
	// "slices"
	"github.com/mmcdole/gofeed"
    // "github.com/k3a/html2text"
)

func feedListView(retrievedFeeds []gofeed.Feed) {
    for n, feed := range retrievedFeeds {
        fmt.Printf("%d (%d) %s -- %s\n", n, len(feed.Items), feed.Title, feed.Description)
    }
    fmt.Println(retrievedFeeds[0].Items[0].Title)
}

func parseUrls(urlsFile string) []gofeed.Feed {
    data, err := os.ReadFile(urlsFile)

    if err != nil {
        panic(err)
    }

    urls := strings.Split(string(data), "\n")
    totalUrls := len(urls) - 1

    retrievedFeeds := make([]gofeed.Feed, 0)

    for n, url := range urls {
        if url == "" {
            continue
        }

        fmt.Printf("[%d/%d] getting data for feed %s ... ", n + 1, totalUrls, url)

        resp, err := http.Get(url)
        if err != nil {
            panic(err)
        }
        defer resp.Body.Close()

        fmt.Printf("%s\n", resp.Status)

        body, err := io.ReadAll(resp.Body)
        body_text := string(body)

        fp := gofeed.NewParser()
        feed, err := fp.ParseString(body_text)
        if err != nil {
            panic(err)
        }
        retrievedFeeds = append(retrievedFeeds, *feed)
    }

    return retrievedFeeds
}

func main() {
    retrievedFeeds := parseUrls("input.txt")

    feedListView(retrievedFeeds)
}
