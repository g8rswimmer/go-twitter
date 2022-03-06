package main

import (
	"context"
	"encoding/json"
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
	ids := flag.String("ids", "", "space ids")
	flag.Parse()

	client := &twitter.Client{
		Authorizer: authorize{
			Token: *token,
		},
		Client: http.DefaultClient,
		Host:   "https://api.twitter.com",
	}
	opts := twitter.SpacesByCreatorLookupOpts{
		Expansions:  []twitter.Expansion{twitter.ExpansionCreatorID},
		SpaceFields: []twitter.SpaceField{twitter.SpaceFieldHostIDs},
		UserFields:  []twitter.UserField{twitter.UserFieldUserName},
	}

	fmt.Println("Callout to spaces by creator callout")

	spaceResponse, err := client.SpacesByCreatorLookup(context.Background(), strings.Split(*ids, ","), opts)
	if err != nil {
		log.Panicf("spaces by creator error: %v", err)
	}

	enc, err := json.MarshalIndent(spaceResponse, "", "    ")
	if err != nil {
		log.Panic(err)
	}
	fmt.Println(string(enc))
}
