package twitter

import (
	"net/http"
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
