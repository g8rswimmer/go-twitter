![](https://img.shields.io/endpoint?url=https%3A%2F%2Ftwbadges.glitch.me%2Fbadges%2Fv2)
[![golangci-lint](https://github.com/g8rswimmer/go-twitter/actions/workflows/golangci-lint.yml/badge.svg)](https://github.com/g8rswimmer/go-twitter/actions/workflows/golangci-lint.yml)
[![go-test](https://github.com/g8rswimmer/go-twitter/actions/workflows/go-test.yml/badge.svg)](https://github.com/g8rswimmer/go-twitter/actions/workflows/go-test.yml)
[![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg)](https://opensource.org/licenses/MIT)

# go-twitter v2

> It seems like every v1 is just a prototype.

One of my co-workers said that when we were working on a project, released the first version and then realizing that there were many improvements that needed to be done.  That seems to ring true, and this library is not exception.  When looking at the improvements that needed to be done, unfortunately these improvements will introduce breaking changes.  This version will focus on giving the caller more information from the callouts to allow for better response handling.  Another factor is to think about improvements, but first delivering the information first then providing more functionality second.

There will be `beta` releases throughout this process and if there are any new functionality requested, it will be discussed to try and get it into this version.  Version 1 will still be maintained as I believe that version 2 will need to be mature before the pervious version is even considered to be sunset.

This [project](https://github.com/g8rswimmer/go-twitter/projects/1) will track the process of this initial version.

```
go get -u github.com/g8rswimmer/go-twitter/v2
```

## Table Of Contents
*  [Changes](#changes) Gives an outline of the changes between `v1` and `v2`
    * [Breaking Changes](#breaking-changes)
*  [Features](#features) Outlines the twitter v2 APIs supported
    * [Tweets](#tweets)
	* [Users](#users)
	* [Spaces](#spaces)
	* [Lists](#lists)
	* [Compliance](#compliance)
*  [Rate Limiting](#rate-limiting) Explains how API rate limits are supported
*  [Error Handling](#error-handling) Explains how the different types of errors are handled by the library
    * [Parameter Errors](#parameter-errors)
	* [Callout Errors](#callout-errors)
	* [Response Decode Errors](#response-decode-errors)
	* [Twitter HTTP Response Errors](#twitter-http-response-errors)
	* [Twitter Partial Errors](#twitter-partial-errors)
*  [Examples](#examples) Brief overview of where the examples are contained

## Changes
The following are changes between `v1` and `v2` of the library.
*  One structure for all endpoint callouts.  In `v1` there were two structures, `Tweet` and `User`, to preform callouts.  This required management of two structures and knowledge which structure contained the desired method callouts.  At the time, the grouping looked logical.  However with the addition of the timeline endpoints, it makes more sense to have a single structure `Client` to handle all of the callouts.  If the user would like to separate the functionality, interfaces can be used to achieve this.
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
#### v2.0.0-beta13
* There was a typo in the user retweet lookup options and it was corrected.
```
UserRetweetLookuoOpts -> UserRetweetLookupOpts
```
#### v2.0.0-beta16
* There was a typo in the user following metadata.
```
UserFollowinghMeta -> UserFollowingMeta
```
* There was a typo in the delete user likes response
```
DeteleUserLikesResponse -> DeleteUserLikesResponse
```
* There was a typo in the `EntityURLObj`
```
EntityURLObj.Desription -> EntityURLObj.Description
```
* There was a typo in the `TweetObj` and the JSON tag
```
TweetObj.PossibySensitive -> TweetObj.PossiblySensitive

json: possiby_sensitive -> possibly_sensitive
```

## Features 
Here are the current twitter `v2` API features supported.

### Tweets
The following APIs are supported, with the examples [here](./_examples/tweets)

* [Lookup](https://developer.twitter.com/en/docs/twitter-api/tweets/lookup/introduction)
* [Counts](https://developer.twitter.com/en/docs/twitter-api/tweets/counts/introduction)
* [Manage Tweets](https://developer.twitter.com/en/docs/twitter-api/tweets/manage-tweets/introduction)
* [Retweets](https://developer.twitter.com/en/docs/twitter-api/tweets/retweets/introduction)
* [Likes](https://developer.twitter.com/en/docs/twitter-api/tweets/likes/introduction)
* [Volume Stream](https://developer.twitter.com/en/docs/twitter-api/tweets/volume-streams/introduction)
* [Filtered Stream](https://developer.twitter.com/en/docs/twitter-api/tweets/filtered-stream/introduction)
* [Timelines](https://developer.twitter.com/en/docs/twitter-api/tweets/timelines/introduction)
* [Hide Replies](https://developer.twitter.com/en/docs/twitter-api/tweets/hide-replies/introduction)
* [Search](https://developer.twitter.com/en/docs/twitter-api/tweets/search/introduction)
* [Quote Tweets](https://developer.twitter.com/en/docs/twitter-api/tweets/quote-tweets/introduction)
* [Bookmarks](https://developer.twitter.com/en/docs/twitter-api/tweets/bookmarks/introduction)

### Users
The following APIs are supported, with the examples [here](./_examples/users)

* [Lookup](https://developer.twitter.com/en/docs/twitter-api/users/lookup/introduction)
* [Blocks](https://developer.twitter.com/en/docs/twitter-api/users/blocks/introduction)
* [Mutes](https://developer.twitter.com/en/docs/twitter-api/users/mutes/introduction)
* [Follows](https://developer.twitter.com/en/docs/twitter-api/users/follows/introduction)

### Spaces
The following APIs are supported, with the examples [here](./_examples/spaces)

* [Spaces Lookup](https://developer.twitter.com/en/docs/twitter-api/spaces/lookup/introduction)
* [Spaces Search](https://developer.twitter.com/en/docs/twitter-api/spaces/search/introduction)

### Lists
The following APIs are supported, with the examples [here](./_examples/lists)

* [List Lookup](https://developer.twitter.com/en/docs/twitter-api/lists/list-lookup/introduction)
* [List Tweets Lookup](https://developer.twitter.com/en/docs/twitter-api/lists/list-tweets/introduction)
* [Manage Lists](https://developer.twitter.com/en/docs/twitter-api/lists/manage-lists/introduction)
* [List Members](https://developer.twitter.com/en/docs/twitter-api/lists/list-members/introduction)
* [Pinned Lists](https://developer.twitter.com/en/docs/twitter-api/lists/pinned-lists/introduction)
* [List Follows](https://developer.twitter.com/en/docs/twitter-api/lists/list-follows/introduction)

### Compliance
The following APIs are supported, with the examples [here](./_examples/compliance)

* [Compliance Batch](https://developer.twitter.com/en/docs/twitter-api/compliance/batch-compliance/introduction)

## Rate Limiting
With each response, the rate limits from the response header are returned.  This allows the caller to manage any limits that are imposed.  Along with the response, errors that are returned may have rate limits as well.  If the error occurs after the request is sent, then rate limits may apply and are returned.

There is an example of rate limiting from a response [here](./_examples/misc/rate-limit/main.go).

This is an example of a twitter callout and if the limits have been reached, then it will back off and try again.
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

## Error Handling
There are different types of error handling within the library.  The library supports errors and partial errors defined by [twitter](https://developer.twitter.com/en/support/twitter-api/error-troubleshooting).

### Parameter Errors
The library does some error checking before a callout.  This checking is very basic, like making sure an id is not an empty string.  If there is an parameter error, it will be wrapped with `ErrParameter`.

```go
	opts := twitter.ListUserMembersOpts{
		MaxResults: 1,
	}
	tweetResponse, err := client.TweetLikesLookup(ctx, id, opts)
	switch {
	case errors.Is(err, twitter.ErrParameter):
		// handle a parameter error
	case err != nil:
		// handle other errors
	default:
		// happy path
	}
```

### Callout Errors
The library will return any errors from when creating and _doing_ the callout.  These errors might be, but not limited to, json encoding error or http request or client error.  These errors are also wrapped to allow for the caller to handle specific errors.

```go
	opts := twitter.ListUserMembersOpts{
		MaxResults: 1,
	}
	tweetResponse, err := client.TweetLikesLookup(ctx, id, opts)
	jErr := &json.UnsupportedValueError{}
	switch {
	case errors.As(err, &jErr):
		// handle a json error
	case err != nil:
		// handle other errors
	default:
		// happy path
	}
```

### Response Decode Errors
The library will return a json decode error, `ResponseDecodeError`, when the response is malformed.  This is done to allow for the rate limits to be part of the error.

```go
	opts := twitter.ListUserMembersOpts{
		MaxResults: 1,
	}
	tweetResponse, err := client.TweetLikesLookup(ctx, id, opts)
	
	rdErr := &twitter.ResponseDecodeError{}
	switch {
	case errors.As(err, &rdErr):
		// handle response decode error
	case err != nil:
		// handle other errors
	default:
		// happy path
	}
```

### Twitter HTTP Response Errors
The library will return a HTTP error, `HTTPError`, when a HTTP status is not successful.  This allows for the twitter error response to be decoded and the rate limits to be part of the error.

```go
	opts := twitter.ListUserMembersOpts{
		MaxResults: 1,
	}
	tweetResponse, err := client.TweetLikesLookup(ctx, id, opts)
	
	httpErr := &twitter.HTTPError{}
	switch {
	case errors.As(err, &httpErr):
		// handle http response error
	case err != nil:
		// handle other errors
	default:
		// happy path
	}
```

### Twitter Partial Errors
The library will return what twitter defines as partial errors.  These errors are not return as an error in the callout, but in the response as the callout was returned as successful.

```go
	opts := twitter.ListUserMembersOpts{
		MaxResults: 1,
	}
	tweetResponse, err := client.TweetLikesLookup(ctx, id, opts)
	if err != nil {
		/// handle error
	}	
	// handle response
	if len(tweetResponse.Raw.Errors) > 0 {
		// handle partial errors
	}
```

## Examples
Much like `v1`, there is an `_example` directory to demonstrate library usage.

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