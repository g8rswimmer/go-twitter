package twitter

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type CreateTweetOps struct {
	DirectMessageDeepLink string           `json:"direct_message_deep_link,omitempty"`
	ForSuperFollowersOnly bool             `json:"for_super_followers_only,omitempty"`
	QuoteTweetID          string           `json:"quote_tweet_id,omitempty"`
	Text                  string           `json:"text,omitemtpy"`
	ReplySettings         string           `json:"reply_settings"`
	Geo                   CreateTweetGeo   `json:"geo"`
	Media                 CreateTweetMedia `json:"media"`
	Poll                  CreateTweetPoll  `json:"poll"`
	Reply                 CreateTweetReply `json:"reply"`
}

func (t CreateTweetOps) validate() error {
	if err := t.Media.validate(); err != nil {
		return fmt.Errorf("create tweet error: %w", err)
	}
	if err := t.Poll.validate(); err != nil {
		return fmt.Errorf("create tweet error: %w", err)
	}
	if err := t.Reply.validate(); err != nil {
		return fmt.Errorf("create tweet error: %w", err)
	}
	if len(t.Media.IDs) == 0 && len(t.Text) == 0 {
		return fmt.Errorf("create tweet text is required if no media ids %w", ErrParameter)
	}
	return nil
}

type CreateTweetGeo struct {
	PlaceID string `json:"place_id"`
}

type CreateTweetMedia struct {
	IDs           []string `json:"media_ids"`
	TaggedUserIDs []string `json:"tagged_user_ids"`
}

func (m CreateTweetMedia) validate() error {
	if len(m.TaggedUserIDs) > 0 && len(m.IDs) == 0 {
		return fmt.Errorf("media ids are required if taged user ids are present %w", ErrParameter)
	}
	return nil
}

type CreateTweetPoll struct {
	DurationMinutes int      `json:"duration_minutes"`
	Options         []string `json:"options"`
}

func (p CreateTweetPoll) validate() error {
	if len(p.Options) > 0 && p.DurationMinutes <= 0 {
		return fmt.Errorf("poll duration minutes are required with options %w", ErrParameter)
	}
	return nil
}

type CreateTweetReply struct {
	ExcludeReplyUserIDs []string `json:"exclude_reply_user_ids"`
	InReplyToTweetID    string   `json:"in_reply_to_tweet_id"`
}

func (r CreateTweetReply) validate() error {
	if len(r.ExcludeReplyUserIDs) > 0 && len(r.InReplyToTweetID) == 0 {
		return fmt.Errorf("reply in reply to tweet is needs to be present it exclude reply user ids are present %w", ErrParameter)
	}
	return nil
}

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

// TweetRecentSearchOpts are the optional parameters for the recent seach API
type TweetRecentSearchOpts struct {
	Expansions  []Expansion
	MediaFields []MediaField
	PlaceFields []PlaceField
	PollFields  []PollField
	TweetFields []TweetField
	UserFields  []UserField
	StartTime   time.Time
	EndTime     time.Time
	MaxResults  int
	NextToken   string
	SinceID     string
	UntilID     string
}

func (t TweetRecentSearchOpts) addQuery(req *http.Request) {
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
	if t.StartTime.IsZero() == false {
		q.Add("start_time", t.StartTime.Format(time.RFC3339))
	}
	if t.EndTime.IsZero() == false {
		q.Add("end_time", t.EndTime.Format(time.RFC3339))
	}
	if t.MaxResults > 0 {
		q.Add("max_results", strconv.Itoa(t.MaxResults))
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

// UserTweetTimelineOpts are the options for the user tweet timeline request
type UserTweetTimelineOpts struct {
	Expansions      []Expansion
	MediaFields     []MediaField
	PlaceFields     []PlaceField
	PollFields      []PollField
	TweetFields     []TweetField
	UserFields      []UserField
	Excludes        []Exclude
	StartTime       time.Time
	EndTime         time.Time
	MaxResults      int
	PaginationToken string
	SinceID         string
	UntilID         string
}

func (t UserTweetTimelineOpts) addQuery(req *http.Request) {
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
	if len(t.Excludes) > 0 {
		q.Add("exclude", strings.Join(excludeStringArray(t.Excludes), ","))
	}
	if t.StartTime.IsZero() == false {
		q.Add("start_time", t.StartTime.Format(time.RFC3339))
	}
	if t.EndTime.IsZero() == false {
		q.Add("end_time", t.EndTime.Format(time.RFC3339))
	}
	if t.MaxResults > 0 {
		q.Add("max_results", strconv.Itoa(t.MaxResults))
	}
	if len(t.PaginationToken) > 0 {
		q.Add("pagination_token", t.PaginationToken)
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

// UserMentionTimelineOpts are the options for the user mention timeline request
type UserMentionTimelineOpts struct {
	Expansions      []Expansion
	MediaFields     []MediaField
	PlaceFields     []PlaceField
	PollFields      []PollField
	TweetFields     []TweetField
	UserFields      []UserField
	StartTime       time.Time
	EndTime         time.Time
	MaxResults      int
	PaginationToken string
	SinceID         string
	UntilID         string
}

func (t UserMentionTimelineOpts) addQuery(req *http.Request) {
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
	if t.StartTime.IsZero() == false {
		q.Add("start_time", t.StartTime.Format(time.RFC3339))
	}
	if t.EndTime.IsZero() == false {
		q.Add("end_time", t.EndTime.Format(time.RFC3339))
	}
	if t.MaxResults > 0 {
		q.Add("max_results", strconv.Itoa(t.MaxResults))
	}
	if len(t.PaginationToken) > 0 {
		q.Add("pagination_token", t.PaginationToken)
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
