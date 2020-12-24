package twitter

// Exclude used in the timeline parameters
type Exclude string

const (
	// ExcludeRetweets will exclude the tweet retweets
	ExcludeRetweets Exclude = "retweets"
	// ExcludeReplies will exclude the tweet replies
	ExcludeReplies Exclude = "replies"
)

func excludetringArray(arr []Exclude) []string {
	strs := make([]string, len(arr))
	for i, field := range arr {
		strs[i] = string(field)
	}
	return strs
}
