package twitter

import (
	"net/http"
	"strconv"
	"strings"
	"time"
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

// UserTimelineOpts are options in the user tweet timeline
type UserTimelineOpts struct {
	Excludes        []Exclude
	Expansions      []Expansion
	MediaFields     []MediaField
	PlaceFields     []PlaceField
	PollFields      []PollField
	TweetFields     []TweetField
	UserFields      []UserField
	SinceID         string
	UntilID         string
	PaginationToken string
	MaxResults      int
	StartTime       time.Time
	EndTime         time.Time
}

func (u UserTimelineOpts) addQuery(req *http.Request) {
	q := req.URL.Query()
	if len(u.Excludes) > 0 {
		q.Add("exclude", strings.Join(excludetringArray(u.Excludes), ","))
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
	if len(u.Expansions) > 0 {
		q.Add("expansions", strings.Join(expansionStringArray(u.Expansions), ","))
	}
	if len(u.SinceID) > 0 {
		q.Add("since_id", u.SinceID)
	}
	if len(u.UntilID) > 0 {
		q.Add("until_id", u.UntilID)
	}
	if u.MaxResults > 0 {
		q.Add("max_results", strconv.Itoa(u.MaxResults))
	}
	if len(u.PaginationToken) > 0 {
		q.Add("pagination_token", u.PaginationToken)
	}
	if !u.EndTime.IsZero() {
		q.Add("end_time", u.EndTime.Format(time.RFC3339))
	}
	if !u.StartTime.IsZero() {
		q.Add("start_time", u.StartTime.Format(time.RFC3339))
	}
	if len(q) > 0 {
		req.URL.RawQuery = q.Encode()
	}
}
