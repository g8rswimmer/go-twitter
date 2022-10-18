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

/**
	In order to run, the user will need to provide the bearer token and a file name to be used as a log.
**/
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
	opts := twitter.TweetSearchStreamOpts{}

	fmt.Println("Callout to tweet search stream callout")

	outputFile, err := os.Create(*output)
	if err != nil {
		log.Panicf("tweet stream output file error %v", err)
	}
	defer outputFile.Close()

	tweetStream, err := client.TweetSearchStream(context.Background(), opts)
	if err != nil {
		log.Panicf("tweet sample callout error: %v", err)
	}
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL)
	func() {
		defer tweetStream.Close()
		for {
			select {
			case <-ch:
				fmt.Println("closing")
				return
			case tm := <-tweetStream.Tweets():
				tmb, err := json.Marshal(tm)
				if err != nil {
					fmt.Printf("error decoding tweet message %v", err)
				}
				outputFile.WriteString(fmt.Sprintf("tweet: %s\n\n", string(tmb)))
				outputFile.Sync()
				fmt.Println("tweet")
			case sm := <-tweetStream.SystemMessages():
				smb, err := json.Marshal(sm)
				if err != nil {
					fmt.Printf("error decoding system message %v", err)
				}
				outputFile.WriteString(fmt.Sprintf("system: %s\n\n", string(smb)))
				outputFile.Sync()
				fmt.Println("system")
			case de := <-tweetStream.DisconnectionError():
				ded, err := json.Marshal(de)
				if err != nil {
					fmt.Printf("error decoding disconnect message %v", err)
				}
				outputFile.WriteString(fmt.Sprintf("disconnect: %s\n\n", string(ded)))
				outputFile.Sync()
				fmt.Println("disconnect")
			case strErr := <-tweetStream.Err():
				outputFile.WriteString(fmt.Sprintf("error: %v\n\n", strErr))
				outputFile.Sync()
				fmt.Println("error")
			default:
			}
			if tweetStream.Connection() == false {
				fmt.Println("connection lost")
				return
			}
		}
	}()
}
