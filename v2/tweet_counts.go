package twitter

// TweetRecentCountsResponse contains all of the information from a tweet recent counts
type TweetRecentCountsResponse struct {
	TweetCounts []*TweetCount          `json:"data"`
	Meta        *TweetRecentCountsMeta `json:"meta"`
	RateLimit   *RateLimit
}

// TweetRecentCountsMeta contains the meta data from the recent counts information
type TweetRecentCountsMeta struct {
	TotalTweetCount int `json:"total_tweet_count"`
}

// TweetCount is the object on the tweet counts endpoints
type TweetCount struct {
	Start      string `json:"start"`
	End        string `json:"end"`
	TweetCount int    `json:"tweet_count"`
}

// TweetAllCountsResponse contain all fo the information from a tweet all counts
type TweetAllCountsResponse struct {
	TweetCounts []*TweetCount       `json:"data"`
	Meta        *TweetAllCountsMeta `json:"meta"`
	RateLimit   *RateLimit
}

// TweetAllCountsMeta is the meta data from the all counts information
type TweetAllCountsMeta struct {
	TotalTweetCount int    `json:"total_tweet_count"`
	NextToken       string `json:"next_token"`
}
