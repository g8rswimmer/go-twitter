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
	In order to run, the user will need to provide the bearer token and the list of user ids.
**/
func main() {
	token := flag.String("token", "", "twitter API token")
	id := flag.String("id", "", "user id")
	flag.Parse()

	client := &twitter.Client{
		Authorizer: authorize{
			Token: *token,
		},
		Client: http.DefaultClient,
		Host:   "https://api.twitter.com",
	}
	opts := twitter.UserFollowingLookupOpts{
		Expansions:  []twitter.Expansion{twitter.ExpansionPinnedTweetID},
		TweetFields: []twitter.TweetField{twitter.TweetFieldContextAnnotations},
	}

	fmt.Println("Callout to user following lookup callout")

	userResponse, err := client.UserFollowingLookup(context.Background(), *id, opts)
	if err != nil {
		log.Panicf("user following lookup error: %v", err)
	}

	dictionaries := userResponse.Raw.UserDictionaries()

	enc, err := json.MarshalIndent(dictionaries, "", "    ")
	if err != nil {
		log.Panic(err)
	}
	fmt.Println(string(enc))

	metaBytes, err := json.MarshalIndent(userResponse.Meta, "", "    ")
	if err != nil {
		log.Panic(err)
	}
	fmt.Println(string(metaBytes))

}
