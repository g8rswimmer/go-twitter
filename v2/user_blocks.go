package twitter

import (
	"net/http"
	"strconv"
	"strings"
)

// UserBlocksLookupOpts are the options for the users blocked API
type UserBlocksLookupOpts struct {
	Expansions      []Expansion
	TweetFields     []TweetField
	UserFields      []UserField
	MaxResults      int
	PaginationToken string
}

func (u UserBlocksLookupOpts) addQuery(req *http.Request) {
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

// UserBlocksLookupResponse is the list of users that are blocked
type UserBlocksLookupResponse struct {
	Raw       *UserRaw
	Meta      *UserBlocksLookupMeta `json:"meta"`
	RateLimit *RateLimit
}

// UserBlocksLookupMeta is the meta associated with the blocked users lookup
type UserBlocksLookupMeta struct {
	ResultCount   int    `json:"result_count"`
	NextToken     string `json:"next_token"`
	PreviousToken string `json:"previous_token"`
}

// UserBlocksData indicates if the user is blocked
type UserBlocksData struct {
	Blocking bool `json:"blocking"`
}

// UserBlocksResponse is when a user blocks another
type UserBlocksResponse struct {
	Data      *UserBlocksData `json:"data"`
	RateLimit *RateLimit
}

// UserDeleteBlocksResponse is when a user unblocks another
type UserDeleteBlocksResponse struct {
	Data      *UserBlocksData `json:"data"`
	RateLimit *RateLimit
}
