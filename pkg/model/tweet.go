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
	Entities           Entities                 `json:"entities"`
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
	WithHeld           WithHeld                 `json:"withheld"`
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
