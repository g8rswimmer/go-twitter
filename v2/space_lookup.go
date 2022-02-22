package twitter

import (
	"net/http"
	"strings"
)

// SpacesLookupOpts are the options for the space lookup
type SpacesLookupOpts struct {
	Expansions  []Expansion
	SpaceFields []SpaceField
	TopicFields []TopicField
	UserFields  []UserField
}

func (s SpacesLookupOpts) addQuery(req *http.Request) {
	q := req.URL.Query()
	if len(s.Expansions) > 0 {
		q.Add("expansions", strings.Join(expansionStringArray(s.Expansions), ","))
	}
	if len(s.SpaceFields) > 0 {
		q.Add("space.fields", strings.Join(spaceFieldStringArray(s.SpaceFields), ","))
	}
	if len(s.TopicFields) > 0 {
		q.Add("topic.fields", strings.Join(topicFieldStringArray(s.TopicFields), ","))
	}
	if len(s.UserFields) > 0 {
		q.Add("user.fields", strings.Join(userFieldStringArray(s.UserFields), ","))
	}
	if len(q) > 0 {
		req.URL.RawQuery = q.Encode()
	}
}

// SpacesLookupResponse is the response for the space lookup
type SpacesLookupResponse struct {
	Raw       *SpacesRaw
	RateLimit *RateLimit
}

type spaceRaw struct {
	Space    *SpaceObj          `json:"data"`
	Includes *SpacesRawIncludes `json:"includes,omitempty"`
	Errors   []*ErrorObj        `json:"errors,omitempty"`
}

// SpacesRaw the raw space objects
type SpacesRaw struct {
	Spaces   []*SpaceObj        `json:"data"`
	Includes *SpacesRawIncludes `json:"includes,omitempty"`
	Errors   []*ErrorObj        `json:"errors,omitempty"`
}

// SpacesRawIncludes are the includes for a space
type SpacesRawIncludes struct {
	Users  []*UserObj  `json:"users,omitempty"`
	Topics []*TopicObj `json:"topics,omitempty"`
}

// SpacesByCreatorLookupOpts are the options for the space by creator
type SpacesByCreatorLookupOpts struct {
	Expansions  []Expansion
	SpaceFields []SpaceField
	TopicFields []TopicField
	UserFields  []UserField
}

func (s SpacesByCreatorLookupOpts) addQuery(req *http.Request) {
	q := req.URL.Query()
	if len(s.Expansions) > 0 {
		q.Add("expansions", strings.Join(expansionStringArray(s.Expansions), ","))
	}
	if len(s.SpaceFields) > 0 {
		q.Add("space.fields", strings.Join(spaceFieldStringArray(s.SpaceFields), ","))
	}
	if len(s.TopicFields) > 0 {
		q.Add("topic.fields", strings.Join(topicFieldStringArray(s.TopicFields), ","))
	}
	if len(s.UserFields) > 0 {
		q.Add("user.fields", strings.Join(userFieldStringArray(s.UserFields), ","))
	}
	if len(q) > 0 {
		req.URL.RawQuery = q.Encode()
	}
}

// SpacesByCreatorLookupResponse is the response to the space by creator
type SpacesByCreatorLookupResponse struct {
	Raw       *SpacesRaw
	Meta      *SpacesByCreatorMeta `json:"meta"`
	RateLimit *RateLimit
}

// SpacesByCreatorMeta the meta for the space by creator
type SpacesByCreatorMeta struct {
	ResultCount int `json:"result_count"`
}

// SpaceBuyersLookupOpts are the options for the space buyer lookup
type SpaceBuyersLookupOpts struct {
	Expansions  []Expansion
	TweetFields []TweetField
	UserFields  []UserField
	MediaFields []MediaField
	PlaceFields []PlaceField
	PollFields  []PollField
}

func (s SpaceBuyersLookupOpts) addQuery(req *http.Request) {
	q := req.URL.Query()
	if len(s.Expansions) > 0 {
		q.Add("expansions", strings.Join(expansionStringArray(s.Expansions), ","))
	}
	if len(s.TweetFields) > 0 {
		q.Add("tweet.fields", strings.Join(tweetFieldStringArray(s.TweetFields), ","))
	}
	if len(s.UserFields) > 0 {
		q.Add("user.fields", strings.Join(userFieldStringArray(s.UserFields), ","))
	}
	if len(s.MediaFields) > 0 {
		q.Add("media.fields", strings.Join(mediaFieldStringArray(s.MediaFields), ","))
	}
	if len(s.PlaceFields) > 0 {
		q.Add("place.fields", strings.Join(placeFieldStringArray(s.PlaceFields), ","))
	}
	if len(s.PollFields) > 0 {
		q.Add("poll.fields", strings.Join(pollFieldStringArray(s.PollFields), ","))
	}
	if len(q) > 0 {
		req.URL.RawQuery = q.Encode()
	}
}

// SpaceBuyersLookupResponse is the space buyers lookup response
type SpaceBuyersLookupResponse struct {
	Raw       *UserRaw
	RateLimit *RateLimit
}

// SpaceTweetsLookupOpts are the options for the space tweets lookup
type SpaceTweetsLookupOpts struct {
	Expansions  []Expansion
	MediaFields []MediaField
	PlaceFields []PlaceField
	PollFields  []PollField
	TweetFields []TweetField
	UserFields  []UserField
}

func (s SpaceTweetsLookupOpts) addQuery(req *http.Request) {
	q := req.URL.Query()
	if len(s.Expansions) > 0 {
		q.Add("expansions", strings.Join(expansionStringArray(s.Expansions), ","))
	}
	if len(s.MediaFields) > 0 {
		q.Add("media.fields", strings.Join(mediaFieldStringArray(s.MediaFields), ","))
	}
	if len(s.PlaceFields) > 0 {
		q.Add("place.fields", strings.Join(placeFieldStringArray(s.PlaceFields), ","))
	}
	if len(s.PollFields) > 0 {
		q.Add("poll.fields", strings.Join(pollFieldStringArray(s.PollFields), ","))
	}
	if len(s.TweetFields) > 0 {
		q.Add("tweet.fields", strings.Join(tweetFieldStringArray(s.TweetFields), ","))
	}
	if len(s.UserFields) > 0 {
		q.Add("user.fields", strings.Join(userFieldStringArray(s.UserFields), ","))
	}
	if len(q) > 0 {
		req.URL.RawQuery = q.Encode()
	}
}

// SpaceTweetsLookupResponse is the response for the space tweets lookup
type SpaceTweetsLookupResponse struct {
	Raw       *TweetRaw
	Meta      *SpaceTweetsLookupMeta `json:"meta"`
	RateLimit *RateLimit
}

// SpaceTweetsLookupMeta is the space tweets lookup meta
type SpaceTweetsLookupMeta struct {
	ResultCount int `json:"result_count"`
}
