# go-twitter v2

> It seems like every v1 is just a prototype.

One of my co-workers said that when we were working on a project, released the first version and then realizing that there were many improvements that needed to be done.  That seems to ring true, and this library is not exception.  When looking at the improvements that needed to be done, unfortunately these improvements will introduce breaking changes.  This version will focus on giving the caller more information from the callouts to allow for better response handling.  Another factor is to think about improvements, but first delivering the information first then providing more funcitonality second.

There will be `beta` releases throughout this process and if there are any new functionality requested, it will be discussed to try and get it into this version.  Version 1 will still be maintained as I believe that version 2 will need to be mature before the pervious version is even considered to be sunset.

This [project](https://github.com/g8rswimmer/go-twitter/projects/1) will track the process of this initial version.

```
go get -u github.com/g8rswimmer/go-twitter/v2
```

## Changes
The following are changes between `v1` and `v2` of the library.
*  One structure for all endpoint callouts.  In `v1` there were two structures, `Tweet` and `User`, to preform callouts.  This required management of two structures and knowledge which structure contained the desired method callouts.  At the time, the grouping looked logical.  However with the addtion of the timeline endpoints, it makes more sense to have a single struture `Client` to handle all of the callouts.  If the user would like to separate the functionality, interfaces can be used to achieve this.
*  Endpoint methods will return the entire response.  One of the major drawbacks of `v1` was the object returned was not the entire response sent by the callout.  For example, the `errors` object in the response is included in `v1` response which does not allow the caller to properly handle partial errors.  In `v2`, the first focus is to return the response from twitter to allow the caller to use it as they see fit.  However, it does not mean that methods can not be added to the response object to provide groupings, etc.

## Features 
Here are the current twitter `v2` API features supported:
*  [Tweet Lookup](https://developer.twitter.com/en/docs/twitter-api/tweets/lookup/introduction)
    * [example](./_examples/tweet-lookup)
* [Tweet Dictionary](https://developer.twitter.com/en/docs/twitter-api/data-dictionary/object-model/tweet) - the tweet and all of its references related to it
*  [User Lookup](https://developer.twitter.com/en/docs/twitter-api/users/lookup/introduction)
    * [example](./_examples/user-lookup)
    * [example: by usernames](./_examples/username-lookup)
    * [example: authorized](./_examples/auth-user-lookup)
* [User Dictionary](https://developer.twitter.com/en/docs/twitter-api/data-dictionary/object-model/user) - the user and all of its references related to it
* [Tweet Counts](https://developer.twitter.com/en/docs/twitter-api/tweets/counts/introduction)
    * [example](./_examples/tweet-recent-counts)
* [Manage Tweet](https://developer.twitter.com/en/docs/twitter-api/tweets/manage-tweets/introduction)
    * [create example](./_examples/tweet-create)
    * [delete example](./_examples/tweet-delete)
* [Manage Retweet](https://developer.twitter.com/en/docs/twitter-api/tweets/retweets/introduction)
    * [retweet example](./_examples/user-retweet)
    * [delete retweet example](./_examples/user-delete-retweet)
    * [retweet lookup example](./_examples/user-retweet-lookup)
* [User Blocks](https://developer.twitter.com/en/docs/twitter-api/users/blocks/introduction)
	* [blocks lookup example](./_examples/user-blocks-lookup)
	* [blocks example](./_examples/user-blocks)
	* [delete blocks example](./_examples/user-delete-blocks)
* [User Mutes](https://developer.twitter.com/en/docs/twitter-api/users/mutes/introduction)
	* [mutes lookup example](./_examples/user-mutes-lookup)
	* [mutes example](./_examples/user-mutes)
	* [delete mutes example](./_examples/user-delete-mutes)
* [User Follows](https://developer.twitter.com/en/docs/twitter-api/users/follows/introduction)
	* [user following](./_examples/user-following-lookup)
	* [user followers](./_examples/user-followers-lookup)
	* [follows example](./_examples/user-follows)
	* [delete follows example](./_examples/user-delete-follows)

## Examples
Much like `v1`, there is an `_example` directory to demostrate library usage.  Refer to the [readme](./_examples) for more information.

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