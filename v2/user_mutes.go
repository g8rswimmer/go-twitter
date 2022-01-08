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

type UserMutesLookupResponse struct {
	Raw  *UserRaw
	Meta *UserMutesLookupMeta `json:"meta"`
}

type UserMutesLookupMeta struct {
	ResultCount   int    `json:"result_count"`
	NextToken     string `json:"next_token"`
	PreviousToken string `json:"previous_token"`
}

type UserMutesData struct {
	Muting bool `json:"muting"`
}

type UserMutesResponse struct {
	Data *UserMutesData `json:"data"`
}

type UserDeleteMutesResponse struct {
	Data *UserMutesData `json:"data"`
}
