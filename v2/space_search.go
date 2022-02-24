package twitter

import (
	"net/http"
	"strings"
)

// SpacesSearchOpts are the space search options
type SpacesSearchOpts struct {
	Expansions  []Expansion
	SpaceFields []SpaceField
	TopicFields []TopicField
	UserFields  []UserField
	State       SpaceState
}

func (s SpacesSearchOpts) addQuery(req *http.Request) {
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
	if len(s.State) > 0 {
		q.Add("state", string(s.State))
	}
	if len(q) > 0 {
		req.URL.RawQuery = q.Encode()
	}
}

// SpacesSearchResponse is the respones from the search
type SpacesSearchResponse struct {
	Raw       *SpacesRaw
	Meta      *SpacesSearchMeta `json:"meta"`
	RateLimit *RateLimit
}

// SpacesSearchMeta is the search meta
type SpacesSearchMeta struct {
	ResultCount int `json:"result_count"`
}
