package twitter

// TweetLookupResponse contains all of the information from a tweet lookup callout
type TweetLookupResponse struct {
	Raw *TweetLookupRaw
}

// TweetLookupRaw is the raw response from the tweet lookup endpoint
type TweetLookupRaw struct {
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

type tweetraw struct {
	Tweet    *TweetObj                `json:"data"`
	Includes *TweetDictionaryIncludes `json:"includes"`
	Errors   []*ErrorObj              `json:"errors"`
}
