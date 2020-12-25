# go-twitter v2

> It seems like every v1 is just a prototype.

One of my co-workers said that when we were working on a project, released the first version and then realizing that there were many improvements that needed to be done.  That seems to ring true, and this library is not exception.  When looking at the improvements that needed to be done, unfortunately these improvements will introduce breaking changes.  This version will focus on giving the caller more information from the callouts to allow for better response handling.  Another factor is to think about improvements, but first delivering the information first then providing more funcitonality second.

There will be `beta` releases throughout this process and if there are any new functionality requested, it will be discussed to try and get it into this version.  Version 1 will still be maintained as I believe that version 2 will need to be mature before the pervious version is even considered to be sunset.

This [project](https://github.com/g8rswimmer/go-twitter/projects/1) will track the process of this initial version.

```
go get -u github.com/g8rswimmer/go-twitter/v2
```

## Examples
Much like `v1`, there is an `_example` directory to demostrate library usage.  Refer to the [readme](./v2/_examples) for more information.

## Simple Usage
```go
import (
	"encoding/json"
	"log"
	"flag"
	"fmt"
	"net/http"
	"context"
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
	In order to run, the user will need to provide the bearer token and the list of ids.
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
		Host:   "https://api.twitter.com",
	}
	opts := twitter.TweetLookupOpts{
		Expansions:  []twitter.Expansion{twitter.ExpansionEntitiesMentionsUserName, twitter.ExpansionAuthorID},
		TweetFields: []twitter.TweetField{twitter.TweetFieldCreatedAt, twitter.TweetFieldConversationID, twitter.TweetFieldAttachments},
	}

	fmt.Println("Callout to tweet lookup callout")

	tweetDictionary, err := client.TweetLookup(context.Background(), strings.Split(*ids, ","), opts)
	if err != nil {
		log.Panicf("tweet lookup error: %v", err)
	}

	enc, err := json.MarshalIndent(tweetDictionary, "", "    ")
	if err != nil {
		log.Panic(err)
	}
	fmt.Println(string(enc))
}

```