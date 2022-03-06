package main

import (
	"context"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"log"
	"net/http"
	"strings"

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
	ids := flag.String("ids", "", "twitter ids")
	flag.Parse()

	client := &twitter.Client{
		Authorizer: authorize{
			Token: *token,
		},
		Client: http.DefaultClient,
		Host:   "https://api.twitter.com/",
	}
	opts := twitter.TweetLookupOpts{
		Expansions:  []twitter.Expansion{twitter.ExpansionEntitiesMentionsUserName, twitter.ExpansionAuthorID},
		TweetFields: []twitter.TweetField{twitter.TweetFieldCreatedAt, twitter.TweetFieldConversationID, twitter.TweetFieldAttachments},
	}

	fmt.Println("Twitter HTTP Error Example")

	_, err := client.TweetLookup(context.Background(), strings.Split(*ids, ","), opts)
	var httpErr *twitter.HTTPError
	switch {
	case err == nil:
		log.Panic("there should be an error")
	case errors.As(err, &httpErr):
		enc, err := json.MarshalIndent(httpErr, "", "    ")
		if err != nil {
			log.Panic(err)
		}
		fmt.Println(string(enc))
	default:
		log.Panicf("wrong error: %v", err)
	}
}
