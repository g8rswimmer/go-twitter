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
	tweetMaxIDs                                     = 100
	userMaxIDs                                      = 100
	spaceMaxIDs                                     = 100
	spaceByCreatorMaxIDs                            = 100
	userMaxNames                                    = 100
	tweetRecentSearchQueryLength                    = 512
	tweetSearchQueryLength                          = 1024
	tweetRecentCountsQueryLength                    = 512
	tweetAllCountsQueryLength                       = 1024
	userBlocksMaxResults                            = 1000
	userMutesMaxResults                             = 1000
	likesMaxResults                                 = 100
	likesMinResults                                 = 10
	sampleStreamMaxBackOffMin                       = 5
	userListMaxResults                              = 100
	listTweetMaxResults                             = 100
	userListMembershipMaxResults                    = 100
	listUserMemberMaxResults                        = 100
	userListFollowedMaxResults                      = 100
	listUserFollowersMaxResults                     = 100
	quoteTweetMaxResults                            = 100
	quoteTweetMinResults                            = 10
	tweetBookmarksMaxResults                        = 100
	userTweetTimelineMinResults                     = 5
	userTweetTimelineMaxResults                     = 100
	userMentionTimelineMinResults                   = 5
	userMentionTimelineMaxResults                   = 100
	userRetweetLookupMaxResults                     = 100
	userTweetReverseChronologicalTimelineMinResults = 1
	userTweetReverseChronologicalTimelineMaxResults = 100
)

// Client is used to make twitter v2 API callouts.
//
// Authorizer is used to add auth to the request
//
// Client is the HTTP client to use for all requests
//
// Host is the base URL to use like, https://api.twitter.com
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

	rl := rateFromHeader(resp.Header)

	if resp.StatusCode != http.StatusCreated {
		e := &ErrorResponse{}
		if err := decoder.Decode(e); err != nil {
			return nil, &HTTPError{
				Status:     resp.Status,
				StatusCode: resp.StatusCode,
				URL:        resp.Request.URL.String(),
				RateLimit:  rl,
			}
		}
		e.StatusCode = resp.StatusCode
		e.RateLimit = rl
		return nil, e
	}

	raw := &CreateTweetResponse{}
	if err := decoder.Decode(raw); err != nil {
		return nil, &ResponseDecodeError{
			Name:      "create tweet",
			Err:       err,
			RateLimit: rl,
		}
	}
	raw.RateLimit = rl
	return raw, nil
}

// DeleteTweet allow the user to delete a specific tweet
func (c *Client) DeleteTweet(ctx context.Context, id string) (*DeleteTweetResponse, error) {
	if len(id) == 0 {
		return nil, fmt.Errorf("delete tweet id is required %w", ErrParameter)
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

	rl := rateFromHeader(resp.Header)

	if resp.StatusCode != http.StatusOK {
		e := &ErrorResponse{}
		if err := decoder.Decode(e); err != nil {
			return nil, &HTTPError{
				Status:     resp.Status,
				StatusCode: resp.StatusCode,
				URL:        resp.Request.URL.String(),
				RateLimit:  rl,
			}
		}
		e.StatusCode = resp.StatusCode
		e.RateLimit = rl
		return nil, e
	}

	raw := &DeleteTweetResponse{}
	if err := decoder.Decode(raw); err != nil {
		return nil, &ResponseDecodeError{
			Name:      "delete tweet",
			Err:       err,
			RateLimit: rl,
		}
	}
	raw.RateLimit = rl
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

	rl := rateFromHeader(resp.Header)

	if resp.StatusCode != http.StatusOK {
		e := &ErrorResponse{}
		if err := decoder.Decode(e); err != nil {
			return nil, &HTTPError{
				Status:     resp.Status,
				StatusCode: resp.StatusCode,
				URL:        resp.Request.URL.String(),
				RateLimit:  rl,
			}
		}
		e.StatusCode = resp.StatusCode
		e.RateLimit = rl
		return nil, e
	}

	raw := &TweetRaw{}
	switch {
	case len(ids) == 1:
		single := &tweetraw{}
		if err := decoder.Decode(single); err != nil {
			return nil, &ResponseDecodeError{
				Name:      "tweet lookup",
				Err:       err,
				RateLimit: rl,
			}
		}
		raw.Tweets = make([]*TweetObj, 1)
		raw.Tweets[0] = single.Tweet
		raw.Includes = single.Includes
		raw.Errors = single.Errors
	default:
		if err := decoder.Decode(raw); err != nil {
			return nil, &ResponseDecodeError{
				Name:      "tweet lookup ",
				Err:       err,
				RateLimit: rl,
			}
		}
	}
	return &TweetLookupResponse{
		Raw:       raw,
		RateLimit: rl,
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

	rl := rateFromHeader(resp.Header)

	if resp.StatusCode != http.StatusOK {
		e := &ErrorResponse{}
		if err := decoder.Decode(e); err != nil {
			return nil, &HTTPError{
				Status:     resp.Status,
				StatusCode: resp.StatusCode,
				URL:        resp.Request.URL.String(),
				RateLimit:  rl,
			}
		}
		e.StatusCode = resp.StatusCode
		e.RateLimit = rl
		return nil, e
	}

	raw := &UserRaw{}
	switch {
	case len(ids) == 1:
		single := &userraw{}
		if err := decoder.Decode(single); err != nil {
			return nil, &ResponseDecodeError{
				Name:      "user lookup",
				Err:       err,
				RateLimit: rl,
			}
		}
		raw.Users = make([]*UserObj, 1)
		raw.Users[0] = single.User
		raw.Includes = single.Includes
		raw.Errors = single.Errors
	default:
		if err := decoder.Decode(raw); err != nil {
			return nil, &ResponseDecodeError{
				Name:      "user lookup",
				Err:       err,
				RateLimit: rl,
			}
		}
	}
	return &UserLookupResponse{
		Raw:       raw,
		RateLimit: rl,
	}, nil
}

// UserRetweetLookup allows you to get information about users that have retweeted a tweet
func (c *Client) UserRetweetLookup(ctx context.Context, tweetID string, opts UserRetweetLookupOpts) (*UserRetweetLookupResponse, error) {
	switch {
	case len(tweetID) == 0:
		return nil, fmt.Errorf("user retweet lookup: an id is required: %w", ErrParameter)
	case opts.MaxResults > userRetweetLookupMaxResults:
		return nil, fmt.Errorf("user retweet lookup: max results is limited to 100 [%d], %w", opts.MaxResults, ErrParameter)
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

	rl := rateFromHeader(resp.Header)

	if resp.StatusCode != http.StatusOK {
		e := &ErrorResponse{}
		if err := decoder.Decode(e); err != nil {
			return nil, &HTTPError{
				Status:     resp.Status,
				StatusCode: resp.StatusCode,
				URL:        resp.Request.URL.String(),
				RateLimit:  rl,
			}
		}
		e.StatusCode = resp.StatusCode
		e.RateLimit = rl
		return nil, e
	}

	raw := struct {
		*UserRetweetRaw
		Meta *UserRetweetMeta `json:"meta"`
	}{}
	if err := decoder.Decode(&raw); err != nil {
		return nil, &ResponseDecodeError{
			Name:      "user retweet lookup",
			Err:       err,
			RateLimit: rl,
		}
	}
	return &UserRetweetLookupResponse{
		Raw:       raw.UserRetweetRaw,
		Meta:      raw.Meta,
		RateLimit: rl,
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

	rl := rateFromHeader(resp.Header)

	if resp.StatusCode != http.StatusOK {
		e := &ErrorResponse{}
		if err := decoder.Decode(e); err != nil {
			return nil, &HTTPError{
				Status:     resp.Status,
				StatusCode: resp.StatusCode,
				URL:        resp.Request.URL.String(),
				RateLimit:  rl,
			}
		}
		e.StatusCode = resp.StatusCode
		e.RateLimit = rl
		return nil, e
	}

	raw := &UserRaw{}
	switch {
	case len(usernames) == 1:
		single := &userraw{}
		if err := decoder.Decode(single); err != nil {
			return nil, &ResponseDecodeError{
				Name:      "username lookup",
				Err:       err,
				RateLimit: rl,
			}
		}
		raw.Users = make([]*UserObj, 1)
		raw.Users[0] = single.User
		raw.Includes = single.Includes
		raw.Errors = single.Errors
	default:
		if err := decoder.Decode(raw); err != nil {
			return nil, &ResponseDecodeError{
				Name:      "username lookup",
				Err:       err,
				RateLimit: rl,
			}
		}
	}
	return &UserLookupResponse{
		Raw:       raw,
		RateLimit: rl,
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

	rl := rateFromHeader(resp.Header)

	if resp.StatusCode != http.StatusOK {
		e := &ErrorResponse{}
		if err := decoder.Decode(e); err != nil {
			return nil, &HTTPError{
				Status:     resp.Status,
				StatusCode: resp.StatusCode,
				URL:        resp.Request.URL.String(),
				RateLimit:  rl,
			}
		}
		e.StatusCode = resp.StatusCode
		e.RateLimit = rl
		return nil, e
	}

	single := &userraw{}
	if err := decoder.Decode(single); err != nil {
		return nil, &ResponseDecodeError{
			Name:      "auth user lookup",
			Err:       err,
			RateLimit: rl,
		}
	}
	raw := &UserRaw{}
	raw.Users = make([]*UserObj, 1)
	raw.Users[0] = single.User
	raw.Includes = single.Includes
	raw.Errors = single.Errors

	return &UserLookupResponse{
		Raw:       raw,
		RateLimit: rl,
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

	rl := rateFromHeader(resp.Header)

	if resp.StatusCode != http.StatusOK {
		e := &ErrorResponse{}
		if err := json.Unmarshal(respBytes, e); err != nil {
			return nil, &HTTPError{
				Status:     resp.Status,
				StatusCode: resp.StatusCode,
				URL:        resp.Request.URL.String(),
				RateLimit:  rl,
			}
		}
		e.StatusCode = resp.StatusCode
		e.RateLimit = rl
		return nil, e
	}

	recentSearch := &TweetRecentSearchResponse{
		Raw:       &TweetRaw{},
		Meta:      &TweetRecentSearchMeta{},
		RateLimit: rl,
	}

	if err := json.Unmarshal(respBytes, recentSearch.Raw); err != nil {
		return nil, &ResponseDecodeError{
			Name:      "tweet recent search",
			Err:       err,
			RateLimit: rl,
		}
	}

	if err := json.Unmarshal(respBytes, recentSearch); err != nil {
		return nil, &ResponseDecodeError{
			Name:      "tweet recent search",
			Err:       err,
			RateLimit: rl,
		}
	}

	return recentSearch, nil
}

// TweetSearch is a full-archive search endpoint returns the complete history of public Tweets matching a search query.
//
// This endpoint is only available to those users who have been approved for Academic Research access.
func (c *Client) TweetSearch(ctx context.Context, query string, opts TweetSearchOpts) (*TweetSearchResponse, error) {
	switch {
	case len(query) == 0:
		return nil, fmt.Errorf("tweet search: a query is required: %w", ErrParameter)
	case len(query) > tweetSearchQueryLength:
		return nil, fmt.Errorf("tweet search: the query over the length (%d): %w", tweetSearchQueryLength, ErrParameter)
	default:
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, tweetSearchEndpoint.url(c.Host), nil)
	if err != nil {
		return nil, fmt.Errorf("tweet search request: %w", err)
	}
	req.Header.Add("Accept", "application/json")
	c.Authorizer.Add(req)
	opts.addQuery(req)
	q := req.URL.Query()
	q.Add("query", query)
	req.URL.RawQuery = q.Encode()

	resp, err := c.Client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("tweet search response: %w", err)
	}
	defer resp.Body.Close()

	decoder := json.NewDecoder(resp.Body)

	rl := rateFromHeader(resp.Header)

	if resp.StatusCode != http.StatusOK {
		e := &ErrorResponse{}
		if err := decoder.Decode(e); err != nil {
			return nil, &HTTPError{
				Status:     resp.Status,
				StatusCode: resp.StatusCode,
				URL:        resp.Request.URL.String(),
				RateLimit:  rl,
			}
		}
		e.StatusCode = resp.StatusCode
		e.RateLimit = rl
		return nil, e
	}

	respBody := struct {
		*TweetRaw
		Meta *TweetSearchMeta `json:"meta"`
	}{}

	if err := decoder.Decode(&respBody); err != nil {
		return nil, &ResponseDecodeError{
			Name:      "tweet search",
			Err:       err,
			RateLimit: rl,
		}
	}

	return &TweetSearchResponse{
		Raw:       respBody.TweetRaw,
		Meta:      respBody.Meta,
		RateLimit: rl,
	}, nil
}

// TweetSearchStreamAddRule will create on or more rules for search sampling.  Set dry run to true to validate the rules before commit
func (c *Client) TweetSearchStreamAddRule(ctx context.Context, rules []TweetSearchStreamRule, dryRun bool) (*TweetSearchStreamAddRuleResponse, error) {
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
	if dryRun {
		q := req.URL.Query()
		q.Add("dry_run", "true")
		req.URL.RawQuery = q.Encode()
	}

	resp, err := c.Client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("tweet search stream add rule http response %w", err)
	}
	defer resp.Body.Close()

	decoder := json.NewDecoder(resp.Body)

	rl := rateFromHeader(resp.Header)

	if resp.StatusCode != http.StatusCreated {
		e := &ErrorResponse{}
		if err := decoder.Decode(e); err != nil {
			return nil, &HTTPError{
				Status:     resp.Status,
				StatusCode: resp.StatusCode,
				URL:        resp.Request.URL.String(),
				RateLimit:  rl,
			}
		}
		e.StatusCode = resp.StatusCode
		e.RateLimit = rl
		return nil, e
	}

	ruleResponse := &TweetSearchStreamAddRuleResponse{}
	if err := decoder.Decode(ruleResponse); err != nil {
		return nil, &ResponseDecodeError{
			Name:      "tweet search stream add rule",
			Err:       err,
			RateLimit: rl,
		}
	}
	ruleResponse.RateLimit = rl
	return ruleResponse, nil
}

// TweetSearchStreamDeleteRuleByID will delete one or more rules for search sampling using the rule ids. Set dry run to true to validate the rules before commit
func (c *Client) TweetSearchStreamDeleteRuleByID(ctx context.Context, ruleIDs []TweetSearchStreamRuleID, dryRun bool) (*TweetSearchStreamDeleteRuleResponse, error) {
	if len(ruleIDs) == 0 {
		return nil, fmt.Errorf("tweet search stream delete rule: rule ids are required: %w", ErrParameter)
	}
	type ids struct {
		IDs tweetSearchStreamRuleIDs `json:"ids"`
	}
	deleteIDs := ids{
		IDs: tweetSearchStreamRuleIDs(ruleIDs),
	}
	if err := deleteIDs.IDs.validate(); err != nil {
		return nil, err
	}
	type requestBody struct {
		Delete ids `json:"delete"`
	}
	body := requestBody{
		Delete: deleteIDs,
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
	if dryRun {
		q := req.URL.Query()
		q.Add("dry_run", "true")
		req.URL.RawQuery = q.Encode()
	}

	resp, err := c.Client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("tweet search stream delete rule http response %w", err)
	}
	defer resp.Body.Close()

	decoder := json.NewDecoder(resp.Body)

	rl := rateFromHeader(resp.Header)

	if resp.StatusCode != http.StatusOK {
		e := &ErrorResponse{}
		if err := decoder.Decode(e); err != nil {
			return nil, &HTTPError{
				Status:     resp.Status,
				StatusCode: resp.StatusCode,
				URL:        resp.Request.URL.String(),
				RateLimit:  rl,
			}
		}
		e.StatusCode = resp.StatusCode
		e.RateLimit = rl
		return nil, e
	}

	ruleResponse := &TweetSearchStreamDeleteRuleResponse{}
	if err := decoder.Decode(ruleResponse); err != nil {
		return nil, &ResponseDecodeError{
			Name:      "tweet search stream delete rule ny id",
			Err:       err,
			RateLimit: rl,
		}
	}
	ruleResponse.RateLimit = rl
	return ruleResponse, nil
}

// TweetSearchStreamDeleteRuleByValue will delete one or more rules for search sampling using the rule values. Set dry run to true to validate the rules before commit
func (c *Client) TweetSearchStreamDeleteRuleByValue(ctx context.Context, ruleValues []string, dryRun bool) (*TweetSearchStreamDeleteRuleResponse, error) {
	if len(ruleValues) == 0 {
		return nil, fmt.Errorf("tweet search stream delete rule: rule values are required: %w", ErrParameter)
	}
	type values struct {
		Values []string `json:"values"`
	}
	type requestBody struct {
		Delete values `json:"delete"`
	}
	body := requestBody{
		Delete: values{
			Values: ruleValues,
		},
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
	if dryRun {
		q := req.URL.Query()
		q.Add("dry_run", "true")
		req.URL.RawQuery = q.Encode()
	}

	resp, err := c.Client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("tweet search stream delete rule http response %w", err)
	}
	defer resp.Body.Close()

	decoder := json.NewDecoder(resp.Body)

	rl := rateFromHeader(resp.Header)

	if resp.StatusCode != http.StatusOK {
		e := &ErrorResponse{}
		if err := decoder.Decode(e); err != nil {
			return nil, &HTTPError{
				Status:     resp.Status,
				StatusCode: resp.StatusCode,
				URL:        resp.Request.URL.String(),
				RateLimit:  rl,
			}
		}
		e.StatusCode = resp.StatusCode
		e.RateLimit = rl
		return nil, e
	}

	ruleResponse := &TweetSearchStreamDeleteRuleResponse{}
	if err := decoder.Decode(ruleResponse); err != nil {
		return nil, &ResponseDecodeError{
			Name:      "tweet search stream delete rule by value",
			Err:       err,
			RateLimit: rl,
		}
	}
	ruleResponse.RateLimit = rl
	return ruleResponse, nil
}

// TweetSearchStreamRules will return a list of rules active on the streaming endpoint
func (c *Client) TweetSearchStreamRules(ctx context.Context, ruleIDs []TweetSearchStreamRuleID) (*TweetSearchStreamRulesResponse, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, tweetSearchStreamRulesEndpoint.url(c.Host), nil)
	if err != nil {
		return nil, fmt.Errorf("tweet search stream rules http request %w", err)
	}
	req.Header.Add("Accept", "application/json")
	c.Authorizer.Add(req)
	if len(ruleIDs) > 0 {
		ruleArr := tweetSearchStreamRuleIDs(ruleIDs)
		if err := ruleArr.validate(); err != nil {
			return nil, err
		}
		q := req.URL.Query()
		q.Add("ids", strings.Join(ruleArr.toStringArray(), ","))
		req.URL.RawQuery = q.Encode()
	}

	resp, err := c.Client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("tweet search stream rules http response %w", err)
	}
	defer resp.Body.Close()

	decoder := json.NewDecoder(resp.Body)

	rl := rateFromHeader(resp.Header)

	if resp.StatusCode != http.StatusOK {
		e := &ErrorResponse{}
		if err := decoder.Decode(e); err != nil {
			return nil, &HTTPError{
				Status:     resp.Status,
				StatusCode: resp.StatusCode,
				URL:        resp.Request.URL.String(),
				RateLimit:  rl,
			}
		}
		e.StatusCode = resp.StatusCode
		e.RateLimit = rl
		return nil, e
	}

	ruleResponse := &TweetSearchStreamRulesResponse{}
	if err := decoder.Decode(ruleResponse); err != nil {
		return nil, &ResponseDecodeError{
			Name:      "tweet search stream rules",
			Err:       err,
			RateLimit: rl,
		}
	}
	ruleResponse.RateLimit = rl
	return ruleResponse, nil
}

// TweetSearchStream will stream in real-time based on a specific set of filter rules
func (c *Client) TweetSearchStream(ctx context.Context, opts TweetSearchStreamOpts) (*TweetStream, error) {
	switch {
	case opts.BackfillMinutes == 0:
	case opts.BackfillMinutes > sampleStreamMaxBackOffMin:
		return nil, fmt.Errorf("tweet search stream: a max back off minutes [%d] is [current: %d]: %w", sampleStreamMaxBackOffMin, opts.BackfillMinutes, ErrParameter)
	default:
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, tweetSearchStreamEndpoint.url(c.Host), nil)
	if err != nil {
		return nil, fmt.Errorf("tweet search stream request: %w", err)
	}
	req.Header.Add("Accept", "application/json")
	c.Authorizer.Add(req)
	opts.addQuery(req)

	resp, err := c.Client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("tweet search stream response: %w", err)
	}

	rl := rateFromHeader(resp.Header)

	if resp.StatusCode != http.StatusOK {
		defer resp.Body.Close()
		e := &ErrorResponse{}
		if err := json.NewDecoder(resp.Body).Decode(e); err != nil {
			return nil, &HTTPError{
				Status:     resp.Status,
				StatusCode: resp.StatusCode,
				URL:        resp.Request.URL.String(),
				RateLimit:  rl,
			}
		}
		e.StatusCode = resp.StatusCode
		e.RateLimit = rl
		return nil, e
	}

	stream := StartTweetStream(resp.Body)
	stream.RateLimit = rl
	return stream, nil
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

	rl := rateFromHeader(resp.Header)

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
				RateLimit:  rl,
			}
		}
		e.StatusCode = resp.StatusCode
		e.RateLimit = rl
		return nil, e
	}

	recentCounts := &TweetRecentCountsResponse{
		TweetCounts: []*TweetCount{},
		Meta:        &TweetRecentCountsMeta{},
	}

	if err := json.Unmarshal(respBytes, recentCounts); err != nil {
		return nil, &ResponseDecodeError{
			Name:      "tweet recent counts",
			Err:       err,
			RateLimit: rl,
		}
	}
	recentCounts.RateLimit = rl
	return recentCounts, nil
}

// TweetAllCounts receive a count of Tweets that match a query
func (c *Client) TweetAllCounts(ctx context.Context, query string, opts TweetAllCountsOpts) (*TweetAllCountsResponse, error) {
	switch {
	case len(query) == 0:
		return nil, fmt.Errorf("tweet all counts: a query is required: %w", ErrParameter)
	case len(query) > tweetAllCountsQueryLength:
		return nil, fmt.Errorf("tweet all counts: the query over the length (%d): %w", tweetAllCountsQueryLength, ErrParameter)
	default:
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, tweetAllCountsEndpoint.url(c.Host), nil)
	if err != nil {
		return nil, fmt.Errorf("tweet all counts request: %w", err)
	}
	req.Header.Add("Accept", "application/json")
	c.Authorizer.Add(req)
	opts.addQuery(req)
	q := req.URL.Query()
	q.Add("query", query)
	req.URL.RawQuery = q.Encode()

	resp, err := c.Client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("tweet all counts response: %w", err)
	}
	defer resp.Body.Close()

	rl := rateFromHeader(resp.Header)

	decoder := json.NewDecoder(resp.Body)

	if resp.StatusCode != http.StatusOK {
		e := &ErrorResponse{}
		if err := decoder.Decode(e); err != nil {
			return nil, &HTTPError{
				Status:     resp.Status,
				StatusCode: resp.StatusCode,
				URL:        resp.Request.URL.String(),
				RateLimit:  rl,
			}
		}
		e.StatusCode = resp.StatusCode
		e.RateLimit = rl
		return nil, e
	}

	allCounts := &TweetAllCountsResponse{
		TweetCounts: []*TweetCount{},
		Meta:        &TweetAllCountsMeta{},
	}

	if err := decoder.Decode(allCounts); err != nil {
		return nil, &ResponseDecodeError{
			Name:      "tweet all counts",
			Err:       err,
			RateLimit: rl,
		}
	}
	allCounts.RateLimit = rl
	return allCounts, nil
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

	rl := rateFromHeader(resp.Header)

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
				RateLimit:  rl,
			}
		}
		e.StatusCode = resp.StatusCode
		e.RateLimit = rl
		return nil, e
	}

	followingLookup := &UserFollowingLookupResponse{
		Raw:  &UserRaw{},
		Meta: &UserFollowingMeta{},
	}

	if err := json.Unmarshal(respBytes, followingLookup.Raw); err != nil {
		return nil, &ResponseDecodeError{
			Name:      "user following lookup",
			Err:       err,
			RateLimit: rl,
		}
	}

	if err := json.Unmarshal(respBytes, followingLookup); err != nil {
		return nil, &ResponseDecodeError{
			Name:      "user following lookup",
			Err:       err,
			RateLimit: rl,
		}
	}
	followingLookup.RateLimit = rl
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

	rl := rateFromHeader(resp.Header)

	if resp.StatusCode != http.StatusOK {
		e := &ErrorResponse{}
		if err := decoder.Decode(e); err != nil {
			return nil, &HTTPError{
				Status:     resp.Status,
				StatusCode: resp.StatusCode,
				URL:        resp.Request.URL.String(),
				RateLimit:  rl,
			}
		}
		e.StatusCode = resp.StatusCode
		e.RateLimit = rl
		return nil, e
	}

	raw := &UserFollowsResponse{}
	if err := decoder.Decode(raw); err != nil {
		return nil, &ResponseDecodeError{
			Name:      "user follows",
			Err:       err,
			RateLimit: rl,
		}
	}
	raw.RateLimit = rl
	return raw, nil
}

// DeleteUserFollows allows a user ID to unfollow another user
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

	rl := rateFromHeader(resp.Header)

	if resp.StatusCode != http.StatusOK {
		e := &ErrorResponse{}
		if err := decoder.Decode(e); err != nil {
			return nil, &HTTPError{
				Status:     resp.Status,
				StatusCode: resp.StatusCode,
				URL:        resp.Request.URL.String(),
				RateLimit:  rl,
			}
		}
		e.StatusCode = resp.StatusCode
		e.RateLimit = rl
		return nil, e
	}

	raw := &UserDeleteFollowsResponse{}
	if err := decoder.Decode(raw); err != nil {
		return nil, &ResponseDecodeError{
			Name:      "delete user follows",
			Err:       err,
			RateLimit: rl,
		}
	}
	raw.RateLimit = rl
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

	rl := rateFromHeader(resp.Header)

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
				RateLimit:  rl,
			}
		}
		e.StatusCode = resp.StatusCode
		e.RateLimit = rl
		return nil, e
	}

	followersLookup := &UserFollowersLookupResponse{
		Raw:  &UserRaw{},
		Meta: &UserFollowershMeta{},
	}

	if err := json.Unmarshal(respBytes, followersLookup.Raw); err != nil {
		return nil, &ResponseDecodeError{
			Name:      "user followers lookup",
			Err:       err,
			RateLimit: rl,
		}
	}

	if err := json.Unmarshal(respBytes, followersLookup); err != nil {
		return nil, &ResponseDecodeError{
			Name:      "user followers lookup",
			Err:       err,
			RateLimit: rl,
		}
	}
	followersLookup.RateLimit = rl
	return followersLookup, nil
}

// UserTweetTimeline will return the user tweet timeline
func (c *Client) UserTweetTimeline(ctx context.Context, userID string, opts UserTweetTimelineOpts) (*UserTweetTimelineResponse, error) {
	switch {
	case len(userID) == 0:
		return nil, fmt.Errorf("user tweet timeline: a query is required: %w", ErrParameter)
	case opts.MaxResults == 0:
	case opts.MaxResults < userTweetTimelineMinResults:
		return nil, fmt.Errorf("user tweet timeline: max results [%d] have a min[%d] %w", opts.MaxResults, userTweetTimelineMinResults, ErrParameter)
	case opts.MaxResults > userTweetTimelineMaxResults:
		return nil, fmt.Errorf("user tweet timeline: max results [%d] have a max[%d] %w", opts.MaxResults, userTweetTimelineMaxResults, ErrParameter)
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

	rl := rateFromHeader(resp.Header)

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
				RateLimit:  rl,
			}
		}
		e.StatusCode = resp.StatusCode
		e.RateLimit = rl
		return nil, e
	}

	timeline := &UserTweetTimelineResponse{
		Raw:  &TweetRaw{},
		Meta: &UserTimelineMeta{},
	}

	if err := json.Unmarshal(respBytes, timeline.Raw); err != nil {
		return nil, &ResponseDecodeError{
			Name:      "user tweet timeline",
			Err:       err,
			RateLimit: rl,
		}
	}

	if err := json.Unmarshal(respBytes, timeline); err != nil {
		return nil, &ResponseDecodeError{
			Name:      "user tweet timeline",
			Err:       err,
			RateLimit: rl,
		}
	}
	timeline.RateLimit = rl
	return timeline, nil
}

// UserMentionTimeline will return the user's mentions timeline
func (c *Client) UserMentionTimeline(ctx context.Context, userID string, opts UserMentionTimelineOpts) (*UserMentionTimelineResponse, error) {
	switch {
	case len(userID) == 0:
		return nil, fmt.Errorf("user mention timeline: a query is required: %w", ErrParameter)
	case opts.MaxResults == 0:
	case opts.MaxResults < userMentionTimelineMinResults:
		return nil, fmt.Errorf("user mention timeline: max results [%d] have a min[%d] %w", opts.MaxResults, userMentionTimelineMinResults, ErrParameter)
	case opts.MaxResults > userMentionTimelineMaxResults:
		return nil, fmt.Errorf("user mention timeline: max results [%d] have a max[%d] %w", opts.MaxResults, userMentionTimelineMaxResults, ErrParameter)
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

	rl := rateFromHeader(resp.Header)

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
				RateLimit:  rl,
			}
		}
		e.StatusCode = resp.StatusCode
		e.RateLimit = rl
		return nil, e
	}

	timeline := &UserMentionTimelineResponse{
		Raw:  &TweetRaw{},
		Meta: &UserTimelineMeta{},
	}

	if err := json.Unmarshal(respBytes, timeline.Raw); err != nil {
		return nil, &ResponseDecodeError{
			Name:      "user mention timeline",
			Err:       err,
			RateLimit: rl,
		}
	}

	if err := json.Unmarshal(respBytes, timeline); err != nil {
		return nil, &ResponseDecodeError{
			Name:      "user mention timeline",
			Err:       err,
			RateLimit: rl,
		}
	}
	timeline.RateLimit = rl
	return timeline, nil
}

// UserTweetReverseChronologicalTimeline allows you to retrieve a collection of the most recent Tweets and Retweets posted by you and users you follow.
// This endpoint returns up to the last 3200 Tweets.
func (c *Client) UserTweetReverseChronologicalTimeline(ctx context.Context, userID string, opts UserTweetReverseChronologicalTimelineOpts) (*UserTweetReverseChronologicalTimelineResponse, error) {
	switch {
	case len(userID) == 0:
		return nil, fmt.Errorf("user tweet reverse chronological timeline: a query is required: %w", ErrParameter)
	case opts.MaxResults == 0:
	case opts.MaxResults < userTweetReverseChronologicalTimelineMinResults:
		return nil, fmt.Errorf("user tweet reverse chronological timeline: max results [%d] have a min[%d] %w", opts.MaxResults, userTweetReverseChronologicalTimelineMinResults, ErrParameter)
	case opts.MaxResults > userTweetReverseChronologicalTimelineMaxResults:
		return nil, fmt.Errorf("user tweet reverse chronological timeline: max results [%d] have a max[%d] %w", opts.MaxResults, userTweetReverseChronologicalTimelineMaxResults, ErrParameter)
	default:
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, userTweetReverseChronologicalTimelineEndpoint.urlID(c.Host, userID), nil)
	if err != nil {
		return nil, fmt.Errorf("user tweet reverse chronological timeline request: %w", err)
	}
	req.Header.Add("Accept", "application/json")
	c.Authorizer.Add(req)
	opts.addQuery(req)

	resp, err := c.Client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("user tweet reverse chronological timeline response: %w", err)
	}
	defer resp.Body.Close()

	rl := rateFromHeader(resp.Header)

	decoder := json.NewDecoder(resp.Body)

	if resp.StatusCode != http.StatusOK {
		e := &ErrorResponse{}
		if err := decoder.Decode(e); err != nil {
			return nil, &HTTPError{
				Status:     resp.Status,
				StatusCode: resp.StatusCode,
				URL:        resp.Request.URL.String(),
				RateLimit:  rl,
			}
		}
		e.StatusCode = resp.StatusCode
		e.RateLimit = rl
		return nil, e
	}

	timeline := struct {
		TweetRaw
		Meta UserReverseChronologicalTimelineMeta `json:"meta"`
	}{}

	if err := decoder.Decode(&timeline); err != nil {
		return nil, &ResponseDecodeError{
			Name:      "user tweet reverse chronological timeline",
			Err:       err,
			RateLimit: rl,
		}
	}

	return &UserTweetReverseChronologicalTimelineResponse{
		Raw:       &timeline.TweetRaw,
		Meta:      &timeline.Meta,
		RateLimit: rl,
	}, nil
}

// TweetHideReplies will hide the replies for a given tweet
func (c Client) TweetHideReplies(ctx context.Context, id string, hide bool) (*TweetHideReplyResponse, error) {
	if len(id) == 0 {
		return nil, fmt.Errorf("tweet hide replies: id must be present %w", ErrParameter)
	}
	type body struct {
		Hidden bool `json:"hidden"`
	}
	rb := body{
		Hidden: hide,
	}
	enc, err := json.Marshal(rb)
	if err != nil {
		return nil, fmt.Errorf("tweet hide replies: request body marshal %w", err)
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPut, tweetHideRepliesEndpoint.urlID(c.Host, id), bytes.NewReader(enc))
	if err != nil {
		return nil, fmt.Errorf("tweet hide replies request: %w", err)
	}
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")
	c.Authorizer.Add(req)

	resp, err := c.Client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("tweet hide replies response: %w", err)
	}
	defer resp.Body.Close()

	rl := rateFromHeader(resp.Header)

	decoder := json.NewDecoder(resp.Body)

	if resp.StatusCode != http.StatusOK {
		errResp := &ErrorResponse{}
		if err := decoder.Decode(errResp); err != nil {
			return nil, &HTTPError{
				Status:     resp.Status,
				StatusCode: resp.StatusCode,
				URL:        resp.Request.URL.String(),
				RateLimit:  rl,
			}
		}
		errResp.StatusCode = resp.StatusCode
		errResp.RateLimit = rl
		return nil, errResp
	}

	rd := &TweetHideReplyResponse{}
	if err := decoder.Decode(rd); err != nil {
		return nil, &ResponseDecodeError{
			Name:      "tweet hide replies",
			Err:       err,
			RateLimit: rl,
		}
	}
	rd.RateLimit = rl
	return rd, nil
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

	rl := rateFromHeader(resp.Header)

	if resp.StatusCode != http.StatusOK {
		e := &ErrorResponse{}
		if err := decoder.Decode(e); err != nil {
			return nil, &HTTPError{
				Status:     resp.Status,
				StatusCode: resp.StatusCode,
				URL:        resp.Request.URL.String(),
				RateLimit:  rl,
			}
		}
		e.StatusCode = resp.StatusCode
		e.RateLimit = rl
		return nil, e
	}

	raw := &UserRetweetResponse{}
	if err := decoder.Decode(raw); err != nil {
		return nil, &ResponseDecodeError{
			Name:      "user retweet",
			Err:       err,
			RateLimit: rl,
		}
	}
	raw.RateLimit = rl
	return raw, nil
}

// DeleteUserRetweet will delete a retweet from a user
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

	rl := rateFromHeader(resp.Header)

	if resp.StatusCode != http.StatusOK {
		e := &ErrorResponse{}
		if err := decoder.Decode(e); err != nil {
			return nil, &HTTPError{
				Status:     resp.Status,
				StatusCode: resp.StatusCode,
				URL:        resp.Request.URL.String(),
				RateLimit:  rl,
			}
		}
		e.StatusCode = resp.StatusCode
		e.RateLimit = rl
		return nil, e
	}

	raw := &DeleteUserRetweetResponse{}
	if err := decoder.Decode(raw); err != nil {
		return nil, &ResponseDecodeError{
			Name:      "delete user retweet",
			Err:       err,
			RateLimit: rl,
		}
	}
	raw.RateLimit = rl
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

	rl := rateFromHeader(resp.Header)

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
				RateLimit:  rl,
			}
		}
		e.StatusCode = resp.StatusCode
		e.RateLimit = rl
		return nil, e
	}

	blockedLookup := &UserBlocksLookupResponse{
		Raw:  &UserRaw{},
		Meta: &UserBlocksLookupMeta{},
	}

	if err := json.Unmarshal(respBytes, blockedLookup.Raw); err != nil {
		return nil, &ResponseDecodeError{
			Name:      "user blocked lookup",
			Err:       err,
			RateLimit: rl,
		}
	}

	if err := json.Unmarshal(respBytes, blockedLookup); err != nil {
		return nil, &ResponseDecodeError{
			Name:      "user blocked lookup",
			Err:       err,
			RateLimit: rl,
		}
	}
	blockedLookup.RateLimit = rl
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

	rl := rateFromHeader(resp.Header)

	if resp.StatusCode != http.StatusOK {
		e := &ErrorResponse{}
		if err := decoder.Decode(e); err != nil {
			return nil, &HTTPError{
				Status:     resp.Status,
				StatusCode: resp.StatusCode,
				URL:        resp.Request.URL.String(),
				RateLimit:  rl,
			}
		}
		e.StatusCode = resp.StatusCode
		e.RateLimit = rl
		return nil, e
	}

	raw := &UserBlocksResponse{}
	if err := decoder.Decode(raw); err != nil {
		return nil, &ResponseDecodeError{
			Name:      "user blocks",
			Err:       err,
			RateLimit: rl,
		}
	}
	raw.RateLimit = rl
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

	rl := rateFromHeader(resp.Header)

	if resp.StatusCode != http.StatusOK {
		e := &ErrorResponse{}
		if err := decoder.Decode(e); err != nil {
			return nil, &HTTPError{
				Status:     resp.Status,
				StatusCode: resp.StatusCode,
				URL:        resp.Request.URL.String(),
				RateLimit:  rl,
			}
		}
		e.StatusCode = resp.StatusCode
		e.RateLimit = rl
		return nil, e
	}

	raw := &UserDeleteBlocksResponse{}
	if err := decoder.Decode(raw); err != nil {
		return nil, &ResponseDecodeError{
			Name:      "delete user blocks",
			Err:       err,
			RateLimit: rl,
		}
	}
	raw.RateLimit = rl
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

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, userMutesEndpoint.urlID(c.Host, userID), nil)
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

	rl := rateFromHeader(resp.Header)

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
				RateLimit:  rl,
			}
		}
		e.StatusCode = resp.StatusCode
		e.RateLimit = rl
		return nil, e
	}

	mutedLookup := &UserMutesLookupResponse{
		Raw:  &UserRaw{},
		Meta: &UserMutesLookupMeta{},
	}

	if err := json.Unmarshal(respBytes, mutedLookup.Raw); err != nil {
		return nil, &ResponseDecodeError{
			Name:      "user muted lookup",
			Err:       err,
			RateLimit: rl,
		}
	}

	if err := json.Unmarshal(respBytes, mutedLookup); err != nil {
		return nil, &ResponseDecodeError{
			Name:      "user muted lookup",
			Err:       err,
			RateLimit: rl,
		}
	}
	mutedLookup.RateLimit = rl
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
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, userMutesEndpoint.urlID(c.Host, userID), bytes.NewReader(enc))
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

	rl := rateFromHeader(resp.Header)

	if resp.StatusCode != http.StatusOK {
		e := &ErrorResponse{}
		if err := decoder.Decode(e); err != nil {
			return nil, &HTTPError{
				Status:     resp.Status,
				StatusCode: resp.StatusCode,
				URL:        resp.Request.URL.String(),
				RateLimit:  rl,
			}
		}
		e.StatusCode = resp.StatusCode
		e.RateLimit = rl
		return nil, e
	}

	raw := &UserMutesResponse{}
	if err := decoder.Decode(raw); err != nil {
		return nil, &ResponseDecodeError{
			Name:      "user mutes",
			Err:       err,
			RateLimit: rl,
		}
	}
	raw.RateLimit = rl
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

	ep := userMutesEndpoint.urlID(c.Host, userID) + "/" + targetUserID
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

	rl := rateFromHeader(resp.Header)

	if resp.StatusCode != http.StatusOK {
		e := &ErrorResponse{}
		if err := decoder.Decode(e); err != nil {
			return nil, &HTTPError{
				Status:     resp.Status,
				StatusCode: resp.StatusCode,
				URL:        resp.Request.URL.String(),
				RateLimit:  rl,
			}
		}
		e.StatusCode = resp.StatusCode
		e.RateLimit = rl
		return nil, e
	}

	raw := &UserDeleteMutesResponse{}
	if err := decoder.Decode(raw); err != nil {
		return nil, &ResponseDecodeError{
			Name:      "user delete mutes",
			Err:       err,
			RateLimit: rl,
		}
	}
	raw.RateLimit = rl
	return raw, nil
}

// TweetLikesLookup gets information about a tweet's liking users.  The response will have at most 100 users who liked the tweet
func (c *Client) TweetLikesLookup(ctx context.Context, tweetID string, opts TweetLikesLookupOpts) (*TweetLikesLookupResponse, error) {
	switch {
	case len(tweetID) == 0:
		return nil, fmt.Errorf("user tweet likes lookup: an id is required: %w", ErrParameter)
	case opts.MaxResults == 0:
	case opts.MaxResults < likesMinResults:
		return nil, fmt.Errorf("tweet tweet likes lookup: a min results [%d] is required [current: %d]: %w", likesMinResults, opts.MaxResults, ErrParameter)
	case opts.MaxResults > likesMaxResults:
		return nil, fmt.Errorf("tweet tweet likes lookup: a max results [%d] is required [current: %d]: %w", likesMaxResults, opts.MaxResults, ErrParameter)
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

	rl := rateFromHeader(resp.Header)

	if resp.StatusCode != http.StatusOK {
		e := &ErrorResponse{}
		if err := decoder.Decode(e); err != nil {
			return nil, &HTTPError{
				Status:     resp.Status,
				StatusCode: resp.StatusCode,
				URL:        resp.Request.URL.String(),
				RateLimit:  rl,
			}
		}
		e.StatusCode = resp.StatusCode
		e.RateLimit = rl
		return nil, e
	}

	respBody := struct {
		*UserRaw
		Meta *TweetLikesMeta `json:"meta"`
	}{}

	if err := decoder.Decode(&respBody); err != nil {
		return nil, &ResponseDecodeError{
			Name:      "tweet likes lookup",
			Err:       err,
			RateLimit: rl,
		}
	}

	return &TweetLikesLookupResponse{
		Raw:       respBody.UserRaw,
		Meta:      respBody.Meta,
		RateLimit: rl,
	}, nil
}

// UserLikesLookup gets information about a user's liked tweets.
func (c *Client) UserLikesLookup(ctx context.Context, userID string, opts UserLikesLookupOpts) (*UserLikesLookupResponse, error) {
	switch {
	case len(userID) == 0:
		return nil, fmt.Errorf("tweet user likes lookup: an id is required: %w", ErrParameter)
	case opts.MaxResults == 0:
	case opts.MaxResults < likesMinResults:
		return nil, fmt.Errorf("tweet user likes lookup: a min results [%d] is required [current: %d]: %w", likesMinResults, opts.MaxResults, ErrParameter)
	case opts.MaxResults > likesMaxResults:
		return nil, fmt.Errorf("tweet user likes lookup: a max results [%d] is required [current: %d]: %w", likesMaxResults, opts.MaxResults, ErrParameter)
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

	rl := rateFromHeader(resp.Header)

	if resp.StatusCode != http.StatusOK {
		e := &ErrorResponse{}
		if err := decoder.Decode(e); err != nil {
			return nil, &HTTPError{
				Status:     resp.Status,
				StatusCode: resp.StatusCode,
				URL:        resp.Request.URL.String(),
				RateLimit:  rl,
			}
		}
		e.StatusCode = resp.StatusCode
		e.RateLimit = rl
		return nil, e
	}

	respBody := struct {
		*TweetRaw
		Meta *UserLikesMeta `json:"meta"`
	}{}

	if err := decoder.Decode(&respBody); err != nil {
		return nil, &ResponseDecodeError{
			Name:      "user likes lookup",
			Err:       err,
			RateLimit: rl,
		}
	}
	return &UserLikesLookupResponse{
		Raw:       respBody.TweetRaw,
		Meta:      respBody.Meta,
		RateLimit: rl,
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

	rl := rateFromHeader(resp.Header)

	if resp.StatusCode != http.StatusOK {
		e := &ErrorResponse{}
		if err := decoder.Decode(e); err != nil {
			return nil, &HTTPError{
				Status:     resp.Status,
				StatusCode: resp.StatusCode,
				URL:        resp.Request.URL.String(),
				RateLimit:  rl,
			}
		}
		e.StatusCode = resp.StatusCode
		e.RateLimit = rl
		return nil, e
	}

	raw := &UserLikesResponse{}
	if err := decoder.Decode(raw); err != nil {
		return nil, &ResponseDecodeError{
			Name:      "user likes",
			Err:       err,
			RateLimit: rl,
		}
	}
	raw.RateLimit = rl
	return raw, nil
}

// DeleteUserLikes will unlike the targeted tweet
func (c *Client) DeleteUserLikes(ctx context.Context, userID, tweetID string) (*DeleteUserLikesResponse, error) {
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

	rl := rateFromHeader(resp.Header)

	if resp.StatusCode != http.StatusOK {
		e := &ErrorResponse{}
		if err := decoder.Decode(e); err != nil {
			return nil, &HTTPError{
				Status:     resp.Status,
				StatusCode: resp.StatusCode,
				URL:        resp.Request.URL.String(),
				RateLimit:  rl,
			}
		}
		e.StatusCode = resp.StatusCode
		e.RateLimit = rl
		return nil, e
	}

	raw := &DeleteUserLikesResponse{}
	if err := decoder.Decode(raw); err != nil {
		return nil, &ResponseDecodeError{
			Name:      "delete user likes",
			Err:       err,
			RateLimit: rl,
		}
	}
	raw.RateLimit = rl
	return raw, nil
}

// TweetSampleStream will return a streamer for streaming 1% of all tweets real-time
func (c *Client) TweetSampleStream(ctx context.Context, opts TweetSampleStreamOpts) (*TweetStream, error) {
	switch {
	case opts.BackfillMinutes == 0:
	case opts.BackfillMinutes > sampleStreamMaxBackOffMin:
		return nil, fmt.Errorf("tweet sample stream: a max back off minutes [%d] is [current: %d]: %w", sampleStreamMaxBackOffMin, opts.BackfillMinutes, ErrParameter)
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

	rl := rateFromHeader(resp.Header)

	if resp.StatusCode != http.StatusOK {
		defer resp.Body.Close()
		e := &ErrorResponse{}
		if err := json.NewDecoder(resp.Body).Decode(e); err != nil {
			return nil, &HTTPError{
				Status:     resp.Status,
				StatusCode: resp.StatusCode,
				URL:        resp.Request.URL.String(),
				RateLimit:  rl,
			}
		}
		e.StatusCode = resp.StatusCode
		e.RateLimit = rl
		return nil, e
	}

	stream := StartTweetStream(resp.Body)
	stream.RateLimit = rl
	return stream, nil
}

// ListLookup returns the details of a specified list
func (c *Client) ListLookup(ctx context.Context, listID string, opts ListLookupOpts) (*ListLookupResponse, error) {
	switch {
	case len(listID) == 0:
		return nil, fmt.Errorf("list lookup: an id is required: %w", ErrParameter)
	default:
	}

	ep := listLookupEndpoint.urlID(c.Host, listID)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, ep, nil)
	if err != nil {
		return nil, fmt.Errorf("list lookup request: %w", err)
	}
	req.Header.Add("Accept", "application/json")
	c.Authorizer.Add(req)
	opts.addQuery(req)

	resp, err := c.Client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("list lookup response: %w", err)
	}
	defer resp.Body.Close()

	decoder := json.NewDecoder(resp.Body)

	rl := rateFromHeader(resp.Header)

	if resp.StatusCode != http.StatusOK {
		e := &ErrorResponse{}
		if err := decoder.Decode(e); err != nil {
			return nil, &HTTPError{
				Status:     resp.Status,
				StatusCode: resp.StatusCode,
				URL:        resp.Request.URL.String(),
				RateLimit:  rl,
			}
		}
		e.StatusCode = resp.StatusCode
		e.RateLimit = rl
		return nil, e
	}

	respBody := struct {
		*ListRaw
	}{}

	if err := decoder.Decode(&respBody); err != nil {
		return nil, &ResponseDecodeError{
			Name:      "list lookup",
			Err:       err,
			RateLimit: rl,
		}
	}

	return &ListLookupResponse{
		Raw:       respBody.ListRaw,
		RateLimit: rl,
	}, nil
}

// UserListLookup returns all lists owned by the specified user
func (c *Client) UserListLookup(ctx context.Context, userID string, opts UserListLookupOpts) (*UserListLookupResponse, error) {
	switch {
	case len(userID) == 0:
		return nil, fmt.Errorf("user list lookup: an id is required: %w", ErrParameter)
	case opts.MaxResults == 0:
	case opts.MaxResults > userListMaxResults:
		return nil, fmt.Errorf("user list lookup: max results [%d] is greater than max [%d]: %w", opts.MaxResults, userListMaxResults, ErrParameter)
	default:
	}

	ep := userListLookupEndpoint.urlID(c.Host, userID)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, ep, nil)
	if err != nil {
		return nil, fmt.Errorf("user list lookup request: %w", err)
	}
	req.Header.Add("Accept", "application/json")
	c.Authorizer.Add(req)
	opts.addQuery(req)

	resp, err := c.Client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("user list lookup response: %w", err)
	}
	defer resp.Body.Close()

	decoder := json.NewDecoder(resp.Body)

	rl := rateFromHeader(resp.Header)

	if resp.StatusCode != http.StatusOK {
		e := &ErrorResponse{}
		if err := decoder.Decode(e); err != nil {
			return nil, &HTTPError{
				Status:     resp.Status,
				StatusCode: resp.StatusCode,
				URL:        resp.Request.URL.String(),
				RateLimit:  rl,
			}
		}
		e.StatusCode = resp.StatusCode
		e.RateLimit = rl
		return nil, e
	}

	respBody := struct {
		*UserListRaw
		Meta *UserListLookupMeta `json:"meta"`
	}{}

	if err := decoder.Decode(&respBody); err != nil {
		return nil, &ResponseDecodeError{
			Name:      "user list lookup",
			Err:       err,
			RateLimit: rl,
		}
	}

	return &UserListLookupResponse{
		Raw:       respBody.UserListRaw,
		Meta:      respBody.Meta,
		RateLimit: rl,
	}, nil
}

// ListTweetLookup returns a list of tweets from the specified list
func (c *Client) ListTweetLookup(ctx context.Context, listID string, opts ListTweetLookupOpts) (*ListTweetLookupResponse, error) {
	switch {
	case len(listID) == 0:
		return nil, fmt.Errorf("list tweet lookup: an id is required: %w", ErrParameter)
	case opts.MaxResults == 0:
	case opts.MaxResults > listTweetMaxResults:
		return nil, fmt.Errorf("list tweet lookup: max results [%d] is greater than max [%d]: %w", opts.MaxResults, listTweetMaxResults, ErrParameter)
	default:
	}

	ep := listTweetLookupEndpoint.urlID(c.Host, listID)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, ep, nil)
	if err != nil {
		return nil, fmt.Errorf("list tweet lookup request: %w", err)
	}
	req.Header.Add("Accept", "application/json")
	c.Authorizer.Add(req)
	opts.addQuery(req)

	resp, err := c.Client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("list tweet lookup response: %w", err)
	}
	defer resp.Body.Close()

	decoder := json.NewDecoder(resp.Body)

	rl := rateFromHeader(resp.Header)

	if resp.StatusCode != http.StatusOK {
		e := &ErrorResponse{}
		if err := decoder.Decode(e); err != nil {
			return nil, &HTTPError{
				Status:     resp.Status,
				StatusCode: resp.StatusCode,
				URL:        resp.Request.URL.String(),
				RateLimit:  rl,
			}
		}
		e.StatusCode = resp.StatusCode
		e.RateLimit = rl
		return nil, e
	}

	respBody := struct {
		*TweetRaw
		Meta *ListTweetLookupMeta `json:"meta"`
	}{}

	if err := decoder.Decode(&respBody); err != nil {
		return nil, &ResponseDecodeError{
			Name:      "list tweet lookup",
			Err:       err,
			RateLimit: rl,
		}
	}

	return &ListTweetLookupResponse{
		Raw:       respBody.TweetRaw,
		Meta:      respBody.Meta,
		RateLimit: rl,
	}, nil
}

// CreateList enables the authenticated user to create a list
func (c *Client) CreateList(ctx context.Context, list ListMetaData) (*ListCreateResponse, error) {
	switch {
	case len(*list.Name) == 0:
		return nil, fmt.Errorf("create list: a name is required: %w", ErrParameter)
	default:
	}

	enc, err := json.Marshal(list)
	if err != nil {
		return nil, fmt.Errorf("create list: unable to encode json request %w", err)
	}

	ep := listCreateEndpoint.url(c.Host)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, ep, bytes.NewReader(enc))
	if err != nil {
		return nil, fmt.Errorf("create list request: %w", err)
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")
	c.Authorizer.Add(req)

	resp, err := c.Client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("create list response: %w", err)
	}
	defer resp.Body.Close()

	decoder := json.NewDecoder(resp.Body)

	rl := rateFromHeader(resp.Header)

	if resp.StatusCode != http.StatusCreated {
		e := &ErrorResponse{}
		if err := decoder.Decode(e); err != nil {
			return nil, &HTTPError{
				Status:     resp.Status,
				StatusCode: resp.StatusCode,
				URL:        resp.Request.URL.String(),
				RateLimit:  rl,
			}
		}
		e.StatusCode = resp.StatusCode
		e.RateLimit = rl
		return nil, e
	}

	respBody := &ListCreateResponse{}

	if err := decoder.Decode(respBody); err != nil {
		return nil, &ResponseDecodeError{
			Name:      "create list",
			Err:       err,
			RateLimit: rl,
		}
	}
	respBody.RateLimit = rl
	return respBody, nil
}

// UpdateList enables the authenticated user to update the meta data of a list
func (c *Client) UpdateList(ctx context.Context, listID string, update ListMetaData) (*ListUpdateResponse, error) {
	switch {
	case len(listID) == 0:
		return nil, fmt.Errorf("update list: an id is required: %w", ErrParameter)
	default:
	}

	enc, err := json.Marshal(update)
	if err != nil {
		return nil, fmt.Errorf("update list: unable to encode json request %w", err)
	}

	ep := listUpdateEndpoint.urlID(c.Host, listID)
	req, err := http.NewRequestWithContext(ctx, http.MethodPut, ep, bytes.NewReader(enc))
	if err != nil {
		return nil, fmt.Errorf("create list request: %w", err)
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")
	c.Authorizer.Add(req)

	resp, err := c.Client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("update list response: %w", err)
	}
	defer resp.Body.Close()

	decoder := json.NewDecoder(resp.Body)

	rl := rateFromHeader(resp.Header)

	if resp.StatusCode != http.StatusOK {
		e := &ErrorResponse{}
		if err := decoder.Decode(e); err != nil {
			return nil, &HTTPError{
				Status:     resp.Status,
				StatusCode: resp.StatusCode,
				URL:        resp.Request.URL.String(),
				RateLimit:  rl,
			}
		}
		e.StatusCode = resp.StatusCode
		e.RateLimit = rl
		return nil, e
	}

	respBody := &ListUpdateResponse{}

	if err := decoder.Decode(respBody); err != nil {
		return nil, &ResponseDecodeError{
			Name:      "update list",
			Err:       err,
			RateLimit: rl,
		}
	}
	respBody.RateLimit = rl
	return respBody, nil
}

// DeleteList enables the authenticated user to delete a list
func (c *Client) DeleteList(ctx context.Context, listID string) (*ListDeleteResponse, error) {
	switch {
	case len(listID) == 0:
		return nil, fmt.Errorf("delete list: an id is required: %w", ErrParameter)
	default:
	}

	ep := listDeleteEndpoint.urlID(c.Host, listID)
	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, ep, nil)
	if err != nil {
		return nil, fmt.Errorf("delete list request: %w", err)
	}
	req.Header.Add("Accept", "application/json")
	c.Authorizer.Add(req)

	resp, err := c.Client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("delete list response: %w", err)
	}
	defer resp.Body.Close()

	decoder := json.NewDecoder(resp.Body)

	rl := rateFromHeader(resp.Header)

	if resp.StatusCode != http.StatusOK {
		e := &ErrorResponse{}
		if err := decoder.Decode(e); err != nil {
			return nil, &HTTPError{
				Status:     resp.Status,
				StatusCode: resp.StatusCode,
				URL:        resp.Request.URL.String(),
				RateLimit:  rl,
			}
		}
		e.StatusCode = resp.StatusCode
		e.RateLimit = rl
		return nil, e
	}

	respBody := &ListDeleteResponse{}

	if err := decoder.Decode(respBody); err != nil {
		return nil, &ResponseDecodeError{
			Name:      "delete list",
			Err:       err,
			RateLimit: rl,
		}
	}
	respBody.RateLimit = rl
	return respBody, nil
}

// AddListMember enables the authenticated user to add a member to a list
func (c *Client) AddListMember(ctx context.Context, listID, userID string) (*ListAddMemberResponse, error) {
	switch {
	case len(listID) == 0:
		return nil, fmt.Errorf("add list member: a list id is required: %w", ErrParameter)
	case len(userID) == 0:
		return nil, fmt.Errorf("add list member: an user id is required: %w", ErrParameter)
	default:
	}

	reqBody := struct {
		UserID string `json:"user_id"`
	}{
		UserID: userID,
	}

	enc, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("add list member: unable to encode json request %w", err)
	}

	ep := listMemberEndpoint.urlID(c.Host, listID)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, ep, bytes.NewReader(enc))
	if err != nil {
		return nil, fmt.Errorf("add list member request: %w", err)
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")
	c.Authorizer.Add(req)

	resp, err := c.Client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("create list member response: %w", err)
	}
	defer resp.Body.Close()

	decoder := json.NewDecoder(resp.Body)

	rl := rateFromHeader(resp.Header)

	if resp.StatusCode != http.StatusOK {
		e := &ErrorResponse{}
		if err := decoder.Decode(e); err != nil {
			return nil, &HTTPError{
				Status:     resp.Status,
				StatusCode: resp.StatusCode,
				URL:        resp.Request.URL.String(),
				RateLimit:  rl,
			}
		}
		e.StatusCode = resp.StatusCode
		e.RateLimit = rl
		return nil, e
	}

	respBody := &ListAddMemberResponse{}

	if err := decoder.Decode(respBody); err != nil {
		return nil, &ResponseDecodeError{
			Name:      "add list member",
			Err:       err,
			RateLimit: rl,
		}
	}
	respBody.RateLimit = rl
	return respBody, nil
}

// RemoveListMember enables the authenticated user to remove a member to a list
func (c *Client) RemoveListMember(ctx context.Context, listID, userID string) (*ListRemoveMemberResponse, error) {
	switch {
	case len(listID) == 0:
		return nil, fmt.Errorf("remove list member: a list id is required: %w", ErrParameter)
	case len(userID) == 0:
		return nil, fmt.Errorf("remove list member: an user id is required: %w", ErrParameter)
	default:
	}

	ep := listMemberEndpoint.urlID(c.Host, listID) + "/" + userID
	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, ep, nil)
	if err != nil {
		return nil, fmt.Errorf("remove list member request: %w", err)
	}
	req.Header.Add("Accept", "application/json")
	c.Authorizer.Add(req)

	resp, err := c.Client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("remove list member response: %w", err)
	}
	defer resp.Body.Close()

	decoder := json.NewDecoder(resp.Body)

	rl := rateFromHeader(resp.Header)

	if resp.StatusCode != http.StatusOK {
		e := &ErrorResponse{}
		if err := decoder.Decode(e); err != nil {
			return nil, &HTTPError{
				Status:     resp.Status,
				StatusCode: resp.StatusCode,
				URL:        resp.Request.URL.String(),
				RateLimit:  rl,
			}
		}
		e.StatusCode = resp.StatusCode
		e.RateLimit = rl
		return nil, e
	}

	respBody := &ListRemoveMemberResponse{}

	if err := decoder.Decode(respBody); err != nil {
		return nil, &ResponseDecodeError{
			Name:      "remove list member",
			Err:       err,
			RateLimit: rl,
		}
	}
	respBody.RateLimit = rl
	return respBody, nil
}

// ListUserMembers returns a list of users who are member of the list
func (c *Client) ListUserMembers(ctx context.Context, listID string, opts ListUserMembersOpts) (*ListUserMembersResponse, error) {
	switch {
	case len(listID) == 0:
		return nil, fmt.Errorf("list user members: an id is required: %w", ErrParameter)
	case opts.MaxResults == 0:
	case opts.MaxResults > listUserMemberMaxResults:
		return nil, fmt.Errorf("list user members: max results [%d] is greater than max [%d]: %w", opts.MaxResults, listUserMemberMaxResults, ErrParameter)
	default:
	}

	ep := listMemberEndpoint.urlID(c.Host, listID)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, ep, nil)
	if err != nil {
		return nil, fmt.Errorf("list user members request: %w", err)
	}
	req.Header.Add("Accept", "application/json")
	c.Authorizer.Add(req)
	opts.addQuery(req)

	resp, err := c.Client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("list user members response: %w", err)
	}
	defer resp.Body.Close()

	decoder := json.NewDecoder(resp.Body)

	rl := rateFromHeader(resp.Header)

	if resp.StatusCode != http.StatusOK {
		e := &ErrorResponse{}
		if err := decoder.Decode(e); err != nil {
			return nil, &HTTPError{
				Status:     resp.Status,
				StatusCode: resp.StatusCode,
				URL:        resp.Request.URL.String(),
				RateLimit:  rl,
			}
		}
		e.StatusCode = resp.StatusCode
		e.RateLimit = rl
		return nil, e
	}

	respBody := struct {
		*UserRaw
		Meta *ListUserMembersMeta `json:"meta"`
	}{}

	if err := decoder.Decode(&respBody); err != nil {
		return nil, &ResponseDecodeError{
			Name:      "list user members",
			Err:       err,
			RateLimit: rl,
		}
	}

	return &ListUserMembersResponse{
		Raw:       respBody.UserRaw,
		Meta:      respBody.Meta,
		RateLimit: rl,
	}, nil
}

// UserListMemberships returns all list a user is a member of
func (c *Client) UserListMemberships(ctx context.Context, userID string, opts UserListMembershipsOpts) (*UserListMembershipsResponse, error) {
	switch {
	case len(userID) == 0:
		return nil, fmt.Errorf("user list membership: an id is required: %w", ErrParameter)
	case opts.MaxResults == 0:
	case opts.MaxResults > userListMembershipMaxResults:
		return nil, fmt.Errorf("user list membership: max results [%d] is greater than max [%d]: %w", opts.MaxResults, userListMembershipMaxResults, ErrParameter)
	default:
	}

	ep := userListMemberEndpoint.urlID(c.Host, userID)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, ep, nil)
	if err != nil {
		return nil, fmt.Errorf("user list membership request: %w", err)
	}
	req.Header.Add("Accept", "application/json")
	c.Authorizer.Add(req)
	opts.addQuery(req)

	resp, err := c.Client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("user list membership response: %w", err)
	}
	defer resp.Body.Close()

	decoder := json.NewDecoder(resp.Body)

	rl := rateFromHeader(resp.Header)

	if resp.StatusCode != http.StatusOK {
		e := &ErrorResponse{}
		if err := decoder.Decode(e); err != nil {
			return nil, &HTTPError{
				Status:     resp.Status,
				StatusCode: resp.StatusCode,
				URL:        resp.Request.URL.String(),
				RateLimit:  rl,
			}
		}
		e.StatusCode = resp.StatusCode
		e.RateLimit = rl
		return nil, e
	}

	respBody := struct {
		*UserListMembershipsRaw
		Meta *UserListMembershipsMeta `json:"meta"`
	}{}

	if err := decoder.Decode(&respBody); err != nil {
		return nil, &ResponseDecodeError{
			Name:      "user list memberships",
			Err:       err,
			RateLimit: rl,
		}
	}

	return &UserListMembershipsResponse{
		Raw:       respBody.UserListMembershipsRaw,
		Meta:      respBody.Meta,
		RateLimit: rl,
	}, nil
}

// UserPinList enables the user to pin a list
func (c *Client) UserPinList(ctx context.Context, userID, listID string) (*UserPinListResponse, error) {
	switch {
	case len(listID) == 0:
		return nil, fmt.Errorf("user pin list: a list id is required: %w", ErrParameter)
	case len(userID) == 0:
		return nil, fmt.Errorf("user pin list: an user id is required: %w", ErrParameter)
	default:
	}

	reqBody := struct {
		ListID string `json:"list_id"`
	}{
		ListID: listID,
	}

	enc, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("user pin list: unable to encode json request %w", err)
	}

	ep := userPinnedListEndpoint.urlID(c.Host, userID)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, ep, bytes.NewReader(enc))
	if err != nil {
		return nil, fmt.Errorf("user pin list request: %w", err)
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")
	c.Authorizer.Add(req)

	resp, err := c.Client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("user pin list response: %w", err)
	}
	defer resp.Body.Close()

	decoder := json.NewDecoder(resp.Body)

	rl := rateFromHeader(resp.Header)

	if resp.StatusCode != http.StatusOK {
		e := &ErrorResponse{}
		if err := decoder.Decode(e); err != nil {
			return nil, &HTTPError{
				Status:     resp.Status,
				StatusCode: resp.StatusCode,
				URL:        resp.Request.URL.String(),
				RateLimit:  rl,
			}
		}
		e.StatusCode = resp.StatusCode
		e.RateLimit = rl
		return nil, e
	}

	respBody := &UserPinListResponse{}

	if err := decoder.Decode(respBody); err != nil {
		return nil, &ResponseDecodeError{
			Name:      "user pin list",
			Err:       err,
			RateLimit: rl,
		}
	}
	respBody.RateLimit = rl
	return respBody, nil
}

// UserUnpinList enables a user to unpin a list
func (c *Client) UserUnpinList(ctx context.Context, userID, listID string) (*UserUnpinListResponse, error) {
	switch {
	case len(listID) == 0:
		return nil, fmt.Errorf("user unpin list: a list id is required: %w", ErrParameter)
	case len(userID) == 0:
		return nil, fmt.Errorf("user unpin list: an user id is required: %w", ErrParameter)
	default:
	}

	ep := userPinnedListEndpoint.urlID(c.Host, userID) + "/" + listID
	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, ep, nil)
	if err != nil {
		return nil, fmt.Errorf("user unpin list request: %w", err)
	}
	req.Header.Add("Accept", "application/json")
	c.Authorizer.Add(req)

	resp, err := c.Client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("user unpin list response: %w", err)
	}
	defer resp.Body.Close()

	decoder := json.NewDecoder(resp.Body)

	rl := rateFromHeader(resp.Header)

	if resp.StatusCode != http.StatusOK {
		e := &ErrorResponse{}
		if err := decoder.Decode(e); err != nil {
			return nil, &HTTPError{
				Status:     resp.Status,
				StatusCode: resp.StatusCode,
				URL:        resp.Request.URL.String(),
				RateLimit:  rl,
			}
		}
		e.StatusCode = resp.StatusCode
		e.RateLimit = rl
		return nil, e
	}

	respBody := &UserUnpinListResponse{}

	if err := decoder.Decode(respBody); err != nil {
		return nil, &ResponseDecodeError{
			Name:      "user unpin list",
			Err:       err,
			RateLimit: rl,
		}
	}
	respBody.RateLimit = rl
	return respBody, nil
}

// UserPinnedLists returns the lists pinned by a user
func (c *Client) UserPinnedLists(ctx context.Context, userID string, opts UserPinnedListsOpts) (*UserPinnedListsResponse, error) {
	switch {
	case len(userID) == 0:
		return nil, fmt.Errorf("user pinned list: an id is required: %w", ErrParameter)
	default:
	}

	ep := userPinnedListEndpoint.urlID(c.Host, userID)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, ep, nil)
	if err != nil {
		return nil, fmt.Errorf("user pinned list request: %w", err)
	}
	req.Header.Add("Accept", "application/json")
	c.Authorizer.Add(req)
	opts.addQuery(req)

	resp, err := c.Client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("user pinned list response: %w", err)
	}
	defer resp.Body.Close()

	decoder := json.NewDecoder(resp.Body)

	rl := rateFromHeader(resp.Header)

	if resp.StatusCode != http.StatusOK {
		e := &ErrorResponse{}
		if err := decoder.Decode(e); err != nil {
			return nil, &HTTPError{
				Status:     resp.Status,
				StatusCode: resp.StatusCode,
				URL:        resp.Request.URL.String(),
				RateLimit:  rl,
			}
		}
		e.StatusCode = resp.StatusCode
		e.RateLimit = rl
		return nil, e
	}

	respBody := struct {
		*UserPinnedListsRaw
		Meta *UserPinnedListsMeta `json:"meta"`
	}{}

	if err := decoder.Decode(&respBody); err != nil {
		return nil, &ResponseDecodeError{
			Name:      "user pinned list",
			Err:       err,
			RateLimit: rl,
		}
	}

	return &UserPinnedListsResponse{
		Raw:       respBody.UserPinnedListsRaw,
		Meta:      respBody.Meta,
		RateLimit: rl,
	}, nil
}

// UserFollowList enables an user to follow a list
func (c *Client) UserFollowList(ctx context.Context, userID, listID string) (*UserFollowListResponse, error) {
	switch {
	case len(listID) == 0:
		return nil, fmt.Errorf("user follow list: a list id is required: %w", ErrParameter)
	case len(userID) == 0:
		return nil, fmt.Errorf("user follow list: an user id is required: %w", ErrParameter)
	default:
	}

	reqBody := struct {
		ListID string `json:"list_id"`
	}{
		ListID: listID,
	}

	enc, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("user follow list: unable to encode json request %w", err)
	}

	ep := userFollowedListEndpoint.urlID(c.Host, userID)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, ep, bytes.NewReader(enc))
	if err != nil {
		return nil, fmt.Errorf("user follow list request: %w", err)
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")
	c.Authorizer.Add(req)

	resp, err := c.Client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("user follow list response: %w", err)
	}
	defer resp.Body.Close()

	decoder := json.NewDecoder(resp.Body)

	rl := rateFromHeader(resp.Header)

	if resp.StatusCode != http.StatusOK {
		e := &ErrorResponse{}
		if err := decoder.Decode(e); err != nil {
			return nil, &HTTPError{
				Status:     resp.Status,
				StatusCode: resp.StatusCode,
				URL:        resp.Request.URL.String(),
				RateLimit:  rl,
			}
		}
		e.StatusCode = resp.StatusCode
		e.RateLimit = rl
		return nil, e
	}

	respBody := &UserFollowListResponse{}

	if err := decoder.Decode(respBody); err != nil {
		return nil, &ResponseDecodeError{
			Name:      "user follow list",
			Err:       err,
			RateLimit: rl,
		}
	}
	respBody.RateLimit = rl
	return respBody, nil
}

// UserUnfollowList enables an user to unfollow a list
func (c *Client) UserUnfollowList(ctx context.Context, userID, listID string) (*UserUnfollowListResponse, error) {
	switch {
	case len(listID) == 0:
		return nil, fmt.Errorf("user unfollow list: a list id is required: %w", ErrParameter)
	case len(userID) == 0:
		return nil, fmt.Errorf("user unfollow list: an user id is required: %w", ErrParameter)
	default:
	}

	ep := userFollowedListEndpoint.urlID(c.Host, userID) + "/" + listID
	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, ep, nil)
	if err != nil {
		return nil, fmt.Errorf("user unfollow list request: %w", err)
	}
	req.Header.Add("Accept", "application/json")
	c.Authorizer.Add(req)

	resp, err := c.Client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("user unfollow list response: %w", err)
	}
	defer resp.Body.Close()

	decoder := json.NewDecoder(resp.Body)

	rl := rateFromHeader(resp.Header)

	if resp.StatusCode != http.StatusOK {
		e := &ErrorResponse{}
		if err := decoder.Decode(e); err != nil {
			return nil, &HTTPError{
				Status:     resp.Status,
				StatusCode: resp.StatusCode,
				URL:        resp.Request.URL.String(),
				RateLimit:  rl,
			}
		}
		e.StatusCode = resp.StatusCode
		e.RateLimit = rl
		return nil, e
	}

	respBody := &UserUnfollowListResponse{}

	if err := decoder.Decode(respBody); err != nil {
		return nil, &ResponseDecodeError{
			Name:      "user unfollow list",
			Err:       err,
			RateLimit: rl,
		}
	}
	respBody.RateLimit = rl
	return respBody, nil
}

// UserFollowedLists returns all list an user follows
func (c *Client) UserFollowedLists(ctx context.Context, userID string, opts UserFollowedListsOpts) (*UserFollowedListsResponse, error) {
	switch {
	case len(userID) == 0:
		return nil, fmt.Errorf("user followed list: an id is required: %w", ErrParameter)
	case opts.MaxResults == 0:
	case opts.MaxResults > userListFollowedMaxResults:
		return nil, fmt.Errorf("user followed list: max results [%d] is greater than max [%d]: %w", opts.MaxResults, userListFollowedMaxResults, ErrParameter)
	default:
	}

	ep := userFollowedListEndpoint.urlID(c.Host, userID)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, ep, nil)
	if err != nil {
		return nil, fmt.Errorf("user followed list request: %w", err)
	}
	req.Header.Add("Accept", "application/json")
	c.Authorizer.Add(req)
	opts.addQuery(req)

	resp, err := c.Client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("user followed list response: %w", err)
	}
	defer resp.Body.Close()

	decoder := json.NewDecoder(resp.Body)

	rl := rateFromHeader(resp.Header)

	if resp.StatusCode != http.StatusOK {
		e := &ErrorResponse{}
		if err := decoder.Decode(e); err != nil {
			return nil, &HTTPError{
				Status:     resp.Status,
				StatusCode: resp.StatusCode,
				URL:        resp.Request.URL.String(),
				RateLimit:  rl,
			}
		}
		e.StatusCode = resp.StatusCode
		e.RateLimit = rl
		return nil, e
	}

	respBody := struct {
		*UserFollowedListsRaw
		Meta *UserFollowedListsMeta `json:"meta"`
	}{}

	if err := decoder.Decode(&respBody); err != nil {
		return nil, &ResponseDecodeError{
			Name:      "user followed list",
			Err:       err,
			RateLimit: rl,
		}
	}

	return &UserFollowedListsResponse{
		Raw:       respBody.UserFollowedListsRaw,
		Meta:      respBody.Meta,
		RateLimit: rl,
	}, nil
}

// ListUserFollowers returns a list of users who are followers of a list
func (c *Client) ListUserFollowers(ctx context.Context, listID string, opts ListUserFollowersOpts) (*ListUserFollowersResponse, error) {
	switch {
	case len(listID) == 0:
		return nil, fmt.Errorf("list user followers: an id is required: %w", ErrParameter)
	case opts.MaxResults == 0:
	case opts.MaxResults > listUserFollowersMaxResults:
		return nil, fmt.Errorf("list user followers: max results [%d] is greater than max [%d]: %w", opts.MaxResults, listUserFollowersMaxResults, ErrParameter)
	default:
	}

	ep := listUserFollowersEndpoint.urlID(c.Host, listID)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, ep, nil)
	if err != nil {
		return nil, fmt.Errorf("list user followers request: %w", err)
	}
	req.Header.Add("Accept", "application/json")
	c.Authorizer.Add(req)
	opts.addQuery(req)

	resp, err := c.Client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("list user followers response: %w", err)
	}
	defer resp.Body.Close()

	decoder := json.NewDecoder(resp.Body)

	rl := rateFromHeader(resp.Header)

	if resp.StatusCode != http.StatusOK {
		e := &ErrorResponse{}
		if err := decoder.Decode(e); err != nil {
			return nil, &HTTPError{
				Status:     resp.Status,
				StatusCode: resp.StatusCode,
				URL:        resp.Request.URL.String(),
				RateLimit:  rl,
			}
		}
		e.StatusCode = resp.StatusCode
		e.RateLimit = rl
		return nil, e
	}

	respBody := struct {
		*UserRaw
		Meta *ListUserFollowersMeta `json:"meta"`
	}{}

	if err := decoder.Decode(&respBody); err != nil {
		return nil, &ResponseDecodeError{
			Name:      "list user followers",
			Err:       err,
			RateLimit: rl,
		}
	}

	return &ListUserFollowersResponse{
		Raw:       respBody.UserRaw,
		Meta:      respBody.Meta,
		RateLimit: rl,
	}, nil
}

// SpacesLookup returns information about a space from the requested ids
func (c *Client) SpacesLookup(ctx context.Context, ids []string, opts SpacesLookupOpts) (*SpacesLookupResponse, error) {
	ep := spaceLookupEndpoint.url(c.Host)
	switch {
	case len(ids) == 0:
		return nil, fmt.Errorf("space lookup: an id is required: %w", ErrParameter)
	case len(ids) > spaceMaxIDs:
		return nil, fmt.Errorf("space lookup: ids %d is greater than max %d: %w", len(ids), spaceMaxIDs, ErrParameter)
	case len(ids) == 1:
		ep += fmt.Sprintf("/%s", ids[0])
	default:
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, ep, nil)
	if err != nil {
		return nil, fmt.Errorf("space lookup request: %w", err)
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
		return nil, fmt.Errorf("space lookup response: %w", err)
	}
	defer resp.Body.Close()

	decoder := json.NewDecoder(resp.Body)

	rl := rateFromHeader(resp.Header)

	if resp.StatusCode != http.StatusOK {
		e := &ErrorResponse{}
		if err := decoder.Decode(e); err != nil {
			return nil, &HTTPError{
				Status:     resp.Status,
				StatusCode: resp.StatusCode,
				URL:        resp.Request.URL.String(),
				RateLimit:  rl,
			}
		}
		e.StatusCode = resp.StatusCode
		e.RateLimit = rl
		return nil, e
	}

	raw := &SpacesRaw{}
	switch {
	case len(ids) == 1:
		single := &spaceRaw{}
		if err := decoder.Decode(single); err != nil {
			return nil, &ResponseDecodeError{
				Name:      "space lookup",
				Err:       err,
				RateLimit: rl,
			}
		}
		raw.Spaces = make([]*SpaceObj, 1)
		raw.Spaces[0] = single.Space
		raw.Includes = single.Includes
		raw.Errors = single.Errors
	default:
		if err := decoder.Decode(raw); err != nil {
			return nil, &ResponseDecodeError{
				Name:      "space lookup ",
				Err:       err,
				RateLimit: rl,
			}
		}
	}
	return &SpacesLookupResponse{
		Raw:       raw,
		RateLimit: rl,
	}, nil
}

// SpacesByCreatorLookup returns live or scheduled spaces created by a specific user ids
func (c *Client) SpacesByCreatorLookup(ctx context.Context, userIDs []string, opts SpacesByCreatorLookupOpts) (*SpacesByCreatorLookupResponse, error) {
	switch {
	case len(userIDs) == 0:
		return nil, fmt.Errorf("space by creator lookup: an id is required: %w", ErrParameter)
	case len(userIDs) > spaceByCreatorMaxIDs:
		return nil, fmt.Errorf("space by creator lookup: ids %d is greater than max %d: %w", len(userIDs), spaceByCreatorMaxIDs, ErrParameter)
	default:
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, spaceByCreatorLookupEndpoint.url(c.Host), nil)
	if err != nil {
		return nil, fmt.Errorf("space by creator lookup request: %w", err)
	}
	req.Header.Add("Accept", "application/json")
	c.Authorizer.Add(req)
	opts.addQuery(req)
	q := req.URL.Query()
	q.Add("user_ids", strings.Join(userIDs, ","))
	req.URL.RawQuery = q.Encode()

	resp, err := c.Client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("space by creator lookup response: %w", err)
	}
	defer resp.Body.Close()

	decoder := json.NewDecoder(resp.Body)

	rl := rateFromHeader(resp.Header)

	if resp.StatusCode != http.StatusOK {
		e := &ErrorResponse{}
		if err := decoder.Decode(e); err != nil {
			return nil, &HTTPError{
				Status:     resp.Status,
				StatusCode: resp.StatusCode,
				URL:        resp.Request.URL.String(),
				RateLimit:  rl,
			}
		}
		e.StatusCode = resp.StatusCode
		e.RateLimit = rl
		return nil, e
	}

	raw := struct {
		*SpacesRaw
		Meta *SpacesByCreatorMeta `json:"meta"`
	}{}

	if err := decoder.Decode(&raw); err != nil {
		return nil, &ResponseDecodeError{
			Name:      "space by creator lookup ",
			Err:       err,
			RateLimit: rl,
		}
	}
	return &SpacesByCreatorLookupResponse{
		Raw:       raw.SpacesRaw,
		Meta:      raw.Meta,
		RateLimit: rl,
	}, nil
}

// SpaceBuyersLookup returns a list of users who purchased a ticket to the requested space
func (c *Client) SpaceBuyersLookup(ctx context.Context, spaceID string, opts SpaceBuyersLookupOpts) (*SpaceBuyersLookupResponse, error) {
	switch {
	case len(spaceID) == 0:
		return nil, fmt.Errorf("space buyers lookup: an id is required: %w", ErrParameter)
	default:
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, spaceBuyersLookupEndpoint.urlID(c.Host, spaceID), nil)
	if err != nil {
		return nil, fmt.Errorf("space buyers lookup request: %w", err)
	}
	req.Header.Add("Accept", "application/json")
	c.Authorizer.Add(req)
	opts.addQuery(req)

	resp, err := c.Client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("space buyers lookup response: %w", err)
	}
	defer resp.Body.Close()

	decoder := json.NewDecoder(resp.Body)

	rl := rateFromHeader(resp.Header)

	if resp.StatusCode != http.StatusOK {
		e := &ErrorResponse{}
		if err := decoder.Decode(e); err != nil {
			return nil, &HTTPError{
				Status:     resp.Status,
				StatusCode: resp.StatusCode,
				URL:        resp.Request.URL.String(),
				RateLimit:  rl,
			}
		}
		e.StatusCode = resp.StatusCode
		e.RateLimit = rl
		return nil, e
	}

	raw := &UserRaw{}

	if err := decoder.Decode(raw); err != nil {
		return nil, &ResponseDecodeError{
			Name:      "space buyers lookup ",
			Err:       err,
			RateLimit: rl,
		}
	}
	return &SpaceBuyersLookupResponse{
		Raw:       raw,
		RateLimit: rl,
	}, nil
}

// SpaceTweetsLookup returns tweets shared in the request space
func (c *Client) SpaceTweetsLookup(ctx context.Context, spaceID string, opts SpaceTweetsLookupOpts) (*SpaceTweetsLookupResponse, error) {
	switch {
	case len(spaceID) == 0:
		return nil, fmt.Errorf("space tweets lookup: an id is required: %w", ErrParameter)
	default:
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, spaceTweetsLookupEndpoint.urlID(c.Host, spaceID), nil)
	if err != nil {
		return nil, fmt.Errorf("space tweets lookup request: %w", err)
	}
	req.Header.Add("Accept", "application/json")
	c.Authorizer.Add(req)
	opts.addQuery(req)

	resp, err := c.Client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("space tweets lookup response: %w", err)
	}
	defer resp.Body.Close()

	decoder := json.NewDecoder(resp.Body)

	rl := rateFromHeader(resp.Header)

	if resp.StatusCode != http.StatusOK {
		e := &ErrorResponse{}
		if err := decoder.Decode(e); err != nil {
			return nil, &HTTPError{
				Status:     resp.Status,
				StatusCode: resp.StatusCode,
				URL:        resp.Request.URL.String(),
				RateLimit:  rl,
			}
		}
		e.StatusCode = resp.StatusCode
		e.RateLimit = rl
		return nil, e
	}

	raw := struct {
		*TweetRaw
		Meta *SpaceTweetsLookupMeta `json:"meta"`
	}{}

	if err := decoder.Decode(&raw); err != nil {
		return nil, &ResponseDecodeError{
			Name:      "space tweets lookup ",
			Err:       err,
			RateLimit: rl,
		}
	}
	return &SpaceTweetsLookupResponse{
		Raw:       raw.TweetRaw,
		Meta:      raw.Meta,
		RateLimit: rl,
	}, nil
}

// SpacesSearch returns live or scheduled spaces matching the specified search terms.
func (c *Client) SpacesSearch(ctx context.Context, query string, opts SpacesSearchOpts) (*SpacesSearchResponse, error) {
	switch {
	case len(query) == 0:
		return nil, fmt.Errorf("space search: a query is required: %w", ErrParameter)
	default:
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, spaceSearchEndpoint.url(c.Host), nil)
	if err != nil {
		return nil, fmt.Errorf("space search request: %w", err)
	}
	req.Header.Add("Accept", "application/json")
	c.Authorizer.Add(req)
	opts.addQuery(req)
	q := req.URL.Query()
	q.Add("query", query)
	req.URL.RawQuery = q.Encode()

	resp, err := c.Client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("space search response: %w", err)
	}
	defer resp.Body.Close()

	decoder := json.NewDecoder(resp.Body)

	rl := rateFromHeader(resp.Header)

	if resp.StatusCode != http.StatusOK {
		e := &ErrorResponse{}
		if err := decoder.Decode(e); err != nil {
			return nil, &HTTPError{
				Status:     resp.Status,
				StatusCode: resp.StatusCode,
				URL:        resp.Request.URL.String(),
				RateLimit:  rl,
			}
		}
		e.StatusCode = resp.StatusCode
		e.RateLimit = rl
		return nil, e
	}

	raw := struct {
		*SpacesRaw
		Meta *SpacesSearchMeta `json:"meta"`
	}{}

	if err := decoder.Decode(&raw); err != nil {
		return nil, &ResponseDecodeError{
			Name:      "space search",
			Err:       err,
			RateLimit: rl,
		}
	}

	return &SpacesSearchResponse{
		Raw:       raw.SpacesRaw,
		Meta:      raw.Meta,
		RateLimit: rl,
	}, nil
}

// CreateComplianceBatchJob creates a new compliance job for tweet or user IDs.
func (c *Client) CreateComplianceBatchJob(ctx context.Context, jobType ComplianceBatchJobType, opts CreateComplianceBatchJobOpts) (*CreateComplianceBatchJobResponse, error) {
	switch {
	case len(jobType) == 0:
		return nil, fmt.Errorf("create compliance batch job: a type is required: %w", ErrParameter)
	default:
	}

	request := struct {
		Type      ComplianceBatchJobType `json:"type"`
		Name      string                 `json:"name,omitempty"`
		Resumable bool                   `json:"resumable,omitempty"`
	}{
		Type:      jobType,
		Name:      opts.Name,
		Resumable: opts.Resumable,
	}

	enc, err := json.Marshal(request)
	if err != nil {
		return nil, fmt.Errorf("create compliance batch job request encode: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, complianceJobsEndpoint.url(c.Host), bytes.NewReader(enc))
	if err != nil {
		return nil, fmt.Errorf("create compliance batch job request: %w", err)
	}
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")

	resp, err := c.Client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("create compliance batch job response: %w", err)
	}
	defer resp.Body.Close()

	decoder := json.NewDecoder(resp.Body)

	rl := rateFromHeader(resp.Header)

	if resp.StatusCode != http.StatusOK {
		e := &ErrorResponse{}
		if err := decoder.Decode(e); err != nil {
			return nil, &HTTPError{
				Status:     resp.Status,
				StatusCode: resp.StatusCode,
				URL:        resp.Request.URL.String(),
				RateLimit:  rl,
			}
		}
		e.StatusCode = resp.StatusCode
		e.RateLimit = rl
		return nil, e
	}

	raw := &ComplianceBatchJobRaw{}

	if err := decoder.Decode(&raw); err != nil {
		return nil, &ResponseDecodeError{
			Name:      "create compliance batch job",
			Err:       err,
			RateLimit: rl,
		}
	}
	raw.Job.client = c.Client

	return &CreateComplianceBatchJobResponse{
		Raw:       raw,
		RateLimit: rl,
	}, nil
}

// ComplianceBatchJob returns a single compliance job
func (c *Client) ComplianceBatchJob(ctx context.Context, id string) (*ComplianceBatchJobResponse, error) {
	switch {
	case len(id) == 0:
		return nil, fmt.Errorf("compliance batch job: a type is required: %w", ErrParameter)
	default:
	}

	ep := complianceJobsEndpoint.url(c.Host) + "/" + id
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, ep, nil)
	if err != nil {
		return nil, fmt.Errorf("compliance batch job request: %w", err)
	}
	req.Header.Add("Accept", "application/json")

	resp, err := c.Client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("compliance batch job response: %w", err)
	}
	defer resp.Body.Close()

	decoder := json.NewDecoder(resp.Body)

	rl := rateFromHeader(resp.Header)

	if resp.StatusCode != http.StatusOK {
		e := &ErrorResponse{}
		if err := decoder.Decode(e); err != nil {
			return nil, &HTTPError{
				Status:     resp.Status,
				StatusCode: resp.StatusCode,
				URL:        resp.Request.URL.String(),
				RateLimit:  rl,
			}
		}
		e.StatusCode = resp.StatusCode
		e.RateLimit = rl
		return nil, e
	}

	raw := &ComplianceBatchJobRaw{}

	if err := decoder.Decode(&raw); err != nil {
		return nil, &ResponseDecodeError{
			Name:      "compliance batch job",
			Err:       err,
			RateLimit: rl,
		}
	}
	raw.Job.client = c.Client

	return &ComplianceBatchJobResponse{
		Raw:       raw,
		RateLimit: rl,
	}, nil
}

// ComplianceBatchJobLookup returns a list of compliance jobs
func (c *Client) ComplianceBatchJobLookup(ctx context.Context, jobType ComplianceBatchJobType, opts ComplianceBatchJobLookupOpts) (*ComplianceBatchJobLookupResponse, error) {
	switch {
	case len(jobType) == 0:
		return nil, fmt.Errorf("compliance batch job lookup: a type is required: %w", ErrParameter)
	default:
	}

	ep := complianceJobsEndpoint.url(c.Host)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, ep, nil)
	if err != nil {
		return nil, fmt.Errorf("compliance batch job lookup request: %w", err)
	}
	req.Header.Add("Accept", "application/json")

	opts.addQuery(req)
	q := req.URL.Query()
	q.Add("type", string(jobType))
	req.URL.RawQuery = q.Encode()

	resp, err := c.Client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("compliance batch job lookup response: %w", err)
	}
	defer resp.Body.Close()

	decoder := json.NewDecoder(resp.Body)

	rl := rateFromHeader(resp.Header)

	if resp.StatusCode != http.StatusOK {
		e := &ErrorResponse{}
		if err := decoder.Decode(e); err != nil {
			return nil, &HTTPError{
				Status:     resp.Status,
				StatusCode: resp.StatusCode,
				URL:        resp.Request.URL.String(),
				RateLimit:  rl,
			}
		}
		e.StatusCode = resp.StatusCode
		e.RateLimit = rl
		return nil, e
	}

	raw := &ComplianceBatchJobsRaw{}

	if err := decoder.Decode(&raw); err != nil {
		return nil, &ResponseDecodeError{
			Name:      "compliance batch job lookup",
			Err:       err,
			RateLimit: rl,
		}
	}

	for i := range raw.Jobs {
		raw.Jobs[i].client = c.Client
	}

	return &ComplianceBatchJobLookupResponse{
		Raw:       raw,
		RateLimit: rl,
	}, nil
}

// QuoteTweetsLookup returns quote tweets for a tweet specified by the requested tweet id
func (c *Client) QuoteTweetsLookup(ctx context.Context, tweetID string, opts QuoteTweetsLookupOpts) (*QuoteTweetsLookupResponse, error) {
	switch {
	case len(tweetID) == 0:
		return nil, fmt.Errorf("quote tweets lookup: an id is required: %w", ErrParameter)
	case opts.MaxResults == 0:
	case opts.MaxResults < quoteTweetMinResults:
		return nil, fmt.Errorf("quote tweets lookup: a min results [%d] is required [current: %d]: %w", quoteTweetMinResults, opts.MaxResults, ErrParameter)
	case opts.MaxResults > quoteTweetMaxResults:
		return nil, fmt.Errorf("quote tweets lookup: a max results [%d] is required [current: %d]: %w", quoteTweetMaxResults, opts.MaxResults, ErrParameter)
	default:
	}

	ep := quoteTweetLookupEndpoint.urlID(c.Host, tweetID)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, ep, nil)
	if err != nil {
		return nil, fmt.Errorf("quote tweets lookup request: %w", err)
	}
	req.Header.Add("Accept", "application/json")
	c.Authorizer.Add(req)
	opts.addQuery(req)

	resp, err := c.Client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("quote tweets lookup response: %w", err)
	}
	defer resp.Body.Close()

	decoder := json.NewDecoder(resp.Body)

	rl := rateFromHeader(resp.Header)

	if resp.StatusCode != http.StatusOK {
		e := &ErrorResponse{}
		if err := decoder.Decode(e); err != nil {
			return nil, &HTTPError{
				Status:     resp.Status,
				StatusCode: resp.StatusCode,
				URL:        resp.Request.URL.String(),
				RateLimit:  rl,
			}
		}
		e.StatusCode = resp.StatusCode
		e.RateLimit = rl
		return nil, e
	}

	respBody := struct {
		*TweetRaw
		Meta *QuoteTweetsLookupMeta `json:"meta"`
	}{}

	if err := decoder.Decode(&respBody); err != nil {
		return nil, &ResponseDecodeError{
			Name:      "quote tweets lookup",
			Err:       err,
			RateLimit: rl,
		}
	}

	return &QuoteTweetsLookupResponse{
		Raw:       respBody.TweetRaw,
		Meta:      respBody.Meta,
		RateLimit: rl,
	}, nil
}

// TweetBookmarksLookup allows you to get an authenticated user's 800 most recent bookmarked Tweets
func (c *Client) TweetBookmarksLookup(ctx context.Context, userID string, opts TweetBookmarksLookupOpts) (*TweetBookmarksLookupResponse, error) {
	switch {
	case len(userID) == 0:
		return nil, fmt.Errorf("tweet bookmarks lookup: an id is required: %w", ErrParameter)
	case opts.MaxResults == 0:
	case opts.MaxResults > tweetBookmarksMaxResults:
		return nil, fmt.Errorf("tweet bookmarks lookup: a max results [%d] is required [current: %d]: %w", tweetBookmarksMaxResults, opts.MaxResults, ErrParameter)
	default:
	}

	ep := tweetBookmarksEndpoint.urlID(c.Host, userID)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, ep, nil)
	if err != nil {
		return nil, fmt.Errorf("tweet bookmarks lookup request: %w", err)
	}
	req.Header.Add("Accept", "application/json")
	c.Authorizer.Add(req)
	opts.addQuery(req)

	resp, err := c.Client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("tweet bookmarks lookup response: %w", err)
	}
	defer resp.Body.Close()

	decoder := json.NewDecoder(resp.Body)

	rl := rateFromHeader(resp.Header)

	if resp.StatusCode != http.StatusOK {
		e := &ErrorResponse{}
		if err := decoder.Decode(e); err != nil {
			return nil, &HTTPError{
				Status:     resp.Status,
				StatusCode: resp.StatusCode,
				URL:        resp.Request.URL.String(),
				RateLimit:  rl,
			}
		}
		e.StatusCode = resp.StatusCode
		e.RateLimit = rl
		return nil, e
	}

	respBody := struct {
		*TweetRaw
		Meta *TweetBookmarksLookupMeta `json:"meta"`
	}{}

	if err := decoder.Decode(&respBody); err != nil {
		return nil, &ResponseDecodeError{
			Name:      "tweet bookmarks lookup",
			Err:       err,
			RateLimit: rl,
		}
	}

	return &TweetBookmarksLookupResponse{
		Raw:       respBody.TweetRaw,
		Meta:      respBody.Meta,
		RateLimit: rl,
	}, nil
}

// AddTweetBookmark causes the user ID identified in the path parameter to Bookmark the target Tweet provided in the request body
func (c *Client) AddTweetBookmark(ctx context.Context, userID, tweetID string) (*AddTweetBookmarkResponse, error) {
	switch {
	case len(userID) == 0:
		return nil, fmt.Errorf("tweet bookmarks add: an user id is required: %w", ErrParameter)
	case len(tweetID) == 0:
		return nil, fmt.Errorf("tweet bookmarks add: a tweet id is required: %w", ErrParameter)
	default:
	}

	rb := struct {
		TweetID string `json:"tweet_id"`
	}{
		TweetID: tweetID,
	}
	enc, err := json.Marshal(rb)
	if err != nil {
		return nil, fmt.Errorf("tweet bookmarks add body encoding: %w", err)
	}

	ep := tweetBookmarksEndpoint.urlID(c.Host, userID)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, ep, bytes.NewReader(enc))
	if err != nil {
		return nil, fmt.Errorf("tweet bookmarks add request: %w", err)
	}
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")
	c.Authorizer.Add(req)

	resp, err := c.Client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("tweet bookmarks add response: %w", err)
	}
	defer resp.Body.Close()

	decoder := json.NewDecoder(resp.Body)

	rl := rateFromHeader(resp.Header)

	if resp.StatusCode != http.StatusOK {
		e := &ErrorResponse{}
		if err := decoder.Decode(e); err != nil {
			return nil, &HTTPError{
				Status:     resp.Status,
				StatusCode: resp.StatusCode,
				URL:        resp.Request.URL.String(),
				RateLimit:  rl,
			}
		}
		e.StatusCode = resp.StatusCode
		e.RateLimit = rl
		return nil, e
	}

	respBody := &AddTweetBookmarkResponse{}

	if err := decoder.Decode(respBody); err != nil {
		return nil, &ResponseDecodeError{
			Name:      "tweet bookmarks add",
			Err:       err,
			RateLimit: rl,
		}
	}

	respBody.RateLimit = rl

	return respBody, nil
}

// RemoveTweetBookmark allows a user or authenticated user ID to remove a Bookmark of a Tweet
func (c *Client) RemoveTweetBookmark(ctx context.Context, userID, tweetID string) (*RemoveTweetBookmarkResponse, error) {
	switch {
	case len(userID) == 0:
		return nil, fmt.Errorf("tweet bookmarks remove: an user id is required: %w", ErrParameter)
	case len(tweetID) == 0:
		return nil, fmt.Errorf("tweet bookmarks remove: a tweet id is required: %w", ErrParameter)
	default:
	}

	ep := tweetBookmarksEndpoint.urlID(c.Host, userID) + "/" + tweetID

	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, ep, nil)
	if err != nil {
		return nil, fmt.Errorf("tweet bookmarks remove request: %w", err)
	}
	req.Header.Add("Accept", "application/json")
	c.Authorizer.Add(req)

	resp, err := c.Client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("tweet bookmarks remove response: %w", err)
	}
	defer resp.Body.Close()

	decoder := json.NewDecoder(resp.Body)

	rl := rateFromHeader(resp.Header)

	if resp.StatusCode != http.StatusOK {
		e := &ErrorResponse{}
		if err := decoder.Decode(e); err != nil {
			return nil, &HTTPError{
				Status:     resp.Status,
				StatusCode: resp.StatusCode,
				URL:        resp.Request.URL.String(),
				RateLimit:  rl,
			}
		}
		e.StatusCode = resp.StatusCode
		e.RateLimit = rl
		return nil, e
	}

	respBody := &RemoveTweetBookmarkResponse{}

	if err := decoder.Decode(respBody); err != nil {
		return nil, &ResponseDecodeError{
			Name:      "tweet bookmarks remove",
			Err:       err,
			RateLimit: rl,
		}
	}

	respBody.RateLimit = rl

	return respBody, nil
}
