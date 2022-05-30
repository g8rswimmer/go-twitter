# Twitter v2 Tweets Examples
This directory contains examples for the APIs under `Tweets` in the Developer Platform.

## Examples
The examples can be run my providing some options, including the authorization token.

### [Tweets Lookup](https://developer.twitter.com/en/docs/twitter-api/tweets/lookup/introduction)

* [Returns a variety of information about the tweet specified by the list of ids](./lookup/tweet-lookup/main.go)

### [Manage Tweets](https://developer.twitter.com/en/docs/twitter-api/tweets/manage-tweets/introduction)

* [Creates a tweet on behalf of the user](./manage/tweet-create/main.go)
* [Allows a user to delete a tweet](./manage/tweet-delete/main.go)

### [Timelines](https://developer.twitter.com/en/docs/twitter-api/tweets/timelines/introduction)

* [Returns most recent tweets composed by a user](./timeline/user-tweet-timeline/main.go)
* [Returns most recent tweet mentioning a user](./timeline/user-mention-timeline/main.go)
* [Allows you to retrieve a collection of the most recent Tweets and Retweets posted by you and users you follow.](./timeline/user-tweet-reverse-chronological-timeline/main.go)

### [Search Tweets](https://developer.twitter.com/en/docs/twitter-api/tweets/search/introduction)

* [Search for Tweets published in the last 7 days](./search/tweet-recent-search/main.go)
* [Full-archive search endpoint returns the complete history of public Tweets](./search/tweet-search-all/main.go)
    * This endpoint is only available to those users who have been approved for Academic Research access

### [Tweet Counts](https://developer.twitter.com/en/docs/twitter-api/tweets/counts/introduction)

* [Receive a count of Tweets that match a query in the last 7 days](./counts/tweet-recent-counts/main.go)
* [Full-archive Tweet counts endpoint returns the count of Tweets that match your query from the complete history of public Tweets](./counts/tweet-all-counts/main.go)
    * This endpoint is only available to those users who have been approved for Academic Research access

### [Filtered Stream](https://developer.twitter.com/en/docs/twitter-api/tweets/filtered-stream/introduction)

* [Add rules from your stream](./filtered-stream/tweet-search-stream-add-rule/main.go)
* [Delete rules from your stream](./filtered-stream/tweet-search-stream-delete-rules/main.go)
* [Retrieve your stream's rules](./filtered-stream/tweet-search-stream-rules/main.go)
* [Connect to the stream](./filtered-stream/tweet-search-stream/main.go)

### [Volume Streams](https://developer.twitter.com/en/docs/twitter-api/tweets/volume-streams/introduction)

* [Streams about 1% of all Tweets in real-time](./volume-stream/tweet-sample-stream/main.go)

### [Retweets](https://developer.twitter.com/en/docs/twitter-api/tweets/retweets/introduction)

* [Users who have Retweeted a Tweet](./retweets/user-retweet-lookup/main.go)
* [Allows a user ID to Retweet a Tweet](./retweets/user-retweet/main.go)
* [Allows a user ID to undo a Retweet](./retweets/user-delete-retweet/main.go)

### [Likes](https://developer.twitter.com/en/docs/twitter-api/tweets/likes/introduction)

* [Users who have liked a Tweet](./likes/user-likes-lookup/main.go)
* [Tweets liked by an user](./likes/tweet-likes-lookup/main.go)
* [Allows a user ID to like a Tweet](./likes/user-like-tweet/main.go)
* [Allows a user ID to unlike a Tweet](./likes/user-unlike-tweet/main.go)

### [Hide Replies](https://developer.twitter.com/en/docs/twitter-api/tweets/hide-replies/introduction)

* [Hides or unhides a reply to a Tweet](./hide-replies/tweet-hide-replies/main.go)

### [Quote Tweets](https://developer.twitter.com/en/docs/twitter-api/tweets/quote-tweets/introduction)

* [Returns Quote Tweets for a Tweet specified by the requested Tweet ID](./quote/quote-tweets/main.go)

### [Bookmarks](https://developer.twitter.com/en/docs/twitter-api/tweets/bookmarks/introduction)

* [Allows you to get an authenticated user's 800 most recent bookmarked Tweets](./bookmarks/tweet-bookmarks-lookup/main.go)
* [Causes the user ID identified in the path parameter to Bookmark the target Tweet provided in the request body](./bookmarks/tweet-add-bookmark/main.go)
* [Allows a user or authenticated user ID to remove a Bookmark of a Tweet](./bookmarks/tweet-remove-bookmark/main.go)