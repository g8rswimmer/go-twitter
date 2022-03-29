package twitter

import (
	"net/http"
	"strconv"
	"strings"
)

// TweetBookmarksLookupOpts are the tweet bookmark lookup options
type TweetBookmarksLookupOpts struct {
	MaxResults      int
	PaginationToken string
	Expansions      []Expansion
	MediaFields     []MediaField
	PlaceFields     []PlaceField
	PollFields      []PollField
	TweetFields     []TweetField
	UserFields      []UserField
}

func (t TweetBookmarksLookupOpts) addQuery(req *http.Request) {
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

// TweetBookmarksLookupResponse is the response to the bookmark lookup
type TweetBookmarksLookupResponse struct {
	Raw       *TweetRaw
	Meta      *TweetBookmarksLookupMeta
	RateLimit *RateLimit
}

// TweetBookmarksLookupMeta is the meta for the bookmark lookup
type TweetBookmarksLookupMeta struct {
	ResultCount int    `json:"result_count"`
	NextToken   string `json:"next_token"`
}

// AddTweetBookmarkResponse is the response for adding a bookmark
type AddTweetBookmarkResponse struct {
	Tweet     *TweetBookmarkData `json:"data"`
	RateLimit *RateLimit
}

// RemoveTweetBookmarkResponse is the response for removing a bookmark
type RemoveTweetBookmarkResponse struct {
	Tweet     *TweetBookmarkData `json:"data"`
	RateLimit *RateLimit
}

// TweetBookmarkData is the data from adding or removing a bookmark
type TweetBookmarkData struct {
	Bookmarked bool `json:"bookmarked"`
}
