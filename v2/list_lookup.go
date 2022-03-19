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
	Raw       *ListRaw
	RateLimit *RateLimit
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

// UserListLookupResponse is the raw response with meta
type UserListLookupResponse struct {
	Raw       *UserListRaw
	Meta      *UserListLookupMeta `json:"meta"`
	RateLimit *RateLimit
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
	Raw       *TweetRaw
	Meta      *ListTweetLookupMeta `json:"meta"`
	RateLimit *RateLimit
}

// ListTweetLookupMeta is the meta data associated with the list tweet lookup
type ListTweetLookupMeta struct {
	ResultCount   int    `json:"result_count"`
	PreviousToken string `json:"previous_token"`
	NextToken     string `json:"next_token"`
}

// UserListMembershipsOpts the user list member options
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

// UserListMembershipsRaw the raw data from the user list memberships
type UserListMembershipsRaw struct {
	Lists    []*ListObj       `json:"data"`
	Includes *ListRawIncludes `json:"includes,omitempty"`
	Errors   []*ErrorObj      `json:"errors,omitempty"`
}

// UserListMembershipsMeta the response meta data
type UserListMembershipsMeta struct {
	ResultCount   int    `json:"result_count"`
	PreviousToken string `json:"previous_token"`
	NextToken     string `json:"next_token"`
}

// UserListMembershipsResponse the user list membership response
type UserListMembershipsResponse struct {
	Raw       *UserListMembershipsRaw
	Meta      *UserListMembershipsMeta `json:"meta"`
	RateLimit *RateLimit
}

// ListUserMembersOpts is the list user member options
type ListUserMembersOpts struct {
	Expansions      []Expansion
	TweetFields     []TweetField
	UserFields      []UserField
	MaxResults      int
	PaginationToken string
}

func (l ListUserMembersOpts) addQuery(req *http.Request) {
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

// ListUserMembersMeta is the meta data of the response
type ListUserMembersMeta struct {
	ResultCount   int    `json:"result_count"`
	PreviousToken string `json:"previous_token"`
	NextToken     string `json:"next_token"`
}

// ListUserMembersResponse is the response to the list user members
type ListUserMembersResponse struct {
	Raw       *UserRaw
	Meta      *ListUserMembersMeta `json:"meta"`
	RateLimit *RateLimit
}

// UserPinnedListsOpts pinned list options
type UserPinnedListsOpts struct {
	Expansions []Expansion
	ListFields []ListField
	UserFields []UserField
}

func (l UserPinnedListsOpts) addQuery(req *http.Request) {
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

// UserPinnedListsResponse pinned list response
type UserPinnedListsResponse struct {
	Raw       *UserPinnedListsRaw
	Meta      *UserPinnedListsMeta `json:"meta"`
	RateLimit *RateLimit
}

// UserPinnedListsRaw the raw data for pinned lists
type UserPinnedListsRaw struct {
	Lists    []*ListObj       `json:"data"`
	Includes *ListRawIncludes `json:"includes,omitempty"`
	Errors   []*ErrorObj      `json:"errors,omitempty"`
}

// UserPinnedListsMeta the meta for pinned lists
type UserPinnedListsMeta struct {
	ResultCount int `json:"result_count"`
}

// UserFollowedListsOpts are the options for the user followed lists
type UserFollowedListsOpts struct {
	Expansions      []Expansion
	ListFields      []ListField
	UserFields      []UserField
	MaxResults      int
	PaginationToken string
}

func (l UserFollowedListsOpts) addQuery(req *http.Request) {
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

// UserFollowedListsResponse is the user followed response
type UserFollowedListsResponse struct {
	Raw       *UserFollowedListsRaw
	Meta      *UserFollowedListsMeta `json:"meta"`
	RateLimit *RateLimit
}

// UserFollowedListsRaw is the raw response for the user followed
type UserFollowedListsRaw struct {
	Lists    []*ListObj       `json:"data"`
	Includes *ListRawIncludes `json:"includes,omitempty"`
	Errors   []*ErrorObj      `json:"errors,omitempty"`
}

// UserFollowedListsMeta is the meta for the user followed
type UserFollowedListsMeta struct {
	ResultCount   int    `json:"result_count"`
	PreviousToken string `json:"previous_token"`
	NextToken     string `json:"next_token"`
}

// ListUserFollowersOpts is the list followers options
type ListUserFollowersOpts struct {
	Expansions      []Expansion
	TweetFields     []TweetField
	UserFields      []UserField
	MaxResults      int
	PaginationToken string
}

func (l ListUserFollowersOpts) addQuery(req *http.Request) {
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

// ListUserFollowersMeta is the meta for the list followers
type ListUserFollowersMeta struct {
	ResultCount   int    `json:"result_count"`
	PreviousToken string `json:"previous_token"`
	NextToken     string `json:"next_token"`
}

// ListUserFollowersResponse is the response for the list followers
type ListUserFollowersResponse struct {
	Raw       *UserRaw
	Meta      *ListUserFollowersMeta `json:"meta"`
	RateLimit *RateLimit
}
