package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"time"

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

	fmt.Println("Rate Limiting Example\n\n")

	var rateLimit *twitter.RateLimit

	rateLimit = callout(client)
	fmt.Printf("Rate Limits: %v\n", *rateLimit)
	rateLimit = callout(client)
	fmt.Printf("Rate Limits: %v\n\n", *rateLimit)

	fmt.Printf("Wait for reset %v\n", rateLimit.Reset.Time())
	time.Sleep(time.Until(rateLimit.Reset.Time()))
	fmt.Println("Rate Limits are reset")

	rateLimit = callout(client)
	fmt.Printf("Rate Limits: %v\n", *rateLimit)
}

func callout(client *twitter.Client) *twitter.RateLimit {
	opts := twitter.ListUserMembersOpts{
		MaxResults: 1,
	}
	listResponse, err := client.ListUserMembers(context.Background(), "84839422", opts)
	if err != nil {
		log.Panicf("list members error: %v", err)
	}
	return listResponse.RateLimit
}
