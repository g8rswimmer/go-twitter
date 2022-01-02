package twitter

// RetweetData will be returned by the manage retweet APIs
type RetweetData struct {
	Retweeted bool `json:"retweeted"`
}

// UserRetweetResponse is the response with a user retweet
type UserRetweetResponse struct {
	Data *RetweetData `json:"data"`
}

// DeleteUserRetweetResponse is the response with a uset retweet
type DeleteUserRetweetResponse struct {
	Data *RetweetData `json:"data"`
}
