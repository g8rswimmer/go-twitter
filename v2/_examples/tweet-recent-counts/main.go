package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"

	twitter "github.com/g8rswimmer/go-twitter/v2"
)

type authorize struct {
	Token string
}

func (a authorize) Add(req *http.Request) {
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", a.Token))
}
func main() {
	token := flag.String("token", "", "twitter API token")
	query := flag.String("query", "", "twitter query")
	flag.Parse()

	client := &twitter.Client{
		Authorizer: authorize{
			Token: *token,
		},
		Client: http.DefaultClient,
		Host:   "https://api.twitter.com",
	}
	opts := twitter.TweetRecentCountsOpts{
		Granularity: twitter.GranularityHour,
	}

	fmt.Println("Callout to tweet recent counts callout")

	tweetResponse, err := client.TweetRecentCounts(context.Background(), *query, opts)
	if err != nil {
		log.Panicf("tweet recent counts error: %v", err)
	}

	enc, err := json.MarshalIndent(tweetResponse.TweetCounts, "", "    ")
	if err != nil {
		log.Panic(err)
	}
	fmt.Println(string(enc))

	metaBytes, err := json.MarshalIndent(tweetResponse.Meta, "", "    ")
	if err != nil {
		log.Panic(err)
	}
	fmt.Println(string(metaBytes))
}
