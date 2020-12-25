package twitter

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

const (
	tweetMaxIDs = 100
)

// Client is used to make twitter v2 API callouts.
type Client struct {
	Authorizer Authorizer
	Client     *http.Client
	Host       string
}

// TweetLookup returns information about a tweet or group of tweets specified by a group of tweet ids.
func (c *Client) TweetLookup(ctx context.Context, ids []string, opts TweetLookupOpts) (*TweetDictionary, error) {
	ep := tweetLookupEndpoint.url(c.Host)
	switch {
	case len(ids) == 0:
		return nil, fmt.Errorf("tweet lookup: an id is required: %w", ErrParameter)
	case len(ids) > tweetMaxIDs:
		return nil, fmt.Errorf("tweet lookup: ids %d is greater than max %d: %w", len(ids), tweetMaxIDs, ErrParameter)
	case len(ids) == 1:
		ep += fmt.Sprintf("/%s", ids[0])
	default:
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, ep, nil)
	if err != nil {
		return nil, fmt.Errorf("tweet lookup request: %w", err)
	}
	req.Header.Add("Accept", "application/json")
	c.Authorizer.Add(req)
	opts.addQuery(req)
	if len(ids) > 1 {
		q := req.URL.Query()
		q.Add("ids", strings.Join(ids, ","))
		req.URL.RawQuery = q.Encode()
	}

	resp, err := c.Client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("tweet lookup response: %w", err)
	}
	defer resp.Body.Close()

	decoder := json.NewDecoder(resp.Body)

	if resp.StatusCode != http.StatusOK {
		e := &ErrorResponse{}
		if err := decoder.Decode(e); err != nil {
			return nil, fmt.Errorf("tweet lookup response error decode: %w", err)
		}
		e.StatusCode = resp.StatusCode
		return nil, e
	}

	dictionary := &TweetDictionary{}
	switch {
	case len(ids) == 1:
		single := &tweetdictionary{}
		if err := decoder.Decode(single); err != nil {
			return nil, fmt.Errorf("tweet lookup single dictionary: %w", err)
		}
		dictionary.Tweets = make([]*TweetObj, 1)
		dictionary.Tweets[0] = single.Tweet
		dictionary.Includes = single.Includes
		dictionary.Errors = single.Errors
	default:
		if err := decoder.Decode(dictionary); err != nil {
			return nil, fmt.Errorf("tweet lookup dictionary: %w", err)
		}
	}
	return dictionary, nil
}
