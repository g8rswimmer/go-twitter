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

func tweetFieldStringArray(arr []TweetField) []string {
	strs := make([]string, len(arr))
	for i, field := range arr {
		strs[i] = string(field)
	}
	return strs
}

// TweetObj is the primary object on the tweets endpoints
type TweetObj struct {
	ID                 string                       `json:"id"`
	Text               string                       `json:"text"`
	Attachments        *TweetAttachmentsObj         `json:"attachments,omitempty"`
	AuthorID           string                       `json:"author_id,omitempty"`
	ContextAnnotations []*TweetContextAnnotationObj `json:"context_annotations,omitempty"`
	ConversationID     string                       `json:"conversation_id,omitempty"`
	CreatedAt          string                       `json:"created_at,omitempty"`
	Entities           *EntitiesObj                 `json:"entities,omitempty"`
	Geo                *TweetGeoObj                 `json:"geo,omitempty"`
	InReplyToUserID    string                       `json:"in_reply_to_user_id,omitempty"`
	Language           string                       `json:"lang,omitempty"`
	NonPublicMetrics   *TweetMetricsObj             `json:"non_public_metrics,omitempty"`
	OrganicMetrics     *TweetMetricsObj             `json:"organic_metrics,omitempty"`
	PossibySensitive   bool                         `json:"possiby_sensitive,omitempty"`
	PromotedMetrics    *TweetMetricsObj             `json:"promoted_metrics,omitempty"`
	PublicMetrics      *TweetMetricsObj             `json:"public_metrics,omitempty"`
	ReferencedTweets   []*TweetReferencedTweetObj   `json:"referenced_tweets,omitempty"`
	Source             string                       `json:"source,omitempty"`
	WithHeld           *WithHeldObj                 `json:"withheld,omitempty"`
}

// TweetAttachmentsObj specifics the type of attachment present in the tweet
type TweetAttachmentsObj struct {
	MediaKeys []string `json:"media_keys"`
	PollIDs   []string `json:"poll_ids"`
}

// TweetContextAnnotationObj contain the context annotation
type TweetContextAnnotationObj struct {
	Domain TweetContextObj `json:"domain"`
	Entity TweetContextObj `json:"entity"`
}

// TweetContextObj contains the elements which identify detailed information regarding the domain classificaiton based on the Tweet text
type TweetContextObj struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

// TweetGeoObj contains details about the location tagged by the user in this Tweet, if they specified one.
type TweetGeoObj struct {
	PlaceID     string                 `json:"place_id"`
	Coordinates TweetGeoCoordinatesObj `json:"coordinates"`
}

// TweetGeoCoordinatesObj contains details about the coordinates of the location tagged by the user in this Tweet, if they specified one.
type TweetGeoCoordinatesObj struct {
	Type        string    `json:"type"`
	Coordinates []float64 `json:"coordinates"`
}

// TweetMetricsObj engagement metrics for the Tweet at the time of the request
type TweetMetricsObj struct {
	Impressions       int `json:"impression_count"`
	URLLinkClicks     int `json:"url_link_clicks"`
	UserProfileClicks int `json:"user_profile_clicks"`
	Likes             int `json:"like_count"`
	Replies           int `json:"reply_count"`
	Retweets          int `json:"retweet_count"`
	Quotes            int `json:"quote_count"`
}

// TweetReferencedTweetObj is a Tweet this Tweet refers to
type TweetReferencedTweetObj struct {
	Type string `json:"type"`
	ID   string `json:"id"`
}
