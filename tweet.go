package twitter

// TweetField defines the fields of the basic building block of all things twitter
type TweetField string

const (
	// TweetFieldID is the unique identifier of the requested Tweet.
	TweetFieldID TweetField = "id"
	// TweetFieldText is the actual UTF-8 text of the Tweet. See twitter-text for details on what characters are currently considered valid.
	TweetFieldText TweetField = "text"
	// TweetFieldAttachments specifies the type of attachments (if any) present in this Tweet.
	TweetFieldAttachments TweetField = "attachments"
	// TweetFieldAuthorID is the unique identifier of the User who posted this Tweet
	TweetFieldAuthorID TweetField = "author_id"
	// TweetFieldContextAnnotations contains context annotations for the Tweet.
	TweetFieldContextAnnotations TweetField = "context_annotations"
	// TweetFieldConversationID is the Tweet ID of the original Tweet of the conversation (which includes direct replies, replies of replies).
	TweetFieldConversationID TweetField = "conversation_id"
	// TweetFieldCreatedAt is the creation time of the Tweet.
	TweetFieldCreatedAt TweetField = "created_at"
	// TweetFieldEntities are the entities which have been parsed out of the text of the Tweet. Additionally see entities in Twitter Objects.
	TweetFieldEntities TweetField = "entities"
	// TweetFieldGeo contains details about the location tagged by the user in this Tweet, if they specified one.
	TweetFieldGeo TweetField = "geo"
	// TweetFieldInReplyToUserID will contain the original Tweet’s author ID
	TweetFieldInReplyToUserID TweetField = "in_reply_to_user_id"
	// TweetFieldLanguage is the language of the Tweet, if detected by Twitter. Returned as a BCP47 language tag.
	TweetFieldLanguage TweetField = "lang"
	// TweetFieldNonPublicMetrics are the non-public engagement metrics for the Tweet at the time of the request.
	TweetFieldNonPublicMetrics TweetField = "non_public_metrics"
	// TweetFieldPublicMetrics are the public engagement metrics for the Tweet at the time of the request.
	TweetFieldPublicMetrics TweetField = "public_metrics"
	// TweetFieldOrganicMetrics are the engagement metrics, tracked in an organic context, for the Tweet at the time of the request.
	TweetFieldOrganicMetrics TweetField = "organic_metrics"
	// TweetFieldPromotedMetrics are the engagement metrics, tracked in a promoted context, for the Tweet at the time of the request.
	TweetFieldPromotedMetrics TweetField = "promoted_metrics"
	// TweetFieldPossiblySensitve is an indicator that the URL contained in the Tweet may contain content or media identified as sensitive content.
	TweetFieldPossiblySensitve TweetField = "possibly_sensitive"
	// TweetFieldReferencedTweets is a list of Tweets this Tweet refers to.
	TweetFieldReferencedTweets TweetField = "referenced_tweets"
	// TweetFieldSource is the name of the app the user Tweeted from.
	TweetFieldSource TweetField = "source"
	// TweetFieldWithHeld contains withholding details
	TweetFieldWithHeld TweetField = "withheld"
)

// Tweet is the primary object on the tweets endpoints
type Tweet struct {
	ID                 string                   `json:"id"`
	Text               string                   `json:"text"`
	Attachments        TweetAttachments         `json:"attachments"`
	AuthorID           string                   `json:"author_id"`
	ContextAnnotations []TweetContextAnnotation `json:"context_annotations"`
	ConversationID     string                   `json:"conversation_id"`
	CreatedAt          string                   `json:"created_at"`
	Entities           EntitiesObj              `json:"entities"`
	Geo                TweetGeo                 `json:"geo"`
	InReplyToUserID    string                   `json:"in_reply_to_user_id"`
	Language           string                   `json:"lang"`
	NonPublicMetrics   TweetMetrics             `json:"non_public_metrics"`
	OrganicMetrics     TweetMetrics             `json:"organic_metrics"`
	PossibySensitive   bool                     `json:"possiby_sensitive"`
	PromotedMetrics    TweetMetrics             `json:"promoted_metrics"`
	PublicMetrics      TweetMetrics             `json:"public_metrics"`
	ReferencedTweets   []TweetReferencedTweet   `json:"referenced_tweets"`
	Source             string                   `json:"source"`
	WithHeld           WithHeldObj              `json:"withheld"`
}

// TweetAttachments specifics the type of attachment present in the tweet
type TweetAttachments struct {
	MediaKeys []string `json:"media_keys"`
	PollIDs   []string `json:"poll_ids"`
}

// TweetContextAnnotation contain the context annotation
type TweetContextAnnotation struct {
	Domain TweetContext `json:"domain"`
	Entity TweetContext `json:"entity"`
}

// TweetContext contains the elements which identify detailed information regarding the domain classificaiton based on the Tweet text
type TweetContext struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

// TweetGeo contains details about the location tagged by the user in this Tweet, if they specified one.
type TweetGeo struct {
	PlaceID     string              `json:"place_id"`
	Coordinates TweetGeoCoordinates `json:"coordinates"`
}

// TweetGeoCoordinates contains details about the coordinates of the location tagged by the user in this Tweet, if they specified one.
type TweetGeoCoordinates struct {
	Type        string    `json:"type"`
	Coordinates []float64 `json:"coordinates"`
}

// TweetMetrics engagement metrics for the Tweet at the time of the request
type TweetMetrics struct {
	Impressions       int `json:"impression_count"`
	URLLinkClicks     int `json:"url_link_clicks"`
	UserProfileClicks int `json:"user_profile_clicks"`
	Likes             int `json:"like_count"`
	Replies           int `json:"reply_count"`
	Retweets          int `json:"retweet_count"`
	Quotes            int `json:"quote_count"`
}

// TweetReferencedTweet is a Tweet this Tweet refers to
type TweetReferencedTweet struct {
	Type string `json:"type"`
	ID   string `json:"id"`
}
