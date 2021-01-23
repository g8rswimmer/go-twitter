package twitter

type Exclude string

const (
	ExcludeRetweets Exclude = "retweets"
	ExcludeReplies  Exclude = "replies"
)

func excludeStringArray(arr []Exclude) []string {
	strs := make([]string, len(arr))
	for i, field := range arr {
		strs[i] = string(field)
	}
	return strs
}
