package main

import (
	"fmt"
	"io"
	"bufio"
	"net/http"
	"os"
	"strings"
	"strconv"
	"log"

	"github.com/mmcdole/gofeed"
    // "github.com/k3a/html2text"
)

func feedListView(retrievedFeeds []gofeed.Feed) {
    fmt.Printf("Feeds view, %d feeds\n", len(retrievedFeeds))

    for n, feed := range retrievedFeeds {
        fmt.Printf("%d (%d) %s -- %s\n", n, len(feed.Items), feed.Title, feed.Description)
    }
}

func articleListView(selectedFeed gofeed.Feed) {
    fmt.Printf("Article view for feed %s\n", selectedFeed.Title)

    for n, article := range selectedFeed.Items {
        fmt.Printf("%d %s %s\n", n, article.Published, article.Title)
    }
}

func getFeedsData(urlsFile string) []gofeed.Feed {
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

        fmt.Printf("[%d/%d] parsing data for feed %s ... ", n + 1, totalUrls, url)

        fp := gofeed.NewParser()
        feed, err := fp.ParseString(body_text)
        if err != nil {
            panic(err)
        }
        retrievedFeeds = append(retrievedFeeds, *feed)

        fmt.Printf("done!\n")
    }

    return retrievedFeeds
}

func main() {
    retrievedFeeds := getFeedsData("input.txt")

    context := "feeds"
    lastFeedIndex := 0
    var selectedFeed gofeed.Feed

    // Default action
    feedListView(retrievedFeeds)

    for {
        reader := bufio.NewReader(os.Stdin)

        fmt.Printf("[greed view: %s] ", context)
        text, _ := reader.ReadString('\n')
        text = text[:len(text) - 1] // Get rid of the "\n"

        command, commandArgs, hasArgs := strings.Cut(text, " ")

        // TODO: divide this up into functions.
        switch context {
        case "feeds":
            switch command {
            case "ls":
                feedListView(retrievedFeeds)
            case "open":
                if hasArgs {
                    newCommandArgs, err := strconv.Atoi(commandArgs)
                    if err != nil {
                        log.Fatalf("error: '%s' should be of type int.\n", err)
                    }
                    selectedFeed = retrievedFeeds[newCommandArgs]

                    articleListView(selectedFeed)
                    context = "articles"
                    lastFeedIndex = newCommandArgs
                }
            default:
                fmt.Printf("error: '%s' is not a valid command in the feeds context.\n", command)
            }
        case "articles":
            switch text {
            case "ls":
                selectedFeed = retrievedFeeds[lastFeedIndex]
                articleListView(selectedFeed)
            case "back":
                feedListView(retrievedFeeds)
                context = "feeds"
            default:
                fmt.Printf("error: '%s' is not a valid command in the articles context.\n", command)
            }
        }
    }
}
