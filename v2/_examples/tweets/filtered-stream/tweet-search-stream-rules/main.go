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

/**
	In order to run, the user will need to provide the bearer token and the list of tweet ids.
**/
func main() {
	token := flag.String("token", "", "twitter API token")
	flag.Parse()

	client := &twitter.Client{
		Authorizer: authorize{
			Token: *token,
		},
		Client: http.DefaultClient,
		Host:   "https://api.twitter.com",
	}

	fmt.Println("Callout to tweet search stream rules callout")

	searchStreamRules, err := client.TweetSearchStreamRules(context.Background(), []twitter.TweetSearchStreamRuleID{})
	if err != nil {
		log.Panicf("tweet search stream rule callout error: %v", err)
	}

	enc, err := json.MarshalIndent(searchStreamRules, "", "    ")
	if err != nil {
		log.Panic(err)
	}
	fmt.Println(string(enc))

}
