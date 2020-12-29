package twitter

type UserLookupResponse struct {
	Raw *UserRaw
}

type userraw struct {
	User     *UserObj         `json:"data"`
	Includes *UserRawIncludes `json:"includes"`
	Errors   []*ErrorObj      `json:"errors"`
}

type UserRaw struct {
	Users    []*UserObj       `json:"data"`
	Includes *UserRawIncludes `json:"includes,omitempty"`
	Errors   []*ErrorObj      `json:"errors,omitempty"`
}

type UserRawIncludes struct {
	Tweets []*TweetObj `json:"tweets,omitempty"`
}
