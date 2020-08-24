package main

import (
	"context"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"net/http"
	"strings"

	"github.com/g8rswimmer/go-twitter"
)

type authorize struct {
	Token string
}

func (a authorize) Add(req *http.Request) {
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", a.Token))
}

/**
	In order to run, the user will need to provide the bearer token and the list of ids.
**/
func main() {
	token := flag.String("token", "", "twitter API token")
	names := flag.String("names", "", "user names")
	flag.Parse()

	user := &twitter.User{
		Authorizer: authorize{
			Token: *token,
		},
		Client: http.DefaultClient,
		Host:   "https://api.twitter.com",
	}
	fieldOpts := twitter.UserFieldOptions{
		Expansions:  []twitter.Expansion{twitter.ExpansionEntitiesMentionsUserName, twitter.ExpansionAuthorID},
		TweetFields: []twitter.TweetField{twitter.TweetFieldCreatedAt, twitter.TweetFieldConversationID, twitter.TweetFieldAttachments},
	}

	lookups, err := user.LookupUsername(context.Background(), strings.Split(*names, ","), fieldOpts)
	var tweetErr *twitter.TweetErrorResponse
	switch {
	case errors.As(err, &tweetErr):
		printTweetError(tweetErr)
	case err != nil:
		fmt.Println(err)
	default:
		for _, lookup := range lookups {
			printTweetLookup(lookup)
			fmt.Println()
		}
	}
}

func printTweetLookup(lookup twitter.UserLookup) {
	enc, err := json.MarshalIndent(lookup, "", "    ")
	if err != nil {
		panic(err)
	}
	fmt.Println(string(enc))
}

func printTweetError(tweetErr *twitter.TweetErrorResponse) {
	enc, err := json.MarshalIndent(tweetErr, "", "    ")
	if err != nil {
		panic(err)
	}
	fmt.Println(string(enc))
}
