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
	tweetOpts := twitter.UserTweetOpts{
		TweetFields: []twitter.TweetField{
			twitter.TweetFieldAttachments,
			twitter.TweetFieldAuthorID,
			twitter.TweetFieldContextAnnotations,
			twitter.TweetFieldConversationID,
			twitter.TweetFieldCreatedAt,
			twitter.TweetFieldEntities,
			twitter.TweetFieldGeo,
			twitter.TweetFieldID,
			twitter.TweetFieldInReplyToUserID,
			twitter.TweetFieldLanguage,
			twitter.TweetFieldPossiblySensitve,
			twitter.TweetFieldPublicMetrics,
			twitter.TweetFieldReferencedTweets,
			twitter.TweetFieldSource,
			twitter.TweetFieldText,
		},
		UserFields: []twitter.UserField{
			twitter.UserFieldCreatedAt,
			twitter.UserFieldDescription,
			twitter.UserFieldEntities,
			twitter.UserFieldLocation,
			twitter.UserFieldName,
			twitter.UserFieldPinnedTweetID,
			twitter.UserFieldProfileImageURL,
			twitter.UserFieldProtected,
			twitter.UserFieldURL,
			twitter.UserFieldUserName,
			twitter.UserFieldVerified,
			twitter.UserFieldWithHeld,
		},
		Expansions: []twitter.Expansion{
			twitter.ExpansionAuthorID,
			twitter.ExpansionReferencedTweetsID,
			twitter.ExpansionReferencedTweetsIDAuthorID,
			twitter.ExpansionEntitiesMentionsUserName,
			twitter.ExpansionAttachmentsMediaKeys,
			twitter.ExpansionInReplyToUserID,
			twitter.ExpansionGeoPlaceID,
		},
		PlaceFields: []twitter.PlaceField{
			twitter.PlaceFieldContainedWithin,
			twitter.PlaceFieldCountry,
			twitter.PlaceFieldCountryCode,
			twitter.PlaceFieldFullName,
			twitter.PlaceFieldGeo,
			twitter.PlaceFieldID,
			twitter.PlaceFieldName,
			twitter.PlaceFieldPlaceType,
		},
		PollFields: []twitter.PollField{
			twitter.PollFieldDurationMinutes,
			twitter.PollFieldEndDateTime,
			twitter.PollFieldID,
			twitter.PollFieldOptions,
			twitter.PollFieldVotingStatus,
		},
		MaxResults: 10,
	}

	userTweets, err := user.Tweets(context.Background(), *id, tweetOpts)
	var tweetErr *twitter.TweetErrorResponse
	switch {
	case errors.As(err, &tweetErr):
		printTweetError(tweetErr)
	case err != nil:
		fmt.Println(err)
	default:
		printUserTweets(userTweets)
	}

}
func printUserTweets(userTweets *twitter.UserTweets) {
	enc, err := json.MarshalIndent(userTweets, "", "    ")
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
