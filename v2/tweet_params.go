package twitter

import (
	"net/http"
	"strings"
)

// TweetLookupOpts are the optional paramters that can be passed to the lookup callout
type TweetLookupOpts struct {
	Expansions  []Expansion
	MediaFields []MediaField
	PlaceFields []PlaceField
	PollFields  []PollField
	TweetFields []TweetField
	UserFields  []UserField
}

func (t TweetLookupOpts) addQuery(req *http.Request) {
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
	if len(q) > 0 {
		req.URL.RawQuery = q.Encode()
	}
}
