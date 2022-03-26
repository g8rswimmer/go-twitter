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
	id := flag.String("id", "", "twitter id")
	flag.Parse()

	client := &twitter.Client{
		Authorizer: authorize{
			Token: *token,
		},
		Client: http.DefaultClient,
		Host:   "https://api.twitter.com",
	}
	opts := twitter.QuoteTweetsLookupOpts{
		Expansions:  []twitter.Expansion{twitter.ExpansionAuthorID},
		TweetFields: []twitter.TweetField{twitter.TweetFieldCreatedAt, twitter.TweetFieldConversationID, twitter.TweetFieldAttachments, twitter.TweetFieldAuthorID, twitter.TweetFieldPublicMetrics},
		UserFields:  []twitter.UserField{twitter.UserFieldUserName},
	}

	fmt.Println("Callout to quote tweet lookup callout")

	tweetResponse, err := client.QuoteTweetsLookup(context.Background(), *id, opts)
	if err != nil {
		log.Panicf("tweet quote lookup error: %v", err)
	}

	dictionaries := tweetResponse.Raw.TweetDictionaries()

	enc, err := json.MarshalIndent(dictionaries, "", "    ")
	if err != nil {
		log.Panic(err)
	}
	fmt.Println(string(enc))

	enc, err = json.MarshalIndent(tweetResponse.Meta, "", "    ")
	if err != nil {
		log.Panic(err)
	}
	fmt.Println(string(enc))
}
