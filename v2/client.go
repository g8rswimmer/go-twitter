package twitter

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

const (
	tweetMaxIDs                  = 100
	userMaxIDs                   = 100
	userMaxNames                 = 100
	tweetRecentSearchQueryLength = 512
	tweetRecentCountsQueryLength = 512
	userBlocksMaxResults         = 1000
	userMutesMaxResults          = 1000
	userLikesMaxResults          = 100
	userLikesMinResults          = 10
	sampleStreamMaxBackoffMin    = 5
)

// Client is used to make twitter v2 API callouts.
type Client struct {
	Authorizer Authorizer
	Client     *http.Client
	Host       string
}

// CreateTweet will let a user post polls, quote tweets, tweet with reply setting, tweet with geo, attach
// perviously uploaded media toa tweet and tag users, tweet to super followers, etc.
func (c *Client) CreateTweet(ctx context.Context, tweet CreateTweetRequest) (*CreateTweetResponse, error) {
	if err := tweet.validate(); err != nil {
		return nil, err
	}
	body, err := json.Marshal(tweet)
	if err != nil {
		return nil, fmt.Errorf("create tweet marshal error %w", err)
	}
	ep := tweetCreateEndpoint.url(c.Host)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, ep, bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("create tweet request: %w", err)
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")
	c.Authorizer.Add(req)

	resp, err := c.Client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("create tweet response: %w", err)
	}
	defer resp.Body.Close()

	decoder := json.NewDecoder(resp.Body)

	if resp.StatusCode != http.StatusCreated {
		e := &ErrorResponse{}
		if err := decoder.Decode(e); err != nil {
			return nil, &HTTPError{
				Status:     resp.Status,
				StatusCode: resp.StatusCode,
				URL:        resp.Request.URL.String(),
			}
		}
		e.StatusCode = resp.StatusCode
		return nil, e
	}

	raw := &CreateTweetResponse{}
	if err := decoder.Decode(raw); err != nil {
		return nil, fmt.Errorf("create tweet decode response %w", err)
	}
	return raw, nil
}

// DeleteTweet allow the user to delete a specific tweet
func (c *Client) DeleteTweet(ctx context.Context, id string) (*DeleteTweetResponse, error) {
	if len(id) == 0 {
		return nil, fmt.Errorf("delete tweet id is required")
	}
	ep := tweetDeleteEndpoint.urlID(c.Host, id)

	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, ep, nil)
	if err != nil {
		return nil, fmt.Errorf("delete tweet request: %w", err)
	}
	req.Header.Add("Accept", "application/json")
	c.Authorizer.Add(req)

	resp, err := c.Client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("delete tweet response: %w", err)
	}
	defer resp.Body.Close()

	decoder := json.NewDecoder(resp.Body)

	if resp.StatusCode != http.StatusOK {
		e := &ErrorResponse{}
		if err := decoder.Decode(e); err != nil {
			return nil, &HTTPError{
				Status:     resp.Status,
				StatusCode: resp.StatusCode,
				URL:        resp.Request.URL.String(),
			}
		}
		e.StatusCode = resp.StatusCode
		return nil, e
	}

	raw := &DeleteTweetResponse{}
	if err := decoder.Decode(raw); err != nil {
		return nil, fmt.Errorf("delete tweet decode response %w", err)
	}
	return raw, nil
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
			return nil, &HTTPError{
				Status:     resp.Status,
				StatusCode: resp.StatusCode,
				URL:        resp.Request.URL.String(),
			}
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
			return nil, &HTTPError{
				Status:     resp.Status,
				StatusCode: resp.StatusCode,
				URL:        resp.Request.URL.String(),
			}
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

// UserRetweetLookup allows you to get information about users that have retweeted a tweet
func (c *Client) UserRetweetLookup(ctx context.Context, tweetID string, opts UserRetweetLookuoOpts) (*UserRetweetLookupResponse, error) {
	switch {
	case len(tweetID) == 0:
		return nil, fmt.Errorf("user retweet lookup: an id is required: %w", ErrParameter)
	default:
	}

	ep := userRetweetLookupEndpoint.urlID(c.Host, tweetID)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, ep, nil)
	if err != nil {
		return nil, fmt.Errorf("user retweet lookup request: %w", err)
	}
	req.Header.Add("Accept", "application/json")
	c.Authorizer.Add(req)
	opts.addQuery(req)

	resp, err := c.Client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("user retweet lookup response: %w", err)
	}
	defer resp.Body.Close()

	decoder := json.NewDecoder(resp.Body)

	if resp.StatusCode != http.StatusOK {
		e := &ErrorResponse{}
		if err := decoder.Decode(e); err != nil {
			return nil, &HTTPError{
				Status:     resp.Status,
				StatusCode: resp.StatusCode,
				URL:        resp.Request.URL.String(),
			}
		}
		e.StatusCode = resp.StatusCode
		return nil, e
	}

	raw := struct {
		*UserRetweetRaw
		Meta *UserRetweetMeta `json:"meta"`
	}{}
	if err := decoder.Decode(&raw); err != nil {
		return nil, fmt.Errorf("user retweet lookup dictionary: %w", err)
	}
	return &UserRetweetLookupResponse{
		Raw:  raw.UserRetweetRaw,
		Meta: raw.Meta,
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
			return nil, &HTTPError{
				Status:     resp.Status,
				StatusCode: resp.StatusCode,
				URL:        resp.Request.URL.String(),
			}
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

// AuthUserLookup will return the authorized user lookup
func (c *Client) AuthUserLookup(ctx context.Context, opts UserLookupOpts) (*UserLookupResponse, error) {
	ep := userAuthLookupEndpoint.url(c.Host)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, ep, nil)
	if err != nil {
		return nil, fmt.Errorf("auth user lookup request: %w", err)
	}
	req.Header.Add("Accept", "application/json")
	c.Authorizer.Add(req)
	opts.addQuery(req)

	resp, err := c.Client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("auth user lookup response: %w", err)
	}
	defer resp.Body.Close()

	decoder := json.NewDecoder(resp.Body)

	if resp.StatusCode != http.StatusOK {
		e := &ErrorResponse{}
		if err := decoder.Decode(e); err != nil {
			return nil, &HTTPError{
				Status:     resp.Status,
				StatusCode: resp.StatusCode,
				URL:        resp.Request.URL.String(),
			}
		}
		e.StatusCode = resp.StatusCode
		return nil, e
	}

	single := &userraw{}
	if err := decoder.Decode(single); err != nil {
		return nil, fmt.Errorf("user lookup single dictionary: %w", err)
	}
	raw := &UserRaw{}
	raw.Users = make([]*UserObj, 1)
	raw.Users[0] = single.User
	raw.Includes = single.Includes
	raw.Errors = single.Errors

	return &UserLookupResponse{
		Raw: raw,
	}, nil
}

// TweetRecentSearch will return a recent search based of a query
func (c *Client) TweetRecentSearch(ctx context.Context, query string, opts TweetRecentSearchOpts) (*TweetRecentSearchResponse, error) {
	switch {
	case len(query) == 0:
		return nil, fmt.Errorf("tweet recent search: a query is required: %w", ErrParameter)
	case len(query) > tweetRecentSearchQueryLength:
		return nil, fmt.Errorf("tweet recent search: the query over the length (%d): %w", tweetRecentSearchQueryLength, ErrParameter)
	default:
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, tweetRecentSearchEndpoint.url(c.Host), nil)
	if err != nil {
		return nil, fmt.Errorf("tweet recent search request: %w", err)
	}
	req.Header.Add("Accept", "application/json")
	c.Authorizer.Add(req)
	opts.addQuery(req)
	q := req.URL.Query()
	q.Add("query", query)
	req.URL.RawQuery = q.Encode()

	resp, err := c.Client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("tweet recent search response: %w", err)
	}
	defer resp.Body.Close()

	respBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("tweet recent search response read: %w", err)
	}
	if resp.StatusCode != http.StatusOK {
		e := &ErrorResponse{}
		if err := json.Unmarshal(respBytes, e); err != nil {
			return nil, &HTTPError{
				Status:     resp.Status,
				StatusCode: resp.StatusCode,
				URL:        resp.Request.URL.String(),
			}
		}
		e.StatusCode = resp.StatusCode
		return nil, e
	}

	recentSearch := &TweetRecentSearchResponse{
		Raw:  &TweetRaw{},
		Meta: &TweetRecentSearchMeta{},
	}

	if err := json.Unmarshal(respBytes, recentSearch.Raw); err != nil {
		return nil, fmt.Errorf("tweet recent search raw response error decode: %w", err)
	}

	if err := json.Unmarshal(respBytes, recentSearch); err != nil {
		return nil, fmt.Errorf("tweet recent search meta response error decode: %w", err)
	}

	return recentSearch, nil
}

func (c *Client) TweetSearchStreamAddRule(ctx context.Context, rules []TweetSearchStreamRule) (*TweetSearchStreamAddRuleResponse, error) {
	if len(rules) == 0 {
		return nil, fmt.Errorf("tweet search stream add rule: rules are required: %w", ErrParameter)
	}
	body := struct {
		Add tweetSearchStreamRules `json:"add"`
	}{
		Add: tweetSearchStreamRules(rules),
	}
	if err := body.Add.validate(); err != nil {
		return nil, err
	}
	enc, err := json.Marshal(body)
	if err != nil {
		return nil, fmt.Errorf("tweet search stream add rule body encoding %w", err)
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, tweetSearchStreamRulesEndpoint.url(c.Host), bytes.NewReader(enc))
	if err != nil {
		return nil, fmt.Errorf("tweet search stream add rule http request %w", err)
	}
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")
	c.Authorizer.Add(req)

	resp, err := c.Client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("tweet search stream add rule http response %w", err)
	}
	defer resp.Body.Close()

	decoder := json.NewDecoder(resp.Body)

	if resp.StatusCode != http.StatusOK {
		e := &ErrorResponse{}
		if err := decoder.Decode(e); err != nil {
			return nil, &HTTPError{
				Status:     resp.Status,
				StatusCode: resp.StatusCode,
				URL:        resp.Request.URL.String(),
			}
		}
		e.StatusCode = resp.StatusCode
		return nil, e
	}

	ruleResponse := &TweetSearchStreamAddRuleResponse{}
	if err := decoder.Decode(ruleResponse); err != nil {
		return nil, fmt.Errorf("tweet search stream add rule json response %w", err)
	}
	return ruleResponse, nil
}

func (c *Client) TweetSearchStreamDeleteRule(ctx context.Context, ruleIDs []TweetSearchStreamRuleID) (*TweetSearchStreamDeleteRuleResponse, error) {
	if len(ruleIDs) == 0 {
		return nil, fmt.Errorf("tweet search stream delete rule: rule ids are required: %w", ErrParameter)
	}
	body := struct {
		Delete tweetSeachStreamRuleIDs `json:"delete"`
	}{
		Delete: tweetSeachStreamRuleIDs(ruleIDs),
	}
	if err := body.Delete.validate(); err != nil {
		return nil, err
	}
	enc, err := json.Marshal(body)
	if err != nil {
		return nil, fmt.Errorf("tweet search stream delete rule body encoding %w", err)
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, tweetSearchStreamRulesEndpoint.url(c.Host), bytes.NewReader(enc))
	if err != nil {
		return nil, fmt.Errorf("tweet search stream delete rule http request %w", err)
	}
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")
	c.Authorizer.Add(req)

	resp, err := c.Client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("tweet search stream delete rule http response %w", err)
	}
	defer resp.Body.Close()

	decoder := json.NewDecoder(resp.Body)

	if resp.StatusCode != http.StatusOK {
		e := &ErrorResponse{}
		if err := decoder.Decode(e); err != nil {
			return nil, &HTTPError{
				Status:     resp.Status,
				StatusCode: resp.StatusCode,
				URL:        resp.Request.URL.String(),
			}
		}
		e.StatusCode = resp.StatusCode
		return nil, e
	}

	ruleResponse := &TweetSearchStreamDeleteRuleResponse{}
	if err := decoder.Decode(ruleResponse); err != nil {
		return nil, fmt.Errorf("tweet search stream delete rule json response %w", err)
	}
	return ruleResponse, nil
}

// TweetRecentCounts will return a recent tweet counts based of a query
func (c *Client) TweetRecentCounts(ctx context.Context, query string, opts TweetRecentCountsOpts) (*TweetRecentCountsResponse, error) {
	switch {
	case len(query) == 0:
		return nil, fmt.Errorf("tweet recent counts: a query is required: %w", ErrParameter)
	case len(query) > tweetRecentCountsQueryLength:
		return nil, fmt.Errorf("tweet recent counts: the query over the length (%d): %w", tweetRecentCountsQueryLength, ErrParameter)
	default:
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, tweetRecentCountsEndpoint.url(c.Host), nil)
	if err != nil {
		return nil, fmt.Errorf("tweet recent counts request: %w", err)
	}
	req.Header.Add("Accept", "application/json")
	c.Authorizer.Add(req)
	opts.addQuery(req)
	q := req.URL.Query()
	q.Add("query", query)
	req.URL.RawQuery = q.Encode()

	resp, err := c.Client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("tweet recent counts response: %w", err)
	}
	defer resp.Body.Close()

	respBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("tweet recent counts response read: %w", err)
	}
	if resp.StatusCode != http.StatusOK {
		e := &ErrorResponse{}
		if err := json.Unmarshal(respBytes, e); err != nil {
			return nil, &HTTPError{
				Status:     resp.Status,
				StatusCode: resp.StatusCode,
				URL:        resp.Request.URL.String(),
			}
		}
		e.StatusCode = resp.StatusCode
		return nil, e
	}

	recentCounts := &TweetRecentCountsResponse{
		TweetCounts: []*TweetCount{},
		Meta:        &TweetRecentCountsMeta{},
	}

	if err := json.Unmarshal(respBytes, recentCounts); err != nil {
		return nil, fmt.Errorf("tweet recent counts response error decode: %w", err)
	}

	return recentCounts, nil
}

// UserFollowingLookup will return a user's following users
func (c *Client) UserFollowingLookup(ctx context.Context, id string, opts UserFollowingLookupOpts) (*UserFollowingLookupResponse, error) {
	if len(id) == 0 {
		return nil, fmt.Errorf("user following lookup: id is required: %w", ErrParameter)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, userFollowingEndpoint.urlID(c.Host, id), nil)
	if err != nil {
		return nil, fmt.Errorf("user following lookup request: %w", err)
	}
	req.Header.Add("Accept", "application/json")
	c.Authorizer.Add(req)
	opts.addQuery(req)
	q := req.URL.Query()
	req.URL.RawQuery = q.Encode()

	resp, err := c.Client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("user following lookup response: %w", err)
	}
	defer resp.Body.Close()

	respBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("user following lookup response read: %w", err)
	}
	if resp.StatusCode != http.StatusOK {
		e := &ErrorResponse{}
		if err := json.Unmarshal(respBytes, e); err != nil {
			return nil, &HTTPError{
				Status:     resp.Status,
				StatusCode: resp.StatusCode,
				URL:        resp.Request.URL.String(),
			}
		}
		e.StatusCode = resp.StatusCode
		return nil, e
	}

	followingLookup := &UserFollowingLookupResponse{
		Raw:  &UserRaw{},
		Meta: &UserFollowinghMeta{},
	}

	if err := json.Unmarshal(respBytes, followingLookup.Raw); err != nil {
		return nil, fmt.Errorf("user following lookup raw response error decode: %w", err)
	}

	if err := json.Unmarshal(respBytes, followingLookup); err != nil {
		return nil, fmt.Errorf("user following lookup meta response error decode: %w", err)
	}

	return followingLookup, nil
}

// UserFollows allows a user ID to follow another user
func (c *Client) UserFollows(ctx context.Context, userID, targetUserID string) (*UserFollowsResponse, error) {
	switch {
	case len(userID) == 0:
		return nil, fmt.Errorf("user follows: user id is required %w", ErrParameter)
	case len(targetUserID) == 0:
		return nil, fmt.Errorf("user follows: target user id is required %w", ErrParameter)
	default:
	}

	reqBody := struct {
		TargetUserID string `json:"target_user_id"`
	}{
		TargetUserID: targetUserID,
	}
	enc, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("user follows: json marshal %w", err)
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, userFollowingEndpoint.urlID(c.Host, userID), bytes.NewReader(enc))
	if err != nil {
		return nil, fmt.Errorf("user follows request: %w", err)
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")
	c.Authorizer.Add(req)

	resp, err := c.Client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("user follows response: %w", err)
	}
	defer resp.Body.Close()

	decoder := json.NewDecoder(resp.Body)

	if resp.StatusCode != http.StatusOK {
		e := &ErrorResponse{}
		if err := decoder.Decode(e); err != nil {
			return nil, &HTTPError{
				Status:     resp.Status,
				StatusCode: resp.StatusCode,
				URL:        resp.Request.URL.String(),
			}
		}
		e.StatusCode = resp.StatusCode
		return nil, e
	}

	raw := &UserFollowsResponse{}
	if err := decoder.Decode(raw); err != nil {
		return nil, fmt.Errorf("user follows decode response %w", err)
	}
	return raw, nil
}

// DeleteUserFollows llows a user ID to unfollow another user
func (c *Client) DeleteUserFollows(ctx context.Context, userID, targetUserID string) (*UserDeleteFollowsResponse, error) {
	switch {
	case len(userID) == 0:
		return nil, fmt.Errorf("user delete follows: user id is required %w", ErrParameter)
	case len(targetUserID) == 0:
		return nil, fmt.Errorf("user delete follows: target user id is required %w", ErrParameter)
	default:
	}

	ep := userFollowingEndpoint.urlID(c.Host, userID) + "/" + targetUserID
	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, ep, nil)
	if err != nil {
		return nil, fmt.Errorf("user delete follows request: %w", err)
	}
	req.Header.Add("Accept", "application/json")
	c.Authorizer.Add(req)

	resp, err := c.Client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("user delete follows response: %w", err)
	}
	defer resp.Body.Close()

	decoder := json.NewDecoder(resp.Body)

	if resp.StatusCode != http.StatusOK {
		e := &ErrorResponse{}
		if err := decoder.Decode(e); err != nil {
			return nil, &HTTPError{
				Status:     resp.Status,
				StatusCode: resp.StatusCode,
				URL:        resp.Request.URL.String(),
			}
		}
		e.StatusCode = resp.StatusCode
		return nil, e
	}

	raw := &UserDeleteFollowsResponse{}
	if err := decoder.Decode(raw); err != nil {
		return nil, fmt.Errorf("user delete follows decode response %w", err)
	}
	return raw, nil
}

// UserFollowersLookup will return a user's followers
func (c *Client) UserFollowersLookup(ctx context.Context, id string, opts UserFollowersLookupOpts) (*UserFollowersLookupResponse, error) {
	if len(id) == 0 {
		return nil, fmt.Errorf("user followers lookup: id is required: %w", ErrParameter)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, userFollowersEndpoint.urlID(c.Host, id), nil)
	if err != nil {
		return nil, fmt.Errorf("user followers lookup request: %w", err)
	}
	req.Header.Add("Accept", "application/json")
	c.Authorizer.Add(req)
	opts.addQuery(req)
	q := req.URL.Query()
	req.URL.RawQuery = q.Encode()

	resp, err := c.Client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("user followers lookup response: %w", err)
	}
	defer resp.Body.Close()

	respBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("user followers lookup response read: %w", err)
	}
	if resp.StatusCode != http.StatusOK {
		e := &ErrorResponse{}
		if err := json.Unmarshal(respBytes, e); err != nil {
			return nil, &HTTPError{
				Status:     resp.Status,
				StatusCode: resp.StatusCode,
				URL:        resp.Request.URL.String(),
			}
		}
		e.StatusCode = resp.StatusCode
		return nil, e
	}

	followersLookup := &UserFollowersLookupResponse{
		Raw:  &UserRaw{},
		Meta: &UserFollowershMeta{},
	}

	if err := json.Unmarshal(respBytes, followersLookup.Raw); err != nil {
		return nil, fmt.Errorf("user followers lookup raw response error decode: %w", err)
	}

	if err := json.Unmarshal(respBytes, followersLookup); err != nil {
		return nil, fmt.Errorf("user followers lookup meta response error decode: %w", err)
	}

	return followersLookup, nil
}

// UserTweetTimeline will return the user tweet timeline
func (c *Client) UserTweetTimeline(ctx context.Context, userID string, opts UserTweetTimelineOpts) (*UserTweetTimelineResponse, error) {
	switch {
	case len(userID) == 0:
		return nil, fmt.Errorf("user tweet timeline: a query is required: %w", ErrParameter)
	default:
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, userTweetTimelineEndpoint.urlID(c.Host, userID), nil)
	if err != nil {
		return nil, fmt.Errorf("user tweet timeline request: %w", err)
	}
	req.Header.Add("Accept", "application/json")
	c.Authorizer.Add(req)
	opts.addQuery(req)

	resp, err := c.Client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("user tweet timeline response: %w", err)
	}
	defer resp.Body.Close()

	respBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("user tweet timeline response read: %w", err)
	}
	if resp.StatusCode != http.StatusOK {
		e := &ErrorResponse{}
		if err := json.Unmarshal(respBytes, e); err != nil {
			return nil, &HTTPError{
				Status:     resp.Status,
				StatusCode: resp.StatusCode,
				URL:        resp.Request.URL.String(),
			}
		}
		e.StatusCode = resp.StatusCode
		return nil, e
	}

	timeline := &UserTweetTimelineResponse{
		Raw:  &TweetRaw{},
		Meta: &UserTimelineMeta{},
	}

	if err := json.Unmarshal(respBytes, timeline.Raw); err != nil {
		return nil, fmt.Errorf("user tweet timeline raw response error decode: %w", err)
	}

	if err := json.Unmarshal(respBytes, timeline); err != nil {
		return nil, fmt.Errorf("user tweet timeline meta response error decode: %w", err)
	}

	return timeline, nil
}

// UserMentionTimeline will return the user's mentions timeline
func (c *Client) UserMentionTimeline(ctx context.Context, userID string, opts UserMentionTimelineOpts) (*UserMentionTimelineResponse, error) {
	switch {
	case len(userID) == 0:
		return nil, fmt.Errorf("user mention timeline: a query is required: %w", ErrParameter)
	default:
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, userMentionTimelineEndpoint.urlID(c.Host, userID), nil)
	if err != nil {
		return nil, fmt.Errorf("user mention timeline request: %w", err)
	}
	req.Header.Add("Accept", "application/json")
	c.Authorizer.Add(req)
	opts.addQuery(req)

	resp, err := c.Client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("user mention timeline response: %w", err)
	}
	defer resp.Body.Close()

	respBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("user mention timeline response read: %w", err)
	}
	if resp.StatusCode != http.StatusOK {
		e := &ErrorResponse{}
		if err := json.Unmarshal(respBytes, e); err != nil {
			return nil, &HTTPError{
				Status:     resp.Status,
				StatusCode: resp.StatusCode,
				URL:        resp.Request.URL.String(),
			}
		}
		e.StatusCode = resp.StatusCode
		return nil, e
	}

	timeline := &UserMentionTimelineResponse{
		Raw:  &TweetRaw{},
		Meta: &UserTimelineMeta{},
	}

	if err := json.Unmarshal(respBytes, timeline.Raw); err != nil {
		return nil, fmt.Errorf("user mention timeline raw response error decode: %w", err)
	}

	if err := json.Unmarshal(respBytes, timeline); err != nil {
		return nil, fmt.Errorf("user mention timeline meta response error decode: %w", err)
	}

	return timeline, nil
}

// TweetHideReplies will hide the replies for a given tweet
func (c Client) TweetHideReplies(ctx context.Context, id string, hide bool) error {
	if len(id) == 0 {
		return fmt.Errorf("tweet hide replies: id must be present %w", ErrParameter)
	}
	type body struct {
		Hidden bool `json:"hidden"`
	}
	rb := body{
		Hidden: hide,
	}
	enc, err := json.Marshal(rb)
	if err != nil {
		return fmt.Errorf("tweet hide replies: request body marshal %w", err)
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPut, tweetHideRepliesEndpoint.urlID(c.Host, id), bytes.NewReader(enc))
	if err != nil {
		return fmt.Errorf("tweet hide replies request: %w", err)
	}
	req.Header.Add("Accept", "application/json")
	c.Authorizer.Add(req)

	resp, err := c.Client.Do(req)
	if err != nil {
		return fmt.Errorf("tweet hide replies response: %w", err)
	}
	defer resp.Body.Close()

	respBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("tweet hide replies response read: %w", err)
	}
	if resp.StatusCode != http.StatusOK {
		errResp := &ErrorResponse{}
		if err := json.Unmarshal(respBytes, errResp); err != nil {
			return &HTTPError{
				Status:     resp.Status,
				StatusCode: resp.StatusCode,
				URL:        resp.Request.URL.String(),
			}
		}
		errResp.StatusCode = resp.StatusCode
		return errResp
	}

	type responseData struct {
		Data body `json:"data"`
	}
	rd := &responseData{}
	if err := json.Unmarshal(respBytes, rd); err != nil {
		return fmt.Errorf("tweet hide replies response error decode: %w", err)
	}
	if rd.Data.Hidden != hide {
		return fmt.Errorf("tweet hide replies response unable to hide %v", hide)
	}
	return nil
}

// UserRetweet will retweet a tweet for a user
func (c *Client) UserRetweet(ctx context.Context, userID, tweetID string) (*UserRetweetResponse, error) {
	switch {
	case len(userID) == 0:
		return nil, fmt.Errorf("user retweet: user id is required %w", ErrParameter)
	case len(tweetID) == 0:
		return nil, fmt.Errorf("user retweet: tweet id is required %w", ErrParameter)
	default:
	}

	reqBody := struct {
		TweetID string `json:"tweet_id"`
	}{
		TweetID: tweetID,
	}
	enc, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("user retweet: json marshal %w", err)
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, userManageRetweetEndpoint.urlID(c.Host, userID), bytes.NewReader(enc))
	if err != nil {
		return nil, fmt.Errorf("user retweet request: %w", err)
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")
	c.Authorizer.Add(req)

	resp, err := c.Client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("user retweet response: %w", err)
	}
	defer resp.Body.Close()

	decoder := json.NewDecoder(resp.Body)

	if resp.StatusCode != http.StatusOK {
		e := &ErrorResponse{}
		if err := decoder.Decode(e); err != nil {
			return nil, &HTTPError{
				Status:     resp.Status,
				StatusCode: resp.StatusCode,
				URL:        resp.Request.URL.String(),
			}
		}
		e.StatusCode = resp.StatusCode
		return nil, e
	}

	raw := &UserRetweetResponse{}
	if err := decoder.Decode(raw); err != nil {
		return nil, fmt.Errorf("user retweet decode response %w", err)
	}
	return raw, nil
}

// DeleteUserRetweet will return a retweet from a user
func (c *Client) DeleteUserRetweet(ctx context.Context, userID, tweetID string) (*DeleteUserRetweetResponse, error) {
	switch {
	case len(userID) == 0:
		return nil, fmt.Errorf("user delete retweet: user id is required %w", ErrParameter)
	case len(tweetID) == 0:
		return nil, fmt.Errorf("user delete retweet: tweet id is required %w", ErrParameter)
	default:
	}

	ep := userManageRetweetEndpoint.urlID(c.Host, userID) + "/" + tweetID
	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, ep, nil)
	if err != nil {
		return nil, fmt.Errorf("user delete retweet request: %w", err)
	}
	req.Header.Add("Accept", "application/json")
	c.Authorizer.Add(req)

	resp, err := c.Client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("user delete retweet response: %w", err)
	}
	defer resp.Body.Close()

	decoder := json.NewDecoder(resp.Body)

	if resp.StatusCode != http.StatusOK {
		e := &ErrorResponse{}
		if err := decoder.Decode(e); err != nil {
			return nil, &HTTPError{
				Status:     resp.Status,
				StatusCode: resp.StatusCode,
				URL:        resp.Request.URL.String(),
			}
		}
		e.StatusCode = resp.StatusCode
		return nil, e
	}

	raw := &DeleteUserRetweetResponse{}
	if err := decoder.Decode(raw); err != nil {
		return nil, fmt.Errorf("user delete retweet decode response %w", err)
	}
	return raw, nil
}

// UserBlocksLookup returns a list of users who are blocked by the user ID
func (c *Client) UserBlocksLookup(ctx context.Context, userID string, opts UserBlocksLookupOpts) (*UserBlocksLookupResponse, error) {
	switch {
	case len(userID) == 0:
		return nil, fmt.Errorf("user blocked lookup: user id is required: %w", ErrParameter)
	case opts.MaxResults > userBlocksMaxResults:
		return nil, fmt.Errorf("user blocked lookup: max results can't be above %d: %w", userBlocksMaxResults, ErrParameter)
	default:
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, userBlocksEndpoint.urlID(c.Host, userID), nil)
	if err != nil {
		return nil, fmt.Errorf("user blocked lookup request: %w", err)
	}
	req.Header.Add("Accept", "application/json")
	c.Authorizer.Add(req)
	opts.addQuery(req)
	q := req.URL.Query()
	req.URL.RawQuery = q.Encode()

	resp, err := c.Client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("user blocked lookup response: %w", err)
	}
	defer resp.Body.Close()

	respBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("user blocked lookup response read: %w", err)
	}
	if resp.StatusCode != http.StatusOK {
		e := &ErrorResponse{}
		if err := json.Unmarshal(respBytes, e); err != nil {
			return nil, &HTTPError{
				Status:     resp.Status,
				StatusCode: resp.StatusCode,
				URL:        resp.Request.URL.String(),
			}
		}
		e.StatusCode = resp.StatusCode
		return nil, e
	}

	blockedLookup := &UserBlocksLookupResponse{
		Raw:  &UserRaw{},
		Meta: &UserBlocksLookupMeta{},
	}

	if err := json.Unmarshal(respBytes, blockedLookup.Raw); err != nil {
		return nil, fmt.Errorf("user blocked lookup raw response error decode: %w", err)
	}

	if err := json.Unmarshal(respBytes, blockedLookup); err != nil {
		return nil, fmt.Errorf("user blocked lookup meta response error decode: %w", err)
	}

	return blockedLookup, nil
}

// UserBlocks will have the user block the targeted user ID
func (c *Client) UserBlocks(ctx context.Context, userID, targetUserID string) (*UserBlocksResponse, error) {
	switch {
	case len(userID) == 0:
		return nil, fmt.Errorf("user blocks: user id is required %w", ErrParameter)
	case len(targetUserID) == 0:
		return nil, fmt.Errorf("user blocks: target user id is required %w", ErrParameter)
	default:
	}

	reqBody := struct {
		TargetUserID string `json:"target_user_id"`
	}{
		TargetUserID: targetUserID,
	}
	enc, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("user blocks: json marshal %w", err)
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, userBlocksEndpoint.urlID(c.Host, userID), bytes.NewReader(enc))
	if err != nil {
		return nil, fmt.Errorf("user blocks request: %w", err)
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")
	c.Authorizer.Add(req)

	resp, err := c.Client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("user blocks response: %w", err)
	}
	defer resp.Body.Close()

	decoder := json.NewDecoder(resp.Body)

	if resp.StatusCode != http.StatusOK {
		e := &ErrorResponse{}
		if err := decoder.Decode(e); err != nil {
			return nil, &HTTPError{
				Status:     resp.Status,
				StatusCode: resp.StatusCode,
				URL:        resp.Request.URL.String(),
			}
		}
		e.StatusCode = resp.StatusCode
		return nil, e
	}

	raw := &UserBlocksResponse{}
	if err := decoder.Decode(raw); err != nil {
		return nil, fmt.Errorf("user blocks decode response %w", err)
	}
	return raw, nil
}

// DeleteUserBlocks will remove the target user block
func (c *Client) DeleteUserBlocks(ctx context.Context, userID, targetUserID string) (*UserDeleteBlocksResponse, error) {
	switch {
	case len(userID) == 0:
		return nil, fmt.Errorf("user delete blocks: user id is required %w", ErrParameter)
	case len(targetUserID) == 0:
		return nil, fmt.Errorf("user delete blocks: target user id is required %w", ErrParameter)
	default:
	}

	ep := userBlocksEndpoint.urlID(c.Host, userID) + "/" + targetUserID
	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, ep, nil)
	if err != nil {
		return nil, fmt.Errorf("user delete blocks request: %w", err)
	}
	req.Header.Add("Accept", "application/json")
	c.Authorizer.Add(req)

	resp, err := c.Client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("user delete blocks response: %w", err)
	}
	defer resp.Body.Close()

	decoder := json.NewDecoder(resp.Body)

	if resp.StatusCode != http.StatusOK {
		e := &ErrorResponse{}
		if err := decoder.Decode(e); err != nil {
			return nil, &HTTPError{
				Status:     resp.Status,
				StatusCode: resp.StatusCode,
				URL:        resp.Request.URL.String(),
			}
		}
		e.StatusCode = resp.StatusCode
		return nil, e
	}

	raw := &UserDeleteBlocksResponse{}
	if err := decoder.Decode(raw); err != nil {
		return nil, fmt.Errorf("user delete blocks decode response %w", err)
	}
	return raw, nil
}

// UserMutesLookup returns a list of users who are muted by the user ID
func (c *Client) UserMutesLookup(ctx context.Context, userID string, opts UserMutesLookupOpts) (*UserMutesLookupResponse, error) {
	switch {
	case len(userID) == 0:
		return nil, fmt.Errorf("user muted lookup: user id is required: %w", ErrParameter)
	case opts.MaxResults > userBlocksMaxResults:
		return nil, fmt.Errorf("user muted lookup: max results can't be above %d: %w", userMutesMaxResults, ErrParameter)
	default:
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, userMutesEndpont.urlID(c.Host, userID), nil)
	if err != nil {
		return nil, fmt.Errorf("user muted lookup request: %w", err)
	}
	req.Header.Add("Accept", "application/json")
	c.Authorizer.Add(req)
	opts.addQuery(req)
	q := req.URL.Query()
	req.URL.RawQuery = q.Encode()

	resp, err := c.Client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("user muted lookup response: %w", err)
	}
	defer resp.Body.Close()

	respBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("user muted lookup response read: %w", err)
	}
	if resp.StatusCode != http.StatusOK {
		e := &ErrorResponse{}
		if err := json.Unmarshal(respBytes, e); err != nil {
			return nil, &HTTPError{
				Status:     resp.Status,
				StatusCode: resp.StatusCode,
				URL:        resp.Request.URL.String(),
			}
		}
		e.StatusCode = resp.StatusCode
		return nil, e
	}

	mutedLookup := &UserMutesLookupResponse{
		Raw:  &UserRaw{},
		Meta: &UserMutesLookupMeta{},
	}

	if err := json.Unmarshal(respBytes, mutedLookup.Raw); err != nil {
		return nil, fmt.Errorf("user muted lookup raw response error decode: %w", err)
	}

	if err := json.Unmarshal(respBytes, mutedLookup); err != nil {
		return nil, fmt.Errorf("user muted lookup meta response error decode: %w", err)
	}

	return mutedLookup, nil
}

// UserMutes allows an authenticated user ID to mute the target user
func (c *Client) UserMutes(ctx context.Context, userID, targetUserID string) (*UserMutesResponse, error) {
	switch {
	case len(userID) == 0:
		return nil, fmt.Errorf("user mutes: user id is required %w", ErrParameter)
	case len(targetUserID) == 0:
		return nil, fmt.Errorf("user mutes: target user id is required %w", ErrParameter)
	default:
	}

	reqBody := struct {
		TargetUserID string `json:"target_user_id"`
	}{
		TargetUserID: targetUserID,
	}
	enc, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("user mutes: json marshal %w", err)
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, userMutesEndpont.urlID(c.Host, userID), bytes.NewReader(enc))
	if err != nil {
		return nil, fmt.Errorf("user mutes request: %w", err)
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")
	c.Authorizer.Add(req)

	resp, err := c.Client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("user mutes response: %w", err)
	}
	defer resp.Body.Close()

	decoder := json.NewDecoder(resp.Body)

	if resp.StatusCode != http.StatusOK {
		e := &ErrorResponse{}
		if err := decoder.Decode(e); err != nil {
			return nil, &HTTPError{
				Status:     resp.Status,
				StatusCode: resp.StatusCode,
				URL:        resp.Request.URL.String(),
			}
		}
		e.StatusCode = resp.StatusCode
		return nil, e
	}

	raw := &UserMutesResponse{}
	if err := decoder.Decode(raw); err != nil {
		return nil, fmt.Errorf("user mutes decode response %w", err)
	}
	return raw, nil
}

// DeleteUserMutes allows an authenticated user ID to unmute the target user
func (c *Client) DeleteUserMutes(ctx context.Context, userID, targetUserID string) (*UserDeleteMutesResponse, error) {
	switch {
	case len(userID) == 0:
		return nil, fmt.Errorf("user delete mutes: user id is required %w", ErrParameter)
	case len(targetUserID) == 0:
		return nil, fmt.Errorf("user delete mutes: target user id is required %w", ErrParameter)
	default:
	}

	ep := userMutesEndpont.urlID(c.Host, userID) + "/" + targetUserID
	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, ep, nil)
	if err != nil {
		return nil, fmt.Errorf("user delete mutes request: %w", err)
	}
	req.Header.Add("Accept", "application/json")
	c.Authorizer.Add(req)

	resp, err := c.Client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("user delete mutes response: %w", err)
	}
	defer resp.Body.Close()

	decoder := json.NewDecoder(resp.Body)

	if resp.StatusCode != http.StatusOK {
		e := &ErrorResponse{}
		if err := decoder.Decode(e); err != nil {
			return nil, &HTTPError{
				Status:     resp.Status,
				StatusCode: resp.StatusCode,
				URL:        resp.Request.URL.String(),
			}
		}
		e.StatusCode = resp.StatusCode
		return nil, e
	}

	raw := &UserDeleteMutesResponse{}
	if err := decoder.Decode(raw); err != nil {
		return nil, fmt.Errorf("user delete mutes decode response %w", err)
	}
	return raw, nil
}

// TweetLikesLookup gets information about a tweet's liking users.  The response will have at most 100 users who liked the tweet
func (c *Client) TweetLikesLookup(ctx context.Context, tweetID string, opts TweetLikesLookupOpts) (*TweetLikesLookupResponse, error) {
	switch {
	case len(tweetID) == 0:
		return nil, fmt.Errorf("user tweet likes lookup: an id is required: %w", ErrParameter)
	default:
	}

	ep := tweetLikesEndpoint.urlID(c.Host, tweetID)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, ep, nil)
	if err != nil {
		return nil, fmt.Errorf("user tweet likes lookup request: %w", err)
	}
	req.Header.Add("Accept", "application/json")
	c.Authorizer.Add(req)
	opts.addQuery(req)

	resp, err := c.Client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("user tweet likes lookup response: %w", err)
	}
	defer resp.Body.Close()

	decoder := json.NewDecoder(resp.Body)

	if resp.StatusCode != http.StatusOK {
		e := &ErrorResponse{}
		if err := decoder.Decode(e); err != nil {
			return nil, &HTTPError{
				Status:     resp.Status,
				StatusCode: resp.StatusCode,
				URL:        resp.Request.URL.String(),
			}
		}
		e.StatusCode = resp.StatusCode
		return nil, e
	}

	raw := &UserRaw{}

	if err := decoder.Decode(&raw); err != nil {
		return nil, fmt.Errorf("user tweet likes lookup dictionary: %w", err)
	}
	return &TweetLikesLookupResponse{
		Raw: raw,
	}, nil
}

// UserLikesLookup gets informaiton about a user's liked tweets.
func (c *Client) UserLikesLookup(ctx context.Context, userID string, opts UserLikesLookupOpts) (*UserLikesLookupResponse, error) {
	switch {
	case len(userID) == 0:
		return nil, fmt.Errorf("tweet user likes lookup: an id is required: %w", ErrParameter)
	case opts.MaxResults == 0:
	case opts.MaxResults < userLikesMinResults:
		return nil, fmt.Errorf("tweet user likes lookup: a min results [%d] is required [current: %d]: %w", userLikesMinResults, opts.MaxResults, ErrParameter)
	case opts.MaxResults > userLikesMaxResults:
		return nil, fmt.Errorf("tweet user likes lookup: a max results [%d] is required [current: %d]: %w", userLikesMaxResults, opts.MaxResults, ErrParameter)
	default:
	}

	ep := userLikedTweetEndpoint.urlID(c.Host, userID)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, ep, nil)
	if err != nil {
		return nil, fmt.Errorf("tweet user likes lookup request: %w", err)
	}
	req.Header.Add("Accept", "application/json")
	c.Authorizer.Add(req)
	opts.addQuery(req)

	resp, err := c.Client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("tweet user likes lookup response: %w", err)
	}
	defer resp.Body.Close()

	decoder := json.NewDecoder(resp.Body)

	if resp.StatusCode != http.StatusOK {
		e := &ErrorResponse{}
		if err := decoder.Decode(e); err != nil {
			return nil, &HTTPError{
				Status:     resp.Status,
				StatusCode: resp.StatusCode,
				URL:        resp.Request.URL.String(),
			}
		}
		e.StatusCode = resp.StatusCode
		return nil, e
	}

	respBody := struct {
		*TweetRaw
		Meta *UserLikesMeta `json:"meta"`
	}{}

	if err := decoder.Decode(&respBody); err != nil {
		return nil, fmt.Errorf("tweet user likes lookup dictionary: %w", err)
	}
	return &UserLikesLookupResponse{
		Raw:  respBody.TweetRaw,
		Meta: respBody.Meta,
	}, nil
}

// UserLikes will like the targeted tweet
func (c *Client) UserLikes(ctx context.Context, userID, tweetID string) (*UserLikesResponse, error) {
	switch {
	case len(userID) == 0:
		return nil, fmt.Errorf("user likes: user id is required %w", ErrParameter)
	case len(tweetID) == 0:
		return nil, fmt.Errorf("user likes: tweet id is required %w", ErrParameter)
	default:
	}

	reqBody := struct {
		TweetID string `json:"tweet_id"`
	}{
		TweetID: tweetID,
	}
	enc, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("user likes: json marshal %w", err)
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, userLikesEndpoint.urlID(c.Host, userID), bytes.NewReader(enc))
	if err != nil {
		return nil, fmt.Errorf("user likes request: %w", err)
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")
	c.Authorizer.Add(req)

	resp, err := c.Client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("user likes response: %w", err)
	}
	defer resp.Body.Close()

	decoder := json.NewDecoder(resp.Body)

	if resp.StatusCode != http.StatusOK {
		e := &ErrorResponse{}
		if err := decoder.Decode(e); err != nil {
			return nil, &HTTPError{
				Status:     resp.Status,
				StatusCode: resp.StatusCode,
				URL:        resp.Request.URL.String(),
			}
		}
		e.StatusCode = resp.StatusCode
		return nil, e
	}

	raw := &UserLikesResponse{}
	if err := decoder.Decode(raw); err != nil {
		return nil, fmt.Errorf("user likes decode response %w", err)
	}
	return raw, nil
}

// DeleteUserLikes will unlike the targeted tweet
func (c *Client) DeleteUserLikes(ctx context.Context, userID, tweetID string) (*DeteleUserLikesResponse, error) {
	switch {
	case len(userID) == 0:
		return nil, fmt.Errorf("user delete likes: user id is required %w", ErrParameter)
	case len(tweetID) == 0:
		return nil, fmt.Errorf("user delete likes: tweet id is required %w", ErrParameter)
	default:
	}

	ep := userLikesEndpoint.urlID(c.Host, userID) + "/" + tweetID
	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, ep, nil)
	if err != nil {
		return nil, fmt.Errorf("user delete likes request: %w", err)
	}
	req.Header.Add("Accept", "application/json")
	c.Authorizer.Add(req)

	resp, err := c.Client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("user delete likes response: %w", err)
	}
	defer resp.Body.Close()

	decoder := json.NewDecoder(resp.Body)

	if resp.StatusCode != http.StatusOK {
		e := &ErrorResponse{}
		if err := decoder.Decode(e); err != nil {
			return nil, &HTTPError{
				Status:     resp.Status,
				StatusCode: resp.StatusCode,
				URL:        resp.Request.URL.String(),
			}
		}
		e.StatusCode = resp.StatusCode
		return nil, e
	}

	raw := &DeteleUserLikesResponse{}
	if err := decoder.Decode(raw); err != nil {
		return nil, fmt.Errorf("user likes retweet decode response %w", err)
	}
	return raw, nil
}

// TweetSampleStream will return a streamer for streaming 1% of all tweets real-time
func (c *Client) TweetSampleStream(ctx context.Context, opts TweetSampleStreamOpts) (*TweetStream, error) {
	switch {
	case opts.BackfillMinutes == 0:
	case opts.BackfillMinutes > sampleStreamMaxBackoffMin:
		return nil, fmt.Errorf("tweet sample stream: a max backoff minutes [%d] is [current: %d]: %w", sampleStreamMaxBackoffMin, opts.BackfillMinutes, ErrParameter)
	default:
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, tweetSampleStreamEndpoint.url(c.Host), nil)
	if err != nil {
		return nil, fmt.Errorf("tweet sample stream request: %w", err)
	}
	req.Header.Add("Accept", "application/json")
	c.Authorizer.Add(req)
	opts.addQuery(req)

	resp, err := c.Client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("tweet sample stream response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		defer resp.Body.Close()
		e := &ErrorResponse{}
		if err := json.NewDecoder(resp.Body).Decode(e); err != nil {
			return nil, &HTTPError{
				Status:     resp.Status,
				StatusCode: resp.StatusCode,
				URL:        resp.Request.URL.String(),
			}
		}
		e.StatusCode = resp.StatusCode
		return nil, e
	}

	return StartTweetStream(resp.Body), nil
}
