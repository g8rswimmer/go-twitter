package twitter

// Expansion can expand objects referenced in the payload
type Expansion string

const (
	// ExpansionAttachmentsPollIDs returns a poll object containing metadata for the poll included in the Tweet
	ExpansionAttachmentsPollIDs Expansion = "attachments.poll_ids"
	// ExpansionAttachmentsMediaKeys returns a media object representing the images, videos, GIFs included in the Tweet
	ExpansionAttachmentsMediaKeys Expansion = "attachments.media_keys"
	// ExpansionAuthorID returns a user object representing the Tweet’s author
	ExpansionAuthorID Expansion = "author_id"
	// ExpansionEntitiesMentionsUserName returns a user object for the user mentioned in the Tweet
	ExpansionEntitiesMentionsUserName Expansion = "entities.mentions.username"
	// ExpansionGeoPlaceID returns a place object containing metadata for the location tagged in the Tweet
	ExpansionGeoPlaceID Expansion = "geo.place_id"
	// ExpansionInReplyToUserID returns a user object representing the Tweet author this requested Tweet is a reply of
	ExpansionInReplyToUserID Expansion = "in_reply_to_user_id"
	// ExpansionReferencedTweetsID returns a Tweet object that this Tweet is referencing (either as a Retweet, Quoted Tweet, or reply)
	ExpansionReferencedTweetsID Expansion = "referenced_tweets.id"
	// ExpansionReferencedTweetsIDAuthorID returns a user object for the author of the referenced Tweet
	ExpansionReferencedTweetsIDAuthorID Expansion = "referenced_tweets.id.author_id"
	// ExpansionPinnedTweetID returns a Tweet object representing the Tweet pinned to the top of the user’s profile
	ExpansionPinnedTweetID Expansion = "pinned_tweet_id"
)

func expansionStringArray(arr []Expansion) []string {
	strs := make([]string, len(arr))
	for i, expansion := range arr {
		strs[i] = string(expansion)
	}
	return strs
}
