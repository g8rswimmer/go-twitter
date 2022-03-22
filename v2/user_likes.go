package twitter

import (
	"net/http"
	"strconv"
	"strings"
)

// UserLikesResponse the response for the user likes
type UserLikesResponse struct {
	Data      *UserLikesData `json:"data"`
	RateLimit *RateLimit
}

// DeleteUserLikesResponse the response for the user unlike
type DeleteUserLikesResponse struct {
	Data      *UserLikesData `json:"data"`
	RateLimit *RateLimit
}

// UserLikesData is the data from the user like management
type UserLikesData struct {
	Liked bool `json:"liked"`
}

// UserLikesLookupResponse is the tweets from the user likes
type UserLikesLookupResponse struct {
	Raw       *TweetRaw
	Meta      *UserLikesMeta `json:"meta"`
	RateLimit *RateLimit
}

// UserLikesMeta is the meta data from the response
type UserLikesMeta struct {
	ResultCount   int    `json:"result_count"`
	NextToken     string `json:"next_token"`
	PreviousToken string `json:"previous_token"`
}

// UserLikesLookupOpts the tweet like lookup options
type UserLikesLookupOpts struct {
	Expansions      []Expansion
	MediaFields     []MediaField
	PlaceFields     []PlaceField
	PollFields      []PollField
	TweetFields     []TweetField
	UserFields      []UserField
	MaxResults      int
	PaginationToken string
}

func (t UserLikesLookupOpts) addQuery(req *http.Request) {
	q := req.URL.Query()
	if len(t.Expansions) > 0 {
		q.Add("expansions", strings.Join(expansionStringArray(t.Expansions), ","))
	}
	if len(t.MediaFields) > 0 {
		q.Add("media.fields", strings.Join(mediaFieldStringArray(t.MediaFields), ","))
	}
	if len(t.PlaceFields) > 0 {
		q.Add("place.fields", strings.Join(placeFieldStringArray(t.PlaceFields), ","))
	}
	if len(t.PollFields) > 0 {
		q.Add("poll.fields", strings.Join(pollFieldStringArray(t.PollFields), ","))
	}
	if len(t.TweetFields) > 0 {
		q.Add("tweet.fields", strings.Join(tweetFieldStringArray(t.TweetFields), ","))
	}
	if len(t.UserFields) > 0 {
		q.Add("user.fields", strings.Join(userFieldStringArray(t.UserFields), ","))
	}
	if t.MaxResults > 0 {
		q.Add("max_results", strconv.Itoa(t.MaxResults))
	}
	if len(t.PaginationToken) > 0 {
		q.Add("pagination_token", t.PaginationToken)
	}
	if len(q) > 0 {
		req.URL.RawQuery = q.Encode()
	}
}
