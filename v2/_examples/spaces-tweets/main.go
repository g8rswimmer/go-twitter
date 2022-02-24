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
	id := flag.String("id", "", "space id")
	flag.Parse()

	client := &twitter.Client{
		Authorizer: authorize{
			Token: *token,
		},
		Client: http.DefaultClient,
		Host:   "https://api.twitter.com",
	}
	opts := twitter.SpaceTweetsLookupOpts{
		TweetFields: []twitter.TweetField{twitter.TweetFieldAuthorID},
	}

	fmt.Println("Callout to spaces tweets callout")

	spaceResponse, err := client.SpaceTweetsLookup(context.Background(), *id, opts)
	if err != nil {
		log.Panicf("spaces tweets error: %v", err)
	}

	enc, err := json.MarshalIndent(spaceResponse, "", "    ")
	if err != nil {
		log.Panic(err)
	}
	fmt.Println(string(enc))
}
