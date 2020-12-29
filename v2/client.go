package twitter

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

const (
	tweetMaxIDs  = 100
	userMaxIDs   = 100
	userMaxNames = 100
)

// Client is used to make twitter v2 API callouts.
type Client struct {
	Authorizer Authorizer
	Client     *http.Client
	Host       string
}

// TweetLookup returns information about a tweet or group of tweets specified by a group of tweet ids.
func (c *Client) TweetLookup(ctx context.Context, ids []string, opts TweetLookupOpts) (*TweetLookupResponse, error) {
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

	raw := &TweetRaw{}
	switch {
	case len(ids) == 1:
		single := &tweetraw{}
		if err := decoder.Decode(single); err != nil {
			return nil, fmt.Errorf("tweet lookup single dictionary: %w", err)
		}
		raw.Tweets = make([]*TweetObj, 1)
		raw.Tweets[0] = single.Tweet
		raw.Includes = single.Includes
		raw.Errors = single.Errors
	default:
		if err := decoder.Decode(raw); err != nil {
			return nil, fmt.Errorf("tweet lookup dictionary: %w", err)
		}
	}
	return &TweetLookupResponse{
		Raw: raw,
	}, nil
}

// UserLookup returns information about an user or group of users specified by a group of user ids.
func (c *Client) UserLookup(ctx context.Context, ids []string, opts UserLookupOpts) (*UserLookupResponse, error) {
	ep := userLookupEndpoint.url(c.Host)
	switch {
	case len(ids) == 0:
		return nil, fmt.Errorf("user lookup: an id is required: %w", ErrParameter)
	case len(ids) > userMaxIDs:
		return nil, fmt.Errorf("user lookup: ids %d is greater than max %d: %w", len(ids), userMaxIDs, ErrParameter)
	case len(ids) == 1:
		ep += fmt.Sprintf("/%s", ids[0])
	default:
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, ep, nil)
	if err != nil {
		return nil, fmt.Errorf("user lookup request: %w", err)
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
		return nil, fmt.Errorf("user lookup response: %w", err)
	}
	defer resp.Body.Close()

	decoder := json.NewDecoder(resp.Body)

	if resp.StatusCode != http.StatusOK {
		e := &ErrorResponse{}
		if err := decoder.Decode(e); err != nil {
			return nil, fmt.Errorf("user lookup response error decode: %w", err)
		}
		e.StatusCode = resp.StatusCode
		return nil, e
	}

	raw := &UserRaw{}
	switch {
	case len(ids) == 1:
		single := &userraw{}
		if err := decoder.Decode(single); err != nil {
			return nil, fmt.Errorf("user lookup single dictionary: %w", err)
		}
		raw.Users = make([]*UserObj, 1)
		raw.Users[0] = single.User
		raw.Includes = single.Includes
		raw.Errors = single.Errors
	default:
		if err := decoder.Decode(raw); err != nil {
			return nil, fmt.Errorf("user lookup dictionary: %w", err)
		}
	}
	return &UserLookupResponse{
		Raw: raw,
	}, nil
}

// UserNameLookup returns information about an user or group of users specified by a group of usernames.
func (c *Client) UserNameLookup(ctx context.Context, usernames []string, opts UserLookupOpts) (*UserLookupResponse, error) {
	ep := userNameLookupEndpoint.url(c.Host)
	switch {
	case len(usernames) == 0:
		return nil, fmt.Errorf("username lookup: an username is required: %w", ErrParameter)
	case len(usernames) > userMaxIDs:
		return nil, fmt.Errorf("username lookup: usernames %d is greater than max %d: %w", len(usernames), userMaxNames, ErrParameter)
	case len(usernames) == 1:
		ep += fmt.Sprintf("/username/%s", usernames[0])
	default:
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, ep, nil)
	if err != nil {
		return nil, fmt.Errorf("username lookup request: %w", err)
	}
	req.Header.Add("Accept", "application/json")
	c.Authorizer.Add(req)
	opts.addQuery(req)
	if len(usernames) > 1 {
		q := req.URL.Query()
		q.Add("usernames", strings.Join(usernames, ","))
		req.URL.RawQuery = q.Encode()
	}

	resp, err := c.Client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("username lookup response: %w", err)
	}
	defer resp.Body.Close()

	decoder := json.NewDecoder(resp.Body)

	if resp.StatusCode != http.StatusOK {
		e := &ErrorResponse{}
		if err := decoder.Decode(e); err != nil {
			return nil, fmt.Errorf("username lookup response error decode: %w", err)
		}
		e.StatusCode = resp.StatusCode
		return nil, e
	}

	raw := &UserRaw{}
	switch {
	case len(usernames) == 1:
		single := &userraw{}
		if err := decoder.Decode(single); err != nil {
			return nil, fmt.Errorf("username lookup single dictionary: %w", err)
		}
		raw.Users = make([]*UserObj, 1)
		raw.Users[0] = single.User
		raw.Includes = single.Includes
		raw.Errors = single.Errors
	default:
		if err := decoder.Decode(raw); err != nil {
			return nil, fmt.Errorf("username lookup dictionary: %w", err)
		}
	}
	return &UserLookupResponse{
		Raw: raw,
	}, nil
}
