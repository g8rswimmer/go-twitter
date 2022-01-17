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
	userManageRetweetEndpoint   endpoint = "2/users/{id}/retweets"
	userBlocksEndpoint          endpoint = "2/users/{id}/blocking"
	userMutesEndpont            endpoint = "2/users/{id}/muting"
	userRetweetLookupEndpoint   endpoint = "2/tweets/{id}/retweeted_by"
	tweetRecentSearchEndpoint   endpoint = "2/tweets/search/recent"
	tweetRecentCountsEndpoint   endpoint = "2/tweets/counts/recent"
	userFollowingEndpoint       endpoint = "2/users/{id}/following"
	userFollowersEndpoint       endpoint = "2/users/{id}/followers"
	userTweetTimelineEndpoint   endpoint = "2/users/{id}/tweets"
	userMentionTimelineEndpoint endpoint = "2/users/{id}/mentions"
	tweetHideRepliesEndpoint    endpoint = "2/tweets/{id}/hidden"
	userTweetLikesEndpoint      endpoint = "2/tweets/{id}/liking_users"
	tweetUserLikesEndpoint      endpoint = "2/users/{id}/liked_tweets"
	userLikesEndpoint           endpoint = "2/users/{id}/likes"

	idTag = "{id}"
)

func (e endpoint) url(host string) string {
	return fmt.Sprintf("%s/%s", host, string(e))
}

func (e endpoint) urlID(host, id string) string {
	u := fmt.Sprintf("%s/%s", host, string(e))
	return strings.ReplaceAll(u, idTag, id)
}
