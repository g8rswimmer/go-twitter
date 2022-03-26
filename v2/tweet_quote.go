package twitter

import (
	"net/http"
	"strconv"
	"strings"
)

// QuoteTweetsLookupOpts are the options for the quote tweets
type QuoteTweetsLookupOpts struct {
	MaxResults      int
	PaginationToken string
	Expansions      []Expansion
	MediaFields     []MediaField
	PlaceFields     []PlaceField
	PollFields      []PollField
	TweetFields     []TweetField
	UserFields      []UserField
}

func (qt QuoteTweetsLookupOpts) addQuery(req *http.Request) {
	q := req.URL.Query()
	if len(qt.Expansions) > 0 {
		q.Add("expansions", strings.Join(expansionStringArray(qt.Expansions), ","))
	}
	if len(qt.MediaFields) > 0 {
		q.Add("media.fields", strings.Join(mediaFieldStringArray(qt.MediaFields), ","))
	}
	if len(qt.PlaceFields) > 0 {
		q.Add("place.fields", strings.Join(placeFieldStringArray(qt.PlaceFields), ","))
	}
	if len(qt.PollFields) > 0 {
		q.Add("poll.fields", strings.Join(pollFieldStringArray(qt.PollFields), ","))
	}
	if len(qt.TweetFields) > 0 {
		q.Add("tweet.fields", strings.Join(tweetFieldStringArray(qt.TweetFields), ","))
	}
	if len(qt.UserFields) > 0 {
		q.Add("user.fields", strings.Join(userFieldStringArray(qt.UserFields), ","))
	}
	if qt.MaxResults > 0 {
		q.Add("max_results", strconv.Itoa(qt.MaxResults))
	}
	if len(qt.PaginationToken) > 0 {
		q.Add("pagination_token", qt.PaginationToken)
	}
	if len(q) > 0 {
		req.URL.RawQuery = q.Encode()
	}
}

// QuoteTweetsLookupResponse is the response from the quote tweet
type QuoteTweetsLookupResponse struct {
	Raw       *TweetRaw
	Meta      *QuoteTweetsLookupMeta
	RateLimit *RateLimit
}

// QuoteTweetsLookupMeta is the meta data from the response
type QuoteTweetsLookupMeta struct {
	ResultCount int    `json:"result_count"`
	NextToken   string `json:"next_token"`
}
