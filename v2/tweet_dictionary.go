package twitter

// TweetLookupResponse contains all of the information from a tweet lookup callout
type TweetLookupResponse struct {
	Raw *TweetRaw
}

// TweetRaw is the raw response from the tweet lookup endpoint
type TweetRaw struct {
	Tweets   []*TweetObj       `json:"data"`
	Includes *TweetRawIncludes `json:"includes,omitempty"`
	Errors   []*ErrorObj       `json:"errors,omitempty"`
}

// TweetRawIncludes contains any additional information from the tweet callout
type TweetRawIncludes struct {
	Tweets []*TweetObj `json:"tweets,omitempty"`
	Users  []*UserObj  `json:"users,omitempty"`
	Places []*PlaceObj `json:"places,omitempty"`
	Media  []*MediaObj `json:"media,omitempty"`
	Polls  []*PollObj  `json:"polls,omitempty"`
}

type tweetraw struct {
	Tweet    *TweetObj         `json:"data"`
	Includes *TweetRawIncludes `json:"includes"`
	Errors   []*ErrorObj       `json:"errors"`
}
