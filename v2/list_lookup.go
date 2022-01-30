package twitter

import (
	"net/http"
	"strings"
)

// ListLookupOpts are the options for the list lookup
type ListLookupOpts struct {
	Expansions []Expansion
	ListFields []ListField
	UserFields []UserField
}

func (l ListLookupOpts) addQuery(req *http.Request) {
	q := req.URL.Query()
	if len(l.Expansions) > 0 {
		q.Add("expansions", strings.Join(expansionStringArray(l.Expansions), ","))
	}
	if len(l.ListFields) > 0 {
		q.Add("list.fields", strings.Join(listFieldStringArray(l.ListFields), ","))
	}
	if len(l.UserFields) > 0 {
		q.Add("user.fields", strings.Join(userFieldStringArray(l.UserFields), ","))
	}
	if len(q) > 0 {
		req.URL.RawQuery = q.Encode()
	}
}

// ListRaw the raw list response
type ListRaw struct {
	List     *ListObj         `json:"data"`
	Includes *ListRawIncludes `json:"includes,omitempty"`
	Errors   []*ErrorObj      `json:"errors,omitempty"`
}

// ListRawIncludes the data include from the expansion
type ListRawIncludes struct {
	Users []*UserObj `json:"users,omitempty"`
}

// ListLookupResponse is the response from the list lookup
type ListLookupResponse struct {
	Raw *ListRaw
}
