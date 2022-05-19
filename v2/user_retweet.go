package twitter

import (
	"net/http"
	"strconv"
	"strings"
)

// RetweetData will be returned by the manage retweet APIs
type RetweetData struct {
	Retweeted bool `json:"retweeted"`
}

// UserRetweetResponse is the response with a user retweet
type UserRetweetResponse struct {
	Data      *RetweetData `json:"data"`
	RateLimit *RateLimit
}

// DeleteUserRetweetResponse is the response with a user retweet
type DeleteUserRetweetResponse struct {
	Data      *RetweetData `json:"data"`
	RateLimit *RateLimit
}

// UserRetweetLookupResponse os the response that contains the users
type UserRetweetLookupResponse struct {
	Raw       *UserRetweetRaw
	Meta      *UserRetweetMeta `json:"meta"`
	RateLimit *RateLimit
}

// UserRetweetMeta is the meta data returned by the retweet user lookup
type UserRetweetMeta struct {
	ResultCount   int    `json:"result_count"`
	NextToken     string `json:"next_token"`
	PreviousToken string `json:"previous_token"`
}

// UserRetweetRaw is the raw data and includes from the response
type UserRetweetRaw struct {
	Users    []*UserObj              `json:"data"`
	Includes *UserRetweetRawIncludes `json:"includes,omitempty"`
	Errors   []*ErrorObj             `json:"errors,omitempty"`
}

// UserRetweetRawIncludes are the includes from the options
type UserRetweetRawIncludes struct {
	Tweets []*TweetObj `json:"tweets,omitempty"`
}

// UserRetweetLookupOpts are the options for the user retweet loopup
type UserRetweetLookupOpts struct {
	Expansions      []Expansion
	TweetFields     []TweetField
	UserFields      []UserField
	MediaFields     []MediaField
	PlaceFields     []PlaceField
	PollFields      []PollField
	MaxResults      int
	PaginationToken string
}

func (u UserRetweetLookupOpts) addQuery(req *http.Request) {
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
	if u.MaxResults > 0 {
		q.Add("max_results", strconv.Itoa(u.MaxResults))
	}
	if len(u.PaginationToken) > 0 {
		q.Add("pagination_token", u.PaginationToken)
	}
	if len(q) > 0 {
		req.URL.RawQuery = q.Encode()
	}
}
