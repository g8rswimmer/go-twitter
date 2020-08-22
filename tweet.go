package twitter

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"
)

const (
	tweetLookupEndpoint            = "2/tweets"
	tweetRecentSearchEndpoint      = "2/tweets/search/recent"
	tweetSearchStreamRulesEndpoint = "/2/tweets/search/stream/rules"
	tweetSearchStreamEndpoint      = "/2/tweets/search/stream"
	tweetMaxIDs                    = 100
	tweetQuerySize                 = 512
)

// TweetLookups is a map of tweet lookups
type TweetLookups map[string]TweetLookup

func (t TweetLookups) lookup(decoder *json.Decoder) error {
	type include struct {
		Media []*MediaObj `json:"medias"`
		Place []*PlaceObj `json:"places"`
		Poll  []*PollObj  `json:"polls"`
		User  []*UserObj  `json:"users"`
	}
	type body struct {
		Data    TweetObj `json:"data"`
		Include include  `json:"includes"`
	}
	b := &body{}
	if err := decoder.Decode(b); err != nil {
		return fmt.Errorf("tweet lookup decode error %w", err)
	}

	tl := TweetLookup{
		Tweet: b.Data,
	}
	if len(b.Include.Media) > 0 {
		tl.Media = b.Include.Media[0]
	}
	if len(b.Include.Place) > 0 {
		tl.Place = b.Include.Place[0]
	}
	if len(b.Include.Poll) > 0 {
		tl.Poll = b.Include.Poll[0]
	}
	if len(b.Include.User) > 0 {
		tl.User = b.Include.User[0]
	}
	t[b.Data.ID] = tl

	return nil
}

func (t TweetLookups) lookups(decoder *json.Decoder) error {
	type include struct {
		Media []*MediaObj `json:"medias"`
		Place []*PlaceObj `json:"places"`
		Poll  []*PollObj  `json:"polls"`
		User  []*UserObj  `json:"users"`
	}
	type body struct {
		Data    []TweetObj `json:"data"`
		Include include    `json:"includes"`
	}
	b := &body{}
	if err := decoder.Decode(b); err != nil {
		return fmt.Errorf("tweet lookup decode error %w", err)
	}

	for i, tweet := range b.Data {
		tl := TweetLookup{
			Tweet: tweet,
		}
		if i < len(b.Include.Media) {
			tl.Media = b.Include.Media[i]
		}
		if i < len(b.Include.Place) {
			tl.Place = b.Include.Place[i]
		}
		if i < len(b.Include.Poll) {
			tl.Poll = b.Include.Poll[i]
		}
		if i < len(b.Include.User) {
			tl.User = b.Include.User[i]
		}
		t[tweet.ID] = tl
	}
	return nil
}

// TweetLookup is a complete tweet objects
type TweetLookup struct {
	Tweet TweetObj
	Media *MediaObj
	Place *PlaceObj
	Poll  *PollObj
	User  *UserObj
}

// TweetRecentSearchMeta is the media data returned from the recent search
type TweetRecentSearchMeta struct {
	NewestID    string `json:"newest_id"`
	OldestID    string `json:"oldest_id"`
	ResultCount int    `json:"result_count"`
	NextToken   string `json:"next_token"`
}

// TweetRecentSearch is what is returned from the callout
type TweetRecentSearch struct {
	LookUps TweetLookups
	Meta    TweetRecentSearchMeta
}

// UnmarshalJSON will unmarshal the byte array
func (t *TweetRecentSearch) UnmarshalJSON(b []byte) error {
	type meta struct {
		Meta TweetRecentSearchMeta `json:"meta"`
	}
	m := &meta{}
	if err := json.Unmarshal(b, m); err != nil {
		return err
	}
	t.Meta = m.Meta

	tl := TweetLookups{}
	tl.lookups(json.NewDecoder(bytes.NewReader(b)))
	t.LookUps = tl
	return nil
}

// TweetError is the group of errors in a response
type TweetError struct {
	Parameters interface{} `json:"parameters"`
	Message    string      `json:"message"`
}

// TweetErrorResponse is the error message from the callout
type TweetErrorResponse struct {
	StatusCode int
	Errors     []TweetError `json:"errors"`
	Title      string       `json:"title"`
	Detail     string       `json:"detail"`
	Type       string       `json:"type"`
}

func (t *TweetErrorResponse) Error() string {
	return fmt.Sprintf("status %d %s:%s", t.StatusCode, t.Title, t.Detail)
}

// TweetLookupParameters are the query parameters for the tweet lookup
type TweetLookupParameters struct {
	ids         []string
	Expansions  []Expansion
	MediaFields []MediaField
	PlaceFields []PlaceField
	PollFields  []PollField
	TweetFields []TweetField
	UserFields  []UserField
}

func (t TweetLookupParameters) encode(req *http.Request) {
	q := req.URL.Query()
	if len(t.ids) > 0 {
		q.Add("ids", strings.Join(t.ids, ","))
	}
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

// TweetStreamSearchParameters are the search tweet get parameters
type TweetStreamSearchParameters struct {
	Expansions  []Expansion
	MediaFields []MediaField
	PlaceFields []PlaceField
	PollFields  []PollField
	TweetFields []TweetField
	UserFields  []UserField
}

func (t TweetStreamSearchParameters) encode(req *http.Request) {
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

// TweetRecentSearchParameters handles all of the recent search parameters
type TweetRecentSearchParameters struct {
	query       string
	StartTime   time.Time
	EndTime     time.Time
	MaxResult   int
	NextToken   string
	SinceID     string
	UntilID     string
	Expansions  []Expansion
	MediaFields []MediaField
	PlaceFields []PlaceField
	PollFields  []PollField
	TweetFields []TweetField
	UserFields  []UserField
}

func (t TweetRecentSearchParameters) encode(req *http.Request) {
	q := req.URL.Query()
	q.Add("query", t.query)

	if t.StartTime.IsZero() == false {
		q.Add("start_time", t.StartTime.Format(time.RFC3339))
	}
	if t.EndTime.IsZero() == false {
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

// TweetSearchStreamAddRule are the rules to add the search stream
type TweetSearchStreamAddRule struct {
	Value string `json:"value"`
	Tag   string `json:"tag,omitempty"`
}

// TweetSearchStreamDeleteRule lists the search rule ids to remove
type TweetSearchStreamDeleteRule struct {
	IDs []string `json:"ids"`
}

// TweetSearchStreamRule are the rules to add and/or delete
type TweetSearchStreamRule struct {
	Add    []*TweetSearchStreamAddRule  `json:"add,omitempty"`
	Delete *TweetSearchStreamDeleteRule `json:"delete,omitempty"`
}

// TweetSearchStreamRuleData are the rules that where added
type TweetSearchStreamRuleData struct {
	ID    string `json:"id"`
	Value string `json:"value"`
	Tag   string `json:"tag"`
}

// TweetSearchStreamRuleMeta is the meta data for the search rules
type TweetSearchStreamRuleMeta struct {
	Sent    string                       `json:"sent"`
	Summary TweetSearchStreamRuleSummary `json:"summary"`
}

// TweetSearchStreamRuleSummary is the summary of the rules
type TweetSearchStreamRuleSummary struct {
	Created    int `json:"created"`
	NotCreated int `json:"not_created"`
	Deleted    int `json:"deleted"`
	NotDeleted int `json:"not_deleted"`
}

// TweetSearchStreamRules is the returned set of rules
type TweetSearchStreamRules struct {
	Data []TweetSearchStreamRuleData `json:"data"`
	Meta TweetSearchStreamRuleMeta   `json:"meta"`
}

// Tweet represents the Tweet v2 APIs
type Tweet struct {
	Authorizer Authorizer
	Client     *http.Client
	Host       string
}

// Lookup will return a tweet or tweets from a set of ids
func (t *Tweet) Lookup(ctx context.Context, ids []string, parameters TweetLookupParameters) (TweetLookups, error) {
	ep := tweetLookupEndpoint
	switch {
	case len(ids) == 0:
		return nil, fmt.Errorf("tweet lookup an id is required")
	case len(ids) > tweetMaxIDs:
		return nil, fmt.Errorf("tweet lookup: ids %d is greater than max %d", len(ids), tweetMaxIDs)
	case len(ids) == 1:
		ep += fmt.Sprintf("/%s", ids[0])
	default:
		parameters.ids = ids
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, fmt.Sprintf("%s/%s", t.Host, ep), nil)
	if err != nil {
		return nil, fmt.Errorf("tweet lookup request: %w", err)
	}
	req.Header.Add("Accept", "application/json")
	t.Authorizer.Add(req)
	parameters.encode(req)

	resp, err := t.Client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("tweet lookup response: %w", err)
	}
	defer resp.Body.Close()

	decoder := json.NewDecoder(resp.Body)

	if resp.StatusCode != http.StatusOK {
		e := &TweetErrorResponse{}
		if err := decoder.Decode(e); err != nil {
			return nil, fmt.Errorf("tweet lookup response error decode: %w", err)
		}
		e.StatusCode = resp.StatusCode
		return nil, e
	}

	tl := TweetLookups{}
	if len(ids) == 1 {
		if err := tl.lookup(decoder); err != nil {
			return nil, err
		}
		return tl, nil
	}

	if err := tl.lookups(decoder); err != nil {
		return nil, err
	}
	return tl, nil
}

// RecentSearch will query the recent search
func (t *Tweet) RecentSearch(ctx context.Context, query string, parameters TweetRecentSearchParameters) (*TweetRecentSearch, error) {
	switch {
	case len(query) == 0:
		return nil, fmt.Errorf("tweet recent search query must be present")
	case len(query) > tweetQuerySize:
		return nil, fmt.Errorf("tweet recent search query size %d greater than max %d", len(query), tweetQuerySize)
	case parameters.MaxResult > 0 && (parameters.MaxResult < 10 || parameters.MaxResult > 100):
		return nil, fmt.Errorf("tweet resent search max result needs to be between 10 -100 (%d", parameters.MaxResult)
	default:
		parameters.query = query
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, fmt.Sprintf("%s/%s", t.Host, tweetRecentSearchEndpoint), nil)
	if err != nil {
		return nil, fmt.Errorf("tweet recent search request: %w", err)
	}
	req.Header.Add("Accept", "application/json")
	t.Authorizer.Add(req)
	parameters.encode(req)

	resp, err := t.Client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("tweet recent search response: %w", err)
	}
	defer resp.Body.Close()

	decoder := json.NewDecoder(resp.Body)

	if resp.StatusCode != http.StatusOK {
		e := &TweetErrorResponse{}
		if err := decoder.Decode(e); err != nil {
			return nil, fmt.Errorf("tweet recent search response error decode: %w", err)
		}
		e.StatusCode = resp.StatusCode
		return nil, e
	}

	tr := &TweetRecentSearch{}
	if err := decoder.Decode(tr); err != nil {
		return nil, fmt.Errorf("tweet recent search response decode: %w", err)
	}
	return tr, nil
}

// ApplyFilteredStreamRules will add and/or remove rules from the seach stream
func (t *Tweet) ApplyFilteredStreamRules(ctx context.Context, rules TweetSearchStreamRule, validate bool) (*TweetSearchStreamRules, error) {
	if len(rules.Add) == 0 && rules.Delete == nil {
		return nil, errors.New("tweet search stream rules: there must be add or delete rules")
	}
	for _, add := range rules.Add {
		if len(add.Value) == 0 {
			return nil, errors.New("tweet search stream rules: add value is required")
		}
	}
	if rules.Delete != nil && len(rules.Delete.IDs) == 0 {
		return nil, errors.New("tweet search stream rules: delete ids are required")
	}
	enc, err := json.Marshal(rules)
	if err != nil {
		return nil, fmt.Errorf("tweet search stream rules: rules encoding %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, fmt.Sprintf("%s/%s", t.Host, tweetSearchStreamRulesEndpoint), bytes.NewReader(enc))
	if err != nil {
		return nil, fmt.Errorf("tweet rsearch stream rules request: %w", err)
	}

	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-type", "application/json")
	t.Authorizer.Add(req)
	if validate {
		q := req.URL.Query()
		q.Add("dry_run", "true")
		req.URL.RawQuery = q.Encode()
	}

	resp, err := t.Client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("tweet search stream rules response: %w", err)
	}
	defer resp.Body.Close()

	decoder := json.NewDecoder(resp.Body)

	if resp.StatusCode != http.StatusOK {
		e := &TweetErrorResponse{}
		if err := decoder.Decode(e); err != nil {
			return nil, fmt.Errorf("tweet search stream rules response error decode: %w", err)
		}
		e.StatusCode = resp.StatusCode
		return nil, e
	}

	tr := &TweetSearchStreamRules{}
	if err := decoder.Decode(tr); err != nil {
		return nil, fmt.Errorf("tweet search stream rules response decode: %w", err)
	}
	return tr, nil
}

// FilteredStreamRules will return the rules from the ids
func (t *Tweet) FilteredStreamRules(ctx context.Context, ids []string) (*TweetSearchStreamRules, error) {
	if len(ids) == 0 {
		return nil, errors.New("tweet search stream rules: there must be ids")
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, fmt.Sprintf("%s/%s", t.Host, tweetSearchStreamRulesEndpoint), nil)
	if err != nil {
		return nil, fmt.Errorf("tweet rsearch stream rules request: %w", err)
	}

	q := req.URL.Query()
	q.Add("ids", strings.Join(ids, ","))
	req.URL.RawQuery = q.Encode()

	req.Header.Add("Accept", "application/json")
	t.Authorizer.Add(req)

	resp, err := t.Client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("tweet search stream rules response: %w", err)
	}
	defer resp.Body.Close()

	decoder := json.NewDecoder(resp.Body)

	if resp.StatusCode != http.StatusOK {
		e := &TweetErrorResponse{}
		if err := decoder.Decode(e); err != nil {
			return nil, fmt.Errorf("tweet search stream rules response error decode: %w", err)
		}
		e.StatusCode = resp.StatusCode
		return nil, e
	}

	tr := &TweetSearchStreamRules{}
	if err := decoder.Decode(tr); err != nil {
		return nil, fmt.Errorf("tweet search stream rules response decode: %w", err)
	}
	return tr, nil
}

// FilteredStream allows to stream some tweets on a specific set of filter rules
func (t *Tweet) FilteredStream(ctx context.Context, parameters TweetStreamSearchParameters) (TweetLookups, error) {

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, fmt.Sprintf("%s/%s", t.Host, tweetSearchStreamEndpoint), nil)
	if err != nil {
		return nil, fmt.Errorf("tweet lookup request: %w", err)
	}
	req.Header.Add("Accept", "application/json")
	t.Authorizer.Add(req)
	parameters.encode(req)

	resp, err := t.Client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("tweet lookup response: %w", err)
	}
	defer resp.Body.Close()

	decoder := json.NewDecoder(resp.Body)

	if resp.StatusCode != http.StatusOK {
		e := &TweetErrorResponse{}
		if err := decoder.Decode(e); err != nil {
			return nil, fmt.Errorf("tweet lookup response error decode: %w", err)
		}
		e.StatusCode = resp.StatusCode
		return nil, e
	}

	tl := TweetLookups{}
	if err := tl.lookup(decoder); err != nil {
		return nil, err
	}
	return tl, nil
}
