package twitter

import (
	"fmt"
	"strings"
)

type endpoint string

const (
	tweetLookupEndpoint         endpoint = "2/tweets"
	tweetCreateEndpoint         endpoint = "2/tweets"
	tweetDeleteEndpoint         endpoint = "2/tweets/{id}"
	userLookupEndpoint          endpoint = "2/users"
	userNameLookupEndpoint      endpoint = "2/users/by"
	userAuthLookupEndpoint      endpoint = "2/users/me"
	tweetRecentSearchEndpoint   endpoint = "2/tweets/search/recent"
	tweetRecentCountsEndpoint   endpoint = "2/tweets/counts/recent"
	userFollowingEndpoint       endpoint = "2/users/{id}/following"
	userFollowersEndpoint       endpoint = "2/users/{id}/followers"
	userTweetTimelineEndpoint   endpoint = "2/users/{id}/tweets"
	userMentionTimelineEndpoint endpoint = "2/users/{id}/mentions"
	tweetHideRepliesEndpoint    endpoint = "2/tweets/{id}/hidden"

	idTag = "{id}"
)

func (e endpoint) url(host string) string {
	return fmt.Sprintf("%s/%s", host, string(e))
}

func (e endpoint) urlID(host, id string) string {
	u := fmt.Sprintf("%s/%s", host, string(e))
	return strings.ReplaceAll(u, idTag, id)
}
