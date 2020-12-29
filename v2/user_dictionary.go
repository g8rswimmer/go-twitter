package twitter

type UserDictionary struct {
	User        UserObj
	PinnedTweet *TweetObj
}

func CreateUserDictionary(user UserObj, includes *UserRawIncludes) *UserDictionary {
	dictionary := &UserDictionary{
		User: user,
	}
	if includes == nil {
		return dictionary
	}

	tweets := includes.TweetsByID()
	if tweet, has := tweets[user.PinnedTweetID]; has {
		dictionary.PinnedTweet = tweet
	}

	return dictionary
}
