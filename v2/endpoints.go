package twitter

import (
	"fmt"
	"strings"
)

type endpoint string

const (
	tweetLookupEndpoint         endpoint = "2/tweets"
	userLookupEndpoint          endpoint = "2/users"
	userNameLookupEndpoint      endpoint = "2/users/by"
	tweetRecentSearchEndpoint   endpoint = "2/tweets/search/recent"
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
