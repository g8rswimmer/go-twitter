package twitter

// UserField defines the twitter user account metadata fields
type UserField string

const (
	// UserFieldCreatedAt is the UTC datetime that the user account was created on Twitter.
	UserFieldCreatedAt UserField = "created_at"
	// UserFieldDescription is the text of this user's profile description (also known as bio), if the user provided one.
	UserFieldDescription UserField = "description"
	// UserFieldEntities contains details about text that has a special meaning in the user's description.
	UserFieldEntities UserField = "entities"
	// UserFieldID is the unique identifier of this user.
	UserFieldID UserField = "id"
	// UserFieldLocation is the location specified in the user's profile, if the user provided one.
	UserFieldLocation UserField = "location"
	// UserFieldName is the name of the user, as they’ve defined it on their profile
	UserFieldName UserField = "name"
	// UserFieldPinnedTweetID is the unique identifier of this user's pinned Tweet.
	UserFieldPinnedTweetID UserField = "pinned_tweet_id"
	// UserFieldProfileImageURL is the URL to the profile image for this user, as shown on the user's profile.
	UserFieldProfileImageURL UserField = "profile_image_url"
	// UserFieldProtected indicates if this user has chosen to protect their Tweets (in other words, if this user's Tweets are private).
	UserFieldProtected UserField = "protected"
	// UserFieldPublicMetrics contains details about activity for this user.
	UserFieldPublicMetrics UserField = "public_metrics"
	// UserFieldURL is the URL specified in the user's profile, if present.
	UserFieldURL UserField = "url"
	// UserFieldUserName is the Twitter screen name, handle, or alias that this user identifies themselves with
	UserFieldUserName UserField = "username"
	// UserFieldVerified indicates if this user is a verified Twitter User.
	UserFieldVerified UserField = "verified"
	// UserFieldWithHeld contains withholding details
	UserFieldWithHeld UserField = "withheld"
)

func userFieldStringArray(arr []UserField) []string {
	strs := make([]string, len(arr))
	for i, field := range arr {
		strs[i] = string(field)
	}
	return strs
}

// UserObj contains Twitter user account metadata describing the referenced user
type UserObj struct {
	ID              string          `json:"id"`
	Name            string          `json:"name"`
	UserName        string          `json:"username"`
	CreatedAt       string          `json:"created_at,omitempty"`
	Description     string          `json:"description,omitempty"`
	Entities        *EntitiesObj    `json:"entities,omitempty"`
	Location        string          `json:"location,omitempty"`
	PinnedTweetID   string          `json:"pinned_tweet_id,omitempty"`
	ProfileImageURL string          `json:"profile_image_url,omitempty"`
	Protected       bool            `json:"protected,omitempty"`
	PublicMetrics   *UserMetricsObj `json:"public_metrics,omitempty"`
	URL             string          `json:"url,omitempty"`
	Verified        bool            `json:"verified,omitempty"`
	WithHeld        *WithHeldObj    `json:"withheld,omitempty"`
}

// UserMetricsObj contains details about activity for this user
type UserMetricsObj struct {
	Followers int `json:"followers_count"`
	Following int `json:"following_count"`
	Tweets    int `json:"tweet_count"`
	Listed    int `json:"listed_count"`
}
