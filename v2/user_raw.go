package twitter

type UserLookupResponse struct {
	Raw *UserRaw
}

type userraw struct {
	User     *UserObj         `json:"data"`
	Includes *UserRawIncludes `json:"includes"`
	Errors   []*ErrorObj      `json:"errors"`
}

type UserRaw struct {
	Users        []*UserObj       `json:"data"`
	Includes     *UserRawIncludes `json:"includes,omitempty"`
	Errors       []*ErrorObj      `json:"errors,omitempty"`
	dictionaries map[string]*UserDictionary
}

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

func (t *UserRawIncludes) tweetsByID() map[string]*TweetObj {
	t.pinnedTweets = map[string]*TweetObj{}
	for _, tweet := range t.Tweets {
		t.pinnedTweets[tweet.ID] = tweet
	}
	return t.pinnedTweets
}
