package twitter

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"
)

// TweetSearchSortOrder specifies the order the tweets are returned
type TweetSearchSortOrder string

const (
	// TweetSearchSortOrderRecency will return the tweets in order of recency
	TweetSearchSortOrderRecency TweetSearchSortOrder = "recency"
	// TweetSearchSortOrderRelevancy will return the tweets in order of relevancy
	TweetSearchSortOrderRelevancy TweetSearchSortOrder = "relevancy"
)

// TweetRecentSearchOpts are the optional parameters for the recent search API
type TweetRecentSearchOpts struct {
	Expansions  []Expansion
	MediaFields []MediaField
	PlaceFields []PlaceField
	PollFields  []PollField
	TweetFields []TweetField
	UserFields  []UserField
	StartTime   time.Time
	EndTime     time.Time
	SortOrder   TweetSearchSortOrder
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
	if len(t.SortOrder) > 0 {
		q.Add("sort_order", string(t.SortOrder))
	}
	if len(q) > 0 {
		req.URL.RawQuery = q.Encode()
	}
}

// TweetRecentSearchResponse contains all of the information from a tweet recent search
type TweetRecentSearchResponse struct {
	Raw       *TweetRaw
	Meta      *TweetRecentSearchMeta `json:"meta"`
	RateLimit *RateLimit
}

// TweetRecentSearchMeta contains the recent search information
type TweetRecentSearchMeta struct {
	NewestID    string `json:"newest_id"`
	OldestID    string `json:"oldest_id"`
	ResultCount int    `json:"result_count"`
	NextToken   string `json:"next_token"`
}

// TweetSearchOpts are the tweet search options
type TweetSearchOpts struct {
	Expansions  []Expansion
	MediaFields []MediaField
	PlaceFields []PlaceField
	PollFields  []PollField
	TweetFields []TweetField
	UserFields  []UserField
	StartTime   time.Time
	EndTime     time.Time
	SortOrder   TweetSearchSortOrder
	MaxResults  int
	NextToken   string
	SinceID     string
	UntilID     string
}

func (t TweetSearchOpts) addQuery(req *http.Request) {
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
	if len(t.SortOrder) > 0 {
		q.Add("sort_order", string(t.SortOrder))
	}
	if len(q) > 0 {
		req.URL.RawQuery = q.Encode()
	}
}

// TweetSearchResponse is the tweet search response
type TweetSearchResponse struct {
	Raw       *TweetRaw
	Meta      *TweetSearchMeta `json:"meta"`
	RateLimit *RateLimit
}

// TweetSearchMeta is the tweet search meta data
type TweetSearchMeta struct {
	NewestID    string `json:"newest_id"`
	OldestID    string `json:"oldest_id"`
	ResultCount int    `json:"result_count"`
	NextToken   string `json:"next_token"`
}

// TweetSearchStreamRule is the search stream filter rule
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

// TweetSearchStreamRuleID is the filter rule id
type TweetSearchStreamRuleID string

func (t TweetSearchStreamRuleID) validate() error {
	if len(t) == 0 {
		return fmt.Errorf("tweet search rule id is required %w", ErrParameter)
	}
	return nil
}

type tweetSearchStreamRuleIDs []TweetSearchStreamRuleID

func (t tweetSearchStreamRuleIDs) validate() error {
	for _, id := range t {
		if err := id.validate(); err != nil {
			return err
		}
	}
	return nil
}

func (t tweetSearchStreamRuleIDs) toStringArray() []string {
	ids := make([]string, len(t))
	for i, id := range t {
		ids[i] = string(id)
	}
	return ids
}

// TweetSearchStreamRuleEntity is the search filter rule entity
type TweetSearchStreamRuleEntity struct {
	ID TweetSearchStreamRuleID `json:"id"`
	TweetSearchStreamRule
}

// TweetSearchStreamRulesResponse is the response to getting the search rules
type TweetSearchStreamRulesResponse struct {
	Rules     []*TweetSearchStreamRuleEntity `json:"data"`
	Meta      *TweetSearchStreamRuleMeta     `json:"meta"`
	Errors    []*ErrorObj                    `json:"errors,omitempty"`
	RateLimit *RateLimit
}

// TweetSearchStreamAddRuleResponse is the response from adding rules
type TweetSearchStreamAddRuleResponse struct {
	Rules     []*TweetSearchStreamRuleEntity `json:"data"`
	Meta      *TweetSearchStreamRuleMeta     `json:"meta"`
	Errors    []*ErrorObj                    `json:"errors,omitempty"`
	RateLimit *RateLimit
}

// TweetSearchStreamDeleteRuleResponse is the response from deleting rules
type TweetSearchStreamDeleteRuleResponse struct {
	Meta      *TweetSearchStreamRuleMeta `json:"meta"`
	Errors    []*ErrorObj                `json:"errors,omitempty"`
	RateLimit *RateLimit
}

// TweetSearchStreamRuleMeta is the meta data object from the request
type TweetSearchStreamRuleMeta struct {
	Sent    time.Time                    `json:"sent"`
	Summary TweetSearchStreamRuleSummary `json:"summary"`
}

// TweetSearchStreamRuleSummary is the summary of the search filters
type TweetSearchStreamRuleSummary struct {
	Created    int `json:"created"`
	NotCreated int `json:"not_created"`
	Deleted    int `json:"deleted"`
	NotDeleted int `json:"not_deleted"`
}
