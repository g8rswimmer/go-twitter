package model

// Tweet is the primary object on the tweets endpoints
type Tweet struct {
	ID                 string                   `json:"id"`
	Text               string                   `json:"text"`
	Attachments        TweetAttachments         `json:"attachments"`
	AuthorID           string                   `json:"author_id"`
	ContextAnnotations []TweetContextAnnotation `json:"context_annotations"`
	ConversationID     string                   `json:"conversation_id"`
	CreatedAt          string                   `json:"created_at"`
	Entities           TweetEntities            `json:"entities"`
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
	WithHeld           TweetWithHeld            `json:"withheld"`
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

// TweetEntities contains details about text that has a special meaning in a Tweet.
type TweetEntities struct {
	Annotations []TweetEntityAnnotation `json:"annotations"`
	URLs        []TweetEntityURL        `json:"urls"`
	HashTags    []TweetEntityTag        `json:"hashtags"`
	Mentions    []TweetEntityMention    `json:"mentions"`
	CashTags    []TweetEntityTag        `json:"cashtags"`
}

// TweetEntity contains the start and end positions of the text used to annotate the Tweet
type TweetEntity struct {
	Start int `json:"start"`
	End   int `json:"end"`
}

// TweetEntityAnnotation contains details about annotations relative to the text within a Tweet.
type TweetEntityAnnotation struct {
	TweetEntity
	Probability    float64 `json:"probability"`
	Type           string  `json:"type"`
	NormalizedText string  `json:"normalized_text"`
}

// TweetEntityURL contains details about text recognized as a URL.
type TweetEntityURL struct {
	TweetEntity
	URL         string `json:"url"`
	ExpandedURL string `json:"expanded_url"`
	DisplayURL  string `json:"display_url"`
	Status      string `json:"status"`
	Title       string `json:"title"`
	Desription  string `json:"description"`
	UnwoundURL  string `json:"unwound_url"`
}

// TweetEntityTag contains details about text recognized as a tag
type TweetEntityTag struct {
	TweetEntity
	Tag string `json:"tag"`
}

// TweetEntityMention contains details about text recognized as a user mention.
type TweetEntityMention struct {
	TweetEntity
	UserName string `json:"username"`
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

// TweetWithHeld contains withholding details
type TweetWithHeld struct {
	Copyright    bool     `json:"copyright"`
	CountryCodes []string `json:"country_codes"`
}
