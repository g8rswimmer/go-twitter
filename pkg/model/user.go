package model

// User contains Twitter user account metadata describing the referenced user
type User struct {
	ID              string      `json:"id"`
	Name            string      `json:"name"`
	UserName        string      `json:"username"`
	CreatedAt       string      `json:"created_at"`
	Description     string      `json:"description"`
	Entities        Entities    `json:"entities"`
	Location        string      `json:"location"`
	PinnedTweetID   string      `json:"pinned_tweet_id"`
	ProfileImageURL string      `json:"profile_image_url"`
	Protected       bool        `json:"protected"`
	PublicMetrics   UserMetrics `json:"public_metrics"`
	URL             string      `json:"url"`
	Verified        bool        `json:"verified"`
	WithHeld        WithHeld    `json:"withheld"`
}

// UserMetrics contains details about activity for this user
type UserMetrics struct {
	Followers int `json:"followers_count"`
	Following int `json:"following_count"`
	Tweets    int `json:"tweet_count"`
	Listed    int `json:"listed_count"`
}
