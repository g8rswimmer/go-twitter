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
	tweetMaxIDs         = 100
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
