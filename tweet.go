package twitter

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

const (
	tweetLookupEndpoint = "2/tweets"
	// tweetRecentSearchEndpoint = "2/tweets/search/recent"
	// tweetSearchStreamRules    = "2/tweets/search/stream/rules"
	// tweetSearchStream         = "2/tweets/search/stream"
	// tweetSampledStream        = "2/tweets/sample/stream"
	tweetMaxIDs = 100
	// tweetQuerySize            = 512
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

// type TweetRecentSearch struct {
// }

// type TweetSearchStreamRules struct {
// }

// type TweetSearchStreamRule struct {
// }

// type TweetSearchStream struct {
// }

// type TweetSampledStream struct {
// }

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

// type TweetRecentSearchParameters struct {
// 	query     string
// 	StartTime time.Time
// 	EndTime   time.Time
// 	MaxResult int
// 	NextToken string
// 	SinceID   string
// 	UntilID   string
// }

// func (t TweetRecentSearchParameters) Encode(req *http.Request) {
// }

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

// type TweetSearchStreamParameters struct {
// }

// func (t TweetSearchStreamParameters) Encode(req *http.Request) {

// }

// type TweetSampledStreamParameters struct {
// }

// func (t TweetSampledStreamParameters) Encode(req *http.Request) {

// }

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

// func (t *Tweet) RecentSearch(ctx context.Context, query string, parameters TweetRecentSearchParameters) (*TweetRecentSearch, error) {
// 	switch {
// 	case len(query) == 0:
// 		return nil, fmt.Errorf("tweet recent search query must be present")
// 	case len(query) > tweetQuerySize:
// 		return nil, fmt.Errorf("tweet recent search query size %d greater than max %d", len(query), tweetQuerySize)
// 	default:
// 		parameters.query = query
// 	}

// 	req, err := http.NewRequestWithContext(ctx, http.MethodGet, fmt.Sprintf("%s/%s", t.Host, tweetRecentSearchEndpoint), nil)
// 	if err != nil {
// 		return nil, fmt.Errorf("tweet recent search request: %w", err)
// 	}
// 	req.Header.Add("Accept", "application/json")
// 	t.Authorizer.Add(req)
// 	parameters.Encode(req)

// 	resp, err := t.Client.Do(req)
// 	if err != nil {
// 		return nil, fmt.Errorf("tweet recent search response: %w", err)
// 	}
// 	defer resp.Body.Close()

// 	decoder := json.NewDecoder(resp.Body)

// 	if resp.StatusCode != http.StatusOK {
// 		e := &TweetErrorResponse{}
// 		if err := decoder.Decode(e); err != nil {
// 			return nil, fmt.Errorf("tweet recent search response error decode: %w", err)
// 		}
// 		e.StatusCode = resp.StatusCode
// 		return nil, e
// 	}

// 	tr := &TweetRecentSearch{}
// 	if err := decoder.Decode(tr); err != nil {
// 		return nil, fmt.Errorf("tweet recent search response decode: %w", err)
// 	}
// 	return tr, nil
// }

// func (t *Tweet) SearchStreamRules(ctx context.Context, ids []string) (*TweetSearchStreamRules, error) {
// 	if len(ids) > tweetMaxIDs {
// 		return nil, fmt.Errorf("tweet search stream rules: ids %d is greater than max %d", len(ids), tweetMaxIDs)
// 	}

// 	req, err := http.NewRequestWithContext(ctx, http.MethodGet, fmt.Sprintf("%s/%s", t.Host, tweetSearchStreamRules), nil)
// 	if err != nil {
// 		return nil, fmt.Errorf("tweet rsearch stream rules request: %w", err)
// 	}
// 	req.Header.Add("Accept", "application/json")
// 	t.Authorizer.Add(req)
// 	if len(ids) > 0 {
// 		q := req.URL.Query()
// 		q.Add("ids", strings.Join(ids, ","))
// 		req.URL.RawQuery = q.Encode()
// 	}

// 	resp, err := t.Client.Do(req)
// 	if err != nil {
// 		return nil, fmt.Errorf("tweet search stream rules response: %w", err)
// 	}
// 	defer resp.Body.Close()

// 	decoder := json.NewDecoder(resp.Body)

// 	if resp.StatusCode != http.StatusOK {
// 		e := &TweetErrorResponse{}
// 		if err := decoder.Decode(e); err != nil {
// 			return nil, fmt.Errorf("tweet search stream rules response error decode: %w", err)
// 		}
// 		e.StatusCode = resp.StatusCode
// 		return nil, e
// 	}

// 	tr := &TweetSearchStreamRules{}
// 	if err := decoder.Decode(tr); err != nil {
// 		return nil, fmt.Errorf("tweet search stream rules response decode: %w", err)
// 	}
// 	return tr, nil
// }

// func (t *Tweet) AddSearchStreamRules(ctx context.Context, rules []TweetSearchStreamRule, validate bool) (*TweetSearchStreamRules, error) {
// 	if len(rules) == 0 {
// 		return nil, errors.New("tweet search stream add rules need new rules can not be zero")
// 	}
// 	add := map[string][]TweetSearchStreamRule{
// 		"add": rules,
// 	}
// 	enc, err := json.Marshal(add)
// 	if err != nil {
// 		return nil, fmt.Errorf("tweet search stream add rules encode error %w", err)
// 	}

// 	return t.updateSearchStreamRules(ctx, bytes.NewReader(enc), validate)

// }

// func (t *Tweet) DeleteSearchStreamRules(ctx context.Context, ids []string, validate bool) (*TweetSearchStreamRules, error) {
// 	if len(ids) == 0 {
// 		return nil, errors.New("tweet search stream add rules need new ids can not be zero")
// 	}
// 	delete := map[string]map[string][]string{
// 		"delete": map[string][]string{
// 			"ids": ids,
// 		},
// 	}
// 	enc, err := json.Marshal(delete)
// 	if err != nil {
// 		return nil, fmt.Errorf("tweet search stream add rules encode error %w", err)
// 	}
// 	return t.updateSearchStreamRules(ctx, bytes.NewReader(enc), validate)
// }

// func (t *Tweet) updateSearchStreamRules(ctx context.Context, body io.Reader, validate bool) (*TweetSearchStreamRules, error) {

// 	req, err := http.NewRequestWithContext(ctx, http.MethodPost, fmt.Sprintf("%s/%s", t.Host, tweetSearchStreamRules), body)
// 	if err != nil {
// 		return nil, fmt.Errorf("tweet rsearch stream rules request: %w", err)
// 	}

// 	req.Header.Add("Accept", "application/json")
// 	req.Header.Add("Content-type", "application/json")
// 	t.Authorizer.Add(req)
// 	if validate {
// 		q := req.URL.Query()
// 		q.Add("dry_run", "true")
// 		req.URL.RawQuery = q.Encode()
// 	}

// 	resp, err := t.Client.Do(req)
// 	if err != nil {
// 		return nil, fmt.Errorf("tweet search stream rules response: %w", err)
// 	}
// 	defer resp.Body.Close()

// 	decoder := json.NewDecoder(resp.Body)

// 	if resp.StatusCode != http.StatusOK {
// 		e := &TweetErrorResponse{}
// 		if err := decoder.Decode(e); err != nil {
// 			return nil, fmt.Errorf("tweet search stream rules response error decode: %w", err)
// 		}
// 		e.StatusCode = resp.StatusCode
// 		return nil, e
// 	}

// 	tr := &TweetSearchStreamRules{}
// 	if err := decoder.Decode(tr); err != nil {
// 		return nil, fmt.Errorf("tweet search stream rules response decode: %w", err)
// 	}
// 	return tr, nil
// }

// func (t *Tweet) Sampledtream(ctx context.Context, parameters TweetSampledStreamParameters) (*TweetSampledStream, error) {

// 	req, err := http.NewRequestWithContext(ctx, http.MethodGet, fmt.Sprintf("%s/%s", t.Host, tweetSampledStream), nil)
// 	if err != nil {
// 		return nil, fmt.Errorf("tweet rsearch stream request: %w", err)
// 	}
// 	req.Header.Add("Accept", "application/json")
// 	t.Authorizer.Add(req)
// 	parameters.Encode(req)

// 	resp, err := t.Client.Do(req)
// 	if err != nil {
// 		return nil, fmt.Errorf("tweet search stream response: %w", err)
// 	}
// 	defer resp.Body.Close()

// 	decoder := json.NewDecoder(resp.Body)

// 	if resp.StatusCode != http.StatusOK {
// 		e := &TweetErrorResponse{}
// 		if err := decoder.Decode(e); err != nil {
// 			return nil, fmt.Errorf("tweet search stream response error decode: %w", err)
// 		}
// 		e.StatusCode = resp.StatusCode
// 		return nil, e
// 	}

// 	ts := &TweetSampledStream{}
// 	if err := decoder.Decode(ts); err != nil {
// 		return nil, fmt.Errorf("tweet search stream response decode: %w", err)
// 	}
// 	return ts, nil
// }
