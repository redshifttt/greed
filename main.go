package main

import (
	"fmt"
	"io"
	"bufio"
	"net/http"
	"os"
	"strings"
	"strconv"

	"github.com/mmcdole/gofeed"
    // "github.com/k3a/html2text"
)

func feedListView(retrievedFeeds []gofeed.Feed) {
    for n, feed := range retrievedFeeds {
        fmt.Printf("%d (%d) %s -- %s\n", n, len(feed.Items), feed.Title, feed.Description)
    }
}

func articleListView(selectedFeed []*gofeed.Item) {
    for n, article := range selectedFeed {
        fmt.Printf("%d %s\n", n, article.Title)
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
    retrievedFeeds := getFeedsData("input.txt")

    context := "feeds"
    lastItem := 0

    // Default action
    feedListView(retrievedFeeds)

    for {
        reader := bufio.NewReader(os.Stdin)

        fmt.Printf("greed (ctx: %s): ", context)
        text, _ := reader.ReadString('\n')
        text = text[:len(text) - 1] // Get rid of the "\n"

        command, commandArgs, hasArgs := strings.Cut(text, " ")

        switch context {
        case "feeds":
            switch command {
            case "ls":
                feedListView(retrievedFeeds)
            case "open":
                if hasArgs {
                    // God I love how Go forces you to think about errors.
                    // Means this can basically be type-checked.
                    newCommandArgs, err := strconv.Atoi(commandArgs)
                    if err != nil {
                        fmt.Println("error: argument should be int")
                        panic(err)
                    }
                    articleListView(retrievedFeeds[newCommandArgs].Items)
                    context = "articles"
                    lastItem = newCommandArgs
                }
            }
        case "articles":
            switch text {
            case "ls":
                // To get to this point we would have had to selected a feed so
                // we can deterministically use lastItem, knowing that it would
                // have been set for a feed, for printing the current feed aka
                // the one we're in. This is totally fool proof right? No idea!
                // Maybe for sanity-sake we should have a variable for feed and
                // article. Could get confusing otherwise.
                articleListView(retrievedFeeds[lastItem].Items)
            }
        }
    }
}
