package twitter 

import (
	"net/http"
	"strings"
)

// UserLikesResponse the response for the user likes
type UserLikesResponse struct {
	Data *UserLikesData `json:"data"`
}

// DeteleUserLikesResponse the response for the user unlike
type DeteleUserLikesResponse struct {
	Data *UserLikesData `json:"data"`
}

// UserLikesData is the data from the user like management
type UserLikesData struct {
	Liked bool `json:"liked"`
}

// UserTweetLikesLookupResponse is the user from the tweet likes
type UserTweetLikesLookupResponse struct {
	Raw  *UserTweetLikesRaw
}

// UserTweetLikesRaw is the raw data from the user tweet likes
type UserTweetLikesRaw struct {
	Users    []*UserObj              `json:"data"`
	Includes *UserTweetLikesRawIncludes `json:"includes,omitempty"`
	Errors   []*ErrorObj             `json:"errors,omitempty"`
}

// UserTweetLikesRawIncludes is the objects that relate to the user
type UserTweetLikesRawIncludes struct {
	Tweets []*TweetObj `json:"tweets,omitempty"`
}

// UserTweetLikesLookupOpts the user like lookup options
type UserTweetLikesLookupOpts struct {
	Expansions  []Expansion
	TweetFields []TweetField
	UserFields  []UserField
	MediaFields []MediaField
	PlaceFields []PlaceField
	PollFields  []PollField
}

func (u UserTweetLikesLookupOpts) addQuery(req *http.Request) {
	q := req.URL.Query()
	if len(u.Expansions) > 0 {
		q.Add("expansions", strings.Join(expansionStringArray(u.Expansions), ","))
	}
	if len(u.TweetFields) > 0 {
		q.Add("tweet.fields", strings.Join(tweetFieldStringArray(u.TweetFields), ","))
	}
	if len(u.UserFields) > 0 {
		q.Add("user.fields", strings.Join(userFieldStringArray(u.UserFields), ","))
	}
	if len(u.MediaFields) > 0 {
		q.Add("media.fields", strings.Join(mediaFieldStringArray(u.MediaFields), ","))
	}
	if len(u.PlaceFields) > 0 {
		q.Add("place.fields", strings.Join(placeFieldStringArray(u.PlaceFields), ","))
	}
	if len(u.PollFields) > 0 {
		q.Add("poll.fields", strings.Join(pollFieldStringArray(u.PollFields), ","))
	}
	if len(q) > 0 {
		req.URL.RawQuery = q.Encode()
	}
}
