package twitter

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"
)

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
	if !t.StartTime.IsZero() {
		q.Add("start_time", t.StartTime.Format(time.RFC3339))
	}
	if !t.EndTime.IsZero() {
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

// TweetRecentSearchResponse contains all of the information from a tweet recent search
type TweetRecentSearchResponse struct {
	Raw  *TweetRaw
	Meta *TweetRecentSearchMeta `json:"meta"`
}

// TweetRecentSearchMeta contains the recent search information
type TweetRecentSearchMeta struct {
	NewestID    string `json:"newest_id"`
	OldestID    string `json:"oldest_id"`
	ResultCount int    `json:"result_count"`
	NextToken   string `json:"next_token"`
}

type TweetSearchStreamRule struct {
	Value string `json:"value"`
	Tag   string `json:"tag,omitempty"`
}

func (t TweetSearchStreamRule) validate() error {
	if len(t.Value) == 0 {
		return fmt.Errorf("tweet search stream rule value is required: %w", ErrParameter)
	}
	return nil
}

type tweetSearchStreamRules []TweetSearchStreamRule

func (t tweetSearchStreamRules) validate() error {
	for _, rule := range t {
		if err := rule.validate(); err != nil {
			return err
		}
	}
	return nil
}

type TweetSearchStreamRuleID string

func (t TweetSearchStreamRuleID) validate() error {
	if len(t) == 0 {
		return fmt.Errorf("tweet search rule id is required %w", ErrParameter)
	}
	return nil
}

type tweetSeachStreamRuleIDs []TweetSearchStreamRuleID

func (t tweetSeachStreamRuleIDs) validate() error {
	for _, id := range t {
		if err := id.validate(); err != nil {
			return err
		}
	}
	return nil
}

type TweetSearchStreamRuleEntity struct {
	ID TweetSearchStreamRuleID `json:"id"`
	TweetSearchStreamRule
}

type TweetSearchStreamAddRuleResponse struct {
	Rules  []*TweetSearchStreamRuleEntity `json:"data"`
	Meta   *TweetSearchStreamRuleMeta     `json:"meta"`
	Errors []*ErrorObj                    `json:"errors,omitempty"`
}

type TweetSearchStreamDeleteRuleResponse struct {
	Meta   *TweetSearchStreamRuleMeta `json:"meta"`
	Errors []*ErrorObj                `json:"errors,omitempty"`
}

type TweetSearchStreamRuleMeta struct {
	Sent    time.Time                    `json:"sent"`
	Summary TweetSearchStreamRuleSummary `json:"summary"`
}

type TweetSearchStreamRuleSummary struct {
	Created    int `json:"created"`
	NotCreated int `json:"not_created"`
	Deleted    int `json:"deleted"`
	NotDeleted int `json:"not_deleted"`
}
