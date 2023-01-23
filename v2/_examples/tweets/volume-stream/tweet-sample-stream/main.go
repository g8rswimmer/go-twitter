package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	twitter "github.com/g8rswimmer/go-twitter/v2"
)

type authorize struct {
	Token string
}

func (a authorize) Add(req *http.Request) {
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", a.Token))
}

/*
In order to run, the user will need to provide the bearer token and a file name to be used as a log.
*/
func main() {
	token := flag.String("token", "", "twitter API token")
	output := flag.String("output", "", "output")
	flag.Parse()

	client := &twitter.Client{
		Authorizer: authorize{
			Token: *token,
		},
		Client: http.DefaultClient,
		Host:   "https://api.twitter.com",
	}
	opts := twitter.TweetSampleStreamOpts{}

	fmt.Println("Callout to tweet sample stream callout")

	outputFile, err := os.Create(*output)
	if err != nil {
		log.Panicf("tweet stream output file error %v", err)
	}
	defer outputFile.Close()

	ctx := context.Background()
	tweetStream, err := client.TweetSampleStreamV2(ctx, opts)
	if err != nil {
		log.Panicf("tweet sample callout error: %v", err)
	}

	ctx, cancel := signal.NotifyContext(ctx, syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	err = tweetStream.Run(ctx, twitter.TweetStreamV2Options{
		OnTweet: func(tm *twitter.StreamedTweet) {
			tmb, err := json.Marshal(tm)
			if err != nil {
				fmt.Printf("error decoding tweet message %v", err)
			}
			outputFile.WriteString(fmt.Sprintf("tweet: %s\n\n", string(tmb)))
			outputFile.Sync()
			fmt.Println("tweet")
		},
		OnSystemMessage: func(kind twitter.SystemMessageType, msg *twitter.SystemMessage) {
			smb, err := json.Marshal(msg)
			if err != nil {
				fmt.Printf("error decoding system message %v", err)
			}
			outputFile.WriteString(fmt.Sprintf("system[%s]: %s\n\n", kind, string(smb)))
			outputFile.Sync()
			fmt.Println("system")
		},
		OnTransientError: func(err error) {
			outputFile.WriteString(fmt.Sprintf("error: %v\n\n", err))
			outputFile.Sync()
			fmt.Println("error")
		},
	})

	if ctx.Err() != nil {
		fmt.Println("closing")
		return
	}

	if err != nil {
		log.Panicf("tweet stream ended: %s", err)
	}
}
