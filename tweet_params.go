package twitter

import (
	"fmt"
	"net/http"
	"strings"
	"time"
)

// TweetFieldOptions are the tweet options for the response
type TweetFieldOptions struct {
	Expansions  []Expansion
	MediaFields []MediaField
	PlaceFields []PlaceField
	PollFields  []PollField
	TweetFields []TweetField
	UserFields  []UserField
}

func (t TweetFieldOptions) addQuery(req *http.Request) {
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

// TweetRecentSearchOptions handles all of the recent search parameters
type TweetRecentSearchOptions struct {
	query     string
	StartTime time.Time
	EndTime   time.Time
	MaxResult int
	NextToken string
	SinceID   string
	UntilID   string
}

func (t TweetRecentSearchOptions) addQuery(req *http.Request) {
	q := req.URL.Query()
	q.Add("query", t.query)

	if !t.StartTime.IsZero() {
		q.Add("start_time", t.StartTime.Format(time.RFC3339))
	}
	if !t.EndTime.IsZero() {
		q.Add("end_time", t.EndTime.Format(time.RFC3339))
	}
	if t.MaxResult >= 10 {
		q.Add("max_results", fmt.Sprintf("%d", t.MaxResult))
	}
	if len(t.NextToken) > 0 {
		q.Add("next_token", t.NextToken)
	}
	if len(t.SinceID) > 0 {
		q.Add("since_id", t.SinceID)
	}
	if len(t.UntilID) > 0 {
		q.Add("until_id", t.UntilID)
	}
	if len(q) > 0 {
		req.URL.RawQuery = q.Encode()
	}
}
