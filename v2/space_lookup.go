package twitter

import (
	"net/http"
	"strings"
)

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

type SpacesLookupResponse struct {
	Raw       *SpacesRaw
	RateLimit *RateLimit
}

type spaceRaw struct {
	Space    *SpaceObj          `json:"data"`
	Includes *SpacesRawIncludes `json:"includes,omitempty"`
	Errors   []*ErrorObj        `json:"errors,omitempty"`
}

type SpacesRaw struct {
	Spaces   []*SpaceObj        `json:"data"`
	Includes *SpacesRawIncludes `json:"includes,omitempty"`
	Errors   []*ErrorObj        `json:"errors,omitempty"`
}

type SpacesRawIncludes struct {
	Users  []*UserObj  `json:"users,omitempty"`
	Topics []*TopicObj `json:"topics,omitempty"`
}

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

type SpacesByCreatorLookupResponse struct {
	Raw       *SpacesRaw
	Meta      *SpacesByCreatorMeta `json:"meta"`
	RateLimit *RateLimit
}

type SpacesByCreatorMeta struct {
	ResultCount int `json:"result_count"`
}

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

type SpaceBuyersLookupResponse struct {
	Raw       *UserRaw
	RateLimit *RateLimit
}

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

type SpaceTweetsLookupResponse struct {
	Raw       *TweetRaw
	Meta      *SpaceTweetsLookupMeta `json:"meta"`
	RateLimit *RateLimit
}

type SpaceTweetsLookupMeta struct {
	ResultCount int `json:"result_count"`
}
