package twitter

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

const (
	tweetMaxIDs                  = 100
	userMaxIDs                   = 100
	userMaxNames                 = 100
	tweetRecentSearchQueryLength = 512
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

	respBytes, err := ioutil.ReadAll(resp.Body)
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

	respBytes, err := ioutil.ReadAll(resp.Body)
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

	respBytes, err := ioutil.ReadAll(resp.Body)
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

	respBytes, err := ioutil.ReadAll(resp.Body)
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

	respBytes, err := ioutil.ReadAll(resp.Body)
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

	respBytes, err := ioutil.ReadAll(resp.Body)
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
