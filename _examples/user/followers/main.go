package main

import (
	"context"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"net/http"

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
	id := flag.String("id", "", "user id")
	flag.Parse()

	user := &twitter.User{
		Authorizer: authorize{
			Token: *token,
		},
		Client: http.DefaultClient,
		Host:   "https://api.twitter.com",
	}
	fieldOpts := twitter.UserFollowOptions{
		Expansions:  []twitter.Expansion{twitter.ExpansionPinnedTweetID},
		TweetFields: []twitter.TweetField{twitter.TweetFieldCreatedAt, twitter.TweetFieldContextAnnotations},
		UserFields:  []twitter.UserField{twitter.UserFieldCreatedAt},
		MaxResults:  10,
	}

	userFollowLookup, err := user.LookupFollowers(context.Background(), *id, fieldOpts)
	var tweetErr *twitter.TweetErrorResponse
	switch {
	case errors.As(err, &tweetErr):
		printTweetError(tweetErr)
	case err != nil:
		fmt.Println(err)
	default:
		for _, lookup := range userFollowLookup.Lookups {
			printTweetLookup(lookup)
			fmt.Println()
		}
		printTweetLookupErrors(userFollowLookup.Errors)
		printMeta(userFollowLookup.Meta)
	}

}

func printMeta(meta *twitter.UserFollowMeta) {
	enc, err := json.MarshalIndent(meta, "", "    ")
	if err != nil {
		panic(err)
	}
	fmt.Println(string(enc))
}

func printTweetLookupErrors(errs []twitter.ErrorObj) {
	for _, err := range errs {
		enc, e := json.MarshalIndent(err, "", "    ")
		if e != nil {
			panic(e)
		}
		fmt.Println(string(enc))
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
