package twitter

// TweetDictionary contains all of the information from a tweet callout
type TweetDictionary struct {
	Tweets   []*TweetObj              `json:"data"`
	Includes *TweetDictionaryIncludes `json:"includes,omitempty"`
	Errors   []*ErrorObj              `json:"errors,omitempty"`
}

// TweetDictionaryIncludes contains any additional information from the tweet callout
type TweetDictionaryIncludes struct {
	Tweets []*TweetObj `json:"tweets,omitempty"`
	Users  []*UserObj  `json:"users,omitempty"`
	Places []*PlaceObj `json:"places,omitempty"`
	Media  []*MediaObj `json:"media,omitempty"`
	Polls  []*PollObj  `json:"polls,omitempty"`
}

type tweetdictionary struct {
	Tweet    *TweetObj                `json:"data"`
	Includes *TweetDictionaryIncludes `json:"includes"`
	Errors   []*ErrorObj              `json:"errors"`
}
