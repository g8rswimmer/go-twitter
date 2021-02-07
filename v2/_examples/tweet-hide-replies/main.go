package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"strconv"

	twitter "github.com/g8rswimmer/go-twitter/v2"
)

type authorize struct {
	Token string
}

func (a authorize) Add(req *http.Request) {
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", a.Token))
}

/**
	In order to run, the user will need to provide the bearer token and the list of tweet ids.
**/
func main() {
	token := flag.String("token", "", "twitter API token")
	ids := flag.String("id", "", "twitter id")
	hide := flag.String("hide", "", "Hide replies")
	flag.Parse()

	hideBool, err := strconv.ParseBool(hide)
	if err != nil {
		log.Panicf("tweet hide error: %v", err)
	}

	client := &twitter.Client{
		Authorizer: authorize{
			Token: *token,
		},
		Client: http.DefaultClient,
		Host:   "https://api.twitter.com",
	}

	fmt.Println("Callout to tweet hide replies")

	if err := client.TweetHideReplies(context.Background(), *id, hideBool); err != nil {
		log.Panicf("tweet hide replies error: %v", err)
	}

	fmt.Printf("tweet %s hide replies %v", *ids, hideBool)
}
