package twitter

// UserLookupResponse contains all of the information from an user lookup callout
type UserLookupResponse struct {
	Raw *UserRaw
}

// UserFollowingLookupResponse is the response for the user following API
type UserFollowingLookupResponse struct {
	Raw  *UserRaw
	Meta *UserFollowinghMeta `json:"meta"`
}

// UserFollowinghMeta is the meta data returned by the user following API
type UserFollowinghMeta struct {
	ResultCount   int    `json:"result_count"`
	NextToken     string `json:"next_token"`
	PreviousToken string `json:"previous_token"`
}

// UserFollowersLookupResponse is the response for the user followers API
type UserFollowersLookupResponse struct {
	Raw  *UserRaw
	Meta *UserFollowershMeta `json:"meta"`
}

// UserFollowershMeta is the meta data returned by the user followers API
type UserFollowershMeta struct {
	ResultCount   int    `json:"result_count"`
	NextToken     string `json:"next_token"`
	PreviousToken string `json:"previous_token"`
}

type userraw struct {
	User     *UserObj         `json:"data"`
	Includes *UserRawIncludes `json:"includes"`
	Errors   []*ErrorObj      `json:"errors"`
}

// UserRaw is the raw response from the user lookup endpoint
type UserRaw struct {
	Users        []*UserObj       `json:"data"`
	Includes     *UserRawIncludes `json:"includes,omitempty"`
	Errors       []*ErrorObj      `json:"errors,omitempty"`
	dictionaries map[string]*UserDictionary
}

// UserDictionaries create a map of user dictionaries from the raw user response
func (u *UserRaw) UserDictionaries() map[string]*UserDictionary {
	if u.dictionaries != nil {
		return u.dictionaries
	}

	u.dictionaries = map[string]*UserDictionary{}
	for _, user := range u.Users {
		u.dictionaries[user.ID] = CreateUserDictionary(*user, u.Includes)
	}
	return u.dictionaries
}

// UserRawIncludes contains any additional information from the user callout
type UserRawIncludes struct {
	Tweets       []*TweetObj `json:"tweets,omitempty"`
	pinnedTweets map[string]*TweetObj
}

// TweetsByID will return a map of tweet ids to object
func (u *UserRawIncludes) TweetsByID() map[string]*TweetObj {
	switch {
	case u.pinnedTweets == nil:
		return u.tweetsByID()
	default:
		return u.pinnedTweets
	}
}

func (u *UserRawIncludes) tweetsByID() map[string]*TweetObj {
	u.pinnedTweets = map[string]*TweetObj{}
	for _, tweet := range u.Tweets {
		u.pinnedTweets[tweet.ID] = tweet
	}
	return u.pinnedTweets
}
