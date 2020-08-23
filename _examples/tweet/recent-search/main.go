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
	query := flag.String("query", "", "twitter query")
	flag.Parse()

	tweet := &twitter.Tweet{
		Authorizer: authorize{
			Token: *token,
		},
		Client: http.DefaultClient,
		Host:   "https://api.twitter.com",
	}
	parameters := twitter.TweetRecentSearchOptions{
		TweetFields: []twitter.TweetField{twitter.TweetFieldCreatedAt, twitter.TweetFieldConversationID, twitter.TweetFieldLanguage},
	}

	recentSearch, err := tweet.RecentSearch(context.Background(), *query, parameters)
	var tweetErr *twitter.TweetErrorResponse
	switch {
	case errors.As(err, &tweetErr):
		printTweetError(tweetErr)
	case err != nil:
		fmt.Println(err)
	default:
		printRecentSearch(recentSearch)
	}
}

func printRecentSearch(recentSearch *twitter.TweetRecentSearch) {
	for _, lookup := range recentSearch.LookUps {
		enc, err := json.MarshalIndent(lookup, "", "    ")
		if err != nil {
			panic(err)
		}
		fmt.Println(string(enc))
	}
	enc, err := json.MarshalIndent(recentSearch.Meta, "", "    ")
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
