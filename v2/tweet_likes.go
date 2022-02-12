package twitter

import (
	"net/http"
	"strconv"
	"strings"
)

// TweetLikesLookupResponse is the user from the tweet likes
type TweetLikesLookupResponse struct {
	Raw       *UserRaw
	Meta      *TweetLikesMeta `json:"meta"`
	RateLimit *RateLimit
}

// TweetLikesMeta is the meta data from the response
type TweetLikesMeta struct {
	ResultCount   int    `json:"result_count"`
	NextToken     string `json:"next_token"`
	PreviousToken string `json:"previous_token"`
}

// TweetLikesLookupOpts the user like lookup options
type TweetLikesLookupOpts struct {
	Expansions      []Expansion
	TweetFields     []TweetField
	UserFields      []UserField
	MediaFields     []MediaField
	PlaceFields     []PlaceField
	PollFields      []PollField
	MaxResults      int
	PaginationToken string
}

func (u TweetLikesLookupOpts) addQuery(req *http.Request) {
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
