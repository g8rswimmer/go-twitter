package twitter

// TweetDictionary is a struct of a tweet and all of the reference objects
type TweetDictionary struct {
	Tweet            TweetObj
	Author           *UserObj
	InReplyUser      *UserObj
	Place            *PlaceObj
	AttachmentPolls  []*PollObj
	AttachmentMedia  []*MediaObj
	Mentions         []*TweetMention
	ReferencedTweets []*TweetReference
}

// TweetMention is the mention and the user associated with it
type TweetMention struct {
	Mention *EntityMentionObj
	User    *UserObj
}

// TweetReference is the tweet referenced and it's dictionary
type TweetReference struct {
	Reference       *TweetReferencedTweetObj
	TweetDictionary *TweetDictionary
}

// CreateTweetDictionary will create a dictionary from a tweet and the includes
func CreateTweetDictionary(tweet TweetObj, includes *TweetRawIncludes) *TweetDictionary {
	dictionary := &TweetDictionary{
		Tweet:            tweet,
		AttachmentMedia:  []*MediaObj{},
		AttachmentPolls:  []*PollObj{},
		Mentions:         []*TweetMention{},
		ReferencedTweets: []*TweetReference{},
	}
	if includes == nil {
		return dictionary
	}

	userIDs := includes.UsersByID()
	if user, has := userIDs[tweet.AuthorID]; has {
		dictionary.Author = user
	}
	if user, has := userIDs[tweet.InReplyToUserID]; has {
		dictionary.InReplyUser = user
	}

	if tweet.Entities != nil {
		userNames := includes.UsersByUserName()

		mentions := []*TweetMention{}
		for i, entity := range tweet.Entities.Mentions {
			if user, has := userNames[entity.UserName]; has {
				mention := &TweetMention{
					Mention: &tweet.Entities.Mentions[i],
					User:    user,
				}
				mentions = append(mentions, mention)
			}
		}
		dictionary.Mentions = mentions
	}

	if tweet.Attachments != nil {
		pollIDs := includes.PollsByID()

		attachmentPolls := []*PollObj{}
		for _, id := range tweet.Attachments.PollIDs {
			if poll, has := pollIDs[id]; has {
				attachmentPolls = append(attachmentPolls, poll)
			}
		}
		dictionary.AttachmentPolls = attachmentPolls

		mediaKeys := includes.MediaByKeys()

		attachmentMedia := []*MediaObj{}
		for _, key := range tweet.Attachments.MediaKeys {
			if media, has := mediaKeys[key]; has {
				attachmentMedia = append(attachmentMedia, media)
			}
		}
		dictionary.AttachmentMedia = attachmentMedia
	}
	if tweet.Geo != nil {
		placeIDs := includes.PlacesByID()
		if place, has := placeIDs[tweet.Geo.PlaceID]; has {
			dictionary.Place = place
		}
	}

	tweets := includes.TweetsByID()
	tweetReferences := []*TweetReference{}

	for i, rt := range tweet.ReferencedTweets {
		if t, has := tweets[rt.ID]; has {
			ref := &TweetReference{
				Reference:       tweet.ReferencedTweets[i],
				TweetDictionary: CreateTweetDictionary(*t, includes),
			}
			tweetReferences = append(tweetReferences, ref)
		}
	}
	dictionary.ReferencedTweets = tweetReferences

	return dictionary
}
