package twitter

import (
	"net/http"
	"strconv"
	"strings"
)

// UserFieldOptions are the tweet options for the response
type UserFieldOptions struct {
	Expansions  []Expansion
	TweetFields []TweetField
	UserFields  []UserField
}

func (u UserFieldOptions) addQuery(req *http.Request) {
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
	if len(q) > 0 {
		req.URL.RawQuery = q.Encode()
	}
}

// UserFollowOptions are the options that can be passed for the following APIs
type UserFollowOptions struct {
	Expansions      []Expansion
	TweetFields     []TweetField
	UserFields      []UserField
	MaxResults      int
	PaginationToken string
}

func (u UserFollowOptions) addQuery(req *http.Request) {
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
