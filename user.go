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
	userLookupEndpoint           = "2/users"
	userNameLookupEndpoint       = "2/users/by/username"
	userNamesLookupEndpoint      = "2/users/by"
	userFollowingLookupEndpoint  = "2/users/{id}/following"
	userFollowersLookupEndpoint  = "2/users/{id}/followers"
	userTimelineTweetsEndpoint   = "2/users/{id}/tweets"
	userTimelineMentionsEndpoint = "2/users/{id}/mentions"
	userID                       = "{id}"
	userMaxIDs                   = 100
	userMaxNames                 = 100
)

// UserLookups is a map of user lookups
type UserLookups map[string]UserLookup

func (t UserLookups) lookup(decoder *json.Decoder) error {
	type include struct {
		Tweet []*TweetObj `json:"tweets"`
	}
	type body struct {
		Data    UserObj `json:"data"`
		Include include `json:"includes"`
	}
	b := &body{}
	if err := decoder.Decode(b); err != nil {
		return fmt.Errorf("tweet lookup decode error %w", err)
	}

	ul := UserLookup{
		User: b.Data,
	}
	if len(b.Include.Tweet) > 0 {
		ul.Tweet = b.Include.Tweet[0]
	}
	t[b.Data.ID] = ul

	return nil
}

func (t UserLookups) lookups(decoder *json.Decoder) error {
	type include struct {
		Tweet []*TweetObj `json:"tweets"`
	}
	type body struct {
		Data    []UserObj `json:"data"`
		Include include   `json:"includes"`
	}
	b := &body{}
	if err := decoder.Decode(b); err != nil {
		return fmt.Errorf("tweet lookup decode error %w", err)
	}

	pinnedTweets := map[string]*TweetObj{}
	for _, tweet := range b.Include.Tweet {
		pinnedTweets[tweet.ID] = tweet
	}

	for _, user := range b.Data {
		ul := UserLookup{
			User: user,
		}
		if tweet, has := pinnedTweets[user.PinnedTweetID]; has {
			ul.Tweet = tweet
		}
		t[user.ID] = ul
	}
	return nil
}

// UserLookup is a complete user objects
type UserLookup struct {
	User  UserObj
	Tweet *TweetObj
}

// UserFollowLookup contains all of the user following information
type UserFollowLookup struct {
	Lookups UserLookups
	Meta    *UserFollowMeta
	Errors  []ErrorObj
}

// UserFollowMeta the meta that is returned for the following APIs
type UserFollowMeta struct {
	ResultCount   int    `json:"result_count"`
	PreviousToken string `json:"previous_token"`
	NextToken     string `json:"next_token"`
}

// UserTimeline is the response to the user tweet timeline API
type UserTimeline struct {
	Tweets   []TweetObj            `json:"data"`
	Includes *UserTimelineIncludes `json:"includes"`
	Errors   []ErrorObj            `json:"errors"`
	Meta     UserTimelineMeta      `json:"meta"`
}

// UserTimelineIncludes will contain the optional response objects
type UserTimelineIncludes struct {
	Medias []MediaObj `json:"media"`
	Users  []UserObj  `json:"users"`
	Tweets []TweetObj `json:"tweets"`
	Places []PlaceObj `json:"places"`
	Polls  string     `json:"polls"`
}

// UserTimelineMeta is the meta data of the response
type UserTimelineMeta struct {
	OldestID      string `json:"oldest_id"`
	NewestID      string `json:"newest_id"`
	ResultCount   int    `json:"result_count"`
	NextToken     string `json:"next_token"`
	PreviousToken string `json:"previous_token"`
}

// User represents the User v2 APIs
type User struct {
	Authorizer Authorizer
	Client     *http.Client
	Host       string
}

// Lookup can be used to look up a user by their id
func (u *User) Lookup(ctx context.Context, ids []string, fieldOpts UserFieldOptions) (UserLookups, error) {
	ep := userLookupEndpoint
	switch {
	case len(ids) == 0:
		return nil, fmt.Errorf("user lookup an id is required")
	case len(ids) > userMaxIDs:
		return nil, fmt.Errorf("user lookup: ids %d is greater than max %d", len(ids), userMaxIDs)
	case len(ids) == 1:
		ep += fmt.Sprintf("/%s", ids[0])
	default:
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, fmt.Sprintf("%s/%s", u.Host, ep), nil)
	if err != nil {
		return nil, fmt.Errorf("user lookup request: %w", err)
	}
	req.Header.Add("Accept", "application/json")
	u.Authorizer.Add(req)
	fieldOpts.addQuery(req)
	if len(ids) > 1 {
		q := req.URL.Query()
		q.Add("ids", strings.Join(ids, ","))
		req.URL.RawQuery = q.Encode()
	}

	resp, err := u.Client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("user lookup response: %w", err)
	}
	defer resp.Body.Close()

	decoder := json.NewDecoder(resp.Body)

	if resp.StatusCode != http.StatusOK {
		e := &TweetErrorResponse{}
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

	ul := UserLookups{}
	if len(ids) == 1 {
		if err := ul.lookup(decoder); err != nil {
			return nil, err
		}
		return ul, nil
	}

	if err := ul.lookups(decoder); err != nil {
		return nil, err
	}
	return ul, nil
}

// LookupUsername will retuen the user information from its user names
func (u *User) LookupUsername(ctx context.Context, usernames []string, fieldOpts UserFieldOptions) (UserLookups, error) {
	ep := userNamesLookupEndpoint
	switch {
	case len(usernames) == 0:
		return nil, fmt.Errorf("user lookup name is required")
	case len(usernames) > userMaxNames:
		return nil, fmt.Errorf("user lookup: names %d is greater than max %d", len(usernames), userMaxNames)
	case len(usernames) == 1:
		ep = fmt.Sprintf("%s/%s", userNameLookupEndpoint, usernames[0])
	default:
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, fmt.Sprintf("%s/%s", u.Host, ep), nil)
	if err != nil {
		return nil, fmt.Errorf("user lookup request: %w", err)
	}
	req.Header.Add("Accept", "application/json")
	u.Authorizer.Add(req)
	fieldOpts.addQuery(req)
	if len(usernames) > 1 {
		q := req.URL.Query()
		q.Add("usernames", strings.Join(usernames, ","))
		req.URL.RawQuery = q.Encode()
	}

	resp, err := u.Client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("user lookup response: %w", err)
	}
	defer resp.Body.Close()

	decoder := json.NewDecoder(resp.Body)

	if resp.StatusCode != http.StatusOK {
		e := &TweetErrorResponse{}
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

	ul := UserLookups{}
	if len(usernames) == 1 {
		if err := ul.lookup(decoder); err != nil {
			return nil, err
		}
		return ul, nil
	}

	if err := ul.lookups(decoder); err != nil {
		return nil, err
	}
	return ul, nil
}

// LookupFollowing handles the user following callout
func (u *User) LookupFollowing(ctx context.Context, id string, followOpts UserFollowOptions) (*UserFollowLookup, error) {
	switch {
	case len(id) == 0:
		return nil, fmt.Errorf("user id must be present for following lookup")
	case followOpts.MaxResults < 0 || followOpts.MaxResults > 1000:
		return nil, fmt.Errorf("user max results for following lookup must be between 1-1000: %d", followOpts.MaxResults)
	default:
	}

	ep := fmt.Sprintf("%s/%s", u.Host, userFollowingLookupEndpoint)
	ep = strings.Replace(ep, userID, id, -1)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, ep, nil)
	if err != nil {
		return nil, fmt.Errorf("user lookup following request: %w", err)
	}
	req.Header.Add("Accept", "application/json")
	u.Authorizer.Add(req)
	followOpts.addQuery(req)

	resp, err := u.Client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("user lookup response: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("user lookup following reading body: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		e := &TweetErrorResponse{}
		if err := json.Unmarshal(body, e); err != nil {
			return nil, &HTTPError{
				Status:     resp.Status,
				StatusCode: resp.StatusCode,
				URL:        resp.Request.URL.String(),
			}
		}
		e.StatusCode = resp.StatusCode
		return nil, e
	}

	ul := UserLookups{}
	if err := ul.lookups(json.NewDecoder(bytes.NewReader(body))); err != nil {
		return nil, fmt.Errorf("user lookup response lookup decode: %w", err)
	}
	type extra struct {
		Meta   *UserFollowMeta `json:"meta"`
		Errors []ErrorObj      `json:"errors"`
	}
	ufm := &extra{}
	if err := json.Unmarshal(body, ufm); err != nil {
		return nil, fmt.Errorf("user lookup response meta decode: %w", err)
	}
	return &UserFollowLookup{
		Lookups: ul,
		Meta:    ufm.Meta,
		Errors:  ufm.Errors,
	}, nil
}

// LookupFollowers will return a users followers
func (u *User) LookupFollowers(ctx context.Context, id string, followOpts UserFollowOptions) (*UserFollowLookup, error) {
	switch {
	case len(id) == 0:
		return nil, fmt.Errorf("user id must be present for following lookup")
	case followOpts.MaxResults < 0 || followOpts.MaxResults > 1000:
		return nil, fmt.Errorf("user max results for following lookup must be between 1-1000: %d", followOpts.MaxResults)
	default:
	}

	ep := fmt.Sprintf("%s/%s", u.Host, userFollowersLookupEndpoint)
	ep = strings.Replace(ep, userID, id, -1)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, ep, nil)
	if err != nil {
		return nil, fmt.Errorf("user lookup following request: %w", err)
	}
	req.Header.Add("Accept", "application/json")
	u.Authorizer.Add(req)
	followOpts.addQuery(req)

	resp, err := u.Client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("user lookup response: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("user lookup following reading body: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		e := &TweetErrorResponse{}
		if err := json.Unmarshal(body, e); err != nil {
			return nil, &HTTPError{
				Status:     resp.Status,
				StatusCode: resp.StatusCode,
				URL:        resp.Request.URL.String(),
			}
		}
		e.StatusCode = resp.StatusCode
		return nil, e
	}

	ul := UserLookups{}
	if err := ul.lookups(json.NewDecoder(bytes.NewReader(body))); err != nil {
		return nil, fmt.Errorf("user lookup response lookup decode: %w", err)
	}
	type extra struct {
		Meta   *UserFollowMeta `json:"meta"`
		Errors []ErrorObj      `json:"errors"`
	}
	ufm := &extra{}
	if err := json.Unmarshal(body, ufm); err != nil {
		return nil, fmt.Errorf("user lookup response meta decode: %w", err)
	}
	return &UserFollowLookup{
		Lookups: ul,
		Meta:    ufm.Meta,
		Errors:  ufm.Errors,
	}, nil
}

// Tweets is the user timeline tweets
func (u *User) Tweets(ctx context.Context, id string, tweetOpts UserTimelineOpts) (*UserTimeline, error) {
	switch {
	case len(id) == 0:
		return nil, fmt.Errorf("user id must be present for timeline tweets")
	case tweetOpts.MaxResults < 0 || tweetOpts.MaxResults > 100:
		return nil, fmt.Errorf("user max results for timeline tweets must be between 1-1000: %d", tweetOpts.MaxResults)
	default:
	}

	ep := fmt.Sprintf("%s/%s", u.Host, userTimelineTweetsEndpoint)
	ep = strings.Replace(ep, userID, id, -1)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, ep, nil)
	if err != nil {
		return nil, fmt.Errorf("user lookup following request: %w", err)
	}
	req.Header.Add("Accept", "application/json")
	u.Authorizer.Add(req)
	tweetOpts.addQuery(req)

	resp, err := u.Client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("user lookup response: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("user lookup following reading body: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		e := &TweetErrorResponse{}
		if err := json.Unmarshal(body, e); err != nil {
			return nil, &HTTPError{
				Status:     resp.Status,
				StatusCode: resp.StatusCode,
				URL:        resp.Request.URL.String(),
			}
		}
		e.StatusCode = resp.StatusCode
		return nil, e
	}

	result := &UserTimeline{}
	if err := json.Unmarshal(body, result); err != nil {
		return nil, fmt.Errorf("user tweet timeline response decode: %w", err)
	}
	return result, nil
}

// Mentions will return back the user tweets mentions timeline
func (u *User) Mentions(ctx context.Context, id string, tweetOpts UserTimelineOpts) (*UserTimeline, error) {
	switch {
	case len(id) == 0:
		return nil, fmt.Errorf("user id must be present for timeline tweets")
	case tweetOpts.MaxResults < 0 || tweetOpts.MaxResults > 100:
		return nil, fmt.Errorf("user max results for timeline tweets must be between 1-1000: %d", tweetOpts.MaxResults)
	default:
	}

	ep := fmt.Sprintf("%s/%s", u.Host, userTimelineMentionsEndpoint)
	ep = strings.Replace(ep, userID, id, -1)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, ep, nil)
	if err != nil {
		return nil, fmt.Errorf("user lookup following request: %w", err)
	}
	req.Header.Add("Accept", "application/json")
	u.Authorizer.Add(req)
	tweetOpts.addQuery(req)

	resp, err := u.Client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("user lookup response: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("user lookup following reading body: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		e := &TweetErrorResponse{}
		if err := json.Unmarshal(body, e); err != nil {
			return nil, &HTTPError{
				Status:     resp.Status,
				StatusCode: resp.StatusCode,
				URL:        resp.Request.URL.String(),
			}
		}
		e.StatusCode = resp.StatusCode
		return nil, e
	}

	result := &UserTimeline{}
	if err := json.Unmarshal(body, result); err != nil {
		return nil, fmt.Errorf("user tweet timeline response decode: %w", err)
	}
	return result, nil
}
