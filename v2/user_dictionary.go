package twitter

// UserDictionary is a struct of an user and all of the reference objects
type UserDictionary struct {
	User        UserObj
	PinnedTweet *TweetObj
}

// CreateUserDictionary will create a dictionary from a user and the includes
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
