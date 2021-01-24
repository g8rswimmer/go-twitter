package twitter

// TweetLookupResponse contains all of the information from a tweet lookup callout
type TweetLookupResponse struct {
	Raw *TweetRaw
}

// UserMentionTimelineResponse contains the information from the user mention timelint callout
type UserMentionTimelineResponse struct {
	Raw  *TweetRaw
	Meta *UserTimelineMeta `json:"meta"`
}

// UserTweetTimelineResponse contains the information from the user tweet timeline callout
type UserTweetTimelineResponse struct {
	Raw  *TweetRaw
	Meta *UserTimelineMeta `json:"meta"`
}

// UserTimelineMeta contains the meta data from the timeline callout
type UserTimelineMeta struct {
	ResultCount   int    `json:"result_count"`
	NewestID      string `json:"newest_id"`
	OldestID      string `json:"oldest_id"`
	NextToken     string `json:"next_token"`
	PreviousToken string `json:"previous_token"`
}

// TweetRecentSearchResponse contains all of the information from a tweet recent search
type TweetRecentSearchResponse struct {
	Raw  *TweetRaw
	Meta *TweetRecentSearchMeta `json:"meta"`
}

// TweetRecentSearchMeta contains the recent search information
type TweetRecentSearchMeta struct {
	NewestID    string `json:"newest_id"`
	OldestID    string `json:"oldest_id"`
	ResultCount int    `json:"result_count"`
	NextToken   string `json:"next_token"`
}

type tweetraw struct {
	Tweet    *TweetObj         `json:"data"`
	Includes *TweetRawIncludes `json:"includes"`
	Errors   []*ErrorObj       `json:"errors"`
}

// TweetRaw is the raw response from the tweet lookup endpoint
type TweetRaw struct {
	Tweets       []*TweetObj       `json:"data"`
	Includes     *TweetRawIncludes `json:"includes,omitempty"`
	Errors       []*ErrorObj       `json:"errors,omitempty"`
	dictionaries map[string]*TweetDictionary
}

// TweetDictionaries create a map of tweet dictionaries from the raw tweet response
func (t *TweetRaw) TweetDictionaries() map[string]*TweetDictionary {
	if t.dictionaries != nil {
		return t.dictionaries
	}

	t.dictionaries = map[string]*TweetDictionary{}
	for _, tweet := range t.Tweets {
		t.dictionaries[tweet.ID] = CreateTweetDictionary(*tweet, t.Includes)
	}
	return t.dictionaries
}

// TweetRawIncludes contains any additional information from the tweet callout
type TweetRawIncludes struct {
	Tweets          []*TweetObj `json:"tweets,omitempty"`
	Users           []*UserObj  `json:"users,omitempty"`
	Places          []*PlaceObj `json:"places,omitempty"`
	Media           []*MediaObj `json:"media,omitempty"`
	Polls           []*PollObj  `json:"polls,omitempty"`
	userIDs         map[string]*UserObj
	userNames       map[string]*UserObj
	pollIDs         map[string]*PollObj
	mediaKeys       map[string]*MediaObj
	placeIDs        map[string]*PlaceObj
	referenceTweets map[string]*TweetObj
}

// UsersByID will return a map of user ids to object
func (t *TweetRawIncludes) UsersByID() map[string]*UserObj {
	switch {
	case t.userIDs == nil:
		return t.usersByID()
	default:
		return t.userIDs
	}
}

func (t *TweetRawIncludes) usersByID() map[string]*UserObj {
	t.userIDs = map[string]*UserObj{}
	for _, user := range t.Users {
		t.userIDs[user.ID] = user
	}
	return t.userIDs
}

// UsersByUserName will return a map of user names to object
func (t *TweetRawIncludes) UsersByUserName() map[string]*UserObj {
	switch {
	case t.userNames == nil:
		return t.usersByUserName()
	default:
		return t.userNames
	}
}

func (t *TweetRawIncludes) usersByUserName() map[string]*UserObj {
	t.userNames = map[string]*UserObj{}
	for _, user := range t.Users {
		t.userNames[user.UserName] = user
	}
	return t.userNames
}

// PollsByID will return a map of poll ids to object
func (t *TweetRawIncludes) PollsByID() map[string]*PollObj {
	switch {
	case t.pollIDs == nil:
		return t.pollsByID()
	default:
		return t.pollIDs
	}
}

func (t *TweetRawIncludes) pollsByID() map[string]*PollObj {

	t.pollIDs = map[string]*PollObj{}
	for _, poll := range t.Polls {
		t.pollIDs[poll.ID] = poll
	}
	return t.pollIDs
}

// MediaByKeys will return a map of media keys to object
func (t *TweetRawIncludes) MediaByKeys() map[string]*MediaObj {
	switch {
	case t.mediaKeys == nil:
		return t.mediaByKeys()
	default:
		return t.mediaKeys
	}
}

func (t *TweetRawIncludes) mediaByKeys() map[string]*MediaObj {
	t.mediaKeys = map[string]*MediaObj{}
	for _, m := range t.Media {
		t.mediaKeys[m.Key] = m
	}
	return t.mediaKeys
}

// PlacesByID will return a map of place ids to object
func (t *TweetRawIncludes) PlacesByID() map[string]*PlaceObj {
	switch {
	case t.placeIDs == nil:
		return t.placesByID()
	default:
		return t.placeIDs
	}
}

func (t *TweetRawIncludes) placesByID() map[string]*PlaceObj {
	t.placeIDs = map[string]*PlaceObj{}
	for _, place := range t.Places {
		t.placeIDs[place.ID] = place
	}
	return t.placeIDs
}

// TweetsByID will return a map of tweet ids to object
func (t *TweetRawIncludes) TweetsByID() map[string]*TweetObj {
	switch {
	case t.referenceTweets == nil:
		return t.tweetsByID()
	default:
		return t.referenceTweets
	}
}

func (t *TweetRawIncludes) tweetsByID() map[string]*TweetObj {
	t.referenceTweets = map[string]*TweetObj{}
	for _, tweet := range t.Tweets {
		t.referenceTweets[tweet.ID] = tweet
	}
	return t.referenceTweets
}
