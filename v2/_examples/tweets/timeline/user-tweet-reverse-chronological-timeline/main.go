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
	userID := flag.String("user_id", "", "user id")
	flag.Parse()

	client := &twitter.Client{
		Authorizer: authorize{
			Token: *token,
		},
		Client: http.DefaultClient,
		Host:   "https://api.twitter.com",
	}
	opts := twitter.UserTweetChronologicalReverseTimelineOpts{
		TweetFields: []twitter.TweetField{twitter.TweetFieldCreatedAt, twitter.TweetFieldAuthorID, twitter.TweetFieldConversationID, twitter.TweetFieldPublicMetrics, twitter.TweetFieldContextAnnotations},
		UserFields:  []twitter.UserField{twitter.UserFieldUserName},
		Expansions:  []twitter.Expansion{twitter.ExpansionAuthorID},
		MaxResults:  5,
	}

	fmt.Println("Callout to tweet user reverse chronological timeline callout")

	timeline, err := client.UserTweetReverseChronologicalTimeline(context.Background(), *userID, opts)
	if err != nil {
		log.Panicf("user reverse chronological timeline error: %v", err)
	}

	dictionaries := timeline.Raw.TweetDictionaries()

	enc, err := json.MarshalIndent(dictionaries, "", "    ")
	if err != nil {
		log.Panic(err)
	}
	fmt.Println(string(enc))

	metaBytes, err := json.MarshalIndent(timeline.Meta, "", "    ")
	if err != nil {
		log.Panic(err)
	}
	fmt.Println(string(metaBytes))
}
