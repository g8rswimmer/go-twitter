package twitter

import (
	"net/http"
	"strconv"
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

//UserListLookupOpts are the response field options
type UserListLookupOpts struct {
	Expansions      []Expansion
	ListFields      []ListField
	UserFields      []UserField
	MaxResults      int
	PaginationToken string
}

func (l UserListLookupOpts) addQuery(req *http.Request) {
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
	if l.MaxResults > 0 {
		q.Add("max_results", strconv.Itoa(l.MaxResults))
	}
	if len(l.PaginationToken) > 0 {
		q.Add("pagination_token", l.PaginationToken)
	}
	if len(q) > 0 {
		req.URL.RawQuery = q.Encode()
	}
}

// UserListRaw is the raw response
type UserListRaw struct {
	Lists    []*ListObj       `json:"data"`
	Includes *ListRawIncludes `json:"includes,omitempty"`
	Errors   []*ErrorObj      `json:"errors,omitempty"`
}

// UserListLookupResponse is the raw ressponse with meta
type UserListLookupResponse struct {
	Raw  *UserListRaw
	Meta *UserListLookupMeta `json:"meta"`
}

// UserListLookupMeta is the meta data for the lists
type UserListLookupMeta struct {
	ResultCount   int    `json:"result_count"`
	PreviousToken string `json:"previous_token"`
	NextToken     string `json:"next_token"`
}

//ListTweetLookupOpts are the response field options
type ListTweetLookupOpts struct {
	Expansions      []Expansion
	TweetFields     []TweetField
	UserFields      []UserField
	MaxResults      int
	PaginationToken string
}

func (l ListTweetLookupOpts) addQuery(req *http.Request) {
	q := req.URL.Query()
	if len(l.Expansions) > 0 {
		q.Add("expansions", strings.Join(expansionStringArray(l.Expansions), ","))
	}
	if len(l.TweetFields) > 0 {
		q.Add("tweet.fields", strings.Join(tweetFieldStringArray(l.TweetFields), ","))
	}
	if len(l.UserFields) > 0 {
		q.Add("user.fields", strings.Join(userFieldStringArray(l.UserFields), ","))
	}
	if l.MaxResults > 0 {
		q.Add("max_results", strconv.Itoa(l.MaxResults))
	}
	if len(l.PaginationToken) > 0 {
		q.Add("pagination_token", l.PaginationToken)
	}
	if len(q) > 0 {
		req.URL.RawQuery = q.Encode()
	}
}

// ListTweetLookupResponse is the response to the list tweet lookup
type ListTweetLookupResponse struct {
	Raw  *TweetRaw
	Meta *ListTweetLookupMeta `json:"meta"`
}

// ListTweetLookupMeta is the meta data associated with the list tweet lookup
type ListTweetLookupMeta struct {
	ResultCount   int    `json:"result_count"`
	PreviousToken string `json:"previous_token"`
	NextToken     string `json:"next_token"`
}

type UserListMembershipsOpts struct {
	Expansions      []Expansion
	ListFields      []ListField
	UserFields      []UserField
	MaxResults      int
	PaginationToken string
}

func (l UserListMembershipsOpts) addQuery(req *http.Request) {
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
	if l.MaxResults > 0 {
		q.Add("max_results", strconv.Itoa(l.MaxResults))
	}
	if len(l.PaginationToken) > 0 {
		q.Add("pagination_token", l.PaginationToken)
	}
	if len(q) > 0 {
		req.URL.RawQuery = q.Encode()
	}
}

type UserListMembershipsRaw struct {
	Lists    []*ListObj       `json:"data"`
	Includes *ListRawIncludes `json:"includes,omitempty"`
	Errors   []*ErrorObj      `json:"errors,omitempty"`
}

type UserListMembershipsMeta struct {
	ResultCount   int    `json:"result_count"`
	PreviousToken string `json:"previous_token"`
	NextToken     string `json:"next_token"`
}

type UserListMembershipsResponse struct {
	Raw  *UserListMembershipsRaw
	Meta *UserListMembershipsMeta `json:"meta"`
}
