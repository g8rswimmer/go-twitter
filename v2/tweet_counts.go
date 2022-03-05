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

type TweetAllCountsResponse struct {
	TweetCounts []*TweetCount       `json:"data"`
	Meta        *TweetAllCountsMeta `json:"meta"`
	RateLimit   *RateLimit
}

type TweetAllCountsMeta struct {
	TotalTweetCount int    `json:"total_tweet_count"`
	NextToken       string `json:"next_token"`
}
