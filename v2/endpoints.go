package twitter

import (
	"fmt"
	"strings"
)

type endpoint string

const (
	tweetLookupEndpoint                           endpoint = "2/tweets"
	tweetCreateEndpoint                           endpoint = "2/tweets"
	tweetDeleteEndpoint                           endpoint = "2/tweets/{id}"
	userLookupEndpoint                            endpoint = "2/users"
	userNameLookupEndpoint                        endpoint = "2/users/by"
	userAuthLookupEndpoint                        endpoint = "2/users/me"
	userManageRetweetEndpoint                     endpoint = "2/users/{id}/retweets"
	userBlocksEndpoint                            endpoint = "2/users/{id}/blocking"
	userMutesEndpoint                             endpoint = "2/users/{id}/muting"
	userRetweetLookupEndpoint                     endpoint = "2/tweets/{id}/retweeted_by"
	tweetRecentSearchEndpoint                     endpoint = "2/tweets/search/recent"
	tweetSearchEndpoint                           endpoint = "2/tweets/search/all"
	tweetRecentCountsEndpoint                     endpoint = "2/tweets/counts/recent"
	tweetAllCountsEndpoint                        endpoint = "2/tweets/counts/all"
	userFollowingEndpoint                         endpoint = "2/users/{id}/following"
	userFollowersEndpoint                         endpoint = "2/users/{id}/followers"
	userTweetTimelineEndpoint                     endpoint = "2/users/{id}/tweets"
	userMentionTimelineEndpoint                   endpoint = "2/users/{id}/mentions"
	userTweetReverseChronologicalTimelineEndpoint endpoint = "2/users/{id}/timelines/reverse_chronological"
	tweetHideRepliesEndpoint                      endpoint = "2/tweets/{id}/hidden"
	tweetLikesEndpoint                            endpoint = "2/tweets/{id}/liking_users"
	userLikedTweetEndpoint                        endpoint = "2/users/{id}/liked_tweets"
	userLikesEndpoint                             endpoint = "2/users/{id}/likes"
	tweetSampleStreamEndpoint                     endpoint = "2/tweets/sample/stream"
	tweetSearchStreamRulesEndpoint                endpoint = "2/tweets/search/stream/rules"
	tweetSearchStreamEndpoint                     endpoint = "2/tweets/search/stream"
	listLookupEndpoint                            endpoint = "2/lists/{id}"
	userListLookupEndpoint                        endpoint = "2/users/{id}/owned_lists"
	listTweetLookupEndpoint                       endpoint = "2/lists/{id}/tweets"
	listCreateEndpoint                            endpoint = "2/lists"
	listUpdateEndpoint                            endpoint = "2/lists/{id}"
	listDeleteEndpoint                            endpoint = "2/lists/{id}"
	listMemberEndpoint                            endpoint = "2/lists/{id}/members"
	userListMemberEndpoint                        endpoint = "2/users/{id}/list_memberships"
	userPinnedListEndpoint                        endpoint = "2/users/{id}/pinned_lists"
	userFollowedListEndpoint                      endpoint = "2/users/{id}/followed_lists"
	listUserFollowersEndpoint                     endpoint = "2/lists/{id}/followers"
	spaceLookupEndpoint                           endpoint = "2/spaces"
	spaceByCreatorLookupEndpoint                  endpoint = "2/spaces/by/creator_ids"
	spaceBuyersLookupEndpoint                     endpoint = "2/spaces/{id}/buyers"
	spaceTweetsLookupEndpoint                     endpoint = "2/spaces/{id}/tweets"
	spaceSearchEndpoint                           endpoint = "2/spaces/search"
	complianceJobsEndpoint                        endpoint = "2/compliance/jobs"
	quoteTweetLookupEndpoint                      endpoint = "2/tweets/{id}/quote_tweets"
	tweetBookmarksEndpoint                        endpoint = "2/users/{id}/bookmarks"

	idTag = "{id}"
)

func (e endpoint) url(host string) string {
	return fmt.Sprintf("%s/%s", host, string(e))
}

func (e endpoint) urlID(host, id string) string {
	u := fmt.Sprintf("%s/%s", host, string(e))
	return strings.ReplaceAll(u, idTag, id)
}
