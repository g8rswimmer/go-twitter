package twitter

// Exclude are the exclusions in the request
type Exclude string

const (
	// ExcludeRetweets will exclude a tweet's retweets
	ExcludeRetweets Exclude = "retweets"
	// ExcludeReplies will exclude a tweet's replies
	ExcludeReplies Exclude = "replies"
)

func excludeStringArray(arr []Exclude) []string {
	strs := make([]string, len(arr))
	for i, field := range arr {
		strs[i] = string(field)
	}
	return strs
}
