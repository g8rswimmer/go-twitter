package twitter

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

const (
	tweetLookupEndpoint       = "2/tweets"
	tweetRecentSearchEndpoint = "2/tweets/search/recent"
	tweetSearchStreamRules    = "2/tweets/search/stream/rules"
	tweetSearchStream         = "2/tweets/search/stream"
	tweetSampledStream        = "2/tweets/sample/stream"
	tweetMaxIDs               = 100
	tweetQuerySize            = 512
)

type TweetLookup struct {
}

type TweetRecentSearch struct {
}

type TweetSearchStreamRules struct {
}

type TweetSearchStreamRule struct {
}

type TweetSearchStream struct {
}

type TweetSampledStream struct {
}

type TweetError struct {
	StatusCode int
}

func (t *TweetError) Error() string {
	return ""
}

type TweetRecentSearchParameters struct {
	query     string
	StartTime time.Time
	EndTime   time.Time
	MaxResult int
	NextToken string
	SinceID   string
	UntilID   string
}

func (t TweetRecentSearchParameters) Encode(req *http.Request) {
}

type TweetLookupParameters struct {
	ids []string
}

func (t TweetLookupParameters) Encode(req *http.Request) {
}

type TweetSearchStreamParameters struct {
}

func (t TweetSearchStreamParameters) Encode(req *http.Request) {

}

type TweetSampledStreamParameters struct {
}

func (t TweetSampledStreamParameters) Encode(req *http.Request) {

}

type Tweet struct {
	Authorizer Authorizer
	Client     *http.Client
	Host       string
}

func (t *Tweet) Lookup(ctx context.Context, ids []string, parameters TweetLookupParameters) (*TweetLookup, error) {
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
	parameters.Encode(req)

	resp, err := t.Client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("tweet lookup response: %w", err)
	}
	defer resp.Body.Close()

	decoder := json.NewDecoder(resp.Body)

	if resp.StatusCode != http.StatusOK {
		e := &TweetError{}
		if err := decoder.Decode(e); err != nil {
			return nil, fmt.Errorf("tweet lookup response error decode: %w", err)
		}
		e.StatusCode = resp.StatusCode
		return nil, e
	}

	tr := &TweetLookup{}
	if err := decoder.Decode(tr); err != nil {
		return nil, fmt.Errorf("tweet lookup response decode: %w", err)
	}
	return tr, nil
}

func (t *Tweet) RecentSearch(ctx context.Context, query string, parameters TweetRecentSearchParameters) (*TweetRecentSearch, error) {
	switch {
	case len(query) == 0:
		return nil, fmt.Errorf("tweet recent search query must be present")
	case len(query) > tweetQuerySize:
		return nil, fmt.Errorf("tweet recent search query size %d greater than max %d", len(query), tweetQuerySize)
	default:
		parameters.query = query
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, fmt.Sprintf("%s/%s", t.Host, tweetRecentSearchEndpoint), nil)
	if err != nil {
		return nil, fmt.Errorf("tweet recent search request: %w", err)
	}
	req.Header.Add("Accept", "application/json")
	t.Authorizer.Add(req)
	parameters.Encode(req)

	resp, err := t.Client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("tweet recent search response: %w", err)
	}
	defer resp.Body.Close()

	decoder := json.NewDecoder(resp.Body)

	if resp.StatusCode != http.StatusOK {
		e := &TweetError{}
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

func (t *Tweet) SearchStreamRules(ctx context.Context, ids []string) (*TweetSearchStreamRules, error) {
	if len(ids) > tweetMaxIDs {
		return nil, fmt.Errorf("tweet search stream rules: ids %d is greater than max %d", len(ids), tweetMaxIDs)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, fmt.Sprintf("%s/%s", t.Host, tweetSearchStreamRules), nil)
	if err != nil {
		return nil, fmt.Errorf("tweet rsearch stream rules request: %w", err)
	}
	req.Header.Add("Accept", "application/json")
	t.Authorizer.Add(req)
	if len(ids) > 0 {
		q := req.URL.Query()
		q.Add("ids", strings.Join(ids, ","))
		req.URL.RawQuery = q.Encode()
	}

	resp, err := t.Client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("tweet search stream rules response: %w", err)
	}
	defer resp.Body.Close()

	decoder := json.NewDecoder(resp.Body)

	if resp.StatusCode != http.StatusOK {
		e := &TweetError{}
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

func (t *Tweet) AddSearchStreamRules(ctx context.Context, rules []TweetSearchStreamRule, validate bool) (*TweetSearchStreamRules, error) {
	if len(rules) == 0 {
		return nil, errors.New("tweet search stream add rules need new rules can not be zero")
	}
	add := map[string][]TweetSearchStreamRule{
		"add": rules,
	}
	enc, err := json.Marshal(add)
	if err != nil {
		return nil, fmt.Errorf("tweet search stream add rules encode error %w", err)
	}

	return t.updateSearchStreamRules(ctx, bytes.NewReader(enc), validate)

}

func (t *Tweet) DeleteSearchStreamRules(ctx context.Context, ids []string, validate bool) (*TweetSearchStreamRules, error) {
	if len(ids) == 0 {
		return nil, errors.New("tweet search stream add rules need new ids can not be zero")
	}
	delete := map[string]map[string][]string{
		"delete": map[string][]string{
			"ids": ids,
		},
	}
	enc, err := json.Marshal(delete)
	if err != nil {
		return nil, fmt.Errorf("tweet search stream add rules encode error %w", err)
	}
	return t.updateSearchStreamRules(ctx, bytes.NewReader(enc), validate)
}

func (t *Tweet) updateSearchStreamRules(ctx context.Context, body io.Reader, validate bool) (*TweetSearchStreamRules, error) {

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, fmt.Sprintf("%s/%s", t.Host, tweetSearchStreamRules), body)
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
		e := &TweetError{}
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

func (t *Tweet) Sampledtream(ctx context.Context, parameters TweetSampledStreamParameters) (*TweetSampledStream, error) {

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, fmt.Sprintf("%s/%s", t.Host, tweetSampledStream), nil)
	if err != nil {
		return nil, fmt.Errorf("tweet rsearch stream request: %w", err)
	}
	req.Header.Add("Accept", "application/json")
	t.Authorizer.Add(req)
	parameters.Encode(req)

	resp, err := t.Client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("tweet search stream response: %w", err)
	}
	defer resp.Body.Close()

	decoder := json.NewDecoder(resp.Body)

	if resp.StatusCode != http.StatusOK {
		e := &TweetError{}
		if err := decoder.Decode(e); err != nil {
			return nil, fmt.Errorf("tweet search stream response error decode: %w", err)
		}
		e.StatusCode = resp.StatusCode
		return nil, e
	}

	ts := &TweetSampledStream{}
	if err := decoder.Decode(ts); err != nil {
		return nil, fmt.Errorf("tweet search stream response decode: %w", err)
	}
	return ts, nil
}
