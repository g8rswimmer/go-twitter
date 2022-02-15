![](https://img.shields.io/endpoint?url=https%3A%2F%2Ftwbadges.glitch.me%2Fbadges%2Fv2)
[![golangci-lint](https://github.com/g8rswimmer/go-twitter/actions/workflows/golangci-lint.yml/badge.svg)](https://github.com/g8rswimmer/go-twitter/actions/workflows/golangci-lint.yml)
[![go-test](https://github.com/g8rswimmer/go-twitter/actions/workflows/go-test.yml/badge.svg)](https://github.com/g8rswimmer/go-twitter/actions/workflows/go-test.yml)
[![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg)](https://opensource.org/licenses/MIT)

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

### Breaking Changes
These are some breaking changes and what release they are from.  These type of changes will try to be avoided but if necessary for a better library, it will be done.

#### v2.0.0-beta11
* The client callout for tweet hide reply, `TweetHideReplies`, has been changed to return the response instead of just an error.  This allow for the data and the rate limits of the callout to be returned. 
##### Migration
```go
	// old way
	err := client.TweetHideReplies(context.Background(), id, true)
	if err != nil {
		// handle error
	}
```
```go
	// new way
	hideResponse, err := client.TweetHideReplies(context.Background(), id, true)
	if err != nil {
		// handle error
	}
```

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
* [Tweet Likes](https://developer.twitter.com/en/docs/twitter-api/tweets/likes/introduction)
	* [user likes lookup](./_examples/user-likes-lookup)
	* [tweets likes lookup](./_examples/tweet-likes-lookup)
* [Tweet Sample Stream](https://developer.twitter.com/en/docs/twitter-api/tweets/volume-streams/introduction)
	* [tweet sample stream](./_examples/tweet-sample-stream)
* [Tweet Search Stream](https://developer.twitter.com/en/docs/twitter-api/tweets/filtered-stream/introduction)
	* [tweet search stream add rules](./_examples/tweet-search-stream-add-rule)
	* [tweet search stream delete rules](./_examples/tweet-search-stream-delete-rule)
	* [tweet search stream rules](./_examples/tweet-search-stream-rules)
	* [tweet search stream](./_examples/tweet-search-stream)
* [List Lookup](https://developer.twitter.com/en/docs/twitter-api/lists/list-lookup/introduction)
	* [list lookup](./_examples/list-lookup)
	* [user list lookup](./_examples/user-list-lookup)
* [List Tweets Lookup](https://developer.twitter.com/en/docs/twitter-api/lists/list-tweets/introduction)
	* [list tweet lookup](./_examples/list-tweet-lookup)
* [Manage Lists](https://developer.twitter.com/en/docs/twitter-api/lists/manage-lists/introduction)
	* [create list](./_examples/list-create)
	* [update list](./_examples/list-update)
	* [delete list](./_examples/list-delete)
* [List Members](https://developer.twitter.com/en/docs/twitter-api/lists/list-members/introduction)
	* [list add members](./_examples/list-add-member)
	* [list remove members](./_examples/list-remove-member)
	* [list members](./_examples/list-members)
	* [list membersships](./_examples/list-memberships)
* [Pinned Lists](https://developer.twitter.com/en/docs/twitter-api/lists/pinned-lists/introduction)
	* [user pin list](./examples/user-pin-list)
	* [user unpin list](./examples/user-unpin-list)
	* [user pinned lists](./examples/user-pinned-lists)
* [List Follows](https://developer.twitter.com/en/docs/twitter-api/lists/list-follows/introduction)
	* [user follow list](./examples/user-follow-list)
	* [user unfollow list](./examples/user-unfollow-list)
	* [user followed lists](./examples/user-followed-lists)
	* [list followers](./examples/list-followers)

## Rate Limiting
With each response, the rate limits from the response header is returned.  This allows the caller to manage any limits that are imposed.  Along with the response, errors that are returned may have rate limits as well.  If the error occurs after the request is sent, then rate limits may apply and are returned.

There is an example of rate limiting from a response [here](./examples/rate-limit).

This is an example of a twitter callout and if the limits have been reached, then it will backoff and try again.
```go
func TweetLikes(ctx context.Context, id string, client *twitter.Client) (*twitter.TweetLikesLookupResponse, error) {
	var er *ErrorResponse

	opts := twitter.ListUserMembersOpts{
		MaxResults: 1,
	}
	tweetResponse, err := client.TweetLikesLookup(ctx, id, opts)

	if rateLimit, has := twitter.RateLimitFromError(err); has && rateLimit.Remaining == 0 {
		time.Sleep(time.Until(rateLimit.Reset.Time()))
		return client.TweetLikesLookup(ctx, id, opts)
	}
	return tweetResponse, err
}
```

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