package twitter

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

const (
	userLookupEndpoint		= "2/users"
	userNameLookupEndpoint 	= "2/users/by/username"
	userNamesLookupEndpoint	= "2/users/by"
	userMaxIDs				= 100
	userMaxNames			= 100
)

// UserLookups is a map of user lookups
type UserLookups map[string]UserLookup

func (t UserLookups) lookup(decoder *json.Decoder) error {
	type include struct {
		Tweet  []*TweetObj  `json:"tweets"`
	}
	type body struct {
		Data    UserObj `json:"data"`
		Include include  `json:"includes"`
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
		Tweet  []*TweetObj  `json:"tweets"`
	}
	type body struct {
		Data    []UserObj `json:"data"`
		Include include    `json:"includes"`
	}
	b := &body{}
	if err := decoder.Decode(b); err != nil {
		return fmt.Errorf("tweet lookup decode error %w", err)
	}

	for i, user := range b.Data {
		ul := UserLookup{
			User: user,
		}
		if i < len(b.Include.Tweet) {
			ul.Tweet = b.Include.Tweet[i]
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
			return nil, fmt.Errorf("user lookup response error decode: %w", err)
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
			return nil, fmt.Errorf("user lookup response error decode: %w", err)
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