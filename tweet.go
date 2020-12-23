package twitter

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
)

const (
	tweetLookupEndpoint              = "2/tweets"
	tweetRecentSearchEndpoint        = "2/tweets/search/recent"
	tweetFilteredStreamRulesEndpoint = "/2/tweets/search/stream/rules"
	tweetFilteredStreamEndpoint      = "/2/tweets/search/stream"
	tweetSampledStreamEndpoint       = "/2/tweets/sample/stream"
	tweetHideEndpoint                = "/2/tweets/{id}/hidden"
	tweetID                          = "{id}"
	tweetMaxIDs                      = 100
	tweetQuerySize                   = 512
)

// TweetLookups is .a map of tweet lookups
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
func (t *Tweet) Lookup(ctx context.Context, ids []string, options TweetFieldOptions) (TweetLookups, error) {
	ep := tweetLookupEndpoint
	switch {
	case len(ids) == 0:
		return nil, fmt.Errorf("tweet lookup an id is required")
	case len(ids) > tweetMaxIDs:
		return nil, fmt.Errorf("tweet lookup: ids %d is greater than max %d", len(ids), tweetMaxIDs)
	case len(ids) == 1:
		ep += fmt.Sprintf("/%s", ids[0])
	default:
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, fmt.Sprintf("%s/%s", t.Host, ep), nil)
	if err != nil {
		return nil, fmt.Errorf("tweet lookup request: %w", err)
	}
	req.Header.Add("Accept", "application/json")
	t.Authorizer.Add(req)
	options.addQuery(req)
	if len(ids) > 1 {
		q := req.URL.Query()
		q.Add("ids", strings.Join(ids, ","))
		req.URL.RawQuery = q.Encode()
	}

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
func (t *Tweet) RecentSearch(ctx context.Context, query string, searchOpts TweetRecentSearchOptions, fieldOpts TweetFieldOptions) (*TweetRecentSearch, error) {
	switch {
	case len(query) == 0:
		return nil, fmt.Errorf("tweet recent search query must be present")
	case len(query) > tweetQuerySize:
		return nil, fmt.Errorf("tweet recent search query size %d greater than max %d", len(query), tweetQuerySize)
	case searchOpts.MaxResult > 0 && (searchOpts.MaxResult < 10 || searchOpts.MaxResult > 100):
		return nil, fmt.Errorf("tweet resent search max result needs to be between 10 -100 (%d", searchOpts.MaxResult)
	default:
		searchOpts.query = query
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, fmt.Sprintf("%s/%s", t.Host, tweetRecentSearchEndpoint), nil)
	if err != nil {
		return nil, fmt.Errorf("tweet recent search request: %w", err)
	}
	req.Header.Add("Accept", "application/json")
	t.Authorizer.Add(req)
	searchOpts.addQuery(req)
	fieldOpts.addQuery(req)

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

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, fmt.Sprintf("%s/%s", t.Host, tweetFilteredStreamRulesEndpoint), bytes.NewReader(enc))
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

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, fmt.Sprintf("%s/%s", t.Host, tweetFilteredStreamRulesEndpoint), nil)
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
func (t *Tweet) FilteredStream(ctx context.Context, options TweetFieldOptions) (TweetLookups, error) {

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, fmt.Sprintf("%s/%s", t.Host, tweetFilteredStreamEndpoint), nil)
	if err != nil {
		return nil, fmt.Errorf("tweet lookup request: %w", err)
	}
	req.Header.Add("Accept", "application/json")
	t.Authorizer.Add(req)
	options.addQuery(req)

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

// SampledStream will stream about 1% of all tweets
func (t *Tweet) SampledStream(ctx context.Context, options TweetFieldOptions) (TweetLookups, error) {

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, fmt.Sprintf("%s/%s", t.Host, tweetSampledStreamEndpoint), nil)
	if err != nil {
		return nil, fmt.Errorf("tweet lookup request: %w", err)
	}
	req.Header.Add("Accept", "application/json")
	t.Authorizer.Add(req)
	options.addQuery(req)

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

// HideReplies will hide a tweet id replies
func (t *Tweet) HideReplies(ctx context.Context, id string, hidden bool) error {
	if len(id) == 0 {
		return errors.New("tweet hidden: id can not be empty")
	}

	hb := struct {
		Hidden bool `json:"hidden"`
	}{
		Hidden: hidden,
	}
	enc, _ := json.Marshal(hb)

	ep := strings.ReplaceAll(tweetHideEndpoint, tweetID, id)
	req, err := http.NewRequestWithContext(ctx, http.MethodPut, fmt.Sprintf("%s/%s", t.Host, ep), bytes.NewReader(enc))
	if err != nil {
		return fmt.Errorf("tweet lookup request: %w", err)
	}
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-type", "application/json")
	t.Authorizer.Add(req)

	resp, err := t.Client.Do(req)
	if err != nil {
		return fmt.Errorf("tweet lookup response: %w", err)
	}
	defer resp.Body.Close()

	decoder := json.NewDecoder(resp.Body)

	if resp.StatusCode != http.StatusOK {
		e := &TweetErrorResponse{}
		if err := decoder.Decode(e); err != nil {
			return fmt.Errorf("tweet lookup response error decode: %w", err)
		}
		e.StatusCode = resp.StatusCode
		return e
	}

	type responseData struct {
		Hidden bool `json:"hidden"`
	}
	type response struct {
		Data responseData `json:"data"`
	}
	r := &response{}
	if err := decoder.Decode(r); err != nil {
		return fmt.Errorf("tweet hidden: response decode err %w", err)
	}
	if r.Data.Hidden != hidden {
		return fmt.Errorf("tweet hidden: expected response (%v) does not match hidden (%v)", r.Data.Hidden, hidden)
	}
	return nil
}
