package twitter

import (
	"net/http"
	"strconv"
	"strings"
)

// UserMutesLookupOpts are the options for the users muted API
type UserMutesLookupOpts struct {
	Expansions      []Expansion
	TweetFields     []TweetField
	UserFields      []UserField
	MaxResults      int
	PaginationToken string
}

func (u UserMutesLookupOpts) addQuery(req *http.Request) {
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

// UserMutesLookupResponse is the list of users that are muted
type UserMutesLookupResponse struct {
	Raw       *UserRaw
	Meta      *UserMutesLookupMeta `json:"meta"`
	RateLimit *RateLimit
}

// UserMutesLookupMeta is the meta associated with the muted users lookup
type UserMutesLookupMeta struct {
	ResultCount   int    `json:"result_count"`
	NextToken     string `json:"next_token"`
	PreviousToken string `json:"previous_token"`
}

// UserMutesData indicates if the user is muted
type UserMutesData struct {
	Muting bool `json:"muting"`
}

// UserMutesResponse is when a user mutes another
type UserMutesResponse struct {
	Data      *UserMutesData `json:"data"`
	RateLimit *RateLimit
}

// UserDeleteMutesResponse is when a user unmutes another
type UserDeleteMutesResponse struct {
	Data      *UserMutesData `json:"data"`
	RateLimit *RateLimit
}
